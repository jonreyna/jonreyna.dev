package render

import (
	"bytes"
	"flag"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/lint"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

//nolint:gochecknoglobals // test flags are package-level by design
var update = flag.Bool("update", false, "rewrite the golden files")

// TestGolden renders every variant of the fixture in every format and compares the
// result byte for byte with testdata/golden. The cases come from the fixture itself
// rather than a table, so a new variant or format is covered without editing this test.
// After an intended change run `go test ./... -update` and review the diff.
func TestGolden(t *testing.T) {
	t.Parallel()
	r, err := resume.Load("../../testdata/resume.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if issues := lint.Check(r); lint.HasErrors(issues, true) {
		t.Fatalf("fixture has issues: %v", issues)
	}
	for _, name := range slices.Sorted(maps.Keys(r.Variants)) {
		for _, format := range Formats() {
			t.Run(name+"."+format, func(t *testing.T) {
				t.Parallel()
				checkGolden(t, r, name, format)
			})
		}
	}
}

// checkGolden renders one variant in one format and compares it with its golden file,
// or rewrites the file when -update is set.
func checkGolden(t *testing.T, r *resume.Resume, name, format string) {
	t.Helper()
	resolved, err := variant.Apply(r, name)
	if err != nil {
		t.Fatal(err)
	}
	renderer, err := NewRenderer(format)
	if err != nil {
		t.Fatal(err)
	}
	got, err := renderer.Render(resolved)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join("..", "..", "testdata", "golden", "resume."+name+"."+format)
	if *update {
		if err = os.WriteFile(path, got, 0o600); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%v (run go test ./... -update to create it)", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("differs from %s; got:\n%s", path, got)
	}
}
