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
