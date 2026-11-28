# jonreyna.dev

Source for [jonreyna.dev](https://jonreyna.dev), Jon Reyna’s website: essays, selected work, interactive pieces and a résumé, all built from this repository.

The site is static. It’s built with [Jaspr](https://jaspr.site) and [jaspr_content](https://pub.dev/packages/jaspr_content) from Markdown and served by Cloudflare Workers Static Assets. Reading pages ship no JavaScript; math is written in TeX (`$…$`, `$$…$$`) and converted to MathML, and code blocks are highlighted, at build time with Pandoc (see `site/DESIGN.md`).

## Layout

```
site/        Jaspr project (static mode)
  content/   Markdown pages, essays (content/writing/) and data (content/_data/)
  lib/       layout and Markdown components
  web/       CSS, fonts, favicon, response headers
  DESIGN.md  the design system: color, type, components, typesetting rules
resume/      the résumé as one YAML file (resume.yaml), the single source for every version
tools/vitae/ Go CLI that validates resume.yaml and renders each variant (see its README)
tools/card/  Go CLIs that draw the Open Graph image (the picture beside a shared link) and the favicons
play/        interactive pieces, one directory each (see play/README.md)
Makefile     the same targets locally and in CI
wrangler.jsonc
```

## Building

Needs Dart 3.13 (Jaspr 0.23.5 doesn’t resolve on Dart 3.14 yet), Go for the résumé tools, TeX Live with latexmk for the résumé PDFs (the packages are listed in `resume/latex/texlive-packages.txt`), Python 3 and, for deploying, Node. `make` downloads a pinned Pandoc release for the Word files, the math and code highlighting, and a pinned D3 release for Arbor Ecclesiae, and installs fontTools (pinned with hashes in `site/tool/requirements.txt`) into `build/venv` to subset the fonts.

```bash
make setup    # dart pub get, the Jaspr CLI, and the git hooks (no binaries in git)
make serve    # build (the résumé page included), stage dist/, serve on http://localhost:8767
make site-test   # the site's Dart tests (TeX to MathML, code highlighting)
make deploy   # upload dist/ to Cloudflare
make resume   # render every variant of resume/resume.yaml into build/resume
make resume-pdf   # ...and build a PDF of each (fails if one runs past its max_pages)
make resume-word  # ...and a Word file of each, with Pandoc
make resume-release  # stage them under their download names in build/resume-release
```

`make help` lists every target.

## CI/CD

`.github/workflows/site.yml` builds, analyzes and checks the site (`make verify`: scripts, titles, descriptions, stylesheet names that match their content's hash, and that the Atom feed and sitemap `make dist` generates list only published pages) on every push and pull request, and weekly. A second job lints and tests `vitae`, validates the résumé strictly and lints the workflow (`make workflow-lint`, actionlint). The site build renders every public variant itself: the résumé page shows `resume.full.json`, and `/resume/` serves each rendered file for download. It also checks that each résumé PDF reads back as its text, every passage intact and in order, the way an applicant-tracking system extracts it (`make resume-text-check`, with pdftotext). A push to `main` publishes to Cloudflare; a pull request gets a preview version, with its URL commented on the pull request.

Every run keeps the résumé files (PDF, Word, plain text, Markdown and JSON for each variant, under the names the site's downloads save as) as an artifact named `resume-<commit>`, for 30 days. To keep a version for good, such as one you sent, tag the commit and push the tag; CI attaches the same files to a GitHub release:

```bash
git tag resume-v2026.11.22 && git push origin resume-v2026.11.22
```

The Cloudflare steps run once two repository secrets exist:

- `CLOUDFLARE_ACCOUNT_ID`: the account ID from the Cloudflare dashboard.
- `CLOUDFLARE_API_TOKEN`: an API token made from the "Edit Cloudflare Workers" template, limited to this account.

The Worker serves the bare domain, `jonreyna.dev` (a custom domain in `wrangler.jsonc`). `www.jonreyna.dev` is set up in the Cloudflare dashboard, not in this repository: a proxied `AAAA` record for `www` pointing at the placeholder address `100::`, and a Redirect Rule matching `http*://www.jonreyna.dev/*` that sends a 301 to `https://jonreyna.dev/${2}`, keeping the query string. Requests to `www` never reach the Worker.

## License

- **Writing** (everything in `site/content/` and `resume/`: essays, page text and the résumé, and the text of the Play pieces in `play/`, such as Arbor Ecclesiae's descriptions): © Jon Reyna, licensed under [Creative Commons Attribution 4.0](https://creativecommons.org/licenses/by/4.0/) ([`LICENSE-CONTENT`](LICENSE-CONTENT)). Reuse it with credit and a link.
- **Fonts** in `site/web/fonts/`: their own licenses, kept beside them. The font files themselves aren't in git; `make fonts` downloads each from a pinned source, checks its SHA-256 and subsets it. All four families are under the SIL Open Font License 1.1: Alegreya ([`Alegreya-OFL.txt`](site/web/fonts/Alegreya-OFL.txt)), Alegreya SC ([`AlegreyaSC-OFL.txt`](site/web/fonts/AlegreyaSC-OFL.txt)), Courier Prime ([`CourierPrime-OFL.txt`](site/web/fonts/CourierPrime-OFL.txt)) and Libertinus Math ([`LibertinusMath-OFL.txt`](site/web/fonts/LibertinusMath-OFL.txt)).
- **D3**, which Arbor Ecclesiae draws with: ISC, by Mike Bostock ([`D3-LICENSE.txt`](play/arbor-ecclesiae/D3-LICENSE.txt), served beside it). `make play` downloads the pinned release and checks its SHA-256; it isn't in git.
- **Everything else** (code, styles, build files, `DESIGN.md`): [MIT](LICENSE).
