package main

import (
	"fmt"
	"go/format"
	"slices"
	"strings"
)

// fix inserts a blank line after each finding's line and returns the gofmt-formatted
// source. Findings may repeat a line; each line gets one blank line.
func fix(src []byte, findings []finding) ([]byte, error) {

	after := make([]int, 0, len(findings))
	for _, f := range findings {
		after = append(after, f.line)
	}

	slices.Sort(after)
	after = slices.Compact(after)

	lines := strings.Split(string(src), "\n")
	// Insert from the bottom up, so earlier line numbers stay valid.
	for _, line := range slices.Backward(after) {
		lines = slices.Insert(lines, line, "")
	}

	out, err := format.Source([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return nil, fmt.Errorf("format fixed source: %w", err)
	}

	return out, nil
}
