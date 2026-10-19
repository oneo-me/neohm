package main

import (
	"math"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/ha-go/source/internal/ha"
	"github.com/oneo-me/ha-go/source/locales"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// listPage is where the panel opens: the air conditioners, each opening its
// controls. radius rounds the bottom corners the list scrolls within.
func (a *application) listPage(c *ui.Context, radius float32) {
	pageHeader(c, locales.Devices(), func() {
		if circleButton(c, lucide.Settings, locales.SettingsTitle(), 36, false).Clicked() {
			a.showSettings()
		}
	})
	if !a.configured() || !a.listed {
		a.setupView(c)
		return
	}
	// A failure after the list was read leaves the list as last read.
	if e := a.errorDisplay(); e != "" {
		ui.Row(c).FillWidth().Shrink(0).Gap(8).Padding(0, 22, 10, 22).AlignItems(ui.Center).Children(func() {
			components.Icon(c, lucide.TriangleAlert).Size(16, 16).TextColor(danger).Shrink(0)
			components.Text(c, e).FontSize(12).Wrap().TextColor(inkMuted).Grow(1).MinWidth(0)
		})
	}
	if len(a.devices) == 0 {
		ui.Column(c).Fill().Center().Gap(10).Padding(0, 24, 48, 24).Children(func() {
			components.Icon(c, lucide.AirVent).Size(44, 44).TextColor(inkMuted)
			components.Text(c, locales.NoDevices()).FontSize(17).Bold().TextColor(ink)
			components.Text(c, locales.NoDevicesBody()).Wrap().TextAlign(ui.Center).FontSize(13).TextColor(inkMuted).MaxWidth(260)
		})
		return
	}
	components.ScrollView(c, components.ScrollViewOptions{}).Fill().Radius(0, 0, radius, radius).
		ScrollbarInsets(0, 4, radius).Children(func() {
		ui.Column(c).FillWidth().Padding(2, 18, 18, 18).Gap(10).Children(func() {
			for _, d := range a.devices {
				a.deviceRow(c, a.shown(d))
			}
		})
	})
}

// setupView stands in for the list until Home Assistant lists the devices:
// it waits, asks for the connection or says what went wrong.
func (a *application) setupView(c *ui.Context) {
	ui.Column(c).Fill().Center().Gap(14).Padding(0, 24, 48, 24).Children(func() {
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
		components.Button(c, locales.Settings(), components.ButtonPrimary).OnClick(a.showSettings)
	})
}

// deviceRow shows an air conditioner's mode, name and temperatures, and opens
// its controls.
func (a *application) deviceRow(c *ui.Context, d ha.Climate) {
	on := d.On()
	detail := a.deviceDetail(d)
	r := ui.ButtonBase(c.Key(d.EntityID)).FillWidth().Height(72).Radius(20).Padding(0, 14, 0, 14).Gap(12).
		Justify(ui.Start).Label(d.Name + ", " + detail)
	face := tileFace
	switch {
	case r.Pressed():
		face = ui.RGBA(255, 255, 255, 0.12)
	case r.Hovered():
		face = tileHover
	}
	r.Background(face).Border(1, tileEdge)
	r.Children(func() {
		icon, tint, disc := lucide.Power, inkMuted, ui.RGBA(255, 255, 255, 0.06)
		if on {
			color := modeColor(d.Mode)
			icon, tint, disc = modeIcon(d.Mode), color, color.Alpha(0.16)
		}
		ui.Box(c).Size(40, 40).Radius(20).Shrink(0).Background(disc).Center().Children(func() {
			components.Icon(c, icon).Size(20, 20).TextColor(tint)
		})
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(3).Children(func() {
			components.Text(c, d.Name).FontSize(15).FontWeight(600).TextColor(ink).SingleLine().Ellipsis("…")
			components.Text(c, detail).FontSize(12).TextColor(inkMuted).SingleLine().Ellipsis("…")
		})
		if on && d.HasTarget {
			components.Text(c, formatTemp(d.Target)+"°").FontSize(24).FontWeight(300).TextColor(ink).Shrink(0)
		}
		components.Icon(c, lucide.ChevronRight).Size(16, 16).TextColor(inkMuted).Shrink(0)
	})
	if r.Clicked() {
		a.openDevice(d)
	}
}

// deviceDetail is a row's second line: the mode and the room temperature.
func (a *application) deviceDetail(d ha.Climate) string {
	if !d.Available() {
		return locales.Unavailable()
	}
	text := optionText("mode", d.Mode)
	if d.HasCurrent {
		text += " · " + locales.Room() + " " + formatTemp(math.Round(d.Current*10)/10) + a.unit
	}
	return text
}
