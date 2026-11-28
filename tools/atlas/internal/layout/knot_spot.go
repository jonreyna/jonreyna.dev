package layout

import "github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"

// knotSpot is a knot while the layout is worked out: its person, anchor year, lanes
// and label width, and the x it is given.
type knotSpot struct {
	person *plate.Person
	anchor float64
	lo, hi int
	width  float64
	x      float64
}

// overlaps reports whether two knots share a lane, and so need room between them.
func (k *knotSpot) overlaps(o *knotSpot) bool { return k.lo <= o.hi && o.lo <= k.hi }

// spans reports whether the knot covers lane i.
func (k *knotSpot) spans(i int) bool { return k.lo <= i && i <= k.hi }

// distance returns how many lanes lie between lane i and the knot's nearest lane.
func (k *knotSpot) distance(i int) int {

	switch {
	case i < k.lo:
		return k.lo - i
	case i > k.hi:
		return i - k.hi
	default:
		return 0
	}
}
