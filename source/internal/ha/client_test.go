package ha

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

const climateJSON = `{"entity_id":"climate.bedroom","state":"cool","attributes":{
	"friendly_name":"Bedroom AC","hvac_modes":["off","cool","heat","dry","fan_only","auto"],
	"min_temp":16,"max_temp":30,"target_temp_step":0.5,"current_temperature":24,"temperature":16,
	"fan_modes":["auto","low","medium","high"],"fan_mode":"low",
	"swing_modes":["on","off"],"swing_mode":"on",
	"preset_modes":["none","eco","sleep"],"preset_mode":"none"}}`

func TestClient(t *testing.T) {
	var calls []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch {
		case r.URL.Path == "/api/states/climate.bedroom":
			w.Write([]byte(climateJSON))
		case r.URL.Path == "/api/states":
			w.Write([]byte(`[{"entity_id":"light.desk","state":"on","attributes":{}},` + climateJSON + `]`))
		case r.URL.Path == "/api/config":
			w.Write([]byte(`{"unit_system":{"temperature":"°C"}}`))
		case strings.HasPrefix(r.URL.Path, "/api/services/"):
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			calls = append(calls, r.URL.Path+" "+body["entity_id"].(string))
			w.Write([]byte(`[]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	ctx := context.Background()
	c := NewClient(srv.URL+"/ ", "secret")

	s, err := c.State(ctx, "climate.bedroom")
	if err != nil {
		t.Fatal(err)
	}
	got := ParseClimate(s)
	if got.Name != "Bedroom AC" || got.Mode != "cool" || got.Target != 16 || got.Current != 24 || got.Step != 0.5 || !got.HasCurrent {
		t.Fatalf("parsed %+v", got)
	}
	if !got.On() || got.OnMode() != "cool" || got.NeutralPreset() != "none" || len(got.ActiveModes()) != 5 {
		t.Fatalf("helpers: %+v", got)
	}
	if got.Clamp(10) != 16 || got.Clamp(40) != 30 || got.Clamp(21.3) != 21.5 {
		t.Fatal("clamp")
	}
	states, err := c.States(ctx)
	if err != nil || len(states) != 2 || states[0].Domain() != "light" || states[1].Name() != "Bedroom AC" {
		t.Fatalf("states: %v %v", states, err)
	}
	if u, err := c.TemperatureUnit(ctx); err != nil || u != "°C" {
		t.Fatalf("unit %q %v", u, err)
	}
	if err := c.Call(ctx, "climate", "set_temperature", "climate.bedroom", map[string]any{"temperature": 21.5}); err != nil {
		t.Fatal(err)
	}
	if err := c.Call(ctx, "light", "turn_off", "light.desk", nil); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(calls, []string{"/api/services/climate/set_temperature climate.bedroom", "/api/services/light/turn_off light.desk"}) {
		t.Fatalf("calls %v", calls)
	}
	if _, err := NewClient(srv.URL, "wrong").State(ctx, "climate.bedroom"); err != ErrUnauthorized {
		t.Fatalf("bad token: %v", err)
	}
}

// frame is a frame as a server sends it, unmasked.
func frame(fin bool, op byte, payload []byte) []byte {
	b := byte(op)
	if fin {
		b |= 0x80
	}
	out := []byte{b}
	switch n := len(payload); {
	case n < 126:
		out = append(out, byte(n))
	case n <= 0xFFFF:
		out = binary.BigEndian.AppendUint16(append(out, 126), uint16(n))
	default:
		out = binary.BigEndian.AppendUint64(append(out, 127), uint64(n))
	}
	return append(out, payload...)
}

// registryServer answers the WebSocket API as Home Assistant does: it asks
// for the token, then answers the three registry commands out of order,
// pinging, fragmenting a message and sending one too long for a short
// frame on the way.
func registryServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/websocket" || !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
			http.NotFound(w, r)
			return
		}
		conn, brw, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		brw.WriteString("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n" +
			"Sec-WebSocket-Accept: " + acceptKey(r.Header.Get("Sec-WebSocket-Key")) + "\r\n\r\n")
		brw.Flush()
		ws := &webSocket{r: bufio.NewReader(brw), conn: conn}
		ws.writeJSON(map[string]any{"type": "auth_required", "ha_version": "2026.10.0"})
		var auth struct {
			Type        string
			AccessToken string `json:"access_token"`
		}
		if err := ws.readJSON(&auth); err != nil || auth.Type != "auth" {
			t.Errorf("auth: %+v %v", auth, err)
			return
		}
		if auth.AccessToken != "secret" {
			ws.writeJSON(map[string]any{"type": "auth_invalid", "message": "Invalid access token"})
			return
		}
		ws.writeJSON(map[string]any{"type": "auth_ok"})
		ids := map[string]int{}
		for range 3 {
			var cmd struct {
				ID   int
				Type string
			}
			if err := ws.readJSON(&cmd); err != nil {
				t.Error(err)
				return
			}
			ids[cmd.Type] = cmd.ID
		}
		result := func(cmd, body string) []byte {
			b, _ := json.Marshal(map[string]any{"id": ids[cmd], "type": "result", "success": true, "result": json.RawMessage(body)})
			return b
		}
		conn.Write(frame(true, opPing, []byte("hi")))
		ws.writeJSON(map[string]any{"type": "event", "event": map[string]any{}})
		areas := result("config/area_registry/list", `[{"area_id":"office","name":"Office"},{"area_id":"kitchen","name":"Kitchen"}]`)
		conn.Write(append(frame(false, opText, areas[:10]), frame(true, opContinue, areas[10:])...))
		devices := result("config/device_registry/list", `[
			{"id":"d1","name":"Desk","name_by_user":"My desk","area_id":"office","entry_type":null,"padding":"`+strings.Repeat("x", 70000)+`"},
			{"id":"d2","name":"Sun","area_id":null,"entry_type":"service"}]`)
		ws.write(opText, devices)
		ws.write(opText, result("config/entity_registry/list_for_display", `{"entity_categories":{"0":"config","1":"diagnostic"},"entities":[
			{"ei":"light.desk","di":"d1","pl":"hue"},
			{"ei":"switch.desk_led","di":"d1","pl":"hue","ec":0},
			{"ei":"sensor.sun_elevation","di":"d2","pl":"sun","dp":1},
			{"ei":"light.spare","pl":"hue","hb":true},
			{"ei":"light.kitchen","ai":"kitchen","pl":"template"}]}`))
		// The client answers the ping before it reads what follows it.
		ws.read()
	}))
}

func TestRegistry(t *testing.T) {
	srv := registryServer(t)
	defer srv.Close()
	ctx := context.Background()
	reg, err := NewClient(srv.URL, "secret").Registry(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]RegistryEntry{
		"light.desk":           {Device: "d1"},
		"switch.desk_led":      {Device: "d1", Hidden: true},
		"sensor.sun_elevation": {Device: "d2", Precision: 1, HasPrecision: true},
		"light.spare":          {Hidden: true},
		"light.kitchen":        {Area: "kitchen"},
	}
	if len(reg.Entities) != len(want) {
		t.Fatalf("entities %+v", reg.Entities)
	}
	for id, e := range want {
		if reg.Entities[id] != e {
			t.Errorf("%s: %+v", id, reg.Entities[id])
		}
	}
	if reg.Devices["d1"] != (Device{Name: "My desk", Area: "office"}) || !reg.Devices["d2"].Service {
		t.Fatalf("devices %+v", reg.Devices)
	}
	if !slices.Equal(reg.Areas, []Area{{"office", "Office"}, {"kitchen", "Kitchen"}}) {
		t.Fatalf("areas %+v", reg.Areas)
	}
	if p, ok := reg.Precision("sensor.sun_elevation"); !ok || p != 1 {
		t.Fatal("precision")
	}
	if _, ok := (*Registry)(nil).Precision("x"); ok {
		t.Fatal("precision without a registry")
	}
	if _, err := NewClient(srv.URL, "wrong").Registry(ctx); err != ErrUnauthorized {
		t.Fatalf("bad token: %v", err)
	}
}

func parseStates(t *testing.T, list string) []State {
	t.Helper()
	var states []State
	if err := json.Unmarshal([]byte(list), &states); err != nil {
		t.Fatal(err)
	}
	return states
}

func TestArrange(t *testing.T) {
	states := parseStates(t, `[
		{"entity_id":"automation.morning","state":"on","attributes":{}},
		{"entity_id":"switch.strip_2","state":"on","attributes":{"friendly_name":"Strip 2"}},
		{"entity_id":"switch.strip_1","state":"off","attributes":{"friendly_name":"Strip 1"}},
		{"entity_id":"sensor.strip_power","state":"12","attributes":{"device_class":"power"}},
		{"entity_id":"button.strip_restart","state":"unknown","attributes":{}},
		{"entity_id":"sensor.ac_outdoor","state":"31","attributes":{"device_class":"temperature"}},
		{"entity_id":"switch.ac_display","state":"on","attributes":{}},
		{"entity_id":"climate.ac","state":"cool","attributes":{"friendly_name":"Bedroom AC"}},
		{"entity_id":"light.lamp","state":"off","attributes":{"friendly_name":"Bedroom Lamp"}},
		{"entity_id":"switch.lamp_led","state":"on","attributes":{}},
		{"entity_id":"sensor.th_battery","state":"80","attributes":{"device_class":"battery"}},
		{"entity_id":"sensor.th_humidity","state":"45","attributes":{"device_class":"humidity"}},
		{"entity_id":"sensor.th_temperature","state":"24","attributes":{"device_class":"temperature"}},
		{"entity_id":"binary_sensor.door","state":"on","attributes":{"device_class":"door"}},
		{"entity_id":"sensor.phone_steps","state":"5000","attributes":{}},
		{"entity_id":"sensor.sun_elevation","state":"30","attributes":{"unit_of_measurement":"°C"}},
		{"entity_id":"sensor.average","state":"23","attributes":{"device_class":"temperature"}},
		{"entity_id":"light.spare","state":"off","attributes":{}},
		{"entity_id":"light.group","state":"on","attributes":{"friendly_name":"All lights"}},
		{"entity_id":"button.bell","state":"unknown","attributes":{"friendly_name":"Chime"}}
	]`)
	reg := &Registry{
		Entities: map[string]RegistryEntry{
			"switch.strip_1": {Device: "strip"}, "switch.strip_2": {Device: "strip"},
			"sensor.strip_power": {Device: "strip"}, "button.strip_restart": {Device: "strip"},
			"sensor.ac_outdoor": {Device: "ac"}, "switch.ac_display": {Device: "ac"}, "climate.ac": {Device: "ac"},
			"light.lamp": {Device: "lamp"}, "switch.lamp_led": {Device: "lamp", Hidden: true},
			"sensor.th_battery": {Device: "th"}, "sensor.th_humidity": {Device: "th"}, "sensor.th_temperature": {Device: "th"},
			"binary_sensor.door": {Device: "door", Area: "hall"},
			"sensor.phone_steps": {Device: "phone"}, "sensor.sun_elevation": {Device: "sun"},
			"light.spare": {Hidden: true}, "button.bell": {Device: "bell"},
		},
		Devices: map[string]Device{
			"strip": {Name: "Power strip", Area: "living"}, "ac": {Name: "Bedroom AC", Area: "bedroom"},
			"lamp": {Name: "Lamp", Area: "bedroom"}, "th": {Name: "Thermometer", Area: "bedroom"},
			"door": {Name: "Front door", Area: "living"}, "phone": {Name: "Phone"},
			"sun": {Name: "Sun", Service: true}, "bell": {Name: "Bell"},
		},
		Areas: []Area{{"living", "Living Room"}, {"bedroom", "Bedroom"}, {"hall", "Hall"}, {"attic", "Attic"}},
	}
	type tile struct{ id, name, readings string }
	want := map[string][]tile{
		"Living Room": {{"switch.strip_2", "Strip 2", ""}, {"switch.strip_1", "Strip 1", ""}},
		"Bedroom": {{"climate.ac", "Bedroom AC", ""}, {"light.lamp", "Bedroom Lamp", ""},
			{"sensor.th_temperature", "Thermometer", "sensor.th_humidity"}},
		"Hall": {{"binary_sensor.door", "Front door", ""}},
		"":     {{"light.group", "All lights", ""}, {"button.bell", "Chime", ""}},
	}
	sections := Arrange(states, reg)
	var areas []string
	for _, s := range sections {
		areas = append(areas, s.Area)
		var got []tile
		for _, x := range s.Tiles {
			var readings []string
			for _, r := range x.Readings {
				readings = append(readings, r.EntityID)
			}
			got = append(got, tile{x.Entity.EntityID, x.Name, strings.Join(readings, ",")})
		}
		if !slices.Equal(got, want[s.Area]) {
			t.Errorf("%q: %v", s.Area, got)
		}
	}
	if !slices.Equal(areas, []string{"Living Room", "Bedroom", "Hall", ""}) {
		t.Errorf("areas %q", areas)
	}

	// Without the registries, every switch and main control shows alone.
	sections = Arrange(states, nil)
	var ids []string
	for _, x := range sections[0].Tiles {
		ids = append(ids, x.Entity.EntityID)
	}
	if len(sections) != 1 || !slices.Equal(ids, []string{"climate.ac", "light.group", "light.lamp", "light.spare",
		"switch.strip_1", "switch.strip_2", "switch.ac_display", "switch.lamp_led"}) {
		t.Errorf("no registry: %v", ids)
	}
}

func TestActions(t *testing.T) {
	type action struct {
		name, service, next      string
		current, disabled, other bool
	}
	for _, tc := range []struct {
		state  string
		active bool
		want   []action
	}{
		{`{"entity_id":"light.a","state":"on"}`, true, []action{{"turn_off", "turn_off", "off", false, false, false}}},
		{`{"entity_id":"fan.a","state":"off"}`, false, []action{{"turn_on", "turn_on", "on", false, false, false}}},
		{`{"entity_id":"light.a","state":"unavailable"}`, false, []action{{"turn_on", "turn_on", "on", false, false, false}}},
		{`{"entity_id":"cover.a","state":"open","attributes":{"supported_features":15}}`, true, []action{
			{"open", "open_cover", "opening", true, false, false}, {"stop", "stop_cover", "", false, false, false},
			{"close", "close_cover", "closing", false, false, false}}},
		{`{"entity_id":"cover.a","state":"closed","attributes":{"supported_features":3}}`, false, []action{
			{"open", "open_cover", "opening", false, false, false}, {"close", "close_cover", "closing", true, false, false}}},
		{`{"entity_id":"valve.a","state":"closed","attributes":{"supported_features":3}}`, false, []action{
			{"open", "open_valve", "opening", false, false, false}, {"close", "close_valve", "closing", true, false, false}}},
		{`{"entity_id":"lock.a","state":"locked","attributes":{"supported_features":1}}`, false, []action{
			{"lock", "lock", "locking", true, false, false}, {"unlock", "unlock", "unlocking", false, false, false},
			{"open_door", "open", "opening", false, false, false}}},
		{`{"entity_id":"lock.a","state":"unlocked","attributes":{"code_format":"^\\d{4}$"}}`, true, []action{
			{"lock", "lock", "locking", false, true, false}, {"unlock", "unlock", "unlocking", true, true, false}}},
		{`{"entity_id":"media_player.a","state":"playing","attributes":{"supported_features":16817}}`, true, []action{
			{"pause", "media_pause", "paused", false, false, false}, {"previous", "media_previous_track", "", false, false, false},
			{"next", "media_next_track", "", false, false, false}, {"turn_off", "turn_off", "off", false, false, false}}},
		{`{"entity_id":"media_player.a","state":"standby","attributes":{"supported_features":16817}}`, false, []action{
			{"turn_on", "turn_on", "on", false, false, false}}},
		{`{"entity_id":"vacuum.a","state":"docked","attributes":{"supported_features":12316}}`, false, []action{
			{"start", "start", "cleaning", false, false, false}}},
		{`{"entity_id":"vacuum.a","state":"cleaning","attributes":{"supported_features":12828}}`, true, []action{
			{"pause", "pause", "paused", false, false, false}, {"stop", "stop", "idle", false, false, false},
			{"return", "return_to_base", "returning", false, false, false}, {"locate", "locate", "", false, false, false}}},
		{`{"entity_id":"lawn_mower.a","state":"mowing","attributes":{"supported_features":7}}`, true, []action{
			{"pause", "pause", "paused", false, false, false}, {"dock", "dock", "returning", false, false, false}}},
		{`{"entity_id":"alarm_control_panel.a","state":"armed_away","attributes":{"supported_features":3,"code_format":"number","code_arm_required":false}}`, true, []action{
			{"disarm", "alarm_disarm", "disarmed", false, true, false}, {"arm_home", "alarm_arm_home", "arming", false, false, false},
			{"arm_away", "alarm_arm_away", "arming", true, false, false}}},
		{`{"entity_id":"water_heater.a","state":"eco","attributes":{"supported_features":10,"operation_list":["eco","electric"]}}`, true, []action{
			{"eco", "set_operation_mode", "eco", true, false, true}, {"electric", "set_operation_mode", "electric", false, false, true},
			{"turn_off", "turn_off", "off", false, false, false}}},
		{`{"entity_id":"select.a","state":"low","attributes":{"options":["low","high"]}}`, true, []action{
			{"low", "select_option", "low", true, false, true}, {"high", "select_option", "high", false, false, true}}},
		{`{"entity_id":"siren.a","state":"off"}`, false, []action{
			{"turn_on", "turn_on", "on", false, false, false}, {"turn_off", "turn_off", "off", true, false, false}}},
		{`{"entity_id":"button.a","state":"2026-10-19T08:00:00+00:00"}`, true, []action{{"press", "press", "", false, false, false}}},
		{`{"entity_id":"camera.a","state":"streaming"}`, true, nil},
		{`{"entity_id":"climate.a","state":"cool"}`, true, nil},
		{`{"entity_id":"binary_sensor.a","state":"off"}`, false, nil},
	} {
		var s State
		if err := json.Unmarshal([]byte(tc.state), &s); err != nil {
			t.Fatal(err)
		}
		if s.Active() != tc.active {
			t.Errorf("%s: active %v", tc.state, s.Active())
		}
		var got []action
		for _, a := range s.Actions() {
			got = append(got, action{a.Name, a.Service, a.Next, a.Current, a.Disabled, a.Option})
		}
		if !slices.Equal(got, tc.want) {
			t.Errorf("%s: %v", tc.state, got)
		}
	}
	for state, active := range map[string]bool{"locked": false, "jammed": true, "unknown": false} {
		if (State{EntityID: "lock.a", State: state}).Active() != active {
			t.Errorf("lock %s", state)
		}
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
