// Package render writes a resolved résumé ([variant.Resolved]) in each output format.
//
// Each format is a [Renderer]; [For] returns the one for a format name and [Formats]
// lists them. So far there is one: "json", the [Document] that the website reads, in
// which every text field carries both plain text and HTML ([Text]) so no consumer has
// to parse markup. LaTeX, plain text and Markdown are next; PDF and Word are built from
// those by make, not here.
//
// Output is deterministic: no timestamps, and keys in a fixed order. That makes
// golden-file tests exact.
//
// The degree rule holds here too: a [Document] carries a degree only when it was
// conferred.
package render
