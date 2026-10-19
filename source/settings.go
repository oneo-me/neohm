package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/ha-go/source/locales"
	"github.com/oneo-me/onevra/source/components"
)

func (a *application) settingsWindowView(c *ui.Context) {
	a.scope(c)
	w, _ := c.Size()
	a.settingsView(c, w)
}

func (a *application) settingsView(c *ui.Context, width float32) {
	t := c.Theme()
	items := []components.SettingsItem{
		{ID: "address", Title: locales.Address(), Description: locales.AddressHint(), Stacked: true, Build: func(c *ui.Context) {
			components.TextBox(c, &a.address, components.TextBoxOptions{}).FillWidth().Label(locales.Address()).
				Placeholder("http://homeassistant.local:8123").OnSubmit(a.saveConnection)
		}},
		{ID: "token", Title: locales.Token(), Description: locales.TokenHint(), Stacked: true, Build: func(c *ui.Context) {
			components.TextBox(c, &a.token, components.TextBoxOptions{}).FillWidth().Label(locales.Token()).
				Password().OnSubmit(a.saveConnection)
		}},
		{ID: "connect", Title: locales.Connect(), Content: func(c *ui.Context) {
			ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(3)).Children(func() {
				b := components.Button(c, locales.Connect(), components.ButtonPrimary).Disabled(a.connecting || a.address == "" || a.token == "")
				if b.Clicked() {
					a.saveConnection()
				}
				switch {
				case a.connecting:
					components.Spinner(c, components.SpinnerRing).Size(16, 16)
					components.Text(c, locales.Connecting()).TextColor(t.TextMuted)
				case a.connectMsg != "":
					components.Text(c, a.connectMsg).Wrap().Grow(1).Basis(0).TextColor(t.TextMuted)
				}
			})
		}},
		{ID: "entity", Title: locales.Entity(), Description: locales.EntityHint(), Build: func(c *ui.Context) {
			a.entitySelect(c)
		}},
	}
	ui.Column(c).Fill().Children(func() {
		components.Settings(c, &a.settingsState, []components.SettingsCategory{{
			ID: "connection", Title: locales.Connection(),
			Groups: []components.SettingsGroup{{ID: "connection", Items: items}},
		}}, components.SettingsOptions{Narrow: width < 620, HideSearch: true, HideSidebar: true})
		if e := a.errorDisplay(); e != "" {
			ui.Column(c).FillWidth().Padding(0, t.Space(6), t.Space(4), t.Space(6)).Children(func() {
				components.Text(c, e).Wrap().FontSize(12).TextColor(t.Danger)
			})
		}
	})
}

func (a *application) entitySelect(c *ui.Context) {
	selected := a.config.Entity
	name := func(id string) string {
		for _, e := range a.entities {
			if e.EntityID == id {
				return e.Name
			}
		}
		if id == "" {
			return locales.EntityNone()
		}
		return id
	}
	s := components.Select(c, &selected)
	s.Trigger.Label(locales.Entity()).MinWidth(200).Disabled(len(a.entities) == 0).
		Children(func() { components.Text(c, name(selected)).SingleLine() })
	s.Popup(func(ui.Element) {
		for _, e := range a.entities {
			label := e.Name + "  ·  " + e.EntityID
			s.Item(e.EntityID).Label(label).Children(func() { components.Text(c, label) })
		}
	})
	if selected != a.config.Entity {
		a.selectEntity(selected)
	}
}
