package layout

import (
	"fmt"

	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"
)

// Layout is a plate placed on the chart, ready to draw.
type Layout struct {
	// Width is the chart's width.
	Width float64
	// Height is the chart's height, axes included.
	Height float64
	// Bands are the lanes, top to bottom.
	Bands []Band
	// Knots are the people, in the plate's order.
	Knots []Knot
	// Segments are the threads' paths: carrier stretches and branches.
	Segments []Segment
	// Ticks are the century marks.
	Ticks []Tick
	// Crossings counts where a branch crosses another thread's carrier.
	Crossings int
}

// New lays out a plate that [plate.Validate] has passed. It fails if the knots don't
// fit the width or a branch has no room to reach its knot.
func New(p *plate.Plate, o *Options) (*Layout, error) {

	spots, scale, err := placeKnots(p, o)
	if err != nil {
		return nil, fmt.Errorf("place knots: %w", err)
	}

	carriers, err := buildCarriers(p, spots, o)
	if err != nil {
		return nil, fmt.Errorf("build carriers: %w", err)
	}

	a := newArrangement(p, spots, carriers, o)
	a.anneal(o.Iterations, o.Seed)

	l := &Layout{
		Width: o.Width, Height: a.bands[len(a.bands)-1].Bottom + o.Axis,
		Bands: a.bands, Ticks: centuryTicks(p.Span, scale), Crossings: a.crossings(),
	}
	l.Knots = knots(a)
	l.Segments = segments(a)

	return l, nil
}

// knots turns the placed knots into bars, with an end mark for every link arriving.
func knots(a *arrangement) []Knot {

	out := make([]Knot, len(a.spots))
	for i := range a.spots {
		s := &a.spots[i]
		sp := a.spans[i]
		out[i] = Knot{
			Person: s.person.ID, Label: s.person.Short(), X: s.x,
			Top: sp.top, Bottom: sp.bottom, StemBottom: a.bands[s.hi].knotBottom(a.o),
			LabelY: a.bands[s.hi].Bottom - a.o.LabelRow/2 + a.o.LabelRow/6,
		}
	}

	for _, c := range a.carriers {
		for _, t := range c.targets {
			out[t.spot].Ends = append(out[t.spot].Ends, End{
				Link: t.link.ID, Thread: c.thread, From: t.link.From, Type: t.link.Type, X: t.left, Y: t.port,
			})
		}
	}

	return out
}
