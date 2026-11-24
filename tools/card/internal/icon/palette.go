package icon

import (
	"fmt"
	"image/color"
)

// The site's colors (site/DESIGN.md, "Color"), as 0xRRGGBB.
const (
	hexLightGround = 0xf6f0e1 // --bg, parchment
	hexLightLetter = 0x9a2a17 // --red, rubric red
	hexDarkGround  = 0x1a1612 // --bg in dark mode, leather
	hexDarkLetter  = 0xe47d63 // --red in dark mode, copper

	// The bit offsets of the red and green bytes in 0xRRGGBB, and one byte's mask.
	redShift   = 16
	greenShift = 8
	byteMask   = 0xff
)

// Palette is an icon's two colors.
type Palette struct {
	// Ground fills the tile.
	Ground color.RGBA
	// Letter fills the letter.
	Letter color.RGBA
}

// Light is the light theme's palette: red on parchment.
func Light() Palette { return Palette{Ground: rgb(hexLightGround), Letter: rgb(hexLightLetter)} }

// Dark is the dark theme's palette: copper on leather.
func Dark() Palette { return Palette{Ground: rgb(hexDarkGround), Letter: rgb(hexDarkLetter)} }

// rgb is the opaque color 0xRRGGBB.
func rgb(hex uint32) color.RGBA {

	return color.RGBA{
		R: uint8(hex >> redShift & byteMask),
		G: uint8(hex >> greenShift & byteMask),
		B: uint8(hex & byteMask),
		A: byteMask,
	}
}

// cssColor is c as a CSS hex color, such as "#9a2a17". It ignores alpha.
func cssColor(c color.RGBA) string { return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B) }
