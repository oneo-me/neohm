package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/neohm/source/locales"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// settingsPage edits the connection to Home Assistant in the panel, over the
// page it covers. radius rounds the bottom corners it scrolls within.
func (a *application) settingsPage(c *ui.Context, radius float32) {
	t := c.Theme()
	pageHeader(c, locales.SettingsTitle(), func() {
		if circleButton(c, lucide.X, locales.CloseSettings(), 36, false).Clicked() {
			a.closeSettings()
		}
	})
	components.ScrollView(c, components.ScrollViewOptions{}).Fill().Radius(0, 0, radius, radius).
		ScrollbarInsets(0, 4, radius).Children(func() {
		ui.Column(c).FillWidth().Padding(4, 22, 22, 22).Gap(18).Children(func() {
			components.Text(c, locales.Connection()).FontSize(12).Bold().TextColor(t.TextMuted)
			settingStacked(c, locales.Address(), locales.AddressHint(), func() {
				components.TextBox(c, &a.address, components.TextBoxOptions{}).FillWidth().Label(locales.Address()).
					Placeholder("http://homeassistant.local:8123").OnSubmit(a.saveConnection)
			})
			settingStacked(c, locales.Token(), locales.TokenHint(), func() {
				components.TextBox(c, &a.token, components.TextBoxOptions{}).FillWidth().Label(locales.Token()).
					Password().OnSubmit(a.saveConnection)
			})
			ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(3)).Children(func() {
				b := components.Button(c, locales.Connect(), components.ButtonPrimary).Disabled(a.connecting || a.address == "" || a.token == "")
				if b.Clicked() {
					a.saveConnection()
				}
				switch {
				case a.connecting:
					components.Spinner(c, components.SpinnerRing).Size(16, 16)
					components.Text(c, locales.Connecting()).TextColor(t.TextMuted)
				case a.connected:
					components.Text(c, locales.ConnectedTo(a.tileCount())).Wrap().Grow(1).Basis(0).TextColor(t.TextMuted)
				}
			})
			if e := a.errorDisplay(); e != "" {
				components.Text(c, e).Wrap().FontSize(12).TextColor(t.Danger)
			}
		})
	})
}

// settingStacked puts a control below its title and description.
func settingStacked(c *ui.Context, title, description string, build func()) {
	ui.Column(c).FillWidth().Gap(8).Children(func() {
		ui.Column(c).FillWidth().Gap(2).Children(func() {
			components.Text(c, title).FontWeight(600)
			components.Text(c, description).FontSize(12).TextColor(c.Theme().TextMuted).Wrap()
		})
		build()
	})
}
