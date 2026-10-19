# NeoHM

NeoHM is a menu bar app for macOS, and a tray app for Windows and Linux, that controls Home Assistant devices, built with Onevra components on MyGo. Read [ARCHITECTURE.md](ARCHITECTURE.md) before changing code.

## Principles

- Keep all Go code, the MyGo packaging configuration and the message catalogs in `source/`; run Go and MyGo commands there. Keep documentation at the repository root.
- Keep the website's React and Vite sources in `website/`; run Bun commands there.
- Use Onevra components instead of MyGo's styled widgets, in the panel's own theme.
- Add user-visible text to every catalog in `source/locales/`, then run `go tool mygo generate`.

## Working with the user

- Do not commit, tag, push or publish without explicit authorization.
- An explicit request to publish a named version authorizes the [Releasing](README.md#releasing) steps, including the release commit, tag and pushes. Requests to prepare a release or edit its notes do not.

## Verification

Build and test from `source/`, following [Development](README.md#development). A green headless run does not prove native behavior: review UI changes in a desktop session.

## Documentation

- Update a document only when a fact its reader needs changes.
- Keep each document within its budget and written for its reader.
- Repair a document you touch instead of appending to it.
- Keep `README.md` and `README.zh-CN.md` equivalent; everything else, including code comments, is English.
- Keep documentation on the code branch and reconcile it during merges.
- In summaries and release notes, describe only what was added, changed or removed; omit testing details unless asked.

## Changelog

[CHANGELOG.md](CHANGELOG.md) is the release notes on GitHub and in the update popover, written for someone deciding whether to install a version.

- Record each completed outcome a user notices, or must act on, under `Unreleased`.
- Write one short clause per outcome, starting with a verb. Do not extend it with a colon, semicolon or "which" clause.
- Omit refactors, tests, tooling, CI, the website, styling and copy tweaks. The website deploys on its own, outside versions.
- Merge related entries and refine an existing one instead of appending. Never present an unreleased experiment as a change or a removal.
- Until the first release ships, describe what the app does, not how it evolved.
- Prefix breaking changes with `Breaking:` and name the action required.
- Order bullets as `Breaking:`, then Add, then other changes, then Fix. Keep one list per version, without subheadings.
- At release, rename `[Unreleased]` to the version, as [Releasing](README.md#releasing) describes. Change published sections only to correct facts.
