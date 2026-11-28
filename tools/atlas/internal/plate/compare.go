package plate

// Compare is the earlier text a link's card sets beside the later one, when the
// evidence is a parallel in the wording.
type Compare struct {
	// Cite is the earlier passage's reference.
	Cite string `yaml:"cite"`
	// Quote is the earlier passage, if the card quotes it.
	Quote string `yaml:"quote"`
	// Source is the translation's key in Plate.Sources; empty means the link's own.
	Source string `yaml:"source"`
	// URL is where the earlier passage can be read.
	URL string `yaml:"url"`
}
