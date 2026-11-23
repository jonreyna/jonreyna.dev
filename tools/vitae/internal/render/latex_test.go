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

func TestLaTeXMinors(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name   string
		minors []string
		want   string
	}{
		{name: "none", minors: nil, want: ""},
		{name: "one", minors: []string{"mathematics"}, want: "Minor in mathematics"},
		{name: "two", minors: []string{"mathematics", "statistics"}, want: "Minors in mathematics and statistics"},
		{name: "three, no serial comma", minors: []string{"a", "b", "c"}, want: "Minors in a, b and c"},
		{name: "escaped", minors: []string{"R&D"}, want: `Minor in R\&D`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			if got := latexMinors(tc.minors); got != tc.want {
				t.Errorf("latexMinors = %q, want %q", got, tc.want)
			}
		})
	}
}
