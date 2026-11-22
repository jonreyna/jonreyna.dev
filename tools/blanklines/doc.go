// Command blanklines checks, and with -w adds, the blank lines that rules 13 and 14 of
// the go-best-practices skill require:
//
//   - Rule 13: a function or method body longer than three lines starts with a blank
//     line after the signature's opening brace. Function literals count too.
//   - Rule 14: the closing brace of a block (a function body, or an if, for, switch or
//     select body) is followed by a blank line, unless the next line closes something
//     too (it starts with }, ) or ]). A line such as "} else {" opens another block and
//     is left alone. Composite literals and type definitions aren't blocks.
//
// No standard linter checks these rules, and gofumpt removes the first kind of blank
// line, so format with gofmt and run this alongside golangci-lint.
//
// Usage:
//
//	blanklines [-w] [path ...]
//
// Each path is a Go file or a directory, searched recursively; the default is the
// current directory. Directories named testdata or vendor, hidden directories and
// generated files are skipped. Without -w, each missing blank line is printed as
// "file:line: message". With -w, the blank lines are added and gofmt checks the result
// before the file is written.
//
// Exit codes: 0 nothing missing (or everything fixed), 1 blank lines missing or a file
// couldn't be read, parsed or written, 2 usage error.
package main
