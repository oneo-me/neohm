package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// config holds the preferences saved in ~/.ha-go, the access token
// included.
type config struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	// Entity is the device opened last, which the menu bar title follows.
	Entity string `json:"entity"`
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
