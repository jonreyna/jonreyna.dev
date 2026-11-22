package render

import (
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

func TestDates(t *testing.T) {

	t.Parallel()
	oct2023 := resume.Month{Year: 2023, Month: 10}
	tests := []struct {
		name  string
		start resume.Month
		end   *resume.Month
		want  string
	}{
		{name: "current", start: oct2023, end: nil, want: "Oct 2023 – Present"},
		{name: "ended", start: oct2023, end: &resume.Month{Year: 2026, Month: 12}, want: "Oct 2023 – Dec 2026"},
		{name: "september", start: resume.Month{Year: 2013, Month: 9}, end: &resume.Month{Year: 2015, Month: 1}, want: "Sep 2013 – Jan 2015"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			if got := Dates(tc.start, tc.end); got != tc.want {
				t.Errorf("Dates = %q, want %q", got, tc.want)
			}
		})
	}
}
