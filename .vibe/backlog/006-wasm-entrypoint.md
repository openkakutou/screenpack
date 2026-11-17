---
status: todo
depends_on: [003, 004]
---
# WASM entrypoint with a release smoke gate

## Description
Roadmap decision `029`: a `cmd/wasm` entrypoint exposing load + resolve to JavaScript (same contract style as `sff`/`character`: byte buffers in, JSON out), and `cmd/wasm/smoke.mjs` run by the release workflow on the built module.

## Acceptance Criteria
- [ ] `game`'s web build can load and resolve a screenpack through the WASM module
- [ ] The release workflow fails if the smoke test fails
