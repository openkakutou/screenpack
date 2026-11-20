---
status: todo
---
# Format gaps found by the native game (game item 030)

## Description
Drawing the game's screens from the default pack (game ADR 025) hit points the format does not settle:
- `pressStart.blinkMs`: half-period or full period? The game reads it as the half-period (visible `blinkMs`, hidden `blinkMs`).
- No warning/error style: flagged rows (conflicting bindings) and notices use a fixed red in the game. Add an optional `warningStyle` (root or menu).
- Rows with two values (controls: key and button) have one `valueOffset`; add `valueOffsets` (one per column) or a column layout in the `controls` part.
- No error screen (content missing, unsupported mode): define one with a wrapping text box part.
- A single text's `align` against its own box has no effect: say what `align` means for a lone text element (relative to the anchor point?).
- `tabs` on the options screen: the game's options are pages (list, then one screen per page), not tabs; say whether `tabs` is required to be drawn.

## Acceptance Criteria
- [ ] Each point above is settled in `docs/format.md` and the schemas, and the default pack uses the result
