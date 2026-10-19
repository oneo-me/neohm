package main

import (
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/neohm/source/internal/ha"
	"github.com/oneo-me/neohm/source/locales"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// devicePage shows the device of tile t, in area, on a page of its own, laid
// out as an air conditioner's controls: the device drawn, its main value,
// the slider of its level and what else it does or reads along the bottom,
// below a header with its power button. radius rounds the corners it
// scrolls within.
func (a *application) devicePage(c *ui.Context, t ha.Tile, area string, radius float32) {
	s := a.held(t.Entity)
	v := a.tileView(t, s)
	name := tileName(t.Name, area)
	l, hasLevel := s.Level()
	hasLevel = hasLevel && s.Available()
	value := 0.0
	if hasLevel {
		value = a.level(s, l)
	}
	pageFrame(c, radius, name, func() {
		a.titledHeader(c, name, func() {
			if problem := a.problem(s.Available()); problem != "" {
				headerSubtitle(c, problem, danger)
			} else if area != "" {
				headerSubtitle(c, area, inkMuted)
			}
		}, func() { a.devicePower(c, s, v, name) })
	}, func(top float32) {
		view := components.ScrollView(c, components.ScrollViewOptions{}).Fill().ScrollbarInsets(top, 0, radius)
		// The controls keep to the bottom, as an air conditioner's, unless
		// they scroll.
		fill := view.LayoutBounds().H
		view.Children(func() {
			ui.Column(c).FillWidth().MinHeight(fill).Padding(top+pageFade, 18, 18, 18).Children(func() {
				deviceScene(c, s, v, a.sceneAmount(s, l, hasLevel, value))
				a.readout(c, t, s, v, l, hasLevel, value)
				if hasLevel {
					title := message("level_" + l.Name)
					sliderRow(c, slider{
						Value: value, Min: l.Min, Max: l.Max, Step: l.Step,
						Label: title, Lower: locales.Decrease(title), Raise: locales.Increase(title),
						Color: v.color,
						Bound: func(f float64) string { return strconv.Itoa(int(math.Round(f))) + "%" },
						Set:   func(f float64) { a.setLevel(s, l, f) },
					})
				}
				a.deviceTiles(c, t, s, v)
			})
		})
	})
}

// problem is what keeps a device from being controlled, if anything: the
// error shown, the connection lost or the device unavailable.
func (a *application) problem(available bool) string {
	switch {
	case a.errorDisplay() != "":
		return a.errorDisplay()
	case !a.online:
		return locales.Offline()
	case !available:
		return locales.Unavailable()
	}
	return ""
}

// devicePower is the power button of a device that switches on and off,
// as an air conditioner's. Covers and valves open and close from their
// tiles instead.
func (a *application) devicePower(c *ui.Context, s ha.State, v tileView, name string) {
	act, ok := s.Toggle()
	if !ok || act.Name != "turn_on" && act.Name != "turn_off" {
		return
	}
	if powerButton(c, locales.QuickSwitch(name), v.lit, !s.Available(), v.color).Role(ui.RoleSwitch).Clicked() {
		a.toggle(s)
	}
}

// sceneAmount is how far the level of entity s goes, 0 to 1, for its
// drawing: the level set, else whether it is on or open.
func (a *application) sceneAmount(s ha.State, l ha.Level, hasLevel bool, value float64) float64 {
	if hasLevel && l.Max > l.Min {
		return (value - l.Min) / (l.Max - l.Min)
	}
	if s.Active() {
		return 1
	}
	return 0
}

// readout shows the main value of a device large, as an air conditioner's
// target temperature: the level it is set to, its main reading or its
// state, with a detail, as the title playing, below it. What the value is
// shows in a pill beside it, as an air conditioner's Eco Mode.
func (a *application) readout(c *ui.Context, t ha.Tile, s ha.State, v tileView, l ha.Level, hasLevel bool, value float64) {
	var caption, number, unit, text, detail string
	switch {
	case s.Domain() == "sensor" || s.Domain() == "binary_sensor":
		// One Home Assistant names only by its ID goes by the device's name.
		if s.Str("friendly_name") != "" {
			caption = tileName(s.Name(), t.Name)
		}
		number, unit = a.readingParts(s)
		if _, err := strconv.ParseFloat(number, 64); err != nil {
			text, number, unit = a.readingText(s), "", ""
		}
	case hasLevel && (l.Name != "brightness" || s.Active()):
		caption = message("level_" + l.Name)
		number, unit = strconv.Itoa(int(math.Round(value))), "%"
	default:
		text = stateText(s)
		if rest, ok := strings.CutPrefix(v.status, text+" · "); ok {
			detail = rest
		}
	}
	ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(10).Margin(0, 0, 14, 0).Children(func() {
		ui.Column(c).Grow(1).MinWidth(0).Children(func() {
			if number != "" {
				ui.Row(c).AlignItems(ui.End).Children(func() {
					components.Text(c, number).FontSize(68).FontWeight(300).TextColor(ink).LetterSpacing(-2).SingleLine()
					components.Text(c, unit).FontSize(20).FontWeight(500).TextColor(v.color).Margin(0, 0, 12, 0).MarginStart(4).SingleLine()
				})
			} else {
				components.Text(c, text).FontSize(40).FontWeight(300).TextColor(ink).MaxLines(2).Ellipsis("…")
			}
			if detail != "" {
				components.Text(c, detail).FontSize(14).TextColor(inkMuted.Mix(v.color, 0.4)).SingleLine().Ellipsis("…")
			}
		})
		if caption == "" {
			return
		}
		ui.Row(c).Shrink(1).MinWidth(0).AlignItems(ui.Center).Gap(8).Padding(10, 14).Radius(18).
			Background(tileFace).Border(1, tileEdge).Children(func() {
			components.Icon(c, v.icon).Size(18, 18).TextColor(v.color.Mix(ink, 0.15)).Shrink(0)
			components.Text(c, caption).FontSize(14).FontWeight(600).TextColor(ink).SingleLine().Ellipsis("…")
		})
	})
}

// deviceTiles are what a device does besides switching on and off, or what
// else it reads, as tiles along the bottom of its page. Options, as a
// select's, make one tile whose menu chooses among them.
func (a *application) deviceTiles(c *ui.Context, t ha.Tile, s ha.State, v tileView) {
	if d := s.Domain(); d == "sensor" || d == "binary_sensor" {
		tileRows(c, len(t.Readings), func(i int) { a.readingTile(c, t, t.Readings[i]) })
		return
	}
	power := false
	if act, ok := s.Toggle(); ok {
		power = act.Name == "turn_on" || act.Name == "turn_off"
	}
	var actions, options []ha.Action
	for _, act := range s.Actions() {
		switch {
		case power && (act.Name == "turn_on" || act.Name == "turn_off"):
		case act.Option:
			options = append(options, act)
		default:
			actions = append(actions, act)
		}
	}
	n := len(actions)
	if len(options) > 0 {
		n++
	}
	tileRows(c, n, func(i int) {
		if i < len(actions) {
			act := actions[i]
			if actionTile(c.Key(act.Name), act, s, v).Clicked() {
				a.perform(s, act)
			}
			return
		}
		a.optionTile(c, s, v, options)
	})
}

// actionTile is a tile that asks for act, its dot lit while the device is
// where it leads.
func actionTile(c *ui.Context, act ha.Action, s ha.State, v tileView) ui.Element {
	label := actionText(act, s)
	tile := tileBase(c, label, act.Disabled || !s.Available()).Checked(act.Current)
	tile.Children(func() {
		ui.Row(c).FillWidth().Children(func() {
			tileIcon(c, actionIcon(act.Name), v.color)
			ui.Spacer(c)
			if act.Current {
				tileDot(c, true, v.color)
			}
		})
		ui.Spacer(c)
		components.Text(c, label).FontSize(12).Bold().TextColor(ink).MaxLines(2).Ellipsis("…")
	})
	return tile
}

// optionTile chooses among the options of entity s from its menu.
func (a *application) optionTile(c *ui.Context, s ha.State, v tileView, options []ha.Action) {
	kind, title := "select", locales.Option()
	if s.Domain() == "water_heater" {
		kind, title = "water", locales.Mode()
	}
	names := make([]string, len(options))
	current := ""
	for i, o := range options {
		names[i] = o.Name
		if o.Current {
			current = o.Name
		}
	}
	menuTile(c, v.icon, v.color, title, optionText(kind, current), kind, current, names, !s.Available(), func(name string) {
		for _, o := range options {
			if o.Name == name {
				a.perform(s, o)
			}
		}
	})
}

// readingTile shows reading r of the sensor device of tile t, named
// without the device's name.
func (a *application) readingTile(c *ui.Context, t ha.Tile, r ha.State) {
	name, value := tileName(r.Name(), t.Name), a.readingText(r)
	look := a.sensorTile(ha.Tile{}, r)
	ui.Column(c.Key(r.EntityID)).Grow(1).Basis(0).MinWidth(0).Height(controlHeight).Radius(20).Padding(12, 8, 12, 10).
		Background(tileFace).Border(1, tileEdge).Label(name + ": " + value).Children(func() {
		tileIcon(c, look.icon, look.color)
		ui.Spacer(c)
		tileText(c, name, value)
	})
}

// actionIcon is the icon of the tile of the action named name.
func actionIcon(name string) func() *ui.SVG {
	switch name {
	case "turn_on", "turn_off":
		return lucide.Power
	case "open":
		return lucide.ArrowUpToLine
	case "close":
		return lucide.ArrowDownToLine
	case "stop":
		return lucide.Square
	case "lock":
		return lucide.Lock
	case "unlock":
		return lucide.LockOpen
	case "open_door":
		return lucide.DoorOpen
	case "play", "start":
		return lucide.Play
	case "pause":
		return lucide.Pause
	case "previous":
		return lucide.SkipBack
	case "next":
		return lucide.SkipForward
	case "return", "dock":
		return lucide.House
	case "locate":
		return lucide.LocateFixed
	case "disarm":
		return lucide.ShieldOff
	case "arm_home":
		return lucide.House
	case "arm_away":
		return lucide.Plane
	case "arm_night":
		return lucide.Moon
	case "arm_vacation":
		return lucide.TreePalm
	case "arm_custom_bypass":
		return lucide.ShieldPlus
	case "press":
		return lucide.Pointer
	}
	return lucide.CircleDot
}
