package main

import (
	"log"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/neohm/source/locales"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// settingsPage edits the language, opening at login, the connection to Home
// Assistant and updates in the panel, over the page it covers. radius rounds the corners
// it scrolls within.
func (a *application) settingsPage(c *ui.Context, radius float32) {
	pageFrame(c, radius, locales.SettingsTitle(), func() {
		pageHeader(c, locales.SettingsTitle(), func() {
			if circleButton(c, lucide.X, locales.CloseSettings(), 36, false).Clicked() {
				a.closeSettings()
			}
		})
	}, func(top float32) {
		scrollContent(c, nil, top, radius, 18, func() {
			settingsGroup(c, locales.General(), "", func() {
				a.languageRow(c)
				rowDivider(c)
				a.openAtLoginRow(c)
			})
			a.connectionGroup(c)
			a.updatesGroup(c)
			settingsGroup(c, "", "", func() {
				if actionRow(c, lucide.LogOut, danger, locales.Quit(), "", inkMuted, false, false) {
					a.quit()
				}
			})
		})
	})
}

// connectionGroup edits the address of Home Assistant and the access token,
// and connects to it, saying how that went.
func (a *application) connectionGroup(c *ui.Context) {
	settingsGroup(c, locales.Connection(), locales.TokenHint(), func() {
		fieldRow(c, lucide.Server, cyan, locales.Address(), &a.address, "http://homeassistant.local:8123", false, a.saveConnection)
		rowDivider(c)
		fieldRow(c, lucide.KeyRound, amber, locales.Token(), &a.token, "", true, a.saveConnection)
		rowDivider(c)
		detail, color := "", inkMuted
		switch e := a.errorDisplay(); {
		case a.connecting:
			detail = locales.Connecting()
		case e != "":
			detail, color = e, danger
		case a.connected:
			detail, color = locales.ConnectedTo(a.tileCount()), green
		}
		if actionRow(c, lucide.Link, accent, locales.Connect(), detail, color, a.connecting, a.address == "" || a.token == "") {
			a.saveConnection()
		}
	})
}

// languageRow chooses the language the app shows, from the menu of the
// language shown at its end: the system's, or one of the catalogs', each
// named in itself.
func (a *application) languageRow(c *ui.Context) {
	current := locales.LanguageSystem()
	if a.config.Language != "" {
		current = i18n.Autonym(a.config.Language)
	}
	staticRow(c).Children(func() {
		rowBody(c, lucide.Languages, lavender, locales.Language(), "", inkMuted, func() { a.languageMenu(c, current) })
	})
}

// languageMenu is the language current, which opens the menu of languages
// below it: a capsule that lightens while pressed, its text where the row's
// other controls end.
func (a *application) languageMenu(c *ui.Context, current string) {
	b := ui.ButtonBase(c).Label(locales.Language() + ": " + current).Row().AlignItems(ui.Center).Gap(4).
		Height(30).PaddingX(10).MarginEnd(-10).Radius(15).Shrink(1).MinWidth(0)
	if b.Pressed() {
		b.Background(ui.RGBA(255, 255, 255, 0.1))
	}
	b.Children(func() {
		components.Text(c, current).FontSize(14).TextColor(inkMuted).SingleLine().Ellipsis("…").MinWidth(0)
		components.Icon(c, lucide.ChevronsUpDown).Size(14, 14).TextColor(inkMuted).Shrink(0)
	})
	b.Menu(func(m *ui.Menu) {
		if m.Item(locales.LanguageSystem()).Checked(a.config.Language == "").Chosen() {
			a.setLanguage("")
		}
		for _, tag := range locales.Catalog.Languages() {
			if m.Item(i18n.Autonym(tag)).Checked(a.config.Language == tag).Chosen() {
				a.setLanguage(tag)
			}
		}
	})
}

// openAtLoginRow turns opening the app when the user logs in on or off,
// which the system keeps rather than the configuration file.
func (a *application) openAtLoginRow(c *ui.Context) {
	on := a.atLogin
	detail, color := locales.OpenAtLoginHint(), inkMuted
	if a.loginFailed {
		detail, color = locales.OpenAtLoginError(), danger
	}
	switchRow(c, lucide.Power, green, locales.OpenAtLogin(), detail, color, &on)
	if on != a.atLogin {
		a.changeOpenAtLogin(on)
	}
}

// changeOpenAtLogin asks the system to open the app at login, or not, and
// shows what it then reports.
func (a *application) changeOpenAtLogin(on bool) {
	err := a.setOpenAtLogin(on)
	if err != nil {
		log.Printf("openAtLoginError: %v", err)
	}
	a.atLogin, a.loginFailed = a.openAtLogin(), err != nil
	a.redraw()
}

// Settings' own controls: groups of rows, each group one pane of glass and
// its rows divided by hairlines rather than set on surfaces of their own.

const (
	// groupRadius rounds a group's pane.
	groupRadius = 22
	// rowInset is where a row's text starts, past its badge, and where the
	// hairline above it starts.
	rowInset = 14 + badgeSize + 12
	// badgeSize is the width of the disc at the start of a row.
	badgeSize = 30
)

// settingsGroup shows the rows rows builds on a pane of glass under title,
// and footer below it, either left out when empty.
func settingsGroup(c *ui.Context, title, footer string, rows func()) {
	ui.Column(c).FillWidth().Gap(8).Children(func() {
		if title != "" {
			areaHeading(c, title)
		}
		ui.Column(c).FillWidth().Radius(groupRadius).Clip().
			Material(glass.Glass{Tint: ui.RGBA(255, 255, 255, 0.04)}).Children(rows)
		if footer != "" {
			components.Text(c, footer).FontSize(12).TextColor(inkMuted).Wrap().Margin(0, 14)
		}
	})
}

// rowDivider is the hairline between two rows of a group, from below the
// text of the row above to the group's edge.
func rowDivider(c *ui.Context) {
	ui.Box(c).FillWidth().Height(1).MarginStart(rowInset).Background(ui.RGBA(255, 255, 255, 0.08))
}

// rowBody lays a row out: icon on a disc of color, title and, unless empty,
// detail below it in detailColor, and what trailing builds at the end. A
// nil icon spins in place of one.
func rowBody(c *ui.Context, icon func() *ui.SVG, color ui.Color, title, detail string, detailColor ui.Color, trailing func()) {
	badge(c, icon, color)
	ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(2).Children(func() {
		components.Text(c, title).FontSize(14).FontWeight(600).TextColor(ink).Wrap()
		if detail != "" {
			components.Text(c, detail).FontSize(12).TextColor(detailColor).Wrap()
		}
	})
	if trailing != nil {
		trailing()
	}
}

// badge is the disc of color at the start of a row, with icon on it, dark
// while the disc is light. A translucent color makes a quiet disc.
func badge(c *ui.Context, icon func() *ui.SVG, color ui.Color) {
	glyph := ink
	if color.A == 255 && luminance(color) > 0.4 {
		glyph = discInk
	}
	ui.Box(c).Size(badgeSize, badgeSize).Radius(badgeSize / 2).Shrink(0).Background(color).Center().Children(func() {
		if icon == nil {
			components.Spinner(c, components.SpinnerRing).Size(16, 16)
			return
		}
		components.Icon(c, icon).Size(16, 16).TextColor(glyph)
	})
}

// quietBadge is the color of a badge that only tells where something got.
var quietBadge = ui.RGBA(255, 255, 255, 0.12)

// staticRow is a row that only shows something.
func staticRow(c *ui.Context) ui.Element {
	return ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(12).Padding(11, 14).MinHeight(54)
}

// actionRow is a row that does what title says, in color, which detail
// below it says more of in detailColor, and reports a click on it. While
// busy it spins at its end and ignores clicks, undimmed.
func actionRow(c *ui.Context, icon func() *ui.SVG, color ui.Color, title, detail string, detailColor ui.Color, busy, disabled bool) bool {
	return pressRow(c, title, disabled && !busy, func() {
		badge(c, icon, color)
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(2).Children(func() {
			components.Text(c, title).FontSize(14).FontWeight(600).TextColor(color.Mix(ink, 0.25)).Wrap()
			if detail != "" {
				components.Text(c, detail).FontSize(12).TextColor(detailColor).Wrap()
			}
		})
		if busy {
			components.Spinner(c, components.SpinnerRing).Size(16, 16)
		}
	}) && !busy
}

// pressRow is a row labeled label, of what content builds, that reports a
// click on it. It lightens while pressed and while it has the keyboard
// focus, in place of a ring the group's edge would clip, and dims while
// disabled.
func pressRow(c *ui.Context, label string, disabled bool, content func()) bool {
	row := ui.ButtonBase(c).Label(label).Disabled(disabled).
		FillWidth().Row().AlignItems(ui.Center).Gap(12).Padding(11, 14).MinHeight(54).FocusRing(false)
	pressed, focus := row.Pressed() && !disabled, row.FocusVisible()
	row.Draw(func(p *ui.Painter, r ui.Rect) {
		if pressed {
			p.Fill(r, ui.RGBA(255, 255, 255, 0.07), 0)
		}
		if focus {
			p.Fill(r, accent.Alpha(0.16), 0)
		}
	})
	if disabled {
		row.Opacity(0.45)
	}
	row.Children(content)
	return row.Clicked()
}

// fieldRow edits *value, named title above it, on a row; submit runs on
// Enter. The title lights while the field has the keyboard.
func fieldRow(c *ui.Context, icon func() *ui.SVG, color ui.Color, title string, value *string, placeholder string, password bool, submit func()) {
	row := staticRow(c)
	name := inkMuted
	if row.FocusWithin() {
		name = accent.Mix(ink, 0.25)
	}
	row.Children(func() {
		badge(c, icon, color)
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(3).Children(func() {
			components.Text(c, title).FontSize(12).FontWeight(600).TextColor(name).SingleLine().Ellipsis("…")
			in := ui.TextInputBase(c, value).FillWidth().FontSize(14).TextColor(ink).Label(title).
				Placeholder(placeholder).SelectionColor(accent.Alpha(0.4)).OnSubmit(submit)
			if password {
				in.Password()
			}
		})
	})
}

// switchRow turns *on on or off with its switch, at its end.
func switchRow(c *ui.Context, icon func() *ui.SVG, color ui.Color, title, detail string, detailColor ui.Color, on *bool) {
	staticRow(c).Children(func() {
		rowBody(c, icon, color, title, detail, detailColor, func() { glassSwitch(c, on, title, detail) })
	})
}

// glassSwitch turns *on on or off, named label and described by
// description: a track filled in accent while on, and a white thumb that,
// while held, stretches into clear glass the track shows through, as macOS
// 26 draws it.
func glassSwitch(c *ui.Context, on *bool, label, description string) ui.Element {
	const width, height, inset = 48, 28, 2
	e := ui.SwitchBase(c, on).Size(width, height).Radius(height / 2).Shrink(0).Label(label).Description(description)
	lit := e.Animate("on", map[bool]float32{false: 0, true: 1}[*on], 220*time.Millisecond)
	grow := e.Animate("held", map[bool]float32{false: 0, true: 1}[e.Pressed()], 150*time.Millisecond)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, ui.RGBA(255, 255, 255, 0.16).Mix(accent, lit), r.H/2)
		size := r.H - 2*inset
		w := size + 8*grow
		x := r.X + inset + (r.W-2*inset-w)*lit
		k := ui.Rect{X: x, Y: r.Y + inset, W: w, H: size}
		p.Shadow(k, size/2, 0, 2, 6, 0, ui.RGBA(0, 0, 0, 0.25*(1-grow)))
		if grow < 1 {
			p.Fill(k, ui.RGBA(255, 255, 255, 1-grow), size/2)
		}
		if grow > 0 {
			glass.Paint(p, k, size/2, glass.Glass{Style: glass.Clear, Tint: ui.RGBA(255, 255, 255, 0.16*grow)})
		}
	})
	return e
}

// progressBar shows how far, 0 to 1, a download got.
func progressBar(c *ui.Context, done float64) {
	ui.Box(c).FillWidth().Height(6).Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, ui.RGBA(255, 255, 255, 0.12), r.H/2)
		if done > 0 {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.H, r.W*float32(min(done, 1))), H: r.H}, accent, r.H/2)
		}
	})
}
