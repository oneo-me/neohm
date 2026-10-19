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

// devicePage shows the device of tile t, in area, on a page of its own: its
// state, the level it is set to, what else it does and what it reads, below
// a header with its quick switch. radius rounds the corners it scrolls
// within.
func (a *application) devicePage(c *ui.Context, t ha.Tile, area string, radius float32) {
	s := a.held(t.Entity)
	v := a.tileView(t, s)
	name := tileName(t.Name, area)
	pageFrame(c, radius, name, func() {
		a.backButton(c)
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(2).MarginStart(4).Children(func() {
			components.Text(c, name).FontSize(16).Bold().TextColor(ink).SingleLine().Ellipsis("…")
			if area != "" {
				components.Text(c, area).FontSize(12).TextColor(inkMuted).SingleLine().Ellipsis("…")
			}
		})
		a.quickSwitch(c, s, v, name).MarginEnd(8)
	}, func(top float32) {
		scrollContent(c, top, radius, func() {
			deviceHero(c, s, v)
			if l, ok := s.Level(); ok && s.Available() {
				a.levelCard(c, s, l)
			}
			a.actionGrid(c, s, v)
			a.readingList(c, t)
		})
	})
}

// deviceHero shows what the device is and its main state, glowing in its
// color while it is lit.
func deviceHero(c *ui.Context, s ha.State, v tileView) {
	hero := ui.Column(c).FillWidth().AlignItems(ui.Center).Gap(14).Padding(12, 0, 4, 0)
	level := float32(0)
	if v.lit {
		level = 1
	}
	level = hero.Animate("lit", level, 220*time.Millisecond)
	if !s.Available() {
		hero.Opacity(0.45)
	}
	hero.Children(func() {
		deviceDisc(c, v, level, 76).Shadow(0, 10, 36, 0, v.color.Alpha(0.45*level))
		components.Text(c, v.status).FontSize(18).Bold().TextColor(ink).Wrap().TextAlign(ui.Center)
	})
}

// card is a rounded surface holding one part of a device's page.
func card(c *ui.Context) ui.Element {
	return ui.Column(c).FillWidth().Radius(20).Background(tileFace).Border(1, tileEdge)
}

// levelCard sets level l of entity s, as a light's brightness, with a
// slider.
func (a *application) levelCard(c *ui.Context, s ha.State, l ha.Level) {
	value := a.level(s, l)
	title := message("level_" + l.Name)
	card(c).Gap(12).Padding(14, 16).Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			components.Text(c, title).FontSize(13).TextColor(inkMuted)
			ui.Spacer(c)
			shown := strconv.Itoa(int(math.Round(value))) + "%"
			if l.Name == "brightness" && !s.Active() {
				shown = stateText(s)
			}
			components.Text(c, shown).FontSize(13).Bold().TextColor(ink)
		})
		slider := components.Slider(c, &value, l.Min, l.Max, components.SliderOptions{Step: l.Step})
		slider.Label(title)
		slider.FillWidth()
		if slider.Changed() {
			a.setLevel(s, l, value)
		}
	})
}

// actionGrid offers what the device does besides what its quick switch
// does: a cover's stop, a media player's tracks, a select's options.
func (a *application) actionGrid(c *ui.Context, s ha.State, v tileView) {
	_, switches := s.Toggle()
	var actions []ha.Action
	for _, act := range s.Actions() {
		if !switches || act.Name != "turn_on" && act.Name != "turn_off" {
			actions = append(actions, act)
		}
	}
	if len(actions) == 0 {
		return
	}
	ui.Column(c).FillWidth().Gap(10).Children(func() {
		areaHeading(c, locales.Controls())
		ui.Grid(c).Columns(2).Gap(8).Children(func() {
			for _, act := range actions {
				if actionButton(c.Key(act.Name), actionText(act, s), v.color, act.Current, act.Disabled || !s.Available()).Clicked() {
					a.perform(s, act)
				}
			}
		})
	})
}

// actionButton is a button of a device's page, lit in color while the
// device is where it leads.
func actionButton(c *ui.Context, label string, color ui.Color, current, disabled bool) ui.Element {
	b := ui.ButtonBase(c).Height(46).Radius(16).PaddingX(12).Center().Label(label).Checked(current).Disabled(disabled)
	face, edge := tileFace, tileEdge
	switch {
	case disabled:
	case b.Pressed():
		face = ui.RGBA(255, 255, 255, 0.12)
	case b.Hovered():
		face = tileHover
	}
	if current {
		face, edge = color.Alpha(0.22), color.Alpha(0.5)
	}
	b.Background(face).Border(1, edge)
	b.Children(func() {
		components.Text(c, label).FontSize(13).FontWeight(600).TextColor(ink).SingleLine().Ellipsis("…")
	})
	if disabled {
		b.Opacity(0.45)
	}
	return b
}

// readingList shows each reading of a device that only measures, named
// without the device's name.
func (a *application) readingList(c *ui.Context, t ha.Tile) {
	if d := t.Entity.Domain(); d != "sensor" && d != "binary_sensor" {
		return
	}
	ui.Column(c).FillWidth().Gap(10).Children(func() {
		areaHeading(c, locales.Readings())
		card(c).PaddingX(16).Children(func() {
			for i, r := range append([]ha.State{t.Entity}, t.Readings...) {
				if i > 0 {
					ui.Box(c.Key(i)).FillWidth().Height(1).Background(tileEdge)
				}
				ui.Row(c.Key(r.EntityID)).FillWidth().AlignItems(ui.Center).Gap(12).Padding(13, 0).Children(func() {
					components.Text(c, tileName(r.Name(), t.Name)).FontSize(13).TextColor(inkMuted).SingleLine().Ellipsis("…").Grow(1).MinWidth(0)
					components.Text(c, a.readingText(r)).FontSize(13).Bold().TextColor(ink).Shrink(0)
				})
			}
		})
	})
}
