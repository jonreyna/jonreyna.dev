package markup

// Span is a run of text of one [Kind].
type Span struct {
	// Kind is what the run is.
	Kind Kind
	// Text is the run's text with markup and escapes removed. Ties are U+00A0.
	Text string
	// URL is the link target; it is set only when Kind is [Link].
	URL string
}
