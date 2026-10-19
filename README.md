# HA Climate

**English** | [简体中文](README.zh-CN.md)

HA Climate is a macOS menu bar app that controls the air conditioners (`climate.*` entities) of a Home Assistant instance. It is built in Go with [MyGo](https://github.com/egoist/mygo) and [Onevra](https://github.com/oneo-me/onevra).

Left-click the menu bar icon to open the panel. It lists the air conditioners, and choosing one opens its controls. Settings opens in the same panel from its gear button. Right-click the icon for Quit.

## Features

- List the air conditioners with their mode, room temperature and target temperature.
- Control power, temperature, mode, fan speed, airflow, Eco Mode and Sleep, as far as each unit offers them.
- Show the target temperature of the air conditioner opened last beside the menu bar icon.

## Install

No release is published yet. Build from source with the [Development](#development) commands. It needs macOS 13 or later, Go 1.27 or later and a Home Assistant instance the Mac can reach.

## Getting started

1. In Home Assistant, open your profile and create a **long-lived access token**.
2. Start HA Climate. The panel opens on Settings until Home Assistant is set up; later, open Settings from the panel's gear button or with ⌘,.
3. Enter the address, for example `http://homeassistant.local:8123`, and the token, then click **Connect**.
4. Close Settings and choose an air conditioner from the list.

The address and the air conditioner opened last are saved in `~/.ha-go/config.json`. The token is kept in the system keychain.

## Development

The Go module and the MyGo packaging configuration live in `source/`. Run Go and MyGo commands from that directory:

```sh
cd source
go tool mygo dev      # run with live reload
go test ./...         # headless tests; never opens a window
go tool mygo build    # .app and .dmg in ../dist
```

See [ARCHITECTURE.md](ARCHITECTURE.md) for the layout and [AGENTS.md](AGENTS.md) for the working rules.
