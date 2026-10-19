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
	// notify tells the user a newer version is out, in a notification that
	// opens Settings.
	notify func(version string)
}

var mygoUpdates = updates{
	check: mygo.Updater.Check,
	install: func(ctx context.Context, up *mygo.Update, progress func(downloaded, total int64)) error {
		return up.Install(ctx, progress)
	},
	relaunch: mygo.App.Relaunch,
	notify: func(version string) {
		n := mygo.NewNotification(mygo.NotificationOptions{
			ID: updateNotification, Title: locales.UpdateAvailable(version), Body: locales.UpdateNotificationBody(),
		})
		// Show waits for the user to allow notifications the first time.
		go func() {
			if err := n.Show(); err != nil {
				log.Printf("updateNotifyError: %v", err)
			}
		}()
	},
}

// autoCheckEvery is how often the app checks for a newer version by
// itself, unless Settings turned that off.
const autoCheckEvery = 24 * time.Hour

// updateNotification identifies the notification of a newer version, which
// opens Settings when clicked.
const updateNotification = "update"

// updateState is what Settings' updates show: the phase, the newer version
// found and how far its download got.
type updateState struct {
	phase    updatePhase
	found    *mygo.Update
	progress float64
	// installFailed tells a failed install, which Try Again installs again,
	// from a failed check.
	installFailed bool
	notes         *markdown.Options
	// checked is when the last check started, on the wall clock, which
	// keeps running while the Mac sleeps; notified is the version the user
	// was last told of.
	checked  time.Time
	notified string
}

// updatesGroup, in Settings, says whether a newer version is out, checks for
// one and installs it, and turns checking by itself on or off.
func (a *application) updatesGroup(c *ui.Context) {
	u := &a.upd
	settingsGroup(c, locales.Updates(), func() {
		ui.Column(c).FillWidth().Gap(12).Children(func() { a.updateView(c) })
		switch {
		case u.phase == updateFailed && u.installFailed:
			if components.Button(c, locales.UpdateRetry(), components.ButtonSecondary).Clicked() {
				a.installUpdate()
			}
		case u.phase == updateFailed:
			if components.Button(c, locales.UpdateRetry(), components.ButtonSecondary).Clicked() {
				a.checkForUpdates(false)
			}
		case u.phase != updateAvailable && u.phase != updateInstalling && u.phase != updateDisabled:
			b := components.Button(c, locales.CheckForUpdates(), components.ButtonSecondary).Disabled(u.phase == updateChecking)
			if b.Clicked() {
				a.checkForUpdates(false)
			}
		}
		auto := !a.config.ManualUpdates
		ui.Column(c).FillWidth().Gap(6).Children(func() {
			components.Switch(c, &auto, locales.AutoCheckUpdates())
			if auto == a.config.ManualUpdates {
				a.setAutoCheck(auto)
			}
			components.Text(c, locales.AutoCheckUpdatesHint()).FontSize(12).TextColor(c.Theme().TextMuted).Wrap()
		})
	})
}

func (a *application) updateView(c *ui.Context) {
	u := &a.upd
	current := mygo.App.Version()
	switch u.phase {
	case updateIdle:
		updateHeading(c, lucide.CircleArrowUp, inkMuted, locales.Title(), locales.UpdateVersion(current))
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

// updateHeading says where updates got, after an icon in color, above a
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

// autoCheck checks for a newer version when checking by itself is on and
// the last check is a day old, unless one was found or is being checked.
func (a *application) autoCheck() {
	u := &a.upd
	if a.config.ManualUpdates || time.Since(u.checked) < autoCheckEvery {
		return
	}
	switch u.phase {
	case updateIdle, updateCurrent, updateDisabled:
		a.checkForUpdates(true)
	case updateFailed:
		if !u.installFailed {
			a.checkForUpdates(true)
		}
	}
}

// watchUpdates checks for updates soon after launch and, every half hour,
// whether the last check is a day old.
func (a *application) watchUpdates() {
	time.Sleep(10 * time.Second)
	t := time.NewTicker(30 * time.Minute)
	defer t.Stop()
	for {
		a.post(a.autoCheck)
		<-t.C
	}
}

// setAutoCheck turns checking for updates by itself on or off, checking at
// once when it turns on.
func (a *application) setAutoCheck(on bool) {
	a.config.ManualUpdates = !on
	a.save()
	if on {
		a.upd.checked = time.Time{}
		a.autoCheck()
	}
	a.redraw()
}

// checkForUpdates asks GitHub for a version newer than this one. A check of
// the app's own fails quietly, and tells the user of a version they were not
// told of yet.
func (a *application) checkForUpdates(auto bool) {
	u := &a.upd
	before := u.phase
	u.phase, u.found, u.installFailed, u.checked = updateChecking, nil, false, time.Now().Round(0)
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
				if auto {
					u.phase = before
				}
			case found == nil:
				u.phase = updateCurrent
			default:
				u.phase, u.found = updateAvailable, found
				if auto && u.notified != found.Version {
					u.notified = found.Version
					a.updates.notify(found.Version)
				}
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
