package render

import (
	_ "embed"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/markup"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// latexTemplate is the LaTeX document's structure. Its delimiters are << and >>,
// since LaTeX is full of braces.
//
//go:embed latex.tmpl
var latexTemplate string

// LaTeX renders r as a LaTeX document for the jrresume class
// (resume/latex/jrresume.cls), which make turns into a PDF with latexmk and pdfLaTeX.
// The document holds the content and its order; the class holds every design decision.
func LaTeX(r *variant.Resolved) ([]byte, error) { return latexStyle().render(r) }

// latexStyle is the LaTeX format's [textStyle].
func latexStyle() *textStyle {

	return &textStyle{
		name: "latex", source: latexTemplate, left: "<<", right: ">>",
		inline: markup.LaTeX, sep: `\sep{}`,
	}
}
