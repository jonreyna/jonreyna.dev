package lint

import (
	"maps"
	"slices"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

// MaxHighlight is the longest a highlight's plain text may be before a warning. The
// longest on the 2026 résumé is about 620 characters.
const MaxHighlight = 650

// Check returns every issue in r, errors and warnings, in file order.
func Check(r *resume.Resume) []Issue {

	c := &checkState{r: r, used: map[string]bool{}}
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
