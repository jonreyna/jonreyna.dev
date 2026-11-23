// Package card draws jonreyna.dev's Open Graph image: the picture LinkedIn and other
// sites show beside a shared link.
//
// The image is the site's masthead set as a title page on parchment (site/DESIGN.md):
// a red double rule, the name in Alegreya SC, the tagline in Alegreya italic, the
// site's address in small capitals, and a single red rule, centered on a
// [Width]×[Height] canvas in the light theme's colors.
//
// A [Card] holds the three lines of text and [Fonts] holds the two typefaces;
// [Card.Render] draws them. [ParseFonts] reads the fonts from TrueType or OpenType
// data, so the caller decides where the files come from (the build downloads pinned
// copies; tests use the Go fonts).
//
// Rendering is deterministic: the same text and fonts always give the same pixels.
// Kerning comes from a font's kern table when it has one; fonts that kern only
// through GPOS, as Alegreya does, are set without kerning, which the letter-spaced
// small capitals don't need.
package card
