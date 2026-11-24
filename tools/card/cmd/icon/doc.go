// Command icon writes jonreyna.dev's favicons: a red letter in Alegreya SC on a
// parchment tile, drawn from the font's outline so no browser needs the font.
//
// Usage:
//
//	icon -font FONT -letter J -svg FILE.svg -ico FILE.ico -touch FILE.png
//
// The flags are all required:
//
//	-font    the font the letter is set in (Alegreya SC Medium)
//	-letter  the one character to draw
//	-svg     the SVG favicon to write: the letter as a path, switching to copper on
//	         leather when the viewer prefers a dark theme
//	-ico     the favicon.ico to write, holding 16- and 32-pixel PNGs, for browsers and
//	         tools that ask for that file
//	-touch   the 180-pixel apple-touch-icon PNG to write, with a smaller letter, since
//	         phones round its corners
//
// `make favicon` runs it with the pinned font it downloads, writing into site/web/
// (git-ignored), which the site build then publishes. On success it prints one line
// per file written, such as "wrote site/web/favicon.ico (16×16, 32×32)".
//
// Exit codes: 0 on success, 1 if an icon could not be made or written (a font that
// doesn't parse, a letter it has no glyph for), 2 for a usage error.
package main
