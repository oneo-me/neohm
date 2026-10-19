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
	"github.com/oneo-me/onevra/source/components"
)

// How often the climate entity is read while the panel shows, and while it
// is closed and only the menu bar title needs it.
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
	panel, settings *mygo.Window
	settingsState   components.SettingsState

	climate    ha.Climate
	hasClimate bool
	unit       string
	online     bool
	errorKey   string
	errorText  string
	entities   []ha.Climate
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
	if a.tray == nil || !a.configured() {
		a.openSettings()
	}
}

func (a *application) configured() bool {
	return a.config.Address != "" && a.token != ""
}

// connect points the poller at the saved address and token, and with list
// reads the climate entities for Settings.
func (a *application) connect(list bool) {
	a.mu.Lock()
	if a.configured() {
		a.client = ha.NewClient(a.config.Address, a.token)
	} else {
		a.client = nil
	}
	a.entity = a.config.Entity
	client := a.client
	a.mu.Unlock()
	a.hasClimate = a.hasClimate && a.climate.EntityID == a.config.Entity
	a.wake()
	if !list || client == nil {
		return
	}
	a.connecting, a.connectMsg = true, ""
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		list, err := client.Climates(ctx)
		a.post(func() {
			a.connecting = false
			if err != nil {
				a.connectMsg = a.errorMessage(err)
			} else {
				a.entities = list
				a.connectMsg = locales.ConnectedTo(len(list))
				a.chooseEntity()
			}
			a.redraw()
		})
	}()
}

// chooseEntity keeps the saved entity if it still exists, else picks the
// first climate entity.
func (a *application) chooseEntity() {
	if slices.ContainsFunc(a.entities, func(c ha.Climate) bool { return c.EntityID == a.config.Entity }) || len(a.entities) == 0 {
		return
	}
	a.selectEntity(a.entities[0].EntityID)
}

func (a *application) selectEntity(id string) {
	a.config.Entity = id
	a.save()
	a.hasClimate = false
	a.connect(false)
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

// poll reads the climate entity on a background goroutine and hands each
// result to the main thread.
func (a *application) poll() {
	var unitFor *ha.Client
	for {
		a.mu.Lock()
		client, entity := a.client, a.entity
		a.mu.Unlock()
		if client != nil && entity != "" {
			edits := a.edits.Load()
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			c, err := client.Climate(ctx, entity)
			unit := ""
			if err == nil && unitFor != client {
				if u, uerr := client.TemperatureUnit(ctx); uerr == nil && u != "" {
					unit, unitFor = u, client
				}
			}
			cancel()
			a.post(func() { a.applyState(client, edits, c, unit, err) })
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

func (a *application) applyState(client *ha.Client, edits int64, c ha.Climate, unit string, err error) {
	a.mu.Lock()
	current := a.client == client && (err != nil || a.entity == c.EntityID)
	a.mu.Unlock()
	if !current {
		return
	}
	if err != nil {
		a.online = false
		a.report("requestError", err)
		if errors.Is(err, ha.ErrUnauthorized) {
			a.errorKey = "unauthorized"
		}
		a.redraw()
		return
	}
	a.online = true
	if a.errorKey == "requestError" || a.errorKey == "unauthorized" {
		a.clearError()
	}
	if unit != "" {
		a.unit = unit
	}
	// A local change made after the read started wins until Home Assistant
	// reports it.
	if !a.hasClimate || a.edits.Load() == edits && a.pending == 0 && !a.tempPending && !time.Now().Before(a.holdUntil) {
		a.climate, a.hasClimate = c, true
	}
	a.redraw()
}

func (a *application) errorMessage(err error) string {
	if errors.Is(err, ha.ErrUnauthorized) {
		return locales.Unauthorized()
	}
	return locales.RequestError() + " " + err.Error()
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
				a.tempPending = false
				a.call("set_temperature", map[string]any{"temperature": a.climate.Target})
			}
		})
	})
}

// redraw refreshes the windows and the menu bar title after a change.
func (a *application) redraw() {
	for _, w := range []*mygo.Window{a.panel, a.settings} {
		if w != nil {
			w.Invalidate()
		}
	}
	if a.tray != nil {
		a.tray.SetTitle(a.trayTitle())
	}
}

// trayTitle is the target temperature shown beside the menu bar icon.
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
	c.OnShortcut(ui.Cmd, ui.KeyComma, a.openSettings)
}
