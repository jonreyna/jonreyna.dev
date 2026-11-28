// Command atlas checks an atlas plate and builds it into a Play piece's markup and
// stylesheet.
//
// Usage:
//
//	atlas validate [-require-checked] PLATE.yaml
//	atlas build [-require-checked] [-max-crossings N] -o MARKUP.html -css RULES.css PLATE.yaml
//
// validate reads the plate and prints every problem it finds (see package plate),
// or one line such as "validated athens-and-jerusalem.yaml: 26 people, 8 threads,
// 55 links (55 unchecked)".
//
// build validates the plate, lays it out (package layout) and writes the piece's
// markup, everything inside <body>, and the CSS rules generated from its threads and
// people (package render). The play layout serves the markup; `make play` puts the
// rules after the piece's own stylesheet. It prints "wrote FILE" for each file and
// the number of thread crossings.
//
// Flags:
//
//	-require-checked  report every link Jon hasn't checked, as production builds must
//	-max-crossings    fail if the layout has more crossings than this (default: no limit)
//	-o                the markup file to write
//	-css              the stylesheet to write
//
// Exit codes: 0 on success, 1 if the plate has problems or a file can't be read or
// written, 2 for a usage error.
package main
