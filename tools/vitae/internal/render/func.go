package render

import "github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"

// Func adapts a function to [Interface], the way http.HandlerFunc adapts one to
// http.Handler.
type Func func(r *variant.Resolved) ([]byte, error)

// Render calls f(r).
func (f Func) Render(r *variant.Resolved) ([]byte, error) { return f(r) }
