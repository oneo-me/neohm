package ha

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const climateJSON = `{"entity_id":"climate.bedroom","state":"cool","attributes":{
	"friendly_name":"Bedroom AC","hvac_modes":["off","cool","heat","dry","fan_only","auto"],
	"min_temp":16,"max_temp":30,"target_temp_step":0.5,"current_temperature":24,"temperature":16,
	"fan_modes":["auto","low","medium","high"],"fan_mode":"low",
	"swing_modes":["on","off"],"swing_mode":"on",
	"preset_modes":["none","eco","sleep"],"preset_mode":"none"}}`

func TestClient(t *testing.T) {
	var calls []map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/api/states/climate.bedroom":
			w.Write([]byte(climateJSON))
		case "/api/states":
			w.Write([]byte(`[{"entity_id":"light.desk","state":"on","attributes":{}},` + climateJSON + `]`))
		case "/api/config":
			w.Write([]byte(`{"unit_system":{"temperature":"°C"}}`))
		case "/api/services/climate/set_temperature":
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			calls = append(calls, body)
			w.Write([]byte(`[]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	ctx := context.Background()
	c := NewClient(srv.URL+"/ ", "secret")

	got, err := c.Climate(ctx, "climate.bedroom")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Bedroom AC" || got.Mode != "cool" || got.Target != 16 || got.Current != 24 || got.Step != 0.5 || !got.HasCurrent {
		t.Fatalf("parsed %+v", got)
	}
	if !got.On() || got.OnMode() != "cool" || got.NeutralPreset() != "none" || len(got.ActiveModes()) != 5 {
		t.Fatalf("helpers: %+v", got)
	}
	if got.Clamp(10) != 16 || got.Clamp(40) != 30 || got.Clamp(21.3) != 21.5 {
		t.Fatal("clamp")
	}
	list, err := c.Climates(ctx)
	if err != nil || len(list) != 1 || list[0].EntityID != "climate.bedroom" {
		t.Fatalf("climates: %v %v", list, err)
	}
	if u, err := c.TemperatureUnit(ctx); err != nil || u != "°C" {
		t.Fatalf("unit %q %v", u, err)
	}
	if err := c.Call(ctx, "set_temperature", "climate.bedroom", map[string]any{"temperature": 21.5}); err != nil {
		t.Fatal(err)
	}
	if len(calls) != 1 || calls[0]["entity_id"] != "climate.bedroom" || calls[0]["temperature"] != 21.5 {
		t.Fatalf("calls %v", calls)
	}
	if _, err := NewClient(srv.URL, "wrong").Climate(ctx, "climate.bedroom"); err != ErrUnauthorized {
		t.Fatalf("bad token: %v", err)
	}
}

func TestNormalizeURL(t *testing.T) {
	for in, want := range map[string]string{
		" homeassistant.local:8123/ ": "http://homeassistant.local:8123",
		"https://ha.example.com//":    "https://ha.example.com",
		"":                            "",
	} {
		if got := NormalizeURL(in); got != want {
			t.Errorf("%q: %q", in, got)
		}
	}
}
