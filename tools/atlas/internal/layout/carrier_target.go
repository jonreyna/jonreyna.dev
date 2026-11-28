package layout

import "github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"

// carrierTarget is one link a carrier delivers: the later knot, where the branch
// leaves the carrier, where it meets the knot, and the height it arrives at.
type carrierTarget struct {
	link    *plate.Link
	spot    int
	branchX float64
	left    float64
	port    float64
}
