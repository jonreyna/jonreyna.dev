package icon

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"
)

// svgSide is the SVG's viewBox side; the path's coordinates are in these units.
const svgSide = 64

// half finds a midpoint.
const half = 0.5

// Errors a caller can test for with errors.Is.
var (
	// ErrNoGlyph reports a letter the font has no glyph for.
	ErrNoGlyph = errors.New("the font has no glyph for the letter")
	// ErrHeight reports an Icon.Height outside (0, 1].
	ErrHeight = errors.New("height must be more than 0 and at most 1")
	// ErrSize reports a PNG size below 1 pixel.
	ErrSize = errors.New("size must be at least 1 pixel")
)

// Icon is a square tile with one letter centered on it.
type Icon struct {
	// Letter is the character drawn, such as 'J'.
	Letter rune
	// Font is the typeface the letter is set in, such as Alegreya SC Medium.
	Font *sfnt.Font
	// Height is the letter's ink height as a fraction of the tile's side, more than 0
	// and at most 1. Its width follows from the glyph's proportions.
	Height float64
}

// point is a position on the tile, in the tile's units, with y increasing down.
type point struct{ x, y float64 }

// segment is one step of an outline: a move, a line, or a quadratic or cubic curve,
// using the first one, two or three of pts.
type segment struct {
	op  sfnt.SegmentOp
	pts [3]point
}

// SVG returns the icon as an SVG document, svgSide units square: a tile in light's
// ground with the letter's outline as a path in light's letter color, switching to
// dark's colors when the viewer prefers a dark theme.
func (ic *Icon) SVG(light, dark Palette) ([]byte, error) {

	segs, err := ic.outline(svgSide)
	if err != nil {
		return nil, err
	}

	var path strings.Builder
	for i, s := range segs {
		if s.op == sfnt.SegmentOpMoveTo && i > 0 {
			path.WriteString("Z")
		}

		path.WriteString(map[sfnt.SegmentOp]string{
			sfnt.SegmentOpMoveTo: "M", sfnt.SegmentOpLineTo: "L",
			sfnt.SegmentOpQuadTo: "Q", sfnt.SegmentOpCubeTo: "C",
		}[s.op])
		for j, p := range s.pts[:args(s.op)] {
			if j > 0 {
				path.WriteString(" ")
			}

			path.WriteString(number(p.x) + " " + number(p.y))
		}
	}

	path.WriteString("Z")

	return fmt.Appendf(nil, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %[1]d %[1]d">`+
		`<style>.g{fill:%[2]s}.l{fill:%[3]s}`+
		`@media (prefers-color-scheme:dark){.g{fill:%[4]s}.l{fill:%[5]s}}</style>`+
		`<rect class="g" fill="%[2]s" width="%[1]d" height="%[1]d"/>`+
		`<path class="l" fill="%[3]s" d="%[6]s"/></svg>`+"\n",
		svgSide, cssColor(light.Ground), cssColor(light.Letter),
		cssColor(dark.Ground), cssColor(dark.Letter), path.String()), nil
}

// Render draws the icon size pixels square in p's colors.
func (ic *Icon) Render(size int, p Palette) (*image.RGBA, error) {

	if size < 1 {
		return nil, fmt.Errorf("%w (got %d)", ErrSize, size)
	}

	segs, err := ic.outline(float64(size))
	if err != nil {
		return nil, err
	}

	r := vector.NewRasterizer(size, size)
	for i, s := range segs {
		pt := func(k int) (float32, float32) { return float32(s.pts[k].x), float32(s.pts[k].y) }
		switch s.op {
		case sfnt.SegmentOpMoveTo:
			if i > 0 {
				r.ClosePath()
			}

			r.MoveTo(pt(0))
		case sfnt.SegmentOpLineTo:
			r.LineTo(pt(0))
		case sfnt.SegmentOpQuadTo:
			x0, y0 := pt(0)
			x1, y1 := pt(1)
			r.QuadTo(x0, y0, x1, y1)
		case sfnt.SegmentOpCubeTo:
			x0, y0 := pt(0)
			x1, y1 := pt(1)
			x2, y2 := pt(2) //nolint:mnd // the cubic's end point, its third
			r.CubeTo(x0, y0, x1, y1, x2, y2)
		}
	}

	r.ClosePath()

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), image.NewUniform(p.Ground), image.Point{}, draw.Src)
	r.Draw(img, img.Bounds(), image.NewUniform(p.Letter), image.Point{})

	return img, nil
}

// outline is the letter's outline scaled and centered on a tile side units wide.
func (ic *Icon) outline(side float64) ([]segment, error) {

	if ic.Height <= 0 || ic.Height > 1 {
		return nil, fmt.Errorf("%w (got %g)", ErrHeight, ic.Height)
	}

	var buf sfnt.Buffer
	idx, err := ic.Font.GlyphIndex(&buf, ic.Letter)
	if err != nil {
		return nil, fmt.Errorf("find %q: %w", ic.Letter, err)
	}

	if idx == 0 {
		return nil, fmt.Errorf("%q: %w", ic.Letter, ErrNoGlyph)
	}

	// Load at one pixel per font unit, the font's own precision.
	ppem := fixed.I(int(ic.Font.UnitsPerEm()))
	bounds, _, err := ic.Font.GlyphBounds(&buf, idx, ppem, font.HintingNone)
	if err != nil {
		return nil, fmt.Errorf("measure %q: %w", ic.Letter, err)
	}

	segs, err := ic.Font.LoadGlyph(&buf, idx, ppem, nil)
	if err != nil {
		return nil, fmt.Errorf("load %q: %w", ic.Letter, err)
	}

	minX, minY := fromFixed(bounds.Min.X), fromFixed(bounds.Min.Y)
	maxX, maxY := fromFixed(bounds.Max.X), fromFixed(bounds.Max.Y)
	scale := ic.Height * side / (maxY - minY)
	dx := (side - (minX+maxX)*scale) * half
	dy := (side - (minY+maxY)*scale) * half

	out := make([]segment, len(segs))
	for i, s := range segs {
		out[i].op = s.Op
		for j, a := range s.Args {
			out[i].pts[j] = point{fromFixed(a.X)*scale + dx, fromFixed(a.Y)*scale + dy}
		}
	}

	return out, nil
}

// args is how many points each segment op uses.
func args(op sfnt.SegmentOp) int {

	switch op {
	case sfnt.SegmentOpMoveTo, sfnt.SegmentOpLineTo:
		return 1
	case sfnt.SegmentOpQuadTo:
		return 2 //nolint:mnd // a quadratic curve: one control point and the end
	case sfnt.SegmentOpCubeTo:
		return 3 //nolint:mnd // a cubic curve: two control points and the end
	}

	return 0
}

// fromFixed is v as a float.
func fromFixed(v fixed.Int26_6) float64 { return float64(v) / float64(fixed.I(1)) }

// number formats v for a path: at most two decimals, without trailing zeros.
func number(v float64) string {

	s := strconv.FormatFloat(v, 'f', 2, 64)
	s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	if s == "-0" {
		return "0"
	}

	return s
}
