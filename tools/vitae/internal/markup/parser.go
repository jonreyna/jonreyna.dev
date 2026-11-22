package markup

import (
	"fmt"
	"strings"
)

// excerptLen is how much of the text an error quotes.
const excerptLen = 20

// parser holds the state of one [Parse] call.
type parser struct {
	rs    []rune
	spans []Span
	buf   strings.Builder
}

// step consumes the markup that starts at rs[i] and returns the index of its last rune.
func (p *parser) step(i int) (int, error) {
	switch c := p.rs[i]; c {
	case '\\':
		if i+1 == len(p.rs) {
			return 0, ErrTrailingBackslash
		}
		p.buf.WriteRune(p.rs[i+1])
		return i + 1, nil
	case '~':
		p.buf.WriteString(nbsp)
		return i, nil
	case '*', '`':
		return p.delimited(i, c)
	case '[':
		return p.link(i)
	default:
		p.buf.WriteRune(c)
		return i, nil
	}
}

// delimited reads *emphasis* or `code` starting at rs[i].
func (p *parser) delimited(i int, delim rune) (int, error) {
	end := indexRune(p.rs, i+1, delim)
	switch {
	case end < 0:
		return 0, fmt.Errorf("%w: %c at %q", ErrUnclosed, delim, excerpt(p.rs, i))
	case end == i+1:
		return 0, fmt.Errorf("%w: %c%c", ErrEmpty, delim, delim)
	}
	p.flush()
	if delim == '`' {
		p.spans = append(p.spans, Span{Kind: Code, Text: unescape(p.rs[i+1:end], false)})
	} else {
		p.spans = append(p.spans, Span{Kind: Emph, Text: unescape(p.rs[i+1:end], true)})
	}
	return end, nil
}

// link reads [text](url) starting at rs[i].
func (p *parser) link(i int) (int, error) {
	closeBracket := indexRune(p.rs, i+1, ']')
	if closeBracket < 0 || closeBracket+1 >= len(p.rs) || p.rs[closeBracket+1] != '(' {
		return 0, fmt.Errorf("%w: [ without (url) at %q", ErrBadLink, excerpt(p.rs, i))
	}
	urlStart := closeBracket + len("](")
	closeParen := indexRune(p.rs, urlStart, ')')
	if closeParen < 0 {
		return 0, fmt.Errorf("%w: unclosed URL at %q", ErrBadLink, excerpt(p.rs, i))
	}
	url := string(p.rs[urlStart:closeParen])
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "mailto:") {
		return 0, fmt.Errorf("%w: URL %q is neither https nor mailto", ErrBadLink, url)
	}
	p.flush()
	p.spans = append(p.spans, Span{Kind: Link, Text: unescape(p.rs[i+1:closeBracket], true), URL: url})
	return closeParen, nil
}

// flush ends the current text span, if any.
func (p *parser) flush() {
	if p.buf.Len() > 0 {
		p.spans = append(p.spans, Span{Kind: Text, Text: p.buf.String()})
		p.buf.Reset()
	}
}

// indexRune returns the index of the first unescaped r at or after from, or -1.
func indexRune(rs []rune, from int, r rune) int {
	for j := from; j < len(rs); j++ {
		switch rs[j] {
		case '\\':
			j++
		case r:
			return j
		}
	}
	return -1
}

// unescape removes backslash escapes and, when ties is true, turns ~ into U+00A0.
func unescape(rs []rune, ties bool) string {
	var b strings.Builder
	for j := 0; j < len(rs); j++ {
		switch {
		case rs[j] == '\\' && j+1 < len(rs):
			j++
			b.WriteRune(rs[j])
		case rs[j] == '~' && ties:
			b.WriteString(nbsp)
		default:
			b.WriteRune(rs[j])
		}
	}
	return b.String()
}

// excerpt returns up to excerptLen runes of rs from i, for error messages.
func excerpt(rs []rune, i int) string {
	return string(rs[i:min(i+excerptLen, len(rs))])
}
