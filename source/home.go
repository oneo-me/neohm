package main

import (
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/neohm/source/internal/ha"
	"github.com/oneo-me/neohm/source/locales"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// tileHeight is the height of a tile, about the width of a third of the
// panel.
const tileHeight = 112

// homePage is where the panel opens: Home Assistant's devices as tiles, by
// area. radius rounds the bottom corners the tiles scroll within.
func (a *application) homePage(c *ui.Context, radius float32) {
	pageHeader(c, locales.Devices(), func() {
		if circleButton(c, lucide.Settings, locales.SettingsTitle(), 36, false).Clicked() {
			a.showSettings()
		}
	})
	if !a.configured() || !a.listed {
		a.setupView(c)
		return
	}
	// A failure after the devices were read leaves them as last read.
	if e := a.errorDisplay(); e != "" {
		ui.Row(c).FillWidth().Shrink(0).Gap(8).Padding(0, 22, 10, 22).AlignItems(ui.Center).Children(func() {
			components.Icon(c, lucide.TriangleAlert).Size(16, 16).TextColor(danger).Shrink(0)
			components.Text(c, e).FontSize(12).Wrap().TextColor(inkMuted).Grow(1).MinWidth(0)
		})
	}
	if len(a.home) == 0 {
		ui.Column(c).Fill().Center().Gap(10).Padding(0, 24, 48, 24).Children(func() {
			components.Icon(c, lucide.House).Size(44, 44).TextColor(inkMuted)
			components.Text(c, locales.NoDevices()).FontSize(17).Bold().TextColor(ink)
			components.Text(c, locales.NoDevicesBody()).Wrap().TextAlign(ui.Center).FontSize(13).TextColor(inkMuted).MaxWidth(260)
		})
		return
	}
	components.ScrollView(c, components.ScrollViewOptions{}).Fill().Radius(0, 0, radius, radius).
		ScrollbarInsets(0, 4, radius).Children(func() {
		ui.Column(c).FillWidth().Padding(2, 18, 18, 18).Gap(18).Children(func() {
			for i, s := range a.home {
				ui.Column(c.Key(i)).FillWidth().Gap(10).Children(func() {
					switch {
					case s.Area != "":
						areaHeading(c, s.Area)
					case len(a.home) > 1:
						areaHeading(c, locales.OtherDevices())
					}
					ui.Grid(c).Columns(3).Gap(8).Children(func() {
						for _, t := range s.Tiles {
							a.tile(c, t, s.Area)
						}
					})
				})
			}
		})
	})
}

func areaHeading(c *ui.Context, name string) {
	components.Text(c, name).FontSize(13).Bold().TextColor(inkMuted).SingleLine().Ellipsis("…").Margin(0, 4)
}

// setupView stands in for the tiles until Home Assistant lists the devices:
// it waits, asks for the connection or says what went wrong.
func (a *application) setupView(c *ui.Context) {
	ui.Column(c).Fill().Center().Gap(14).Padding(0, 24, 48, 24).Children(func() {
		components.Icon(c, lucide.House).Size(48, 48).TextColor(accent)
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
		components.Button(c, locales.Settings(), components.ButtonPrimary).OnClick(a.showSettings)
	})
}

// tile shows one device: what it is, its name and its main state, lit in
// its color while it is on, open or at work. A click does what the device
// is for: it switches a light, opens a menu of a lock's actions or opens an
// air conditioner's controls.
func (a *application) tile(c *ui.Context, t ha.Tile, area string) {
	s := a.held(t.Entity)
	v := a.tileView(t, s)
	name := tileName(t.Name, area)
	var e ui.Element
	if v.kind == tileStatic {
		e = ui.Box(c.Key(s.EntityID))
	} else {
		e = ui.ButtonBase(c.Key(s.EntityID))
	}
	disabled := v.kind != tileStatic && v.kind != tilePage && (!s.Available() || len(v.actions) == 0)
	e.Height(tileHeight).Radius(20).Padding(12).Column().AlignItems(ui.Start).Label(name + ", " + v.status).Disabled(disabled)
	if v.kind == tileSwitch {
		e.Role(ui.RoleToggleButton).Checked(v.lit)
	}
	level := float32(0)
	if v.lit {
		level = 1
	}
	level = e.Animate("lit", level, 220*time.Millisecond)
	// The face follows the pointer, as the panel's other tiles do.
	offA, litA := float32(0.045), float32(0.2)
	if v.kind != tileStatic && !disabled {
		switch {
		case e.Pressed():
			offA, litA = 0.12, 0.15
		case e.Hovered():
			offA, litA = 0.08, 0.26
		}
	}
	off := ui.RGBA(255, 255, 255, offA)
	e.Gradient(off.Mix(v.color.Alpha(litA+0.1), level), off.Mix(v.color.Alpha(litA-0.08), level), 135).
		Border(1, tileEdge.Mix(v.color.Alpha(0.45), level)).
		Shadow(0, 8, 24, 0, v.color.Alpha(0.28*level))
	if !s.Available() {
		e.Opacity(0.45)
	}
	e.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Start).Children(func() {
			disc, glyph := ui.RGBA(255, 255, 255, 0.07).Mix(v.color, level), inkMuted.Mix(discInk, level)
			if v.neutral {
				disc, glyph = v.color.Alpha(0.16), v.color
			}
			ui.Box(c).Size(36, 36).Radius(18).Shrink(0).Background(disc).Center().Children(func() {
				components.Icon(c, v.icon).Size(19, 19).TextColor(glyph)
			})
			ui.Spacer(c)
			switch v.kind {
			case tileMenu:
				components.Icon(c, lucide.ChevronDown).Size(15, 15).TextColor(inkMuted).Margin(2, 0, 0, 0)
			case tilePage:
				components.Icon(c, lucide.ChevronRight).Size(15, 15).TextColor(inkMuted).Margin(2, 0, 0, 0)
			}
		})
		ui.Spacer(c)
		components.Text(c, name).FontSize(13).FontWeight(600).TextColor(ink).MaxLines(2)
		components.Text(c, v.status).FontSize(12).TextColor(inkMuted.Mix(v.color.Mix(ink, 0.4), level)).
			SingleLine().Ellipsis("…").Margin(2, 0, 0, 0)
	})
	switch v.kind {
	case tileSwitch, tilePress:
		if e.Clicked() {
			a.perform(s, v.actions[0])
		}
	case tileMenu:
		e.Menu(func(m *ui.Menu) {
			for _, act := range v.actions {
				if m.Item(actionText(act, s)).Checked(act.Current).Disabled(act.Disabled).Chosen() {
					a.perform(s, act)
				}
			}
		})
	case tilePage:
		if e.Clicked() {
			a.openClimate(ha.ParseClimate(t.Entity))
		}
	}
}

// tileName drops the area's name from the front of a device's, which the
// heading above already says: "Lamp" for "Bedroom Lamp" in the bedroom.
func tileName(name, area string) string {
	if area == "" || len(name) <= len(area) || !strings.EqualFold(name[:len(area)], area) {
		return name
	}
	if rest := strings.TrimLeft(name[len(area):], " -_·:："); rest != "" {
		return rest
	}
	return name
}
