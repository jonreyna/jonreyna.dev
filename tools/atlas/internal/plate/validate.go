package plate

// Validate checks a parsed plate and returns every problem it finds, in the order of
// the file: the plate's fields, then lanes, threads, people and links, then threads
// and people that nothing uses. An empty result means the plate can be drawn. With
// requireChecked, a link Jon hasn't checked is a problem too.
func Validate(p *Plate, requireChecked bool) []Issue {

	v := newValidation(p)
	v.plate()
	v.lanes()
	v.threads()

	for i := range p.People {
		v.person(&p.People[i])
	}

	for i := range p.Links {
		v.link(&p.Links[i], requireChecked)
	}

	v.unused()

	return v.issues
}
