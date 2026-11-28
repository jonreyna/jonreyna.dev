package render

import (
	"html"
	"strings"
	"unicode"
)

// smart escapes prose for HTML and sets its straight quotes as curly ones: a quote
// opens at the start or after a space or opening bracket, and closes elsewhere; an
// apostrophe between letters is a right single quote.
func smart(s string) string {

	var b strings.Builder

	prev := ' '
	for _, r := range s {
		opening := unicode.IsSpace(prev) || strings.ContainsRune("([{‘“—–", prev)

		switch {
		case r == '"' && opening:
			b.WriteRune('“')
		case r == '"':
			b.WriteRune('”')
		case r == '\'' && opening:
			b.WriteRune('‘')
		case r == '\'':
			b.WriteRune('’')
		default:
			b.WriteString(html.EscapeString(string(r)))
		}

		prev = r
	}

	return b.String()
}

// attr escapes a value for an HTML attribute.
func attr(s string) string { return html.EscapeString(s) }
