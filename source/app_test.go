package main

import (
	"bytes"
	"encoding/json"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/ha-go/source/internal/ha"
	"github.com/oneo-me/ha-go/source/locales"
	"github.com/oneo-me/onevra/source/theme"
)

type memorySecrets map[string][]byte

func (m memorySecrets) Get(k string) ([]byte, error) { return m[k], nil }
func (m memorySecrets) Set(k string, v []byte) error { m[k] = v; return nil }

// demoApp is an application showing a sample unit, without a connection:
// service calls change only what it shows.
func demoApp(t *testing.T) *application {
	var s ha.State
	if err := json.Unmarshal([]byte(`{"entity_id":"climate.bedroom","state":"cool","attributes":{
		"friendly_name":"Bedroom AC","hvac_modes":["off","cool","heat","dry","fan_only","auto"],
		"min_temp":16,"max_temp":30,"target_temp_step":1,"current_temperature":24,"temperature":16,
		"fan_modes":["auto","low","medium","high"],"fan_mode":"low",
		"swing_modes":["on","off"],"swing_mode":"on",
		"preset_modes":["none","eco","sleep"],"preset_mode":"none"}}`), &s); err != nil {
		t.Fatal(err)
	}
	a := newApplication(filepath.Join(t.TempDir(), "config.json"), memorySecrets{})
	a.post = func(fn func()) {}
	a.config = config{Address: "http://ha.test", Entity: "climate.bedroom"}
	a.token = "token"
	a.climate, a.hasClimate, a.online = ha.ParseClimate(s), true, true
	return a
}

func TestPanelControls(t *testing.T) {
	a := demoApp(t)
	tt := ui.NewTester(a.panelView, panelWidth, panelHeight)
	for _, s := range []string{locales.Room(), "24°C", "16", locales.Eco(), locales.ModeCool(), locales.FanLow(), locales.SwingOn()} {
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
}

// A read started before a local change must not undo it on screen.
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

// TestSnapshot writes the panel to $HA_SNAPSHOT_DIR for review.
func TestSnapshot(t *testing.T) {
	dir := os.Getenv("HA_SNAPSHOT_DIR")
	if dir == "" {
		t.Skip("set HA_SNAPSHOT_DIR to write snapshots")
	}
	a := demoApp(t)
	tt := ui.NewTester(a.panelView, panelWidth, panelHeight)
	tt.SetScale(2)
	tt.Advance(800 * time.Millisecond)
	write := func(name string) {
		f, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err := png.Encode(f, tt.Image()); err != nil {
			t.Fatal(err)
		}
	}
	write("panel.png")
	a.climate.Mode = "heat"
	a.climate.Preset = "sleep"
	tt.Frame()
	tt.Advance(800 * time.Millisecond)
	write("panel-heat.png")
}

func bytesReader(b []byte) *bytes.Reader { return bytes.NewReader(b) }

func TestSettingsView(t *testing.T) {
	a := newApplication(filepath.Join(t.TempDir(), "config.json"), memorySecrets{})
	a.post = func(fn func()) { fn() }
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(theme.Palette(c.Theme().Dark))
		a.settingsWindowView(c)
	}, 560, 520)
	for _, s := range []string{locales.Address(), locales.Token(), locales.Entity()} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
}
