package lint

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/markup"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

// MaxHighlight is the longest a highlight's plain text may be before a warning. The
// longest on the 2026 résumé is about 620 characters.
const MaxHighlight = 650

// Check returns every issue in r, errors and warnings, in file order.
func Check(r *resume.Resume) []Issue {
	c := &checker{r: r, used: map[string]bool{}}
	c.basics()
	c.summary()
	c.experience()
	c.text("earlier", r.Earlier)
	c.skills()
	c.education()
	c.variants()
	for _, tag := range slices.Sorted(maps.Keys(r.Tags)) {
		if !c.used[tag] {
			c.warnf("tags."+tag, "declared but never used")
		}
	}
	return c.issues
}

// HasErrors reports whether issues fail validation. In strict mode warnings count too.
func HasErrors(issues []Issue, strict bool) bool {
	return slices.ContainsFunc(issues, func(i Issue) bool { return strict || i.Severity == Error })
}

// checker collects issues for one [Check] call.
type checker struct {
	r      *resume.Resume
	issues []Issue
	used   map[string]bool // tags seen anywhere, for the unused-tag warning
}

func (c *checker) errorf(where, format string, args ...any) {
	c.issues = append(c.issues, Issue{Severity: Error, Where: where, Msg: fmt.Sprintf(format, args...)})
}

func (c *checker) warnf(where, format string, args ...any) {
	c.issues = append(c.issues, Issue{Severity: Warning, Where: where, Msg: fmt.Sprintf(format, args...)})
}

// tags reports undeclared tags and records the ones it sees as used.
func (c *checker) tags(where string, tags []string) {
	for _, t := range tags {
		if _, ok := c.r.Tags[t]; !ok {
			c.errorf(where, "unknown tag %q; declare it under tags", t)
		}
		c.used[t] = true
	}
}

// text parses s's inline markup and checks its typography. It returns the spans, or nil
// if s is empty or doesn't parse.
func (c *checker) text(where, s string) []markup.Span {
	if s == "" {
		return nil
	}
	spans, err := markup.Parse(s)
	if err != nil {
		c.errorf(where, "markup: %v", err)
		return nil
	}
	plain := markup.Plain(spans)
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
