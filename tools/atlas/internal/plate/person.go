package plate

import (
	"fmt"
	"strconv"
)

// midpoint weights two years equally.
const midpoint = 0.5

// Person is one knot on the chart: a writer, a school, or a book whose author is
// unknown.
type Person struct {
	// ID is the person's key, used by links and in anchors (#k-ID).
	ID string `yaml:"id"`
	// Name is the full name, used in cards and the outline.
	Name string `yaml:"name"`
	// Label is the short name drawn on the chart; empty means Name.
	Label string `yaml:"label"`
	// Lanes are the lanes the knot spans, which must be next to each other.
	Lanes []string `yaml:"lanes"`
	// From is the first year of the person's life or activity, or of a book's
	// composition window. BC years are negative.
	From int `yaml:"from"`
	// To is the last such year.
	To int `yaml:"to"`
	// At, when not zero, places the knot at this year instead of the middle of From
	// and To, for a book whose window is long.
	At int `yaml:"at"`
	// Circa marks From and To as approximate.
	Circa bool `yaml:"circa"`
	// Dating says whose dating From and To follow, for a book or a contested date.
	Dating string `yaml:"dating"`
	// Members lists the people a group knot stands for, with their dates.
	Members string `yaml:"members"`
	// Card is the sentence or two the knot's card opens with.
	Card string `yaml:"card"`
	// Wiki is the English Wikipedia article's title, with underscores.
	Wiki string `yaml:"wiki"`
}

// Short returns the label drawn on the chart.
func (p *Person) Short() string {

	if p.Label != "" {
		return p.Label
	}

	return p.Name
}

// Anchor returns the year the knot is drawn at: At if set, otherwise the middle of
// the person's years that fall inside the span.
func (p *Person) Anchor(span [2]int) float64 {

	if p.At != 0 {
		return float64(p.At)
	}

	from, to := max(p.From, span[0]), min(p.To, span[1])

	return float64(from+to) * midpoint
}

// Dates returns the person's years as the card shows them, such as "c. 428–348 BC",
// "c. 20 BC – AD 50" or "354–430".
func (p *Person) Dates() string {

	prefix := ""
	if p.Circa {
		prefix = "c. "
	}

	switch {
	case p.To < 0:
		return fmt.Sprintf("%s%d–%d BC", prefix, -p.From, -p.To)
	case p.From < 0:
		return fmt.Sprintf("%s%d BC – AD %d", prefix, -p.From, p.To)
	default:
		return prefix + strconv.Itoa(p.From) + "–" + strconv.Itoa(p.To)
	}
}
