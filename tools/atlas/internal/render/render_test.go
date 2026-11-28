package render

import (
	"bytes"
	"flag"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/layout"
	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"
)

// update rewrites the golden files instead of comparing with them.
//
//nolint:gochecknoglobals // test flags are package-level by design
var update = flag.Bool("update", false, "rewrite testdata/*.golden")

// build renders the shared test plate.
func build(t *testing.T) (string, string) {

	t.Helper()

	data, err := os.ReadFile("../../testdata/minimal.yaml")
	if err != nil {
		t.Fatal(err)
	}

	p, err := plate.Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	o := layout.Defaults()

	l, err := layout.New(p, &o)
	if err != nil {
		t.Fatal(err)
	}

	var h, c bytes.Buffer
	if err = HTML(&h, p, l); err != nil {
		t.Fatal(err)
	}

	if err = CSS(&c, p); err != nil {
		t.Fatal(err)
	}

	return h.String(), c.String()
}

// TestGolden compares the markup and rules with the golden files: one end-to-end
// render, so a table would only obscure it. Run with -update after a deliberate change.
func TestGolden(t *testing.T) {

	t.Parallel()

	markup, rules := build(t)
	for path, got := range map[string]string{
		"testdata/minimal.html.golden": markup, "testdata/minimal.css.golden": rules,
	} {
		if *update {
			if err := os.WriteFile(path, []byte(got), 0o600); err != nil {
				t.Fatal(err)
			}

			continue
		}

		want, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		if got != string(want) {
			t.Errorf("%s differs; run go test ./internal/render -update and review the diff", path)
		}
	}
}

func TestHTMLRules(t *testing.T) {

	t.Parallel()

	markup, _ := build(t)
	ids := map[string]bool{}

	for _, m := range regexp.MustCompile(`\bid="([^"]+)"`).FindAllStringSubmatch(markup, -1) {
		ids[m[1]] = true
	}

	tests := []struct {
		name string
		ok   func() bool
	}{
		{name: "one h1", ok: func() bool { return strings.Count(markup, "<h1>") == 1 }},
		{name: "no style attributes, which the security policy blocks", ok: func() bool {
			return !strings.Contains(markup, " style=")
		}},
		{name: "no scripts", ok: func() bool { return !strings.Contains(markup, "<script") }},
		{name: "every in-page link has a target", ok: func() bool {

			for _, m := range regexp.MustCompile(`href="#([^"]+)"`).FindAllStringSubmatch(markup, -1) {
				if !ids[m[1]] {
					return false
				}
			}

			return true
		}},
		{name: "unchecked link is marked", ok: func() bool { return strings.Count(markup, "Not yet checked") == 1 }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			if !tc.ok() {
				t.Fail()
			}
		})
	}
}

func TestSmart(t *testing.T) {

	t.Parallel()

	tests := []struct {
		in, want string
	}{
		{`"In the beginning"`, "“In the beginning”"},
		{`Cal's card`, "Cal’s card"},
		{`as ('Moses') said`, "as (‘Moses’) said"},
		{`a < b & "c"`, "a &lt; b &amp; “c”"},
	}

	for _, tc := range tests {
		if got := smart(tc.in); got != tc.want {
			t.Errorf("smart(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
