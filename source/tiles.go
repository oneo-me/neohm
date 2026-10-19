package main

import (
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/ha-go/source/internal/ha"
	"github.com/oneo-me/ha-go/source/locales"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// tileKind is what a click on a tile does.
type tileKind int

const (
	tileStatic tileKind = iota // nothing: the tile only reads
	tileSwitch                 // switches the device on or off
	tilePress                  // presses it, as a button
	tileMenu                   // opens a menu of its actions
	tilePage                   // opens its controls, an air conditioner's
)

// tileView is how a tile shows its device.
type tileView struct {
	kind tileKind
	icon func() *ui.SVG
	// color tints the icon, and the tile while lit: on, open or at work.
	// A neutral device, one never on or off as a camera, keeps its icon
	// in color instead.
	color   ui.Color
	lit     bool
	neutral bool
	// status is the main state, as "On · 80%".
	status  string
	actions []ha.Action
}

// tileView shows entity s, as held, of tile t.
func (a *application) tileView(t ha.Tile, s ha.State) tileView {
	v := tileView{kind: tileMenu, icon: lucide.CircleDot, color: accent, lit: s.Active(), status: stateText(s), actions: s.Actions()}
	switch s.Domain() {
	case "climate":
		return a.climateTile(s)
	case "sensor", "binary_sensor":
		return a.sensorTile(t, s)
	case "light":
		v.kind, v.icon, v.color = tileSwitch, lucide.Lightbulb, amber
		if c, ok := lightColor(s); ok {
			v.color = c
		}
		v.status = withDetail(s, percent(s, "brightness", 255))
	case "switch":
		v.kind, v.icon, v.color = tileSwitch, lucide.Power, green
		if s.Str("device_class") == "outlet" {
			v.icon = lucide.Plug
		}
	case "fan":
		v.kind, v.icon, v.color = tileSwitch, lucide.Fan, cyan
		v.status = withDetail(s, percent(s, "percentage", 100))
	case "humidifier":
		v.kind, v.icon = tileSwitch, lucide.Droplets
		v.status = withDetail(s, percent(s, "humidity", 100))
	case "water_heater":
		v.icon, v.color = lucide.ShowerHead, orange
		if s.Available() && s.State != "unknown" {
			v.status = optionText("water", s.State)
			if target, ok := s.Num("temperature"); ok && s.Active() {
				v.status += " · " + formatTemp(target) + "°"
			}
		}
	case "cover":
		v.icon, v.color = coverIcon(s), lavender
		v.status = withDetail(s, position(s))
	case "valve":
		v.icon = lucide.Faucet
		if s.Str("device_class") == "gas" {
			v.icon, v.color = lucide.Flame, orange
		}
		v.status = withDetail(s, position(s))
	case "lock":
		v.icon, v.color = lucide.LockOpen, orange
		switch s.State {
		case "locked", "locking":
			v.icon = lucide.Lock
		case "jammed":
			v.color = danger
		}
	case "media_player":
		v.icon, v.color = lucide.Speaker, pink
		switch s.Str("device_class") {
		case "tv":
			v.icon = lucide.Tv
		case "receiver":
			v.icon = lucide.Radio
		}
		if s.State == "playing" || s.State == "paused" {
			v.status = withDetail(s, s.Str("media_title"))
		}
	case "vacuum":
		v.icon, v.color = lucide.RobotVacuum, teal
	case "lawn_mower":
		v.icon, v.color = lucide.Tractor, teal
	case "alarm_control_panel":
		v.icon, v.color = lucide.ShieldCheck, green
		switch s.State {
		case "disarmed":
			v.icon = lucide.Shield
		case "triggered":
			v.icon, v.color = lucide.ShieldAlert, danger
		case "pending", "arming", "disarming":
			v.icon, v.color = lucide.ShieldEllipsis, orange
		}
	case "siren":
		v.icon, v.color = lucide.Siren, danger
	case "camera":
		v.kind, v.icon, v.lit, v.neutral = tileStatic, lucide.Cctv, false, true
	case "select":
		v.icon, v.lit, v.neutral = lucide.ListChecks, false, true
		if s.Available() && s.State != "unknown" {
			v.status = s.State
		}
	case "button":
		v.kind, v.icon, v.lit, v.neutral = tilePress, lucide.CircleDot, false, true
		if s.Available() {
			v.status = message("action_press")
		}
	}
	return v
}

// climateTile shows an air conditioner as its controls do, local changes
// included; a click opens them.
func (a *application) climateTile(s ha.State) tileView {
	cl := a.shown(ha.ParseClimate(s))
	v := tileView{kind: tilePage, icon: lucide.AirVent, color: modeColor(cl.Mode), lit: cl.On()}
	switch {
	case !cl.Available():
		v.status = locales.Unavailable()
	case cl.On():
		v.icon, v.status = modeIcon(cl.Mode), optionText("mode", cl.Mode)
		if cl.HasTarget {
			v.status += " · " + formatTemp(cl.Target) + "°"
		}
	default:
		v.status = message("state_off")
	}
	return v
}

// sensorTile shows what a device that only measures reads: its most telling
// reading and the next one. A door left open or motion lights it.
func (a *application) sensorTile(t ha.Tile, s ha.State) tileView {
	v := tileView{kind: tileStatic, neutral: true, status: a.readingText(s)}
	if len(t.Readings) > 0 {
		v.status += " · " + a.readingText(t.Readings[0])
	}
	class := s.Str("device_class")
	if s.Domain() == "binary_sensor" {
		v.icon, v.color, v.lit, v.neutral = binaryIcon(s), orange, s.Active(), false
		switch class {
		case "moisture", "smoke", "gas", "carbon_monoxide", "safety", "heat", "cold":
			v.color = danger
		}
		return v
	}
	switch class {
	case "humidity", "moisture", "water":
		v.icon, v.color = lucide.Droplet, accent
	case "illuminance":
		v.icon, v.color = lucide.Sun, amber
	case "pressure", "atmospheric_pressure":
		v.icon, v.color = lucide.Gauge, lavender
	case "power", "energy":
		v.icon, v.color = lucide.Zap, amber
	case "gas":
		v.icon, v.color = lucide.Flame, orange
	case "pm25", "pm10", "pm1", "carbon_dioxide", "aqi", "volatile_organic_compounds",
		"volatile_organic_compounds_parts", "carbon_monoxide", "nitrogen_dioxide", "ozone":
		v.icon, v.color = lucide.Wind, teal
	default:
		v.icon, v.color = lucide.Thermometer, orange
	}
	return v
}

func binaryIcon(s ha.State) func() *ui.SVG {
	on := s.State == "on"
	switch s.Str("device_class") {
	case "door", "opening":
		if on {
			return lucide.DoorOpen
		}
		return lucide.DoorClosed
	case "garage_door":
		return lucide.Warehouse
	case "window":
		return lucide.AppWindow
	case "lock":
		if on {
			return lucide.LockOpen
		}
		return lucide.Lock
	case "motion", "occupancy", "presence":
		return lucide.PersonStanding
	case "moisture":
		return lucide.Droplets
	case "smoke", "carbon_monoxide":
		return lucide.AlarmSmoke
	case "gas":
		return lucide.Flame
	case "safety":
		return lucide.ShieldAlert
	case "vibration":
		return lucide.Vibrate
	case "sound":
		return lucide.AudioLines
	case "heat":
		return lucide.ThermometerSun
	case "cold":
		return lucide.ThermometerSnowflake
	case "light":
		return lucide.Sun
	case "plug":
		return lucide.Plug
	case "power":
		return lucide.Zap
	case "moving":
		return lucide.Move
	}
	return lucide.Activity
}

func coverIcon(s ha.State) func() *ui.SVG {
	switch s.Str("device_class") {
	case "garage":
		return lucide.Warehouse
	case "door", "gate":
		if s.State == "closed" {
			return lucide.DoorClosed
		}
		return lucide.DoorOpen
	case "window":
		return lucide.AppWindow
	}
	return lucide.Blinds
}

// binaryText names, per device class, what a binary sensor's off and on
// mean; others say Off and On.
var binaryText = map[string][2]string{
	"door": {"state_closed", "state_open"}, "garage_door": {"state_closed", "state_open"},
	"window": {"state_closed", "state_open"}, "opening": {"state_closed", "state_open"},
	"lock":   {"state_locked", "state_unlocked"},
	"motion": {"sensor_clear", "sensor_detected"}, "occupancy": {"sensor_clear", "sensor_detected"},
	"presence": {"sensor_clear", "sensor_detected"}, "vibration": {"sensor_clear", "sensor_detected"},
	"sound": {"sensor_clear", "sensor_detected"}, "smoke": {"sensor_clear", "sensor_detected"},
	"gas": {"sensor_clear", "sensor_detected"}, "carbon_monoxide": {"sensor_clear", "sensor_detected"},
	"moisture": {"sensor_dry", "sensor_wet"},
	"safety":   {"sensor_safe", "sensor_unsafe"},
	"heat":     {"sensor_normal", "sensor_hot"}, "cold": {"sensor_normal", "sensor_cold"},
	"light":  {"sensor_dark", "sensor_light"},
	"plug":   {"sensor_unplugged", "sensor_plugged"},
	"moving": {"sensor_still", "sensor_moving"}, "running": {"sensor_stopped", "sensor_running"},
}

// readingText is what sensor s reads: a binary sensor's state as its device
// class names it, else the value with its unit, rounded as Home Assistant
// says.
func (a *application) readingText(s ha.State) string {
	switch {
	case !s.Available():
		return locales.Unavailable()
	case s.State == "unknown":
		return "—"
	case s.Domain() == "binary_sensor":
		keys, ok := binaryText[s.Str("device_class")]
		if !ok {
			keys = [2]string{"state_off", "state_on"}
		}
		if s.State == "on" {
			return message(keys[1])
		}
		return message(keys[0])
	}
	value := s.State
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		if p, ok := a.registry.Precision(s.EntityID); ok {
			value = strconv.FormatFloat(f, 'f', p, 64)
		} else {
			value = strconv.FormatFloat(math.Round(f*100)/100, 'f', -1, 64)
		}
	}
	switch unit := s.Str("unit_of_measurement"); unit {
	case "":
		return value
	case "%", "°":
		return value + unit
	default:
		return value + " " + unit
	}
}

// stateText names the state of entity s, as "Open", in the app's language;
// states it has no message for show as they are.
func stateText(s ha.State) string {
	if !s.Available() {
		return locales.Unavailable()
	}
	return optionText("state", s.State)
}

// withDetail is the state of s followed by detail, as "On · 80%", while it
// is active and there is a detail.
func withDetail(s ha.State, detail string) string {
	if detail == "" || !s.Active() {
		return stateText(s)
	}
	return stateText(s) + " · " + detail
}

// percent is attribute k, out of full, as a percentage.
func percent(s ha.State, k string, full float64) string {
	v, ok := s.Num(k)
	if !ok || full <= 0 {
		return ""
	}
	return strconv.Itoa(int(math.Round(v/full*100))) + "%"
}

// position is how far a cover or valve is open, unless all the way.
func position(s ha.State) string {
	if v, ok := s.Num("current_position"); ok && v > 0 && v < 100 {
		return strconv.Itoa(int(v)) + "%"
	}
	return ""
}

// lightColor is the color a light shines in, if it says, lightened to show
// on the panel.
func lightColor(s ha.State) (ui.Color, bool) {
	rgb, _ := s.Attributes["rgb_color"].([]any)
	if len(rgb) != 3 {
		return ui.Color{}, false
	}
	var c [3]float64
	top := 0.0
	for i, v := range rgb {
		f, ok := v.(float64)
		if !ok {
			return ui.Color{}, false
		}
		c[i], top = f, max(top, f)
	}
	if top <= 0 {
		return ui.Color{}, false
	}
	k := 255 / top
	return ui.RGB(uint8(c[0]*k), uint8(c[1]*k), uint8(c[2]*k)).Mix(ui.RGB(255, 255, 255), 0.35), true
}

// actionText names action act of entity s in a menu: an option as the
// entity names it, else what the action does.
func actionText(act ha.Action, s ha.State) string {
	switch {
	case !act.Option:
		return message("action_" + act.Name)
	case s.Domain() == "water_heater":
		return optionText("water", act.Name)
	}
	return act.Name
}
