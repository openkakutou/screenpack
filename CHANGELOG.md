# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.6.0] - 2026-11-20

### Added
- `screenpack.DefaultFS` returns the default screenpack built into the program, so a game can fall back on it without a folder on disk.

## [0.5.0] - 2026-11-19

### Added
- `screenpack.Resolve` completes a screenpack from the default one: a missing screen is the default's, a part a screen leaves out is the default screen's part, and `Source` says which pack each screen and part comes from.

## [0.4.0] - 2026-11-19

### Added
- The default screenpack in `default/`: every screen of the format, original images and sounds (MIT, made by `tools/default-assets/generate.py`) and Noto Sans fonts (OFL-1.1), each asset's license recorded in `default/ASSETS.md`.

## [0.3.0] - 2026-11-19

### Added
- `screenpack.Load` reads a screenpack folder (or an in-memory file set) into a typed model: every file checked against its schema, then every referenced file, font, text style, sprite-sheet frame and animation checked to exist. All problems are reported at once, each with its file, JSON path and reason; unknown files in `screens/` are warnings, not errors.

## [0.2.0] - 2026-11-19

### Added
- JSON Schemas for format version 1 in `schema/v1/`: `screenpack.json`, sprite sheets and every screen, published at `https://raw.githubusercontent.com/openkakutou/screenpack/main/schema/v1/<kind>.schema.json` for editor completion. The `schema` Go package validates a JSONC document against them and names the JSON path of each problem.

## [0.1.0] - 2026-11-19

### Added
- Screenpack format specification (`docs/format.md`), format version 1: root `screenpack.json`, one JSONC file per screen with per-screen fallback to `default`, anchor + offset layout in screen-height units, sprite sheets, text styles, transitions, menus, music and sounds, every screen's parts, a complete example and the validation rules.

[Unreleased]: https://github.com/openkakutou/screenpack/compare/v0.6.0...HEAD
[0.6.0]: https://github.com/openkakutou/screenpack/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/openkakutou/screenpack/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/openkakutou/screenpack/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/openkakutou/screenpack/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/openkakutou/screenpack/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/openkakutou/screenpack/releases/tag/v0.1.0
