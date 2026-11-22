package resume

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fixture = "../../testdata/resume.yaml"

// writeFile writes data to name in a fresh temporary directory and returns the path.
func writeFile(t *testing.T, name, data string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func fixtureWith(t *testing.T, old, replacement string) string {
	t.Helper()
	data, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}
	return writeFile(t, "resume.yaml", strings.Replace(string(data), old, replacement, 1))
}

func TestLoad(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		path    func(t *testing.T) string
		wantErr string
	}{
		{name: "fixture", path: func(*testing.T) string { return fixture }},
		{name: "missing file", path: func(t *testing.T) string {
			t.Helper()
			return filepath.Join(t.TempDir(), "none.yaml")
		}, wantErr: "no such file"},
		{name: "unknown field", path: func(t *testing.T) string {
			t.Helper()
			return fixtureWith(t, "  label: Staff Engineer", "  lable: Staff Engineer")
		}, wantErr: "lable"},
		{name: "bad month", path: func(t *testing.T) string {
			t.Helper()
			return fixtureWith(t, "start: 2024-01", "start: 2024-1")
		}, wantErr: "not a YYYY-MM month"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r, err := Load(tc.path(t))
			if tc.wantErr == "" {
				if err != nil || r.Basics.Name != "Ada Example" {
					t.Fatalf("Load() = %v, %v", r, err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Load() error = %v, want it to mention %q", err, tc.wantErr)
			}
		})
	}
}

func TestApplyPrivate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		overlay   string
		wantPhone string
		wantErr   bool
	}{
		{name: "sets phone", overlay: "basics:\n  phone: \"+1 555 0100\"\n", wantPhone: "+1 555 0100"},
		{name: "empty keeps phone", overlay: "basics: {}\n", wantPhone: ""},
		{name: "unknown field", overlay: "basics:\n  address: 1 Main St\n", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r, err := Load(fixture)
			if err != nil {
				t.Fatal(err)
			}
			err = ApplyPrivate(r, writeFile(t, "private.yaml", tc.overlay))
			if (err != nil) != tc.wantErr {
				t.Fatalf("ApplyPrivate() error = %v, wantErr %v", err, tc.wantErr)
			}
			if r.Basics.Phone != tc.wantPhone {
				t.Errorf("phone = %q, want %q", r.Basics.Phone, tc.wantPhone)
			}
		})
	}
}
