package lint

import (
	"strings"
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/model"
)

func load(t *testing.T, path string) *model.Resume {
	t.Helper()
	r, err := model.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// The real résumé must pass strict validation; CI runs this too.
func TestRealResume(t *testing.T) {
	r := load(t, "../../../../resume/resume.yaml")
	for _, i := range Check(r) {
		t.Error(i)
	}
}

func TestFixtureIsClean(t *testing.T) {
	for _, i := range Check(load(t, "../../testdata/resume.yaml")) {
		t.Error(i)
	}
}

// Each case breaks the fixture one way and expects an issue mentioning want.
func TestRules(t *testing.T) {
	cases := []struct {
		name   string
		break_ func(r *model.Resume)
		want   string
		sev    Severity
	}{
		{"schema", func(r *model.Resume) { r.Schema = 2 }, "must be 1", Error},
		{"phone in public file", func(r *model.Resume) { r.Basics.Phone = "555" }, "resume.private.yaml", Error},
		{"insecure link", func(r *model.Resume) { r.Basics.Links[0].URL = "http://x" }, "https://", Error},
		{"no default summary", func(r *model.Resume) { delete(r.Summary, "default") }, "summary.default", Error},
		{"unknown tag", func(r *model.Resume) { r.Experience[0].Highlights[0].Tags = []string{"rust"} }, `unknown tag "rust"`, Error},
		{"two current roles", func(r *model.Resume) { r.Experience[1].End = nil }, "only the first role", Error},
		{"out of order", func(r *model.Resume) { r.Experience[1].Start = model.Month{Year: 2025, Month: 1} }, "newest first", Error},
		{"ends before start", func(r *model.Resume) { r.Experience[1].End = &model.Month{Year: 2014, Month: 1} }, "before it starts", Error},
		{"duplicate id", func(r *model.Resume) { r.Experience[1].ID = "now" }, "duplicate id", Error},
		{"bad weight", func(r *model.Resume) { r.Experience[0].Highlights[0].Weight = 0 }, "weight", Error},
		{"bad markup", func(r *model.Resume) { r.Experience[0].Highlights[0].Text = "*open" }, "markup", Error},
		{"conferred unstated", func(r *model.Resume) { r.Education[0].Conferred = nil }, "conferred must be stated", Error},
		{"degree claimed", func(r *model.Resume) { r.Education[0].Program = "B.S. in Computer Science" }, "not conferred", Error},
		{"degree in honors", func(r *model.Resume) { r.Education[0].Honors = []string{"Graduated cum laude"} }, "not conferred", Error},
		{"no full variant", func(r *model.Resume) { delete(r.Variants, "full") }, "named full", Error},
		{"unknown section", func(r *model.Resume) {
			v := r.Variants["full"]
			v.Sections = append(v.Sections, "hobbies")
			r.Variants["full"] = v
		}, `unknown section "hobbies"`, Error},
		{"variant selects nothing", func(r *model.Resume) {
			v := r.Variants["focused"]
			v.MinWeight = 1000
			r.Variants["focused"] = v
		}, "selects no highlights", Error},
		{"long highlight", func(r *model.Resume) {
			r.Experience[0].Highlights[0].Text = strings.Repeat("word ", 140)
		}, "characters", Warning},
		{"straight quote", func(r *model.Resume) { r.Earlier = "it's" }, "straight quote", Warning},
		{"hyphen as dash", func(r *model.Resume) { r.Earlier = "1999 - 2009" }, "dash", Warning},
		{"unused tag", func(r *model.Resume) { r.Tags["unused"] = "x" }, "never used", Warning},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := load(t, "../../testdata/resume.yaml")
			tc.break_(r)
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
