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

func vitae(t *testing.T, args ...string) (stdout string, code int) {
	t.Helper()
	var out, errOut bytes.Buffer
	err := run(args, &out, &errOut)
	if err == nil {
		return out.String(), 0
	}
	var ee *exitError
	if errors.As(err, &ee) {
		return out.String() + errOut.String(), ee.code
	}
	return out.String() + errOut.String() + err.Error(), 1
}

func TestValidate(t *testing.T) {
	if out, code := vitae(t, "validate", "--strict", "-i", fixture); code != 0 {
		t.Fatalf("exit %d: %s", code, out)
	}
}

func TestRenderMatchesGolden(t *testing.T) {
	out, code := vitae(t, "render", "-i", fixture, "--variant", "focused", "--format", "json", "-o", "-")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, out)
	}
	want, err := os.ReadFile("../../testdata/golden/resume.focused.json")
	if err != nil {
		t.Fatal(err)
	}
	if out != string(want) {
		t.Errorf("render -o - differs from the golden file")
	}
}

func TestUsageErrors(t *testing.T) {
	for _, args := range [][]string{
		{},
		{"frobnicate"},
		{"render", "-i", fixture, "--variant", "nope"},
		{"render", "-i", fixture, "-o", "-"}, // all variants to one stream
		{"list", "-i", fixture},
	} {
		if _, code := vitae(t, args...); code != 2 {
			t.Errorf("%v: exit %d, want 2", args, code)
		}
	}
}

func TestValidationFailureExits1(t *testing.T) {
	bad := filepath.Join(t.TempDir(), "bad.yaml")
	data, _ := os.ReadFile(fixture)
	broken := strings.Replace(string(data), "program: Computer Science major", "program: B.S. in Computer Science", 1)
	if err := os.WriteFile(bad, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, code := vitae(t, "validate", "-i", bad); code != 1 || !strings.Contains(out, "not conferred") {
		t.Errorf("exit %d, output %q", code, out)
	}
	if _, code := vitae(t, "render", "-i", bad, "-o", t.TempDir()); code != 1 {
		t.Errorf("render of an invalid résumé: exit %d, want 1", code)
	}
}

func TestPrivateOverlay(t *testing.T) {
	overlay := filepath.Join(t.TempDir(), "p.yaml")
	os.WriteFile(overlay, []byte("basics:\n  phone: \"+1 555 0100\"\n"), 0o644)
	out, code := vitae(t, "render", "-i", fixture, "--private", overlay, "--variant", "full", "--format", "json", "-o", "-")
	if code != 0 || !strings.Contains(out, `"phone": "+1 555 0100"`) {
		t.Errorf("phone missing with --private (exit %d)", code)
	}
	out, _ = vitae(t, "render", "-i", fixture, "--variant", "full", "--format", "json", "-o", "-")
	if strings.Contains(out, "phone") {
		t.Errorf("phone present without --private")
	}
}

func TestUnknownFieldRejected(t *testing.T) {
	bad := filepath.Join(t.TempDir(), "typo.yaml")
	data, _ := os.ReadFile(fixture)
	os.WriteFile(bad, []byte(strings.Replace(string(data), "  label: Staff Engineer", "  lable: Staff Engineer", 1)), 0o644)
	if out, code := vitae(t, "validate", "-i", bad); code != 1 || !strings.Contains(out, "lable") {
		t.Errorf("typo not caught: exit %d, %q", code, out)
	}
}
