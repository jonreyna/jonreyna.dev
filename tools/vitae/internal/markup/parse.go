package markup

// nbsp is the non-breaking space a ~ stands for.
const nbsp = " "

// Parse splits s into spans. Adjacent text is merged into one [Text] span.
func Parse(s string) ([]Span, error) {
	p := &parser{rs: []rune(s)}
	for i := 0; i < len(p.rs); i++ {
		last, err := p.step(i)
		if err != nil {
			return nil, err
		}
		i = last
	}
	p.flush()
	return p.spans, nil
}

// MustParse is [Parse] for text already validated; it panics on malformed markup.
func MustParse(s string) []Span {
	spans, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return spans
}
