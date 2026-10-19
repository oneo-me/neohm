package ha

import (
	"slices"
	"strings"
)

// Domain is the kind of entity id names, the part before its dot, as
// "light".
func Domain(id string) string {
	d, _, _ := strings.Cut(id, ".")
	return d
}

// Domain is the kind of entity s is, as "light".
func (s State) Domain() string { return Domain(s.EntityID) }

// Name is the entity's friendly name, or its ID when it has none.
func (s State) Name() string {
	if n := str(s.Attributes, "friendly_name"); n != "" {
		return n
	}
	return s.EntityID
}

// Available reports whether Home Assistant reaches the entity. Its state may
// still be "unknown", as a cover's that does not report its position.
func (s State) Available() bool { return s.State != "" && s.State != "unavailable" }

// Active reports whether the entity is on, open or otherwise at work, as
// Home Assistant's dashboards tell it: an unlocked lock and an armed alarm
// are active, a docked vacuum is not.
func (s State) Active() bool {
	if !s.Available() || s.State == "unknown" || s.State == "off" {
		return false
	}
	switch s.Domain() {
	case "alarm_control_panel":
		return s.State != "disarmed"
	case "cover", "valve":
		return s.State != "closed"
	case "lock":
		return s.State != "locked"
	case "media_player":
		return s.State != "standby"
	case "vacuum":
		return !slices.Contains([]string{"idle", "docked", "paused"}, s.State)
	case "lawn_mower":
		return s.State == "mowing" || s.State == "error"
	case "camera":
		return s.State == "streaming"
	}
	return true
}

// Str returns the string attribute k, or "".
func (s State) Str(k string) string { return str(s.Attributes, k) }

// Num returns the number attribute k, and whether there is one.
func (s State) Num(k string) (float64, bool) { return num(s.Attributes, k) }

// Strs returns the list of strings attribute k.
func (s State) Strs(k string) []string { return strs(s.Attributes, k) }

// Supports reports whether the entity offers feature, a bit of its
// supported_features.
func (s State) Supports(feature int) bool {
	v, _ := num(s.Attributes, "supported_features")
	return int(v)&feature != 0
}

// Action is something the panel can ask of an entity: a service of its
// domain, and the state it leads to.
type Action struct {
	// Name says what the action does, as "turn_on" or "lock", or is the
	// option it chooses when Option says so, as a select's.
	Name   string
	Option bool
	// Service is in the entity's domain, as "turn_on"; Data are its fields
	// besides the entity.
	Service string
	Data    map[string]any
	// Next is the state to show until Home Assistant reports one; empty
	// keeps the state read.
	Next string
	// Current reports that the entity already is where the action leads.
	Current bool
	// Disabled reports an action the panel cannot ask for, as one that
	// needs a code.
	Disabled bool
}

// Features of the entities, the bits of their supported_features the panel
// reads.
const (
	coverOpen, coverClose, coverStop = 1, 2, 8
	lockOpen                         = 1

	mediaPause, mediaPrevious, mediaNext = 1, 16, 32
	mediaTurnOn, mediaTurnOff, mediaPlay = 128, 256, 16384

	vacuumPause, vacuumStop, vacuumReturn = 4, 8, 16
	vacuumLocate, vacuumStart             = 512, 8192

	mowerStart, mowerPause, mowerDock = 1, 2, 4

	alarmHome, alarmAway, alarmNight = 1, 2, 4
	alarmBypass, alarmVacation       = 16, 32

	waterOperation, waterOnOff = 2, 8
)

// Actions are what the panel offers for the entity, in the order it lists
// them: the one that switches it on or off, else its commands or options.
// Climate entities have controls of their own.
func (s State) Actions() []Action {
	on := s.Active()
	switch s.Domain() {
	case "light", "switch", "fan", "humidifier":
		if on {
			return []Action{{Name: "turn_off", Service: "turn_off", Next: "off"}}
		}
		return []Action{{Name: "turn_on", Service: "turn_on", Next: "on"}}
	case "siren":
		return []Action{
			{Name: "turn_on", Service: "turn_on", Next: "on", Current: on},
			{Name: "turn_off", Service: "turn_off", Next: "off", Current: s.State == "off"},
		}
	case "button":
		return []Action{{Name: "press", Service: "press"}}
	case "cover":
		return s.openClose("cover")
	case "valve":
		return s.openClose("valve")
	case "lock":
		return s.lockActions()
	case "media_player":
		return s.mediaActions()
	case "vacuum":
		return s.vacuumActions()
	case "lawn_mower":
		return s.mowerActions()
	case "alarm_control_panel":
		return s.alarmActions()
	case "water_heater":
		return s.waterActions()
	case "select":
		var out []Action
		for _, o := range s.Strs("options") {
			out = append(out, Action{Name: o, Option: true, Service: "select_option",
				Data: map[string]any{"option": o}, Next: o, Current: s.State == o})
		}
		return out
	}
	return nil
}

// openClose opens, stops and closes a cover or a valve, the domain whose
// services end in suffix.
func (s State) openClose(suffix string) []Action {
	var out []Action
	if s.Supports(coverOpen) {
		out = append(out, Action{Name: "open", Service: "open_" + suffix, Next: "opening", Current: s.State == "open"})
	}
	if s.Supports(coverStop) {
		out = append(out, Action{Name: "stop", Service: "stop_" + suffix})
	}
	if s.Supports(coverClose) {
		out = append(out, Action{Name: "close", Service: "close_" + suffix, Next: "closing", Current: s.State == "closed"})
	}
	return out
}

// needsCode reports whether the entity asks for a code, which the panel has
// no way to enter.
func (s State) needsCode() bool { return s.Attributes["code_format"] != nil }

func (s State) lockActions() []Action {
	code := s.needsCode()
	out := []Action{
		{Name: "lock", Service: "lock", Next: "locking", Current: s.State == "locked", Disabled: code},
		{Name: "unlock", Service: "unlock", Next: "unlocking", Current: s.State == "unlocked", Disabled: code},
	}
	if s.Supports(lockOpen) {
		out = append(out, Action{Name: "open_door", Service: "open", Next: "opening", Current: s.State == "open", Disabled: code})
	}
	return out
}

func (s State) mediaActions() []Action {
	if s.State == "off" || s.State == "standby" {
		if s.Supports(mediaTurnOn) {
			return []Action{{Name: "turn_on", Service: "turn_on", Next: "on"}}
		}
		return nil
	}
	var out []Action
	switch {
	case s.State == "playing" && s.Supports(mediaPause):
		out = append(out, Action{Name: "pause", Service: "media_pause", Next: "paused"})
	case s.State != "playing" && s.Supports(mediaPlay):
		out = append(out, Action{Name: "play", Service: "media_play", Next: "playing"})
	}
	if s.Supports(mediaPrevious) {
		out = append(out, Action{Name: "previous", Service: "media_previous_track"})
	}
	if s.Supports(mediaNext) {
		out = append(out, Action{Name: "next", Service: "media_next_track"})
	}
	if s.Supports(mediaTurnOff) {
		out = append(out, Action{Name: "turn_off", Service: "turn_off", Next: "off"})
	}
	return out
}

func (s State) vacuumActions() []Action {
	var out []Action
	busy := s.State == "cleaning" || s.State == "returning"
	if s.Supports(vacuumStart) && s.State != "cleaning" {
		out = append(out, Action{Name: "start", Service: "start", Next: "cleaning"})
	}
	if s.Supports(vacuumPause) && s.State == "cleaning" {
		out = append(out, Action{Name: "pause", Service: "pause", Next: "paused"})
	}
	if s.Supports(vacuumStop) && busy {
		out = append(out, Action{Name: "stop", Service: "stop", Next: "idle"})
	}
	if s.Supports(vacuumReturn) && s.State != "docked" && s.State != "returning" {
		out = append(out, Action{Name: "return", Service: "return_to_base", Next: "returning"})
	}
	if s.Supports(vacuumLocate) {
		out = append(out, Action{Name: "locate", Service: "locate"})
	}
	return out
}

func (s State) mowerActions() []Action {
	var out []Action
	if s.Supports(mowerStart) && s.State != "mowing" {
		out = append(out, Action{Name: "start", Service: "start_mowing", Next: "mowing"})
	}
	if s.Supports(mowerPause) && s.State == "mowing" {
		out = append(out, Action{Name: "pause", Service: "pause", Next: "paused"})
	}
	if s.Supports(mowerDock) && s.State != "docked" && s.State != "returning" {
		out = append(out, Action{Name: "dock", Service: "dock", Next: "returning"})
	}
	return out
}

func (s State) alarmActions() []Action {
	code := s.needsCode()
	// Arming asks for the code too unless the panel says it does not.
	armCode := code
	if required, ok := s.Attributes["code_arm_required"].(bool); ok && !required {
		armCode = false
	}
	out := []Action{{Name: "disarm", Service: "alarm_disarm", Next: "disarmed", Current: s.State == "disarmed", Disabled: code}}
	for _, m := range []struct {
		feature int
		mode    string
	}{{alarmHome, "home"}, {alarmAway, "away"}, {alarmNight, "night"}, {alarmVacation, "vacation"}, {alarmBypass, "custom_bypass"}} {
		if s.Supports(m.feature) {
			out = append(out, Action{Name: "arm_" + m.mode, Service: "alarm_arm_" + m.mode, Next: "arming",
				Current: s.State == "armed_"+m.mode, Disabled: armCode})
		}
	}
	return out
}

// waterActions choose a water heater's operation, which may be "off", and
// switch it on and off when it has a switch of its own.
func (s State) waterActions() []Action {
	var out []Action
	modes := s.Strs("operation_list")
	if s.Supports(waterOperation) {
		for _, m := range modes {
			out = append(out, Action{Name: m, Option: true, Service: "set_operation_mode",
				Data: map[string]any{"operation_mode": m}, Next: m, Current: s.State == m})
		}
	}
	if s.Supports(waterOnOff) && !slices.Contains(modes, "off") {
		if s.Active() {
			out = append(out, Action{Name: "turn_off", Service: "turn_off", Next: "off"})
		} else {
			out = append(out, Action{Name: "turn_on", Service: "turn_on"})
		}
	}
	return out
}
