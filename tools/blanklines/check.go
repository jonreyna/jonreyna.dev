package main

import (
	"go/ast"
	"go/token"
	"strings"
)

// minBodyLines is the shortest body, in lines, that must start with a blank line.
const minBodyLines = 4

// Messages for the two rules.
const (
	msgSignature = "add a blank line after the function signature (rule 13)"
	msgClose     = "add a blank line after the closing brace (rule 14)"
)

// check returns the missing blank lines in a parsed file. lines is the file's source
// split on newlines, so lines[n-1] is line n.
func check(fset *token.FileSet, file string, f *ast.File, lines []string) []finding {

	var out []finding
	ast.Inspect(f, func(n ast.Node) bool {

		switch n := n.(type) {
		case *ast.FuncDecl:
			if n.Body != nil && signatureNeedsBlank(fset, n.Body, lines) {
				out = append(out, finding{file: file, line: fset.Position(n.Body.Lbrace).Line, msg: msgSignature})
			}

		case *ast.FuncLit:
			if signatureNeedsBlank(fset, n.Body, lines) {
				out = append(out, finding{file: file, line: fset.Position(n.Body.Lbrace).Line, msg: msgSignature})
			}

		case *ast.BlockStmt:
			if closeNeedsBlank(fset, n, lines) {
				out = append(out, finding{file: file, line: fset.Position(n.Rbrace).Line, msg: msgClose})
			}
		}

		return true
	})

	return out
}

// signatureNeedsBlank reports whether body is long enough to need a blank first line
// and doesn't have one.
func signatureNeedsBlank(fset *token.FileSet, body *ast.BlockStmt, lines []string) bool {

	open, closing := fset.Position(body.Lbrace).Line, fset.Position(body.Rbrace).Line
	inner := closing - open - 1
	if inner <= 0 {
		return false
	}

	if blank(lines[open]) {
		return false
	}

	return inner >= minBodyLines
}

// closeNeedsBlank reports whether the block's closing brace needs a blank line after
// it and doesn't have one.
func closeNeedsBlank(fset *token.FileSet, b *ast.BlockStmt, lines []string) bool {

	rbrace := fset.Position(b.Rbrace)
	if rbrace.Line == fset.Position(b.Lbrace).Line || rbrace.Line >= len(lines) {
		return false
	}

	// "} else {" and the like open the next block on the same line.
	if rest := lines[rbrace.Line-1][rbrace.Column:]; strings.Contains(rest, "{") {
		return false
	}

	next := strings.TrimSpace(lines[rbrace.Line])

	return next != "" && !strings.HasPrefix(next, "}") && !strings.HasPrefix(next, ")") &&
		!strings.HasPrefix(next, "]")
}

// blank reports whether a line holds only whitespace.
func blank(line string) bool { return strings.TrimSpace(line) == "" }
