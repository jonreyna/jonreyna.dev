package main

import (
	"bytes"
	"errors"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/goregular"
)

func TestRun(t *testing.T) {

	t.Parallel()
	dir := t.TempDir()
	sc, it, bad := filepath.Join(dir, "sc.ttf"), filepath.Join(dir, "it.ttf"), filepath.Join(dir, "bad.ttf")
	for path, data := range map[string][]byte{sc: goregular.TTF, it: goitalic.TTF, bad: []byte("x")} {
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
	}

	text := []string{"-name", "Jon Reyna", "-motto", "a tagline", "-site", "jonreyna.dev"}
	tests := []struct {
		name      string
		args      []string
		wantUsage bool
		wantErr   string // empty for success
	}{
		{"writes the card", append([]string{"-smallcaps", sc, "-italic", it}, text...), false, ""},
		{"a missing flag", append([]string{"-smallcaps", sc}, text...), true, "-italic is required"},
		{"an unknown flag", []string{"-nope"}, true, "not defined"},
		{"a missing font file", append([]string{"-smallcaps", filepath.Join(dir, "none"), "-italic", it}, text...),
			false, "read small-capitals font"},
		{"a bad font", append([]string{"-smallcaps", sc, "-italic", bad}, text...), false, "parse italic font"},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			out := filepath.Join(dir, tc.name+".png")
			var stdout, stderr bytes.Buffer
			err := run(append(tc.args, "-o", out), &stdout, &stderr)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) || errors.Is(err, errUsage) != tc.wantUsage {
					t.Errorf("case %d: run error %v, want %q (usage %v)", i, err, tc.wantErr, tc.wantUsage)
				}

				return
			}

			if err != nil {
				t.Fatalf("run: %v", err)
			}

			if want := "wrote " + out + " (1200×630)\n"; stdout.String() != want {
				t.Errorf("stdout %q, want %q", stdout.String(), want)
			}

			data, err := os.ReadFile(out)
			if err != nil {
				t.Fatal(err)
			}

			cfg, err := png.DecodeConfig(bytes.NewReader(data))
			if err != nil || cfg.Width != 1200 || cfg.Height != 630 {
				t.Errorf("PNG %+v, %v; want 1200×630", cfg, err)
			}
		})
	}
}
