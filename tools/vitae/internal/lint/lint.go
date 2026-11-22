// Package lint checks a résumé's structure and content before anything is rendered.
package lint

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/filter"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/inline"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/model"
)

// Severity of an issue. Warnings fail only with --strict.
type Severity int

// Severities.
const (
	Warning Severity = iota
	Error
)

func (s Severity) String() string {
	if s == Error {
		return "error"
	}
	return "warning"
}

// Issue is one problem, located by a path such as "experience[equinix].highlights[2]".
type Issue struct {
	Severity Severity
	Where    string
	Msg      string
}

func (i Issue) String() string { return fmt.Sprintf("%s: %s: %s", i.Severity, i.Where, i.Msg) }

// MaxHighlight is the longest a highlight's plain text may be before a warning. The
// longest on the 2026 résumé is about 620 characters.
const MaxHighlight = 650

// degreeWords are the words that, outside the note, would claim a degree as earned.
var degreeWords = regexp.MustCompile(`(?i)\b(B\.?S\.?|B\.?A\.?|bachelor|degree|graduated)\b`)

type checker struct {
	r      *model.Resume
	issues []Issue
	used   map[string]bool
}

func (c *checker) errorf(where, format string, args ...any) {
	c.issues = append(c.issues, Issue{Error, where, fmt.Sprintf(format, args...)})
}

func (c *checker) warnf(where, format string, args ...any) {
	c.issues = append(c.issues, Issue{Warning, where, fmt.Sprintf(format, args...)})
}

// Check returns every issue in r, errors and warnings, in file order.
func Check(r *model.Resume) []Issue {
	c := &checker{r: r, used: map[string]bool{}}
	c.basics()
	c.summary()
	c.experience()
	c.text("earlier", r.Earlier)
	c.skills()
	c.education()
	c.variants()
	for _, tag := range sortedKeys(r.Tags) {
		if !c.used[tag] {
			c.warnf("tags."+tag, "declared but never used")
		}
	}
	return c.issues
}

// HasErrors reports whether issues would fail validation (warnings count when strict).
func HasErrors(issues []Issue, strict bool) bool {
	for _, i := range issues {
		if i.Severity == Error || strict {
			return true
		}
	}
	return false
}

func (c *checker) basics() {
	b := c.r.Basics
	if c.r.Schema != 1 {
		c.errorf("schema", "must be 1, not %d", c.r.Schema)
	}
	if b.Name == "" {
		c.errorf("basics.name", "required")
	}
	if !strings.Contains(b.Email, "@") {
		c.errorf("basics.email", "required, and must be an address")
	}
	if b.Phone != "" {
		c.errorf("basics.phone", "belongs in resume.private.yaml, never in the public file")
	}
	for i, l := range b.Links {
		where := fmt.Sprintf("basics.links[%d]", i)
		if !strings.HasPrefix(l.URL, "https://") {
			c.errorf(where, "url must start with https://")
		}
		if l.Kind == "" || l.Text == "" {
			c.errorf(where, "kind and text are required")
		}
	}
}

func (c *checker) summary() {
	if c.r.Summary["default"] == "" {
		c.errorf("summary.default", "required")
	}
	for _, k := range sortedKeys(c.r.Summary) {
		c.text("summary."+k, c.r.Summary[k])
	}
}

func (c *checker) experience() {
	if len(c.r.Experience) == 0 {
		c.errorf("experience", "at least one role is required")
	}
	seen := map[string]bool{}
	for i, role := range c.r.Experience {
		where := fmt.Sprintf("experience[%s]", role.ID)
		switch {
		case role.ID == "":
			where = fmt.Sprintf("experience[%d]", i)
			c.errorf(where, "id is required")
		case seen[role.ID]:
			c.errorf(where, "duplicate id")
		}
		seen[role.ID] = true
		if role.Org == "" || role.Title == "" {
			c.errorf(where, "org and title are required")
		}
		if role.Start == (model.Month{}) {
			c.errorf(where, "start is required")
		}
		if role.End == nil && i > 0 {
			c.errorf(where, "only the first role can be current (end: null)")
		}
		if role.End != nil && role.End.Before(role.Start) {
			c.errorf(where, "ends (%s) before it starts (%s)", role.End, role.Start)
		}
		if i > 0 && c.r.Experience[i-1].Start.Before(role.Start) {
			c.errorf(where, "roles must be newest first; this one starts after the one above")
		}
		if len(role.Highlights) == 0 {
			c.warnf(where, "no highlights")
		}
		for j, h := range role.Highlights {
			hw := fmt.Sprintf("%s.highlights[%d]", where, j)
			if h.Text == "" {
				c.errorf(hw, "text is required")
			}
			c.text(hw+".lead", h.Lead)
			if spans := c.text(hw, h.Text); spans != nil {
				if n := len([]rune(inline.Plain(spans))); n > MaxHighlight {
					c.warnf(hw, "%d characters; keep highlights under %d", n, MaxHighlight)
				}
			}
			if h.Weight < 1 || h.Weight > 100 {
				c.errorf(hw, "weight must be 1 to 100, not %d", h.Weight)
			}
			c.tags(hw, h.Tags)
		}
	}
}

func (c *checker) skills() {
	for i, g := range c.r.Skills {
		where := fmt.Sprintf("skills[%d]", i)
		if g.Group == "" || len(g.Items) == 0 {
			c.errorf(where, "group and items are required")
		}
		for _, item := range g.Items {
			c.text(where, item)
		}
	}
}

func (c *checker) education() {
	for i, s := range c.r.Education {
		where := fmt.Sprintf("education[%d]", i)
		if s.ID != "" {
			where = fmt.Sprintf("education[%s]", s.ID)
		}
		if s.Org == "" || s.Program == "" {
			c.errorf(where, "org and program are required")
		}
		if s.Start == 0 || s.End < s.Start {
			c.errorf(where, "start and end years are required, start first")
		}
		switch {
		case s.Conferred == nil:
			c.errorf(where, "conferred must be stated: true or false")
		case !*s.Conferred:
			if s.Note == "" {
				c.errorf(where, "an unconferred degree needs a note saying so")
			}
			// Only the note may mention the degree; every other printed field must not.
			fields := [][2]string{{"program", s.Program}, {"gpa", s.GPA}}
			for _, h := range s.Honors {
				fields = append(fields, [2]string{"honors", h})
			}
			for _, m := range s.Minors {
				fields = append(fields, [2]string{"minors", m})
			}
			for _, f := range fields {
				if degreeWords.MatchString(f[1]) {
					c.errorf(where+"."+f[0], "claims a degree that was not conferred: %q", f[1])
				}
			}
		}
		c.text(where+".note", s.Note)
		for _, course := range s.Courses {
			c.tags(where+".courses", course.Tags)
		}
	}
}

func (c *checker) variants() {
	if _, ok := c.r.Variants["full"]; !ok {
		c.errorf("variants", "a variant named full is required")
	}
	for _, name := range sortedKeys(c.r.Variants) {
		v := c.r.Variants[name]
		where := "variants." + name
		if len(v.Sections) == 0 {
			c.errorf(where, "sections are required")
		}
		for i, s := range v.Sections {
			if !slices.Contains(model.Sections, s) {
				c.errorf(where, "unknown section %q (known: %s)", s, strings.Join(model.Sections, ", "))
			}
			if slices.Contains(v.Sections[:i], s) {
				c.errorf(where, "section %q listed twice", s)
			}
		}
		if v.Summary != "" && c.r.Summary[v.Summary] == "" {
			c.errorf(where, "summary %q is not in summary", v.Summary)
		}
		c.tags(where+".boost", v.Boost)
		c.tags(where+".exclude", v.Exclude)
		c.tags(where+".show_courses", v.ShowCourses)
		if v.MaxRoles < 0 || v.MaxHighlights < 0 || v.MinWeight < 0 || v.BoostBy < 0 {
			c.errorf(where, "limits can't be negative")
		}
		if out, err := filter.Apply(c.r, name); err == nil {
			n := 0
			for _, role := range out.Experience {
				n += len(role.Highlights)
			}
			if n == 0 {
				c.errorf(where, "selects no highlights")
			}
		}
	}
}

func (c *checker) tags(where string, tags []string) {
	for _, t := range tags {
		if _, ok := c.r.Tags[t]; !ok {
			c.errorf(where, "unknown tag %q; declare it under tags", t)
		}
		c.used[t] = true
	}
}

// text parses inline markup and checks typography; it returns nil if s doesn't parse.
func (c *checker) text(where, s string) []inline.Span {
	if s == "" {
		return nil
	}
	spans, err := inline.Parse(s)
	if err != nil {
		c.errorf(where, "markup: %v", err)
		return nil
	}
	plain := inline.Plain(spans)
	if strings.ContainsAny(plain, `'"`) {
		c.warnf(where, "straight quote; use ’ “ ” instead")
	}
	if strings.Contains(plain, "  ") {
		c.warnf(where, "double space")
	}
	if strings.Contains(plain, " -- ") || strings.Contains(plain, " - ") {
		c.warnf(where, "hyphen used as a dash; use – for ranges or — for breaks")
	}
	return spans
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
