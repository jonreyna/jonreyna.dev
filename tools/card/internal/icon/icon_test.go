package icon

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"math"
	"regexp"
	"strings"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/sfnt"
)

// newTestIcon returns an icon of the letter in Go Regular, in place of Alegreya SC,
// so tests need no download.
func newTestIcon(t *testing.T, letter rune, height float64) *Icon {

	t.Helper()
	f, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("parse Go Regular: %v", err)
	}

	return &Icon{Letter: letter, Font: f, Height: height}
}

// inkBounds is the smallest rectangle holding every pixel of img that isn't ground.
func inkBounds(img *image.RGBA, ground color.RGBA) image.Rectangle {

	var ink image.Rectangle
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if img.RGBAAt(x, y) != ground {
				ink = ink.Union(image.Rect(x, y, x+1, y+1))
			}
		}
	}

	return ink
}

func TestRender(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name   string
		letter rune
		height float64
		size   int
	}{
		{"a capital on a touch icon", 'J', 0.62, 180},
		{"a capital filling most of a tab icon", 'J', 0.8, 64},
		{"a letter with two contours", 'B', 0.7, 120},
		{"the whole height", 'H', 1, 50},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			p := Light()
			img, err := newTestIcon(t, tc.letter, tc.height).Render(tc.size, p)
			if err != nil {
				t.Fatalf("Render: %v", err)
			}

			if got := img.Bounds(); got != image.Rect(0, 0, tc.size, tc.size) {
				t.Fatalf("bounds %v, want %d square", got, tc.size)
			}

			if got := img.RGBAAt(0, 0); got != p.Ground && tc.height < 1 {
				t.Errorf("corner is %v, want the ground %v", got, p.Ground)
			}

			ink := inkBounds(img, p.Ground)
			if got, want := float64(ink.Dy()), tc.height*float64(tc.size); math.Abs(got-want) > 2 {
				t.Errorf("ink is %v px tall, want about %v", got, want)
			}

			// Centered: the margins on opposite sides differ by at most a pixel or two
			// of antialiasing.
			left, right := ink.Min.X, tc.size-ink.Max.X
			top, bottom := ink.Min.Y, tc.size-ink.Max.Y
			if abs(left-right) > 2 || abs(top-bottom) > 2 {
				t.Errorf("ink %v isn't centered on a %d tile", ink, tc.size)
			}

			if !hasColor(img, p.Letter) {
				t.Errorf("no pixel has the letter color %v", p.Letter)
			}
		})
	}
}

func TestRenderIsDeterministic(t *testing.T) {

	t.Parallel()
	ic := newTestIcon(t, 'J', 0.8)
	a, errA := ic.Render(32, Light())
	b, errB := ic.Render(32, Light())
	if errA != nil || errB != nil {
		t.Fatalf("Render: %v, %v", errA, errB)
	}

	if !bytes.Equal(a.Pix, b.Pix) {
		t.Error("two renders of the same icon differ")
	}
}

func TestSVG(t *testing.T) {

	t.Parallel()
	svg, err := newTestIcon(t, 'J', 0.8).SVG(Light(), Dark())
	if err != nil {
		t.Fatalf("SVG: %v", err)
	}

	s := string(svg)
	tests := []struct {
		name string
		want string // a regular expression s must match
	}{
		{"is a 64-unit square", `^<svg xmlns="http://www\.w3\.org/2000/svg" viewBox="0 0 64 64">`},
		{"light ground and letter", `\.g\{fill:#f6f0e1\}\.l\{fill:#9a2a17\}`},
		{"dark ground and letter", `@media \(prefers-color-scheme:dark\)\{\.g\{fill:#1a1612\}\.l\{fill:#e47d63\}\}`},
		{"a tile under the letter", `<rect class="g" width="64" height="64"/><path class="l" d="M`},
		{"a closed path", `Z"/></svg>\n$`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			if !regexp.MustCompile(tc.want).MatchString(s) {
				t.Errorf("SVG doesn't match %s:\n%s", tc.want, s)
			}
		})
	}

	d := regexp.MustCompile(` d="([^"]*)"`).FindStringSubmatch(s)
	if d == nil {
		t.Fatal("SVG has no path data")
	}

	if long := regexp.MustCompile(`\d\.\d{3}|[eE]`).FindString(d[1]); long != "" {
		t.Errorf("SVG has a number with more than two decimals or an exponent: %s", long)
	}
}

func TestErrors(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name   string
		letter rune
		height float64
		size   int
		want   error
	}{
		{"a letter the font lacks", '\U0001F600', 0.8, 32, ErrNoGlyph},
		{"no height", 'J', 0, 32, ErrHeight},
		{"taller than the tile", 'J', 1.2, 32, ErrHeight},
		{"no pixels", 'J', 0.8, 0, ErrSize},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			ic := newTestIcon(t, tc.letter, tc.height)
			_, err := ic.Render(tc.size, Light())
			if !errors.Is(err, tc.want) {
				t.Errorf("Render error %v, want %v", err, tc.want)
			}

			if errors.Is(tc.want, ErrSize) {
				return
			}

			_, err = ic.SVG(Light(), Dark())
			if !errors.Is(err, tc.want) {
				t.Errorf("SVG error %v, want %v", err, tc.want)
			}
		})
	}
}

func TestNumber(t *testing.T) {

	t.Parallel()
	tests := []struct {
		in   float64
		want string
	}{
		{12, "12"},
		{12.5, "12.5"},
		{12.345, "12.35"},
		{0.004, "0"},
		{-0.004, "0"},
		{-3.1, "-3.1"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {

			t.Parallel()
			if got := number(tc.in); got != tc.want {
				t.Errorf("number(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestCSSColor(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name string
		in   color.RGBA
		want string
	}{
		{"light letter", Light().Letter, "#9a2a17"},
		{"dark ground", Dark().Ground, "#1a1612"},
		{"pads bytes", color.RGBA{R: 1, G: 2, B: 3, A: 255}, "#010203"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			if got := cssColor(tc.in); !strings.EqualFold(got, tc.want) {
				t.Errorf("cssColor = %q, want %q", got, tc.want)
			}
		})
	}
}

// hasColor reports whether any pixel of img is exactly c.
func hasColor(img *image.RGBA, c color.RGBA) bool {

	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if img.RGBAAt(x, y) == c {
				return true
			}
		}
	}

	return false
}

func abs(n int) int { return max(n, -n) }
