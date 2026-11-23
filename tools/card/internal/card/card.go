package card

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Width and Height are the card's size in pixels, the 1.91:1 shape Open Graph
// consumers expect.
const (
	Width  = 1200
	Height = 630
)

// The light theme's colors (site/DESIGN.md, "Color"), as 0xRRGGBB.
const (
	hexBG    = 0xf6f0e1 // --bg, parchment
	hexText  = 0x2a2119 // --text, the name
	hexMuted = 0x6d5e4a // --muted, the tagline
	hexRed   = 0x9a2a17 // --red, rules and address

	// The bit offsets of the red and green bytes in 0xRRGGBB, and one byte's mask.
	redShift   = 16
	greenShift = 8
	byteMask   = 0xff
)

// rgb is the opaque color 0xRRGGBB.
func rgb(hex uint32) color.RGBA {

	return color.RGBA{
		R: uint8(hex >> redShift & byteMask),
		G: uint8(hex >> greenShift & byteMask),
		B: uint8(hex & byteMask),
		A: byteMask,
	}
}

// The layout, in pixels: a column 840 wide between the rules, with the masthead's
// proportions at about three times the site's desktop size.
const (
	columnLeft  = 180
	columnRight = Width - columnLeft
	ruleLine    = 3   // thickness of each rule line
	topRule     = 150 // top of the double rule: two lines with a gap of one line
	bottomRule  = 490 // top of the single rule

	nameSize     = 120
	nameTracking = 0.1 // em, as the site's masthead name
	nameBaseline = 300

	mottoSize     = 46
	mottoBaseline = 375

	siteSize     = 40
	siteTracking = 0.06 // em, as the site's nav
	siteBaseline = 445

	// dpi makes a font size in points a size in pixels.
	dpi = 72
)

// errTooWide reports a line of text that doesn't fit between the rules.
var errTooWide = errors.New("does not fit the card")

// Card is the text of an Open Graph image.
type Card struct {
	// Name is the largest line, set in small capitals: "Jon Reyna".
	Name string
	// Motto is the italic tagline under it.
	Motto string
	// Site is the address on the last line, such as "jonreyna.dev".
	Site string
}

// Render draws the card with fonts. It fails if a line is too wide for the column
// between the rules.
func (c *Card) Render(fonts *Fonts) (*image.RGBA, error) {

	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	draw.Draw(img, img.Bounds(), image.NewUniform(rgb(hexBG)), image.Point{}, draw.Src)

	fillRule(img, topRule)
	fillRule(img, topRule+2*ruleLine)
	fillRule(img, bottomRule)

	lines := []struct {
		what     string
		text     string
		font     *opentype.Font
		size     float64
		tracking float64
		baseline int
		color    color.Color
	}{
		{"name", c.Name, fonts.SmallCaps, nameSize, nameTracking, nameBaseline, rgb(hexText)},
		{"motto", c.Motto, fonts.Italic, mottoSize, 0, mottoBaseline, rgb(hexMuted)},
		{"site", c.Site, fonts.SmallCaps, siteSize, siteTracking, siteBaseline, rgb(hexRed)},
	}
	for _, l := range lines {
		face, err := opentype.NewFace(l.font, &opentype.FaceOptions{Size: l.size, DPI: dpi, Hinting: font.HintingNone})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", l.what, err)
		}

		err = drawCentered(img, face, l.text, l.tracking*l.size, l.baseline, l.color)
		if err != nil {
			return nil, fmt.Errorf("%s %q: %w", l.what, l.text, err)
		}
	}

	return img, nil
}

// fillRule draws one red rule line across the column, its top at y.
func fillRule(img *image.RGBA, y int) {
	r := image.Rect(columnLeft, y, columnRight, y+ruleLine)
	draw.Draw(img, r, image.NewUniform(rgb(hexRed)), image.Point{}, draw.Src)
}

// drawCentered sets text on the baseline, centered on the card, with tracking pixels
// added between letters.
func drawCentered(img *image.RGBA, face font.Face, text string, tracking float64, baseline int, c color.Color) error {

	track := fixed.Int26_6(tracking * float64(fixed.I(1)))
	width := advance(face, text, track)
	if width.Ceil() > columnRight-columnLeft {
		return errTooWide
	}

	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: face}
	d.Dot = fixed.Point26_6{X: fixed.I(Width)/2 - width/2, Y: fixed.I(baseline)}
	prev := rune(-1)
	for _, r := range text {
		if prev >= 0 {
			d.Dot.X += face.Kern(prev, r) + track
		}

		d.DrawString(string(r))
		prev = r
	}

	return nil
}

// advance is text's width in face, with track added between letters.
func advance(face font.Face, text string, track fixed.Int26_6) fixed.Int26_6 {

	var width fixed.Int26_6
	prev := rune(-1)
	for _, r := range text {
		if prev >= 0 {
			width += face.Kern(prev, r) + track
		}

		a, _ := face.GlyphAdvance(r)
		width += a
		prev = r
	}

	return width
}
