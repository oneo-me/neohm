# NeoHM

**English** | [简体中文](README.zh-CN.md)

NeoHM is a menu bar app for macOS, and a tray app for Windows and Linux, that controls the devices of a Home Assistant instance, with full controls for air conditioners. It is built in Go with [MyGo](https://github.com/egoist/mygo) and [Onevra](https://github.com/oneo-me/onevra).

Click the menu bar icon on macOS, or the tray icon on Windows, to open the panel; on Linux, choose **Open panel** from the tray icon's menu. It shows the devices as tiles by room, lit while they are on. A click on a tile turns its device on or off, and the arrow in its corner opens the device's page; a device that does not switch opens its page from the click. Settings opens in the same panel from its gear button. Right-click the icon for Quit NeoHM, which Linux lists in the same menu, or press ⌘Q (Ctrl+Q on Windows and Linux) while the panel is open.

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
- Show the interface in English, Simplified Chinese, Traditional Chinese, Japanese or Korean, following the system or the language chosen in Settings.

## Install

NeoHM needs a Home Assistant instance the computer can reach. Download the latest version from [Releases](https://github.com/oneo-me/neohm/releases/latest):

- **macOS 13 or later**: open the `.dmg` and drag NeoHM to Applications. NeoHM is not notarized by Apple, so macOS blocks its first launch. Allow it once in Terminal:

  ```sh
  xattr -dr com.apple.quarantine /Applications/NeoHM.app
  ```

- **Windows**: run `NeoHM.Setup.<version>.amd64.exe`, or the `arm64` one on ARM. NeoHM is not signed, so SmartScreen warns before its first launch: click **More info**, then **Run anyway**.
- **Linux**: install it for your user with the command below. Its tray icon needs `libayatana-appindicator3`, and the script names any other library missing. The `.deb` package installs it for every user instead, but does not update itself.

  ```sh
  curl -fsSL https://github.com/oneo-me/neohm/releases/latest/download/install.sh | sh
  ```

Later versions install from Settings: NeoHM checks once a day and notifies you of a new one, or click **Check for Updates** under Updates, then **Install and Relaunch**. NeoHM updates itself only where it can write: Applications on macOS, or where the installer or `install.sh` put it.

## Getting started

1. In Home Assistant, open your profile and create a **long-lived access token**.
2. Start NeoHM. The panel opens on Settings until Home Assistant is set up; later, open Settings from the panel's gear button or with ⌘, (Ctrl+, on Windows and Linux).
3. Enter the address, for example `http://homeassistant.local:8123`, and the token, then click **Connect**.
4. Close Settings to see your devices. Click a tile to switch its device, or its arrow to open the device's page.

The address, the token, the language chosen, whether updates are checked automatically and the air conditioner opened last are saved in `~/.neohm/config.json`, which only your user account can read.

## Contributing

NeoHM has been tried on only a few kinds of devices, so testing yours is the most useful help. Install it, open each of your devices, and [open an issue](https://github.com/oneo-me/neohm/issues) for any that shows or behaves wrongly, or for kinds that are missing. Include:

1. The kind of device, its brand and model, and the Home Assistant integration it uses.
2. Its entity's state and attributes, copied from **Developer tools → States** in Home Assistant.
3. What you did, what you expected and what NeoHM showed, with a screenshot.

Issues confirming that a kind of device works are welcome too. To change the code yourself, start with [Development](#development).

## Development

The Go module and the MyGo packaging configuration live in `source/`. Building needs Go 1.27 or later. Run Go and MyGo commands from that directory:

```sh
cd source
go tool mygo dev      # run with live reload
go test ./...         # headless tests; never opens a window
go tool mygo build    # this system's app and installer in ../dist
```

### Releasing

Pushing a tag `v<version>` builds the apps for macOS, Windows and Linux on GitHub Actions and publishes them, with their updates, as a release whose notes are that version's section of [CHANGELOG.md](CHANGELOG.md):

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
