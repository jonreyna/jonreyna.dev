// Package model holds the résumé as written in resume.yaml, and loads it.
package model

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"go.yaml.in/yaml/v3"
)

// Resume is the whole of resume.yaml.
type Resume struct {
	Schema     int                `yaml:"schema"`
	Basics     Basics             `yaml:"basics"`
	Summary    map[string]string  `yaml:"summary"`
	Tags       map[string]string  `yaml:"tags"`
	Experience []Role             `yaml:"experience"`
	Earlier    string             `yaml:"earlier"`
	Skills     []SkillGroup       `yaml:"skills"`
	Education  []School           `yaml:"education"`
	Variants   map[string]Variant `yaml:"variants"`
}

// Basics is the header: who, where, and how to reach them.
type Basics struct {
	Name     string `yaml:"name"`
	Label    string `yaml:"label"`
	Location string `yaml:"location"`
	Remote   bool   `yaml:"remote"`
	Email    string `yaml:"email"`
	Phone    string `yaml:"phone"` // private overlay only
	Links    []Link `yaml:"links"`
}

// Link is a profile or site link.
type Link struct {
	Kind string `yaml:"kind" json:"kind"`
	URL  string `yaml:"url" json:"url"`
	Text string `yaml:"text" json:"text"`
}

// Role is one job.
type Role struct {
	ID         string      `yaml:"id"`
	Org        string      `yaml:"org"`
	Title      string      `yaml:"title"`
	Location   string      `yaml:"location"`
	Start      Month       `yaml:"start"`
	End        *Month      `yaml:"end"` // nil: present
	Highlights []Highlight `yaml:"highlights"`
}

// Highlight is one bullet. Lead is the optional bold lead-in.
type Highlight struct {
	Lead   string   `yaml:"lead"`
	Text   string   `yaml:"text"`
	Tags   []string `yaml:"tags"`
	Weight int      `yaml:"weight"`
}

// SkillGroup is one labeled line of skills.
type SkillGroup struct {
	Group string   `yaml:"group"`
	Items []string `yaml:"items"`
}

// School is one education entry. Conferred is a pointer so that leaving it out is an
// error rather than a silent "false" or "true".
type School struct {
	ID        string   `yaml:"id"`
	Org       string   `yaml:"org"`
	Program   string   `yaml:"program"`
	Degree    string   `yaml:"degree"`
	Conferred *bool    `yaml:"conferred"`
	Minors    []string `yaml:"minors"`
	Start     int      `yaml:"start"`
	End       int      `yaml:"end"`
	Note      string   `yaml:"note"`
	GPA       string   `yaml:"gpa"`
	Honors    []string `yaml:"honors"`
	Courses   []Course `yaml:"courses"`
}

// Course is one course named in the coursework line.
type Course struct {
	Name string   `yaml:"name"`
	Tags []string `yaml:"tags"`
}

// Variant selects sections and highlights. It never rewords anything.
type Variant struct {
	Description   string   `yaml:"description"`
	Sections      []string `yaml:"sections"`
	Summary       string   `yaml:"summary"` // key into Resume.Summary; "default" if empty
	Label         string   `yaml:"label"`   // replaces Basics.Label if set
	Since         *Month   `yaml:"since"`   // drop roles that ended before this
	MaxRoles      int      `yaml:"max_roles"`
	MaxHighlights int      `yaml:"max_highlights"`
	MinWeight     int      `yaml:"min_weight"`
	Boost         []string `yaml:"boost"`
	BoostBy       int      `yaml:"boost_by"` // 20 if zero
	Exclude       []string `yaml:"exclude"`
	ShowCourses   []string `yaml:"show_courses"` // course tags to show; all if empty
}

// Sections a variant may list, in no particular order.
var Sections = []string{"summary", "skills", "experience", "earlier", "education"}

// Month is a calendar month, written "YYYY-MM".
type Month struct {
	Year, Month int
}

// ParseMonth reads "YYYY-MM".
func ParseMonth(s string) (Month, error) {
	if len(s) != 7 || s[4] != '-' {
		return Month{}, fmt.Errorf("%q is not a YYYY-MM month", s)
	}
	y, err1 := strconv.Atoi(s[:4])
	m, err2 := strconv.Atoi(s[5:])
	if err1 != nil || err2 != nil || m < 1 || m > 12 || y < 1900 {
		return Month{}, fmt.Errorf("%q is not a YYYY-MM month", s)
	}
	return Month{y, m}, nil
}

func (m Month) String() string { return fmt.Sprintf("%04d-%02d", m.Year, m.Month) }

// Before reports whether m is earlier than o.
func (m Month) Before(o Month) bool {
	return m.Year < o.Year || (m.Year == o.Year && m.Month < o.Month)
}

// UnmarshalYAML accepts a "YYYY-MM" scalar.
func (m *Month) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.ScalarNode {
		return fmt.Errorf("line %d: a month must be a YYYY-MM string", n.Line)
	}
	v, err := ParseMonth(n.Value)
	if err != nil {
		return fmt.Errorf("line %d: %w", n.Line, err)
	}
	*m = v
	return nil
}

// MarshalJSON writes "YYYY-MM".
func (m Month) MarshalJSON() ([]byte, error) { return []byte(`"` + m.String() + `"`), nil }

// Private is the git-ignored overlay: fields that never appear in public outputs.
type Private struct {
	Basics struct {
		Phone string `yaml:"phone"`
	} `yaml:"basics"`
}

// Load reads resume.yaml, rejecting unknown fields so typos fail loudly.
func Load(path string) (*Resume, error) {
	var r Resume
	if err := decodeFile(path, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// ApplyPrivate merges the private overlay at path into r.
func ApplyPrivate(r *Resume, path string) error {
	var p Private
	if err := decodeFile(path, &p); err != nil {
		return err
	}
	if p.Basics.Phone != "" {
		r.Basics.Phone = p.Basics.Phone
	}
	return nil
}

func decodeFile(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}
