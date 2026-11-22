package render

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/filter"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/lint"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/model"
)

var update = flag.Bool("update", false, "rewrite the golden files")

// TestGolden renders every variant of the fixture and compares it byte for byte with
// testdata/golden. Run `go test ./... -update` after an intended change, and review the
// diff.
func TestGolden(t *testing.T) {
	r, err := model.Load("../../testdata/resume.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if issues := lint.Check(r); lint.HasErrors(issues, true) {
		t.Fatalf("fixture has issues: %v", issues)
	}
	for name := range r.Variants {
		t.Run(name, func(t *testing.T) {
			resolved, err := filter.Apply(r, name)
			if err != nil {
				t.Fatal(err)
			}
			got, err := JSON(resolved)
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join("../../testdata/golden", "resume."+name+".json")
			if *update {
				if err := os.WriteFile(path, got, 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("%v (run go test ./... -update to create it)", err)
			}
			if !bytes.Equal(got, want) {
				t.Errorf("%s differs from the golden file; got:\n%s", name, got)
			}
		})
	}
}

func TestDates(t *testing.T) {
	start := model.Month{Year: 2023, Month: 10}
	if got := Dates(start, nil); got != "Oct 2023 – Present" {
		t.Errorf("got %q", got)
	}
	end := model.Month{Year: 2026, Month: 12}
	if got := Dates(start, &end); got != "Oct 2023 – Dec 2026" {
		t.Errorf("got %q", got)
	}
}

// The degree must never reach the output when it wasn't conferred.
func TestUnconferredDegreeOmitted(t *testing.T) {
	no, yes := false, true
	r := &filter.Resolved{Sections: []string{"education"}, Education: []model.School{
		{Org: "U", Program: "CS major", Degree: "B.S.", Conferred: &no, Start: 1, End: 2},
	}}
	if d := Build(r); d.Education[0].Degree != "" {
		t.Errorf("unconferred degree rendered: %q", d.Education[0].Degree)
	}
	r.Education[0].Conferred = &yes
	if d := Build(r); d.Education[0].Degree != "B.S." {
		t.Errorf("conferred degree missing")
	}
}
