// Package layout places a plate's knots and threads on a woven chart, in SVG user
// units.
//
// [New] works in four steps:
//
//  1. Time scale. Each knot sits at its anchor year (see plate.Person.Anchor). The
//     scale is piecewise linear: as many units per year as fit the width, stretched
//     wherever two knots in a lane would otherwise overlap, so crowded centuries get
//     more room. The century ticks show where it stretches.
//  2. Carriers. For each thread and each knot with links on it, one carrier runs
//     right from the knot, in the lane nearest its targets, and a branch leaves it
//     for each target. A branch that ends at a knot that passes the idea on lands on
//     that knot's own carrier, so lineages read as one continuous thread.
//  3. Tracks. Carriers in a lane are packed into tracks first-fit, in an order found
//     by simulated annealing with a fixed seed, to keep down the number of places
//     where a branch crosses a carrier and the number of tracks. The same plate and
//     options always give the same layout.
//  4. Geometry. Lanes become bands sized to their tracks, knots become bars spanning
//     their bands, and every carrier and branch becomes an SVG path.
//
// The result is a [Layout]: [Band]s, [Knot]s with their [End] marks, path
// [Segment]s and century [Tick]s, plus [Layout.Crossings], the count the annealing
// minimized, which a build can hold to a budget. Package render draws it.
package layout
