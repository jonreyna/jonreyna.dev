package layout

// knotSpan is one knot's active vertical extent, trimmed to the ports and carriers
// that actually touch it while keeping multi-lane knots across their lane boundaries.
type knotSpan struct {
	top    float64
	bottom float64
}
