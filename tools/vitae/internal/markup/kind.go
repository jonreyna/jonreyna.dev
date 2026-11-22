package markup

// Kind is the type of a [Span].
type Kind int

// The span kinds.
const (
	// Text is ordinary text. Non-breaking spaces in it are kept as U+00A0.
	Text Kind = iota
	// Emph is emphasized text.
	Emph
	// Code is code.
	Code
	// Link is link text with a URL.
	Link
)

// String returns the kind's name.
func (k Kind) String() string {

	switch k {
	case Text:
		return "text"
	case Emph:
		return "emph"
	case Code:
		return "code"
	case Link:
		return "link"
	default:
		return "unknown"
	}
}
