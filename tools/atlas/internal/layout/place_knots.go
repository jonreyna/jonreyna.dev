package layout

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"unicode/utf8"

	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"
)

// bisectSteps is how many halvings find the units per year; 60 is far below a
// hundredth of a unit.
const bisectSteps = 60

// half halves a length.
const half = 0.5

// errTooNarrow marks a chart too narrow for its knots' labels.
var errTooNarrow = errors.New("the chart is too narrow")

// placeKnots gives every knot its x and returns the knots, in the plate's order, and
// the time scale through them. It fails if the knots don't fit the width even with
// no room for time between them.
func placeKnots(p *plate.Plate, o *Options) ([]knotSpot, *timeScale, error) {

	spots := make([]knotSpot, len(p.People))
	for i := range p.People {
		person := &p.People[i]
		spots[i] = knotSpot{
			person: person, anchor: person.Anchor(p.Span),
			lo: p.LaneIndex(person.Lanes[0]), hi: p.LaneIndex(person.Lanes[len(person.Lanes)-1]),
			width: max(float64(utf8.RuneCountInString(person.Short()))*o.CharWidth+o.LabelPad, o.KnotWidth),
		}
	}

	byYear := make([]*knotSpot, len(spots))
	for i := range spots {
		byYear[i] = &spots[i]
	}

	slices.SortStableFunc(byYear, func(a, b *knotSpot) int { return cmp.Compare(a.anchor, b.anchor) })

	left, right := o.Gutter, o.Width-o.Right
	span := [2]float64{float64(p.Span[0]), float64(p.Span[1])}

	if end := spread(byYear, span, left, 0); end > right {
		return nil, nil, fmt.Errorf("%w: knots need a width of %.0f, not %.0f", errTooNarrow, end+o.Right, o.Width)
	}

	lo, hi := 0.0, (right-left)/(span[1]-span[0])
	for range bisectSteps {
		mid := (lo + hi) * half
		if spread(byYear, span, left, mid) > right {
			hi = mid
		} else {
			lo = mid
		}
	}

	end := spread(byYear, span, left, lo)

	scale := &timeScale{}
	scale.add(span[0], left)

	for _, k := range byYear {
		scale.add(k.anchor, k.x)
	}

	scale.add(span[1], end)

	return spots, scale, nil
}

// spread places the knots, in year order, at perYear units per year, pushing each
// right as far as its lanes need, and returns where the span ends.
func spread(byYear []*knotSpot, span [2]float64, left, perYear float64) float64 {

	prevYear, prevX := span[0], left
	for i, k := range byYear {
		x := prevX + perYear*(k.anchor-prevYear)
		if i == 0 {
			x = max(x, left+k.width/2)
		}

		for _, j := range byYear[:i] {
			if j.overlaps(k) {
				x = max(x, j.x+(j.width+k.width)/2)
			}
		}

		k.x = x
		prevYear, prevX = k.anchor, x
	}

	end := prevX + perYear*(span[1]-prevYear)
	if n := len(byYear); n > 0 {
		end = max(end, byYear[n-1].x+byYear[n-1].width/2)
	}

	return end
}
