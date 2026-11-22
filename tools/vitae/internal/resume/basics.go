package resume

// Basics is the résumé's header.
type Basics struct {
	// Name is the person's name as it heads the résumé.
	Name string `yaml:"name"`
	// Label is the headline under the name, such as a job title. A variant may
	// replace it.
	Label string `yaml:"label"`
	// Location is where the person lives, as printed ("Tampa Bay, FL").
	Location string `yaml:"location"`
	// Remote reports whether the person works remotely; renderers print "Remote".
	Remote bool `yaml:"remote"`
	// Email is the public contact address.
	Email string `yaml:"email"`
	// Phone is set only from the private overlay (see [ApplyPrivate]); the public file
	// must leave it empty.
	Phone string `yaml:"phone"`
	// Links are profile and site links, in print order.
	Links []Link `yaml:"links"`
}
