package resume

// The sections a variant can list in [Variant.Sections].
const (
	// SectionSummary is the summary paragraph.
	SectionSummary = "summary"
	// SectionSkills is the skills block.
	SectionSkills = "skills"
	// SectionExperience is the list of roles.
	SectionExperience = "experience"
	// SectionEarlier is the one-line note on earlier work.
	SectionEarlier = "earlier"
	// SectionEducation is the education block.
	SectionEducation = "education"
)

// Sections returns every section name a variant may list.
func Sections() []string {
	return []string{SectionSummary, SectionSkills, SectionExperience, SectionEarlier, SectionEducation}
}

// Variant selects sections and highlights for one version of the résumé. It never
// rewords anything.
type Variant struct {
	// Description says what the variant is for; `vitae list variants` prints it.
	Description string `yaml:"description"`
	// Sections are the sections to show, in print order (see [Sections]).
	Sections []string `yaml:"sections"`
	// Summary is the key into [Resume.Summary] to use. Empty means "default".
	Summary string `yaml:"summary"`
	// Label replaces [Basics.Label] when set.
	Label string `yaml:"label"`
	// Since drops roles that ended before this month. Nil keeps every role.
	Since *Month `yaml:"since"`
	// MaxRoles keeps at most this many roles, newest first. Zero means no limit.
	MaxRoles int `yaml:"max_roles"`
	// MaxHighlights keeps at most this many highlights per role, the highest scoring.
	// Zero means no limit.
	MaxHighlights int `yaml:"max_highlights"`
	// MinWeight drops highlights scoring below it.
	MinWeight int `yaml:"min_weight"`
	// Boost lists tags whose highlights score higher.
	Boost []string `yaml:"boost"`
	// BoostBy is what each boosted tag adds to a highlight's score. Zero means the
	// default, 20.
	BoostBy int `yaml:"boost_by"`
	// Exclude lists tags whose highlights are dropped.
	Exclude []string `yaml:"exclude"`
	// ShowCourses lists the course tags to show. Empty shows every course.
	ShowCourses []string `yaml:"show_courses"`
}
