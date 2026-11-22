package render

import (
	"strconv"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

// DocumentEducation is one school in a [Document].
type DocumentEducation struct {
	// Org is the school's name.
	Org string `json:"org"`
	// Program is what was studied, as printed.
	Program string `json:"program"`
	// Degree is present only when the degree was conferred.
	Degree string `json:"degree,omitempty"`
	// Minors are the completed minors.
	Minors []string `json:"minors,omitempty"`
	// Dates is the span of study as printed ("2010 – 2014").
	Dates string `json:"dates"`
	// Note is the sentence about the program, such as why no degree was conferred.
	Note *Text `json:"note,omitempty"`
	// GPA is the grade point average as printed, if given.
	GPA string `json:"gpa,omitempty"`
	// Honors are honors and societies.
	Honors []string `json:"honors,omitempty"`
	// Courses are the course names the variant kept.
	Courses []string `json:"courses,omitempty"`
}

func newDocumentEducation(s *resume.School) DocumentEducation {
	e := DocumentEducation{
		Org: s.Org, Program: s.Program, Minors: s.Minors,
		Dates: strconv.Itoa(s.Start) + " – " + strconv.Itoa(s.End),
		Note:  newText(s.Note), GPA: s.GPA, Honors: s.Honors,
	}
	if s.Conferred != nil && *s.Conferred {
		e.Degree = s.Degree
	}
	for _, c := range s.Courses {
		e.Courses = append(e.Courses, c.Name)
	}
	return e
}
