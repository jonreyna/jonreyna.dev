package plate

// Located says what kind of evidence a link's card shows.
type Located string

// The kinds of evidence.
const (
	// LocatedPassage means the quote was copied from the linked translation at the
	// cited place.
	LocatedPassage Located = "passage"
	// LocatedParaphrase means there is no public-domain translation to quote, so the
	// card paraphrases.
	LocatedParaphrase Located = "paraphrase"
	// LocatedScholarship means the link rests on the scholarship it names.
	LocatedScholarship Located = "scholarship"
)

// Valid reports whether l is one of the kinds of evidence.
func (l Located) Valid() bool {
	return l == LocatedPassage || l == LocatedParaphrase || l == LocatedScholarship
}
