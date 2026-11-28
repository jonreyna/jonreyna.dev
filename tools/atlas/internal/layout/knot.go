package layout

// Knot is one person drawn as a vertical bar with its label below.
type Knot struct {
	// Person is the person's id.
	Person string
	// Label is the text drawn under the bar.
	Label string
	// X is the bar's center.
	X float64
	// Top is the bar's upper end.
	Top float64
	// Bottom is the bar's lower end.
	Bottom float64
	// StemBottom is the foot of the knot's band above its label, where a quiet drop
	// stem lands when Bottom sits above it.
	StemBottom float64
	// LabelY is the label's baseline.
	LabelY float64
	// Ends are the marks where links arrive at the bar.
	Ends []End
}
