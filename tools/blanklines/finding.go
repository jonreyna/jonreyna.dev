package main

import "fmt"

// finding is one missing blank line: it belongs after line in file.
type finding struct {
	file string
	line int
	msg  string
}

// String formats the finding as "file:line: message", the form editors link to.
func (f finding) String() string { return fmt.Sprintf("%s:%d: %s", f.file, f.line, f.msg) }
