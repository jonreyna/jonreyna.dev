package renderer

import "github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"

// Func adapts a function to the [Renderer] interface, the way http.HandlerFunc adapts
// one to http.Handler. Read it as "renderer function": the package name supplies the
// prefix, so it doesn't repeat it.
type Func func(r *variant.Resolved) ([]byte, error)

// Render calls f(r).
func (f Func) Render(r *variant.Resolved) ([]byte, error) { return f(r) }
