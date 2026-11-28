package layout

// Tick is one century mark on the time axis.
type Tick struct {
	// Year is the year marked, BC negative.
	Year int
	// X is where it falls.
	X float64
	// Label is its text, such as "300 BC" or "AD 200", or empty where labels would
	// crowd.
	Label string
}
