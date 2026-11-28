package plate

// Thread is one idea, drawn as one color across the chart.
type Thread struct {
	// ID is the thread's key, used by Link.Thread and in CSS classes (t-ID).
	ID string `yaml:"id"`
	// Name is the thread's label in the legend and outline.
	Name string `yaml:"name"`
	// Color names the pigment the stylesheet defines for it, such as "gold"
	// (the CSS custom property --gold).
	Color string `yaml:"color"`
	// Summary is the one sentence the legend and outline give for the thread.
	Summary string `yaml:"summary"`
}
