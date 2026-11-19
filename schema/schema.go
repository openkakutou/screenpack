// Package schema holds the JSON Schemas of the screenpack format
// (docs/format.md) and checks a JSONC document against them.
package schema

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"github.com/tailscale/hujson"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:embed v1/*.schema.json
var files embed.FS

// BaseURL is where the schemas are published: each schema's $id is BaseURL
// followed by its file name.
const BaseURL = "https://raw.githubusercontent.com/openkakutou/screenpack/main/schema/v1/"

const suffix = ".schema.json"

// Kinds that are not screens: the root file, the sprite-sheet description,
// and the definitions shared by every schema.
const (
	kindRoot   = "screenpack"
	kindSheet  = "spriteSheet"
	kindCommon = "common"
)

var compiled = sync.OnceValues(compile)

func kinds() []string {
	entries, err := fs.ReadDir(files, "v1")
	if err != nil {
		panic(err) // the embedded folder always exists
	}
	var out []string
	for _, e := range entries {
		out = append(out, strings.TrimSuffix(e.Name(), suffix))
	}
	return out
}

func compile() (map[string]*jsonschema.Schema, error) {
	c := jsonschema.NewCompiler()
	all := kinds()
	for _, kind := range all {
		data, err := files.ReadFile("v1/" + kind + suffix)
		if err != nil {
			return nil, err
		}
		doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("%s%s: %w", kind, suffix, err)
		}
		if err := c.AddResource(BaseURL+kind+suffix, doc); err != nil {
			return nil, err
		}
	}
	out := map[string]*jsonschema.Schema{}
	for _, kind := range all {
		if kind == kindCommon {
			continue
		}
		s, err := c.Compile(BaseURL + kind + suffix)
		if err != nil {
			return nil, err
		}
		out[kind] = s
	}
	return out, nil
}

// Screens lists the screen names, sorted: the files a screenpack may hold
// under screens/, without their .json extension.
func Screens() []string {
	var out []string
	for _, kind := range kinds() {
		if kind != kindRoot && kind != kindSheet && kind != kindCommon {
			out = append(out, kind)
		}
	}
	return out
}

// Validate checks a JSONC document against the schema of kind: "screenpack"
// for screenpack.json, "spriteSheet" for a sprite-sheet description, or a
// screen name from Screens. The error names the JSON path of each problem.
func Validate(kind string, jsonc []byte) error {
	schemas, err := compiled()
	if err != nil {
		return err
	}
	s, ok := schemas[kind]
	if !ok {
		return fmt.Errorf("unknown document kind %q", kind)
	}
	std, err := hujson.Standardize(slices.Clone(jsonc))
	if err != nil {
		return fmt.Errorf("invalid JSONC: %w", err)
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(std))
	if err != nil {
		return err
	}
	err = s.Validate(doc)
	var verr *jsonschema.ValidationError
	if errors.As(err, &verr) {
		return issues(verr)
	}
	return err
}

// Issue is one way a document breaks its schema.
type Issue struct {
	// Path is the JSON Pointer of the faulty value, "/" for the document.
	Path   string
	Reason string
}

// Invalid lists every Issue of a document that breaks its schema.
type Invalid []Issue

func (v Invalid) Error() string {
	lines := make([]string, len(v))
	for i, issue := range v {
		lines[i] = issue.Path + ": " + issue.Reason
	}
	return strings.Join(lines, "\n")
}

var printer = message.NewPrinter(language.English)

// issues keeps the leaves of the error tree: the composite ones (allOf,
// if/then…) only say that a nested check failed.
func issues(e *jsonschema.ValidationError) Invalid {
	if len(e.Causes) > 0 {
		var out Invalid
		for _, c := range e.Causes {
			out = append(out, issues(c)...)
		}
		return out
	}
	reason := e.ErrorKind.LocalizedString(printer)
	if _, ok := e.ErrorKind.(*kind.FalseSchema); ok {
		reason = "not allowed here"
	}
	return Invalid{{Path: pointer(e.InstanceLocation), Reason: reason}}
}

func pointer(tokens []string) string {
	if len(tokens) == 0 {
		return "/"
	}
	escaped := make([]string, len(tokens))
	for i, t := range tokens {
		escaped[i] = strings.NewReplacer("~", "~0", "/", "~1").Replace(t)
	}
	return "/" + strings.Join(escaped, "/")
}
