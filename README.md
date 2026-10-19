# HA Climate

A macOS menu bar app that controls a Home Assistant air conditioner (`climate.*` entity), built in Go with [MyGo](https://github.com/egoist/mygo) and [Onevra](https://github.com/oneo-me/onevra).

The menu bar shows the target temperature beside the icon. Click it for the control panel. Right-click it for Refresh, Settings and Quit.

## Setup

1. In Home Assistant, open your profile and create a **long-lived access token**.
2. Open **Settings…**, then enter the address (for example `http://homeassistant.local:8123`) and the token.
3. Click **Connect**, then choose the air conditioner.

The address and entity are saved in `~/Library/Application Support/ha-climate/config.json`. The token is kept in the system keychain.

## What the panel controls

| Control | Home Assistant service |
| --- | --- |
| Power | `climate.turn_off` / `climate.set_hvac_mode` |
| Temperature slider, − and + | `climate.set_temperature` (sent once the slider rests) |
| Mode | `climate.set_hvac_mode` |
| Fan speed | `climate.set_fan_mode` |
| Airflow | `climate.set_swing_mode` |
| Eco Mode, Sleep | `climate.set_preset_mode` with `eco` / `sleep`, back to `none` |

A control is disabled when the entity does not offer it. The state is read every 3 seconds while the panel is open and every 20 seconds otherwise.

## Development

```sh
cd source
go tool mygo dev      # run with live reload
go test ./...         # headless tests; never opens a window
HA_SNAPSHOT_DIR=/tmp go test -run Snapshot .   # render the panel to PNG
go tool mygo build    # .app and .dmg in ../dist
```
