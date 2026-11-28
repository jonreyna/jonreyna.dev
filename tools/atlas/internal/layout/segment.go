package layout

import "github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"

// Segment is one SVG path of a thread: a stretch of carrier, or the branch that
// carries one link into its later knot.
type Segment struct {
	// Thread is the thread's id.
	Thread string
	// Link is the link's id for a branch, or empty for a stretch of carrier.
	Link string
	// Type is the link's type for a branch; a carrier has none.
	Type plate.LinkType
	// Strength is how firmly the segment is drawn (see plate.LinkType.Strength): a
	// carrier stretch takes the firmest link it carries.
	Strength int
	// People are the ids of the knots the segment belongs to: the earlier knot and
	// every later knot it carries a link to. The page lights these when a knot is
	// hovered.
	People []string
	// D is the SVG path data.
	D string
}
