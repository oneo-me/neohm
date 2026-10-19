package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"runtime"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/ha-go/source/locales"
	"github.com/oneo-me/onevra/source/shells/plain"
	"github.com/oneo-me/onevra/source/theme"
)

// The panel's size in DIPs, its arrow included.
const (
	panelWidth  = 388
	panelHeight = 572
)

func (a *application) createTray() error {
	var err error
	a.tray, err = mygo.NewTray(mygo.TrayOptions{Icon: trayIcon(), IconIsTemplate: true, ToolTip: locales.Title()})
	if err != nil {
		return err
	}
	if runtime.GOOS != "linux" {
		a.tray.OnClick(a.togglePanel)
		a.tray.OnRightClick(func() {
			a.hidePanel()
			a.tray.SetMenu(a.trayMenu)
			defer a.tray.SetMenu(nil)
			a.tray.PopUpMenu()
		})
	}
	return nil
}

func (a *application) installMenus() {
	item := func(label string, fn func()) *mygo.MenuItem {
		return &mygo.MenuItem{Label: label, Click: func(*mygo.MenuItem, *mygo.Window) { fn() }}
	}
	settings := item(locales.Settings(), a.openSettings)
	settings.Accelerator = "CmdOrCtrl+,"
	if runtime.GOOS == "darwin" {
		mygo.App.SetMenu(mygo.NewMenu([]*mygo.MenuItem{
			{Role: mygo.RoleAppMenu, Submenu: []*mygo.MenuItem{settings, mygo.Separator(), {Role: mygo.RoleQuit}}},
			{Role: mygo.RoleEditMenu},
		}))
	}
	if a.tray != nil {
		a.trayMenu = mygo.NewMenu([]*mygo.MenuItem{
			item(locales.OpenPanel(), a.showPanel),
			item(locales.Refresh(), a.wake),
			item(locales.Settings(), a.openSettings),
			mygo.Separator(),
			item(locales.Quit(), mygo.App.Quit),
		})
		// AppIndicator only reports menu actions, not button events.
		if runtime.GOOS == "linux" {
			a.tray.SetMenu(a.trayMenu)
		}
	}
}

func (a *application) reopen() {
	if a.quitting {
		return
	}
	if a.tray != nil && a.configured() {
		a.showPanel()
	} else {
		a.openSettings()
	}
}

func (a *application) togglePanel() {
	if a.panel != nil {
		a.hidePanel()
	} else {
		a.showPanel()
	}
}

func (a *application) hidePanel() {
	if a.panel != nil {
		a.panel.Close()
	}
}

// showPanel opens the climate panel under the menu bar icon. It draws its own
// dark surface and arrow rather than the system popover's material.
func (a *application) showPanel() {
	if a.tray == nil {
		a.openSettings()
		return
	}
	if a.panel != nil {
		return
	}
	panel := mygo.NewFlyout(mygo.FlyoutOptions{
		Tray: a.tray, Placement: mygo.PlacementBottom, Gap: 2,
		Width: panelWidth, Height: panelHeight, Focusable: true, Shadow: true,
		Content: ui.View(a.panelView),
	})
	a.panel = panel
	a.panelOpen.Store(true)
	a.wake()
	panel.OnClosed(func() {
		if a.panel == panel {
			a.panel = nil
			a.panelOpen.Store(false)
		}
	})
}

func (a *application) openSettings() {
	a.hidePanel()
	if a.settings == nil {
		shell := plain.NewShell(plain.ShellOptions{
			Title: locales.Settings(), Theme: theme.Onevra,
			Width: 560, Height: 520, MinWidth: 420, MinHeight: 420,
			Hidden: true, Content: a.settingsWindowView,
		})
		a.settings = shell.Window()
		a.settings.OnClose(func(e *mygo.CloseEvent) {
			if a.quitting {
				return
			}
			e.PreventDefault()
			if a.tray == nil {
				mygo.App.Quit()
			} else {
				a.settings.Hide()
			}
		})
	}
	if a.configured() && a.entities == nil && !a.connecting {
		a.connect(true)
	}
	if a.settings.IsMinimized() {
		a.settings.Restore()
	}
	a.settings.Show()
	a.settings.Focus()
}

func (a *application) destroyTray() {
	if a.tray != nil {
		a.tray.Destroy()
		a.tray = nil
	}
}

// trayIcon draws an air conditioner blowing air, as a template image: black
// on transparent, 18 points at 2x.
func trayIcon() []byte {
	const size = 36
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	ink := color.NRGBA{A: 255}
	if runtime.GOOS != "darwin" {
		ink = color.NRGBA{R: 77, G: 141, B: 255, A: 255}
	}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			px, py := float64(x)+0.5, float64(y)+0.5
			d := math.Abs(roundRectDist(px, py, 18, 12, 14, 7, 4)) - 1.3
			d = min(d, segmentDist(px, py, 9, 15, 27, 15)-0.9)
			for _, l := range [][4]float64{{12, 23, 10, 31}, {18, 23, 18, 31}, {24, 23, 26, 31}} {
				d = min(d, segmentDist(px, py, l[0], l[1], l[2], l[3])-1.2)
			}
			if cov := max(0, min(1, 0.5-d)); cov > 0 {
				c := ink
				c.A = uint8(float64(ink.A) * cov)
				img.SetNRGBA(x, y, c)
			}
		}
	}
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		panic(err)
	}
	return b.Bytes()
}

// roundRectDist is the signed distance from (px, py) to a rounded rectangle
// centered at (cx, cy) with half sizes hw and hh.
func roundRectDist(px, py, cx, cy, hw, hh, r float64) float64 {
	qx, qy := math.Abs(px-cx)-hw+r, math.Abs(py-cy)-hh+r
	return math.Hypot(max(qx, 0), max(qy, 0)) + min(max(qx, qy), 0) - r
}

func segmentDist(px, py, ax, ay, bx, by float64) float64 {
	dx, dy := bx-ax, by-ay
	t := max(0, min(1, ((px-ax)*dx+(py-ay)*dy)/(dx*dx+dy*dy)))
	return math.Hypot(px-ax-t*dx, py-ay-t*dy)
}
