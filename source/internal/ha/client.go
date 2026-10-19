// Package ha talks to Home Assistant's REST API.
package ha

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ErrUnauthorized reports a rejected access token.
var ErrUnauthorized = errors.New("home assistant rejected the access token")

// Client calls one Home Assistant instance with a long-lived access token.
type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

// NewClient returns a client for base, such as "http://homeassistant.local:8123".
func NewClient(base, token string) *Client {
	return &Client{BaseURL: NormalizeURL(base), Token: strings.TrimSpace(token), HTTP: &http.Client{Timeout: 10 * time.Second}}
}

// NormalizeURL trims spaces and trailing slashes and adds a missing scheme.
func NormalizeURL(base string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	if base != "" && !strings.Contains(base, "://") {
		base = "http://" + base
	}
	return base
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	if c.BaseURL == "" || c.Token == "" {
		return errors.New("home assistant is not configured")
	}
	if _, err := url.ParseRequestURI(c.BaseURL); err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	var r io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		r = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, r)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	switch {
	case res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden:
		return ErrUnauthorized
	case res.StatusCode >= 300:
		msg, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return fmt.Errorf("%s %s: %s %s", method, path, res.Status, strings.TrimSpace(string(msg)))
	}
	if out == nil {
		_, err = io.Copy(io.Discard, res.Body)
		return err
	}
	return json.NewDecoder(res.Body).Decode(out)
}

// TemperatureUnit returns the instance's temperature unit, as "°C".
func (c *Client) TemperatureUnit(ctx context.Context) (string, error) {
	var cfg struct {
		UnitSystem struct {
			Temperature string `json:"temperature"`
		} `json:"unit_system"`
	}
	if err := c.do(ctx, http.MethodGet, "/api/config", nil, &cfg); err != nil {
		return "", err
	}
	return cfg.UnitSystem.Temperature, nil
}

// States reads the state of every entity, in Home Assistant's order.
func (c *Client) States(ctx context.Context) ([]State, error) {
	var states []State
	if err := c.do(ctx, http.MethodGet, "/api/states", nil, &states); err != nil {
		return nil, err
	}
	return states, nil
}

// State reads the state of one entity.
func (c *Client) State(ctx context.Context, entity string) (State, error) {
	var s State
	err := c.do(ctx, http.MethodGet, "/api/states/"+url.PathEscape(entity), nil, &s)
	return s, err
}

// Call calls a service of domain, such as "climate" and "set_temperature",
// on entity with the extra fields of data.
func (c *Client) Call(ctx context.Context, domain, service, entity string, data map[string]any) error {
	body := map[string]any{"entity_id": entity}
	for k, v := range data {
		body[k] = v
	}
	return c.do(ctx, http.MethodPost, "/api/services/"+url.PathEscape(domain)+"/"+url.PathEscape(service), body, nil)
}
