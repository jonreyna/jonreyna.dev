package variant

import (
	"slices"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

// Resolved is a résumé with one variant applied. Sections the variant leaves out are
// empty.
type Resolved struct {
	// Variant is the name of the variant applied.
	Variant string
	// Sections are the sections to print, in order.
	Sections []string
	// Basics is the header, with the variant's label if it sets one.
	Basics resume.Basics
	// Summary is the chosen summary paragraph, in inline markup.
	Summary string
	// Experience holds the kept roles, each with only its kept highlights.
	Experience []resume.Role
	// Earlier is the note on earlier work, in inline markup.
	Earlier string
	// Skills are the skill lines.
	Skills []resume.SkillGroup
	// Education holds the schools, each with only its kept courses.
	Education []resume.School
}

// Has reports whether the variant shows the named section.
func (r *Resolved) Has(section string) bool { return slices.Contains(r.Sections, section) }
