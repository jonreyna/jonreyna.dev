package render

import "github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"

// DocumentBasics is the header in a [Document].
type DocumentBasics struct {
	// Name is the person's name.
	Name string `json:"name"`
	// Label is the headline under the name, after any variant override.
	Label string `json:"label"`
	// Location is where the person lives, as printed.
	Location string `json:"location"`
	// Remote reports whether the person works remotely.
	Remote bool `json:"remote"`
	// Email is the public contact address.
	Email string `json:"email"`
	// Phone is present only in private renders (vitae render --private).
	Phone string `json:"phone,omitempty"`
	// Links are profile and site links, in print order. Never null.
	Links []resume.Link `json:"links"`
}
