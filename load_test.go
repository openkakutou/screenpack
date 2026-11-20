package screenpack

import (
	"errors"
	"strings"
	"testing"
	"testing/fstest"
)

const root = `{
	"formatVersion": 1, // JSONC
	"name": "Neon", "author": "Jane Doe", "version": "1.0.0",
	"fonts": { "body": "fonts/Inter.woff2" },
	"textStyles": {
		"default": { "font": "body", "size": 4, "color": "#e8e8ff" },
		"active": { "font": "body", "size": 5, "color": "#ffffff", "align": "center" },
	},
	"sounds": { "move": "sounds/move.ogg", "confirm": "sounds/confirm.ogg",
	            "cancel": "sounds/cancel.ogg", "error": "sounds/error.ogg" }
}`

const cursorSheet = `{
	"image": "cursor.png",
	"frames": { "c0": [0, 0, 32, 32], "c1": [32, 0, 32, 32] },
	"animations": { "blink": { "frames": [{ "frame": "c0", "ms": 300 }, { "frame": "c1", "ms": 300 }] } }
}`

func file(s string) *fstest.MapFile { return &fstest.MapFile{Data: []byte(s)} }

// valid is a small screenpack that uses every kind of cross-file reference.
func valid() fstest.MapFS {
	return fstest.MapFS{
		"screenpack.json":    file(root),
		"fonts/Inter.woff2":  file(""),
		"sounds/move.ogg":    file(""),
		"sounds/confirm.ogg": file(""),
		"sounds/cancel.ogg":  file(""),
		"sounds/error.ogg":   file(""),
		"sounds/title.ogg":   file(""),
		"images/cursor.json": file(cursorSheet),
		"images/cursor.png":  file(""),
		"images/logo.webp":   file(""),
		"screens/title.json": file(`{
			"formatVersion": 1,
			"music": { "file": "sounds/title.ogg", "loopStartMs": 4200 },
			"textStyles": { "big": { "font": "body", "size": 9, "color": "#ff3cac" } },
			"logo": { "type": "image", "image": "images/logo.webp", "anchor": "top", "size": [90, null] },
			"menu": { "style": "big", "activeStyle": "active", "order": ["versus", "options"],
			          "cursor": { "type": "image", "sheet": "images/cursor.json", "animation": "blink" } }
		}`),
		"screens/pause.json": file(`{ "formatVersion": 1, "music": null,
			"title": { "style": "default", "offset": [0, -20] } }`),
	}
}

func TestLoad_ReturnsTheTypedModel(t *testing.T) {
	sp, err := Load(valid())
	if err != nil {
		t.Fatal(err)
	}
	if sp.Name != "Neon" || sp.Fonts["body"] != "fonts/Inter.woff2" || sp.TextStyles["active"].Align != "center" {
		t.Errorf("root = %+v", sp)
	}
	title := sp.Screens.Title
	if title == nil || title.Music == nil || title.Music.LoopStartMs != 4200 || title.Silent {
		t.Fatalf("title = %+v", title)
	}
	if title.Logo.Image != "images/logo.webp" || title.Logo.Size[1] != nil || *title.Logo.Size[0] != 90 {
		t.Errorf("logo = %+v", title.Logo)
	}
	if strings.Join(title.Menu.Order, ",") != "versus,options" || title.Menu.Cursor.Animation != "blink" {
		t.Errorf("menu = %+v", title.Menu)
	}
	if title.TextStyles["big"].Size != 9 {
		t.Errorf("screen text styles = %+v", title.TextStyles)
	}
	if sp.Screens.Pause == nil || !sp.Screens.Pause.Silent || sp.Screens.Pause.Title.Style != "default" {
		t.Errorf("pause = %+v", sp.Screens.Pause)
	}
	if sp.Screens.Versus != nil || sp.Screens.CharacterSelect != nil {
		t.Error("screens the pack does not have must stay nil (fallback to default)")
	}
	sheet := sp.Sheets["images/cursor.json"]
	if sheet == nil || sheet.Frames["c1"] != [4]int{32, 0, 32, 32} || len(sheet.Animations["blink"].Frames) != 2 {
		t.Errorf("sheet = %+v", sheet)
	}
	if len(sp.Warnings) != 0 {
		t.Errorf("warnings = %v", sp.Warnings)
	}
}

func TestLoad_RootFileAlone(t *testing.T) {
	fsys := valid()
	for name := range fsys {
		if strings.HasPrefix(name, "screens/") || strings.HasPrefix(name, "images/") {
			delete(fsys, name)
		}
	}
	sp, err := Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	if sp.Screens.Title != nil || len(sp.Sheets) != 0 {
		t.Errorf("screens = %+v sheets = %v", sp.Screens, sp.Sheets)
	}
}

func TestLoad_ReportsUnknownScreenFilesWithoutFailing(t *testing.T) {
	fsys := valid()
	fsys["screens/credits.json"] = file(`{}`)
	fsys["screens/notes.txt"] = file("x")
	sp, err := Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	if len(sp.Warnings) != 2 || sp.Warnings[0].File != "screens/credits.json" || sp.Warnings[1].File != "screens/notes.txt" {
		t.Fatalf("warnings = %v", sp.Warnings)
	}
	if !strings.Contains(sp.Warnings[0].Reason, "unknown screen") {
		t.Errorf("reason = %q", sp.Warnings[0].Reason)
	}
}

// problems loads fsys, which must fail, and returns its problems keyed by
// "file path".
func problems(t *testing.T, fsys fstest.MapFS) map[string]string {
	t.Helper()
	sp, err := Load(fsys)
	if err == nil {
		t.Fatalf("loaded %+v, want problems", sp)
	}
	var ps Problems
	if !errors.As(err, &ps) {
		t.Fatalf("got %T %v, want Problems", err, err)
	}
	out := map[string]string{}
	for _, p := range ps {
		out[p.File+" "+p.Path] = p.Reason
	}
	return out
}

func TestLoad_ReportsEveryProblemAcrossFiles(t *testing.T) {
	fsys := valid()
	delete(fsys, "fonts/Inter.woff2")
	delete(fsys, "images/cursor.png")
	fsys["screens/pause.json"] = file(`{ "formatVersion": 1, "title": { "style": "huge" } }`)
	fsys["screens/versus.json"] = file(`{ "formatVersion": 1, "durationMs": -1 }`)
	fsys["screens/title.json"] = file(`{ "formatVersion": 1,
		"textStyles": { "big": { "font": "serif", "size": 9, "color": "#ffffff" } },
		"menu": { "cursor": { "type": "image", "sheet": "images/cursor.json", "animation": "spin" } } }`)
	got := problems(t, fsys)
	want := map[string]string{
		"screenpack.json /fonts/body":               "missing file fonts/Inter.woff2",
		"images/cursor.json /image":                 "missing file images/cursor.png",
		"screens/pause.json /title/style":           `unknown text style "huge"`,
		"screens/versus.json /durationMs":           "",
		"screens/title.json /textStyles/big/font":   `unknown font "serif"`,
		"screens/title.json /menu/cursor/animation": `unknown animation "spin" in images/cursor.json`,
	}
	for key, reason := range want {
		r, ok := got[key]
		if !ok {
			t.Errorf("no problem at %s; got %v", key, got)
			continue
		}
		if reason != "" && r != reason {
			t.Errorf("%s: %q, want %q", key, r, reason)
		}
	}
	if len(got) != len(want) {
		t.Errorf("got %d problems, want %d: %v", len(got), len(want), got)
	}
}

func TestLoad_ReportsBrokenFiles(t *testing.T) {
	fsys := valid()
	fsys["screens/pause.json"] = file(`{ "formatVersion": 1,`)
	fsys["images/cursor.json"] = file(`{ "image": "cursor.png", "frames": {} }`)
	got := problems(t, fsys)
	if _, ok := got["screens/pause.json "]; !ok {
		t.Errorf("JSONC syntax error not reported against the file: %v", got)
	}
	if _, ok := got["images/cursor.json /frames"]; !ok {
		t.Errorf("sprite sheet schema error not reported: %v", got)
	}
	if _, ok := got["screens/title.json /menu/cursor/animation"]; ok {
		t.Errorf("a broken sheet must not also report every reference into it: %v", got)
	}
}

func TestLoad_FailsWithoutRootFile(t *testing.T) {
	fsys := valid()
	delete(fsys, "screenpack.json")
	got := problems(t, fsys)
	if !strings.Contains(got["screenpack.json "], "missing") {
		t.Errorf("problems = %v", got)
	}
}

func TestLoad_FrameAndButtonIconReferences(t *testing.T) {
	fsys := valid()
	fsys["screens/moveList.json"] = file(`{ "formatVersion": 1,
		"buttons": { "sheet": "images/cursor.json", "frames": { "a": "c0", "b": "c9" } },
		"scrollbar": { "thumb": { "type": "image", "sheet": "images/cursor.json", "frame": "c7" } } }`)
	got := problems(t, fsys)
	if got["screens/moveList.json /buttons/frames/b"] != `unknown frame "c9" in images/cursor.json` {
		t.Errorf("button icon: %v", got)
	}
	if got["screens/moveList.json /scrollbar/thumb/frame"] != `unknown frame "c7" in images/cursor.json` {
		t.Errorf("element frame: %v", got)
	}
	if len(got) != 2 {
		t.Errorf("problems = %v", got)
	}
}

func TestLoad_AnimationFramesMustExistInTheirSheet(t *testing.T) {
	fsys := valid()
	fsys["images/cursor.json"] = file(`{ "image": "cursor.png", "frames": { "c0": [0, 0, 1, 1] },
		"animations": { "blink": { "frames": [{ "frame": "c0", "ms": 1 }, { "frame": "cX", "ms": 1 }] } } }`)
	got := problems(t, fsys)
	if got["images/cursor.json /animations/blink/frames/1/frame"] != `unknown frame "cX"` {
		t.Errorf("problems = %v", got)
	}
}

func TestProblems_ErrorListsEachOne(t *testing.T) {
	err := Problems{{File: "a.json", Path: "/x", Reason: "bad"}, {File: "b.json", Reason: "missing"}}
	if err.Error() != "a.json: /x: bad\nb.json: missing" {
		t.Errorf("Error() = %q", err.Error())
	}
}

func TestLoad_ReadsWarningColorAndChecksFlaggedStyle(t *testing.T) {
	fsys := valid()
	fsys["screenpack.json"] = file(strings.Replace(root, `"sounds"`, `"warningColor": "#aa0000", "sounds"`, 1))
	fsys["screens/pause.json"] = file(`{ "formatVersion": 1, "menu": { "flaggedStyle": "active" } }`)
	fsys["screens/error.json"] = file(`{ "formatVersion": 1, "message": { "style": "default", "width": 110 } }`)
	sp, err := Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	if sp.WarningColor != "#aa0000" {
		t.Errorf("warning color = %q", sp.WarningColor)
	}
	if sp.Screens.Pause.Menu.FlaggedStyle != "active" {
		t.Errorf("flagged style = %q", sp.Screens.Pause.Menu.FlaggedStyle)
	}
	if sp.Screens.Error == nil || sp.Screens.Error.Message.Width != 110 {
		t.Errorf("error screen = %+v", sp.Screens.Error)
	}
}

func TestLoad_ReportsUnknownFlaggedStyle(t *testing.T) {
	fsys := valid()
	fsys["screens/pause.json"] = file(`{ "formatVersion": 1, "menu": { "flaggedStyle": "nope" } }`)
	got := problems(t, fsys)
	if got["screens/pause.json /menu/flaggedStyle"] == "" {
		t.Errorf("problems = %v", got)
	}
}
