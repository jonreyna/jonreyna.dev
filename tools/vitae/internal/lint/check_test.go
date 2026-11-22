package lint

import (
	"strings"
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

const fixture = "../../testdata/resume.yaml"

func load(t *testing.T, path string) *resume.Resume {

	t.Helper()
	r, err := resume.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	return r
}

// The fixture and the real résumé must both pass strict validation; CI relies on the
// second.
func TestCleanFiles(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name, path string
	}{
		{name: "fixture", path: fixture},
		{name: "real résumé", path: "../../../../resume/resume.yaml"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			for _, i := range Check(load(t, tc.path)) {
				t.Error(i)
			}
		})
	}
}

// Each case breaks the fixture one way and expects an issue mentioning want.
func TestRules(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name    string
		breakIt func(r *resume.Resume)
		want    string
		sev     Severity
	}{
		{"schema", func(r *resume.Resume) { r.Schema = 2 }, "must be 1", Error},
		{"phone in public file", func(r *resume.Resume) { r.Basics.Phone = "555" }, "resume.private.yaml", Error},
		{"insecure link", func(r *resume.Resume) { r.Basics.Links[0].URL = "http://x" }, "https://", Error},
		{"no default summary", func(r *resume.Resume) { delete(r.Summary, "default") }, "summary.default", Error},
		{"unknown tag", func(r *resume.Resume) { r.Experience[0].Highlights[0].Tags = []string{"rust"} }, `unknown tag "rust"`, Error},
		{"two current roles", func(r *resume.Resume) { r.Experience[1].End = nil }, "only the first role", Error},
		{"out of order", func(r *resume.Resume) { r.Experience[1].Start = resume.Month{Year: 2025, Month: 1} }, "newest first", Error},
		{"ends before start", func(r *resume.Resume) { r.Experience[1].End = &resume.Month{Year: 2014, Month: 1} }, "before it starts", Error},
		{"duplicate id", func(r *resume.Resume) { r.Experience[1].ID = "now" }, "duplicate id", Error},
		{"bad weight", func(r *resume.Resume) { r.Experience[0].Highlights[0].Weight = 0 }, "weight", Error},
		{"bad markup", func(r *resume.Resume) { r.Experience[0].Highlights[0].Text = "*open" }, "markup", Error},
		{"conferred unstated", func(r *resume.Resume) { r.Education[0].Conferred = nil }, "conferred must be stated", Error},
		{"degree in program", func(r *resume.Resume) { r.Education[0].Program = "B.S. in Computer Science" }, "not conferred", Error},
		{"degree in honors", func(r *resume.Resume) { r.Education[0].Honors = []string{"Graduated cum laude"} }, "not conferred", Error},
		{"no full variant", func(r *resume.Resume) { delete(r.Variants, "full") }, "named full", Error},
		{"unknown section", func(r *resume.Resume) {
			v := r.Variants["full"]
			v.Sections = append(v.Sections, "hobbies")
			r.Variants["full"] = v
		}, `unknown section "hobbies"`, Error},

		{"variant selects nothing", func(r *resume.Resume) {
			v := r.Variants["focused"]
			v.MinWeight = 1000
			r.Variants["focused"] = v
		}, "selects no highlights", Error},

		{"long highlight", func(r *resume.Resume) {
			r.Experience[0].Highlights[0].Text = strings.Repeat("word ", 140)
		}, "characters", Warning},

		{"straight quote", func(r *resume.Resume) { r.Earlier = "it's" }, "straight quote", Warning},
		{"hyphen as dash", func(r *resume.Resume) { r.Earlier = "1999 - 2009" }, "dash", Warning},
		{"unused tag", func(r *resume.Resume) { r.Tags["unused"] = "x" }, "never used", Warning},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			r := load(t, fixture)
			tc.breakIt(r)
			issues := Check(r)
			for _, i := range issues {
				if strings.Contains(i.String(), tc.want) && i.Severity == tc.sev {
					return
				}
			}

			t.Errorf("no %v mentioning %q; got %v", tc.sev, tc.want, issues)
		})
	}
}

func TestHasErrors(t *testing.T) {

	t.Parallel()
	warning := Issue{Severity: Warning}
	errIssue := Issue{Severity: Error}
	tests := []struct {
		name   string
		issues []Issue
		strict bool
		want   bool
	}{
		{name: "none", issues: nil, want: false},
		{name: "warning only", issues: []Issue{warning}, want: false},
		{name: "warning when strict", issues: []Issue{warning}, strict: true, want: true},
		{name: "error", issues: []Issue{warning, errIssue}, want: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			if got := HasErrors(tc.issues, tc.strict); got != tc.want {
				t.Errorf("HasErrors = %v, want %v", got, tc.want)
			}
		})
	}
}
