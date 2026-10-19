package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/egoist/mygo"
	"github.com/oneo-me/ha-go/source/locales"
)

func main() {
	mygo.App.SetName(locales.Title())
	if !mygo.App.RequestSingleInstanceLock() {
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Join(home, ".ha-go")
	mygo.App.SetPath(mygo.PathUserData, dir)
	a := newApplication(filepath.Join(dir, "config.json"))
	a.config, err = loadConfig(a.path)
	if err != nil {
		a.report("loadError", err)
	}
	a.address, a.token = a.config.Address, a.config.Token
	mygo.App.SetActivationPolicy(mygo.ActivationPolicyAccessory)
	mygo.App.OnWindowAllClosed(func() {})
	mygo.App.OnBeforeQuit(func(*mygo.QuitEvent) { a.quitting = true })
	mygo.App.OnWillQuit(func(*mygo.QuitEvent) { a.destroyTray() })
	mygo.App.OnSecondInstance(func([]string, string) {
		if a.ready {
			a.reopen()
		}
	})
	mygo.App.OnLocaleChanged(func(string) {
		if a.ready {
			a.relabel()
		}
	})
	mygo.App.WhenReady(a.start)
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
