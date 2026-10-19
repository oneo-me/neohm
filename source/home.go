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
// area. radius rounds the corners the tiles scroll within.
func (a *application) homePage(c *ui.Context, radius float32) {
	pageFrame(c, radius, locales.Title(), func() {
		pageHeader(c, locales.Title(), func() {
			ui.Row(c).Gap(8).Shrink(0).Children(func() {
				a.updateButton(c)
				if circleButton(c, lucide.Settings, locales.SettingsTitle(), 36, false).Clicked() {
					a.showSettings()
				}
			})
		})
	}, func(top float32) {
		switch {
		case !a.configured() || !a.listed:
			a.setupView(c, top)
		case len(a.home) == 0:
			ui.Column(c).Fill().Center().Gap(10).Padding(top, 24, 48, 24).Children(func() {
				components.Icon(c, lucide.House).Size(44, 44).TextColor(inkMuted)
				components.Text(c, locales.NoDevices()).FontSize(17).Bold().TextColor(ink)
				components.Text(c, locales.NoDevicesBody()).Wrap().TextAlign(ui.Center).FontSize(13).TextColor(inkMuted).MaxWidth(260)
			})
		default:
			scrollContent(c, top, radius, 18, func() {
				// A failure after the devices were read leaves them as last
				// read.
				if e := a.errorDisplay(); e != "" {
					ui.Row(c).FillWidth().Gap(8).Padding(0, 4).AlignItems(ui.Center).Children(func() {
						components.Icon(c, lucide.TriangleAlert).Size(16, 16).TextColor(danger).Shrink(0)
						components.Text(c, e).FontSize(12).Wrap().TextColor(inkMuted).Grow(1).MinWidth(0)
					})
				}
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
		}
	})
}

func areaHeading(c *ui.Context, name string) ui.Element {
	return components.Text(c, name).FontSize(13).Bold().TextColor(inkMuted).SingleLine().Ellipsis("…").Margin(0, 4)
}

// setupView stands in for the tiles until Home Assistant lists the devices:
// it waits, asks for the connection or says what went wrong. top keeps it
// clear of the header.
func (a *application) setupView(c *ui.Context, top float32) {
	ui.Column(c).Fill().Center().Gap(14).Padding(top, 24, 48, 24).Children(func() {
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
// its color while it is on, open or at work. A click switches a device that
// switches on and off, or opens and closes, whose page then opens from the
// button in the corner; it opens the page of any other, an air
// conditioner's controls.
func (a *application) tile(c *ui.Context, t ha.Tile, area string) {
	s := a.held(t.Entity)
	v := a.tileView(t, s)
	name := tileName(t.Name, area)
	_, switches := s.Toggle()
	ui.Box(c.Key(s.EntityID)).Height(tileHeight).Children(func() {
		e := ui.ButtonBase(c).Fill().Radius(20).Padding(12).Column().AlignItems(ui.Start)
		if switches {
			e.Label(locales.QuickSwitch(name)).Disabled(!s.Available())
			e.Role(ui.RoleSwitch).Checked(v.lit).Value(v.status)
		} else {
			e.Label(name + ", " + v.status)
		}
		level := float32(0)
		if v.lit {
			level = 1
		}
		level = e.Animate("lit", level, 220*time.Millisecond)
		// The face follows the pointer, as the panel's other tiles do.
		offA, litA := float32(0.045), float32(0.2)
		switch {
		case e.Pressed():
			offA, litA = 0.12, 0.15
		case e.Hovered():
			offA, litA = 0.08, 0.26
		}
		off := ui.RGBA(255, 255, 255, offA)
		e.Gradient(off.Mix(v.color.Alpha(litA+0.1), level), off.Mix(v.color.Alpha(litA-0.08), level), 135).
			Border(1, tileEdge.Mix(v.color.Alpha(0.45), level)).
			Shadow(0, 8, 24, 0, v.color.Alpha(0.28*level))
		if !s.Available() {
			e.Opacity(0.45)
		}
		e.Children(func() {
			deviceDisc(c, v, level, 36)
			ui.Spacer(c)
			components.Text(c, name).FontSize(13).FontWeight(600).TextColor(ink).MaxLines(2)
			components.Text(c, v.status).FontSize(12).TextColor(inkMuted.Mix(v.color.Mix(ink, 0.4), level)).
				SingleLine().Ellipsis("…").Margin(2, 0, 0, 0)
		})
		switch {
		case !e.Clicked():
		case switches:
			a.toggle(s)
		default:
			a.open(t)
		}
		if switches {
			if circleButton(c, lucide.ChevronRight, locales.DeviceDetails(name), 28, false).Absolute().Top(16).Right(12).Clicked() {
				a.open(t)
			}
		}
	})
}

// deviceDisc is the round icon of a device, filled with its color while it
// is lit by level, size DIPs across.
func deviceDisc(c *ui.Context, v tileView, level, size float32) ui.Element {
	disc, glyph := ui.RGBA(255, 255, 255, 0.07).Mix(v.color, level), inkMuted.Mix(discInk, level)
	if v.neutral {
		disc, glyph = v.color.Alpha(0.16), v.color
	}
	icon := size * 19 / 36
	return ui.Box(c).Size(size, size).Radius(size / 2).Shrink(0).Background(disc).Center().Children(func() {
		components.Icon(c, v.icon).Size(icon, icon).TextColor(glyph)
	})
}

// quickSwitch switches the device of entity s, as held, on and off, or
// opens and closes it, if it does either; its element is empty otherwise.
func (a *application) quickSwitch(c *ui.Context, s ha.State, v tileView, name string) ui.Element {
	if _, ok := s.Toggle(); !ok {
		return ui.Element{}
	}
	on := v.lit
	sw := components.Switch(c, &on, locales.QuickSwitch(name), components.SwitchOptions{HideLabel: true})
	sw.Disabled(!s.Available())
	if on != v.lit {
		a.toggle(s)
	}
	return sw
}

// open shows the page of the device of tile t: an air conditioner's
// controls, else what its page shows of any other.
func (a *application) open(t ha.Tile) {
	if t.Entity.Domain() == "climate" {
		a.openClimate(ha.ParseClimate(a.held(t.Entity)))
		return
	}
	a.openDevice(t.Entity.EntityID)
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
