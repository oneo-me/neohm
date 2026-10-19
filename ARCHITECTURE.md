# HA Climate architecture

HA Climate is a macOS menu bar app that reads and controls the `climate` entities of one Home Assistant instance through its REST API. It is a single Go module under `source/`, built with Onevra components on MyGo. Its interface is one panel anchored to the menu bar icon, whose pages are the device list, one device's controls and Settings.

## Code map

| Path | Responsibility |
| --- | --- |
| `source/main.go` | Entry point: the MyGo lifecycle, the menu bar activation policy, the configuration directory and the keychain. |
| `source/app.go` | Application state, the poller, merging reads with local changes and sending service calls. |
| `source/tray.go` | The menu bar icon and menus, the panel's lifecycle and navigation, and the window used without a menu bar icon. |
| `source/panel.go` | The panel's frame and theme, and the controls of one air conditioner. |
| `source/list.go` | The device list and what stands in for it before Home Assistant answers. |
| `source/settings.go` | The Settings page. |
| `source/config.go` | Reading and atomically saving the configuration file. |
| `source/internal/ha/` | The Home Assistant REST client and the parsed climate entity. |
| `source/locales/` | MyGo message catalogs, `<tag>.json`, and the accessors `go tool mygo generate` writes from them. |

## Boundaries

- Present the app from the menu bar: a left click opens the panel and a right click offers Quit.
- Keep Settings a page of the panel, never a window of its own. A window shows the panel's pages only when the menu bar icon cannot be created.
- Open the panel on the device list, or on Settings until Home Assistant is set up. Closing the panel returns it there.
- Run network work off the main thread and apply its results on the main thread. Builders never start work.
- Keep the access token in the system keychain, never in the configuration file.
- Show a local change at once, and keep it over reads that started before it until Home Assistant reports it.

## Core flows

**Polling.** One goroutine reads Home Assistant: every climate entity while the device list shows, otherwise only the device opened last, which the menu bar title follows. It reads every 3 seconds while the panel is open and every 20 seconds otherwise; opening the panel, navigating and connecting wake it at once.

**Control.** Each control calls one `climate` service on the device shown: `turn_off`, `set_hvac_mode`, `set_temperature`, `set_fan_mode`, `set_swing_mode` or `set_preset_mode`. Eco Mode and Sleep switch their preset on, and back to the neutral one. The temperature is sent once the slider rests, and a temperature still waiting goes to its own device when another one opens.

## Data and state

The [configuration file](README.md#getting-started) holds the address and the device opened last, and each save replaces it atomically. A list read that no longer contains that device forgets it. The access token lives only in the system keychain. Device states stay in memory and are read again after each launch.

## Verification

Tests never open a window. Logic tests cover the REST client, the configuration file and how reads merge with local changes; headless UI tests drive the panel's pages with MyGo's tester. From `source/`, `HA_SNAPSHOT_DIR=<dir> go test -run Snapshot .` renders each page to PNG for review. The menu bar icon, the panel's placement, keyboard focus and the keychain need manual review in a desktop session.

## Subjects

No subject documents exist yet. List each one added under `docs/architecture/` here.

| Document | When to read |
| --- | --- |
