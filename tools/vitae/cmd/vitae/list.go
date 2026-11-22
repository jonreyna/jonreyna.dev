package main

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

// list prints the variants or the tags, as args[0] asks.
func list(r *resume.Resume, args []string, w io.Writer) error {
	if len(args) != 1 {
		return errorUsage("list needs one of: variants, tags")
	}
	var b strings.Builder
	switch args[0] {
	case "variants":
		for _, name := range slices.Sorted(maps.Keys(r.Variants)) {
			fmt.Fprintf(&b, "%-10s %s\n", name, r.Variants[name].Description)
		}
	case "tags":
		counts := map[string]int{}
		for _, role := range r.Experience {
			for _, h := range role.Highlights {
				for _, t := range h.Tags {
					counts[t]++
				}
			}
		}
		for _, t := range slices.Sorted(maps.Keys(r.Tags)) {
			fmt.Fprintf(&b, "%-14s %2d  %s\n", t, counts[t], r.Tags[t])
		}
	default:
		return errorUsage("list needs one of: variants, tags")
	}
	_, err := io.WriteString(w, b.String())
	return err //nolint:wrapcheck // a failed write to stdout needs no context
}
