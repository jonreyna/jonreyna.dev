// Package render writes a resolved résumé in each output format.
package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/filter"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/inline"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/model"
)

// Text is a text field in every form a consumer needs, so no consumer parses markup.
type Text struct {
	Plain string `json:"plain"`
	HTML  string `json:"html"`
}

func text(s string) *Text {
	if s == "" {
		return nil
	}
	spans := inline.MustParse(s)
	return &Text{Plain: inline.Plain(spans), HTML: inline.HTML(spans)}
}

// Document is the JSON output, resume.<variant>.json. The site reads it (Phase 5d);
// its shape is versioned by Schema.
type Document struct {
	Schema     int         `json:"schema"`
	Variant    string      `json:"variant"`
	Sections   []string    `json:"sections"`
	Basics     basics      `json:"basics"`
	Summary    *Text       `json:"summary,omitempty"`
	Experience []role      `json:"experience,omitempty"`
	Earlier    *Text       `json:"earlier,omitempty"`
	Skills     []skills    `json:"skills,omitempty"`
	Education  []education `json:"education,omitempty"`
}

type basics struct {
	Name     string       `json:"name"`
	Label    string       `json:"label"`
	Location string       `json:"location"`
	Remote   bool         `json:"remote"`
	Email    string       `json:"email"`
	Phone    string       `json:"phone,omitempty"`
	Links    []model.Link `json:"links"`
}

type role struct {
	ID         string       `json:"id"`
	Org        string       `json:"org"`
	Title      string       `json:"title"`
	Location   string       `json:"location"`
	Start      model.Month  `json:"start"`
	End        *model.Month `json:"end"`
	Dates      string       `json:"dates"`
	Highlights []highlight  `json:"highlights"`
}

type highlight struct {
	Lead string   `json:"lead,omitempty"`
	Text Text     `json:"text"`
	Tags []string `json:"tags"`
}

type skills struct {
	Group string   `json:"group"`
	Items []string `json:"items"`
}

// education has no degree field unless the degree was conferred.
type education struct {
	Org     string   `json:"org"`
	Program string   `json:"program"`
	Degree  string   `json:"degree,omitempty"`
	Minors  []string `json:"minors,omitempty"`
	Dates   string   `json:"dates"`
	Note    *Text    `json:"note,omitempty"`
	GPA     string   `json:"gpa,omitempty"`
	Honors  []string `json:"honors,omitempty"`
	Courses []string `json:"courses,omitempty"`
}

// Build turns a resolved résumé into the JSON document.
func Build(r *filter.Resolved) Document {
	d := Document{
		Schema:   1,
		Variant:  r.Variant,
		Sections: r.Sections,
		Basics: basics{
			Name: r.Basics.Name, Label: r.Basics.Label, Location: r.Basics.Location,
			Remote: r.Basics.Remote, Email: r.Basics.Email, Phone: r.Basics.Phone,
			Links: r.Basics.Links,
		},
		Summary: text(r.Summary),
		Earlier: text(r.Earlier),
	}
	for _, ro := range r.Experience {
		out := role{
			ID: ro.ID, Org: ro.Org, Title: ro.Title, Location: ro.Location,
			Start: ro.Start, End: ro.End, Dates: Dates(ro.Start, ro.End),
			Highlights: []highlight{},
		}
		for _, h := range ro.Highlights {
			out.Highlights = append(out.Highlights, highlight{Lead: h.Lead, Text: *text(h.Text), Tags: h.Tags})
		}
		d.Experience = append(d.Experience, out)
	}
	for _, g := range r.Skills {
		d.Skills = append(d.Skills, skills{Group: g.Group, Items: g.Items})
	}
	for _, s := range r.Education {
		e := education{
			Org: s.Org, Program: s.Program, Minors: s.Minors,
			Dates: fmt.Sprintf("%d – %d", s.Start, s.End), Note: text(s.Note),
			GPA: s.GPA, Honors: s.Honors,
		}
		if s.Conferred != nil && *s.Conferred {
			e.Degree = s.Degree
		}
		for _, c := range s.Courses {
			e.Courses = append(e.Courses, c.Name)
		}
		d.Education = append(d.Education, e)
	}
	return d
}

// JSON renders the document with two-space indentation, keys in a fixed order and no
// HTML escaping, so the output is stable byte for byte.
func JSON(r *filter.Resolved) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(Build(r)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var months = strings.Fields("Jan Feb Mar Apr May Jun Jul Aug Sep Oct Nov Dec")

// Dates formats a role's span as the résumé prints it: "Oct 2023 – Present".
func Dates(start model.Month, end *model.Month) string {
	f := func(m model.Month) string { return fmt.Sprintf("%s %d", months[m.Month-1], m.Year) }
	if end == nil {
		return f(start) + " – Present"
	}
	return f(start) + " – " + f(*end)
}
