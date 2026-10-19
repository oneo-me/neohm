package main

import (
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/neohm/source/locales"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// settingsPage edits the connection to Home Assistant and the language in
// the panel, over the page it covers. radius rounds the corners it scrolls
// within.
func (a *application) settingsPage(c *ui.Context, radius float32) {
	t := c.Theme()
	pageFrame(c, radius, locales.SettingsTitle(), func() {
		pageHeader(c, locales.SettingsTitle(), func() {
			if circleButton(c, lucide.X, locales.CloseSettings(), 36, false).Clicked() {
				a.closeSettings()
			}
		})
	}, func(top float32) {
		scrollContent(c, top, radius, headerEdge, func() {
			settingsGroup(c, locales.Language(), func() { a.languageSelect(c) })
			settingsGroup(c, locales.Connection(), func() {
				settingStacked(c, locales.Address(), "", func() {
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
					components.Text(c, e).Wrap().FontSize(12).TextColor(t.Danger).Margin(0, headerText-headerEdge)
				}
			})
		})
	})
}

// languageSelect chooses the language the app shows: the system's, or one
// of the catalogs', each named in itself.
func (a *application) languageSelect(c *ui.Context) {
	options := []components.NativeSelectOption[string]{{Value: "", Label: locales.LanguageSystem()}}
	for _, tag := range locales.Catalog.Languages() {
		options = append(options, components.NativeSelectOption[string]{Value: tag, Label: i18n.Autonym(tag)})
	}
	tag := a.config.Language
	s := components.NativeSelect(c, &tag, options)
	s.Trigger.FillWidth().Label(locales.Language())
	if s.Changed() {
		a.setLanguage(tag)
	}
}

// settingsGroup shows related settings under title, its text in line with
// the header's title and its controls with the header's buttons.
func settingsGroup(c *ui.Context, title string, build func()) {
	ui.Column(c).FillWidth().Gap(12).Padding(0, 0, 10).Children(func() {
		areaHeading(c, title).Margin(0, headerText-headerEdge)
		ui.Column(c).FillWidth().Gap(18).Children(build)
	})
}

// settingStacked puts a control below its title and description, if any.
func settingStacked(c *ui.Context, title, description string, build func()) {
	ui.Column(c).FillWidth().Gap(8).Children(func() {
		ui.Column(c).FillWidth().Gap(2).Padding(0, headerText-headerEdge).Children(func() {
			components.Text(c, title).FontWeight(600)
			if description != "" {
				components.Text(c, description).FontSize(12).TextColor(c.Theme().TextMuted).Wrap()
			}
		})
		build()
	})
}
