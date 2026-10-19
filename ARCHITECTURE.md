# HA Climate architecture

HA Climate is a macOS menu bar app that reads and controls the devices of one Home Assistant instance. It is a single Go module under `source/`, built with Onevra components on MyGo. Its interface is one panel anchored to the menu bar icon, whose pages are Home, an air conditioner's controls and Settings. Home shows the devices as tiles by area.

## Code map

| Path | Responsibility |
| --- | --- |
| `source/main.go` | Entry point: the MyGo lifecycle, the menu bar activation policy and the configuration directory. |
| `source/app.go` | Application state, the poller, merging reads with local changes and sending service calls. |
| `source/tray.go` | The menu bar icon and menus, the panel's lifecycle and navigation, and the window used without a menu bar icon. |
| `source/panel.go` | The panel's frame, theme and palette, and the controls of one air conditioner. |
| `source/home.go` | Home's tiles by area, and what stands in for them before Home Assistant answers. |
| `source/tiles.go` | How each kind of device shows on a tile, and what a click on it does. |
| `source/settings.go` | The Settings page. |
| `source/config.go` | Reading and atomically saving the configuration file. |
| `source/internal/ha/` | The REST and WebSocket clients, what entities' states and actions mean, which tiles Home shows, and the parsed climate entity. |
| `source/locales/` | MyGo message catalogs, `<tag>.json`, and the accessors `go tool mygo generate` writes from them. Their `mygo.*` keys name the app and translate MyGo's own strings where MyGo has no translation. |
| `source/resources/darwin/` | Translations of `Info.plist` strings, such as the reason macOS shows for local network access, in `<tag>.lproj/InfoPlist.strings`. |

## Boundaries

- Present the app from the menu bar: a left click opens the panel and a right click offers Quit.
- Keep Settings a page of the panel, never a window of its own. A window shows the panel's pages only when the menu bar icon cannot be created.
- Open the panel on Home, or on Settings until Home Assistant is set up. Closing the panel returns it there.
- Run network work off the main thread and apply its results on the main thread. Builders never start work.
- Keep the access token in the configuration file, readable only by the user, rather than the keychain.
- When the language changes, set again what was set only once: the menus, the menu bar icon's tooltip and the window's title.
- Show a local change at once, and keep it over reads that started before it until Home Assistant reports it.
- Keep what Home Assistant's domains, states and services mean in `internal/ha`, and how the panel shows them in the main package.
- Treat the registries as optional: without them, Home still shows each switch and main control, ungrouped.

## Core flows

**Polling.** One goroutine reads Home Assistant: every entity while Home shows, otherwise only the air conditioner opened last, which the menu bar title follows. It reads every 3 seconds while the panel is open and every 20 seconds otherwise; opening the panel, navigating and connecting wake it at once.

**Arranging Home.** Which device and area each entity belongs to, and what Home Assistant's dashboards hide, come only from the registries of its WebSocket API. The poller reads them in one short session per connection, and again every five minutes. `ha.Arrange` then shows each device's main controls, or one tile of a sensor device's readings, by area.

**Control.** A tile calls one service of its entity's domain, at once for a switch or a button, else from a menu of the entity's actions. An air conditioner's controls call `climate` services on the device shown. Its temperature is sent once the slider rests, and a temperature still waiting goes to its own device when another one opens.

## Data and state

The [configuration file](README.md#getting-started) holds the address, the access token and the air conditioner opened last. Each save replaces it atomically with a file only the user can read. A read of every entity that no longer contains that air conditioner forgets it. Device states and the registries stay in memory and are read again after each launch.

## Verification

Tests never open a window. Logic tests cover the REST and WebSocket clients, which tiles Home shows, the configuration file, the catalogs and how reads merge with local changes. Headless UI tests drive the panel's pages with MyGo's tester.

From `source/`, `HA_SNAPSHOT_DIR=<dir> go test -run Snapshot .` renders each page and the menu bar icon to PNG for review. Prefix `MYGO_LOCALE=<tag>` to render another language. `go tool mygo dev -lang <tag>` runs the app in one. The menu bar icon, the panel's placement, keyboard focus, each language's menus and system prompts, and a live Home Assistant need manual review in a desktop session.

## Subjects

No subject documents exist yet. List each one added under `docs/architecture/` here.

| Document | When to read |
| --- | --- |
