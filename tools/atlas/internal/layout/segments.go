package layout

import (
	"math"
	"strconv"
	"strings"
)

// A branch that arrives within sameHeight of its carrier is drawn straight; num
// rounds to tenths.
const (
	sameHeight = 0.5
	tenths     = 10
)

// segments draws every carrier as stretches between its branch points, each as firm
// as the firmest link still ahead of it, and every branch as a curve from the
// carrier into its knot.
func segments(a *arrangement) []Segment {

	var out []Segment

	for _, c := range a.carriers {
		source := a.spots[c.source].person.ID
		x := c.x0

		for i, t := range c.targets {
			ahead := c.targets[i:]
			if t.branchX > x {
				out = append(out, Segment{
					Thread: c.thread, Strength: firmest(ahead), People: people(source, ahead, a),
					D: "M" + num(x) + " " + num(c.y) + "H" + num(t.branchX),
				})
				x = t.branchX
			}

			out = append(out, Segment{
				Thread: c.thread, Link: t.link.ID, Type: t.link.Type, Strength: t.link.Type.Strength(),
				People: []string{source, a.spots[t.spot].person.ID}, D: branch(t, c.y, a.o.Bend),
			})
		}
	}

	return out
}

// branch returns the path from the carrier at height y into the knot: straight when
// it arrives at the same height, otherwise an S-curve.
func branch(t carrierTarget, y, bend float64) string {

	if math.Abs(t.port-y) < sameHeight {
		return "M" + num(t.branchX) + " " + num(y) + "H" + num(t.left)
	}

	mid := num(t.branchX + bend/2)

	return "M" + num(t.branchX) + " " + num(y) +
		"C" + mid + " " + num(y) + " " + mid + " " + num(t.port) + " " + num(t.left) + " " + num(t.port)
}

// firmest returns the strength of the firmest link among targets.
func firmest(targets []carrierTarget) int {

	s := 0
	for _, t := range targets {
		s = max(s, t.link.Type.Strength())
	}

	return s
}

// people lists the source knot and the knots of targets, for a carrier stretch.
func people(source string, targets []carrierTarget, a *arrangement) []string {

	out := make([]string, 0, 1+len(targets))
	out = append(out, source)
	for _, t := range targets {
		out = append(out, a.spots[t.spot].person.ID)
	}

	return out
}

// num writes a coordinate with one decimal place at most.
func num(v float64) string {

	s := strconv.FormatFloat(math.Round(v*tenths)/tenths, 'f', 1, 64)

	return strings.TrimSuffix(s, ".0")
}
