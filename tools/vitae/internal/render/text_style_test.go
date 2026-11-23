package render

import (
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

func TestSchoolProgram(t *testing.T) {

	t.Parallel()
	yes, no := true, false
	tests := []struct {
		name      string
		conferred *bool
		want      string
	}{
		{name: "conferred shows the degree", conferred: &yes, want: "B.S. in Computer Science"},
		{name: "not conferred shows the program", conferred: &no, want: "Computer Science major"},
		{name: "unstated shows the program", conferred: nil, want: "Computer Science major"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			s := &resume.School{Program: "Computer Science major", Degree: "B.S. in Computer Science", Conferred: tc.conferred}
			if got := schoolProgram(s); got != tc.want {
				t.Errorf("schoolProgram = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestMinors(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name   string
		style  *textStyle
		minors []string
		want   string
	}{
		{name: "none", style: plainTextStyle(), minors: nil, want: ""},
		{name: "one", style: plainTextStyle(), minors: []string{"mathematics"}, want: "Minor in mathematics"},
		{name: "two", style: plainTextStyle(), minors: []string{"math", "stats"}, want: "Minors in math and stats"},
		{name: "three, no serial comma", style: plainTextStyle(), minors: []string{"a", "b", "c"}, want: "Minors in a, b and c"},
		{name: "escaped for LaTeX", style: latexStyle(), minors: []string{"R&D"}, want: `Minor in R\&D`},
		{name: "escaped for Markdown", style: markdownStyle(), minors: []string{"C*"}, want: `Minor in C\*`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			if got := tc.style.minors(tc.minors); got != tc.want {
				t.Errorf("minors = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSchoolDetails(t *testing.T) {

	t.Parallel()
	school := &resume.School{
		Note: "Two courses short of the~degree.", GPA: "3.9", Honors: []string{"Alpha Chi"},
		Courses: []resume.Course{{Name: "linear algebra"}, {Name: "pottery"}},
	}
	tests := []struct {
		name  string
		style *textStyle
		want  string
	}{
		{name: "plain text", style: plainTextStyle(), want: "Two courses short of the degree. GPA 3.9; Alpha Chi. Coursework included linear algebra and pottery."},
		{name: "latex ties", style: latexStyle(), want: "Two courses short of the~degree. GPA~3.9; Alpha Chi. Coursework included linear algebra and pottery."},
		{name: "markdown keeps U+00A0", style: markdownStyle(), want: "Two courses short of the\u00a0degree. GPA\u00a03.9; Alpha Chi. Coursework included linear algebra and pottery."},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			if got := tc.style.schoolDetails(school); got != tc.want {
				t.Errorf("schoolDetails = %q, want %q", got, tc.want)
			}
		})
	}
}
