package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"

	"github.com/oneo-me/neohm/source/locales"
)

// config holds the preferences saved in ~/.neohm, the access token
// included.
type config struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	// Entity is the air conditioner opened last, which the panel shows
	// in place of Home.
	Entity string `json:"entity"`
	// Language is the catalog the app shows, as "zh-Hans", or "" for the
	// one that suits the system's languages best.
	Language string `json:"language,omitempty"`
	// ManualUpdates checks for updates only when Settings asks, not by
	// itself every day.
	ManualUpdates bool `json:"manualUpdates,omitempty"`
	// AreaOrder lists Home's areas by ID, "" for the devices of no area,
	// as the user dragged them; the areas it misses follow in Home
	// Assistant's order. TileOrder lists each area's tiles by entity the
	// same way, and Collapsed the areas whose tiles Home folds away.
	// Pinned lists the entities whose tiles Home shows above the areas,
	// in the order it shows them.
	AreaOrder []string            `json:"areaOrder,omitempty"`
	TileOrder map[string][]string `json:"tileOrder,omitempty"`
	Collapsed []string            `json:"collapsed,omitempty"`
	Pinned    []string            `json:"pinned,omitempty"`
}

func loadConfig(path string) (config, error) {
	var c config
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return c, nil
	}
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return config{}, err
	}
	if !slices.Contains(locales.Catalog.Languages(), c.Language) {
		c.Language = ""
	}
	return c, nil
}

// saveConfig replaces the file atomically so an interrupted save leaves the
// old preferences. The file is the user's alone to read, as os.CreateTemp
// makes it, for it holds the token.
func saveConfig(path string, c config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".config-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
