package plate

// Source is a translation that quotes are copied from.
type Source struct {
	// Name is the translation and its date, as the bibliography lists it.
	Name string `yaml:"name"`
	// Via is the site the text was read on, such as "New Advent".
	Via string `yaml:"via"`
}
