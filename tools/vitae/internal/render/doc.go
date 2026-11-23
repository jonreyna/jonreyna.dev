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
//
// Plain text and Markdown are next; Word is built from them by make, not here.
//
// Output is deterministic: no timestamps, and keys in a fixed order. That makes
// golden-file tests exact.
//
// The degree rule holds here too: a [Document] carries a degree only when it was
// conferred, and the LaTeX prints the program instead of a degree that wasn't.
package render
