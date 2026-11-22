package resume

// School is one education entry.
type School struct {
	// ID identifies the entry, unique within the résumé.
	ID string `yaml:"id"`
	// Org is the school's name.
	Org string `yaml:"org"`
	// Program is what was studied, as printed ("Computer Science major"). It must not
	// claim a degree that wasn't conferred.
	Program string `yaml:"program"`
	// Degree is the degree's name. Renderers print it only when Conferred is true.
	Degree string `yaml:"degree"`
	// Conferred reports whether the degree was awarded. It is a pointer so that
	// leaving it out is an error rather than a silent default.
	Conferred *bool `yaml:"conferred"`
	// Minors are the completed minors, in lowercase as they read in a sentence.
	Minors []string `yaml:"minors"`
	// Start is the first year of study.
	Start int `yaml:"start"`
	// End is the last year of study.
	End int `yaml:"end"`
	// Note is a sentence in inline markup. When the degree wasn't conferred it says so,
	// and it is the only field allowed to mention the degree.
	Note string `yaml:"note"`
	// GPA is the grade point average as printed ("3.91"), or empty.
	GPA string `yaml:"gpa"`
	// Honors are honors and societies, in print order.
	Honors []string `yaml:"honors"`
	// Courses are the courses named in the coursework line.
	Courses []Course `yaml:"courses"`
}
