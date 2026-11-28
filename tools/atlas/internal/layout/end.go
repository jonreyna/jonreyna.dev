package layout

import "github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"

// End is the mark where one link arrives at its later knot: a dot, or a bar for a
// rejection.
type End struct {
	// Link is the link's id.
	Link string
	// Thread is the link's thread id.
	Thread string
	// From is the earlier person's id.
	From string
	// Type is the link's type, which sets the mark.
	Type plate.LinkType
	// X is where the mark sits, on the knot's left edge.
	X float64
	// Y is the mark's height.
	Y float64
}
