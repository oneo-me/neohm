package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/egoist/mygo"
	"github.com/oneo-me/ha-go/source/locales"
)

const tokenKey = "access-token"

func main() {
	mygo.App.SetName(locales.Title())
	if !mygo.App.RequestSingleInstanceLock() {
		return
	}
	base, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Join(base, "ha-climate")
	mygo.App.SetPath(mygo.PathUserData, dir)
	a := newApplication(filepath.Join(dir, "config.json"), keychain{})
	a.config, err = loadConfig(a.path)
	if err != nil {
		a.report("loadError", err)
	}
	if token, err := a.secrets.Get(tokenKey); err == nil {
		a.token = string(token)
	} else if !errors.Is(err, mygo.ErrSecretNotFound) {
		a.report("secretError", err)
	}
	a.address = a.config.Address
	mygo.App.SetActivationPolicy(mygo.ActivationPolicyAccessory)
	mygo.App.OnWindowAllClosed(func() {})
	mygo.App.OnBeforeQuit(func(*mygo.QuitEvent) { a.quitting = true })
	mygo.App.OnWillQuit(func(*mygo.QuitEvent) { a.destroyTray() })
	mygo.App.OnSecondInstance(func([]string, string) {
		if a.ready {
			a.reopen()
		}
	})
	mygo.App.WhenReady(a.start)
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

// keychain stores the access token in the system's credential store.
type keychain struct{}

func (keychain) Get(key string) ([]byte, error)     { return mygo.Secrets.Get(key) }
func (keychain) Set(key string, value []byte) error { return mygo.Secrets.Set(key, value) }
