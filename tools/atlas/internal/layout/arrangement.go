package layout

import (
	"cmp"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"
)

// The annealing's cost weights, temperatures and knot-bar geometry.
const (
	crossingCost   = 10
	passBehindCost = 2
	trackCost      = 3
	startTemp      = 8.0
	endTemp        = 0.05
	portStep       = 8.0
	knotPad        = 7.0
	minKnotHeight  = 20.0
	minLaneOverlap = 8.0
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
	minSlots []int
	bands    []Band
	spans    []knotSpan
	outgoing map[spotThread]int
}

// newArrangement starts with each lane's carriers ordered by where they head (up
// first, down last) and then by where they start.
func newArrangement(p *plate.Plate, spots []knotSpot, carriers []*carrier, o *Options) *arrangement {

	a := &arrangement{
		o: o, p: p, spots: spots, carriers: carriers,
		order: make([][]int, len(p.Lanes)), minSlots: make([]int, len(p.Lanes)),
		outgoing: map[spotThread]int{},
	}

	incoming := make([]int, len(spots))
	heading := make([]float64, len(carriers))

	for i, c := range carriers {
		a.outgoing[spotThread{c.source, c.thread}] = i
		a.order[c.lane] = append(a.order[c.lane], i)

		for _, t := range c.targets {
			incoming[t.spot]++
			heading[i] += float64(spots[t.spot].lo + spots[t.spot].hi - 2*c.lane)
		}

		heading[i] /= float64(len(c.targets))
	}

	for i, n := range incoming {
		if spots[i].lo == spots[i].hi && n > a.minSlots[spots[i].lo] {
			a.minSlots[spots[i].lo] = n
		}
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
// the bands, and sets every carrier's height, branch arrival height and knot span.
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
	a.spans = a.knotSpans()
}

// sizeBands stacks the bands from the top axis down, each tall enough for its tracks
// and for the busiest single-lane knot's arriving ports.
func (a *arrangement) sizeBands(counts []int) {

	a.bands = make([]Band, len(counts))
	top := a.o.Axis

	for i, n := range counts {
		slots := max(n, a.minSlots[i], a.o.MinTracks)
		height := a.o.BandPad + a.o.BandPad + a.o.LabelRow + float64(slots-1)*a.o.TrackGap
		a.bands[i] = Band{
			Lane: a.p.Lanes[i].ID, Name: a.p.Lanes[i].Name, Top: top, Bottom: top + height, Tracks: n,
		}
		top += height
	}
}

// ports sets the height each branch arrives at, ordering a knot's incoming links by
// their preferred height and spreading them apart so endpoints never collide.
func (a *arrangement) ports() {

	bySpot := make([][]incomingPort, len(a.spots))
	for _, c := range a.carriers {
		for j := range c.targets {
			t := &c.targets[j]
			want := a.preferredPort(c, t)
			bySpot[t.spot] = append(bySpot[t.spot], incomingPort{target: t, want: want, fromY: c.y, fromX: c.x0})
		}
	}

	for spot, list := range bySpot {
		if len(list) == 0 {
			continue
		}

		slices.SortStableFunc(list, func(x, y incomingPort) int {
			return cmp.Or(cmp.Compare(x.want, y.want), cmp.Compare(x.fromY, y.fromY), cmp.Compare(x.fromX, y.fromX))
		})

		k := &a.spots[spot]
		spreadPorts(list, a.bands[k.lo].knotTop(a.o), a.bands[k.hi].knotBottom(a.o))
	}
}

// preferredPort picks a branch's natural arrival height before collision spreading.
func (a *arrangement) preferredPort(c *carrier, t *carrierTarget) float64 {

	k := &a.spots[t.spot]
	if id, ok := a.outgoing[spotThread{t.spot, c.thread}]; ok && k.spans(a.carriers[id].lane) {
		return a.carriers[id].y
	}

	switch {
	case k.spans(c.lane):
		return c.y
	case c.lane < k.lo:
		return a.bands[k.lo].knotTop(a.o) + portStep
	default:
		return a.bands[k.hi].knotBottom(a.o) - portStep
	}
}

// spreadPorts places sorted incoming ports inside [lo, hi], separated by at least
// portStep (or evenly when the knot has more ports than full steps fit).
func spreadPorts(list []incomingPort, lo, hi float64) {

	n := len(list)
	step := portStep

	if n > 1 && float64(n-1)*step > hi-lo {
		step = (hi - lo) / float64(n-1)
	}

	ys := make([]float64, n)
	for i, p := range list {
		ys[i] = min(max(p.want, lo), hi)
	}

	for i := 1; i < n; i++ {
		if ys[i] < ys[i-1]+step {
			ys[i] = ys[i-1] + step
		}
	}

	last := n - 1
	if ys[last] > hi {
		ys[last] = hi
	}

	for i := last - 1; i >= 0; i-- {
		if ys[i] > ys[i+1]-step {
			ys[i] = ys[i+1] - step
		}
	}

	for i, p := range list {
		p.target.port = ys[i]
	}
}

// knotSpans trims each knot's bar to the vertical range of the ports and carriers
// that touch it, while keeping multi-lane knots spanning across their lane borders.
func (a *arrangement) knotSpans() []knotSpan {

	minY := make([]float64, len(a.spots))
	maxY := make([]float64, len(a.spots))

	for i := range a.spots {
		minY[i], maxY[i] = math.Inf(1), math.Inf(-1)
	}

	for _, c := range a.carriers {
		minY[c.source] = min(minY[c.source], c.y)
		maxY[c.source] = max(maxY[c.source], c.y)
		for _, t := range c.targets {
			minY[t.spot] = min(minY[t.spot], t.port)
			maxY[t.spot] = max(maxY[t.spot], t.port)
		}
	}

	out := make([]knotSpan, len(a.spots))
	for i := range a.spots {
		out[i] = a.clampSpan(&a.spots[i], minY[i], maxY[i])
	}

	return out
}

// clampSpan pads one knot's active [loY, hiY] range and clamps it inside its lanes.
func (a *arrangement) clampSpan(s *knotSpot, loY, hiY float64) knotSpan {

	bandTop := a.bands[s.lo].knotTop(a.o)
	bandBot := a.bands[s.hi].knotBottom(a.o)

	if math.IsInf(loY, 1) {
		return knotSpan{top: bandTop, bottom: bandBot}
	}

	top := max(bandTop, loY-knotPad)
	bot := min(bandBot, hiY+knotPad)

	if s.lo < s.hi {
		top = min(top, a.bands[s.lo].knotBottom(a.o)-minLaneOverlap)
		bot = max(bot, a.bands[s.hi].knotTop(a.o)+minLaneOverlap)
	}

	if bot-top < minKnotHeight {
		mid := (top + bot) * half
		top = max(bandTop, mid-minKnotHeight*half)
		bot = min(bandBot, top+minKnotHeight)
		top = max(bandTop, bot-minKnotHeight)
	}

	return knotSpan{top: top, bottom: bot}
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

// passBehinds counts carriers passing behind the active bar of an unrelated knot.
func (a *arrangement) passBehinds() int {

	n := 0
	for _, c := range a.carriers {
		for i := range a.spots {
			if i == c.source || !c.active(a.spots[i].x) || c.reaches(i) {
				continue
			}

			if a.spans[i].top <= c.y && c.y <= a.spans[i].bottom {
				n++
			}
		}
	}

	return n
}

// cost weighs crossings, pass-behinds and the tracks used.
func (a *arrangement) cost() int {

	tracks := 0
	for _, b := range a.bands {
		tracks += b.Tracks
	}

	return crossingCost*a.crossings() + passBehindCost*a.passBehinds() + trackCost*tracks
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
