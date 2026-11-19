# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- JSON Schemas for format version 1 in `schema/v1/`: `screenpack.json`, sprite sheets and every screen, published at `https://raw.githubusercontent.com/openkakutou/screenpack/main/schema/v1/<kind>.schema.json` for editor completion. The `schema` Go package validates a JSONC document against them and names the JSON path of each problem.

## [0.1.0] - 2026-11-19

### Added
- Screenpack format specification (`docs/format.md`), format version 1: root `screenpack.json`, one JSONC file per screen with per-screen fallback to `default`, anchor + offset layout in screen-height units, sprite sheets, text styles, transitions, menus, music and sounds, every screen's parts, a complete example and the validation rules.

[Unreleased]: https://github.com/openkakutou/screenpack/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/openkakutou/screenpack/releases/tag/v0.1.0
