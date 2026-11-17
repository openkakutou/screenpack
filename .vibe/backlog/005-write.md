---
status: todo
depends_on: [003]
---
# Write a screenpack back to disk

## Description
Roadmap decision `029`: serialize the model back to `screenpack.json` + `screens/*.json` for `screenpack-editor`, keeping unchanged files byte-exact (comments included) and rewriting only what changed.

## Acceptance Criteria
- [ ] Round trip of an unchanged screenpack is byte-identical
- [ ] An edit rewrites only the affected file and lines
- [ ] A new screenpack is written in a valid, formatted layout
