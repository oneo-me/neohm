package main

import (
	"bytes"
	"encoding/json"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/neohm/source/internal/ha"
	"github.com/oneo-me/neohm/source/locales"
)

func parseStates(t *testing.T, list string) []ha.State {
	t.Helper()
	var states []ha.State
	if err := json.Unmarshal([]byte(list), &states); err != nil {
		t.Fatal(err)
	}
	return states
}

// demoStates are a sample home: two air conditioners, the bedroom's first,
// lights, a curtain, a TV, a lock, a plug, a thermometer, a window sensor,
// a vacuum and a fan.
const demoStates = `[
	{"entity_id":"climate.bedroom","state":"cool","attributes":{
		"friendly_name":"Bedroom AC","hvac_modes":["off","cool","heat","dry","fan_only","auto"],
		"min_temp":16,"max_temp":30,"target_temp_step":1,"current_temperature":24,"temperature":16,
		"fan_modes":["auto","low","medium","high"],"fan_mode":"low",
		"swing_modes":["on","off"],"swing_mode":"on",
		"preset_modes":["none","eco","sleep"],"preset_mode":"none"}},
	{"entity_id":"climate.living_room","state":"off","attributes":{
		"friendly_name":"Living Room Cabinet AC","hvac_modes":["off","cool","heat"],
		"min_temp":16,"max_temp":30,"current_temperature":26.5,"temperature":24}},
	{"entity_id":"light.ceiling","state":"on","attributes":{"friendly_name":"Living Room Ceiling Light","brightness":204,"rgb_color":[255,180,107],"color_mode":"color_temp","color_temp_kelvin":2604,"min_color_temp_kelvin":2000,"max_color_temp_kelvin":6535,"supported_color_modes":["color_temp","hs"]}},
	{"entity_id":"cover.curtain","state":"open","attributes":{"friendly_name":"Living Room Curtain","device_class":"curtain","current_position":60,"supported_features":15}},
	{"entity_id":"media_player.tv","state":"playing","attributes":{"friendly_name":"Living Room TV","device_class":"tv","media_title":"Interstellar","supported_features":16817}},
	{"entity_id":"lock.front_door","state":"locked","attributes":{"friendly_name":"Front Door","supported_features":0}},
	{"entity_id":"light.lamp","state":"off","attributes":{"friendly_name":"Bedroom Lamp"}},
	{"entity_id":"switch.heater","state":"off","attributes":{"friendly_name":"Bedroom Heater","device_class":"outlet"}},
	{"entity_id":"sensor.thermo_temperature","state":"24.27","attributes":{"friendly_name":"Bedroom Thermometer Temperature","device_class":"temperature","unit_of_measurement":"°C"}},
	{"entity_id":"sensor.thermo_humidity","state":"45","attributes":{"friendly_name":"Bedroom Thermometer Humidity","device_class":"humidity","unit_of_measurement":"%"}},
	{"entity_id":"sensor.thermo_battery","state":"88","attributes":{"device_class":"battery","unit_of_measurement":"%"}},
	{"entity_id":"binary_sensor.window","state":"on","attributes":{"device_class":"window"}},
	{"entity_id":"vacuum.robot","state":"docked","attributes":{"friendly_name":"Robot Vacuum","supported_features":12316}},
	{"entity_id":"fan.bedroom","state":"on","attributes":{"friendly_name":"Bedroom Fan","percentage":40,"percentage_step":10,"supported_features":1}}
]`

func demoRegistry() *ha.Registry {
	return &ha.Registry{
		Entities: map[string]ha.RegistryEntry{
			"climate.bedroom": {Device: "ac"}, "light.lamp": {Device: "lamp"}, "switch.heater": {Device: "plug"},
			"sensor.thermo_temperature": {Device: "thermo", Precision: 1, HasPrecision: true},
			"sensor.thermo_humidity":    {Device: "thermo"}, "sensor.thermo_battery": {Device: "thermo", Hidden: true},
			"binary_sensor.window": {Device: "window"}, "vacuum.robot": {Device: "vacuum"}, "fan.bedroom": {Device: "fan"},
			"climate.living_room": {Area: "living_room"}, "light.ceiling": {Area: "living_room"},
			"cover.curtain": {Area: "living_room"}, "media_player.tv": {Area: "living_room"},
			"lock.front_door": {Area: "living_room"},
		},
		Devices: map[string]ha.Device{
			"ac": {Name: "Bedroom AC", Area: "bedroom"}, "lamp": {Name: "Bedroom Lamp", Area: "bedroom"},
			"plug": {Name: "Bedroom Plug", Area: "bedroom"}, "thermo": {Name: "Bedroom Thermometer", Area: "bedroom"},
			"window": {Name: "Bedroom Window", Area: "bedroom"}, "vacuum": {Name: "Robot Vacuum"}, "fan": {Name: "Bedroom Fan", Area: "bedroom"},
		},
		Areas: []ha.Area{{ID: "living_room", Name: "Living Room"}, {ID: "bedroom", Name: "Bedroom"}},
	}
}

// demoApp is an application showing the sample home, the bedroom's air
// conditioner opened last, without a connection: service calls change only
// what it shows.
func demoApp(t *testing.T) *application {
	states := parseStates(t, demoStates)
	a := newApplication(filepath.Join(t.TempDir(), "config.json"))
	a.post = func(fn func()) {}
	a.config = config{Address: "http://ha.test", Token: "token", Entity: "climate.bedroom"}
	a.address, a.token = a.config.Address, a.config.Token
	a.registry = demoRegistry()
	a.home, a.listed = ha.Arrange(states, a.registry), true
	a.climate, a.hasClimate, a.online = ha.ParseClimate(states[0]), true, true
	return a
}

// serviceCalls records the service calls of Home Assistant at its URL, as
// "light/turn_off light.ceiling".
func serviceCalls(t *testing.T) (string, func() []string) {
	var mu sync.Mutex
	var calls []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		mu.Lock()
		calls = append(calls, r.URL.Path[len("/api/services/"):]+" "+body["entity_id"].(string))
		mu.Unlock()
		w.Write([]byte(`[]`))
	}))
	t.Cleanup(srv.Close)
	return srv.URL, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return slices.Clone(calls)
	}
}

// returned runs what a service call posts to the main thread once it
// returns.
func returned(t *testing.T, posted chan func()) {
	t.Helper()
	select {
	case fn := <-posted:
		fn()
	case <-time.After(5 * time.Second):
		t.Fatal("no service call returned")
	}
}

// focus moves the keyboard to the control labeled label with Tab.
func focus(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	for range 20 {
		if tt.Focused(label) {
			return
		}
		tt.Key(0, ui.KeyTab)
	}
	t.Fatalf("%q never takes the keyboard", label)
}

func TestPanelControls(t *testing.T) {
	a := demoApp(t)
	a.detail = true
	tt := ui.NewTester(a.panelView, panelWidth, panelHeight)
	for _, s := range []string{"Bedroom AC", locales.Room(), "24°C", "16", locales.Eco(), locales.ModeCool(), locales.FanLow(), locales.SwingOn()} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	if err := tt.Click(locales.Raise()); err != nil {
		t.Fatal(err)
	}
	if a.climate.Target != 17 {
		t.Fatalf("raise: %v", a.climate.Target)
	}
	focus(t, tt, locales.Eco())
	tt.Key(0, ui.KeySpace)
	if a.climate.Preset != "eco" {
		t.Fatalf("eco: %q", a.climate.Preset)
	}
	a.togglePreset("eco")
	if err := tt.Click(locales.Sleep() + ": " + locales.Off()); err != nil {
		t.Fatal(err)
	}
	if a.climate.Preset != "sleep" {
		t.Fatalf("sleep: %q", a.climate.Preset)
	}
	a.togglePreset("sleep")
	if a.climate.Preset != "none" {
		t.Fatalf("sleep off: %q", a.climate.Preset)
	}
	if err := tt.Click(locales.Power()); err != nil {
		t.Fatal(err)
	}
	if a.climate.Mode != "off" {
		t.Fatalf("power off: %q", a.climate.Mode)
	}
	tt.Advance(time.Second)
	if err := tt.Click(locales.Power()); err != nil {
		t.Fatal(err)
	}
	if a.climate.Mode != "cool" {
		t.Fatalf("power on: %q", a.climate.Mode)
	}
	tt.Advance(time.Second)
	if err := tt.Click(locales.Mode() + ": " + locales.ModeCool()); err != nil {
		t.Fatal(err)
	}
	if err := tt.ChooseMenuItem(locales.ModeHeat()); err != nil {
		t.Fatalf("%v: %v", err, tt.Menu())
	}
	if a.climate.Mode != "heat" {
		t.Fatalf("mode: %q", a.climate.Mode)
	}
	if err := tt.Click(locales.Back()); err != nil {
		t.Fatal(err)
	}
	if a.detail || !tt.HasText(locales.Title()) {
		t.Fatalf("back: %v", tt.Texts())
	}
}

// The panel opens on the devices as tiles, by area, each with its state; an
// air conditioner opens its controls and is remembered.
func TestHome(t *testing.T) {
	a := demoApp(t)
	a.climate.Target = 21 // A local change the read lacks.
	tt := ui.NewTester(a.panelView, panelWidth, 900)
	for _, s := range []string{
		locales.Title(), "Living Room", "Bedroom", locales.OtherDevices(),
		"Cabinet AC", locales.ModeOff(), "AC", locales.ModeCool() + " · 21°",
		"Ceiling Light", locales.StateOn() + " · 80%", "Curtain", locales.StateOpen() + " · 60%",
		"TV", locales.StatePlaying() + " · Interstellar", "Front Door", locales.StateLocked(),
		"Thermometer", "24.3 °C · 45%", "Window", "Heater", "Robot Vacuum", locales.StateDocked(),
	} {
		if _, ok := tt.Find(s); !ok {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	if err := tt.Click(locales.DeviceDetails("Cabinet AC")); err != nil {
		t.Fatal(err)
	}
	if !a.detail || a.climate.EntityID != "climate.living_room" || a.entity != "climate.living_room" {
		t.Fatalf("open: %v %q", a.detail, a.climate.EntityID)
	}
	if saved, err := loadConfig(a.path); err != nil || saved.Entity != "climate.living_room" {
		t.Fatalf("saved %+v %v", saved, err)
	}
	if !tt.HasText(locales.Power()) {
		t.Fatalf("detail: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	if a.detail || !tt.HasText("Ceiling Light") {
		t.Fatalf("escape: %v", tt.Texts())
	}
}

// A click on a tile switches its device, and the device's page, opened from
// the tile's corner, offers what else it does; either shows the state an action leads to until a fresh
// read reports it.
func TestTileActions(t *testing.T) {
	url, calls := serviceCalls(t)
	a := demoApp(t)
	posted := make(chan func(), 4)
	a.post = func(fn func()) { posted <- fn }
	client := ha.NewClient(url, "token")
	a.client = client
	tt := ui.NewTester(a.panelView, panelWidth, 900)

	edits := a.edits.Load()
	if err := tt.Click(locales.QuickSwitch("Ceiling Light")); err != nil {
		t.Fatal(err)
	}
	if a.device != "" || !tt.HasText(locales.StateOff()) || tt.HasText(locales.StateOn()+" · 80%") {
		t.Fatalf("turn off: %q %v", a.device, tt.Texts())
	}
	returned(t, posted)
	states := parseStates(t, demoStates)
	a.applyList(client, edits, states, nil, "", nil)
	tt.Frame()
	if tt.HasText(locales.StateOn() + " · 80%") {
		t.Fatal("a read from before the action undid it")
	}
	a.holds["light.ceiling"].until = time.Time{} // As if it settled.
	states[2].State = "off"
	a.applyList(client, a.edits.Load(), states, nil, "", nil)
	if len(a.holds) != 0 {
		t.Fatalf("holds %v", a.holds)
	}

	// The air conditioner opened last switches as its controls do; another
	// turns on in its first mode.
	if err := tt.Click(locales.QuickSwitch("AC")); err != nil {
		t.Fatal(err)
	}
	if a.climate.Mode != "off" {
		t.Fatalf("opened last: %q", a.climate.Mode)
	}
	if err := tt.Click(locales.QuickSwitch("Cabinet AC")); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText(locales.ModeCool() + " · 24°") {
		t.Fatalf("turn on: %v", tt.Texts())
	}
	returned(t, posted)

	// A lock does not switch: its tile opens its page, which offers what it
	// does.
	if _, ok := tt.Find(locales.QuickSwitch("Front Door")); ok {
		t.Fatal("a lock switches from its tile")
	}
	if _, ok := tt.Find(locales.DeviceDetails("Front Door")); ok {
		t.Fatal("a lock's tile has a details button")
	}
	if err := tt.Click("Front Door"); err != nil {
		t.Fatal(err)
	}
	if a.device != "lock.front_door" || !tt.HasText(locales.ActionLock()) {
		t.Fatalf("lock page: %q %v", a.device, tt.Texts())
	}
	if err := tt.Click(locales.ActionUnlock()); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText(locales.StateUnlocking()) {
		t.Fatalf("unlock: %v", tt.Texts())
	}
	returned(t, posted)
	tt.Key(0, ui.KeyEscape)
	if a.device != "" {
		t.Fatalf("escape: %q", a.device)
	}
	if err := tt.Click("Robot Vacuum"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click(locales.ActionStart()); err != nil {
		t.Fatalf("%v: %v", err, tt.Texts())
	}
	returned(t, posted)
	want := []string{"light/turn_off light.ceiling", "climate/set_hvac_mode climate.living_room",
		"lock/unlock lock.front_door", "vacuum/start vacuum.robot"}
	if got := calls(); !slices.Equal(got, want) {
		t.Fatalf("calls %v", got)
	}
}

// A device's page sends the level its slider sets once the slider rests,
// showing it meanwhile, a sensor's page shows its readings and a light's
// switches from its power button.
func TestDevicePage(t *testing.T) {
	var mu sync.Mutex
	var bodies []map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		body["service"] = r.URL.Path[len("/api/services/"):]
		mu.Lock()
		bodies = append(bodies, body)
		mu.Unlock()
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()
	a := demoApp(t)
	posted := make(chan func(), 4)
	a.post = func(fn func()) { posted <- fn }
	a.client = ha.NewClient(srv.URL, "token")
	a.panelOpen.Store(true)
	tt := ui.NewTester(a.panelView, panelWidth, 900)
	if err := tt.Click(locales.DeviceDetails("Curtain")); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"Living Room", locales.LevelPosition(), "60", locales.ActionOpen(), locales.ActionStop(), locales.ActionClose()} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	if !a.listing {
		t.Fatal("a device's page stopped reading every device")
	}
	focus(t, tt, locales.LevelPosition())
	tt.Key(0, ui.KeyEnd)
	if _, ok := tt.Find("100"); !ok {
		t.Fatalf("level not shown: %v", tt.Texts())
	}
	returned(t, posted) // The slider rests.
	returned(t, posted) // The call returns.
	mu.Lock()
	if len(bodies) != 1 || bodies[0]["service"] != "cover/set_cover_position" || bodies[0]["position"] != 100.0 {
		t.Fatalf("calls %v", bodies)
	}
	mu.Unlock()
	// Until Home Assistant reports it, the slider keeps the level set.
	a.applyList(a.client, a.edits.Load(), parseStates(t, demoStates), nil, "", nil)
	tt.Frame()
	if _, ok := tt.Find("100"); !ok {
		t.Fatalf("a read undid the level: %v", tt.Texts())
	}

	tt.Key(0, ui.KeyEscape)
	if err := tt.Click("Thermometer"); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"Temperature", "24.3", "°C", "Humidity", "45%"} {
		if _, ok := tt.Find(s); !ok {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	// A light switches from the power button of its page; a cover has none.
	tt.Key(0, ui.KeyEscape)
	if err := tt.Click(locales.DeviceDetails("Ceiling Light")); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click(locales.QuickSwitch("Ceiling Light")); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText(locales.StateOff()) {
		t.Fatalf("power: %v", tt.Texts())
	}
	returned(t, posted)
	tt.Key(0, ui.KeyEscape)
	if err := tt.Click(locales.DeviceDetails("Curtain")); err != nil {
		t.Fatal(err)
	}
	if _, ok := tt.Find(locales.QuickSwitch("Curtain")); ok {
		t.Fatal("a cover has a power button")
	}
	tt.Key(0, ui.KeyEscape)
	if err := tt.Click("Thermometer"); err != nil {
		t.Fatal(err)
	}
	// The page closes once its device is gone.
	a.applyList(a.client, a.edits.Load(), parseStates(t, demoStates)[:8], nil, "", nil)
	if a.device != "" {
		t.Fatalf("gone device: %q", a.device)
	}
}

// A cover's page sets its quarters a click away, and a light's its colors,
// showing them before Home Assistant reports them.
func TestQuickSettings(t *testing.T) {
	var mu sync.Mutex
	var bodies []map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		body["service"] = r.URL.Path[len("/api/services/"):]
		mu.Lock()
		bodies = append(bodies, body)
		mu.Unlock()
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()
	a := demoApp(t)
	posted := make(chan func(), 4)
	a.post = func(fn func()) { posted <- fn }
	a.client = ha.NewClient(srv.URL, "token")
	a.panelOpen.Store(true)
	tt := ui.NewTester(a.panelView, panelWidth, 900)
	sent := func() map[string]any {
		t.Helper()
		mu.Lock()
		defer mu.Unlock()
		if len(bodies) != 1 {
			t.Fatalf("calls %v", bodies)
		}
		body := bodies[0]
		bodies = nil
		return body
	}

	if err := tt.Click(locales.DeviceDetails("Curtain")); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click(locales.LevelPosition() + ": 25%"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("25") {
		t.Fatalf("position not shown: %v", tt.Texts())
	}
	returned(t, posted) // The level rests.
	returned(t, posted) // The call returns.
	if body := sent(); body["service"] != "cover/set_cover_position" || body["position"] != 25.0 {
		t.Fatalf("stop sent %v", body)
	}

	tt.Key(0, ui.KeyEscape)
	if err := tt.Click(locales.DeviceDetails("Ceiling Light")); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"2600 K", locales.ColorRed(), locales.ChooseColor()} {
		if _, ok := tt.Find(s); !ok {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	if err := tt.Click(locales.ColorBlue()); err != nil {
		t.Fatal(err)
	}
	tile, _, _ := a.tileOf("light.ceiling")
	if shining, ok := a.held(tile.Entity).Shining(); !ok || !shining.Near(ha.Hues[4]) {
		t.Fatalf("color not shown: %+v", shining)
	}
	returned(t, posted)
	if body := sent(); body["service"] != "light/turn_on" || !slices.Equal(body["hs_color"].([]any), []any{215.0, 80.0}) {
		t.Fatalf("color sent %v", body)
	}
	if err := tt.Click(locales.ChooseColor()); err != nil {
		t.Fatal(err)
	}
	if _, ok := tt.Find(locales.ColorTemperature()); !ok {
		t.Fatalf("picker not open: %v", tt.Texts())
	}
	focus(t, tt, locales.ColorTemperature())
	tt.Key(0, ui.KeyHome)
	returned(t, posted) // The slider rests.
	returned(t, posted)
	if body := sent(); body["color_temp_kelvin"] != 2000.0 {
		t.Fatalf("white sent %v", body)
	}
}

// A read keeps a local change to the air conditioner opened last, forgets it
// once it is gone and ignores another connection.
func TestListRead(t *testing.T) {
	a := demoApp(t)
	client := ha.NewClient("http://ha.test", "token")
	a.client, a.detail = client, true
	a.panelOpen.Store(true)
	read := parseStates(t, demoStates)
	edits := a.edits.Load()
	a.climate.Target = 20
	a.edits.Add(1)
	a.applyList(client, edits, read, nil, "", nil)
	cl := ha.ParseClimate(read[0])
	if a.climate.Target != 20 || a.shown(cl).Target != 20 || a.tileCount() != 12 {
		t.Fatalf("stale read replaced a local change: %v, %d tiles", a.climate.Target, a.tileCount())
	}
	a.applyList(client, a.edits.Load(), read[1:], nil, "", nil)
	if a.hasClimate || a.detail || a.tileCount() != 11 || !a.listing || a.config.Entity != "" || a.entity != "" {
		t.Fatalf("gone device: %v %v %v %q", a.hasClimate, a.detail, a.listing, a.config.Entity)
	}
	if saved, err := loadConfig(a.path); err != nil || saved.Entity != "" {
		t.Fatalf("saved %+v %v", saved, err)
	}
	a.applyList(ha.NewClient("http://ha.test", "old"), a.edits.Load(), read, nil, "", nil)
	if a.tileCount() != 11 {
		t.Fatal("read of another connection shown")
	}
	// The registries last read keep arranging the tiles.
	if a.registry == nil || len(a.home) != 3 || a.home[0].Area != "Living Room" {
		t.Fatalf("sections %+v", a.home)
	}
}

// A temperature still waiting for the slider goes to its own device when the
// panel opens another.
func TestOpenDeviceSendsWaitingTarget(t *testing.T) {
	var mu sync.Mutex
	var calls []map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		mu.Lock()
		calls = append(calls, body)
		mu.Unlock()
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()
	a := demoApp(t)
	posted := make(chan func(), 4)
	a.post = func(fn func()) { posted <- fn }
	a.client, a.entity, a.detail = ha.NewClient(srv.URL, "token"), "climate.bedroom", true
	a.setTarget(22)
	a.openClimate(ha.ParseClimate(parseStates(t, demoStates)[1]))
	returned(t, posted)
	mu.Lock()
	defer mu.Unlock()
	if len(calls) != 1 || calls[0]["entity_id"] != "climate.bedroom" || calls[0]["temperature"] != 22.0 {
		t.Fatalf("calls %v", calls)
	}
	if a.tempPending || a.pending != 0 || a.climate.EntityID != "climate.living_room" {
		t.Fatalf("after: %v %v %q", a.tempPending, a.pending, a.climate.EntityID)
	}
}

func TestStaleReadIgnored(t *testing.T) {
	a := demoApp(t)
	client := ha.NewClient("http://ha.test", "token")
	a.client, a.entity = client, "climate.bedroom"
	edits := a.edits.Load()
	old := a.climate
	a.climate.Target = 20
	a.edits.Add(1)
	a.applyState(client, edits, old, "", nil)
	if a.climate.Target != 20 {
		t.Fatal("stale read replaced a local change")
	}
	a.applyState(client, a.edits.Load(), old, "", nil)
	if a.climate.Target != 16 {
		t.Fatal("fresh read ignored")
	}
}

func TestTileName(t *testing.T) {
	for _, tc := range []struct{ name, area, want string }{
		{"Bedroom Lamp", "Bedroom", "Lamp"},
		{"bedroom lamp", "Bedroom", "lamp"},
		{"卧室空调", "卧室", "空调"},
		{"Bedroom", "Bedroom", "Bedroom"},
		{"Bedroom - ", "Bedroom", "Bedroom - "},
		{"Lamp", "Bedroom", "Lamp"},
		{"Lamp", "", "Lamp"},
	} {
		if got := tileName(tc.name, tc.area); got != tc.want {
			t.Errorf("%q in %q: %q", tc.name, tc.area, got)
		}
	}
}

func TestSetupView(t *testing.T) {
	a := newApplication(filepath.Join(t.TempDir(), "config.json"))
	tt := ui.NewTester(a.panelView, panelWidth, panelHeight)
	if !tt.HasText(locales.SetupTitle()) {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click(locales.Settings()); err != nil {
		t.Fatal(err)
	}
	if !a.settingsOpen || !tt.HasText(locales.Address()) {
		t.Fatal(tt.Texts())
	}
	// Until Home Assistant is set up, the panel opens on Settings.
	a.settingsOpen = false
	a.resetPanel()
	if !a.settingsOpen {
		t.Fatal("reset to the home page before setup")
	}
	// Connect saves the token with the address, for the next launch to read.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`[]`)) }))
	defer srv.Close()
	a.post = func(func()) {}
	a.address, a.token = srv.URL, "secret"
	tt.Frame()
	if err := tt.Click(locales.Connect()); err != nil {
		t.Fatal(err)
	}
	if saved, err := loadConfig(a.path); err != nil || saved.Address != srv.URL || saved.Token != "secret" || !a.configured() {
		t.Fatalf("saved %+v %v", saved, err)
	}
}

// Settings shows in the panel, over the page it covers.
func TestSettingsPage(t *testing.T) {
	a := demoApp(t)
	tt := ui.NewTester(a.panelView, panelWidth, panelHeight)
	if err := tt.Click(locales.SettingsTitle()); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{locales.Connection(), locales.Address(), locales.Token(), locales.Connect()} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	// The language follows the system until another is chosen.
	var shown []string
	a.setLocale = func(tag string) { shown = append(shown, tag) }
	if err := tt.Click(locales.LanguageSystem()); err != nil {
		t.Fatal(err)
	}
	if err := tt.ChooseMenuItem("日本語"); err != nil {
		t.Fatalf("%v: %v", err, tt.Menu())
	}
	if saved, err := loadConfig(a.path); err != nil || a.config.Language != "ja" || saved.Language != "ja" || !slices.Equal(shown, []string{"ja"}) {
		t.Fatalf("language %q, saved %+v %v, shown %v", a.config.Language, saved, err, shown)
	}
	if err := tt.Click(locales.CloseSettings()); err != nil {
		t.Fatal(err)
	}
	if a.settingsOpen || !tt.HasText("Cabinet AC") {
		t.Fatal(tt.Texts())
	}
	a.detail = true
	a.showSettings()
	tt.Key(0, ui.KeyEscape)
	if a.settingsOpen || !tt.HasText(locales.Power()) {
		t.Fatal(tt.Texts())
	}
	a.resetPanel()
	if a.settingsOpen || a.detail {
		t.Fatal("reset to Settings or the device")
	}
}

// Command-Q quits from every page, a text field's keyboard included.
func TestQuitShortcut(t *testing.T) {
	a := demoApp(t)
	quits := 0
	a.quit = func() { quits++ }
	tt := ui.NewTester(a.panelView, panelWidth, panelHeight)
	tt.Key(ui.Cmd, ui.KeyQ)
	a.showSettings()
	focus(t, tt, locales.Address())
	tt.Key(ui.Cmd, ui.KeyQ)
	if quits != 2 {
		t.Fatalf("quit %d times, want 2", quits)
	}
}

// The configuration file holds the token, so only the user can read it.
func TestConfigRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x", "config.json")
	want := config{Address: "http://ha.local:8123", Token: "secret", Entity: "climate.a", Language: "zh-Hant"}
	if err := saveConfig(path, want); err != nil {
		t.Fatal(err)
	}
	if got, err := loadConfig(path); err != nil || got != want {
		t.Fatalf("%+v %v", got, err)
	}
	// A language the app is not written in follows the system's.
	if err := saveConfig(path, config{Language: "fr"}); err != nil {
		t.Fatal(err)
	}
	if got, err := loadConfig(path); err != nil || got.Language != "" {
		t.Fatalf("%+v %v", got, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 && runtime.GOOS != "windows" {
		t.Fatalf("mode %v", perm)
	}
}

// Readers of each language the app is written in get its catalog, and every
// catalog says every message.
func TestLanguages(t *testing.T) {
	for system, want := range map[string]string{
		"en-US": "en", "zh-CN": "zh-Hans", "zh-SG": "zh-Hans", "zh-TW": "zh-Hant", "zh-HK": "zh-Hant",
		"ja-JP": "ja", "ko-KR": "ko",
	} {
		if got := i18n.Match([]string{system}, locales.Catalog.Languages()); got != want {
			t.Errorf("%s reads %q, not %q", system, got, want)
		}
	}
	catalog := func(tag string) map[string]any {
		data, err := os.ReadFile(filepath.Join("locales", tag+".json"))
		if err != nil {
			t.Fatal(err)
		}
		var m map[string]any
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatal(err)
		}
		return m
	}
	source := catalog(locales.Catalog.Source())
	for _, tag := range locales.Catalog.Languages() {
		m := catalog(tag)
		for key := range source {
			if _, ok := m[key]; !ok && !strings.HasPrefix(key, "@") {
				t.Errorf("%s lacks %s", tag, key)
			}
		}
	}
}

// The menu bar icon is a template image, centered on the menu bar.
func TestTrayIcon(t *testing.T) {
	img, err := png.Decode(bytes.NewReader(trayIcon()))
	if err != nil || img.Bounds().Dx() != 36 || img.Bounds().Dy() != 36 {
		t.Fatal(err)
	}
	for y := range 36 {
		for x := range 18 {
			r, g, b, a := img.At(x, y).RGBA()
			if _, _, _, mirror := img.At(35-x, y).RGBA(); a != mirror {
				t.Fatalf("not symmetric at %d,%d", x, y)
			}
			if runtime.GOOS == "darwin" && r|g|b != 0 {
				t.Fatalf("not black at %d,%d", x, y)
			}
		}
	}
}

// TestSnapshot writes the panel's pages and the menu bar icon to
// $HA_SNAPSHOT_DIR for review.
func TestSnapshot(t *testing.T) {
	dir := os.Getenv("HA_SNAPSHOT_DIR")
	if dir == "" {
		t.Skip("set HA_SNAPSHOT_DIR to write snapshots")
	}
	if err := os.WriteFile(filepath.Join(dir, "tray.png"), trayIcon(), 0o644); err != nil {
		t.Fatal(err)
	}
	a := demoApp(t)
	tt := ui.NewTester(a.panelView, panelWidth, panelHeight)
	tt.SetScale(2)
	write := func(name string) {
		tt.Frame()
		tt.Advance(800 * time.Millisecond)
		f, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err := png.Encode(f, tt.Image()); err != nil {
			t.Fatal(err)
		}
	}
	write("home.png")
	tt.Scroll(panelWidth/2, 300, 0, 200)
	write("home-scrolled.png")
	a.settingsOpen = true
	write("settings.png")
	a.settingsOpen, a.detail = false, true
	write("panel.png")
	a.climate.Mode = "heat"
	a.climate.Preset = "sleep"
	write("panel-heat.png")
	for _, id := range []string{"light.ceiling", "light.lamp", "cover.curtain", "media_player.tv", "sensor.thermo_temperature",
		"binary_sensor.window", "lock.front_door", "switch.heater", "fan.bedroom", "vacuum.robot"} {
		a.openDevice(id)
		write(strings.ReplaceAll(id, ".", "-") + ".png")
	}
	a.openDevice("light.ceiling")
	tt.Frame()
	if err := tt.Click(locales.ChooseColor()); err != nil {
		t.Fatal(err)
	}
	write("light-ceiling-colors.png")
}
