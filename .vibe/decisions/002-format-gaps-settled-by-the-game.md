---
date: 2026-11-20
status: accepted
---
# Format gaps settled from the native game's first use

**Context:** Backlog 009: drawing every screen of the default pack in the native game (game ADR 025) hit six points the format left open.

**Decision:**
1. `pressStart.blinkMs` is the half-period: the text shows for `blinkMs`, then hides for `blinkMs`.
2. Menus take an optional `flaggedStyle` (style id) for an entry the game flags as a problem (e.g. two actions on one key); without it the game uses `activeStyle`/`style` tinted with the root `warningColor`. The root file takes an optional `warningColor` (default `#dc3c3c`).
3. Menus take an optional `valueOffsets`: one `[x, y]` per value column, for entries with several values (key and button on the controls page); without it, a single `valueOffset` places the last value and the others are spread evenly between the label and it.
4. A new optional screen `error`: `background`, `transition`, `music`, `title` (text part), `message` (text part plus `width` in units: the text wraps inside that width), `hint` (text part). Missing, it resolves from the default pack like every screen.
5. `align` aligns lines and entries with each other; a lone text element is placed by its pivot only (its box is its own width), so `align` has no effect on it.
6. `tabs` on the options screen is optional for a game: a game that shows options as pages (a list, then one page per item) uses `menu` for both and does not draw `tabs`.

**Reason:** Each choice keeps existing packs valid (all additions are optional) and matches what the game already draws, so the default pack can rely on it.

**Rejected alternatives:** A root `styles.warning` text style (a colour is all the game needs, and a whole style duplicates font and size); a `columns` table on the controls part (menus already own value placement); making `error` mandatory (would invalidate every existing pack).
