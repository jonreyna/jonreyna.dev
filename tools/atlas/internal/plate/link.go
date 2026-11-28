package plate

// Link joins two knots on one thread: the later writer (To) took up, named, quoted,
// rejected or may have drawn on the earlier (From).
type Link struct {
	// ID is the link's key, such as "L05", used in anchors and the outline.
	ID string `yaml:"id"`
	// Thread is the id of the thread the link belongs to.
	Thread string `yaml:"thread"`
	// From is the id of the earlier person, the source of the idea.
	From string `yaml:"from"`
	// To is the id of the later person, whose text is the evidence.
	To string `yaml:"to"`
	// Type says how To used From.
	Type LinkType `yaml:"type"`
	// Cite is the reference to the passage in To's text (or in the work that reports
	// it).
	Cite string `yaml:"cite"`
	// Quote is the passage, a line or two from a public-domain translation.
	Quote string `yaml:"quote"`
	// Paraphrase stands in for Quote when there is nothing to quote.
	Paraphrase string `yaml:"paraphrase"`
	// Source is the translation's key in Plate.Sources.
	Source string `yaml:"source"`
	// URL is where the passage can be read.
	URL string `yaml:"url"`
	// Compare is the earlier text to set beside it, if any.
	Compare *Compare `yaml:"compare"`
	// Scholarship lists the scholars who make or dispute the case, each with title,
	// place and year.
	Scholarship []string `yaml:"scholarship"`
	// Note is a sentence or two more for the card.
	Note string `yaml:"note"`
	// Located says what kind of evidence the card shows.
	Located Located `yaml:"located"`
	// Checked is nil until Jon checks the link.
	Checked *Check `yaml:"checked"`
}

// Text returns the quote, or the paraphrase when there is no quote.
func (l *Link) Text() string {

	if l.Quote != "" {
		return l.Quote
	}

	return l.Paraphrase
}
