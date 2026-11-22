package resume

// Course is one course named in a coursework line.
type Course struct {
	// Name is the course's name in lowercase, as it reads in a sentence.
	Name string `yaml:"name"`
	// Tags let a variant show only some courses (see [Variant.ShowCourses]).
	Tags []string `yaml:"tags"`
}
