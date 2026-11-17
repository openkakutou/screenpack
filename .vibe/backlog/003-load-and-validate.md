---
status: todo
depends_on: [002]
---
# Load and validate a screenpack

## Description
Roadmap decision `029`: a Go `Load` reading a screenpack folder (or an in-memory file set for WASM), parsing JSONC, validating each file against its schema, then cross-file checks: referenced assets exist, screen ids are known, fonts and sounds resolve.

## Acceptance Criteria
- [ ] A valid screenpack loads into a typed model
- [ ] Every error names the file, the JSON path and the reason; all errors are reported, not just the first
- [ ] Unknown screen files are reported, not fatal
