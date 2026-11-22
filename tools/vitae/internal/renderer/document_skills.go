package renderer

// DocumentSkills is one labeled skill line in a [Document].
type DocumentSkills struct {
	// Group is the line's label.
	Group string `json:"group"`
	// Items are the skills, in print order.
	Items []string `json:"items"`
}
