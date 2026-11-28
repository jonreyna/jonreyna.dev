package layout

// incomingPort pairs an arriving branch with its desired arrival height and source
// height so a knot's ports can be ordered top-to-bottom without collisions.
type incomingPort struct {
	target *carrierTarget
	want   float64
	fromY  float64
	fromX  float64
}
