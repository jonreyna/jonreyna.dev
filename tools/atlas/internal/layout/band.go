package layout

// Band is one lane drawn as a horizontal band.
type Band struct {
	// Lane is the lane's id.
	Lane string
	// Name is the label drawn in the gutter.
	Name string
	// Top is the band's upper edge.
	Top float64
	// Bottom is the band's lower edge, below its label row.
	Bottom float64
	// Tracks is the number of carrier tracks in the band.
	Tracks int
}

// trackY returns the y of the band's track i, counting from the top.
func (b *Band) trackY(i int, o *Options) float64 {
	return b.Top + o.BandPad + float64(i)*o.TrackGap
}

// knotTop returns the top of the knot bars in the band.
func (b *Band) knotTop(o *Options) float64 { return b.Top + o.BandPad/2 }

// knotBottom returns the bottom of the knot bars in the band, above its label row.
func (b *Band) knotBottom(o *Options) float64 { return b.Bottom - o.LabelRow - o.BandPad/2 }
