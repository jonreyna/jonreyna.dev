# Play

Interactive pieces, one directory each. A piece is served at `/play/<piece>` and listed in the Play section of the home page (`site/content/index.md`).

## Pieces

- **[Arbor Ecclesiae](arbor-ecclesiae/)**, at `/play/arbor-ecclesiae`: the branching of the church from the Council of Jerusalem to about 1700, as a radial tree drawn with D3. Its data (every tradition, its dates, a short description and a Wikipedia link) is at the top of `arbor-ecclesiae.js`.

## How a piece is built

A piece that fills the window, like Arbor Ecclesiae, has three files named for its directory:

- `<piece>.html`: the page's markup, everything inside `<body>`, wrapped in `<main id="main">` and holding the page's one `<h1>`. The `play` layout (`site/lib/layouts/play_layout.dart`) puts it under the site's head: title, description, Open Graph, fonts, favicons and canonical URL.
- `<piece>.css`: its styles. They load after the site's `fonts.css`, without the site's tokens or styles, so set type in `"Alegreya"` and small capitals in `"Alegreya SC"` (weight 500).
- `<piece>.js`: its script, plus any libraries it needs.

Its page is `site/content/play/<piece>.md`, front matter only:

```yaml
title: Arbor Ecclesiae
description: One sentence for search results and shared links, at most about 155 characters.
layout: play
piece: arbor-ecclesiae
scripts: [d3.min.js, arbor-ecclesiae.js]  # in load order, each deferred
```

`make play` stages the files where the build serves them: the scripts in `site/web/play/<piece>/`, beside the page, and the stylesheet as `site/web/play-<piece>.css`, so `make dist` fingerprints it with the site's own (`/css/play-<piece>.HASH.css`). Both are git-ignored. A library is downloaded from a pinned release and checked against its SHA-256, as the fonts are, and never committed; its license is committed beside the piece and staged with it.

The security policy applies here as on every page: scripts, styles and fonts only from the site, no inline `<script>`, and no `style` attributes (use classes; setting styles from script is fine). `make verify` lets a page under `/play/<piece>/` load only the staged scripts in that directory.

Pieces built as Go→Wasm bundles, starting with the Kalman-filter image compression piece, will add their build target, loader and size budget here.
