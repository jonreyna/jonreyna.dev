package layout

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"
)

// errNoRoom marks a later knot drawn too close to its earlier one for a branch.
var errNoRoom = errors.New("no room for a branch")

// buildCarriers makes one carrier for each thread and earlier knot, in the order the
// links first name them, in the lane nearest its targets. It fails if a later knot is
// drawn too close to, or left of, its earlier one for a branch to reach it.
func buildCarriers(p *plate.Plate, spots []knotSpot, o *Options) ([]*carrier, error) {

	index := make(map[string]int, len(spots))
	for i := range spots {
		index[spots[i].person.ID] = i
	}

	byKey := map[[2]string]*carrier{}

	var carriers []*carrier

	for i := range p.Links {
		l := &p.Links[i]
		key := [2]string{l.Thread, l.From}

		c := byKey[key]
		if c == nil {
			s := index[l.From]
			c = &carrier{thread: l.Thread, source: s, x0: spots[s].x + o.halfKnot()}
			byKey[key] = c
			carriers = append(carriers, c)
		}

		t := index[l.To]
		left := spots[t].x - o.halfKnot()

		if left-o.Bend < c.x0 {
			return nil, fmt.Errorf("link %s: %w: %s is drawn too close to %s", l.ID, errNoRoom, l.To, l.From)
		}

		c.targets = append(c.targets, carrierTarget{link: l, spot: t, branchX: left - o.Bend, left: left})
	}

	for _, c := range carriers {
		slices.SortStableFunc(c.targets, func(a, b carrierTarget) int { return cmp.Compare(a.branchX, b.branchX) })
		c.x1 = c.targets[len(c.targets)-1].branchX
		c.lane = nearestLane(&spots[c.source], c.targets, spots)
	}

	return carriers, nil
}

// nearestLane picks, among the source knot's lanes, the one closest in total to the
// lanes of the knots its carrier reaches, preferring the upper on a tie.
func nearestLane(source *knotSpot, targets []carrierTarget, spots []knotSpot) int {

	best, bestCost := source.lo, -1
	for lane := source.lo; lane <= source.hi; lane++ {
		cost := 0
		for _, t := range targets {
			cost += spots[t.spot].distance(lane)
		}

		if bestCost < 0 || cost < bestCost {
			best, bestCost = lane, cost
		}
	}

	return best
}
