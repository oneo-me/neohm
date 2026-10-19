# NeoHM

**English** | [简体中文](README.zh-CN.md)

NeoHM is a macOS menu bar app that controls the devices of a Home Assistant instance, with full controls for air conditioners. It is built in Go with [MyGo](https://github.com/egoist/mygo) and [Onevra](https://github.com/oneo-me/onevra).

Left-click the menu bar icon to open the panel. It shows the devices as tiles by room, lit while they are on. A click on a tile turns its device on or off, and the arrow in its corner opens the device's page; a device that does not switch opens its page from the click. Settings opens in the same panel from its gear button. Right-click the icon for Quit NeoHM, or press ⌘Q while the panel is open.

## Features

- Show the devices as tiles grouped by room, each with its main state and whether it is on.
- Switch lights, plugs, fans, covers and more from their tiles, and dim, color, move, play, lock or start a device from its page, with a light's colors and a cover's quarter positions a click away.
- Control an air conditioner's power, temperature, mode, fan speed, airflow, Eco Mode and Sleep, as far as it offers them.
- Show the target temperature of the air conditioner opened last beside the menu bar icon.
- Show the interface in English, Simplified Chinese, Traditional Chinese, Japanese or Korean, following macOS or the language chosen in Settings.

## Install

No release is published yet. Build from source with the [Development](#development) commands. It needs macOS 13 or later, Go 1.27 or later and a Home Assistant instance the Mac can reach.

## Getting started

1. In Home Assistant, open your profile and create a **long-lived access token**.
2. Start NeoHM. The panel opens on Settings until Home Assistant is set up; later, open Settings from the panel's gear button or with ⌘,.
3. Enter the address, for example `http://homeassistant.local:8123`, and the token, then click **Connect**.
4. Close Settings to see your devices. Click a tile to switch its device, or its arrow to open the device's page.

The address, the token, the language chosen and the air conditioner opened last are saved in `~/.neohm/config.json`, which only your user account can read.

## Development

The Go module and the MyGo packaging configuration live in `source/`. Run Go and MyGo commands from that directory:

```sh
cd source
go tool mygo dev      # run with live reload
go test ./...         # headless tests; never opens a window
go tool mygo build    # .app and .dmg in ../dist
```

See [ARCHITECTURE.md](ARCHITECTURE.md) for the layout and [AGENTS.md](AGENTS.md) for the working rules.
