package lint_test

import (
	"fmt"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/lint"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

func ExampleCheck() {

	no := false
	r := &resume.Resume{
		Schema:     1,
		Basics:     resume.Basics{Name: "Ada", Email: "ada@example.com"},
		Summary:    map[string]string{"default": "Builds systems."},
		Tags:       map[string]string{"go": "Go"},
		Experience: []resume.Role{{ID: "acme", Org: "Acme", Title: "Engineer", Start: resume.Month{Year: 2020, Month: 1}, Highlights: []resume.Highlight{{Text: "Wrote Go.", Tags: []string{"go"}, Weight: 50}}}},
		Education:  []resume.School{{Org: "U", Program: "B.S. in Physics", Conferred: &no, Note: "One course short.", Start: 2010, End: 2014}},
		Variants:   map[string]resume.Variant{"full": {Sections: []string{resume.SectionExperience}}},
	}
	issues := lint.Check(r)
	for _, i := range issues {
		fmt.Println(i)
	}

	fmt.Println("fails:", lint.HasErrors(issues, false))
	// Output:
	// error: education[0].program: claims a degree that was not conferred: "B.S. in Physics"
	// fails: true
}
