package plate

// Lane is one horizontal band of the chart, such as "Philosophers".
type Lane struct {
	// ID is the lane's key, used by Person.Lanes.
	ID string `yaml:"id"`
	// Name is the label drawn at the lane's left edge.
	Name string `yaml:"name"`
}
