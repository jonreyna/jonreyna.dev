package render

import (
	_ "embed"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/markup"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// markdownTemplate is the Markdown résumé's structure.
//
//go:embed markdown.tmpl
var markdownTemplate string

// Markdown renders r as CommonMark. It is a download in its own right, and make turns
// it into the Word version with Pandoc, so its structure is chosen to map onto Word
// styles: the name is the only level-1 heading, sections are level 2, and each role or
// school is a level-3 heading followed by an italic line of place and dates.
func Markdown(r *variant.Resolved) ([]byte, error) { return markdownStyle().render(r) }

// markdownStyle is the Markdown format's [textStyle].
func markdownStyle() *textStyle {

	return &textStyle{
		name: "markdown", source: markdownTemplate, left: "{{", right: "}}",
		inline: markup.Markdown, sep: " · ",
	}
}
