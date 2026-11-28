package layout

import (
	"cmp"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"
)

// The annealing's cost weights and temperatures.
const (
	crossingCost = 10
	trackCost    = 3
	startTemp    = 8.0
	endTemp      = 0.05
	portStep     = 3.0
	// seedMix varies the generator's second seed: the golden ratio's bits.
	seedMix = 0x9e3779b97f4a7c15
)

// arrangement is the carriers' order within each lane, which decides their tracks,
// and the band geometry that follows from it.
type arrangement struct {
	o        *Options
	p        *plate.Plate
	spots    []knotSpot
	carriers []*carrier
	order    [][]int
	bands    []Band
	outgoing map[spotThread]int
}

// newArrangement starts with each lane's carriers ordered by where they head (up
// first, down last) and then by where they start.
func newArrangement(p *plate.Plate, spots []knotSpot, carriers []*carrier, o *Options) *arrangement {

	a := &arrangement{
		o: o, p: p, spots: spots, carriers: carriers,
		order: make([][]int, len(p.Lanes)), outgoing: map[spotThread]int{},
	}

	heading := make([]float64, len(carriers))
	for i, c := range carriers {
		a.outgoing[spotThread{c.source, c.thread}] = i
		a.order[c.lane] = append(a.order[c.lane], i)

		for _, t := range c.targets {
			heading[i] += float64(spots[t.spot].lo + spots[t.spot].hi - 2*c.lane)
		}

		heading[i] /= float64(len(c.targets))
	}

	for _, lane := range a.order {
		slices.SortStableFunc(lane, func(x, y int) int {
			return cmp.Or(cmp.Compare(heading[x], heading[y]), cmp.Compare(carriers[x].x0, carriers[y].x0))
		})
	}

	a.place()

	return a
}

// place packs each lane's carriers into tracks first-fit in the current order, sizes
// the bands, and sets every carrier's height and every branch's arrival height.
func (a *arrangement) place() {

	counts := make([]int, len(a.order))
	for lane, ids := range a.order {
		var ends []float64

		for _, id := range ids {
			c := a.carriers[id]

			track := slices.IndexFunc(ends, func(end float64) bool { return end+a.o.KnotWidth < c.x0 })
			if track < 0 {
				track = len(ends)
				ends = append(ends, 0)
			}

			ends[track] = c.x1
			c.track = track
		}

		counts[lane] = len(ends)
	}

	a.sizeBands(counts)

	for _, c := range a.carriers {
		c.y = a.bands[c.lane].trackY(c.track, a.o)
	}

	a.ports()
}

// sizeBands stacks the bands from the top axis down, each tall enough for its tracks.
func (a *arrangement) sizeBands(counts []int) {

	a.bands = make([]Band, len(counts))
	top := a.o.Axis

	for i, n := range counts {
		height := a.o.BandPad + a.o.BandPad + a.o.LabelRow + float64(max(n, a.o.MinTracks)-1)*a.o.TrackGap
		a.bands[i] = Band{
			Lane: a.p.Lanes[i].ID, Name: a.p.Lanes[i].Name, Top: top, Bottom: top + height, Tracks: n,
		}
		top += height
	}
}

// ports sets the height each branch arrives at: the later knot's own carrier on the
// same thread if it has one in reach, the carrier's height if the knot is in its lane,
// or else a slot at the knot's top or bottom end, stacked in arrival order.
func (a *arrangement) ports() {

	above, below := map[int]int{}, map[int]int{}
	for _, c := range a.carriers {
		for j := range c.targets {
			t := &c.targets[j]
			k := &a.spots[t.spot]

			if id, ok := a.outgoing[spotThread{t.spot, c.thread}]; ok && k.spans(a.carriers[id].lane) {
				t.port = a.carriers[id].y

				continue
			}

			switch {
			case k.spans(c.lane):
				t.port = c.y
			case c.lane < k.lo:
				t.port = a.bands[k.lo].knotTop(a.o) + portStep*float64(1+above[t.spot])
				above[t.spot]++
			default:
				t.port = a.bands[k.hi].knotBottom(a.o) - portStep*float64(1+below[t.spot])
				below[t.spot]++
			}
		}
	}
}

// crossings counts the places where a branch, on its way between carrier and knot,
// crosses another carrier.
func (a *arrangement) crossings() int {

	n := 0
	for _, c := range a.carriers {
		for _, t := range c.targets {
			lo, hi := min(c.y, t.port), max(c.y, t.port)
			if hi-lo < 1 {
				continue
			}

			x := (t.branchX + t.left) * half
			for _, d := range a.carriers {
				if d != c && d.active(x) && lo < d.y && d.y < hi {
					n++
				}
			}
		}
	}

	return n
}

// cost weighs crossings against the tracks used.
func (a *arrangement) cost() int {

	tracks := 0
	for _, b := range a.bands {
		tracks += b.Tracks
	}

	return crossingCost*a.crossings() + trackCost*tracks
}

// anneal swaps carriers within lanes, keeping a swap that lowers the cost and, with
// a probability that falls as the search cools, one that raises it. It ends on the
// best order found.
func (a *arrangement) anneal(iterations int, seed uint64) {

	lanes := []int{}
	for lane, ids := range a.order {
		if len(ids) > 1 {
			lanes = append(lanes, lane)
		}
	}

	if iterations == 0 || len(lanes) == 0 {
		return
	}

	rng := rand.New(rand.NewPCG(seed, seed^seedMix)) //nolint:gosec // G404: layout search, not security
	current := a.cost()
	best, bestOrder := current, a.copyOrder()

	for step := range iterations {
		temp := startTemp * math.Pow(endTemp/startTemp, float64(step)/float64(iterations))
		ids := a.order[lanes[rng.IntN(len(lanes))]]
		i, j := rng.IntN(len(ids)), rng.IntN(len(ids))
		ids[i], ids[j] = ids[j], ids[i]
		a.place()

		next := a.cost()
		if next <= current || rng.Float64() < math.Exp(float64(current-next)/temp) {
			current = next
			if current < best {
				best, bestOrder = current, a.copyOrder()
			}

			continue
		}

		ids[i], ids[j] = ids[j], ids[i]
	}

	a.order = bestOrder
	a.place()
}

// copyOrder returns a copy of each lane's order.
func (a *arrangement) copyOrder() [][]int {

	out := make([][]int, len(a.order))
	for i, ids := range a.order {
		out[i] = slices.Clone(ids)
	}

	return out
}
