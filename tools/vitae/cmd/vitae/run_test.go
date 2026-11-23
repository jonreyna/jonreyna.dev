package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fixture = "../../testdata/resume.yaml"

// vitae runs a command and returns its combined output and exit code.
func vitae(t *testing.T, args ...string) (string, int) {

	t.Helper()
	var out bytes.Buffer
	err := run(args, &out, &out)
	if err == nil {
		return out.String(), 0
	}

	if ee, ok := errors.AsType[*errorExit](err); ok {
		return out.String() + ee.msg, ee.code
	}

	return out.String() + err.Error(), exitFailure
}

// fixtureWith writes the fixture with one replacement to a temporary file.
func fixtureWith(t *testing.T, old, replacement string) string {

	t.Helper()
	data, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}

	return writeFile(t, "resume.yaml", strings.Replace(string(data), old, replacement, 1))
}

func writeFile(t *testing.T, name, data string) string {

	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}

	return path
}

func readFile(t *testing.T, path string) string {

	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}

func TestRun(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name     string
		args     func(t *testing.T) []string
		wantCode int
		check    func(t *testing.T, out string)
	}{
		{
			name:     "validate reports what it validated",
			args:     func(*testing.T) []string { return []string{"validate", "-strict", "-i", fixture} },
			wantCode: 0,
			check: func(t *testing.T, out string) {

				t.Helper()
				if !strings.HasPrefix(out, "validated "+fixture+": 3 roles, 2 variants") {
					t.Errorf("output %q", out)
				}
			},
		},
		{
			name: "render to stdout matches the golden file",
			args: func(*testing.T) []string {
				return []string{"render", "-i", fixture, "-variant", "focused", "-format", "json", "-o", "-"}
			},

			check: func(t *testing.T, out string) {

				t.Helper()
				if out != readFile(t, "../../testdata/golden/resume.focused.json") {
					t.Error("output differs from the golden file")
				}
			},
		},
		{
			name: "render reports each file it wrote",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"render", "-i", fixture, "-o", t.TempDir()}
			},

			check: func(t *testing.T, out string) {

				t.Helper()
				if strings.Count(out, "wrote ") != 4 || !strings.Contains(out, "resume.full.json") ||
					!strings.Contains(out, "resume.full.tex") {
					t.Errorf("output %q", out)
				}
			},
		},
		{
			name: "private overlay adds the phone",
			args: func(t *testing.T) []string {
				t.Helper()
				overlay := writeFile(t, "p.yaml", "basics:\n  phone: \"+1 555 0100\"\n")
				return []string{"render", "-i", fixture, "-private", overlay, "-variant", "full", "-format", "json", "-o", "-"}
			},

			check: func(t *testing.T, out string) {

				t.Helper()
				if !strings.Contains(out, `"phone": "+1 555 0100"`) {
					t.Error("phone missing")
				}
			},
		},
		{
			name: "private overlay adds the phone to the LaTeX",
			args: func(t *testing.T) []string {
				t.Helper()
				overlay := writeFile(t, "p.yaml", "basics:\n  phone: \"+1 555 0100\"\n")
				return []string{"render", "-i", fixture, "-private", overlay, "-variant", "full", "-format", "tex", "-o", "-"}
			},

			check: func(t *testing.T, out string) {

				t.Helper()
				if !strings.Contains(out, `\sep{}+1 555 0100\sep{}`) {
					t.Error("phone missing from the contact line")
				}
			},
		},
		{
			name: "no phone without the overlay",
			args: func(*testing.T) []string {
				return []string{"render", "-i", fixture, "-variant", "full", "-format", "json", "-o", "-"}
			},

			check: func(t *testing.T, out string) {

				t.Helper()
				if strings.Contains(out, "phone") {
					t.Error("phone present")
				}
			},
		},
		{
			name: "invalid résumé fails validation",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"validate", "-i", fixtureWith(t, "program: Computer Science major", "program: B.S. in Computer Science")}
			},

			wantCode: exitFailure,
			check: func(t *testing.T, out string) {

				t.Helper()
				if !strings.Contains(out, "not conferred") {
					t.Errorf("output %q", out)
				}
			},
		},
		{
			name: "invalid résumé isn't rendered",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"render", "-o", t.TempDir(), "-i", fixtureWith(t, "program: Computer Science major", "program: B.S.")}
			},

			wantCode: exitFailure,
			check: func(t *testing.T, out string) {

				t.Helper()
				if strings.Contains(out, "wrote") {
					t.Errorf("output %q", out)
				}
			},
		},
		{
			name: "unknown field is caught",
			args: func(t *testing.T) []string {
				t.Helper()
				return []string{"validate", "-i", fixtureWith(t, "  label: Staff", "  lable: Staff")}
			},

			wantCode: exitFailure,
			check: func(t *testing.T, out string) {

				t.Helper()
				if !strings.Contains(out, "lable") {
					t.Errorf("output %q", out)
				}
			},
		},
		{
			name: "list variants",
			args: func(*testing.T) []string { return []string{"list", "-i", fixture, "variants"} },
			check: func(t *testing.T, out string) {

				t.Helper()
				if !strings.Contains(out, "focused") || !strings.Contains(out, "full") {
					t.Errorf("output %q", out)
				}
			},
		},
		{name: "no command", args: func(*testing.T) []string { return nil }, wantCode: exitUsage},
		{name: "unknown command", args: func(*testing.T) []string { return []string{"frobnicate"} }, wantCode: exitUsage},
		{name: "unknown variant", args: func(*testing.T) []string { return []string{"render", "-i", fixture, "-variant", "nope"} }, wantCode: exitUsage},
		{name: "several files to stdout", args: func(*testing.T) []string { return []string{"render", "-i", fixture, "-o", "-"} }, wantCode: exitUsage},
		{name: "list without a subject", args: func(*testing.T) []string { return []string{"list", "-i", fixture} }, wantCode: exitUsage},
		{name: "bad flag", args: func(*testing.T) []string { return []string{"validate", "-nope"} }, wantCode: exitUsage},
		{name: "version", args: func(*testing.T) []string { return []string{"version"} }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			out, code := vitae(t, tc.args(t)...)
			if code != tc.wantCode {
				t.Fatalf("exit %d, want %d; output %q", code, tc.wantCode, out)
			}

			if tc.check != nil {
				tc.check(t, out)
			}
		})
	}
}
