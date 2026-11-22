package renderer

import "github.com/jonreyna/jonreyna.dev/tools/vitae/internal/markup"

// Text is a text field in every form a consumer needs.
type Text struct {
	// Plain is the text with markup removed and ties as ordinary spaces.
	Plain string `json:"plain"`
	// HTML is the text as escaped HTML, with ties as U+00A0.
	HTML string `json:"html"`
}

// newText converts inline markup that lint has already validated. It returns nil for
// empty input, so the field is left out of the JSON.
func newText(s string) *Text {
	if s == "" {
		return nil
	}
	spans := markup.MustParse(s)
	return &Text{Plain: markup.Plain(spans), HTML: markup.HTML(spans)}
}
