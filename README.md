# jonreyna.dev

Source for [jonreyna.dev](https://jonreyna.dev), Jon Reyna’s website: essays, selected work, interactive pieces and a résumé, all built from this repository.

The site is static. It’s built with [Jaspr](https://jaspr.site) and [jaspr_content](https://pub.dev/packages/jaspr_content) from Markdown and served by Cloudflare Workers Static Assets. Reading pages ship no JavaScript; math is MathML rendered at build time.

## Layout

```
site/        Jaspr project (static mode)
  content/   Markdown pages and data (content/_data/)
  lib/       layout and Markdown components
  web/       CSS, fonts, favicon, response headers
  DESIGN.md  the design system: color, type, components, typesetting rules
Makefile     the same targets locally and in CI
wrangler.jsonc
```

Planned: `resume/` (the résumé as one YAML file), `tools/vitae/` (a Go CLI that renders it to PDF, Word, text, Markdown and JSON) and `play/` (interactive pieces, each a self-contained bundle).

## Building

Needs Dart 3.13 (Jaspr 0.23.5 doesn’t resolve on Dart 3.14 yet) and, for deploying, Node.

```bash
make setup    # dart pub get, and the Jaspr CLI
make serve    # build, stage dist/, serve on http://localhost:8767
make deploy   # upload dist/ to Cloudflare
```

`make help` lists every target.
