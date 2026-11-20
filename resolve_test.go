package screenpack

import (
	"os"
	"reflect"
	"testing"

	"github.com/openkakutou/screenpack/schema"
)

func loadDefault(t *testing.T) *Screenpack {
	t.Helper()
	def, err := Load(os.DirFS(defaultDir))
	if err != nil {
		t.Fatal(err)
	}
	return def
}

// onlyTitle is a pack whose single screen sets its menu and nothing else.
func onlyTitle(t *testing.T) *Screenpack {
	t.Helper()
	fsys := valid()
	delete(fsys, "screens/pause.json")
	fsys["screens/title.json"] = file(`{ "formatVersion": 1, "menu": { "style": "default", "spacing": 3 } }`)
	sp, err := Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	return sp
}

func TestResolve_TakesEveryMissingScreenFromTheDefault(t *testing.T) {
	def := loadDefault(t)
	r, err := Resolve(onlyTitle(t), def)
	if err != nil {
		t.Fatal(err)
	}
	screens := reflect.ValueOf(r.Screens)
	for i := range screens.NumField() {
		if screens.Field(i).IsNil() {
			t.Errorf("%s not resolved", screens.Type().Field(i).Name)
		}
	}
	if r.Screens.Versus != def.Screens.Versus {
		t.Error("a missing screen is the default's own")
	}
	if got := r.Source("versus", ""); got != FromDefault {
		t.Errorf("versus source = %v", got)
	}
	if got := r.Source("title", ""); got != FromPack {
		t.Errorf("title source = %v", got)
	}
}

func TestResolve_FillsMissingPartsOfAPackScreen(t *testing.T) {
	pack, def := onlyTitle(t), loadDefault(t)
	r, err := Resolve(pack, def)
	if err != nil {
		t.Fatal(err)
	}
	title := r.Screens.Title
	if title.Menu.Spacing != 3 || r.Source("title", "menu") != FromPack {
		t.Errorf("the pack's own part is kept: %+v", title.Menu)
	}
	if title.Logo != def.Screens.Title.Logo || r.Source("title", "logo") != FromDefault {
		t.Errorf("logo = %+v from %v, want the default's", title.Logo, r.Source("title", "logo"))
	}
	if len(title.Background) != 0 || title.Music != nil {
		t.Error("a screen's own fields are never merged with the default's")
	}
	if pack.Screens.Title.Logo != nil {
		t.Error("resolving must not change the pack")
	}
}

func TestResolve_DefaultAgainstItselfComesFromThePack(t *testing.T) {
	def := loadDefault(t)
	r, err := Resolve(def, def)
	if err != nil {
		t.Fatal(err)
	}
	if r.Source("pause", "menu") != FromPack || r.Source("continue", "") != FromPack {
		t.Error("everything the pack has comes from the pack")
	}
}

func TestResolve_RefusesAnIncompleteDefault(t *testing.T) {
	pack := onlyTitle(t)
	if _, err := Resolve(pack, nil); err == nil {
		t.Error("want an error without a default")
	}
	if _, err := Resolve(pack, onlyTitle(t)); err == nil {
		t.Error("want an error when the default lacks screens")
	}
}

func TestSource_String(t *testing.T) {
	if FromPack.String() != "pack" || FromDefault.String() != "default" {
		t.Errorf("%v %v", FromPack, FromDefault)
	}
}

func TestScreenName_MatchesTheSchemaScreens(t *testing.T) {
	known := map[string]bool{}
	for _, name := range schema.Screens() {
		known[name] = true
	}
	fields := reflect.TypeFor[Screens]()
	for i := range fields.NumField() {
		if name := screenName(fields.Field(i).Name); !known[name] {
			t.Errorf("field %s gives %q, not a screen", fields.Field(i).Name, name)
		}
	}
}

func TestResolve_BorrowsTheErrorScreenFromTheDefault(t *testing.T) {
	def := loadDefault(t)
	r, err := Resolve(onlyTitle(t), def)
	if err != nil {
		t.Fatal(err)
	}
	if r.Screens.Error != def.Screens.Error {
		t.Error("a pack without an error screen gets the default's")
	}
	if got := r.Source("error", ""); got != FromDefault {
		t.Errorf("error source = %v", got)
	}
}
