package markup

import (
	"html"
	"strings"
)

// Plain writes spans as plain text: markup removed and ties as ordinary spaces. Plain
// text output and length checks use it.
func Plain(spans []Span) string {

	var b strings.Builder
	for _, s := range spans {
		b.WriteString(s.Text)
	}

	return strings.ReplaceAll(b.String(), nbsp, " ")
}

// HTML writes spans as escaped HTML. Ties stay as U+00A0 characters.
func HTML(spans []Span) string {

	var b strings.Builder
	for _, s := range spans {
		text := html.EscapeString(s.Text)
		switch s.Kind {
		case Text:
			b.WriteString(text)
		case Emph:
			b.WriteString("<em>" + text + "</em>")
		case Code:
			b.WriteString("<code>" + text + "</code>")
		case Link:
			b.WriteString(`<a href="` + html.EscapeString(s.URL) + `">` + text + "</a>")
		}
	}

	return b.String()
}

// LaTeX writes spans as LaTeX for pdfLaTeX with inputenc's utf8: special characters
// escaped, ties as ~, emphasis as \emph, code as \texttt and links as hyperref's \href.
func LaTeX(spans []Span) string {

	var b strings.Builder
	for _, s := range spans {
		text := escapeLaTeX(s.Text)
		switch s.Kind {
		case Text:
			b.WriteString(text)
		case Emph:
			b.WriteString(`\emph{` + text + "}")
		case Code:
			b.WriteString(`\texttt{` + text + "}")
		case Link:
			b.WriteString(`\href{` + escapeLaTeXURL(s.URL) + "}{" + text + "}")
		}
	}

	return b.String()
}

// escapeLaTeX escapes the characters LaTeX treats as commands and turns U+00A0 into a
// tie (~). Other Unicode passes through for inputenc.
func escapeLaTeX(s string) string {

	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\textbackslash{}`)
		case '{', '}', '#', '$', '%', '&', '_':
			b.WriteRune('\\')
			b.WriteRune(r)
		case '~':
			b.WriteString(`\textasciitilde{}`)
		case '^':
			b.WriteString(`\textasciicircum{}`)
		case '\u00a0':
			b.WriteRune('~')
		default:
			b.WriteRune(r)
		}
	}

	return b.String()
}

// escapeLaTeXURL escapes the two characters hyperref's \href needs escaped in a URL.
func escapeLaTeXURL(url string) string {
	return strings.NewReplacer("%", `\%`, "#", `\#`).Replace(url)
}

// Markdown writes spans as CommonMark: ASCII punctuation that Markdown (or Pandoc's
// extensions) would read as markup is backslash-escaped, ties stay as U+00A0, emphasis
// is *…*, code is `…` and links are [text](url).
func Markdown(spans []Span) string {

	var b strings.Builder
	for _, s := range spans {
		switch s.Kind {
		case Text:
			b.WriteString(escapeMarkdown(s.Text))
		case Emph:
			b.WriteString("*" + escapeMarkdown(s.Text) + "*")
		case Code:
			fence := "`"
			if strings.Contains(s.Text, "`") {
				fence = "``"
			}

			b.WriteString(fence + s.Text + fence)
		case Link:
			b.WriteString("[" + escapeMarkdown(s.Text) + "](" + escapeMarkdownURL(s.URL) + ")")
		}
	}

	return b.String()
}

// escapeMarkdown backslash-escapes the characters that start markup in CommonMark or
// in Pandoc's Markdown: emphasis, code, links, HTML, math, sub- and superscript, and
// table pipes.
func escapeMarkdown(s string) string {

	var b strings.Builder
	for _, r := range s {
		if strings.ContainsRune("\\`*_[]<>$^~|", r) {
			b.WriteRune('\\')
		}

		b.WriteRune(r)
	}

	return b.String()
}

// escapeMarkdownURL percent-encodes the characters that would end a link destination.
func escapeMarkdownURL(url string) string {
	return strings.NewReplacer(" ", "%20", "(", "%28", ")", "%29", "<", "%3C", ">", "%3E").Replace(url)
}
