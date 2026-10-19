package main

import (
	"context"
	"errors"
	"log"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/ha-go/source/internal/ha"
	"github.com/oneo-me/ha-go/source/locales"
	"github.com/oneo-me/onevra/source/shells/plain"
)

// How often Home Assistant is read while the panel shows, and while it is
// closed and only the menu bar title needs the device opened last.
const (
	pollOpen   = 3 * time.Second
	pollClosed = 20 * time.Second
	// debounce is how long the temperature waits for more slider input
	// before Home Assistant hears it.
	debounce = 450 * time.Millisecond
	// settle is how long after a service call the panel trusts what it
	// shows over what Home Assistant reports.
	settle = 1500 * time.Millisecond
)

type secretStore interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

// application keeps the state of the menu bar app. Everything but the fields
// under mu and the atomics is owned by the main thread.
type application struct {
	path    string
	config  config
	secrets secretStore
	// address and token are the values being edited in Settings.
	address, token string

	ready, quitting bool
	tray            *mygo.Tray
	trayMenu        *mygo.Menu
	// panel is the open menu bar panel; window shows its pages instead when
	// there is no menu bar icon to anchor it to.
	panel  *mygo.Window
	window *plain.Shell
	// The panel opens on the device list: detail shows the device opened in
	// its place, and settingsOpen shows Settings over either.
	detail, settingsOpen bool

	// devices are the climate entities the list shows, as last read, which
	// listed reports.
	devices []ha.Climate
	listed  bool
	// climate is the device opened last, config.Entity, once hasClimate says
	// it is known.
	climate    ha.Climate
	hasClimate bool
	unit       string
	online     bool
	errorKey   string
	errorText  string
	connecting bool
	connectMsg string
	tempGen    int
	// pending counts service calls on their way, and tempPending a
	// temperature waiting for the slider to rest; holdUntil lets the unit
	// report a change. Reads do not replace the state shown meanwhile.
	pending     int
	tempPending bool
	holdUntil   time.Time

	mu     sync.Mutex
	client *ha.Client // nil until configured
	entity string
	// listing has the poller read every device instead of entity alone.
	listing bool

	// edits counts local changes, so that a read started before one does
	// not undo it on screen.
	edits     atomic.Int64
	panelOpen atomic.Bool
	kick      chan struct{}
	// post runs a function on the main thread; tests replace it.
	post func(func())
}

func newApplication(path string, secrets secretStore) *application {
	return &application{path: path, secrets: secrets, kick: make(chan struct{}, 1), unit: "°C", post: mygo.RunOnMain}
}

func (a *application) report(key string, err error) {
	a.errorKey, a.errorText = key, err.Error()
	log.Printf("%s: %v", key, err)
}

func (a *application) clearError() { a.errorKey, a.errorText = "", "" }

func (a *application) start() {
	if err := a.createTray(); err != nil {
		a.report("trayError", err)
	}
	a.installMenus()
	a.ready = true
	a.connect(false)
	go a.poll()
	switch {
	case !a.configured():
		a.openSettings()
	case a.tray == nil:
		a.showPanel()
	}
}

func (a *application) configured() bool {
	return a.config.Address != "" && a.token != ""
}

// connect points the poller at the saved address and token, and with list
// reads the devices at once to report the connection in Settings.
func (a *application) connect(list bool) {
	a.mu.Lock()
	if a.configured() {
		a.client = ha.NewClient(a.config.Address, a.token)
	} else {
		a.client = nil
	}
	client := a.client
	a.mu.Unlock()
	// What the old connection read may not exist on this one.
	a.devices, a.listed, a.hasClimate, a.detail = nil, false, false, false
	a.follow()
	if !list || client == nil {
		return
	}
	a.connecting, a.connectMsg = true, ""
	edits := a.edits.Load()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		list, err := client.Climates(ctx)
		a.post(func() {
			if !a.current(client) {
				return
			}
			a.connecting = false
			if err == nil {
				a.connectMsg = locales.ConnectedTo(len(list))
			}
			a.applyList(client, edits, list, "", err)
		})
	}()
}

// follow points the poller at what the panel shows: every device while the
// list is open, else the device opened last, which the menu bar title
// follows.
func (a *application) follow() {
	a.mu.Lock()
	a.entity = a.config.Entity
	a.listing = a.panelOpen.Load() && !a.detail
	a.mu.Unlock()
	a.wake()
}

// current reports whether client is the connection in use; only its reads
// are shown.
func (a *application) current(client *ha.Client) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.client == client
}

// openDevice shows the controls of device d in place of the list. The menu
// bar title follows the device opened last.
func (a *application) openDevice(d ha.Climate) {
	if d.EntityID != a.config.Entity {
		// A temperature still waiting for the slider goes to its own device.
		a.sendTarget()
		a.config.Entity = d.EntityID
		a.save()
	}
	// The device opened last may show changes the list's read lacks.
	if !a.hasClimate || a.climate.EntityID != d.EntityID {
		a.climate, a.hasClimate = d, true
	}
	a.detail = true
	a.follow()
	a.redraw()
}

// closeDevice returns the panel to the device list.
func (a *application) closeDevice() {
	a.detail = false
	a.follow()
	a.redraw()
}

func (a *application) save() {
	if err := saveConfig(a.path, a.config); err != nil {
		a.report("saveError", err)
	} else if a.errorKey == "saveError" || a.errorKey == "loadError" {
		a.clearError()
	}
}

// saveConnection applies the address and token edited in Settings.
func (a *application) saveConnection() {
	a.config.Address = ha.NormalizeURL(a.address)
	a.address = a.config.Address
	if err := a.secrets.Set(tokenKey, []byte(a.token)); err != nil {
		a.report("secretError", err)
	}
	a.save()
	a.connect(true)
	a.redraw()
}

func (a *application) wake() {
	select {
	case a.kick <- struct{}{}:
	default:
	}
}

// poll reads Home Assistant on a background goroutine and hands each result
// to the main thread.
func (a *application) poll() {
	var unitFor *ha.Client
	for {
		a.mu.Lock()
		client, entity, listing := a.client, a.entity, a.listing
		a.mu.Unlock()
		if client != nil && (listing || entity != "") {
			edits := a.edits.Load()
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			var (
				list []ha.Climate
				c    ha.Climate
				err  error
			)
			if listing {
				list, err = client.Climates(ctx)
			} else {
				c, err = client.Climate(ctx, entity)
			}
			unit := ""
			if err == nil && unitFor != client {
				if u, uerr := client.TemperatureUnit(ctx); uerr == nil && u != "" {
					unit, unitFor = u, client
				}
			}
			cancel()
			if listing {
				a.post(func() { a.applyList(client, edits, list, unit, err) })
			} else {
				a.post(func() { a.applyState(client, edits, c, unit, err) })
			}
		}
		wait := pollClosed
		if a.panelOpen.Load() {
			wait = pollOpen
		}
		select {
		case <-a.kick:
		case <-time.After(wait):
		}
	}
}

// applyList shows a read of every device. The device opened last takes its
// entry as applyState takes a read of it; once it is gone, the panel and the
// menu bar title forget it.
func (a *application) applyList(client *ha.Client, edits int64, list []ha.Climate, unit string, err error) {
	if !a.current(client) || !a.reached(unit, err) {
		return
	}
	a.devices, a.listed = list, true
	if i := slices.IndexFunc(list, func(c ha.Climate) bool { return c.EntityID == a.config.Entity }); i >= 0 {
		a.update(edits, list[i])
	} else if a.config.Entity != "" {
		a.config.Entity, a.hasClimate, a.detail = "", false, false
		a.save()
		a.follow()
	}
	a.redraw()
}

func (a *application) applyState(client *ha.Client, edits int64, c ha.Climate, unit string, err error) {
	a.mu.Lock()
	current := a.client == client && (err != nil || a.entity == c.EntityID)
	a.mu.Unlock()
	if !current || !a.reached(unit, err) {
		return
	}
	a.update(edits, c)
	a.redraw()
}

// reached records whether Home Assistant answered a read, with the unit it
// reported, and reports the failure when it did not.
func (a *application) reached(unit string, err error) bool {
	if err != nil {
		a.online = false
		a.report("requestError", err)
		if errors.Is(err, ha.ErrUnauthorized) {
			a.errorKey = "unauthorized"
		}
		a.redraw()
		return false
	}
	a.online = true
	if a.errorKey == "requestError" || a.errorKey == "unauthorized" {
		a.clearError()
	}
	if unit != "" {
		a.unit = unit
	}
	return true
}

// update takes a read of the device opened last, unless a local change made
// after the read started has yet to reach Home Assistant: that wins until
// Home Assistant reports it.
func (a *application) update(edits int64, c ha.Climate) {
	if !a.hasClimate || a.edits.Load() == edits && a.pending == 0 && !a.tempPending && !time.Now().Before(a.holdUntil) {
		a.climate, a.hasClimate = c, true
	}
}

// shown is what the list shows of device d: the device opened last as the
// panel controls it, with changes Home Assistant may not report yet, else the
// list's read.
func (a *application) shown(d ha.Climate) ha.Climate {
	if a.hasClimate && a.climate.EntityID == d.EntityID {
		return a.climate
	}
	return d
}

// call sends a climate service call. The caller has already changed the
// state the panel shows.
func (a *application) call(service string, data map[string]any) {
	a.edits.Add(1)
	a.mu.Lock()
	client, entity := a.client, a.entity
	a.mu.Unlock()
	a.redraw()
	if client == nil || entity == "" {
		return
	}
	a.pending++
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := client.Call(ctx, service, entity, data)
		a.post(func() {
			a.pending--
			a.edits.Add(1)
			if err != nil {
				a.report("requestError", err)
				if errors.Is(err, ha.ErrUnauthorized) {
					a.errorKey = "unauthorized"
				}
				// Show what the unit really does.
				a.holdUntil = time.Time{}
				a.wake()
			} else {
				// Read again once the unit had a moment to report the change.
				a.holdUntil = time.Now().Add(settle)
				time.AfterFunc(settle, a.wake)
			}
			a.redraw()
		})
	}()
}

func (a *application) togglePower() {
	if a.climate.On() {
		a.climate.Mode = "off"
		a.call("turn_off", nil)
		return
	}
	mode := a.climate.OnMode()
	if mode == "" {
		return
	}
	a.climate.Mode = mode
	a.call("set_hvac_mode", map[string]any{"hvac_mode": mode})
}

func (a *application) setMode(mode string) {
	a.climate.Mode = mode
	a.call("set_hvac_mode", map[string]any{"hvac_mode": mode})
}

func (a *application) setFan(mode string) {
	a.climate.FanMode = mode
	a.call("set_fan_mode", map[string]any{"fan_mode": mode})
}

func (a *application) setSwing(mode string) {
	a.climate.SwingMode = mode
	a.call("set_swing_mode", map[string]any{"swing_mode": mode})
}

// togglePreset switches preset p, as "eco" or "sleep", on or back to the
// neutral preset.
func (a *application) togglePreset(p string) {
	next := p
	if a.climate.Preset == p {
		next = a.climate.NeutralPreset()
	}
	if next == "" {
		return
	}
	a.climate.Preset = next
	a.call("set_preset_mode", map[string]any{"preset_mode": next})
}

// setTarget shows t at once and sends it once the slider rests.
func (a *application) setTarget(t float64) {
	a.climate.Target = a.climate.Clamp(t)
	a.climate.HasTarget = true
	a.edits.Add(1)
	a.tempGen++
	a.tempPending = true
	gen := a.tempGen
	a.redraw()
	time.AfterFunc(debounce, func() {
		a.post(func() {
			if gen == a.tempGen {
				a.sendTarget()
			}
		})
	})
}

// sendTarget sends the temperature waiting for the slider to rest, if any.
func (a *application) sendTarget() {
	if !a.tempPending {
		return
	}
	a.tempPending = false
	a.call("set_temperature", map[string]any{"temperature": a.climate.Target})
}

// redraw refreshes the panel and the menu bar title after a change.
func (a *application) redraw() {
	if a.panel != nil {
		a.panel.Invalidate()
	}
	if a.window != nil {
		a.window.Window().Invalidate()
	}
	if a.tray != nil {
		a.tray.SetTitle(a.trayTitle())
	}
}

// trayTitle is the target temperature of the device opened last, shown
// beside the menu bar icon.
func (a *application) trayTitle() string {
	if !a.hasClimate || !a.climate.On() || !a.climate.HasTarget {
		return ""
	}
	return formatTemp(a.climate.Target) + "°"
}

func formatTemp(t float64) string { return strconv.FormatFloat(t, 'f', -1, 64) }

// optionText names a Home Assistant option, as "fan_only", in the app's
// language; options it has no message for show as they are.
func optionText(kind, value string) string {
	key := kind + "_" + value
	if _, _, ok := locales.Catalog.Message(i18n.Locale(), key); ok {
		return locales.Catalog.Format(i18n.Locale(), key, nil)
	}
	if value == "" {
		return "—"
	}
	return value
}

func (a *application) errorDisplay() string {
	switch a.errorKey {
	case "":
		return ""
	case "unauthorized":
		return locales.Unauthorized()
	}
	return locales.Catalog.Format(i18n.Locale(), a.errorKey, nil) + " " + a.errorText
}

func (a *application) scope(c *ui.Context) {
	c.OnShortcut(ui.Cmd, ui.KeyQ, mygo.App.Quit)
	c.OnShortcut(ui.Cmd, ui.KeyComma, a.showSettings)
}
