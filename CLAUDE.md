# CLAUDE.md — screenpack

> Edit freely. Add `<!-- keep -->` on a section heading to preserve it on regeneration by /vibe:init.

## Project overview

OpenKakutou `screenpack` holds the OpenKakutou-specific screenpack format (roadmap decision `029`): its specification, the JSON Schemas (root `screenpack.json` + one per screen), the default screenpack, and a Go library that loads, validates (schemas + cross-file checks), resolves per-screen fallbacks to the default screenpack and writes screenpacks. It compiles to WASM for `game`'s web build and is imported natively by `game`'s native build and by the planned `screenpack-editor`.

**Stack:** Go 1.26, compiled to WASM (`GOOS=js GOARCH=wasm`) for web consumers, imported as a Go module by native ones.
**Type:** library

## Project language

English — all documentation, backlog items, code comments, and other generated content for this project must be written in English.

## Compatibility scope <!-- keep -->

MUGEN/Ikemen GO screenpacks (`system.def`, `select.def`) are **out of** the org's compatibility MUST: they are not read nor converted (roadmap decision `029`). The format here is OpenKakutou's own: JSONC validated by published JSON Schemas, anchor + offset layout in relative units, standard assets only (PNG/WebP, OGG/Opus, TTF/OTF/WOFF2). Lifebars are not part of a screenpack (roadmap decision `031`). A screenpack describes appearance only, never behavior (no scripting, no new menu actions; roadmap decisions `029`, `030`). The org-wide MUST still applies to combat content, which this repo never touches; for that content Ikemen GO is the verification reference (roadmap decision `038`).

## Release propagation <!-- keep -->

Consumers pin an exact tag (WASM download or `go.mod`), propagated by this repo's own release step, as in `sff`/`character` (roadmap decisions `016`, `025`). After `/vibe:release`, bump every consumer listed in `roadmap/repos.md`, run its tests, commit and push; stop and flag on red.

## Development workflow (Vibe Coding)

The user is the **Product Owner only**. They describe requirements and evaluate outcomes — they do NOT write or manually test code.

### Definition of Done

A task is complete ONLY when ALL of the following are true:

- [ ] Implementation code is written
- [ ] Tests covering the nominal path exist and pass
- [ ] Tests covering edge cases and error paths exist and pass
- [ ] Lint exits with code 0 — `gofmt -l .` reports no files, `go vet ./...` is clean
- [ ] All tests pass — `go test ./...`
- [ ] No debug artifacts left in code

**Never present a result to the user if tests are failing.**

### TDD workflow

1. Write tests first, 2. confirm they fail, 3. write the minimum implementation, 4. confirm they pass, 5. refactor, 6. run `gofmt -l .` and `go vet ./...`, 7. present the result with the test output.

### Self-correction loop

If tests or lint fail: diagnose, fix, re-run, up to 3 attempts; escalate only on an ambiguous requirement.

## Testing conventions

- Tests co-located with source, `*_test.go`, one test file per source file, built-in `go test`
- Test names describe behavior: `TestLoad_ReturnsError_OnMissingFile`
- Always cover: happy path + at least 2 edge cases + error/invalid input
- Round-trip tests for any write-path code

## Constraints

- Never leave dead code or unused imports
- No rendering dependency; must compile cleanly to WASM (`GOOS=js GOARCH=wasm`)
- Read and write paths separated: write-only concerns never leak into read-path types

## Review agents

| Agent | Active | Reason |
|---|---|---|
| `vibe:review-tests` | ✅ | always active |
| `vibe:review-naming` | ✅ | always active |
| `vibe:review-security` | ✅ | always active |
| `vibe:review-dependencies` | ✅ | always active |
| `vibe:review-robustness` | ✅ | always active |
| `vibe:review-hygiene` | ✅ | always active |
| `vibe:review-antipatterns` | ✅ | always active |
| `vibe:review-simplicity` | ✅ | always active |
| `vibe:review-overengineering` | ✅ | always active |
| `vibe:review-correctness` | ✅ | always active |
| `vibe:review-comments` | ✅ | always active |
| `vibe:review-signatures` | ✅ | always active |
| `vibe:review-coordination` | ✅ | always active |
| `vibe:review-docs` | ✅ | active: Markdown documentation beyond the README (`.vibe/`, ADRs) |
| `vibe:review-ddd` | ❌ | inactive: no explicit domain layer |
| `vibe:review-performance` | ❌ | inactive: library |
| `vibe:review-web-security` | ❌ | inactive: no HTTP surface |
| `vibe:review-hexagonal` | ❌ | inactive: no declared hexagonal architecture |
