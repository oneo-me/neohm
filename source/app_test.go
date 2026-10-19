package main

import (
	"bytes"
	"encoding/json"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/ha-go/source/internal/ha"
	"github.com/oneo-me/ha-go/source/locales"
)

type memorySecrets map[string][]byte

func (m memorySecrets) Get(k string) ([]byte, error) { return m[k], nil }
func (m memorySecrets) Set(k string, v []byte) error { m[k] = v; return nil }

func parseClimate(t *testing.T, state string) ha.Climate {
	var s ha.State
	if err := json.Unmarshal([]byte(state), &s); err != nil {
		t.Fatal(err)
	}
	return ha.ParseClimate(s)
}

// demoApp is an application listing two sample units, the first opened last,
// without a connection: service calls change only what it shows.
func demoApp(t *testing.T) *application {
	bedroom := parseClimate(t, `{"entity_id":"climate.bedroom","state":"cool","attributes":{
		"friendly_name":"Bedroom AC","hvac_modes":["off","cool","heat","dry","fan_only","auto"],
		"min_temp":16,"max_temp":30,"target_temp_step":1,"current_temperature":24,"temperature":16,
		"fan_modes":["auto","low","medium","high"],"fan_mode":"low",
		"swing_modes":["on","off"],"swing_mode":"on",
		"preset_modes":["none","eco","sleep"],"preset_mode":"none"}}`)
	living := parseClimate(t, `{"entity_id":"climate.living_room","state":"off","attributes":{
		"friendly_name":"Living Room AC","hvac_modes":["off","cool","heat"],
		"min_temp":16,"max_temp":30,"current_temperature":26.5,"temperature":24}}`)
	a := newApplication(filepath.Join(t.TempDir(), "config.json"), memorySecrets{})
	a.post = func(fn func()) {}
	a.config = config{Address: "http://ha.test", Entity: "climate.bedroom"}
	a.token = "token"
	a.devices, a.listed = []ha.Climate{bedroom, living}, true
	a.climate, a.hasClimate, a.online = bedroom, true, true
	return a
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
	if a.climate.Target != 17 || a.trayTitle() != "17°" {
		t.Fatalf("raise: %v %q", a.climate.Target, a.trayTitle())
	}
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
	if a.climate.Mode != "off" || a.trayTitle() != "" {
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
	if a.detail || !tt.HasText(locales.Devices()) {
		t.Fatalf("back: %v", tt.Texts())
	}
}

// The panel opens on the device list; a device opens its controls and is
// remembered for the menu bar title.
func TestDeviceList(t *testing.T) {
	a := demoApp(t)
	a.climate.Target = 21 // A local change the list's read lacks.
	tt := ui.NewTester(a.panelView, panelWidth, panelHeight)
	for _, s := range []string{locales.Devices(), "Bedroom AC", "21°", "Living Room AC", locales.ModeOff() + " · " + locales.Room() + " 26.5°C"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	if err := tt.Click("Living Room AC"); err != nil {
		t.Fatal(err)
	}
	if !a.detail || a.climate.EntityID != "climate.living_room" || a.entity != "climate.living_room" {
		t.Fatalf("open: %v %q", a.detail, a.climate.EntityID)
	}
	if saved, err := loadConfig(a.path); err != nil || saved.Entity != "climate.living_room" {
		t.Fatalf("saved %+v %v", saved, err)
	}
	if !tt.HasText(locales.Power()) || a.trayTitle() != "" {
		t.Fatalf("detail: %v %q", tt.Texts(), a.trayTitle())
	}
	tt.Key(0, ui.KeyEscape)
	if a.detail || !tt.HasText("Bedroom AC") {
		t.Fatalf("escape: %v", tt.Texts())
	}
}

// A list read keeps a local change to the device opened last, forgets that
// device once it is gone and ignores another connection.
func TestListRead(t *testing.T) {
	a := demoApp(t)
	client := ha.NewClient("http://ha.test", "token")
	a.client, a.detail = client, true
	a.panelOpen.Store(true)
	read := a.devices
	edits := a.edits.Load()
	a.climate.Target = 20
	a.edits.Add(1)
	a.applyList(client, edits, read, "", nil)
	if a.climate.Target != 20 || a.shown(read[0]).Target != 20 {
		t.Fatal("stale read replaced a local change")
	}
	a.applyList(client, a.edits.Load(), read[1:], "", nil)
	if a.hasClimate || a.detail || len(a.devices) != 1 || !a.listing || a.config.Entity != "" || a.entity != "" {
		t.Fatalf("gone device: %v %v %v %q", a.hasClimate, a.detail, a.listing, a.config.Entity)
	}
	if saved, err := loadConfig(a.path); err != nil || saved.Entity != "" {
		t.Fatalf("saved %+v %v", saved, err)
	}
	a.applyList(ha.NewClient("http://ha.test", "old"), a.edits.Load(), read, "", nil)
	if len(a.devices) != 1 {
		t.Fatal("read of another connection shown")
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
	a.openDevice(a.devices[1])
	select {
	case fn := <-posted:
		fn()
	case <-time.After(5 * time.Second):
		t.Fatal("no service call")
	}
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

func TestSetupView(t *testing.T) {
	a := newApplication(filepath.Join(t.TempDir(), "config.json"), memorySecrets{})
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
		t.Fatal("reset to the list before setup")
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
	if err := tt.Click(locales.CloseSettings()); err != nil {
		t.Fatal(err)
	}
	if a.settingsOpen || !tt.HasText("Living Room AC") {
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

func TestConfigRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x", "config.json")
	want := config{Address: "http://ha.local:8123", Entity: "climate.a"}
	if err := saveConfig(path, want); err != nil {
		t.Fatal(err)
	}
	if got, err := loadConfig(path); err != nil || got != want {
		t.Fatalf("%+v %v", got, err)
	}
}

func TestTrayIcon(t *testing.T) {
	img, err := png.Decode(bytesReader(trayIcon()))
	if err != nil || img.Bounds().Dx() != 36 {
		t.Fatal(err)
	}
}

// TestSnapshot writes the panel's pages to $HA_SNAPSHOT_DIR for review.
func TestSnapshot(t *testing.T) {
	dir := os.Getenv("HA_SNAPSHOT_DIR")
	if dir == "" {
		t.Skip("set HA_SNAPSHOT_DIR to write snapshots")
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
	write("list.png")
	a.settingsOpen = true
	write("settings.png")
	a.settingsOpen, a.detail = false, true
	write("panel.png")
	a.climate.Mode = "heat"
	a.climate.Preset = "sleep"
	write("panel-heat.png")
}

func bytesReader(b []byte) *bytes.Reader { return bytes.NewReader(b) }
