package main

import (
	"bytes"
	"errors"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
)

func TestRun(t *testing.T) {

	t.Parallel()
	dir := t.TempDir()
	good, bad := filepath.Join(dir, "good.ttf"), filepath.Join(dir, "bad.ttf")
	for path, data := range map[string][]byte{good: goregular.TTF, bad: []byte("x")} {
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name      string
		args      []string
		wantUsage bool
		wantErr   string // empty for success
	}{
		{"writes the icons", []string{"-font", good, "-letter", "J"}, false, ""},
		{"a missing flag", []string{"-letter", "J"}, true, "-font is required"},
		{"an unknown flag", []string{"-nope"}, true, "not defined"},
		{"two letters", []string{"-font", good, "-letter", "JR"}, true, "one character"},
		{"a missing font file", []string{"-font", filepath.Join(dir, "none"), "-letter", "J"}, false, "read font"},
		{"a bad font", []string{"-font", bad, "-letter", "J"}, false, "parse font"},
		{"a letter the font lacks", []string{"-font", good, "-letter", "\U0001F600"}, false, "no glyph"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			out := filepath.Join(dir, strings.ReplaceAll(tc.name, " ", "-"))
			if err := os.Mkdir(out, 0o700); err != nil {
				t.Fatal(err)
			}

			svg, ico, touch := filepath.Join(out, "f.svg"), filepath.Join(out, "f.ico"), filepath.Join(out, "t.png")
			var stdout, stderr bytes.Buffer
			err := run(append(tc.args, "-svg", svg, "-ico", ico, "-touch", touch), &stdout, &stderr)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) || errors.Is(err, errUsage) != tc.wantUsage {
					t.Errorf("run error %v, want %q (usage %v)", err, tc.wantErr, tc.wantUsage)
				}

				return
			}

			if err != nil {
				t.Fatalf("run: %v", err)
			}

			want := "wrote " + svg + "\nwrote " + ico + " (16×16, 32×32)\nwrote " + touch + " (180×180)\n"
			if stdout.String() != want {
				t.Errorf("stdout %q, want %q", stdout.String(), want)
			}

			checkFiles(t, svg, ico, touch)
		})
	}
}

// checkFiles checks that run wrote an SVG, an ICO and a 180-pixel PNG.
func checkFiles(t *testing.T, svg, ico, touch string) {

	t.Helper()
	for path, prefix := range map[string]string{svg: "<svg ", ico: "\x00\x00\x01\x00\x02\x00", touch: "\x89PNG"} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}

		if !bytes.HasPrefix(data, []byte(prefix)) {
			t.Errorf("%s starts %q, want %q", path, data[:min(len(data), len(prefix))], prefix)
		}
	}

	data, err := os.ReadFile(touch)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := png.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width != touchSize || cfg.Height != touchSize {
		t.Errorf("touch icon is %d×%d (%v), want %d square", cfg.Width, cfg.Height, err, touchSize)
	}
}
