---
status: todo
---
# Write the screenpack format specification

## Description
Roadmap decision `029`: write `docs/format.md`: root `screenpack.json` (name, author, version, format version, defaults: base fonts and sounds), one file per screen under `screens/`, JSONC syntax, anchor + offset layout in relative units, standard assets (PNG/WebP with JSON sprite-sheet descriptions, OGG/Opus, TTF/OTF/WOFF2), per-screen fallback to the default screenpack, what a screenpack may change (visuals, layout, element animations and transitions, text styles, sounds, order and visibility of the menu entries the game offers) and may not (behavior). Screens: those of `game`'s built-in screen types (game#032): title/main menu, character select and stage select (with a sort rule: folder name, display name or author, case-insensitive, default folder name; roadmap decision `034`), match setup, VS, options (gameplay, audio, video, controls, appearance), pause, move list, results, continue, victory. Character portraits and names come from the characters, drawn by the game. No lifebar (roadmap decision `031`).

## Acceptance Criteria
- [ ] The spec covers every point above with one complete example screenpack
- [ ] Each screen lists its elements and their parameters
- [ ] Links roadmap decisions `029`, `030`, `031`, `034`
