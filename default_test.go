package screenpack

import (
	"io/fs"
	"os"
	"reflect"
	"strings"
	"testing"
)

const defaultDir = "default"

func TestDefault_LoadsWithoutProblemsAndCoversEveryScreen(t *testing.T) {
	sp, err := Load(os.DirFS(defaultDir))
	if err != nil {
		t.Fatal(err)
	}
	if len(sp.Warnings) != 0 {
		t.Errorf("warnings: %v", sp.Warnings)
	}
	screens := reflect.ValueOf(sp.Screens)
	for i := range screens.NumField() {
		if screens.Field(i).IsNil() {
			t.Errorf("default screenpack has no %s screen", screens.Type().Field(i).Name)
		}
	}
}

// Every asset of the default screenpack has its license recorded in
// ASSETS.md, and every file ASSETS.md names exists.
func TestDefault_RecordsTheLicenseOfEveryAsset(t *testing.T) {
	record, err := os.ReadFile(defaultDir + "/ASSETS.md")
	if err != nil {
		t.Fatal(err)
	}
	listed := map[string]bool{}
	for _, line := range strings.Split(string(record), "\n") {
		cells := strings.Split(line, "|")
		if len(cells) < 4 || !strings.HasPrefix(strings.TrimSpace(cells[1]), "`") {
			continue
		}
		name := strings.Trim(strings.TrimSpace(cells[1]), "`")
		if strings.TrimSpace(cells[2]) == "" {
			t.Errorf("%s: no license", name)
		}
		listed[name] = true
	}
	assets := fs.FS(os.DirFS(defaultDir))
	err = fs.WalkDir(assets, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !isAsset(p) {
			return err
		}
		if !listed[p] {
			t.Errorf("%s: license not recorded in ASSETS.md", p)
		}
		delete(listed, p)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for name := range listed {
		t.Errorf("ASSETS.md names %s, which does not exist", name)
	}
}

func isAsset(p string) bool {
	for _, dir := range []string{"images/", "sounds/", "fonts/"} {
		if strings.HasPrefix(p, dir) && !strings.HasSuffix(p, ".json") && !strings.HasSuffix(p, ".txt") {
			return true
		}
	}
	return false
}
