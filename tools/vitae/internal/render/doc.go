// Package render writes a resolved résumé ([variant.Resolved]) in each output format.
//
// Each format implements [Interface], with [Func] to adapt a plain function; [For]
// returns the one for a format name and [Formats] lists them:
//
//   - "json" ([JSON]) is the [Document] that the website reads, in which every text
//     field carries both plain text and HTML ([Text]) so no consumer has to parse
//     markup.
//   - "tex" ([LaTeX]) is a LaTeX document for the jrresume class in resume/latex/. It
//     holds the content and its order; the class holds the design. make builds the PDF
//     from it with latexmk, and the class stops the build if the PDF runs past the
//     variant's max_pages.
//   - "txt" ([PlainText]) is plain text for pasting into application forms.
//   - "md" ([Markdown]) is CommonMark, a download of its own and the source make turns
//     into Word with Pandoc.
//
// LaTeX, plain text and Markdown share one template engine (textStyle): each format's
// template lays out the document, and its style says how fields are escaped and how
// the contact line and school details read.
//
// Output is deterministic: no timestamps, and keys in a fixed order. That makes
// golden-file tests exact.
//
// The degree rule holds here too: a [Document] carries a degree only when it was
// conferred, and the text formats print the program instead of a degree that
// wasn't.
package render
