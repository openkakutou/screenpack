package screenpack

import (
	"io/fs"
	"os"
	"testing"
)

func TestDefaultFS_HoldsTheSameFilesAsTheDefaultFolder(t *testing.T) {
	var onDisk, embedded []string
	walk := func(fsys fs.FS, out *[]string) {
		err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				*out = append(*out, p)
			}
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	walk(os.DirFS(defaultDir), &onDisk)
	walk(DefaultFS(), &embedded)
	if len(onDisk) == 0 || len(onDisk) != len(embedded) {
		t.Fatalf("embedded %d files, %s has %d", len(embedded), defaultDir, len(onDisk))
	}
	for i := range onDisk {
		if onDisk[i] != embedded[i] {
			t.Errorf("file %d: embedded %q, on disk %q", i, embedded[i], onDisk[i])
		}
	}
}

func TestDefaultFS_LoadsWithoutProblems(t *testing.T) {
	sp, err := Load(DefaultFS())
	if err != nil {
		t.Fatal(err)
	}
	if len(sp.Warnings) != 0 {
		t.Errorf("warnings: %v", sp.Warnings)
	}
}
