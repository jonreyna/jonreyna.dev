package render

import (
	"fmt"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// Renderer writes a resolved résumé in one format.
type Renderer interface {
	// Render returns the résumé in the renderer's format.
	Render(r *variant.Resolved) ([]byte, error)
}

// Formats returns the names of the formats [NewRenderer] knows, in a fixed order.
func Formats() []string { return []string{"json"} }

// NewRenderer returns the renderer for the named format. The error wraps [ErrUnknownFormat].
//
//nolint:ireturn // NewRenderer chooses among implementations by name; returning the interface is its purpose.
func NewRenderer(format string) (Renderer, error) {
	switch format {
	case "json":
		return RendererFunc(JSON), nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownFormat, format)
	}
}
