package lint

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/markup"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// degreeClaim matches words that, outside an unconferred degree's note, would claim
// the degree as earned.
var degreeClaim = regexp.MustCompile(`(?i)\b(B\.?S\.?|B\.?A\.?|bachelor|degree|graduated)\b`)

const (
	minWeight = 1
	maxWeight = 100
)

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
	for _, k := range slices.Sorted(maps.Keys(c.r.Summary)) {
		c.text("summary."+k, c.r.Summary[k])
	}
}

func (c *checker) experience() {
	if len(c.r.Experience) == 0 {
		c.errorf("experience", "at least one role is required")
	}
	seen := map[string]bool{}
	for i := range c.r.Experience {
		role := &c.r.Experience[i]
		where := fmt.Sprintf("experience[%s]", role.ID)
		switch {
		case role.ID == "":
			where = fmt.Sprintf("experience[%d]", i)
			c.errorf(where, "id is required")
		case seen[role.ID]:
			c.errorf(where, "duplicate id")
		}
		seen[role.ID] = true
		c.role(where, i, role)
		for j, h := range role.Highlights {
			c.highlight(fmt.Sprintf("%s.highlights[%d]", where, j), h)
		}
	}
}

// role checks one role's fields and its place in the list.
func (c *checker) role(where string, i int, role *resume.Role) {
	if role.Org == "" || role.Title == "" {
		c.errorf(where, "org and title are required")
	}
	if role.Start.IsZero() {
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
}

func (c *checker) highlight(where string, h resume.Highlight) {
	if h.Text == "" {
		c.errorf(where, "text is required")
	}
	c.text(where+".lead", h.Lead)
	if spans := c.text(where, h.Text); spans != nil {
		if n := len([]rune(markup.Plain(spans))); n > MaxHighlight {
			c.warnf(where, "%d characters; keep highlights under %d", n, MaxHighlight)
		}
	}
	if h.Weight < minWeight || h.Weight > maxWeight {
		c.errorf(where, "weight must be %d to %d, not %d", minWeight, maxWeight, h.Weight)
	}
	c.tags(where, h.Tags)
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
	for i := range c.r.Education {
		s := &c.r.Education[i]
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
		c.degree(where, s)
		c.text(where+".note", s.Note)
		for _, course := range s.Courses {
			c.tags(where+".courses", course.Tags)
		}
	}
}

// degree enforces the rule for unconferred degrees: say so, and claim it nowhere.
func (c *checker) degree(where string, s *resume.School) {
	switch {
	case s.Conferred == nil:
		c.errorf(where, "conferred must be stated: true or false")
	case !*s.Conferred:
		if s.Note == "" {
			c.errorf(where, "an unconferred degree needs a note saying so")
		}
		// Only the note may mention the degree; every other printed field must not.
		fields := make([][2]string, 0, 2+len(s.Honors)+len(s.Minors))
		fields = append(fields, [2]string{"program", s.Program}, [2]string{"gpa", s.GPA})
		for _, h := range s.Honors {
			fields = append(fields, [2]string{"honors", h})
		}
		for _, m := range s.Minors {
			fields = append(fields, [2]string{"minors", m})
		}
		for _, f := range fields {
			if degreeClaim.MatchString(f[1]) {
				c.errorf(where+"."+f[0], "claims a degree that was not conferred: %q", f[1])
			}
		}
	}
}

func (c *checker) variants() {
	if _, ok := c.r.Variants["full"]; !ok {
		c.errorf("variants", "a variant named full is required")
	}
	for _, name := range slices.Sorted(maps.Keys(c.r.Variants)) {
		v := c.r.Variants[name]
		c.variant("variants."+name, name, &v)
	}
}

func (c *checker) variant(where, name string, v *resume.Variant) {
	c.sections(where, v.Sections)
	if v.Summary != "" && c.r.Summary[v.Summary] == "" {
		c.errorf(where, "summary %q is not in summary", v.Summary)
	}
	c.tags(where+".boost", v.Boost)
	c.tags(where+".exclude", v.Exclude)
	c.tags(where+".show_courses", v.ShowCourses)
	if v.MaxRoles < 0 || v.MaxHighlights < 0 || v.MinWeight < 0 || v.BoostBy < 0 {
		c.errorf(where, "limits can't be negative")
	}
	if out, err := variant.Apply(c.r, name); err == nil && out.Has(resume.SectionExperience) {
		n := 0
		for i := range out.Experience {
			n += len(out.Experience[i].Highlights)
		}
		if n == 0 {
			c.errorf(where, "selects no highlights")
		}
	}
}

// sections checks a variant's section list: known names, each once.
func (c *checker) sections(where string, sections []string) {
	if len(sections) == 0 {
		c.errorf(where, "sections are required")
	}
	for i, s := range sections {
		if !slices.Contains(resume.Sections(), s) {
			c.errorf(where, "unknown section %q (known: %s)", s, strings.Join(resume.Sections(), ", "))
		}
		if slices.Contains(sections[:i], s) {
			c.errorf(where, "section %q listed twice", s)
		}
	}
}
