# screenpack

OpenKakutou `screenpack` holds the OpenKakutou-specific screenpack format (roadmap decision `029`): its specification, the JSON Schemas (root `screenpack.json` + one per screen), the default screenpack, and a Go library that loads, validates (schemas + cross-file checks), resolves per-screen fallbacks to the default screenpack and writes screenpacks. It compiles to WASM for `game`'s web build and is imported natively by `game`'s native build and by the planned `screenpack-editor`.

Part of [OpenKakutou](https://github.com/openkakutou). Early stage: see `.vibe/backlog/` for the planned work.

## License

MIT, see [LICENSE](LICENSE).
