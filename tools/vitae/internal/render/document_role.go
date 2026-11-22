package render

import (
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

// DocumentRole is one role in a [Document].
type DocumentRole struct {
	// ID identifies the role ("equinix").
	ID string `json:"id"`
	// Org is the employer's name.
	Org string `json:"org"`
	// Title is the job title as printed.
	Title string `json:"title"`
	// Location is where the job was.
	Location string `json:"location"`
	// Start is the first month, "YYYY-MM".
	Start resume.Month `json:"start"`
	// End is the last month, "YYYY-MM", or null for a current role.
	End *resume.Month `json:"end"`
	// Dates is the span as printed ("Oct 2023 – Present"); see [Dates].
	Dates string `json:"dates"`
	// Highlights are the bullets the variant kept, in order. Never null.
	Highlights []DocumentHighlight `json:"highlights"`
}

func newDocumentRole(role *resume.Role) DocumentRole {

	out := DocumentRole{
		ID: role.ID, Org: role.Org, Title: role.Title, Location: role.Location,
		Start: role.Start, End: role.End, Dates: Dates(role.Start, role.End),
		Highlights: make([]DocumentHighlight, 0, len(role.Highlights)),
	}
	for _, h := range role.Highlights {
		out.Highlights = append(out.Highlights, DocumentHighlight{Lead: h.Lead, Text: *newText(h.Text), Tags: h.Tags})
	}

	return out
}
