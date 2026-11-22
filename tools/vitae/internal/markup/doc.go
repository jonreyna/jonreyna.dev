// Package markup parses the small inline markup allowed in résumé text fields and
// writes it out for each format.
//
// The syntax is deliberately tiny:
//
//	*emphasis*            emphasized text
//	`code`                code, printed in a monospace font where the format has one
//	[text](https://url)   a link; the URL must be https:// or mailto:
//	~                     a non-breaking space, as in TeX ("RFC~8342")
//	\*  \`  \[  \~  \\    the character itself
//
// Emphasis, code and link text don't nest. [Parse] turns a string into [Span]s and
// reports malformed markup with errors that wrap the sentinels in errors.go. [Plain]
// and [HTML] write spans out; other formats (LaTeX, Markdown) add their own writers
// here as they're built.
package markup
