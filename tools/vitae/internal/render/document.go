package render

import (
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// DocumentSchema is the version of the [Document] shape. Raise it when a field changes
// meaning or disappears, so the website can tell.
const DocumentSchema = 1

// Document is the JSON output, resume.<variant>.json.
type Document struct {
	// Schema is [DocumentSchema].
	Schema int `json:"schema"`
	// Variant is the name of the variant this document shows.
	Variant string `json:"variant"`
	// Sections are the sections to print, in order.
	Sections []string `json:"sections"`
	// Basics is the header.
	Basics DocumentBasics `json:"basics"`
	// Summary is the summary paragraph, if the variant shows one.
	Summary *Text `json:"summary,omitempty"`
	// Experience lists the roles, newest first.
	Experience []DocumentRole `json:"experience,omitempty"`
	// Earlier is the note on earlier work, if the variant shows it.
	Earlier *Text `json:"earlier,omitempty"`
	// Skills lists the skill lines.
	Skills []DocumentSkills `json:"skills,omitempty"`
	// Education lists the schools.
	Education []DocumentEducation `json:"education,omitempty"`
}

// NewDocument builds the JSON document for a resolved résumé whose text lint has
// already validated.
func NewDocument(r *variant.Resolved) Document {
	b := r.Basics
	d := Document{
		Schema:   DocumentSchema,
		Variant:  r.Variant,
		Sections: r.Sections,
		Basics: DocumentBasics{
			Name: b.Name, Label: b.Label, Location: b.Location, Remote: b.Remote,
			Email: b.Email, Phone: b.Phone, Links: append([]resume.Link{}, b.Links...),
		},
		Summary: newText(r.Summary),
		Earlier: newText(r.Earlier),
	}
	for i := range r.Experience {
		d.Experience = append(d.Experience, newDocumentRole(&r.Experience[i]))
	}
	for _, g := range r.Skills {
		d.Skills = append(d.Skills, DocumentSkills{Group: g.Group, Items: g.Items})
	}
	for i := range r.Education {
		d.Education = append(d.Education, newDocumentEducation(&r.Education[i]))
	}
	return d
}
