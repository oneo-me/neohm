# NeoHM

**English** | [简体中文](README.zh-CN.md)

NeoHM is a macOS menu bar app that controls the devices of a Home Assistant instance, with full controls for air conditioners. It is built in Go with [MyGo](https://github.com/egoist/mygo) and [Onevra](https://github.com/oneo-me/onevra).

Left-click the menu bar icon to open the panel. It shows the devices as tiles by room, lit while they are on. A click on a tile turns its device on or off, and the arrow in its corner opens the device's page; a device that does not switch opens its page from the click. Settings opens in the same panel from its gear button. Right-click the icon for Quit NeoHM, or press ⌘Q while the panel is open.

<p align="center">
  <img src="docs/screenshots/home.jpg" width="24%" alt="Devices by room">
  <img src="docs/screenshots/air-conditioner.jpg" width="24%" alt="An air conditioner's controls">
  <img src="docs/screenshots/light.jpg" width="24%" alt="A light's page">
  <img src="docs/screenshots/cover.jpg" width="24%" alt="A curtain's page">
</p>

## Features

- Show the devices as tiles grouped by room, each with its main state and whether it is on.
- Switch lights, plugs, fans, covers and more from their tiles, and dim, color, move, play, lock or start a device from its page, with a light's colors and a cover's quarter positions a click away.
- Control an air conditioner's power, temperature, mode, fan speed, airflow, Eco Mode and Sleep, as far as it offers them.
- Show the interface in English, Simplified Chinese, Traditional Chinese, Japanese or Korean, following macOS or the language chosen in Settings.

## Install

NeoHM needs macOS 13 or later and a Home Assistant instance the Mac can reach.

1. Download the `.dmg` of the latest version from [Releases](https://github.com/oneo-me/ha-go/releases/latest) and drag NeoHM to Applications.
2. NeoHM is not notarized by Apple, so macOS blocks its first launch. Allow it once in Terminal:

   ```sh
   xattr -dr com.apple.quarantine /Applications/NeoHM.app
   ```

Later versions install from the panel: click the arrow button left of the gear to check for one, then **Install and Relaunch**. NeoHM updates itself only from Applications, or another folder you can write to.

## Getting started

1. In Home Assistant, open your profile and create a **long-lived access token**.
2. Start NeoHM. The panel opens on Settings until Home Assistant is set up; later, open Settings from the panel's gear button or with ⌘,.
3. Enter the address, for example `http://homeassistant.local:8123`, and the token, then click **Connect**.
4. Close Settings to see your devices. Click a tile to switch its device, or its arrow to open the device's page.

The address, the token, the language chosen and the air conditioner opened last are saved in `~/.neohm/config.json`, which only your user account can read.

## Development

The Go module and the MyGo packaging configuration live in `source/`. Building needs Go 1.27 or later. Run Go and MyGo commands from that directory:

```sh
cd source
go tool mygo dev      # run with live reload
go test ./...         # headless tests; never opens a window
go tool mygo build    # .app and .dmg in ../dist
```

### Releasing

Pushing a tag `v<version>` builds a universal app on GitHub Actions and publishes it, with its update, as a release whose notes are that version's section of [CHANGELOG.md](CHANGELOG.md):

1. Set `version` in `source/mygo.json`, and rename `[Unreleased]` in the changelog to `[<version>] - <date>`.
2. Commit, then `git tag v<version>` and `git push origin v<version>`.

The workflow signs updates with the `MYGO_UPDATER_PRIVATE_KEY` secret, the key whose public half is `updates.publicKey` in `mygo.json`. Installed apps accept only updates signed with it, so keep a copy outside the repository.

See [ARCHITECTURE.md](ARCHITECTURE.md) for the layout and [AGENTS.md](AGENTS.md) for the working rules.

## Website

The landing page is a React and Vite app in `website/`, built with [Bun](https://bun.sh). Run its commands from that directory:

```sh
cd website
bun install       # once
bun run dev       # local preview
bun run build     # static site in dist
```

Pushes to `main` that change `website/` deploy it to GitHub Pages with Actions. Set Settings → Pages → Source to **GitHub Actions** once.
