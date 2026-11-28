package plate

// LinkType says how the later writer used the earlier one, and sets how the link is
// drawn.
type LinkType string

// The link types. [LinkQuotes] and [LinkNames] are drawn as solid lines,
// [LinkAdapts] as a thinner line, [LinkRejects] as a solid line ending in a bar, and
// [LinkArgued] as a faint dotted line.
const (
	// LinkQuotes means the later writer quotes the earlier.
	LinkQuotes LinkType = "quotes"
	// LinkNames means the later writer names the earlier as a source.
	LinkNames LinkType = "names"
	// LinkAdapts means the later writer takes the idea up without naming the source;
	// the card shows the evidence (a verbal parallel or a scholar).
	LinkAdapts LinkType = "adapts"
	// LinkRejects means the later writer argues against the earlier.
	LinkRejects LinkType = "rejects"
	// LinkArgued means the link is proposed in the scholarship and disputed.
	LinkArgued LinkType = "argued"
)

// How firmly each kind of link is drawn (see Strength).
const (
	strengthArgued = 1
	strengthAdapts = 2
	strengthFirm   = 3
)

// Valid reports whether t is one of the link types.
func (t LinkType) Valid() bool {

	switch t {
	case LinkQuotes, LinkNames, LinkAdapts, LinkRejects, LinkArgued:
		return true
	default:
		return false
	}
}

// Strength ranks how firmly a link is drawn: 3 for quotes, names and rejects, 2 for
// adapts, 1 for argued, 0 for anything else. A shared line is drawn as firmly as
// the firmest link it carries.
func (t LinkType) Strength() int {

	switch t {
	case LinkQuotes, LinkNames, LinkRejects:
		return strengthFirm
	case LinkAdapts:
		return strengthAdapts
	case LinkArgued:
		return strengthArgued
	default:
		return 0
	}
}

// Verb returns the phrase the outline uses between the two names: "Justin names
// Heraclitus", "Paul may draw on the Stoics".
func (t LinkType) Verb() string {

	if t == LinkArgued {
		return "may draw on"
	}

	return string(t)
}
