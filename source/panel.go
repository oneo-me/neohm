package main

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/neohm/source/internal/ha"
	"github.com/oneo-me/neohm/source/locales"
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
	tileEdge    = ui.RGBA(255, 255, 255, 0.08)
	danger      = ui.RGB(255, 99, 99)
	// Devices glow in these while on; discInk draws on them.
	amber    = ui.RGB(255, 190, 70)
	green    = ui.RGB(95, 210, 140)
	cyan     = ui.RGB(90, 205, 245)
	teal     = ui.RGB(84, 214, 200)
	orange   = ui.RGB(255, 150, 82)
	lavender = ui.RGB(168, 140, 255)
	pink     = ui.RGB(255, 120, 180)
	discInk  = ui.RGB(14, 20, 38)
)

const (
	// panelRadius rounds the pages, which keeps them clear of the corners
	// the panel's window rounds.
	panelRadius = 26
	// pageFade is how far past its bar a page's header blurs the content
	// under it: a device's page starts below that.
	pageFade = 24
	// listGap is how far past its bar a scrolling page's content and scroll
	// bar start, where the blur no longer shows on text at rest.
	listGap = 16
)

// pageBar floats a page's header over the top of its content, which blurs
// under it the more the nearer the top.
var pageBar = components.ScrollEdgeBarOptions{Style: components.ScrollEdgeBlur, Inset: 8, Height: 48, Color: panelTop.Alpha(0.5)}

// panelTheme is Onevra's dark palette in the panel's colors, so that
// components such as Switch match it.
func panelTheme() *ui.Theme {
	t := theme.Palette(true)
	t.Background, t.Surface = panelBottom, tileFace
	t.Text, t.TextMuted = ink, inkMuted
	t.Accent, t.AccentHover, t.AccentPressed = accent, accent.Mix(ink, 0.12), accent.Mix(panelBottom, 0.15)
	t.AccentText = ink
	t.Border = ui.RGBA(255, 255, 255, 0.14)
	return t
}

// windowTheme is the panel's theme in a window, whose header continues the
// top of the panel's gradient.
func windowTheme(*ui.Context) *ui.Theme {
	t := panelTheme()
	t.Background = panelTop
	return t
}

// modeColor is the color of the air the unit blows in mode.
func modeColor(mode string) ui.Color {
	switch mode {
	case "heat":
		return orange
	case "dry":
		return teal
	case "fan_only":
		return ui.RGB(170, 190, 230)
	case "auto", "heat_cool":
		return lavender
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
	if c.Shortcut(0, ui.KeyEscape) && !a.back() {
		a.hidePanel()
	}
	if c.Shortcut(ui.Cmd, ui.KeyW) {
		a.hidePanel()
	}
	detail := a.showsDetail()
	unitTop := components.ScrollEdgeBarExtent(c, pageBar) + pageFade
	ui.Column(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
		// The surface, opaque: the blur under the pages' headers sees only
		// what the panel paints. The window rounds its corners.
		p.FillGradient(r, ui.LinearGradient{From: panelTop, To: panelBottom, Angle: 180}, 0)
		if detail {
			// A soft glow behind the unit.
			p.Shadow(ui.Rect{X: r.X + r.W/2 - 70, Y: r.Y + unitTop + 6, W: 140, H: 30}, 15, 0, 0, 60, 0, accent.Alpha(0.10))
		}
	}).Children(func() { a.page(c, panelRadius) })
}

// windowView shows the panel's pages in a window, for when there is no menu
// bar icon to anchor the panel to.
func (a *application) windowView(c *ui.Context) {
	a.scope(c)
	if c.Shortcut(0, ui.KeyEscape) {
		a.back()
	}
	if c.Shortcut(ui.Cmd, ui.KeyW) {
		a.window.Window().Close()
	}
	ui.Column(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
		p.FillGradient(r, ui.LinearGradient{From: panelTop, To: panelBottom, Angle: 180}, 0)
	}).Children(func() { a.page(c, 0) })
}

// page builds Settings, the air conditioner opened, the page of another
// device or the home page, whichever the panel shows. radius rounds the
// corners of the surface they show on.
func (a *application) page(c *ui.Context, radius float32) {
	t, area, device := a.tileOf(a.device)
	switch {
	case a.settingsOpen:
		a.settingsPage(c, radius)
	case a.showsDetail():
		a.climatePage(c, radius)
	case device:
		a.devicePage(c, t, area, radius)
	default:
		a.homePage(c, radius)
	}
}

func (a *application) showsDetail() bool { return a.detail && a.hasClimate && !a.settingsOpen }

// pageFrame builds a page whose content reaches the top of the panel, under
// the bar of its header, which label names. content fills the page and
// keeps its controls clear of the bar by top. The page, the header's blur
// included, keeps within the panel's corners, rounded by radius.
func pageFrame(c *ui.Context, radius float32, label string, header func(), content func(top float32)) {
	ui.Box(c).Fill().Radius(radius).Clip().Children(func() {
		content(components.ScrollEdgeBarExtent(c, pageBar))
		components.ScrollEdgeBar(c, pageBar, header).Label(label)
	})
}

// scrollContent scrolls a page's content across the whole panel, its scroll
// bar along the edge, both clear of the header below top and of the bottom
// corners, rounded by radius, and inset from the sides by side. s keeps how
// far it is scrolled, if not nil.
func scrollContent(c *ui.Context, s *ui.ScrollState, top, radius, side float32, content func()) {
	components.ScrollView(c, components.ScrollViewOptions{}).Fill().TrackScroll(s).ScrollbarInsets(top+listGap, 0, radius).Children(func() {
		ui.Column(c).FillWidth().Padding(top+listGap, side, 18, side).Gap(18).Children(content)
	})
}

// pageHeader titles a page in its bar, with its actions at the right.
func pageHeader(c *ui.Context, title string, actions func()) {
	components.Text(c, title).FontSize(20).Bold().TextColor(ink).SingleLine().Ellipsis("…").Grow(1).MinWidth(0).MarginStart(8)
	actions()
}

// backButton leads from a device's page back to the home page.
func (a *application) backButton(c *ui.Context) {
	if circleButton(c, lucide.ChevronLeft, locales.Back(), headerButton, false).Clicked() {
		a.closeDevice()
	}
}

// headerButton is how wide the buttons of a titled header are: the back
// button at its left and the power button at its right. headerSide is the
// width of each side.
const headerButton, headerSide = 36, 36

// titledHeader centers title, and what subtitle adds below it, between the
// back button and the controls trailing builds. Sides of one width keep it
// centered with or without them.
func (a *application) titledHeader(c *ui.Context, title string, subtitle, trailing func()) {
	ui.Row(c).Width(headerSide).Shrink(0).Children(func() { a.backButton(c) })
	ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(2).Children(func() {
		components.Text(c, title).FontSize(16).Bold().TextColor(ink).SingleLine().Ellipsis("…").FillWidth().TextAlign(ui.Center)
		subtitle()
	})
	ui.Row(c).Width(headerSide).Shrink(0).Justify(ui.End).Children(trailing)
}

// headerSubtitle is a line below a header's title.
func headerSubtitle(c *ui.Context, text string, color ui.Color) {
	components.Text(c, text).FontSize(12).TextColor(color).SingleLine().Ellipsis("…").FillWidth().TextAlign(ui.Center)
}

// climatePage controls the air conditioner opened. radius rounds the
// corners of the surface it shows on.
func (a *application) climatePage(c *ui.Context, radius float32) {
	pageFrame(c, radius, a.climate.Name, func() { a.climateHeader(c) }, func(top float32) {
		ui.Column(c).Fill().Padding(top+pageFade, 18, 18, 18).Children(func() {
			a.unitView(c)
			ui.Spacer(c)
			a.temperatureRow(c)
			a.sliderRow(c)
			a.tiles(c)
		})
	})
}

// climateHeader leads back to the home page and names the air conditioner
// and the room temperature, or what keeps it from being controlled, beside
// the power button.
func (a *application) climateHeader(c *ui.Context) {
	cl := a.climate
	a.titledHeader(c, cl.Name, func() {
		if problem := a.problem(cl.Available()); problem != "" {
			headerSubtitle(c, problem, danger)
			return
		}
		current := "—"
		if cl.HasCurrent {
			current = formatTemp(math.Round(cl.Current*10)/10) + a.unit
		}
		ui.Row(c).Justify(ui.Center).AlignItems(ui.Center).Gap(4).Children(func() {
			components.Icon(c, lucide.Thermometer).Size(13, 13).TextColor(accent)
			components.Text(c, locales.Room()).FontSize(12).TextColor(inkMuted)
			components.Text(c, current).FontSize(12).Bold().TextColor(ink)
		})
	}, func() { a.powerButton(c) })
}

// powerButton turns the air conditioner on or off, glowing in its mode's
// color while it is on.
func (a *application) powerButton(c *ui.Context) {
	cl := a.climate
	tint := accent
	if cl.On() {
		tint = modeColor(cl.Mode)
	}
	if powerButton(c, locales.Power(), lucide.Power, cl.On(), !cl.Available(), tint, 1).Clicked() {
		a.togglePower()
	}
}

// unitView draws the indoor unit and the air it blows. A click on it turns
// the air conditioner on or off, as its power button does.
func (a *application) unitView(c *ui.Context) {
	cl := a.climate
	on := cl.On()
	color := modeColor(cl.Mode)
	target := formatTemp(cl.Target)
	if !cl.HasTarget {
		target = "--"
	}
	box := ui.Box(c).FillWidth().Height(142)
	if cl.Available() {
		box.OnClick(a.togglePower)
	}
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
		u := ui.Rect{X: r.X + (r.W-uw)/2, Y: r.Y, W: uw, H: uh}
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

// temperatureRow shows the target temperature in the middle, and Eco Mode
// below it, which a click switches.
func (a *application) temperatureRow(c *ui.Context) {
	cl := a.climate
	value := "--"
	if cl.HasTarget {
		value = formatTemp(cl.Target)
	}
	ui.Column(c).FillWidth().AlignItems(ui.Center).Margin(0, 0, 14, 0).Children(func() {
		bigNumber(c, value, a.unit, accent)
		if toggleChip(c, lucide.Leaf, locales.Eco(), cl.Preset == "eco", !cl.HasPreset("eco") || !cl.Available(), green).Margin(-6, 0, 0, 0).Clicked() {
			a.togglePreset("eco")
		}
	})
}

func (a *application) sliderRow(c *ui.Context) {
	cl := a.climate
	sliderRow(c, slider{
		Value: cl.Target, Min: cl.Min, Max: cl.Max, Step: cl.Step,
		Label: locales.Target(), Lower: locales.Lower(), Raise: locales.Raise(),
		Color: modeColor(cl.Mode), Disabled: !cl.Available(),
		Bound: func(t float64) string { return formatTemp(t) + "°" },
		Set:   a.setTarget,
	})
}

// tiles are the four controls along the bottom.
func (a *application) tiles(c *ui.Context) {
	cl := a.climate
	off := !cl.Available()
	tileRows(c, 4, func(i int) {
		switch i {
		case 0:
			menuTile(c, modeIcon(cl.Mode), accent, locales.Mode(), optionText("mode", cl.Mode), "mode", cl.Mode, cl.ActiveModes(), off, a.setMode)
		case 1:
			menuTile(c, lucide.Fan, accent, locales.FanSpeed(), optionText("fan", cl.FanMode), "fan", cl.FanMode, cl.FanModes, off, a.setFan)
		case 2:
			a.sleepTile(c, cl)
		case 3:
			menuTile(c, lucide.Wind, accent, locales.Airflow(), optionText("swing", cl.SwingMode), "swing", cl.SwingMode, cl.SwingModes, off, a.setSwing)
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
				color = accent
			}
			tileIcon(c, lucide.Moon, color)
			ui.Spacer(c)
			tileDot(c, on, accent)
		})
		ui.Spacer(c)
		tileText(c, locales.Sleep(), value)
	})
	if t.Clicked() {
		a.togglePreset("sleep")
	}
}
