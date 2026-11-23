package render

import (
	"fmt"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// Interface writes a resolved résumé in one format. It is named Interface, not
// Renderer, so that callers write render.Interface rather than the stutter
// render.Renderer, the way package sort names sort.Interface.
type Interface interface {
	// Render returns the résumé in the implementation's format.
	Render(r *variant.Resolved) ([]byte, error)
}

// Formats returns the names of the formats [For] knows, in a fixed order.
func Formats() []string { return []string{"json", "tex", "txt", "md"} }

// For returns the [Interface] for the named format. The error wraps
// [ErrUnknownFormat].
//
//nolint:ireturn // For chooses among implementations by name; returning the interface is its purpose.
func For(format string) (Interface, error) {

	switch format {
	case "json":
		return Func(JSON), nil
	case "tex":
		return Func(LaTeX), nil
	case "txt":
		return Func(PlainText), nil
	case "md":
		return Func(Markdown), nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownFormat, format)
	}
}
