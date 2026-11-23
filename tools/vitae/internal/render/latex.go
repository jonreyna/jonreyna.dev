package render

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/markup"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// latexTemplate is the document's structure. Its delimiters are << and >>, since LaTeX
// is full of braces.
//
//go:embed latex.tmpl
var latexTemplate string

// LaTeX renders r as a LaTeX document for the jrresume class
// (resume/latex/jrresume.cls), which make turns into a PDF with latexmk and pdfLaTeX.
// The document holds the content and its order; the class holds every design decision.
func LaTeX(r *variant.Resolved) ([]byte, error) {

	tmpl, err := template.New("latex").Delims("<<", ">>").Funcs(latexFuncs()).Parse(latexTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse latex template: %w", err)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, r); err != nil {
		return nil, fmt.Errorf("execute latex template: %w", err)
	}

	return buf.Bytes(), nil
}

// latexFuncs are the template's functions. Each returns LaTeX that is safe to insert:
// plain fields are escaped, and inline markup (already validated by lint) is converted.
func latexFuncs() template.FuncMap {

	return template.FuncMap{
		"escape": latexText,
		"markup": func(s string) string { return markup.LaTeX(markup.MustParse(s)) },
		"dates":  Dates,
		"items": func(items []string) string {

			out := make([]string, len(items))
			for i, item := range items {
				out[i] = markup.LaTeX(markup.MustParse(item))
			}

			return strings.Join(out, ", ")
		},
		"contact":       latexContact,
		"schools":       schools,
		"program":       schoolProgram,
		"minors":        latexMinors,
		"schoolDetails": latexSchoolDetails,
	}
}

// latexText escapes a plain field.
func latexText(s string) string { return markup.LaTeX([]markup.Span{{Kind: markup.Text, Text: s}}) }

// latexLink writes a link with its own text.
func latexLink(url, text string) string {
	return markup.LaTeX([]markup.Span{{Kind: markup.Link, Text: text, URL: url}})
}

// latexContact writes the contact line: location, "Remote", email, the phone if the
// private overlay set one, and the profile links, separated by \sep{}.
func latexContact(r *variant.Resolved) string {

	b := &r.Basics
	parts := []string{latexText(b.Location)}
	if b.Remote {
		parts = append(parts, "Remote")
	}

	parts = append(parts, latexLink("mailto:"+b.Email, b.Email))
	if b.Phone != "" {
		parts = append(parts, latexText(b.Phone))
	}

	for _, l := range b.Links {
		parts = append(parts, latexLink(l.URL, l.Text))
	}

	return strings.Join(parts, `\sep{}`)
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

// latexMinors writes "Minors in a, b and c", or nothing when there are none.
func latexMinors(minors []string) string {

	switch len(minors) {
	case 0:
		return ""
	case 1:
		return "Minor in " + latexText(minors[0])
	default:
		return "Minors in " + latexText(sentenceList(minors))
	}
}

// latexSchoolDetails writes the paragraph under a school: its note, then the GPA and
// honors, then the coursework.
func latexSchoolDetails(s *resume.School) string {

	var parts []string
	if s.Note != "" {
		parts = append(parts, markup.LaTeX(markup.MustParse(s.Note)))
	}

	var standing []string
	if s.GPA != "" {
		standing = append(standing, "GPA~"+latexText(s.GPA))
	}

	for _, h := range s.Honors {
		standing = append(standing, latexText(h))
	}

	if len(standing) > 0 {
		parts = append(parts, strings.Join(standing, "; ")+".")
	}

	if len(s.Courses) > 0 {
		names := make([]string, len(s.Courses))
		for i, c := range s.Courses {
			names[i] = c.Name
		}

		parts = append(parts, "Coursework included "+latexText(sentenceList(names))+".")
	}

	return strings.Join(parts, " ")
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
