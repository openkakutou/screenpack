# screenpack

OpenKakutou `screenpack` holds the OpenKakutou-specific screenpack format (roadmap decision `029`): its specification, the JSON Schemas (root `screenpack.json` + one per screen), the default screenpack, and a Go library that loads, validates (schemas + cross-file checks), resolves per-screen fallbacks to the default screenpack and writes screenpacks. It compiles to WASM for `game`'s web build and is imported natively by `game`'s native build and by the planned `screenpack-editor`.

Part of [OpenKakutou](https://github.com/openkakutou). Early stage: see `.vibe/backlog/` for the planned work.

## Go library

```go
sp, err := screenpack.Load(os.DirFS("screenpacks/neon"))
var problems screenpack.Problems
if errors.As(err, &problems) {
	for _, p := range problems {
		fmt.Println(p.File, p.Path, p.Reason) // e.g. screens/pause.json /title/style unknown text style "huge"
	}
}
```

`sp.Screens` holds the screens the pack has; `sp.Warnings` lists what was ignored, such as an unknown file in `screens/`. `screenpack.Resolve(sp, def)` then fills every missing screen, and every part a pack screen leaves out, from the default screenpack `def`; `Source(screen, part)` tells which pack each one comes from, so its text styles, fonts and assets are looked up there.

## License

MIT, see [LICENSE](LICENSE).

## Documentation

- [Screenpack format](docs/format.md) — the specification, format version 1.
- [JSON Schemas](schema/v1/) — one per file kind, for editors and validation.
- [Default screenpack](default/) — every screen, the fallback for the others; [asset licenses](default/ASSETS.md).
