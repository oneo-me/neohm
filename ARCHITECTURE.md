# NeoHM architecture

NeoHM is a macOS menu bar app that reads and controls the devices of one Home Assistant instance. It is a single Go module under `source/`, built with Onevra components on MyGo. Its interface is one panel anchored to the menu bar icon, whose pages are Home, a device's page, an air conditioner's controls and Settings. Home shows the devices as tiles by area.

## Code map

| Path | Responsibility |
| --- | --- |
| `source/main.go` | Entry point: the MyGo lifecycle, the menu bar activation policy and the configuration directory. |
| `source/app.go` | Application state, the poller, merging reads with local changes and sending service calls. |
| `source/tray.go` | The menu bar icon and menus, the panel's lifecycle and navigation, and the window used without a menu bar icon. |
| `source/panel.go` | The panel's frame, theme and palette, the frame every page shares, and the controls of one air conditioner. |
| `source/home.go` | Home's tiles by area, which switch their devices and open their pages, and what stands in for them before Home Assistant answers. |
| `source/tiles.go` | How each kind of device shows on its tile and page. |
| `source/device.go` | A device's page, laid out as an air conditioner's controls: its drawing, its main value, the slider of its level, a light's colors and its other actions or readings along the bottom. |
| `source/controls.go` | The controls every device's page shares: the power button at the top right, the main value in the middle, the slider with its Liquid Glass thumb and the stops below it, and the tiles along the bottom. |
| `source/colors.go` | A light's swatches below its slider and the picker of any color. |
| `source/scenes.go` | The drawing of each kind of device at the top of its page. |
| `source/settings.go` | The Settings page. |
| `source/update.go` | Settings' updates, which check GitHub for a newer version, by themselves once a day, notify the user of it and install it. |
| `source/config.go` | Reading and atomically saving the configuration file. |
| `source/internal/ha/` | The REST and WebSocket clients, what entities' states, actions, quick switches, levels and their stops, and lights' colors mean, which tiles Home shows, and the parsed climate entity. |
| `source/locales/` | MyGo message catalogs, `<tag>.json`, and the accessors `go tool mygo generate` writes from them. Their `mygo.*` keys name the app and translate MyGo's own strings where MyGo has no translation. |
| `source/resources/darwin/` | Translations of `Info.plist` strings, such as the reason macOS shows for local network access, in `<tag>.lproj/InfoPlist.strings`. |
| `source/icons/` | The Icon Composer app icon `AppIcon.icon`, which `mygo.json`'s `macos.icon` names, its flattened `AppIcon.png` for `mygo.json`'s `icon` and Neorio's project icon, and the unused alternatives `AppIcon-1.icon` and `AppIcon-2.icon`. They stay outside `resources/`, whose contents ship in the app. |
| `website/` | The landing page: a React and Vite site in English, Simplified Chinese, Traditional Chinese, Japanese and Korean, with the app's screenshots and install steps. |
| `.github/workflows/website.yml` | Builds `website/` with Bun and deploys it to GitHub Pages. |
| `.github/workflows/release.yml` | On a `v*` tag, builds the universal app, signs its update and publishes both as a GitHub release, its notes from `CHANGELOG.md`. |

## Boundaries

- Present the app from the menu bar: a left click opens the panel and a right click offers Quit, which ⌘Q also does in the panel.
- Keep Settings a page of the panel, never a window of its own. A window shows the panel's pages only when the menu bar icon cannot be created.
- Build each page in the shared frame: its header floats over content that scrolls to the panel's edges and blurs beneath it.
- Keep the panel's surface opaque and drawn by the panel. The header's blur is computed in MyGo from what the panel paints, so it cannot blur a system material.
- Open the panel on Home, or on Settings until Home Assistant is set up. Closing the panel returns it there.
- Run network work off the main thread and apply its results on the main thread. Builders never start work.
- Keep the access token in the configuration file, readable only by the user, rather than the keychain.
- When the language changes, set again what was set only once: the menus, the menu bar icon's tooltip and the window's title.
- Apply a language only when Settings changes it. MyGo keeps it across launches, and `MYGO_LOCALE` still overrides it.
- Show a local change at once, and keep it over reads that started before it until Home Assistant reports it.
- Keep what Home Assistant's domains, states and services mean in `internal/ha`, and how the panel shows them in the main package.
- Treat the registries as optional: without them, Home still shows each switch and main control, ungrouped.
- Check for updates only from Settings, or once a day unless Settings turns that off. A newer version found by itself is notified once and marks Home's Settings button. Apps are signed ad hoc, not notarized, and accept only updates signed with the key of `updates.publicKey`; a release's version, tag and changelog section must agree.
- Keep the website a static Vite build in `website/`. CI deploys it to GitHub Pages; the app build never includes it.

## Core flows

**Polling.** One goroutine reads Home Assistant while the panel is open, every 3 seconds: every entity while Home or a device's page shows, otherwise only the air conditioner shown. It reads nothing while the panel is closed, as the menu bar icon never changes; opening the panel, navigating and connecting wake it at once.

**Arranging Home.** Which device and area each entity belongs to, and what Home Assistant's dashboards hide, come only from the registries of its WebSocket API. The poller reads them in one short session per connection, and again every five minutes. `ha.Arrange` then shows each device's main controls, or one tile of a sensor device's readings, by area.

**Control.** A click on a tile that switches calls one service of its entity's domain at once, as the power button in a device's page header, or a click on the device drawn below it, does; the tile's arrow, or a click on any other tile, opens the device's page. A device's page calls the others, and sends a level, or a color dragged in the picker, once its control rests; a swatch sends its color at once. An air conditioner's controls call `climate` services on the device shown, and a click on its unit turns it on or off. Its temperature is sent once the slider rests, and a temperature still waiting goes to its own device when another one opens.

## Data and state

The [configuration file](README.md#getting-started) holds the address, the access token, the language chosen, whether updates are checked only by hand and the air conditioner opened last. Each save replaces it atomically with a file only the user can read. A read of every entity that no longer contains that air conditioner forgets it. Device states and the registries stay in memory and are read again after each launch.

## Verification

Tests never open a window. Logic tests cover the REST and WebSocket clients, which tiles Home shows, the configuration file, the catalogs and how reads merge with local changes. Headless UI tests drive the panel's pages with MyGo's tester.

From `source/`, `HA_SNAPSHOT_DIR=<dir> go test -run Snapshot .` renders each page and the menu bar icon to PNG for review. Prefix `MYGO_LOCALE=<tag>` to render another language. `go tool mygo dev -lang <tag>` runs the app in one. Installing an update needs a published release and a packaged app, and so does its notification, so both are reviewed by hand. The menu bar icon, the panel's placement, keyboard focus, each language's menus and system prompts, and a live Home Assistant need manual review in a desktop session.

From `website/`, `bun run build` type-checks and bundles the site. Review layout changes in a browser at every language and viewport size.

## Subjects

No subject documents exist yet. List each one added under `docs/architecture/` here.

| Document | When to read |
| --- | --- |
