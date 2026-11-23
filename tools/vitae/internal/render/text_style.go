package render

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/markup"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// textStyle is one template-driven format: LaTeX, plain text or Markdown. The
// template lays out the document; the style says how every shared piece is written
// (escaped fields, inline markup, links, the contact line, the school details), so the
// formats differ only where they must.
type textStyle struct {
	// name names the template in errors ("latex").
	name string
	// source is the template text.
	source string
	// left and right are the template's action delimiters.
	left, right string
	// inline converts parsed inline markup, escaping it as the format needs.
	inline func([]markup.Span) string
	// sep goes between the items of the contact line.
	sep string
}

// render executes the style's template on r.
func (s *textStyle) render(r *variant.Resolved) ([]byte, error) {

	tmpl, err := template.New(s.name).Delims(s.left, s.right).Funcs(s.funcs()).Parse(s.source)
	if err != nil {
		return nil, fmt.Errorf("parse %s template: %w", s.name, err)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, r); err != nil {
		return nil, fmt.Errorf("execute %s template: %w", s.name, err)
	}

	return buf.Bytes(), nil
}

// funcs are the template's functions. Each returns text that is safe to insert: plain
// fields are escaped, and inline markup (already validated by lint) is converted.
func (s *textStyle) funcs() template.FuncMap {

	return template.FuncMap{
		"escape":        s.text,
		"markup":        s.markup,
		"items":         s.items,
		"contact":       s.contact,
		"minors":        s.minors,
		"schoolDetails": s.schoolDetails,
		"dates":         Dates,
		"schools":       schools,
		"program":       schoolProgram,
	}
}

// text escapes a plain field.
func (s *textStyle) text(v string) string {
	return s.inline([]markup.Span{{Kind: markup.Text, Text: v}})
}

// markup converts a field in inline markup.
func (s *textStyle) markup(v string) string { return s.inline(markup.MustParse(v)) }

// link writes a link with its own text.
func (s *textStyle) link(url, text string) string {
	return s.inline([]markup.Span{{Kind: markup.Link, Text: text, URL: url}})
}

// items writes a skill line's items, each in inline markup, separated by commas.
func (s *textStyle) items(items []string) string {

	out := make([]string, len(items))
	for i, item := range items {
		out[i] = s.markup(item)
	}

	return strings.Join(out, ", ")
}

// contact writes the contact line: location, "Remote", email, the phone if the
// private overlay set one, and the profile links, separated by sep.
func (s *textStyle) contact(r *variant.Resolved) string {

	b := &r.Basics
	parts := []string{s.text(b.Location)}
	if b.Remote {
		parts = append(parts, "Remote")
	}

	parts = append(parts, s.link("mailto:"+b.Email, b.Email))
	if b.Phone != "" {
		parts = append(parts, s.text(b.Phone))
	}

	for _, l := range b.Links {
		parts = append(parts, s.link(l.URL, l.Text))
	}

	return strings.Join(parts, s.sep)
}

// minors writes "Minors in a, b and c", or nothing when there are none.
func (s *textStyle) minors(minors []string) string {

	switch len(minors) {
	case 0:
		return ""
	case 1:
		return "Minor in " + s.text(minors[0])
	default:
		return "Minors in " + s.text(sentenceList(minors))
	}
}

// schoolDetails writes the paragraph under a school: its note, then the GPA and
// honors, then the coursework.
func (s *textStyle) schoolDetails(sc *resume.School) string {

	var parts []string
	if sc.Note != "" {
		parts = append(parts, s.markup(sc.Note))
	}

	var standing []string
	if sc.GPA != "" {
		standing = append(standing, s.text("GPA "+sc.GPA))
	}

	for _, h := range sc.Honors {
		standing = append(standing, s.text(h))
	}

	if len(standing) > 0 {
		parts = append(parts, strings.Join(standing, "; ")+".")
	}

	if len(sc.Courses) > 0 {
		names := make([]string, len(sc.Courses))
		for i, c := range sc.Courses {
			names[i] = c.Name
		}

		parts = append(parts, "Coursework included "+s.text(sentenceList(names))+".")
	}

	return strings.Join(parts, " ")
}

// schools returns pointers to the kept schools, so the template's functions can take
// them without copying.
func schools(r *variant.Resolved) []*resume.School {

	out := make([]*resume.School, len(r.Education))
	for i := range r.Education {
		out[i] = &r.Education[i]
	}

	return out
}

// schoolProgram is what a school entry says was studied: the degree when it was
// conferred, the program otherwise (the degree rule).
func schoolProgram(s *resume.School) string {

	if s.Conferred != nil && *s.Conferred && s.Degree != "" {
		return s.Degree
	}

	return s.Program
}

// sentenceList joins items as a sentence does: "a, b and c", with no serial comma.
func sentenceList(items []string) string {

	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	default:
		return strings.Join(items[:len(items)-1], ", ") + " and " + items[len(items)-1]
	}
}
