package main

import (
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/neohm/source/internal/ha"
	"github.com/oneo-me/neohm/source/locales"
	"github.com/oneo-me/onevra/source/components"
)

// lightColors offers the colors light s shines in a click away, below its
// slider: its whites in a row above its colors. The last swatch opens a
// picker of any color, lit while the light shines in one the others lack.
func (a *application) lightColors(c *ui.Context, s ha.State) {
	whites, colorful := s.Swatches(), s.Colorful()
	_, _, warmth := s.Whites()
	if !warmth && !colorful {
		return
	}
	var hues []ha.Color
	if colorful {
		hues = ha.Hues
	}
	shining, lit := s.Shining()
	other := lit
	for _, col := range append(whites, hues...) {
		if lit && shining.Near(col) {
			other = false
		}
	}
	disabled := !s.Available()
	row := func(key string, colors []ha.Color, picker bool) {
		ui.Row(c.Key(key)).FillWidth().Justify(ui.Center).Gap(6).Children(func() {
			for i, col := range colors {
				if swatch(c.Key(i), colorName(col), col, lit && shining.Near(col), disabled).Clicked() {
					a.perform(s, col.Set())
				}
			}
			if picker {
				a.colorPicker(c, s, shining, other, disabled)
			}
		})
	}
	ui.Column(c).FillWidth().Gap(4).Margin(18, 0, 0, 0).Children(func() {
		row("whites", whites, !colorful)
		if colorful {
			row("hues", hues, true)
		}
	})
}

// colorName names col for the screen reader: a white by its color
// temperature, else as it is called.
func colorName(col ha.Color) string {
	if col.Kelvin > 0 {
		return strconv.Itoa(int(col.Kelvin)) + " K"
	}
	return message("color_" + col.Name)
}

// uiColor is how col looks on the panel.
func uiColor(col ha.Color) ui.Color { return ui.RGB(col.RGB()) }

// swatchSize is how wide a swatch is, the ring around it included.
const swatchSize = 46

// swatchBase is a round button that shows a color, which paint fills,
// ringed while current, as the color the light shines in.
func swatchBase(c *ui.Context, label string, current, disabled bool, paint func(p *ui.Painter, r ui.Rect)) ui.Element {
	b := ui.ButtonBase(c).Size(swatchSize, swatchSize).Radius(swatchSize / 2).Shrink(0).Label(label).Checked(current).Disabled(disabled)
	ring := float32(0)
	if !disabled && current {
		ring = 1
	}
	ring = b.Animate("ring", ring, 150*time.Millisecond)
	press := b.Animate("press", map[bool]float32{false: 0, true: 1}[b.Pressed()], 100*time.Millisecond)
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		inset := 5 + 1.5*press
		face := ui.Rect{X: r.X + inset, Y: r.Y + inset, W: r.W - 2*inset, H: r.H - 2*inset}
		p.Shadow(face, face.W/2, 0, 3, 10, 0, ui.RGBA(0, 0, 0, 0.3))
		paint(p, face)
		p.Stroke(face, ui.RGBA(255, 255, 255, 0.14), face.W/2, 1)
		if ring > 0 {
			o := ui.Rect{X: r.X + 1, Y: r.Y + 1, W: r.W - 2, H: r.H - 2}
			p.Stroke(o, ink.Alpha(ring), o.W/2, 2)
		}
	})
	if disabled {
		b.Opacity(0.45)
	}
	return b
}

// swatch sets the light to col, named label.
func swatch(c *ui.Context, label string, col ha.Color, current, disabled bool) ui.Element {
	face := uiColor(col)
	return swatchBase(c, label, current, disabled, func(p *ui.Painter, r ui.Rect) { p.Fill(r, face, r.W/2) })
}

// hueStops are the hues around a turn, red to red.
var hueStops = func() []ui.GradientStop {
	var out []ui.GradientStop
	for i := 0; i <= 6; i++ {
		out = append(out, ui.GradientStop{Color: uiColor(ha.Color{Hue: float64(i * 60), Saturation: 100}), At: float32(i) / 6})
	}
	return out
}()

// whiteStops are the whites from warm to cool, warm to cool.
func whiteStops(warm, cool float64) []ui.GradientStop {
	const n = 8
	var out []ui.GradientStop
	for i := range n + 1 {
		m := 1e6/warm + (1e6/cool-1e6/warm)*float64(i)/n
		out = append(out, ui.GradientStop{Color: uiColor(ha.White(1e6 / m)), At: float32(i) / n})
	}
	return out
}

// colorPicker opens a picker of any color light s shines in, which shining
// is now: the hues and how saturated, its whites from warm to cool, or
// both. other lights it, as the color shown is none of the swatches.
func (a *application) colorPicker(c *ui.Context, s ha.State, shining ha.Color, other, disabled bool) {
	warm, cool, warmth := s.Whites()
	colorful := s.Colorful()
	b := swatchBase(c.Key("picker"), locales.ChooseColor(), other, disabled, func(p *ui.Painter, r ui.Rect) {
		if colorful {
			p.FillConicGradient(r, ui.ConicGradient{Stops: hueStops}, r.W/2)
			p.FillGradient(r, ui.LinearGradient{From: ui.RGBA(255, 255, 255, 0.55), To: ui.RGBA(255, 255, 255, 0), Angle: 180}, r.W/2)
			return
		}
		// As the whites go, warm at the top.
		p.FillGradient(r, ui.LinearGradient{Angle: 180, Stops: whiteStops(warm, cool)}, r.W/2)
	})
	open := ui.Local(b, "open", func() bool { return false })
	if disabled {
		*open = false
	} else if b.Clicked() {
		*open = !*open
	}
	b.Expanded(*open)
	components.Popover(c, b, open, func(panel ui.Element) {
		panel.Width(272).Label(locales.ChooseColor())
		if colorful {
			a.huePlane(c, s, shining)
		}
		if warmth {
			a.warmthSlider(c, s, shining, warm, cool)
		}
	})
}

// huePlane picks a color of light s, shining now, on a plane of the hues
// across it, fading to white towards the bottom. Arrow keys move along it.
func (a *application) huePlane(c *ui.Context, s ha.State, shining ha.Color) {
	const radius = 14
	col := shining
	if col.Kelvin > 0 {
		col = ha.Color{}
	}
	plane := ui.Box(c).FillWidth().Height(150).Radius(radius).Focusable().FocusRing(false).
		Role(ui.RoleSlider).Label(locales.Color()).Cursor(ui.CursorCrosshair)
	next := col
	if _, _, held := plane.Dragged(); held {
		if x, y, _ := plane.PointerPosition(); plane.LayoutBounds().W > 0 {
			r := plane.LayoutBounds()
			next.Hue = float64(max(0, min(1, x/r.W))) * 360
			next.Saturation = (1 - float64(max(0, min(1, y/r.H)))) * 100
			plane.Focus()
		}
	}
	switch {
	case plane.Shortcut(0, ui.KeyRight):
		next.Hue = math.Mod(next.Hue+5, 360)
	case plane.Shortcut(0, ui.KeyLeft):
		next.Hue = math.Mod(next.Hue+355, 360)
	case plane.Shortcut(0, ui.KeyUp):
		next.Saturation = min(100, next.Saturation+5)
	case plane.Shortcut(0, ui.KeyDown):
		next.Saturation = max(0, next.Saturation-5)
	}
	if next != col {
		a.adjust(s, next.Set())
	}
	plane.Range(0, 360, next.Hue).Step(5).
		Value(strconv.Itoa(int(math.Round(next.Hue))) + "°, " + strconv.Itoa(int(math.Round(next.Saturation))) + "%")
	focus := plane.FocusVisible()
	face := uiColor(next)
	plane.Draw(func(p *ui.Painter, r ui.Rect) {
		p.FillGradient(r, ui.LinearGradient{Angle: 90, Stops: hueStops}, radius)
		p.FillGradient(r, ui.LinearGradient{From: ui.RGBA(255, 255, 255, 0), To: ui.RGB(255, 255, 255), Angle: 180}, radius)
		if focus {
			p.FocusRing(r, [4]float32{radius, radius, radius, radius})
		}
		if shining.Kelvin > 0 {
			return
		}
		k := ui.Rect{X: r.X + r.W*float32(next.Hue/360) - 12, Y: r.Y + r.H*float32(1-next.Saturation/100) - 12, W: 24, H: 24}
		p.Shadow(k, 12, 0, 2, 8, 0, ui.RGBA(0, 0, 0, 0.35))
		p.Fill(k, face, 12)
		p.Stroke(k, ui.RGB(255, 255, 255), 12, 3)
	})
}

// warmthSlider picks a white of light s, shining now, from its warmest,
// warm, to its coolest, cool.
func (a *application) warmthSlider(c *ui.Context, s ha.State, shining ha.Color, warm, cool float64) {
	k := shining.Kelvin
	if k == 0 {
		k = math.Round(math.Sqrt(warm * cool))
	}
	ui.Row(c).FillWidth().Children(func() {
		components.Text(c, locales.ColorTemperature()).FontSize(13).TextColor(inkMuted).Grow(1).MinWidth(0).SingleLine().Ellipsis("…")
		if shining.Kelvin > 0 {
			components.Text(c, colorName(shining)).FontSize(13).Bold().TextColor(ink).SingleLine()
		}
	})
	const height, radius = 30, 15
	value := k
	e := ui.SliderBase(c, &value, warm, cool).Step(50).Label(locales.ColorTemperature()).
		FillWidth().Height(height).PaddingX(radius).FocusRing(false).Value(colorName(ha.White(value)))
	if e.Changed() {
		a.adjust(s, ha.White(value).Set())
	}
	focus := e.FocusVisible()
	stops := whiteStops(warm, cool)
	fraction := float32(0)
	if cool > warm {
		fraction = float32((value - warm) / (cool - warm))
	}
	face := uiColor(ha.White(value))
	shown := shining.Kelvin > 0
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		p.FillGradient(r, ui.LinearGradient{Angle: 90, Stops: stops}, radius)
		if !shown {
			return
		}
		x := r.X + radius + (r.W-2*radius)*fraction
		knob := ui.Rect{X: x - 12, Y: r.Y + r.H/2 - 12, W: 24, H: 24}
		p.Shadow(knob, 12, 0, 2, 8, 0, ui.RGBA(0, 0, 0, 0.35))
		p.Fill(knob, face, 12)
		p.Stroke(knob, ui.RGB(255, 255, 255), 12, 3)
		if focus {
			p.FocusRing(ui.Rect{X: knob.X - 3, Y: knob.Y - 3, W: knob.W + 6, H: knob.H + 6}, [4]float32{15, 15, 15, 15})
		}
	})
}
