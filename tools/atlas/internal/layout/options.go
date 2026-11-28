package layout

// Options sets the chart's dimensions, in SVG user units, and the search.
type Options struct {
	// Width is the whole chart's width.
	Width float64
	// Gutter is the room at the left for the lane names.
	Gutter float64
	// Right is the margin at the right.
	Right float64
	// Axis is the height of the century row above and below the bands.
	Axis float64
	// BandPad is the space inside a band above its first track and below its last.
	BandPad float64
	// TrackGap is the distance between tracks.
	TrackGap float64
	// LabelRow is the height of the row at a band's foot that holds knot labels.
	LabelRow float64
	// MinTracks is the fewest tracks a band is sized for, so a quiet lane's knots
	// still stand tall enough to read.
	MinTracks int
	// KnotWidth is a knot bar's width.
	KnotWidth float64
	// Bend is how far before a knot a same-lane branch leaves its carrier.
	Bend float64
	// CharWidth estimates the width of one character of a knot label.
	CharWidth float64
	// LabelPad is the space added between labels.
	LabelPad float64
	// Iterations is the number of annealing steps; 0 keeps the first packing.
	Iterations int
	// Seed seeds the annealing, so a layout can be reproduced.
	Seed uint64
}

// Defaults returns the options the site's plates are drawn with.
func Defaults() Options {

	//nolint:mnd // these numbers are the chart's chosen dimensions
	return Options{
		Width: 1840, Gutter: 150, Right: 40, Axis: 28, BandPad: 18, TrackGap: 10.5,
		LabelRow: 26, MinTracks: 6, KnotWidth: 5, Bend: 36, CharWidth: 8.2, LabelPad: 20,
		Iterations: 60000, Seed: 1,
	}
}

// halfKnot returns half a knot bar's width, the distance from its center to an edge.
func (o *Options) halfKnot() float64 { return o.KnotWidth * half }
