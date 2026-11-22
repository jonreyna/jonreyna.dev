package render

import "github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"

// RendererFunc adapts a function to the [Renderer] interface, the way http.HandlerFunc
// adapts one to http.Handler.
type RendererFunc func(r *variant.Resolved) ([]byte, error)

// Render calls f(r).
func (f RendererFunc) Render(r *variant.Resolved) ([]byte, error) { return f(r) }
