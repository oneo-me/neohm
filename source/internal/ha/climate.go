package ha

import "slices"

// State is an entity state as Home Assistant reports it.
type State struct {
	EntityID   string         `json:"entity_id"`
	State      string         `json:"state"`
	Attributes map[string]any `json:"attributes"`
}

// Climate is the part of a climate entity the panel shows.
type Climate struct {
	EntityID string
	Name     string
	// Mode is the HVAC mode: "off", "cool", "heat", "heat_cool", "auto",
	// "dry" or "fan_only"; "unavailable" when the device is gone.
	Mode  string
	Modes []string
	// Current and Target are known when HasCurrent and HasTarget say so.
	Current float64
	Target  float64
	Min     float64
	Max     float64
	Step    float64

	HasCurrent, HasTarget bool

	FanMode, SwingMode, Preset string
	FanModes, SwingModes       []string
	Presets                    []string
}

// ParseClimate reads a climate entity's state, filling Home Assistant's
// defaults for missing limits.
func ParseClimate(s State) Climate {
	a := s.Attributes
	c := Climate{
		EntityID: s.EntityID, Mode: s.State,
		Name:  str(a, "friendly_name"),
		Modes: strs(a, "hvac_modes"), FanModes: strs(a, "fan_modes"),
		SwingModes: strs(a, "swing_modes"), Presets: strs(a, "preset_modes"),
		FanMode: str(a, "fan_mode"), SwingMode: str(a, "swing_mode"), Preset: str(a, "preset_mode"),
		Min: 7, Max: 35, Step: 1,
	}
	if c.Name == "" {
		c.Name = s.EntityID
	}
	if v, ok := num(a, "min_temp"); ok {
		c.Min = v
	}
	if v, ok := num(a, "max_temp"); ok && v > c.Min {
		c.Max = v
	}
	if v, ok := num(a, "target_temp_step"); ok && v > 0 {
		c.Step = v
	}
	c.Current, c.HasCurrent = num(a, "current_temperature")
	c.Target, c.HasTarget = num(a, "temperature")
	return c
}

// On reports whether the unit runs.
func (c Climate) On() bool { return c.Mode != "off" && c.Available() }

// Available reports whether Home Assistant can reach the unit.
func (c Climate) Available() bool {
	return c.Mode != "" && c.Mode != "unavailable" && c.Mode != "unknown"
}

// HasPreset reports whether the unit offers preset p, as "eco" or "sleep".
func (c Climate) HasPreset(p string) bool { return slices.Contains(c.Presets, p) }

// NeutralPreset is the preset that turns eco and sleep off: "none" when the
// unit offers it, else the first other than eco and sleep.
func (c Climate) NeutralPreset() string {
	if c.HasPreset("none") {
		return "none"
	}
	for _, p := range c.Presets {
		if p != "eco" && p != "sleep" {
			return p
		}
	}
	return ""
}

// OnMode is the mode turning the unit on goes to: cool if offered, else
// the first mode other than off.
func (c Climate) OnMode() string {
	if slices.Contains(c.Modes, "cool") {
		return "cool"
	}
	for _, m := range c.Modes {
		if m != "off" {
			return m
		}
	}
	return ""
}

// ActiveModes are the modes other than off, which the mode menu offers.
func (c Climate) ActiveModes() []string {
	return slices.DeleteFunc(slices.Clone(c.Modes), func(m string) bool { return m == "off" })
}

// Clamp holds t within the unit's limits, on its step.
func (c Climate) Clamp(t float64) float64 {
	if c.Step > 0 {
		n := (t - c.Min) / c.Step
		t = c.Min + float64(int64(n+0.5))*c.Step
	}
	return max(c.Min, min(c.Max, t))
}

func str(a map[string]any, k string) string {
	s, _ := a[k].(string)
	return s
}

func num(a map[string]any, k string) (float64, bool) {
	v, ok := a[k].(float64)
	return v, ok
}

func strs(a map[string]any, k string) []string {
	list, _ := a[k].([]any)
	var out []string
	for _, v := range list {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
