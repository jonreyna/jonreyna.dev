package card

import (
	"fmt"

	"golang.org/x/image/font/opentype"
)

// Fonts holds the two typefaces a card is set in.
type Fonts struct {
	// SmallCaps sets the name and the site's address: a small-capitals face, such as
	// Alegreya SC Medium, whose lowercase letters are small capitals.
	SmallCaps *opentype.Font
	// Italic sets the tagline, such as Alegreya Italic.
	Italic *opentype.Font
}

// ParseFonts reads the small-capitals and italic faces from TrueType or OpenType
// font data.
func ParseFonts(smallCaps, italic []byte) (*Fonts, error) {

	sc, err := opentype.Parse(smallCaps)
	if err != nil {
		return nil, fmt.Errorf("parse small-capitals font: %w", err)
	}

	it, err := opentype.Parse(italic)
	if err != nil {
		return nil, fmt.Errorf("parse italic font: %w", err)
	}

	return &Fonts{SmallCaps: sc, Italic: it}, nil
}
