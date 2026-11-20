package screenpack

import (
	"embed"
	"io/fs"
)

//go:embed all:default
var defaultFiles embed.FS

// DefaultFS returns the default screenpack's files (the default/ folder),
// built into the program: a game ships it and falls back on it without
// depending on a folder on disk.
func DefaultFS() fs.FS {
	sub, err := fs.Sub(defaultFiles, "default")
	if err != nil {
		panic(err) // "default" is a valid literal path: fs.Sub cannot fail on it
	}
	return sub
}
