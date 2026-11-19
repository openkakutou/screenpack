package screenpack

import (
	"encoding/json"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/openkakutou/screenpack/schema"
)

// jsonFields lists the JSON names a type decodes, embedded structs included.
func jsonFields(t reflect.Type) []string {
	var out []string
	for i := range t.NumField() {
		f := t.Field(i)
		if f.Anonymous {
			out = append(out, jsonFields(f.Type)...)
			continue
		}
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if name != "" && name != "-" {
			out = append(out, name)
		}
	}
	return out
}

// schemaFields lists a schema's top-level properties; screens also take the
// common screen fields.
func schemaFields(t *testing.T, kind string) []string {
	t.Helper()
	read := func(name string) map[string]any {
		data, err := os.ReadFile("schema/v1/" + name + ".schema.json")
		if err != nil {
			t.Fatal(err)
		}
		var s map[string]any
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		return s
	}
	var out []string
	for name := range read(kind)["properties"].(map[string]any) {
		out = append(out, name)
	}
	if kind != "screenpack" {
		defs := read("common")["$defs"].(map[string]any)
		for name := range defs["screen"].(map[string]any)["properties"].(map[string]any) {
			out = append(out, name)
		}
	}
	return out
}

func TestScreens_HaveATypeForEverySchemaScreen(t *testing.T) {
	for _, name := range schema.Screens() {
		var s Screens
		v := s.slot(name)
		if v == nil {
			t.Errorf("no type for screen %q", name)
			continue
		}
		want := schemaFields(t, name)
		got := jsonFields(reflect.TypeOf(v).Elem())
		for _, f := range want {
			if f != "$schema" && f != "formatVersion" && !slices.Contains(got, f) {
				t.Errorf("%s: schema field %q is not decoded", name, f)
			}
		}
	}
	var s Screens
	if s.slot("credits") != nil {
		t.Error("an unknown screen must have no type")
	}
}

func TestRootFile_DecodesEverySchemaField(t *testing.T) {
	got := jsonFields(reflect.TypeFor[rootFile]())
	for _, f := range schemaFields(t, "screenpack") {
		if f != "$schema" && f != "formatVersion" && !slices.Contains(got, f) {
			t.Errorf("schema field %q is not decoded", f)
		}
	}
}
