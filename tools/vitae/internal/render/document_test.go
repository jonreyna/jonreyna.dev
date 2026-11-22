package render

import (
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// The degree must reach the output only when it was conferred.
func TestNewDocumentDegree(t *testing.T) {
	t.Parallel()
	yes, no := true, false
	tests := []struct {
		name      string
		conferred *bool
		want      string
	}{
		{name: "conferred", conferred: &yes, want: "B.S."},
		{name: "not conferred", conferred: &no, want: ""},
		{name: "unstated", conferred: nil, want: ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := &variant.Resolved{
				Sections: []string{resume.SectionEducation},
				Education: []resume.School{
					{Org: "U", Program: "CS major", Degree: "B.S.", Conferred: tc.conferred, Start: 2010, End: 2014},
				},
			}
			if got := NewDocument(r).Education[0].Degree; got != tc.want {
				t.Errorf("degree = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNewDocumentHighlightsNeverNull(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		highlights []resume.Highlight
		want       int
	}{
		{name: "none", highlights: nil, want: 0},
		{name: "one", highlights: []resume.Highlight{{Text: "Did it.", Weight: 50}}, want: 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := &variant.Resolved{Experience: []resume.Role{{ID: "x", Highlights: tc.highlights}}}
			got := NewDocument(r).Experience[0].Highlights
			if got == nil || len(got) != tc.want {
				t.Errorf("highlights = %#v, want %d (non-nil)", got, tc.want)
			}
		})
	}
}
