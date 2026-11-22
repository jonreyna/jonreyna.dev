package resume

// SkillGroup is one labeled line of skills.
type SkillGroup struct {
	// Group is the line's label ("Languages").
	Group string `yaml:"group"`
	// Items are the skills on the line, in print order, in inline markup.
	Items []string `yaml:"items"`
}
