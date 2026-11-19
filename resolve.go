package screenpack

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// Source says which screenpack a screen or a part comes from. Text styles,
// fonts and assets of something from the default screenpack are looked up
// in the default screenpack.
type Source int

const (
	FromPack Source = iota
	FromDefault
)

func (s Source) String() string {
	if s == FromDefault {
		return "default"
	}
	return "pack"
}

// Resolved is a screenpack with every screen present: a screen the pack
// does not have is the default's, and a part a pack screen leaves out is
// the default screen's part.
type Resolved struct {
	Pack    *Screenpack
	Default *Screenpack
	Screens Screens
	// borrowed holds what comes from the default: "screen" for a whole
	// screen, "screen/part" for a part of a pack screen.
	borrowed map[string]bool
}

// Resolve fills pack's missing screens and parts from def, which must have
// every screen. Neither screenpack is changed.
func Resolve(pack, def *Screenpack) (*Resolved, error) {
	if pack == nil || def == nil {
		return nil, errors.New("resolve: a pack and a default screenpack are needed")
	}
	r := &Resolved{Pack: pack, Default: def, borrowed: map[string]bool{}}
	own := reflect.ValueOf(pack.Screens)
	fallback := reflect.ValueOf(def.Screens)
	out := reflect.ValueOf(&r.Screens).Elem()
	for i := range own.NumField() {
		name := screenName(own.Type().Field(i).Name)
		if fallback.Field(i).IsNil() {
			return nil, fmt.Errorf("resolve: the default screenpack has no %s screen", name)
		}
		if own.Field(i).IsNil() {
			out.Field(i).Set(fallback.Field(i))
			r.borrowed[name] = true
			continue
		}
		screen := reflect.New(own.Field(i).Type().Elem())
		screen.Elem().Set(own.Field(i).Elem())
		r.fillParts(name, screen.Elem(), fallback.Field(i).Elem())
		out.Field(i).Set(screen)
	}
	return r, nil
}

// fillParts sets every nil part of screen to the default screen's.
func (r *Resolved) fillParts(name string, screen, def reflect.Value) {
	t := screen.Type()
	for i := range t.NumField() {
		f := t.Field(i)
		if f.Anonymous || f.Type.Kind() != reflect.Pointer {
			continue
		}
		if screen.Field(i).IsNil() && !def.Field(i).IsNil() {
			screen.Field(i).Set(def.Field(i))
			part, _, _ := strings.Cut(f.Tag.Get("json"), ",")
			r.borrowed[name+"/"+part] = true
		}
	}
}

// Source says where a screen comes from (part "") or, for a pack screen,
// where one of its parts comes from; screen and part are the names of
// docs/format.md, e.g. "title" and "logo".
func (r *Resolved) Source(screen, part string) Source {
	if r.borrowed[screen] || r.borrowed[screen+"/"+part] {
		return FromDefault
	}
	return FromPack
}

// screenName turns a Screens field name into its screen name: MoveList is
// moveList.
func screenName(field string) string {
	runes := []rune(field)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
