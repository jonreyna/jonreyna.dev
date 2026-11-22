package resume

import (
	"bytes"
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

// Resume is the whole of resume.yaml: the single source for every format and variant.
type Resume struct {
	// Schema is the file format version. Only version 1 exists.
	Schema int `yaml:"schema"`
	// Basics is the header: name, label, location and how to get in touch.
	Basics Basics `yaml:"basics"`
	// Summary holds summary paragraphs by key, in inline markup. The key "default" is
	// required; a variant can choose another.
	Summary map[string]string `yaml:"summary"`
	// Tags declares every tag that highlights, courses and variants may use, each with
	// a short description of what it marks.
	Tags map[string]string `yaml:"tags"`
	// Experience lists the roles, newest first.
	Experience []Role `yaml:"experience"`
	// Earlier is a one-line note, in inline markup, on work too old to list as roles.
	Earlier string `yaml:"earlier"`
	// Skills lists the labeled skill lines in the order they print.
	Skills []SkillGroup `yaml:"skills"`
	// Education lists schools, newest first.
	Education []School `yaml:"education"`
	// Variants maps each variant's name to its selection rules. The variant "full" is
	// required.
	Variants map[string]Variant `yaml:"variants"`
}

// Load reads and decodes the résumé at path. Unknown fields are errors. Load doesn't
// validate the content; package lint does that.
func Load(path string) (*Resume, error) {
	var r Resume
	if err := decodeFile(path, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// decodeFile decodes the YAML file at path into v, rejecting fields v doesn't have.
func decodeFile(path string, v any) error {
	data, err := os.ReadFile(path) //nolint:gosec // G304: reading the file the user names is the point
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}
