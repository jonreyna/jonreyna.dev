package resume

// Role is one job.
type Role struct {
	// ID identifies the role, unique within the résumé ("equinix").
	ID string `yaml:"id"`
	// Org is the employer's name as printed.
	Org string `yaml:"org"`
	// Title is the job title as printed, including any qualifier such as "(contract)".
	Title string `yaml:"title"`
	// Location is where the job was ("Remote", "Atlanta, GA").
	Location string `yaml:"location"`
	// Start is the first month in the role.
	Start Month `yaml:"start"`
	// End is the last month in the role, or nil if the role is current.
	End *Month `yaml:"end"`
	// Highlights are the role's bullets, in print order.
	Highlights []Highlight `yaml:"highlights"`
}
