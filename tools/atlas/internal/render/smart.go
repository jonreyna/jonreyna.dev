package render

import (
	"html"
	"strings"
	"unicode"
)

// smart escapes prose for HTML, sets straight quotes as curly ones, and turns
// balanced *emphasis* spans into <em> elements (used for book titles and terms).
func smart(s string) string {

	rs := []rune(s)

	var b strings.Builder

	inEm := false

	for i, r := range rs {
		prev := ' '
		if i > 0 {
			prev = rs[i-1]
		}

		opening := unicode.IsSpace(prev) || strings.ContainsRune("([{‘“'\"—–*", prev)

		if r == '*' {
			inEm = writeStar(&b, rs, i, prev, opening, inEm)

			continue
		}

		b.WriteString(quoteRune(r, opening))
	}

	return b.String()
}

// writeStar turns a balanced asterisk into <em> or </em>, or leaves a stray one literal.
func writeStar(b *strings.Builder, rs []rune, i int, prev rune, opening, inEm bool) bool {

	if inEm && !unicode.IsSpace(prev) {
		b.WriteString("</em>")

		return false
	}

	if !inEm && opening && hasClosingStar(rs, i) {
		b.WriteString("<em>")

		return true
	}

	b.WriteByte('*')

	return inEm
}

// quoteRune sets straight quotes as curly ones and escapes other HTML characters.
func quoteRune(r rune, opening bool) string {

	switch r {
	case '"':
		if opening {
			return "“"
		}

		return "”"
	case '\'':
		if opening {
			return "‘"
		}

		return "’"
	default:
		return html.EscapeString(string(r))
	}
}

// hasClosingStar reports whether an opening asterisk at start has a matching close.
func hasClosingStar(rs []rune, start int) bool {

	first := start + 1
	if first >= len(rs) || unicode.IsSpace(rs[first]) || rs[first] == '*' {
		return false
	}

	for j := first + 1; j < len(rs); j++ {
		if rs[j] != '*' || unicode.IsSpace(rs[j-1]) {
			continue
		}

		if j+1 == len(rs) || unicode.IsSpace(rs[j+1]) || strings.ContainsRune(".,;:!?)]}’”'\"—–", rs[j+1]) {
			return true
		}
	}

	return false
}

// attr escapes a value for an HTML attribute.
func attr(s string) string { return html.EscapeString(s) }
