// Package icon draws jonreyna.dev's favicon: a red "J" in Alegreya SC, centered on
// a parchment tile (site/DESIGN.md).
//
// An [Icon] holds the letter, the font and how tall the letter stands on the tile.
// [Icon.SVG] writes the tile as an SVG whose letter is the glyph's outline as a path,
// so no browser needs the font, and which switches to [Dark] colors when the viewer
// prefers a dark theme. [Icon.Render] draws the same outline as pixels in one
// [Palette], for the PNG icons, and [EncodeICO] packs PNGs into a favicon.ico for
// browsers and tools that ask for that file.
//
// The letter is centered by its ink, the glyph's bounding box, not by its advance and
// baseline, since a J's hook drops below the baseline. Output is deterministic: the
// same letter, font and size always give the same bytes.
package icon
