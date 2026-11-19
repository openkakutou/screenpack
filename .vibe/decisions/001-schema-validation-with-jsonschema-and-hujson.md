---
date: 2026-11-19
status: accepted
---
# Schemas validated by santhosh-tekuri/jsonschema, JSONC read by tailscale/hujson, published from the repo

**Context:** Backlog item 002 publishes one JSON Schema per screenpack file kind; item 003's loader must validate every file against them, natively and in WASM.
**Decision:** The schemas live in `schema/v1/` (draft 2020-12, shared definitions in `common.schema.json`), embedded in the `schema` Go package, which validates with `github.com/santhosh-tekuri/jsonschema/v6` after turning JSONC into JSON with `github.com/tailscale/hujson`. Their `$id`, the URL editors use, is `https://raw.githubusercontent.com/openkakutou/screenpack/main/schema/v1/<kind>.schema.json`; a breaking change goes to `v2/`, so `v1/` stays stable.
**Reason:** The schemas are the format's single source of truth: validating with them, instead of re-stating every rule in Go as `game`'s mode loader does, keeps the spec, the schemas and the loader from drifting. Both libraries are pure Go and build for WASM.
**Rejected alternatives:** Hand-written Go validation (duplicates the schemas); a GitHub Pages site for the schema URLs (a deployment to keep for static files the repo already serves raw).
