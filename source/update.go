package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
	"github.com/oneo-me/neohm/source/locales"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// updatePhase is where checking for, or installing, a newer version got.
type updatePhase int

const (
	updateIdle updatePhase = iota
	updateChecking
	updateCurrent
	updateAvailable
	updateInstalling
	updateFailed
	// updateDisabled is a copy that cannot update itself: a development
	// build, or an app where it cannot write.
	updateDisabled
)

// updates checks for and installs newer versions, published as signed
// releases on GitHub; tests replace it.
type updates struct {
	check    func(context.Context) (*mygo.Update, error)
	install  func(context.Context, *mygo.Update, func(downloaded, total int64)) error
	relaunch func()
}

var mygoUpdates = updates{
	check: mygo.Updater.Check,
	install: func(ctx context.Context, up *mygo.Update, progress func(downloaded, total int64)) error {
		return up.Install(ctx, progress)
	},
	relaunch: mygo.App.Relaunch,
}

// updateState is what the update button's popover shows: the phase, the
// newer version found and how far its download got.
type updateState struct {
	open     bool
	phase    updatePhase
	found    *mygo.Update
	progress float64
	// installFailed tells a failed install, which Try Again installs again,
	// from a failed check.
	installFailed bool
	notes         *markdown.Options
}

// updateButton, left of Settings in Home's header, checks for a newer
// version and says what it found in a popover below it.
func (a *application) updateButton(c *ui.Context) {
	u := &a.upd
	b := circleButton(c, lucide.CircleArrowUp, locales.CheckForUpdates(), 36, false)
	if b.Clicked() {
		u.open = !u.open
		if u.open {
			switch u.phase {
			case updateIdle, updateCurrent, updateFailed, updateDisabled:
				a.checkForUpdates()
			}
		}
	}
	b.Expanded(u.open)
	components.Popover(c, b, &u.open, func(panel ui.Element) {
		panel.Width(268).Label(locales.CheckForUpdates())
		a.updateView(c)
	})
}

func (a *application) updateView(c *ui.Context) {
	u := &a.upd
	current := mygo.App.Version()
	switch u.phase {
	case updateChecking:
		ui.Row(c).FillWidth().Gap(10).AlignItems(ui.Center).Children(func() {
			components.Spinner(c, components.SpinnerSpokes).Size(16, 16)
			components.Text(c, locales.UpdateChecking()).FontSize(13).TextColor(ink)
		})
	case updateCurrent:
		updateHeading(c, lucide.CircleCheck, green, locales.UpdateCurrent(), locales.UpdateVersion(current))
	case updateDisabled:
		updateHeading(c, lucide.CircleAlert, inkMuted, locales.UpdateDisabled(), locales.UpdateDisabledBody())
	case updateFailed:
		title := locales.UpdateCheckError()
		if u.installFailed {
			title = locales.UpdateInstallError()
		}
		updateHeading(c, lucide.TriangleAlert, danger, title, locales.UpdateVersion(current))
		if components.Button(c, locales.UpdateRetry(), components.ButtonSecondary).FillWidth().Clicked() {
			if u.installFailed {
				a.installUpdate()
			} else {
				a.checkForUpdates()
			}
		}
	case updateAvailable:
		updateHeading(c, lucide.CircleArrowUp, accent, locales.UpdateAvailable(u.found.Version), locales.UpdateYours(current))
		if u.found.Notes != "" {
			components.ScrollView(c, components.ScrollViewOptions{}).FillWidth().MaxHeight(180).Children(func() {
				ui.Column(c).FillWidth().Padding(0, 2).Children(func() {
					markdown.View(c, u.found.Notes, u.notes)
				})
			})
		}
		if components.Button(c, locales.UpdateInstall(), components.ButtonPrimary).FillWidth().Clicked() {
			a.installUpdate()
		}
	case updateInstalling:
		updateHeading(c, lucide.CircleArrowUp, accent, locales.UpdateInstalling(u.found.Version), locales.UpdateRelaunch())
		components.Progress(c, u.progress)
	}
}

// updateHeading is the popover's title, after an icon in color, above a
// line that says more.
func updateHeading(c *ui.Context, icon func() *ui.SVG, color ui.Color, title, detail string) {
	ui.Row(c).FillWidth().Gap(10).AlignItems(ui.Start).Children(func() {
		components.Icon(c, icon).Size(18, 18).TextColor(color).Shrink(0)
		ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
			components.Text(c, title).FontSize(13).Bold().Wrap().TextColor(ink)
			components.Text(c, detail).FontSize(12).Wrap().TextColor(inkMuted)
		})
	})
}

// checkForUpdates asks GitHub for a version newer than this one.
func (a *application) checkForUpdates() {
	u := &a.upd
	u.phase, u.found, u.installFailed = updateChecking, nil, false
	a.redraw()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		found, err := a.updates.check(ctx)
		a.post(func() {
			switch {
			case errors.Is(err, mygo.ErrUpdatesDisabled):
				u.phase = updateDisabled
			case err != nil:
				log.Printf("updateCheckError: %v", err)
				u.phase = updateFailed
			case found == nil:
				u.phase = updateCurrent
			default:
				u.phase, u.found = updateAvailable, found
			}
			a.redraw()
		})
	}()
}

// installUpdate downloads the version found, replaces the app with it and
// relaunches into it.
func (a *application) installUpdate() {
	u := &a.upd
	if u.found == nil {
		return
	}
	found := u.found
	u.phase, u.progress, u.installFailed = updateInstalling, 0, false
	a.redraw()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		shown := -1
		err := a.updates.install(ctx, found, func(downloaded, total int64) {
			if total <= 0 {
				return
			}
			// Redraw once a percent, not once a chunk.
			if p := int(downloaded * 100 / total); p != shown {
				shown = p
				a.post(func() {
					if u.phase == updateInstalling {
						u.progress = float64(p) / 100
						a.redraw()
					}
				})
			}
		})
		a.post(func() {
			if err != nil {
				log.Printf("updateInstallError: %v", err)
				u.phase, u.installFailed = updateFailed, true
				a.redraw()
				return
			}
			u.progress = 1
			a.updates.relaunch()
		})
	}()
}
