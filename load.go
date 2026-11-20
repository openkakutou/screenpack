package screenpack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"path"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/openkakutou/screenpack/schema"
	"github.com/tailscale/hujson"
)

// RootFile is the file every screenpack has.
const RootFile = "screenpack.json"

const screensDir = "screens"

// Problem is one thing wrong in a screenpack: the file, the JSON Pointer of
// the faulty value ("" for the whole file) and the reason.
type Problem struct {
	File   string
	Path   string
	Reason string
}

// Problems lists every problem found; as an error, one per line.
type Problems []Problem

func (ps Problems) Error() string {
	lines := make([]string, len(ps))
	for i, p := range ps {
		lines[i] = p.File + ": " + p.Reason
		if p.Path != "" {
			lines[i] = p.File + ": " + p.Path + ": " + p.Reason
		}
	}
	return strings.Join(lines, "\n")
}

// Load reads the screenpack at the root of fsys (os.DirFS of its folder,
// or an in-memory file set). It checks every file against its schema, then
// that every referenced file, font, text style, sprite-sheet frame and
// animation exists. On failure the error is Problems, listing all of them.
func Load(fsys fs.FS) (*Screenpack, error) {
	l := &loader{fsys: fsys, sheets: map[Path]*SpriteSheet{}, tried: map[Path]bool{}}
	sp := &Screenpack{}
	var root rootFile
	l.rootOK = l.read(RootFile, "screenpack", &root)
	if l.rootOK {
		sp.Name, sp.Author, sp.Version, sp.WarningColor = root.Name, root.Author, root.Version, root.WarningColor
		sp.Fonts, sp.TextStyles, sp.Sounds, sp.Music = root.Fonts, root.TextStyles, root.Sounds, root.Music
		l.fonts, l.rootStyles = root.Fonts, root.TextStyles
		l.check(RootFile, &root, nil)
	}
	sp.Warnings = l.loadScreens(&sp.Screens)
	if len(l.problems) > 0 {
		return nil, l.problems
	}
	sp.Sheets = l.sheets
	return sp, nil
}

type loader struct {
	fsys     fs.FS
	problems Problems
	// rootOK is false when screenpack.json is broken: fonts and root text
	// styles are then unknown, so references to them are not checked.
	rootOK     bool
	fonts      map[FontID]Path
	rootStyles map[StyleID]TextStyle
	sheets     map[Path]*SpriteSheet
	tried      map[Path]bool
}

func (l *loader) add(file, at, reason string) {
	l.problems = append(l.problems, Problem{File: file, Path: at, Reason: reason})
}

// read validates file against the schema of kind and decodes it into v.
func (l *loader) read(file, kind string, v any) bool {
	data, err := fs.ReadFile(l.fsys, file)
	if errors.Is(err, fs.ErrNotExist) {
		l.add(file, "", "missing file")
		return false
	}
	if err != nil {
		l.add(file, "", err.Error())
		return false
	}
	err = schema.Validate(kind, data)
	var invalid schema.Invalid
	if errors.As(err, &invalid) {
		for _, issue := range invalid {
			l.add(file, issue.Path, issue.Reason)
		}
		return false
	}
	if err != nil {
		l.add(file, "", err.Error())
		return false
	}
	std, err := hujson.Standardize(slices.Clone(data))
	if err == nil {
		err = json.Unmarshal(std, v)
	}
	if err != nil {
		l.add(file, "", err.Error())
		return false
	}
	return true
}

// screenFile is implemented by every screen type through its Screen.
type screenFile interface{ base() *Screen }

func (s *Screen) base() *Screen { return s }

func (l *loader) loadScreens(screens *Screens) Problems {
	entries, err := fs.ReadDir(l.fsys, screensDir)
	if err != nil {
		return nil
	}
	var warnings Problems
	for _, e := range entries {
		file := screensDir + "/" + e.Name()
		name, isJSON := strings.CutSuffix(e.Name(), ".json")
		if e.IsDir() || !isJSON {
			warnings = append(warnings, Problem{File: file, Reason: "unknown screen file, ignored"})
			continue
		}
		v := screens.slot(name)
		if v == nil {
			warnings = append(warnings, Problem{File: file, Reason: fmt.Sprintf("unknown screen %q, ignored", name)})
			continue
		}
		if !l.read(file, name, v) {
			continue
		}
		screen := v.(screenFile).base()
		screen.Silent = l.musicOff(file)
		l.check(file, v, screen.TextStyles)
	}
	return warnings
}

// musicOff reports whether a screen file says "music": null, which decoding
// cannot tell from leaving music out.
func (l *loader) musicOff(file string) bool {
	data, err := fs.ReadFile(l.fsys, file)
	if err != nil {
		return false
	}
	std, err := hujson.Standardize(data)
	if err != nil {
		return false
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(std, &fields) != nil {
		return false
	}
	raw, ok := fields["music"]
	return ok && string(raw) == "null"
}

// check walks a decoded file and reports every reference that does not
// resolve; styles are the file's own text styles, looked up before the
// root's.
func (l *loader) check(file string, v any, styles map[StyleID]TextStyle) {
	c := checker{loader: l, file: file, styles: styles}
	c.walk(reflect.ValueOf(v), "")
}

type checker struct {
	*loader
	file   string
	styles map[StyleID]TextStyle
}

var (
	pathType  = reflect.TypeFor[Path]()
	styleType = reflect.TypeFor[StyleID]()
	fontType  = reflect.TypeFor[FontID]()
)

func (c checker) walk(v reflect.Value, at string) {
	switch v.Kind() {
	case reflect.Pointer:
		if !v.IsNil() {
			c.walk(v.Elem(), at)
		}
	case reflect.Slice:
		for i := range v.Len() {
			c.walk(v.Index(i), at+"/"+strconv.Itoa(i))
		}
	case reflect.Map:
		keys := v.MapKeys()
		slices.SortFunc(keys, func(a, b reflect.Value) int { return strings.Compare(a.String(), b.String()) })
		for _, k := range keys {
			c.walk(v.MapIndex(k), at+"/"+k.String())
		}
	case reflect.Struct:
		c.walkStruct(v, at)
	case reflect.String:
		c.reference(v, at)
	}
}

func (c checker) walkStruct(v reflect.Value, at string) {
	switch x := v.Interface().(type) {
	case Element:
		c.sheetEntry(at+"/animation", x.Sheet, x.Animation, "animation")
		c.sheetEntry(at+"/frame", x.Sheet, x.Frame, "frame")
	case ButtonIcons:
		for _, button := range slices.Sorted(maps.Keys(x.Frames)) {
			c.sheetEntry(at+"/frames/"+button, x.Sheet, x.Frames[button], "frame")
		}
	}
	t := v.Type()
	for i := range t.NumField() {
		f := t.Field(i)
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if !f.IsExported() || name == "-" {
			continue
		}
		if f.Anonymous {
			c.walk(v.Field(i), at)
			continue
		}
		c.walk(v.Field(i), at+"/"+name)
	}
}

func (c checker) reference(v reflect.Value, at string) {
	s := v.String()
	if s == "" {
		return
	}
	switch v.Type() {
	case pathType:
		c.asset(at, s)
	case styleType:
		_, local := c.styles[StyleID(s)]
		_, global := c.rootStyles[StyleID(s)]
		if c.rootOK && !local && !global {
			c.add(c.file, at, fmt.Sprintf("unknown text style %q", s))
		}
	case fontType:
		if _, ok := c.fonts[FontID(s)]; c.rootOK && !ok {
			c.add(c.file, at, fmt.Sprintf("unknown font %q", s))
		}
	}
}

func (c checker) asset(at, p string) {
	if _, err := fs.Stat(c.fsys, p); err != nil {
		c.add(c.file, at, "missing file "+p)
	}
}

// sheetEntry checks that a sheet has the frame or animation id.
func (c checker) sheetEntry(at string, sheet Path, id, what string) {
	if sheet == "" || id == "" {
		return
	}
	s := c.sheet(sheet)
	if s == nil {
		return
	}
	_, isFrame := s.Frames[id]
	_, isAnimation := s.Animations[id]
	if (what == "frame" && !isFrame) || (what == "animation" && !isAnimation) {
		c.add(c.file, at, fmt.Sprintf("unknown %s %q in %s", what, id, sheet))
	}
}

// sheet loads a sprite sheet once; nil when it is missing or broken, which
// is reported once, at the sheet or at the reference to it.
func (l *loader) sheet(p Path) *SpriteSheet {
	if l.tried[p] {
		return l.sheets[p]
	}
	l.tried[p] = true
	file := string(p)
	if _, err := fs.Stat(l.fsys, file); err != nil {
		return nil
	}
	var s SpriteSheet
	if !l.read(file, "spriteSheet", &s) {
		return nil
	}
	image := path.Join(path.Dir(file), string(s.Image))
	if _, err := fs.Stat(l.fsys, image); err != nil {
		l.add(file, "/image", "missing file "+image)
	}
	for _, name := range slices.Sorted(maps.Keys(s.Animations)) {
		for i, fr := range s.Animations[name].Frames {
			if _, ok := s.Frames[fr.Frame]; !ok {
				l.add(file, fmt.Sprintf("/animations/%s/frames/%d/frame", name, i), fmt.Sprintf("unknown frame %q", fr.Frame))
			}
		}
	}
	l.sheets[p] = &s
	return &s
}
