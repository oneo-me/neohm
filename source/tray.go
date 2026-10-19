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
	"github.com/oneo-me/neohm/source/locales"
	"github.com/oneo-me/onevra/source/shells/plain"
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
	if a.tray == nil {
		return
	}
	// The panel holds everything else, Settings included.
	quit := item(locales.Quit(), a.quit)
	quit.Accelerator = "CmdOrCtrl+Q"
	items := []*mygo.MenuItem{quit}
	if runtime.GOOS == "linux" {
		// AppIndicator only reports menu actions, not button events, so the
		// menu opens the panel too.
		items = append([]*mygo.MenuItem{item(locales.OpenPanel(), a.showPanel), mygo.Separator()}, items...)
	}
	a.trayMenu = mygo.NewMenu(items)
	if runtime.GOOS == "linux" {
		a.tray.SetMenu(a.trayMenu)
	}
}

// relabel says again, in the language the app shows now, what was said only
// once: the menus, the menu bar icon's tooltip and the window's title. The
// panel follows by itself, as it is built anew every frame.
func (a *application) relabel() {
	a.installMenus()
	if a.tray != nil {
		a.tray.SetToolTip(locales.Title())
	}
	if a.window != nil {
		a.window.Window().SetTitle(locales.Title())
	}
}

func (a *application) reopen() {
	if !a.quitting {
		a.showPanel()
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

// showPanel opens the panel under the menu bar icon. It draws its own dark
// surface and arrow rather than the system popover's material, which the
// blur under its pages' headers cannot see. Without a menu bar icon, its
// pages show in a window.
func (a *application) showPanel() {
	if a.tray == nil {
		a.showWindow()
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
	a.follow()
	panel.OnClosed(func() {
		if a.panel == panel {
			a.panel = nil
			a.panelOpen.Store(false)
			a.resetPanel()
		}
	})
}

// showWindow shows the panel's pages in a window when there is no menu bar
// icon to anchor the panel to. The window is then the app: closing it quits.
func (a *application) showWindow() {
	if a.window == nil {
		a.window = plain.NewShell(plain.ShellOptions{
			Title: locales.Title(), Theme: windowTheme,
			Width: panelWidth, Height: panelHeight, MinWidth: 360, MinHeight: 480,
			Hidden: true, Content: a.windowView,
		})
		a.window.Window().OnClose(func(*mygo.CloseEvent) {
			if !a.quitting {
				mygo.App.Quit()
			}
		})
	}
	w := a.window.Window()
	if w.IsMinimized() {
		w.Restore()
	}
	w.Show()
	w.Focus()
	a.panelOpen.Store(true)
	a.follow()
}

// resetPanel returns the closed panel to the page it opens on: the home
// page, or Settings until Home Assistant is set up.
func (a *application) resetPanel() {
	a.detail, a.device, a.settingsOpen = false, "", !a.configured()
	a.upd.open = false
	a.follow()
}

// openSettings shows Settings in the panel, opening the panel first when it is
// closed.
func (a *application) openSettings() {
	a.showSettings()
	a.showPanel()
}

// showSettings shows Settings over the page the panel shows.
func (a *application) showSettings() {
	if !a.connecting {
		a.connected = false
	}
	a.settingsOpen = true
	a.redraw()
}

func (a *application) closeSettings() {
	a.settingsOpen = false
	a.redraw()
}

// back leaves Settings, else the device opened, reporting whether the panel
// showed either.
func (a *application) back() bool {
	switch {
	case a.settingsOpen:
		a.closeSettings()
	case a.detail || a.device != "":
		a.closeDevice()
	default:
		return false
	}
	return true
}

func (a *application) destroyTray() {
	if a.tray != nil {
		a.tray.Destroy()
		a.tray = nil
	}
}

// trayIcon draws a house holding four blocks, the tiles of the home page, as
// a template image: black on transparent, 18 points at 2x. Its strokes are 3
// pixels wide, and its walls, floor and blocks fall on whole pixels.
func trayIcon() []byte {
	const size = 36
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	ink := color.NRGBA{A: 255}
	if runtime.GOOS != "darwin" {
		ink = color.NRGBA{R: 77, G: 141, B: 255, A: 255}
	}
	// The outline follows a pentagon rounded by r: its walls at x 5.5 and
	// 30.5, its floor at 30.5 and its eaves at y 14.5.
	const r = 2.5
	house := [][2]float64{{18, 4.5 + 1.25*r}, {30.5 - r, 14.5}, {30.5 - r, 30.5 - r}, {5.5 + r, 30.5 - r}, {5.5 + r, 14.5}}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			px, py := float64(x)+0.5, float64(y)+0.5
			d := math.Abs(polygonDist(px, py, house)-r) - 1.5
			for _, b := range [][2]float64{{14.5, 16.5}, {21.5, 16.5}, {14.5, 23.5}, {21.5, 23.5}} {
				d = min(d, roundRectDist(px, py, b[0], b[1], 2.5, 2.5, 1.2))
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

// polygonDist is the signed distance from (px, py) to the polygon of
// vertices v, negative inside.
func polygonDist(px, py float64, v [][2]float64) float64 {
	d, inside := math.Inf(1), false
	for i, j := 0, len(v)-1; i < len(v); j, i = i, i+1 {
		ax, ay, bx, by := v[i][0], v[i][1], v[j][0], v[j][1]
		ex, ey, wx, wy := bx-ax, by-ay, px-ax, py-ay
		t := max(0, min(1, (wx*ex+wy*ey)/(ex*ex+ey*ey)))
		d = min(d, math.Hypot(wx-ex*t, wy-ey*t))
		if (ay > py) != (by > py) && px < ax+(py-ay)*ex/ey {
			inside = !inside
		}
	}
	if inside {
		return -d
	}
	return d
}
