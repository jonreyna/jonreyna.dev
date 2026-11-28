package plate

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"go.yaml.in/yaml/v3"
)

// ErrEmpty is returned by [Parse] for a file with no plate in it.
var ErrEmpty = errors.New("empty plate")

// Plate is one woven chart and everything drawn on it.
type Plate struct {
	// ID names the plate and its piece directory, such as "athens-and-jerusalem".
	ID string `yaml:"plate"`
	// Title is the page's heading.
	Title string `yaml:"title"`
	// Subtitle is the line under the title.
	Subtitle string `yaml:"subtitle"`
	// Scope says, on the page, what the plate covers and leaves out.
	Scope string `yaml:"scope"`
	// Span is the first and last year on the time axis.
	Span [2]int `yaml:"span"`
	// Sources are the translations quotes are copied from, by key.
	Sources map[string]Source `yaml:"sources"`
	// Lanes are the chart's bands, top to bottom.
	Lanes []Lane `yaml:"lanes"`
	// Threads are the ideas, in legend order.
	Threads []Thread `yaml:"threads"`
	// People are the knots, in the order their cards appear.
	People []Person `yaml:"people"`
	// Links are the citations that join the knots, in outline order.
	Links []Link `yaml:"links"`

	people map[string]int
}

// Parse decodes a plate from YAML. Unknown fields are errors, so a misspelling can't
// silently drop data. It does not check the plate's content; see [Validate].
func Parse(data []byte) (*Plate, error) {

	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)

	var p Plate

	err := dec.Decode(&p)
	if errors.Is(err, io.EOF) {
		return nil, ErrEmpty
	}

	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	p.people = make(map[string]int, len(p.People))
	for i := range p.People {
		p.people[p.People[i].ID] = i
	}

	return &p, nil
}

// Person returns the person with the given id, or nil if there is none.
func (p *Plate) Person(id string) *Person {

	i, ok := p.people[id]
	if !ok {
		return nil
	}

	return &p.People[i]
}

// Thread returns the thread with the given id, or nil if there is none.
func (p *Plate) Thread(id string) *Thread {

	for i := range p.Threads {
		if p.Threads[i].ID == id {
			return &p.Threads[i]
		}
	}

	return nil
}

// LaneIndex returns the position of the lane with the given id, top first, or -1.
func (p *Plate) LaneIndex(id string) int {

	for i := range p.Lanes {
		if p.Lanes[i].ID == id {
			return i
		}
	}

	return -1
}
