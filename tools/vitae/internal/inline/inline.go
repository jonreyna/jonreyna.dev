// Package inline parses the small markup allowed in résumé text fields and writes it out
// for each format:
//
//	*emphasis*   `code`   [text](https://url)   ~ (non-breaking space)
//
// A backslash makes the next character literal (\*, \`, \[, \~, \\). Emphasis, code and
// link text don't nest.
package inline

import (
	"fmt"
	"html"
	"strings"
)

// Kind is the type of a span.
type Kind int

// Span kinds.
const (
	Text Kind = iota
	Emph
	Code
	Link
)

// Span is a run of text with one kind. Ties are kept in Text as U+00A0.
type Span struct {
	Kind Kind
	Text string
	URL  string // Link only
}

const nbsp = " "

// Parse splits s into spans, or reports the first malformed piece of markup.
func Parse(s string) ([]Span, error) {
	var spans []Span
	var buf strings.Builder
	flush := func() {
		if buf.Len() > 0 {
			spans = append(spans, Span{Kind: Text, Text: buf.String()})
			buf.Reset()
		}
	}
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		switch c := rs[i]; c {
		case '\\':
			if i+1 == len(rs) {
				return nil, fmt.Errorf("trailing backslash")
			}
			i++
			buf.WriteRune(rs[i])
		case '~':
			buf.WriteString(nbsp)
		case '*', '`':
			end := indexRune(rs, i+1, c)
			if end < 0 {
				return nil, fmt.Errorf("unclosed %c at %q", c, excerpt(rs, i))
			}
			if end == i+1 {
				return nil, fmt.Errorf("empty %c%c", c, c)
			}
			flush()
			kind := Emph
			if c == '`' {
				kind = Code
			}
			inner := string(rs[i+1 : end])
			if kind == Emph {
				inner = strings.ReplaceAll(inner, "~", nbsp)
			}
			spans = append(spans, Span{Kind: kind, Text: inner})
			i = end
		case '[':
			close := indexRune(rs, i+1, ']')
			if close < 0 || close+1 >= len(rs) || rs[close+1] != '(' {
				return nil, fmt.Errorf("a [ must start a [text](url) link, at %q", excerpt(rs, i))
			}
			paren := indexRune(rs, close+2, ')')
			if paren < 0 {
				return nil, fmt.Errorf("unclosed link URL at %q", excerpt(rs, i))
			}
			url := string(rs[close+2 : paren])
			if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "mailto:") {
				return nil, fmt.Errorf("link URL %q must start with https:// or mailto:", url)
			}
			flush()
			spans = append(spans, Span{Kind: Link, Text: strings.ReplaceAll(string(rs[i+1:close]), "~", nbsp), URL: url})
			i = paren
		default:
			buf.WriteRune(c)
		}
	}
	flush()
	return spans, nil
}

func indexRune(rs []rune, from int, r rune) int {
	for j := from; j < len(rs); j++ {
		if rs[j] == '\\' {
			j++
			continue
		}
		if rs[j] == r {
			return j
		}
	}
	return -1
}

func excerpt(rs []rune, i int) string {
	end := min(i+20, len(rs))
	return string(rs[i:end])
}

// Plain writes text with markup removed and ties as ordinary spaces, for plain-text
// output and for measuring length.
func Plain(spans []Span) string {
	var b strings.Builder
	for _, s := range spans {
		b.WriteString(s.Text)
	}
	return strings.ReplaceAll(b.String(), nbsp, " ")
}

// HTML writes escaped HTML. Ties stay as U+00A0 characters.
func HTML(spans []Span) string {
	var b strings.Builder
	for _, s := range spans {
		t := html.EscapeString(s.Text)
		switch s.Kind {
		case Text:
			b.WriteString(t)
		case Emph:
			b.WriteString("<em>" + t + "</em>")
		case Code:
			b.WriteString("<code>" + t + "</code>")
		case Link:
			b.WriteString(`<a href="` + html.EscapeString(s.URL) + `">` + t + "</a>")
		}
	}
	return b.String()
}

// MustParse is Parse for text that has already been validated.
func MustParse(s string) []Span {
	spans, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return spans
}
