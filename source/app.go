package main

import (
	"context"
	"errors"
	"log"
	"maps"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
	"github.com/oneo-me/neohm/source/internal/ha"
	"github.com/oneo-me/neohm/source/locales"
	"github.com/oneo-me/onevra/source/shells/plain"
)

// How often Home Assistant is read while the panel shows. Nothing is read
// while it is closed.
const (
	pollOpen = 3 * time.Second
	// debounce is how long the temperature waits for more slider input
	// before Home Assistant hears it.
	debounce = 450 * time.Millisecond
	// settle is how long after a service call the panel trusts what it
	// shows over what Home Assistant reports.
	settle = 1500 * time.Millisecond
	// Home Assistant's registries change rarely: they are read again after
	// registryAge, or registryRetry after a failure.
	registryAge   = 5 * time.Minute
	registryRetry = time.Minute
)

// application keeps the state of the menu bar app. Everything but the fields
// under mu and the atomics is owned by the main thread.
type application struct {
	path   string
	config config
	// address and token are the values being edited in Settings.
	address, token string

	ready, quitting bool
	tray            *mygo.Tray
	trayMenu        *mygo.Menu
	// panel is the open menu bar panel; window shows its pages instead when
	// there is no menu bar icon to anchor it to.
	panel  *mygo.Window
	window *plain.Shell
	// The panel opens on Home: detail shows the air conditioner opened last
	// in its place, device the page of another device, by its entity, and
	// settingsOpen shows Settings over any of them.
	detail, settingsOpen bool
	device               string

	// home is what the home page shows, as last read, which listed
	// reports. registry is what Home Assistant's registries said, nil
	// until they could be read.
	home     []ha.Section
	listed   bool
	registry *ha.Registry
	// holds keep what tiles show after an action until Home Assistant
	// reports it.
	holds map[string]*hold
	// climate is the device opened last, config.Entity, once hasClimate says
	// it is known.
	climate    ha.Climate
	hasClimate bool
	unit       string
	online     bool
	errorKey   string
	errorText  string
	connecting bool
	// connected reports that the last Connect reached Home Assistant, which
	// Settings says until it opens again.
	connected bool
	tempGen   int
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
	// The registries were read for regFor, and are read again after
	// regNext.
	regFor  *ha.Client
	regNext time.Time

	// edits counts local changes, so that a read started before one does
	// not undo it on screen.
	edits     atomic.Int64
	panelOpen atomic.Bool
	kick      chan struct{}
	// upd is what Settings' updates found, with updates.
	upd     updateState
	updates updates
	// atLogin is whether the system opens the app at login, as Settings
	// last read it, and loginFailed that Settings could not change it.
	atLogin, loginFailed bool
	// post runs a function on the main thread, setLocale shows a language,
	// quit quits the app, and openAtLogin and setOpenAtLogin read and change
	// whether the system opens it at login; tests replace them.
	post           func(func())
	setLocale      func(string)
	quit           func()
	openAtLogin    func() bool
	setOpenAtLogin func(bool) error
}

func newApplication(path string) *application {
	return &application{path: path, kick: make(chan struct{}, 1), unit: "°C", post: mygo.RunOnMain, setLocale: mygo.App.SetLocale, quit: mygo.App.Quit,
		openAtLogin: mygo.App.OpenAtLogin, setOpenAtLogin: mygo.App.SetOpenAtLogin,
		updates: mygoUpdates, upd: updateState{notes: &markdown.Options{Theme: &markdown.Theme{FontSize: 12}}}}
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
	go a.watchUpdates()
	switch {
	case !a.configured():
		a.openSettings()
	case a.tray == nil:
		a.showPanel()
	}
}

func (a *application) configured() bool {
	return a.config.Address != "" && a.config.Token != ""
}

// connect points the poller at the saved address and token, and with list
// reads the devices at once to report the connection in Settings.
func (a *application) connect(list bool) {
	a.mu.Lock()
	if a.configured() {
		a.client = ha.NewClient(a.config.Address, a.config.Token)
	} else {
		a.client = nil
	}
	client := a.client
	a.mu.Unlock()
	// What the old connection read may not exist on this one.
	a.home, a.listed, a.registry, a.holds, a.hasClimate, a.detail, a.device = nil, false, nil, nil, false, false, ""
	a.follow()
	if !list || client == nil {
		return
	}
	a.connecting, a.connected = true, false
	edits := a.edits.Load()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// The registries too, for Settings to count the tiles Home shows.
		states, reg, err := a.readHome(ctx, client, true)
		a.post(func() {
			if !a.current(client) {
				return
			}
			a.connecting, a.connected = false, err == nil
			a.applyList(client, edits, states, reg, "", err)
		})
	}()
}

// follow points the poller at what the panel shows: every device while the
// list is open, else the air conditioner opened last.
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

// openClimate shows the controls of air conditioner d in place of the home
// page.
func (a *application) openClimate(d ha.Climate) {
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
	a.detail, a.device = true, ""
	a.follow()
	a.redraw()
}

// openDevice shows the page of the device of entity id in place of the home
// page.
func (a *application) openDevice(id string) {
	a.detail, a.device = false, id
	a.follow()
	a.redraw()
}

// closeDevice returns the panel to the home page.
func (a *application) closeDevice() {
	a.detail, a.device = false, ""
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

// setLanguage shows the app in the language of tag, one of the catalogs',
// or follows the system's with "".
func (a *application) setLanguage(tag string) {
	a.config.Language = tag
	a.save()
	a.setLocale(tag)
}

// saveConnection applies the address and token edited in Settings.
func (a *application) saveConnection() {
	a.config.Address, a.config.Token = ha.NormalizeURL(a.address), a.token
	a.address = a.config.Address
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
		open := a.panelOpen.Load()
		if client != nil && open && (listing || entity != "") {
			edits := a.edits.Load()
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			var (
				states []ha.State
				reg    *ha.Registry
				s      ha.State
				err    error
			)
			if listing {
				states, reg, err = a.readHome(ctx, client, false)
			} else {
				s, err = client.State(ctx, entity)
			}
			unit := ""
			if err == nil && unitFor != client {
				if u, uerr := client.TemperatureUnit(ctx); uerr == nil && u != "" {
					unit, unitFor = u, client
				}
			}
			cancel()
			if listing {
				a.post(func() { a.applyList(client, edits, states, reg, unit, err) })
			} else {
				a.post(func() { a.applyState(client, edits, ha.ParseClimate(s), unit, err) })
			}
		}
		if !open {
			<-a.kick
			continue
		}
		select {
		case <-a.kick:
		case <-time.After(pollOpen):
		}
	}
}

// readHome reads every entity, and the registries when always says so or
// this connection has not read them lately; reg is nil when it did not. The
// registries only arrange the home page, so failing to read them is no
// failure.
func (a *application) readHome(ctx context.Context, client *ha.Client, always bool) (states []ha.State, reg *ha.Registry, err error) {
	if states, err = client.States(ctx); err != nil {
		return nil, nil, err
	}
	a.mu.Lock()
	due := always || a.regFor != client || !time.Now().Before(a.regNext)
	if due {
		a.regFor, a.regNext = client, time.Now().Add(registryRetry)
	}
	a.mu.Unlock()
	if !due {
		return states, nil, nil
	}
	r, err := client.Registry(ctx)
	if err != nil {
		log.Printf("registry: %v", err)
		return states, nil, nil
	}
	a.mu.Lock()
	if a.regFor == client {
		a.regNext = time.Now().Add(registryAge)
	}
	a.mu.Unlock()
	return states, &r, nil
}

// applyList shows a read of every entity, arranged with the registries last
// read. The air conditioner opened last takes its state as applyState takes
// a read of it; once it is gone, the panel forgets it.
// The page of a device gone from Home closes.
func (a *application) applyList(client *ha.Client, edits int64, states []ha.State, reg *ha.Registry, unit string, err error) {
	if !a.current(client) || !a.reached(unit, err) {
		return
	}
	if reg != nil {
		a.registry = reg
	}
	a.home, a.listed = ha.Arrange(states, a.registry), true
	if _, _, ok := a.tileOf(a.device); !ok {
		a.device = ""
	}
	a.release(edits)
	if i := slices.IndexFunc(states, func(s ha.State) bool { return s.EntityID == a.config.Entity }); i >= 0 {
		a.update(edits, ha.ParseClimate(states[i]))
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

// shown is what the home page shows of air conditioner d: the one opened
// last as the panel controls it, with changes Home Assistant may not report
// yet, else the read.
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
		err := client.Call(ctx, "climate", service, entity, data)
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

// hold is what a device shows after an action, until a read that started
// once the action's calls returned and settled.
type hold struct {
	// state is the state shown, unless "", attrs attributes shown over
	// those read, as a light's color, and level the level its page shows,
	// as a light's brightness, when hasLevel says so.
	state    string
	attrs    map[string]any
	level    float64
	hasLevel bool
	// pending counts calls on their way, and a level waiting for its
	// slider to rest, which waiting reports and gen numbers; edits is the
	// count of local changes once the last call returned, and until when
	// it settles.
	pending int
	waiting bool
	gen     int
	edits   int64
	until   time.Time
}

// tileOf finds the tile of entity id on the home page, and the area it
// shows in.
func (a *application) tileOf(id string) (ha.Tile, string, bool) {
	if id == "" {
		return ha.Tile{}, "", false
	}
	for _, s := range a.home {
		for _, t := range s.Tiles {
			if t.Entity.EntityID == id {
				return t, s.Area, true
			}
		}
	}
	return ha.Tile{}, "", false
}

// tileCount is how many tiles the home page shows.
func (a *application) tileCount() int {
	n := 0
	for _, s := range a.home {
		n += len(s.Tiles)
	}
	return n
}

// held is entity s as its tile shows it: in the state, and with the
// attributes, its last action leads to, until Home Assistant reports them.
func (a *application) held(s ha.State) ha.State {
	h := a.holds[s.EntityID]
	if h == nil {
		return s
	}
	if h.state != "" {
		s.State = h.state
	}
	if len(h.attrs) > 0 {
		attrs := maps.Clone(s.Attributes)
		if attrs == nil {
			attrs = map[string]any{}
		}
		maps.Copy(attrs, h.attrs)
		s.Attributes = attrs
	}
	return s
}

// show holds what act leads to: its state, unless it keeps the state read,
// and the attributes it shows.
func (h *hold) show(act ha.Action) {
	if act.Next != "" {
		h.state = act.Next
	}
	if len(act.Shows) > 0 {
		if h.attrs == nil {
			h.attrs = map[string]any{}
		}
		maps.Copy(h.attrs, act.Shows)
	}
}

// holdOf is the hold of entity id, made as needed.
func (a *application) holdOf(id string) *hold {
	h := a.holds[id]
	if h == nil {
		h = &hold{}
		if a.holds == nil {
			a.holds = map[string]*hold{}
		}
		a.holds[id] = h
	}
	return h
}

// level is level l of entity s as its page shows it: as last set, until
// Home Assistant reports it.
func (a *application) level(s ha.State, l ha.Level) float64 {
	if h := a.holds[s.EntityID]; h != nil && h.hasLevel {
		return h.level
	}
	return l.Value
}

// setLevel shows level v of entity s at once and sends it once the slider
// rests.
func (a *application) setLevel(s ha.State, l ha.Level, v float64) {
	h := a.holdOf(s.EntityID)
	h.level, h.hasLevel = v, true
	a.adjust(s, l.Set(v))
}

// adjust shows what act leads to at once and asks for it once the control
// setting it rests: of several in a row, only the last.
func (a *application) adjust(s ha.State, act ha.Action) {
	h := a.holdOf(s.EntityID)
	if !h.waiting {
		h.waiting = true
		h.pending++
	}
	h.show(act)
	h.gen++
	gen := h.gen
	a.edits.Add(1)
	a.redraw()
	time.AfterFunc(debounce, func() {
		a.post(func() {
			if a.holds[s.EntityID] != h || h.gen != gen {
				return
			}
			h.waiting = false
			h.pending--
			a.perform(s, act)
		})
	})
}

// toggle does what the quick switch of entity s does. The air conditioner
// opened last switches as its controls switch it.
func (a *application) toggle(s ha.State) {
	if s.Domain() == "climate" && a.hasClimate && a.climate.EntityID == s.EntityID {
		a.togglePower()
		return
	}
	if act, ok := s.Toggle(); ok {
		a.perform(s, act)
	}
}

// perform asks Home Assistant for action act of entity s, showing at once the
// state it leads to.
func (a *application) perform(s ha.State, act ha.Action) {
	a.mu.Lock()
	client := a.client
	a.mu.Unlock()
	if client == nil {
		return
	}
	// An action that leads to a state or shows attributes holds them until
	// Home Assistant reports them; a hold already made, as of a level,
	// waits for the action too.
	h := a.holds[s.EntityID]
	if act.Next != "" || len(act.Shows) > 0 {
		h = a.holdOf(s.EntityID)
		h.show(act)
	}
	if h != nil {
		h.pending++
	}
	a.edits.Add(1)
	a.redraw()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := client.Call(ctx, s.Domain(), act.Service, s.EntityID, act.Data)
		a.post(func() {
			if !a.current(client) {
				return
			}
			if h != nil && a.holds[s.EntityID] == h {
				h.pending--
				h.edits, h.until = a.edits.Add(1), time.Now().Add(settle)
				if err != nil {
					// Show what the device really does.
					delete(a.holds, s.EntityID)
				}
			}
			if err != nil {
				a.report("requestError", err)
				if errors.Is(err, ha.ErrUnauthorized) {
					a.errorKey = "unauthorized"
				}
				a.wake()
			} else {
				time.AfterFunc(settle, a.wake)
			}
			a.redraw()
		})
	}()
}

// release lets go of the holds a read is fresh for: one that started once
// their calls returned and settled.
func (a *application) release(edits int64) {
	now := time.Now()
	for id, h := range a.holds {
		if h.pending == 0 && edits >= h.edits && !now.Before(h.until) {
			delete(a.holds, id)
		}
	}
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

// redraw refreshes the panel after a change. The menu bar icon never
// changes.
func (a *application) redraw() {
	if a.panel != nil {
		a.panel.Invalidate()
	}
	if a.window != nil {
		a.window.Window().Invalidate()
	}
}

func formatTemp(t float64) string { return strconv.FormatFloat(t, 'f', -1, 64) }

// message is the message of key in the app's language, or key itself when
// the catalogs have none.
func message(key string) string {
	if _, _, ok := locales.Catalog.Message(i18n.Locale(), key); ok {
		return locales.Catalog.Format(i18n.Locale(), key, nil)
	}
	return key
}

// optionText names a Home Assistant option, as "fan_only", in the app's
// language; options it has no message for show as they are.
func optionText(kind, value string) string {
	if value == "" {
		return "—"
	}
	if m := message(kind + "_" + value); m != kind+"_"+value {
		return m
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
	c.OnShortcut(ui.Cmd, ui.KeyQ, a.quit)
	c.OnShortcut(ui.Cmd, ui.KeyComma, a.showSettings)
}
