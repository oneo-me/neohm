package main

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/ha-go/source/internal/ha"
	"github.com/oneo-me/ha-go/source/locales"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/theme"
)

// The panel keeps one dark look in both appearances, as a remote does.
var (
	panelTop    = ui.RGB(31, 45, 80)
	panelBottom = ui.RGB(9, 14, 30)
	ink         = ui.RGB(236, 241, 250)
	inkMuted    = ui.RGBA(214, 224, 245, 0.62)
	accent      = ui.RGB(80, 144, 255)
	tileFace    = ui.RGBA(255, 255, 255, 0.045)
	tileHover   = ui.RGBA(255, 255, 255, 0.08)
	tileEdge    = ui.RGBA(255, 255, 255, 0.08)
	danger      = ui.RGB(255, 99, 99)
)

const (
	arrowHeight = 9
	panelRadius = 26
)

// panelTheme is Onevra's dark palette in the panel's colors, so that
// components such as Switch match it.
func panelTheme() *ui.Theme {
	t := theme.Palette(true)
	t.Background, t.Surface = panelBottom, tileFace
	t.Text, t.TextMuted = ink, inkMuted
	t.Accent, t.AccentHover, t.AccentPressed = accent, accent.Mix(ink, 0.12), accent.Mix(panelBottom, 0.15)
	t.Border = ui.RGBA(255, 255, 255, 0.14)
	return t
}

// modeColor is the color of the air the unit blows in mode.
func modeColor(mode string) ui.Color {
	switch mode {
	case "heat":
		return ui.RGB(255, 150, 82)
	case "dry":
		return ui.RGB(84, 214, 200)
	case "fan_only":
		return ui.RGB(170, 190, 230)
	case "auto", "heat_cool":
		return ui.RGB(168, 140, 255)
	}
	return accent
}

func modeIcon(mode string) func() *ui.SVG {
	switch mode {
	case "heat":
		return lucide.Sun
	case "dry":
		return lucide.Droplets
	case "fan_only":
		return lucide.Fan
	case "auto", "heat_cool":
		return lucide.RefreshCw
	case "off":
		return lucide.Power
	}
	return lucide.Snowflake
}

func (a *application) panelView(c *ui.Context) {
	c.SetTheme(panelTheme())
	a.scope(c)
	if c.Shortcut(0, ui.KeyEscape) || c.Shortcut(ui.Cmd, ui.KeyW) {
		a.hidePanel()
	}
	c.Root().Background(ui.Transparent)
	ui.Column(c).Fill().Padding(arrowHeight, 0, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
		// An arrow pointing at the menu bar icon, then the body.
		cx := r.X + r.W/2
		var arrow ui.Path
		arrow.MoveTo(cx-11, r.Y+arrowHeight+1).QuadTo(cx-3, r.Y+1, cx, r.Y).QuadTo(cx+3, r.Y+1, cx+11, r.Y+arrowHeight+1).Close()
		p.FillPath(&arrow, panelTop)
		body := ui.Rect{X: r.X, Y: r.Y + arrowHeight, W: r.W, H: r.H - arrowHeight}
		p.FillGradient(body, ui.LinearGradient{From: panelTop, To: panelBottom, Angle: 180}, panelRadius)
		// A soft glow behind the unit.
		p.Shadow(ui.Rect{X: cx - 70, Y: body.Y + 90, W: 140, H: 30}, 15, 0, 0, 60, 0, accent.Alpha(0.10))
	}).Children(func() {
		ui.Column(c).Fill().Padding(20, 18, 18, 18).Gap(0).Children(func() {
			if !a.configured() || !a.hasClimate {
				a.setupView(c)
				return
			}
			a.headerRow(c)
			a.unitView(c)
			ui.Spacer(c)
			a.temperatureRow(c)
			a.sliderRow(c)
			a.tiles(c)
		})
	})
}

// setupView shows while there is no climate entity to control.
func (a *application) setupView(c *ui.Context) {
	ui.Column(c).Fill().Center().Gap(14).Children(func() {
		components.Icon(c, lucide.AirVent).Size(48, 48).TextColor(accent)
		if a.configured() && a.errorKey == "" {
			components.Spinner(c, components.SpinnerRing)
			return
		}
		components.Text(c, locales.SetupTitle()).FontSize(18).Bold().TextColor(ink)
		msg := locales.SetupBody()
		if e := a.errorDisplay(); e != "" {
			msg = e
		}
		components.Text(c, msg).Wrap().TextAlign(ui.Center).FontSize(13).TextColor(inkMuted).MaxWidth(260)
		components.Button(c, locales.Settings(), components.ButtonPrimary).OnClick(a.openSettings)
	})
}

func (a *application) headerRow(c *ui.Context) {
	cl := a.climate
	ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(8).Children(func() {
		components.Icon(c, lucide.Thermometer).Size(18, 18).TextColor(accent)
		components.Text(c, locales.Room()).FontSize(15).TextColor(inkMuted)
		current := "—"
		if cl.HasCurrent {
			current = formatTemp(math.Round(cl.Current*10)/10) + a.unit
		}
		components.Text(c, current).FontSize(15).Bold().TextColor(ink)
		status, tip := accent, locales.Connected()
		switch {
		case !a.online:
			status, tip = danger, locales.Offline()
		case !cl.Available():
			status, tip = danger, locales.Unavailable()
		}
		if e := a.errorDisplay(); e != "" {
			status, tip = danger, e
		}
		dot := ui.Box(c).Size(8, 8).Radius(4).Background(status).MarginStart(4).Label(tip)
		dot.Shadow(0, 0, 8, 0, status.Alpha(0.7))
		components.Tooltip(c, dot, tip)
		ui.Spacer(c)
		on := cl.On()
		glow := float32(0)
		if on {
			glow = 1
		}
		power := ui.ButtonBase(c).Size(52, 52).Radius(26).Label(locales.Power()).Checked(on).Disabled(!cl.Available())
		glow = power.Animate("glow", glow, 220*time.Millisecond)
		tint := modeColor(cl.Mode)
		if !on {
			tint = accent
		}
		face := ui.RGBA(255, 255, 255, 0.08).Mix(tint, glow)
		if power.Pressed() {
			face = face.Mix(panelBottom, 0.2)
		} else if power.Hovered() {
			face = face.Mix(ink, 0.08)
		}
		power.Background(face).Shadow(0, 4, 22, 0, tint.Alpha(0.55*glow)).Center()
		power.Children(func() { components.Icon(c, lucide.Power).Size(24, 24).TextColor(ink) })
		if power.Clicked() {
			a.togglePower()
		}
	})
}

// unitView draws the indoor unit and the air it blows.
func (a *application) unitView(c *ui.Context) {
	cl := a.climate
	on := cl.On()
	color := modeColor(cl.Mode)
	target := formatTemp(cl.Target)
	if !cl.HasTarget {
		target = "--"
	}
	box := ui.Box(c).FillWidth().Height(150).Margin(4, 0, 0, 0)
	level := float32(0)
	if on {
		level = 1
	}
	level = box.Animate("air", level, 600*time.Millisecond)
	// The faster the fan, the faster the air moves.
	period := 2.2
	switch cl.FanMode {
	case "medium", "middle", "mid":
		period = 1.6
	case "high", "strong", "turbo":
		period = 1.1
	case "quiet", "silent", "low":
		period = 2.6
	}
	swing := cl.SwingMode != "" && cl.SwingMode != "off"
	still := c.Preferences().ReduceMotion
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		uw, uh := float32(128), float32(40)
		u := ui.Rect{X: r.X + (r.W-uw)/2, Y: r.Y + 8, W: uw, H: uh}
		now := float64(p.Now().UnixMilli()) / 1000
		if level > 0 {
			drawAir(p, u, r, color, level, now, period, swing, still)
			if !still {
				p.AnimationFrame()
			}
		}
		// The unit: a white body, its display and the vent.
		p.Shadow(u, 12, 0, 8, 20, 0, ui.RGBA(0, 0, 0, 0.35))
		p.FillGradient(u, ui.LinearGradient{From: ui.RGB(252, 253, 255), To: ui.RGB(214, 221, 236), Angle: 180}, 12)
		p.Fill(ui.Rect{X: u.X + 10, Y: u.Y + u.H - 9, W: u.W - 20, H: 3.5}, ui.RGB(132, 138, 152), 1.75)
		display := color.Alpha(0.35 + 0.65*level)
		label := ui.Span{Text: target, Size: 10, Weight: 600, Color: display}
		w, h := p.MeasureText(0, label)
		p.RichText(u.X+u.W-24-w, u.Y+9, 0, label)
		p.Fill(ui.Rect{X: u.X + u.W - 18, Y: u.Y + 9 + h/2 - 2, W: 4, H: 4}, display, 2)
	})
}

// drawAir draws streams of air fanning out below the unit u, with dashes
// flowing along them.
func drawAir(p *ui.Painter, u, r ui.Rect, color ui.Color, level float32, now, period float64, swing, still bool) {
	const streams = 11
	top := u.Y + u.H + 2
	bottom := r.Y + r.H - 4
	sway := float32(0)
	if swing && !still {
		sway = float32(math.Sin(now*2*math.Pi/5)) * 26
	}
	for i := range streams {
		f := float32(i)/(streams-1)*2 - 1 // -1 … 1 across the vent
		x0 := u.X + u.W/2 + f*(u.W/2-14)
		x1 := r.X + r.W/2 + f*r.W*0.44 + sway*(1-0.4*abs32(f))
		var path ui.Path
		path.MoveTo(x0, top).CubeTo(x0+f*6, top+30, x1-f*12+sway*0.3, bottom-60, x1, bottom)
		fade := 0.18 + 0.12*(1-abs32(f))
		p.StrokePathGradient(&path, 1.2, ui.LinearGradient{From: color.Alpha(fade * level), To: color.Alpha(0), Angle: 180})
		if still {
			continue
		}
		// Two dashes per stream, offset so the streams do not move in step.
		for k := range 2 {
			t := float32(math.Mod(now/period+float64(i)*0.137+float64(k)*0.5, 1))
			dash := path
			alpha := level * (1 - t) * 0.9
			p.StrokePath(dash.Trim(t, min(1, t+0.07)), 1.6, ink.Mix(color, 0.35).Alpha(alpha))
		}
	}
}

func abs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func (a *application) temperatureRow(c *ui.Context) {
	cl := a.climate
	ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(10).Margin(0, 0, 14, 0).Children(func() {
		ui.Row(c).AlignItems(ui.End).Children(func() {
			value := "--"
			if cl.HasTarget {
				value = formatTemp(cl.Target)
			}
			components.Text(c, value).FontSize(68).FontWeight(300).TextColor(ink).LetterSpacing(-2)
			components.Text(c, a.unit).FontSize(20).FontWeight(500).TextColor(accent).Margin(0, 0, 12, 0).MarginStart(4)
		})
		ui.Spacer(c)
		eco := cl.Preset == "eco"
		pill := ui.Row(c).AlignItems(ui.Center).Gap(10).Padding(10, 12, 10, 14).Radius(18).
			Background(tileFace).Border(1, tileEdge)
		pill.Children(func() {
			components.Icon(c, lucide.Leaf).Size(18, 18).TextColor(inkMuted)
			components.Text(c, locales.Eco()).FontSize(14).FontWeight(600).TextColor(ink)
			sw := components.Switch(c, &eco, locales.Eco(), components.SwitchOptions{HideLabel: true})
			sw.Disabled(!cl.HasPreset("eco"))
			if sw.Changed() {
				a.togglePreset("eco")
			}
		})
		if !cl.HasPreset("eco") {
			pill.Opacity(0.45)
		}
	})
}

func (a *application) sliderRow(c *ui.Context) {
	cl := a.climate
	step := cl.Step
	ui.Row(c).FillWidth().AlignItems(ui.Start).Gap(12).Children(func() {
		roundButton(c, lucide.Minus, locales.Lower(), cl.Target <= cl.Min, func() { a.setTarget(cl.Target - step) })
		ui.Column(c).Grow(1).Basis(0).Gap(8).Children(func() {
			a.temperatureSlider(c)
			ui.Row(c).FillWidth().PaddingX(10).Children(func() {
				components.Text(c, formatTemp(cl.Min)+"°").FontSize(13).TextColor(inkMuted)
				ui.Spacer(c)
				components.Text(c, formatTemp(cl.Max)+"°").FontSize(13).TextColor(inkMuted)
			})
		})
		roundButton(c, lucide.Plus, locales.Raise(), cl.Target >= cl.Max, func() { a.setTarget(cl.Target + step) })
	})
}

// roundButton is a circular icon button beside the slider.
func roundButton(c *ui.Context, icon func() *ui.SVG, label string, disabled bool, click func()) {
	b := ui.ButtonBase(c).Size(44, 44).Radius(22).Label(label).Disabled(disabled).Margin(6, 0, 0, 0).Center()
	face := ui.RGBA(255, 255, 255, 0.06)
	switch {
	case disabled:
	case b.Pressed():
		face = ui.RGBA(255, 255, 255, 0.14)
	case b.Hovered():
		face = ui.RGBA(255, 255, 255, 0.1)
	}
	b.Background(face).Border(1, tileEdge)
	b.Children(func() { components.Icon(c, icon).Size(18, 18).TextColor(ink) })
	if disabled {
		b.Opacity(0.4)
	}
	components.Tooltip(c, b, label)
	if b.Clicked() {
		click()
	}
}

// temperatureSlider is a dark track with ticks and a white thumb.
func (a *application) temperatureSlider(c *ui.Context) {
	cl := a.climate
	const height, inset = 56, 4
	thumb := float32(height - 2*inset)
	value := cl.Target
	s := ui.SliderBase(c, &value, cl.Min, cl.Max).Step(cl.Step).Label(locales.Target()).
		FillWidth().Height(height).PaddingX(inset + thumb/2).FocusRing(false)
	if s.Changed() {
		a.setTarget(value)
	}
	_, _, held := s.Dragged()
	grow := s.Animate("held", map[bool]float32{false: 0, true: 1}[held], 150*time.Millisecond)
	focus := s.FocusVisible()
	fraction := float32(0)
	if cl.Max > cl.Min {
		fraction = float32((value - cl.Min) / (cl.Max - cl.Min))
	}
	ticks := 0
	if cl.Step > 0 {
		ticks = int(math.Round((cl.Max - cl.Min) / cl.Step))
	}
	for ticks > 28 {
		ticks /= 2
	}
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, ui.RGBA(0, 0, 0, 0.42), 18)
		p.Stroke(r, ui.RGBA(255, 255, 255, 0.05), 18, 1)
		left, width := r.X+inset+thumb/2, r.W-2*inset-thumb
		x := left + width*fraction
		for i := 0; i <= ticks && ticks > 0; i++ {
			tx := left + width*float32(i)/float32(ticks)
			if tx < x+thumb/2+4 {
				continue
			}
			h := float32(8)
			if i%5 == 0 {
				h = 12
			}
			p.Fill(ui.Rect{X: tx - 0.75, Y: r.Y + r.H/2 - h/2, W: 1.5, H: h}, ui.RGBA(255, 255, 255, 0.22), 0.75)
		}
		size := thumb * (1 + 0.06*grow)
		k := ui.Rect{X: x - size/2, Y: r.Y + r.H/2 - size/2, W: size, H: size}
		p.Shadow(k, 15, 0, 0, 14, 1, accent.Alpha(0.55))
		p.Fill(ui.Rect{X: k.X - 3, Y: k.Y - 3, W: k.W + 6, H: k.H + 6}, accent.Alpha(0.9), 17)
		p.Fill(k, ui.RGB(255, 255, 255), 14)
		if focus {
			p.FocusRing(ui.Rect{X: k.X - 3, Y: k.Y - 3, W: k.W + 6, H: k.H + 6}, [4]float32{17, 17, 17, 17})
		}
	})
}

// tiles are the four controls along the bottom.
func (a *application) tiles(c *ui.Context) {
	cl := a.climate
	ui.Row(c).FillWidth().Gap(8).Margin(18, 0, 0, 0).Children(func() {
		mode := cl.Mode
		a.menuTile(c, modeIcon(mode), locales.Mode(), optionText("mode", mode), "mode", mode, cl.ActiveModes(), a.setMode)
		a.menuTile(c, lucide.Fan, locales.FanSpeed(), optionText("fan", cl.FanMode), "fan", cl.FanMode, cl.FanModes, a.setFan)
		a.sleepTile(c, cl)
		a.menuTile(c, lucide.Wind, locales.Airflow(), optionText("swing", cl.SwingMode), "swing", cl.SwingMode, cl.SwingModes, a.setSwing)
	})
}

func tileBase(c *ui.Context, label string, disabled bool) ui.Element {
	t := ui.ButtonBase(c).Grow(1).Basis(0).Height(104).Radius(20).Padding(12, 8, 12, 10).
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

func tileText(c *ui.Context, title, value string) {
	components.Text(c, title).FontSize(12).TextColor(inkMuted).SingleLine().Ellipsis("…")
	components.Text(c, value).FontSize(12).Bold().TextColor(ink).SingleLine().Ellipsis("…").Margin(2, 0, 0, 0)
}

// menuTile opens a menu of a Home Assistant option list.
func (a *application) menuTile(c *ui.Context, icon func() *ui.SVG, title, value, kind, current string, options []string, choose func(string)) {
	t := tileBase(c, title+": "+value, len(options) == 0 || !a.climate.Available())
	t.Children(func() {
		ui.Row(c).FillWidth().Children(func() {
			components.Icon(c, icon).Size(26, 26).TextColor(accent.Mix(ink, 0.15))
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

func (a *application) sleepTile(c *ui.Context, cl ha.Climate) {
	on := cl.Preset == "sleep"
	value := locales.Off()
	if on {
		value = locales.On()
	}
	t := tileBase(c, locales.Sleep()+": "+value, !cl.HasPreset("sleep") || !cl.Available()).Checked(on)
	t.Children(func() {
		ui.Row(c).FillWidth().Children(func() {
			color := inkMuted
			if on {
				color = accent.Mix(ink, 0.15)
			}
			components.Icon(c, lucide.Moon).Size(26, 26).TextColor(color)
			ui.Spacer(c)
			dot := ui.RGBA(255, 255, 255, 0.3)
			if on {
				dot = accent
			}
			d := ui.Box(c).Size(8, 8).Radius(4).Background(dot).Margin(4, 0, 0, 0).MarginEnd(3)
			if on {
				d.Shadow(0, 0, 8, 0, accent.Alpha(0.8))
			}
		})
		ui.Spacer(c)
		tileText(c, locales.Sleep(), value)
	})
	if t.Clicked() {
		a.togglePreset("sleep")
	}
}
