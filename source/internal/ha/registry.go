package ha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// Registry is what Home Assistant's registries say of the entities: the
// devices and areas they belong to, and which ones its dashboards leave out.
type Registry struct {
	Entities map[string]RegistryEntry
	Devices  map[string]Device
	// Areas are in Home Assistant's order.
	Areas []Area
}

// RegistryEntry is what the entity registry says of one entity.
type RegistryEntry struct {
	// Device and Area are IDs; an entity's own area overrides its device's.
	Device, Area string
	// Hidden reports an entity the user hid, or one that configures or
	// diagnoses its device rather than doing what the device is for.
	Hidden bool
	// Precision is how many decimals a sensor shows, when HasPrecision.
	Precision    int
	HasPrecision bool
}

// Device is an entry of the device registry.
type Device struct {
	Name string
	Area string // ID
	// Service reports a service, as a weather forecast, rather than a
	// device.
	Service bool
}

// Area is an entry of the area registry.
type Area struct{ ID, Name string }

// Precision returns how many decimals sensor id shows, as its settings in
// Home Assistant say, and whether they say.
func (r *Registry) Precision(id string) (int, bool) {
	if r == nil {
		return 0, false
	}
	e := r.Entities[id]
	return e.Precision, e.HasPrecision
}

// Registry reads the entity, device and area registries, which only Home
// Assistant's WebSocket API offers.
func (c *Client) Registry(ctx context.Context) (Registry, error) {
	if c.BaseURL == "" || c.Token == "" {
		return Registry{}, errors.New("home assistant is not configured")
	}
	ws, err := c.dialWebSocket(ctx)
	if err != nil {
		return Registry{}, err
	}
	defer ws.Close()
	defer context.AfterFunc(ctx, func() { ws.Close() })()

	var msg struct{ Type, Message string }
	if err := ws.readJSON(&msg); err != nil {
		return Registry{}, err
	}
	if msg.Type != "auth_required" {
		return Registry{}, fmt.Errorf("websocket: unexpected %q", msg.Type)
	}
	if err := ws.writeJSON(map[string]any{"type": "auth", "access_token": c.Token}); err != nil {
		return Registry{}, err
	}
	if err := ws.readJSON(&msg); err != nil {
		return Registry{}, err
	}
	switch msg.Type {
	case "auth_ok":
	case "auth_invalid":
		return Registry{}, ErrUnauthorized
	default:
		return Registry{}, fmt.Errorf("websocket: unexpected %q", msg.Type)
	}
	commands := []string{"config/entity_registry/list_for_display", "config/device_registry/list", "config/area_registry/list"}
	for i, cmd := range commands {
		if err := ws.writeJSON(map[string]any{"id": i + 1, "type": cmd}); err != nil {
			return Registry{}, err
		}
	}
	results := make([]json.RawMessage, len(commands))
	for missing := len(commands); missing > 0; {
		var res struct {
			ID      int
			Type    string
			Success bool
			Result  json.RawMessage
			Error   struct{ Message string }
		}
		if err := ws.readJSON(&res); err != nil {
			return Registry{}, err
		}
		if res.Type != "result" || res.ID < 1 || res.ID > len(commands) || results[res.ID-1] != nil {
			continue
		}
		if !res.Success {
			return Registry{}, fmt.Errorf("%s: %s", commands[res.ID-1], res.Error.Message)
		}
		results[res.ID-1] = res.Result
		missing--
	}
	return parseRegistry(results[0], results[1], results[2])
}

func parseRegistry(entities, devices, areas json.RawMessage) (Registry, error) {
	var e struct {
		Entities []struct {
			ID        string `json:"ei"`
			Device    string `json:"di"`
			Area      string `json:"ai"`
			Category  *int   `json:"ec"`
			Hidden    bool   `json:"hb"`
			Precision *int   `json:"dp"`
		}
	}
	var d []struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		NameByUser string `json:"name_by_user"`
		Area       string `json:"area_id"`
		EntryType  string `json:"entry_type"`
	}
	var a []struct {
		ID   string `json:"area_id"`
		Name string `json:"name"`
	}
	if err := errors.Join(json.Unmarshal(entities, &e), json.Unmarshal(devices, &d), json.Unmarshal(areas, &a)); err != nil {
		return Registry{}, err
	}
	r := Registry{Entities: make(map[string]RegistryEntry, len(e.Entities)), Devices: make(map[string]Device, len(d))}
	for _, x := range e.Entities {
		entry := RegistryEntry{Device: x.Device, Area: x.Area, Hidden: x.Hidden || x.Category != nil}
		if x.Precision != nil {
			entry.Precision, entry.HasPrecision = *x.Precision, true
		}
		r.Entities[x.ID] = entry
	}
	for _, x := range d {
		name := x.NameByUser
		if name == "" {
			name = x.Name
		}
		r.Devices[x.ID] = Device{Name: name, Area: x.Area, Service: x.EntryType == "service"}
	}
	for _, x := range a {
		r.Areas = append(r.Areas, Area{ID: x.ID, Name: x.Name})
	}
	return r, nil
}
