# NeoHM

NeoHM is a macOS menu bar app that controls Home Assistant devices, built with Onevra components on MyGo. Read [ARCHITECTURE.md](ARCHITECTURE.md) before changing code.

## Principles

- Keep all Go code, the MyGo packaging configuration and the message catalogs in `source/`; run Go and MyGo commands there. Keep documentation at the repository root.
- Keep the website's React and Vite sources in `website/`; run Bun commands there.
- Use Onevra components instead of MyGo's styled widgets, in the panel's own theme.
- Add user-visible text to every catalog in `source/locales/`, then run `go tool mygo generate`.

## Working with the user

- Do not commit, tag, push or publish without explicit authorization.

## Verification

Build and test from `source/`, following [Development](README.md#development). A green headless run does not prove native behavior: review UI changes in a desktop session.

## Documentation

- Update a document only when a fact its reader needs changes.
- Keep each document within its budget and written for its reader.
- Record completed user-visible outcomes as short `Unreleased` bullets in [CHANGELOG.md](CHANGELOG.md), merging related ones.
- Repair a document you touch instead of appending to it.
- Keep `README.md` and `README.zh-CN.md` equivalent; everything else, including code comments, is English.
- Keep documentation on the code branch and reconcile it during merges.
