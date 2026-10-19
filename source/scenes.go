package main

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/neohm/source/internal/ha"
)

// sceneHeight is the height of the drawing at the top of a device's page
// at least, as tall as an air conditioner's unit and the air it blows.
const sceneHeight = 142

// Surfaces of the devices drawn, white as the air conditioner's unit.
var (
	bodyTop    = ui.RGB(252, 253, 255)
	bodyBottom = ui.RGB(214, 221, 236)
	bodyEdge   = ui.RGB(132, 138, 152)
)

// deviceScene draws the device of entity s at the top of its page, lit in
// its color while it is on, open or at work. amount, 0 to 1, is how far its
// level goes: a light's brightness, a fan's speed or how far a cover is
// open.
func deviceScene(c *ui.Context, s ha.State, v tileView, amount float64) {
	// The drawing takes the room the controls leave, in its middle.
	box := ui.Box(c.Key(s.EntityID)).FillWidth().Grow(1).MinHeight(sceneHeight).Shrink(0)
	lit := float32(0)
	if v.lit {
		lit = 1
	}
	lit = box.Animate("lit", lit, 400*time.Millisecond)
	more := box.Animate("amount", float32(amount), 400*time.Millisecond)
	still := c.Preferences().ReduceMotion
	if !s.Available() {
		box.Opacity(0.45)
	}
	color := v.color
	switch s.Domain() {
	case "light":
		box.Draw(func(p *ui.Painter, r ui.Rect) { drawLamp(p, r, color, lit, more) })
	case "fan":
		spin := ui.Local(box, "spin", func() fanSpin { return fanSpin{} })
		box.Draw(func(p *ui.Painter, r ui.Rect) { drawFan(p, r, color, lit, more, spin, still) })
	case "cover":
		curtain := s.Str("device_class") == "curtain"
		box.Draw(func(p *ui.Painter, r ui.Rect) { drawWindow(p, r, color, more, curtain) })
	case "lock":
		box.Draw(func(p *ui.Painter, r ui.Rect) { drawLock(p, r, color, lit) })
	case "vacuum":
		spin := ui.Local(box, "spin", func() fanSpin { return fanSpin{} })
		box.Draw(func(p *ui.Painter, r ui.Rect) { drawVacuum(p, r, color, lit, spin, still) })
	case "media_player":
		playing := s.State == "playing"
		title := s.Str("media_title")
		if s.Str("device_class") == "tv" {
			box.Draw(func(p *ui.Painter, r ui.Rect) { drawScreen(p, r, color, lit, playing, title, still) })
		} else {
			box.Draw(func(p *ui.Painter, r ui.Rect) { drawSpeaker(p, r, color, lit, playing, still) })
		}
	default:
		box.Center().Draw(func(p *ui.Painter, r ui.Rect) { drawRipples(p, r, color, lit, still) })
		box.Children(func() {
			deviceDisc(c, v, lit, 76).Shadow(0, 10, 36, 0, color.Alpha(0.45*lit))
		})
	}
}

// stage is the middle of r, sceneHeight high, where a device is drawn.
func stage(r ui.Rect) ui.Rect {
	return ui.Rect{X: r.X, Y: r.Y + (r.H-sceneHeight)/2, W: r.W, H: sceneHeight}
}

// drawLamp draws a pendant lamp hanging from the top of r, shining in color, as bright as lit and
// amount make it.
func drawLamp(p *ui.Painter, r ui.Rect, color ui.Color, lit, amount float32) {
	ceiling := r.Y
	r = stage(r)
	cx := r.X + r.W/2
	top := r.Y + 18
	const w, h = 112, 44
	bright := lit * (0.3 + 0.7*amount)
	rim := top + h
	// The light falls in a cone to the floor.
	var cone ui.Path
	cone.MoveTo(cx-w/2+8, rim).LineTo(cx+w/2-8, rim).LineTo(r.X+r.W*0.92, r.Y+r.H).LineTo(r.X+r.W*0.08, r.Y+r.H).Close()
	p.FillPathGradient(&cone, ui.LinearGradient{From: color.Alpha(0.32 * bright), To: color.Alpha(0), Angle: 180})
	p.Shadow(ui.Rect{X: cx - 90, Y: r.Y + r.H - 8, W: 180, H: 6}, 3, 0, 0, 22, 0, color.Alpha(0.3*bright))
	p.Line(cx, ceiling, cx, top+2, 2, ui.RGBA(255, 255, 255, 0.35))
	// The bulb shows under the shade.
	bulb := ui.Rect{X: cx - 13, Y: rim - 13, W: 26, H: 26}
	p.Shadow(bulb, 13, 0, 0, 28, 4, color.Alpha(0.8*bright))
	p.Fill(bulb, ui.RGB(120, 128, 148).Mix(color.Mix(ui.RGB(255, 255, 255), 0.4), lit), 13)
	var shade ui.Path
	shade.MoveTo(cx-14, top).LineTo(cx+14, top).
		CubeTo(cx+38, top+4, cx+w/2, top+22, cx+w/2, rim).LineTo(cx-w/2, rim).
		CubeTo(cx-w/2, top+22, cx-38, top+4, cx-14, top).Close()
	p.FillPathGradient(&shade, ui.LinearGradient{From: bodyTop, To: bodyBottom, Angle: 180})
	p.Fill(ui.Rect{X: cx - w/2, Y: rim - 3, W: w, H: 3}, bodyEdge, 1.5)
}

// fanSpin is how far a fan's blades have turned, and when.
type fanSpin struct {
	angle float64
	at    time.Time
}

// drawFan draws a fan whose blades turn the faster the higher amount sets
// its speed, while lit.
func drawFan(p *ui.Painter, r ui.Rect, color ui.Color, lit, amount float32, spin *fanSpin, still bool) {
	r = stage(r)
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	const radius = 62
	now := p.Now()
	speed := float64(lit) * (90 + 630*float64(amount)) // degrees a second
	if !spin.at.IsZero() && !still {
		spin.angle = math.Mod(spin.angle+now.Sub(spin.at).Seconds()*speed, 360)
	}
	spin.at = now
	if speed > 0 && !still {
		p.AnimationFrame()
	}
	face := ui.Rect{X: cx - radius, Y: cy - radius, W: 2 * radius, H: 2 * radius}
	p.Shadow(face, radius, 0, 0, 30, 0, color.Alpha(0.35*lit))
	p.Fill(face, ui.RGBA(255, 255, 255, 0.05).Mix(color.Alpha(0.14), lit*amount), radius)
	for _, g := range []float32{radius * 0.45, radius * 0.75} {
		p.Stroke(ui.Rect{X: cx - g, Y: cy - g, W: 2 * g, H: 2 * g}, ui.RGBA(255, 255, 255, 0.1), g, 1)
	}
	blade := ui.RGB(226, 232, 244).Mix(color.Mix(ui.RGB(255, 255, 255), 0.45), lit)
	for i := range 3 {
		p.Transform(ui.Rotation(float32(spin.angle)+float32(i)*120), cx, cy, func() {
			var b ui.Path
			b.MoveTo(cx, cy-8).CubeTo(cx+28, cy-18, cx+30, cy-52, cx+6, cy-56).CubeTo(cx-14, cy-56, cx-16, cy-30, cx, cy-8).Close()
			p.FillPath(&b, blade.Alpha(0.92))
		})
	}
	p.Stroke(face, ui.RGB(226, 232, 244).Alpha(0.75), radius, 3)
	hub := ui.Rect{X: cx - 13, Y: cy - 13, W: 26, H: 26}
	p.FillGradient(hub, ui.LinearGradient{From: bodyTop, To: bodyBottom, Angle: 180}, 13)
	p.Fill(ui.Rect{X: cx - 4, Y: cy - 4, W: 8, H: 8}, inkMuted.Mix(color, lit), 4)
}

// drawWindow draws a window behind a cover open by amount: curtains drawn
// to the sides, else a blind rolled up.
func drawWindow(p *ui.Painter, r ui.Rect, color ui.Color, amount float32, curtain bool) {
	r = stage(r)
	const w, h = 176, 112
	win := ui.Rect{X: r.X + (r.W-w)/2, Y: r.Y + 18, W: w, H: h}
	glow := 0.25 + 0.75*amount
	p.Shadow(win, 10, 0, 0, 34, 0, ui.RGB(150, 190, 255).Alpha(0.3*glow))
	p.FillGradient(win, ui.LinearGradient{From: ui.RGB(110, 160, 245).Alpha(0.2 + 0.5*glow), To: ui.RGB(255, 196, 150).Alpha(0.1 + 0.35*glow), Angle: 180}, 10)
	p.Line(win.X+w/2, win.Y, win.X+w/2, win.Y+h, 2, ui.RGBA(255, 255, 255, 0.3))
	p.Line(win.X, win.Y+h/2, win.X+w, win.Y+h/2, 2, ui.RGBA(255, 255, 255, 0.3))
	covered := 1 - amount
	fabric := color.Mix(ui.RGB(255, 255, 255), 0.25)
	if curtain {
		// Each curtain keeps a little gathered at its side.
		cw := (w/2 + 14) * (0.16 + 0.84*covered)
		for _, left := range []bool{true, false} {
			x := win.X - 10
			if !left {
				x = win.X + w + 10 - cw
			}
			panel := ui.Rect{X: x, Y: win.Y - 8, W: cw, H: h + 18}
			p.FillGradient(panel, ui.LinearGradient{From: fabric.Alpha(0.95), To: fabric.Mix(panelBottom, 0.35), Angle: 180}, 4)
			folds := max(2, int(cw/14))
			for i := 1; i < folds; i++ {
				fx := panel.X + panel.W*float32(i)/float32(folds)
				p.Line(fx, panel.Y+4, fx, panel.Y+panel.H-2, 2, panelBottom.Alpha(0.18))
			}
		}
		p.Fill(ui.Rect{X: win.X - 22, Y: win.Y - 12, W: w + 44, H: 4}, bodyEdge, 2)
		return
	}
	p.Stroke(win, ui.RGBA(255, 255, 255, 0.55), 10, 3)
	p.Clip(win, 10, func() {
		down := h * covered
		for y := win.Y; y < win.Y+down; y += 9 {
			p.FillGradient(ui.Rect{X: win.X, Y: y, W: w, H: 7.5}, ui.LinearGradient{From: fabric, To: fabric.Mix(panelBottom, 0.3), Angle: 180}, 1.5)
		}
		p.Fill(ui.Rect{X: win.X, Y: win.Y + down - 3, W: w, H: 4}, bodyEdge, 2)
	})
	p.Fill(ui.Rect{X: win.X - 8, Y: win.Y - 10, W: w + 16, H: 10}, bodyBottom, 4)
}

// drawLock draws a padlock, its shackle raised and glowing in color while
// lit, unlocked.
func drawLock(p *ui.Painter, r ui.Rect, color ui.Color, lit float32) {
	r = stage(r)
	cx := r.X + r.W/2
	const w, h = 84, 64
	body := ui.Rect{X: cx - w/2, Y: r.Y + r.H - h - 12, W: w, H: h}
	lift := 18 * lit
	var shackle ui.Path
	sx, top := float32(24), body.Y-40-lift
	shackle.MoveTo(cx-sx, body.Y+4-lift*0.2).LineTo(cx-sx, top+sx).
		CubeTo(cx-sx, top-6, cx+sx, top-6, cx+sx, top+sx).LineTo(cx+sx, body.Y+4-lift)
	p.Shadow(body, 16, 0, 0, 36, 0, color.Alpha(0.5*lit))
	p.StrokePath(&shackle, 10, ui.RGB(196, 204, 222))
	p.Shadow(body, 16, 0, 8, 20, 0, ui.RGBA(0, 0, 0, 0.35))
	p.FillGradient(body, ui.LinearGradient{From: bodyTop, To: bodyBottom, Angle: 180}, 16)
	hole := ui.RGB(60, 68, 90).Mix(color, lit)
	p.Fill(ui.Rect{X: cx - 7, Y: body.Y + 20, W: 14, H: 14}, hole, 7)
	p.Fill(ui.Rect{X: cx - 3, Y: body.Y + 28, W: 6, H: 16}, hole, 3)
}

// drawVacuum draws a robot vacuum from above, its side brushes turning and
// its light in color while lit, at work.
func drawVacuum(p *ui.Painter, r ui.Rect, color ui.Color, lit float32, spin *fanSpin, still bool) {
	r = stage(r)
	cx, cy := r.X+r.W/2, r.Y+r.H/2+4
	const radius = 58
	now := p.Now()
	if !spin.at.IsZero() && !still {
		spin.angle = math.Mod(spin.angle+now.Sub(spin.at).Seconds()*540*float64(lit), 360)
	}
	spin.at = now
	if lit > 0 && !still {
		p.AnimationFrame()
	}
	body := ui.Rect{X: cx - radius, Y: cy - radius, W: 2 * radius, H: 2 * radius}
	p.Shadow(body, radius, 0, 0, 34, 0, color.Alpha(0.4*lit))
	// The side brushes reach out under the front of the body.
	for _, side := range []float32{-1, 1} {
		bx, by := cx+side*40, cy-38
		p.Transform(ui.Rotation(float32(spin.angle)*side), bx, by, func() {
			for k := range 3 {
				p.Transform(ui.Rotation(float32(k)*120), bx, by, func() {
					p.Line(bx, by, bx, by-22, 2, ui.RGBA(255, 255, 255, 0.5))
				})
			}
		})
	}
	p.Shadow(body, radius, 0, 8, 20, 0, ui.RGBA(0, 0, 0, 0.35))
	p.FillGradient(body, ui.LinearGradient{From: bodyTop, To: bodyBottom, Angle: 180}, radius)
	// The bumper along the front, and the turret of its sensor.
	var bumper ui.Path
	bumper.MoveTo(cx-radius+8, cy-18).CubeTo(cx-radius+14, cy-radius-4, cx+radius-14, cy-radius-4, cx+radius-8, cy-18)
	p.StrokePath(&bumper, 3, bodyEdge.Alpha(0.6))
	turret := ui.Rect{X: cx - 17, Y: cy - 17, W: 34, H: 34}
	p.FillGradient(turret, ui.LinearGradient{From: ui.RGB(70, 78, 100), To: ui.RGB(36, 42, 62), Angle: 180}, 17)
	p.Fill(ui.Rect{X: cx - 5, Y: cy + 26, W: 10, H: 10}, ui.RGB(160, 168, 186).Mix(color, lit), 5)
	if lit > 0 {
		p.Shadow(ui.Rect{X: cx - 5, Y: cy + 26, W: 10, H: 10}, 5, 0, 0, 10, 0, color.Alpha(0.9*lit))
	}
}

// drawScreen draws a television, its screen lit in color while on, with a
// moving equalizer and the title while it plays.
func drawScreen(p *ui.Painter, r ui.Rect, color ui.Color, lit float32, playing bool, title string, still bool) {
	r = stage(r)
	cx := r.X + r.W/2
	const w, h = 184, 108
	tv := ui.Rect{X: cx - w/2, Y: r.Y + 6, W: w, H: h}
	p.Shadow(tv, 12, 0, 0, 40, 0, color.Alpha(0.4*lit))
	p.Fill(tv, ui.RGB(18, 24, 42), 12)
	p.Stroke(tv, ui.RGBA(255, 255, 255, 0.14), 12, 1)
	screen := ui.Rect{X: tv.X + 5, Y: tv.Y + 5, W: w - 10, H: h - 10}
	p.FillGradient(screen, ui.LinearGradient{From: ui.RGB(10, 14, 28).Mix(color, 0.55*lit), To: ui.RGB(10, 14, 28).Mix(lavender, 0.4*lit), Angle: 135}, 8)
	if lit > 0 && title != "" {
		label := ui.Span{Text: title, Size: 12, Weight: 600, Color: ink.Alpha(lit)}
		p.Clip(screen, 8, func() { p.RichText(screen.X+12, screen.Y+10, screen.W-24, label) })
	}
	if lit > 0 {
		equalizer(p, ui.Rect{X: screen.X + 12, Y: screen.Y + screen.H - 34, W: 60, H: 24}, ink.Alpha(0.85*lit), playing, still)
	}
	p.Fill(ui.Rect{X: cx - 4, Y: tv.Y + h, W: 8, H: 12}, bodyEdge, 1)
	p.Fill(ui.Rect{X: cx - 34, Y: tv.Y + h + 11, W: 68, H: 4}, bodyBottom, 2)
}

// drawSpeaker draws a speaker, its cones pulsing in color while it plays.
func drawSpeaker(p *ui.Painter, r ui.Rect, color ui.Color, lit float32, playing bool, still bool) {
	r = stage(r)
	cx := r.X + r.W/2
	const w, h = 88, 128
	box := ui.Rect{X: cx - w/2, Y: r.Y + (r.H-h)/2, W: w, H: h}
	p.Shadow(box, 18, 0, 0, 36, 0, color.Alpha(0.4*lit))
	p.FillGradient(box, ui.LinearGradient{From: ui.RGB(48, 56, 80), To: ui.RGB(24, 30, 50), Angle: 180}, 18)
	p.Stroke(box, ui.RGBA(255, 255, 255, 0.14), 18, 1)
	pulse := float32(0)
	if playing && !still {
		pulse = float32(math.Abs(math.Sin(float64(p.Now().UnixMilli()) / 1000 * math.Pi * 2)))
		p.AnimationFrame()
	}
	for i, cone := range []struct{ y, r float32 }{{box.Y + 32, 14}, {box.Y + 84, 26}} {
		rr := cone.r * (1 + 0.05*pulse*float32(i))
		k := ui.Rect{X: cx - rr, Y: cone.y - rr, W: 2 * rr, H: 2 * rr}
		p.Fill(k, ui.RGB(14, 18, 32), rr)
		p.Stroke(k, ui.RGBA(255, 255, 255, 0.18).Mix(color, lit), rr, 2)
		p.Fill(ui.Rect{X: cx - rr*0.3, Y: cone.y - rr*0.3, W: rr * 0.6, H: rr * 0.6}, ui.RGB(70, 78, 100).Mix(color, lit*0.6), rr*0.3)
	}
}

// equalizer draws bars in r that dance while playing.
func equalizer(p *ui.Painter, r ui.Rect, color ui.Color, playing, still bool) {
	const bars = 5
	now := float64(p.Now().UnixMilli()) / 1000
	bw := r.W / (bars*2 - 1)
	for i := range bars {
		f := float32(0.25)
		if playing {
			f = 0.35 + 0.3*float32(math.Sin(float64(i)*1.7)+1)
			if !still {
				f = 0.2 + 0.8*float32(math.Abs(math.Sin(now*(2.2+0.5*float64(i))+float64(i))))
			}
		}
		bh := r.H * f
		p.Fill(ui.Rect{X: r.X + float32(i)*bw*2, Y: r.Y + r.H - bh, W: bw, H: bh}, color, bw/2)
	}
	if playing && !still {
		p.AnimationFrame()
	}
}

// drawRipples draws rings spreading from the middle of r in color while
// lit, behind the device's disc.
func drawRipples(p *ui.Painter, r ui.Rect, color ui.Color, lit float32, still bool) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	for _, g := range []float32{52, 66} {
		p.Stroke(ui.Rect{X: cx - g, Y: cy - g, W: 2 * g, H: 2 * g}, ui.RGBA(255, 255, 255, 0.05), g, 1)
	}
	if lit <= 0 {
		return
	}
	now := float64(p.Now().UnixMilli()) / 1000
	for k := range 3 {
		t := float32(math.Mod(now/2.4+float64(k)/3, 1))
		if still {
			t = float32(k+1) / 4
		}
		g := 40 + 34*t
		p.Stroke(ui.Rect{X: cx - g, Y: cy - g, W: 2 * g, H: 2 * g}, color.Alpha(0.5*lit*(1-t)), g, 1.5)
	}
	if !still {
		p.AnimationFrame()
	}
}
