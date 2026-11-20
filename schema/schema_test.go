package schema

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// specExample is one ```jsonc block of docs/format.md and the document kind
// it illustrates.
type specExample struct {
	kind string
	line int
	text string
}

var (
	rootLabel   = regexp.MustCompile("^`screenpack\\.json`:$")
	screenLabel = regexp.MustCompile("^`screens/(\\w+)\\.json`:$")
	screenTitle = regexp.MustCompile("^### `(\\w+)\\.json`")
)

// specExamples reads every JSONC example of the spec whose kind is clear:
// labelled as a whole file, or the sprite-sheet example.
func specExamples(t *testing.T) []specExample {
	t.Helper()
	f, err := os.Open("../docs/format.md")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var out []specExample
	var heading, label string
	var block *specExample
	n := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		n++
		line := sc.Text()
		if block != nil {
			if line == "```" {
				out = append(out, *block)
				block = nil
				continue
			}
			block.text += line + "\n"
			continue
		}
		switch {
		case strings.HasPrefix(line, "#"):
			heading = line
		case line == "```jsonc":
			kind := ""
			if m := screenLabel.FindStringSubmatch(label); m != nil {
				kind = m[1]
			}
			if rootLabel.MatchString(label) {
				kind = "screenpack"
			}
			if heading == "### Sprite sheets" {
				kind = "spriteSheet"
			}
			if kind == "" {
				t.Fatalf("docs/format.md:%d: a JSONC example with no file label", n)
			}
			block = &specExample{kind: kind, line: n}
		case strings.TrimSpace(line) != "":
			label = line
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestValidate_AcceptsEverySpecExample(t *testing.T) {
	examples := specExamples(t)
	if len(examples) < 5 {
		t.Fatalf("found %d examples in the spec, want at least 5", len(examples))
	}
	for _, ex := range examples {
		if err := Validate(ex.kind, []byte(ex.text)); err != nil {
			t.Errorf("docs/format.md:%d (%s): %v", ex.line, ex.kind, err)
		}
	}
}

func TestValidate_RejectsWithThePath(t *testing.T) {
	cases := []struct {
		name, kind, doc, path string
	}{
		{"wrong type", "pause", `{"formatVersion": 1, "dim": {"color": 3}}`, "/dim/color"},
		{"unknown property", "pause", `{"formatVersion": 1, "menu": {"spacing": 2, "colour": "#fff"}}`, "/menu"},
		{"missing required field", "screenpack", `{"formatVersion": 1, "author": "a", "version": "1",
			"fonts": {}, "textStyles": {"default": {"font": "f", "size": 1, "color": "#ffffff"}},
			"sounds": {"move": "m.ogg", "confirm": "c.ogg", "cancel": "x.ogg", "error": "e.ogg"}}`, "name"},
		{"wrong format version", "title", `{"formatVersion": 2}`, "/formatVersion"},
		{"path leaving the folder", "title", `{"formatVersion": 1, "music": {"file": "../x.ogg"}}`, "/music/file"},
		{"asset in a MUGEN format", "title", `{"formatVersion": 1, "logo": {"type": "image", "image": "logo.sff"}}`, "/logo"},
		{"element without its type's field", "pause", `{"formatVersion": 1, "foreground": [{"type": "rect"}]}`, "/foreground/0"},
		{"bad warning color", "screenpack", `{"formatVersion": 1, "name": "n", "author": "a", "version": "1",
			"fonts": {}, "textStyles": {"default": {"font": "f", "size": 1, "color": "#ffffff"}},
			"sounds": {"move": "m.ogg", "confirm": "c.ogg", "cancel": "x.ogg", "error": "e.ogg"},
			"warningColor": "red"}`, "/warningColor"},
		{"value offsets of a wrong shape", "options", `{"formatVersion": 1, "menu": {"valueOffsets": [[30]]}}`, "/menu/valueOffsets/0"},
		{"empty value offsets", "options", `{"formatVersion": 1, "menu": {"valueOffsets": []}}`, "/menu/valueOffsets"},
		{"flagged style that is not an id", "pause", `{"formatVersion": 1, "menu": {"flaggedStyle": 3}}`, "/menu/flaggedStyle"},
		{"error message without a positive width", "error", `{"formatVersion": 1, "message": {"width": 0}}`, "/message/width"},
		{"unknown part on the error screen", "error", `{"formatVersion": 1, "menu": {}}`, "/menu"},
		{"sound key of another screen", "pause", `{"formatVersion": 1, "sounds": {"tick": "t.ogg"}}`, "/sounds"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := Validate(c.kind, []byte(c.doc))
			if err == nil {
				t.Fatal("accepted, want an error")
			}
			if !strings.Contains(err.Error(), c.path) {
				t.Errorf("error does not name %q:\n%v", c.path, err)
			}
		})
	}
}

func TestValidate_ReadsJSONC(t *testing.T) {
	doc := `{
		// a comment
		"formatVersion": 1, /* another */
		"music": null,
	}`
	if err := Validate("pause", []byte(doc)); err != nil {
		t.Fatal(err)
	}
	if err := Validate("pause", []byte(`{"formatVersion": 1`)); err == nil {
		t.Error("accepted truncated JSONC")
	}
}

func TestValidate_RefusesUnknownKind(t *testing.T) {
	err := Validate("lifebar", []byte(`{}`))
	if err == nil || !strings.Contains(err.Error(), "lifebar") {
		t.Fatalf("got %v, want an error naming the kind", err)
	}
}

// The spec describes one screen per "### `<name>.json`" heading; there must
// be exactly one schema per such screen.
func TestScreens_MatchTheSpec(t *testing.T) {
	spec, err := os.ReadFile("../docs/format.md")
	if err != nil {
		t.Fatal(err)
	}
	var want []string
	for _, line := range strings.Split(string(spec), "\n") {
		if m := screenTitle.FindStringSubmatch(line); m != nil {
			want = append(want, m[1])
		}
	}
	sort.Strings(want)
	got := Screens()
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("schemas %v, spec %v", got, want)
	}
}

func TestSchemas_CarryTheFormatVersion(t *testing.T) {
	err := fs.WalkDir(files, "v1", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := files.ReadFile(path)
		if err != nil {
			return err
		}
		var s struct {
			ID string `json:"$id"`
		}
		if err := json.Unmarshal(data, &s); err != nil {
			t.Errorf("%s: %v", path, err)
			return nil
		}
		if want := BaseURL + strings.TrimPrefix(path, "v1/"); s.ID != want {
			t.Errorf("%s: $id %q, want %q", path, s.ID, want)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, kind := range append(Screens(), "screenpack") {
		err := Validate(kind, []byte(`{"formatVersion": 2}`))
		if err == nil || !strings.Contains(err.Error(), "/formatVersion") {
			t.Errorf("%s: formatVersion 2 not refused: %v", kind, err)
		}
	}
}

func TestValidate_ListsEveryProblemWithItsPathAndReason(t *testing.T) {
	err := Validate("pause", []byte(`{"formatVersion": 1, "dim": {"color": 3}, "title": {"colour": "#fff"}}`))
	var invalid Invalid
	if !errors.As(err, &invalid) {
		t.Fatalf("got %T %v, want Invalid", err, err)
	}
	got := map[string]string{}
	for _, issue := range invalid {
		got[issue.Path] = issue.Reason
	}
	if got["/dim/color"] == "" || got["/title/colour"] == "" || len(got) != 2 {
		t.Fatalf("issues = %+v, want one at /dim/color and one at /title/colour", invalid)
	}
	if got["/title/colour"] != "not allowed here" {
		t.Errorf("unknown property reason = %q", got["/title/colour"])
	}
}

func TestValidate_PointsAtTheRootForAMissingField(t *testing.T) {
	var invalid Invalid
	if !errors.As(Validate("title", []byte(`{}`)), &invalid) || len(invalid) != 1 || invalid[0].Path != "/" {
		t.Fatalf("issues = %+v", invalid)
	}
}

func TestValidate_AcceptsTheNewOptionalFields(t *testing.T) {
	docs := map[string]string{
		"options": `{"formatVersion": 1, "menu": {"valueOffsets": [[30, 0], [50, 0]], "flaggedStyle": "bad"}}`,
		"error": `{"formatVersion": 1, "music": null, "title": {"style": "heading"},
			"message": {"style": "default", "width": 110}, "hint": {"style": "hint"}}`,
	}
	for kind, doc := range docs {
		if err := Validate(kind, []byte(doc)); err != nil {
			t.Errorf("%s: %v", kind, err)
		}
	}
}
