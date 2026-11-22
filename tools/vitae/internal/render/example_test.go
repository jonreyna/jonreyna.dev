package render_test

import (
	"fmt"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/render"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

func ExampleDates() {
	start := resume.Month{Year: 2023, Month: 10}
	end := resume.Month{Year: 2026, Month: 12}
	fmt.Println(render.Dates(start, nil))
	fmt.Println(render.Dates(start, &end))
	// Output:
	// Oct 2023 – Present
	// Oct 2023 – Dec 2026
}

func ExampleFor() {
	r := &variant.Resolved{
		Variant:  "full",
		Sections: []string{resume.SectionSummary},
		Basics:   resume.Basics{Name: "Ada Example", Email: "ada@example.com"},
		Summary:  "Builds *distributed* systems.",
	}
	renderer, err := render.For("json")
	if err != nil {
		panic(err)
	}
	out, err := renderer.Render(r)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(out))
	// Output:
	// {
	//   "schema": 1,
	//   "variant": "full",
	//   "sections": [
	//     "summary"
	//   ],
	//   "basics": {
	//     "name": "Ada Example",
	//     "label": "",
	//     "location": "",
	//     "remote": false,
	//     "email": "ada@example.com",
	//     "links": []
	//   },
	//   "summary": {
	//     "plain": "Builds distributed systems.",
	//     "html": "Builds <em>distributed</em> systems."
	//   }
	// }
}
