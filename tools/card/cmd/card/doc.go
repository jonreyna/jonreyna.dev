// Command card writes jonreyna.dev's Open Graph image, the picture shown beside a
// shared link, as a 1200×630 PNG.
//
// Usage:
//
//	card -smallcaps FONT -italic FONT -name TEXT -motto TEXT -site TEXT -o FILE.png
//
// The flags are all required:
//
//	-smallcaps  the small-capitals font for the name and address (Alegreya SC Medium)
//	-italic     the italic font for the tagline (Alegreya Italic)
//	-name       the name, the largest line
//	-motto      the tagline under the name
//	-site       the address on the last line
//	-o          the PNG to write
//
// `make og-image` runs it with the pinned fonts it downloads, writing
// site/web/og.png (git-ignored), which the site build then publishes. On success it
// prints one line, "wrote FILE (1200×630)".
//
// Exit codes: 0 on success, 1 if the image could not be made or written (a font that
// doesn't parse, a line too wide for the card), 2 for a usage error.
package main
