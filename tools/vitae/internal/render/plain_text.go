package render

import (
	_ "embed"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/markup"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// plainTextTemplate is the plain-text résumé's structure.
//
//go:embed plain_text.tmpl
var plainTextTemplate string

// PlainText renders r as plain text for pasting into application forms: no markup,
// one line per paragraph or bullet (forms rewrap text themselves), section heads in
// capitals and "- " bullets.
func PlainText(r *variant.Resolved) ([]byte, error) { return plainTextStyle().render(r) }

// plainTextStyle is the plain-text format's [textStyle].
func plainTextStyle() *textStyle {

	return &textStyle{
		name: "plain text", source: plainTextTemplate, left: "{{", right: "}}",
		inline: markup.Plain, sep: " · ",
	}
}
