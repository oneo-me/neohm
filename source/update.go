package main

import (
	"context"
	"errors"
	"log"
	"runtime"
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

// updatesGroup, in Settings, is one row that says where updates got and
// does what they allow next, checking, installing or trying again, with the
// notes of a newer version found, or how far its download got, below it. Its
// switch turns checking by itself on or off.
func (a *application) updatesGroup(c *ui.Context) {
	u := &a.upd
	current := mygo.App.Version()
	icon, color, title, detail, detailColor := lucide.RefreshCw, accent, locales.CheckForUpdates(), locales.UpdateVersion(current), inkMuted
	act := func() { a.checkForUpdates(false) }
	switch u.phase {
	case updateChecking:
		icon, color, title, act = nil, quietBadge, locales.UpdateChecking(), nil
	case updateCurrent:
		icon, color, title = lucide.CircleCheck, green, locales.UpdateCurrent()
	case updateDisabled:
		icon, color, title, detail, act = lucide.CircleAlert, quietBadge, locales.UpdateDisabled(), updateDisabledBody(), nil
	case updateFailed:
		icon, color, title, detail, detailColor = lucide.TriangleAlert, danger, locales.UpdateCheckError(), locales.UpdateRetry(), accent
		if u.installFailed {
			title, act = locales.UpdateInstallError(), a.installUpdate
		}
	case updateAvailable:
		icon, title, detail, detailColor, act = lucide.CircleArrowUp, locales.UpdateAvailable(u.found.Version), locales.UpdateInstall(), accent, a.installUpdate
	case updateInstalling:
		icon, title, detail, act = lucide.Download, locales.UpdateInstalling(u.found.Version), locales.UpdateRelaunch(), nil
	}
	settingsGroup(c, locales.Updates(), locales.AutoCheckUpdatesHint(), func() {
		clicked := pressRow(c, title, false, func() {
			rowBody(c, icon, color, title, detail, detailColor, func() {
				auto := !a.config.ManualUpdates
				glassSwitch(c, &auto, locales.AutoCheckUpdates(), locales.AutoCheckUpdatesHint())
				if auto == a.config.ManualUpdates {
					a.setAutoCheck(auto)
				}
			})
		})
		if clicked && act != nil {
			act()
		}
		switch {
		case u.phase == updateAvailable && u.found.Notes != "":
			// The notes continue the row above, below its text.
			components.ScrollView(c, components.ScrollViewOptions{}).FillWidth().MaxHeight(180).Children(func() {
				ui.Column(c).FillWidth().Padding(0, 14, 12, 0).PaddingStart(rowInset).Children(func() {
					markdown.View(c, u.found.Notes, u.notes)
				})
			})
		case u.phase == updateInstalling:
			ui.Column(c).FillWidth().Padding(0, 14, 16, 0).PaddingStart(rowInset).Children(func() { progressBar(c, u.progress) })
		}
	})
}

// updateDisabledBody says how to install the app where it can update itself.
func updateDisabledBody() string {
	switch runtime.GOOS {
	case "windows":
		return locales.UpdateDisabledBodyWindows()
	case "linux":
		return locales.UpdateDisabledBodyLinux()
	}
	return locales.UpdateDisabledBody()
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
