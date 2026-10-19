package main

import (
	"math"
	"time"

	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// The controls every device's page shares: the power button at the top
// right, the slider of its level and the tiles along the bottom.

// powerButton switches a device on and off, glowing in tint while it is on.
// label names it.
func powerButton(c *ui.Context, label string, on, disabled bool, tint ui.Color) ui.Element {
	glow := float32(0)
	if on {
		glow = 1
	}
	power := ui.ButtonBase(c).Size(44, 44).Radius(22).Shrink(0).Label(label).Checked(on).Disabled(disabled)
	glow = power.Animate("glow", glow, 220*time.Millisecond)
	face := ui.RGBA(255, 255, 255, 0.08).Mix(tint, glow)
	switch {
	case disabled:
	case power.Pressed():
		face = face.Mix(panelBottom, 0.2)
	case power.Hovered():
		face = face.Mix(ink, 0.08)
	}
	power.Background(face).Shadow(0, 4, 18, 0, tint.Alpha(0.55*glow)).Center()
	power.Children(func() { components.Icon(c, lucide.Power).Size(21, 21).TextColor(ink) })
	if disabled {
		power.Opacity(0.45)
	}
	return power
}

// circleButton is a circular icon button of size DIPs across.
func circleButton(c *ui.Context, icon func() *ui.SVG, label string, size float32, disabled bool) ui.Element {
	b := ui.ButtonBase(c).Size(size, size).Radius(size / 2).Shrink(0).Label(label).Disabled(disabled).Center()
	face := ui.RGBA(255, 255, 255, 0.06)
	switch {
	case disabled:
	case b.Pressed():
		face = ui.RGBA(255, 255, 255, 0.14)
	case b.Hovered():
		face = ui.RGBA(255, 255, 255, 0.1)
	}
	b.Background(face).Border(1, tileEdge)
	glyph := min(size/2, 18)
	b.Children(func() { components.Icon(c, icon).Size(glyph, glyph).TextColor(ink) })
	if disabled {
		b.Opacity(0.4)
	}
	components.Tooltip(c, b, label)
	return b
}

// slider is a setting of a device that slides between bounds: an air
// conditioner's temperature or a light's brightness.
type slider struct {
	Value, Min, Max, Step float64
	// Label names the slider, Lower and Raise the buttons beside it.
	Label, Lower, Raise string
	// Color fills the track up to the thumb.
	Color    ui.Color
	Disabled bool
	// Bound shows Min and Max below the track, as "16°".
	Bound func(float64) string
	// Set shows a value at once; it is up to it to send it once the slider
	// rests.
	Set func(float64)
}

// sliderRow sets s with the slider between buttons that lower and raise it
// by a step, its bounds below it.
func sliderRow(c *ui.Context, s slider) {
	ui.Row(c).FillWidth().AlignItems(ui.Start).Gap(12).Children(func() {
		if circleButton(c, lucide.Minus, s.Lower, 44, s.Disabled || s.Value <= s.Min).Margin(6, 0, 0, 0).Clicked() {
			s.Set(max(s.Min, s.Value-s.Step))
		}
		ui.Column(c).Grow(1).Basis(0).Gap(8).Children(func() {
			levelSlider(c, s)
			ui.Row(c).FillWidth().PaddingX(10).Children(func() {
				components.Text(c, s.Bound(s.Min)).FontSize(13).TextColor(inkMuted)
				ui.Spacer(c)
				components.Text(c, s.Bound(s.Max)).FontSize(13).TextColor(inkMuted)
			})
		})
		if circleButton(c, lucide.Plus, s.Raise, 44, s.Disabled || s.Value >= s.Max).Margin(6, 0, 0, 0).Clicked() {
			s.Set(min(s.Max, s.Value+s.Step))
		}
	})
}

// levelSlider is a dark track with ticks, filled in s's color up to a thumb
// of liquid glass, which shows the fill's edge bent through it.
func levelSlider(c *ui.Context, s slider) {
	const height, inset, radius = 56, 4, 18
	thumb := float32(height - 2*inset)
	value := s.Value
	e := ui.SliderBase(c, &value, s.Min, s.Max).Step(s.Step).Label(s.Label).Disabled(s.Disabled).
		FillWidth().Height(height).PaddingX(inset + thumb/2).FocusRing(false)
	if e.Changed() {
		s.Set(value)
	}
	_, _, held := e.Dragged()
	grow := e.Animate("held", map[bool]float32{false: 0, true: 1}[held && !s.Disabled], 150*time.Millisecond)
	focus := e.FocusVisible()
	fraction := float32(0)
	if s.Max > s.Min {
		fraction = float32((value - s.Min) / (s.Max - s.Min))
	}
	ticks := 0
	if s.Step > 0 {
		ticks = int(math.Round((s.Max - s.Min) / s.Step))
	}
	for ticks > 28 {
		ticks /= 2
	}
	if s.Disabled {
		e.Opacity(0.45)
	}
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, ui.RGBA(0, 0, 0, 0.42), radius)
		p.Stroke(r, ui.RGBA(255, 255, 255, 0.05), radius, 1)
		left, width := r.X+inset+thumb/2, r.W-2*inset-thumb
		x := left + width*fraction
		fill := ui.Rect{X: r.X + inset, Y: r.Y + inset, W: x - r.X - inset, H: thumb}
		p.Shadow(fill, radius-inset, 0, 0, 16, 0, s.Color.Alpha(0.25))
		p.FillGradient(fill, ui.LinearGradient{From: s.Color.Alpha(0.2), To: s.Color.Alpha(0.85), Angle: 90}, radius-inset)
		for i := 0; i <= ticks && ticks > 0; i++ {
			tx := left + width*float32(i)/float32(ticks)
			if abs32(tx-x) < thumb/2+4 {
				continue
			}
			h, tick := float32(8), ui.RGBA(255, 255, 255, 0.22)
			if i%5 == 0 {
				h = 12
			}
			if tx < x {
				tick = discInk.Alpha(0.3)
			}
			p.Fill(ui.Rect{X: tx - 0.75, Y: r.Y + r.H/2 - h/2, W: 1.5, H: h}, tick, 0.75)
		}
		size := thumb * (1 + 0.1*grow)
		round := (radius - inset) * (1 + 0.1*grow)
		k := ui.Rect{X: x - size/2, Y: r.Y + r.H/2 - size/2, W: size, H: size}
		p.Shadow(k, round, 0, 6, 16, 0, ui.RGBA(0, 0, 0, 0.3+0.1*grow))
		glass.Paint(p, k, round, glass.Glass{Style: glass.Clear, Tint: ui.RGBA(255, 255, 255, 0.1+0.06*grow)})
		// A grip down the middle of the thumb.
		p.Fill(ui.Rect{X: k.X + k.W/2 - 1.5, Y: k.Y + k.H/2 - 9, W: 3, H: 18}, ink.Alpha(0.85), 1.5)
		if focus {
			p.FocusRing(ui.Rect{X: k.X - 3, Y: k.Y - 3, W: k.W + 6, H: k.H + 6}, [4]float32{round + 3, round + 3, round + 3, round + 3})
		}
	})
}

// controlHeight is the height of the tiles along the bottom of a device's
// page.
const controlHeight = 104

// tileBase is a tile along the bottom of a device's page, which a click
// acts on unless it only shows a reading.
func tileBase(c *ui.Context, label string, disabled bool) ui.Element {
	t := ui.ButtonBase(c).Grow(1).Basis(0).MinWidth(0).Height(controlHeight).Radius(20).Padding(12, 8, 12, 10).
		Column().AlignItems(ui.Start).Label(label).Disabled(disabled)
	face := tileFace
	switch {
	case disabled:
	case t.Pressed():
		face = ui.RGBA(255, 255, 255, 0.12)
	case t.Hovered():
		face = tileHover
	}
	t.Background(face).Border(1, tileEdge)
	if disabled {
		t.Opacity(0.45)
	}
	return t
}

// tileText names what a tile sets, then its value.
func tileText(c *ui.Context, title, value string) {
	components.Text(c, title).FontSize(12).TextColor(inkMuted).SingleLine().Ellipsis("…")
	components.Text(c, value).FontSize(12).Bold().TextColor(ink).SingleLine().Ellipsis("…").Margin(2, 0, 0, 0)
}

// tileIcon is the icon at the top left of a tile, lit in color.
func tileIcon(c *ui.Context, icon func() *ui.SVG, color ui.Color) {
	components.Icon(c, icon).Size(26, 26).TextColor(color.Mix(ink, 0.15))
}

// tileDot is the dot at the top right of a tile that switches, glowing in
// color while on.
func tileDot(c *ui.Context, on bool, color ui.Color) {
	dot := ui.RGBA(255, 255, 255, 0.3)
	if on {
		dot = color
	}
	d := ui.Box(c).Size(8, 8).Radius(4).Background(dot).Margin(4, 0, 0, 0).MarginEnd(3)
	if on {
		d.Shadow(0, 0, 8, 0, color.Alpha(0.8))
	}
}

// menuTile opens a menu of options, kind's option current checked, which
// choose takes, below icon and what it sets, title.
func menuTile(c *ui.Context, icon func() *ui.SVG, color ui.Color, title, value, kind, current string, options []string, disabled bool, choose func(string)) {
	t := tileBase(c, title+": "+value, len(options) == 0 || disabled)
	t.Children(func() {
		ui.Row(c).FillWidth().Children(func() {
			tileIcon(c, icon, color)
			ui.Spacer(c)
			components.Icon(c, lucide.ChevronDown).Size(15, 15).TextColor(inkMuted)
		})
		ui.Spacer(c)
		tileText(c, title, value)
	})
	t.Menu(func(m *ui.Menu) {
		for _, o := range options {
			if m.Item(optionText(kind, o)).Checked(o == current).Chosen() {
				choose(o)
			}
		}
	})
}

// tileRows lays n tiles out along the bottom of a device's page, four to a
// row at most and in rows as even as can be, tile building the i-th. A tile
// alone takes half the width, in the middle.
func tileRows(c *ui.Context, n int, tile func(i int)) {
	if n == 0 {
		return
	}
	cols := n
	if n > 4 {
		cols = 4
		if n <= 6 {
			cols = 3
		}
	}
	ui.Column(c).FillWidth().Gap(8).Margin(18, 0, 0, 0).Children(func() {
		if n == 1 {
			ui.Row(c).FillWidth().Gap(8).Children(func() {
				ui.Box(c).Grow(0.5).Basis(0)
				tile(0)
				ui.Box(c).Grow(0.5).Basis(0)
			})
			return
		}
		for row := 0; row < n; row += cols {
			ui.Row(c.Key(row)).FillWidth().Gap(8).Children(func() {
				for i := row; i < row+cols; i++ {
					if i < n {
						tile(i)
					} else {
						// Keeps the last row's tiles as wide as the others.
						ui.Box(c.Key(i)).Grow(1).Basis(0)
					}
				}
			})
		}
	})
}
