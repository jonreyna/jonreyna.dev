package card

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"strings"
	"testing"

	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/goregular"
)

// newTestFonts returns the Go fonts in place of Alegreya, so tests need no download.
func newTestFonts(t *testing.T) *Fonts {

	t.Helper()
	fonts, err := ParseFonts(goregular.TTF, goitalic.TTF)
	if err != nil {
		t.Fatalf("ParseFonts: %v", err)
	}

	return fonts
}

// newTestCard renders the site's own text with the test fonts.
func newTestCard(t *testing.T) *image.RGBA {

	t.Helper()
	c := &Card{Name: "Jon Reyna", Motto: "distributed systems & AI infrastructure", Site: "jonreyna.dev"}
	img, err := c.Render(newTestFonts(t))
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	return img
}

// count is the number of pixels in r that have color c exactly.
func count(img *image.RGBA, r image.Rectangle, c color.RGBA) int {

	n := 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if img.RGBAAt(x, y) == c {
				n++
			}
		}
	}

	return n
}

func TestRender(t *testing.T) {

	t.Parallel()
	img := newTestCard(t)
	row := func(y int) image.Rectangle { return image.Rect(columnLeft, y, columnRight, y+1) }
	tests := []struct {
		name string
		area image.Rectangle
		want color.RGBA
		min  int // at least this many pixels of want in area
	}{
		{"size", image.Rect(0, 0, 1, 1), rgb(hexBG), 1},
		{"corners are parchment", image.Rect(0, 0, 40, 40), rgb(hexBG), 40 * 40},
		{"double rule, first line", image.Rect(columnLeft, topRule, columnRight, topRule+ruleLine), rgb(hexRed), 840 * 3},
		{"double rule, gap", image.Rect(columnLeft, topRule+ruleLine, columnRight, topRule+2*ruleLine), rgb(hexBG), 840 * 3},
		{"double rule, second line", row(topRule + 2*ruleLine), rgb(hexRed), 840},
		{"rules stop at the column", image.Rect(0, topRule, columnLeft, topRule+ruleLine), rgb(hexBG), columnLeft * 3},
		{"single rule", image.Rect(columnLeft, bottomRule, columnRight, bottomRule+ruleLine), rgb(hexRed), 840 * 3},
		{"name in ink", image.Rect(columnLeft, nameBaseline-nameSize, columnRight, nameBaseline), rgb(hexText), 500},
		{"motto in muted", image.Rect(columnLeft, mottoBaseline-mottoSize, columnRight, mottoBaseline), rgb(hexMuted), 200},
		{"site in red", image.Rect(columnLeft, siteBaseline-siteSize, columnRight, siteBaseline), rgb(hexRed), 100},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			if got := count(img, tc.area, tc.want); got < tc.min {
				t.Errorf("%d pixels of %v in %v, want at least %d", got, tc.want, tc.area, tc.min)
			}
		})
	}

	if b := img.Bounds(); b.Dx() != Width || b.Dy() != Height {
		t.Errorf("bounds %v, want %d×%d", b, Width, Height)
	}
}

func TestRenderCentersEachLine(t *testing.T) {

	t.Parallel()
	img := newTestCard(t)
	bands := []struct {
		name     string
		top, bot int
		ink      color.RGBA
	}{
		{"name", nameBaseline - nameSize, nameBaseline, rgb(hexText)},
		{"motto", mottoBaseline - mottoSize, mottoBaseline, rgb(hexMuted)},
		{"site", siteBaseline - siteSize, siteBaseline, rgb(hexRed)},
	}
	for _, b := range bands {
		t.Run(b.name, func(t *testing.T) {

			t.Parallel()
			left, right := Width, 0
			for y := b.top; y < b.bot; y++ {
				for x := range Width {
					if img.RGBAAt(x, y) == b.ink {
						left, right = min(left, x), max(right, x)
					}
				}
			}

			// Glyph side bearings make ink a few pixels off its advance box.
			if mid := (left + right) / 2; mid < Width/2-8 || mid > Width/2+8 {
				t.Errorf("ink spans %d–%d, centered on %d, want about %d", left, right, mid, Width/2)
			}
		})
	}
}

func TestRenderRejects(t *testing.T) {

	t.Parallel()
	long := strings.Repeat("Reyna ", 12)
	tests := []struct {
		name string
		card Card
	}{
		{"a name too wide", Card{Name: long, Motto: "m", Site: "s"}},
		{"a motto too wide", Card{Name: "n", Motto: strings.Repeat(long, 3), Site: "s"}},
		{"an address too wide", Card{Name: "n", Motto: "m", Site: strings.Repeat(long, 3)}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			_, err := tc.card.Render(newTestFonts(t))
			if !errors.Is(err, errTooWide) {
				t.Errorf("Render error %v, want errTooWide", err)
			}
		})
	}
}

// One scenario, so no table: the same input must give the same pixels, or the
// published image would change on every build.
func TestRenderIsDeterministic(t *testing.T) {

	t.Parallel()
	a, b := newTestCard(t), newTestCard(t)
	if !bytes.Equal(a.Pix, b.Pix) {
		t.Error("two renders of the same card differ")
	}
}
