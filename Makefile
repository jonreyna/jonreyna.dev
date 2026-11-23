# Same targets locally and in CI.
#
# Jaspr 0.23.5 needs Dart 3.13 (see site/pubspec.yaml). On a Mac that is Homebrew's
# dart@3.13.7, and `jaspr` refuses Homebrew's bin/dart symlink, so the real SDK goes
# first on PATH. CI sets up Dart itself and DART_HOME stays empty there.
DART_HOME ?= $(firstword $(wildcard /usr/local/opt/dart@3.13.7/libexec /opt/homebrew/opt/dart@3.13.7/libexec))
export PATH := $(if $(DART_HOME),$(DART_HOME)/bin:)$(PATH):$(HOME)/.pub-cache/bin

# The last commit's time (Unix seconds). Passed as SOURCE_DATE_EPOCH so generated files
# (PDFs, Word files, the feed) are the same on every build of a commit.
COMMIT_TIME := $(shell git log -1 --format=%ct 2>/dev/null || echo 0)

# Wrangler, pinned exactly; bump it on purpose.
WRANGLER ?= npx --yes wrangler@4.158.0

.PHONY: help setup no-binaries fonts card og-image site dist serve check site-test verify deploy preview clean vitae workflow-lint go-lint go-test resume-check resume resume-pdf resume-word resume-text-check resume-release site-resume resume-private

help:  ## List targets
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-15s %s\n", $$1, $$2}'

setup:  ## Fetch Dart packages and the Jaspr CLI, and turn on the git hooks
	cd site && dart pub get
	dart pub global activate jaspr_cli 0.23.5
	git config core.hooksPath .githooks

# Binary files never go in git (.githooks/pre-commit refuses them; this catches any
# that got past it). Git's own test decides what is binary.
no-binaries:  ## Fail if git tracks a binary file
	@bins=$$(git diff --numstat $$(git hash-object -t tree /dev/null) HEAD | awk '$$1 == "-" { print $$3 }'); \
	if [ -n "$$bins" ]; then echo "binary files tracked in git:"; echo "$$bins"; exit 1; fi; \
	echo "checked tracked files: no binaries"

# Fonts are binaries, so they stay out of git: make downloads each from a pinned
# release, checks its SHA-256 and unpacks it into site/web/fonts/ (git-ignored).
LIBERTINUS_VERSION := 7.051
LIBERTINUS_SHA256 := 4d9be29b5cb380c35af8ba967abcc752ad1e07be1f738a9789c33e0dd7478c92
LIBERTINUS_ZIP := build/downloads/Libertinus-$(LIBERTINUS_VERSION).zip
MATH_FONT := site/web/fonts/LibertinusMath-Regular.woff2

$(LIBERTINUS_ZIP):
	mkdir -p $(dir $@)
	curl -sSfL -o $@.part https://github.com/alerque/libertinus/releases/download/v$(LIBERTINUS_VERSION)/Libertinus-$(LIBERTINUS_VERSION).zip
	echo "$(LIBERTINUS_SHA256)  $@.part" | shasum -a 256 -c -
	mv $@.part $@

$(MATH_FONT): $(LIBERTINUS_ZIP)
	unzip -o -q -j $< Libertinus-$(LIBERTINUS_VERSION)/static/WOFF2/LibertinusMath-Regular.woff2 -d $(dir $@)
	touch $@

fonts: $(MATH_FONT)  ## Download the pinned fonts into site/web/fonts (not in git)

# The Open Graph image (the picture beside a shared link) is drawn by tools/card in
# the masthead's typefaces, Alegreya SC Medium and Alegreya Italic, downloaded from one
# pinned google/fonts commit and checked like Libertinus. They're used only to draw
# the image, not served. site/web/og.png is git-ignored.
GOOGLE_FONTS_COMMIT := 06bfe0d62a98375cb09c942fa18e68bd32cff057
GOOGLE_FONTS := https://raw.githubusercontent.com/google/fonts/$(GOOGLE_FONTS_COMMIT)/ofl
ALEGREYA_SC := build/downloads/fonts/AlegreyaSC-Medium.ttf
ALEGREYA_SC_SHA256 := fa896222115cc9a41c0695a3ee6feafbaba97932e9bb6240520b3feeb560ea2c
ALEGREYA_ITALIC := build/downloads/fonts/Alegreya-Italic.ttf
ALEGREYA_ITALIC_SHA256 := fa915eec76227935dc5fb678953c94b71287c360928013cfdb441dfe52f5a391
CARD := build/bin/card
OG_IMAGE := site/web/og.png

$(ALEGREYA_SC):
	mkdir -p $(dir $@)
	curl -sSfL -o $@.part $(GOOGLE_FONTS)/alegreyasc/AlegreyaSC-Medium.ttf
	echo "$(ALEGREYA_SC_SHA256)  $@.part" | shasum -a 256 -c -
	mv $@.part $@

$(ALEGREYA_ITALIC):
	mkdir -p $(dir $@)
	curl -sSfL -o $@.part '$(GOOGLE_FONTS)/alegreya/Alegreya-Italic%5Bwght%5D.ttf'
	echo "$(ALEGREYA_ITALIC_SHA256)  $@.part" | shasum -a 256 -c -
	mv $@.part $@

card:  ## Build the card CLI (the Open Graph image) into build/bin
	cd tools/card && go build -o ../../$(CARD) ./cmd/card

og-image: card $(ALEGREYA_SC) $(ALEGREYA_ITALIC)  ## Draw the Open Graph image into site/web/og.png (not in git)
	$(CARD) -smallcaps $(ALEGREYA_SC) -italic $(ALEGREYA_ITALIC) -name 'Jon Reyna' \
		-motto 'distributed systems & AI infrastructure' -site jonreyna.dev -o $(OG_IMAGE)

# PANDOC converts the pages' TeX math to MathML and highlights their code
# (site/lib/tex_math.dart, site/lib/code_highlight.dart). It is
# defined further down, so it can't be a prerequisite here; site-resume's Word step
# downloads it.
site: fonts og-image site-resume  ## Build the static site into site/build/jaspr
	cd site && PANDOC=$(CURDIR)/$(PANDOC) jaspr build

# Copy only the site's own files: the build also leaves packages/, .dart_tool/ and a manifest.
dist: site  ## Stage the deployable files in dist/
	rm -rf dist
	rsync -a --exclude='packages/' --exclude='.dart_tool/' --exclude='.build.manifest' site/build/jaspr/ dist/
	cp site/web/_headers dist/_headers
	# Cloudflare's 404-page handling looks for 404.html; Jaspr writes 404/index.html.
	mv dist/404/index.html dist/404.html && rmdir dist/404
	# The Atom feed and the sitemap, from the staged pages (site/tool/feed_sitemap.py).
	SOURCE_DATE_EPOCH=$(COMMIT_TIME) python3 -I site/tool/feed_sitemap.py dist
	@echo "dist: $$(find dist -type f | wc -l | tr -d ' ') files, $$(du -sh dist | cut -f1)"

serve: dist  ## Serve dist/ on http://localhost:8767
	python3 -m http.server 8767 --directory dist

check:  ## Static analysis
	cd site && dart analyze

# The site's Dart tests. PANDOC lets the math and highlighting tests run end to end; it comes
# from site-resume's Word step, like `make site`.
site-test: resume-word  ## Run the site's Dart tests (TeX to MathML, code highlighting)
	cd site && PANDOC=$(CURDIR)/$(PANDOC) dart test

verify: dist-exists  ## Check the staged pages (scripts, titles, descriptions) and the feed and sitemap
	python3 site/tool/verify_dist.py dist

# deploy and preview upload whatever is in dist/, so CI builds once and the
# same files are what you checked. Run `make dist` first.
deploy: dist-exists  ## Publish dist/ to production (needs `wrangler login` or CLOUDFLARE_API_TOKEN)
	$(WRANGLER) deploy

# A preview version with its own URL; production is untouched. ALIAS names the URL.
ALIAS ?= local
preview: dist-exists  ## Upload dist/ as a preview version at <ALIAS>-jonreyna-dev.<account>.workers.dev
	$(WRANGLER) versions upload --preview-alias $(ALIAS)

.PHONY: dist-exists
dist-exists:
	@test -f dist/index.html || { echo "dist/ is missing: run make dist first"; exit 1; }

# ---- Résumé (tools/vitae, resume/resume.yaml) ----
VITAE := build/bin/vitae

# golangci-lint, pinned exactly and installed per version, so a bump reinstalls it.
GOLANGCI_LINT_VERSION := v2.15.0
GOLANGCI_LINT := build/tools/golangci-lint-$(GOLANGCI_LINT_VERSION)/golangci-lint

$(GOLANGCI_LINT):
	GOBIN=$(CURDIR)/$(dir $(GOLANGCI_LINT)) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

# actionlint checks the GitHub workflows (expressions, action inputs, shell). Pinned and
# built from source like golangci-lint.
ACTIONLINT_VERSION := v1.7.12
ACTIONLINT := build/tools/actionlint-$(ACTIONLINT_VERSION)/actionlint

$(ACTIONLINT):
	GOBIN=$(CURDIR)/$(dir $(ACTIONLINT)) go install github.com/rhysd/actionlint/cmd/actionlint@$(ACTIONLINT_VERSION)

workflow-lint: $(ACTIONLINT)  ## Lint the GitHub workflows with actionlint
	$(ACTIONLINT) .github/workflows/*.yml
	@echo "linted .github/workflows: no issues"

# Go modules under tools/, each linted and tested the same way.
GO_MODULES := tools/vitae tools/card

vitae:  ## Build the vitae CLI into build/bin
	cd tools/vitae && go build -o ../../$(VITAE) ./cmd/vitae

go-lint: $(GOLANGCI_LINT)  ## Check formatting and lint every Go module in tools/ (tools/.golangci.yml)
	for m in $(GO_MODULES); do \
		(cd $$m && $(CURDIR)/$(GOLANGCI_LINT) fmt --diff ./... && $(CURDIR)/$(GOLANGCI_LINT) run ./...) || exit 1; \
	done

go-test:  ## Vet and test every Go module in tools/
	for m in $(GO_MODULES); do (cd $$m && go vet ./... && go test ./...) || exit 1; done

resume-check: vitae go-lint go-test  ## Lint and test the Go tools, and validate resume.yaml strictly
	$(VITAE) validate -strict -i resume/resume.yaml

resume: vitae  ## Render every variant and format into build/resume (public fields only)
	rm -rf build/resume
	$(VITAE) render -i resume/resume.yaml -o build/resume

# PDFs: latexmk and pdfLaTeX turn each resume.<variant>.tex into a PDF beside it. The
# design is resume/latex/jrresume.cls, found through TEXINPUTS. Auxiliary files go to
# build/latex/. SOURCE_DATE_EPOCH (the last commit's time) keeps the PDFs reproducible.
LATEXMK ?= latexmk
define build-pdfs
	for tex in $(1)/resume.*.tex; do \
		name=$$(basename $$tex .tex); \
		SOURCE_DATE_EPOCH=$(COMMIT_TIME) FORCE_SOURCE_DATE=1 TEXINPUTS=$(CURDIR)/resume/latex//: \
			$(LATEXMK) -pdf -silent -interaction=nonstopmode -halt-on-error \
			-outdir=$(CURDIR)/build/latex/$(notdir $(1)) $$tex || exit 1; \
		cp build/latex/$(notdir $(1))/$$name.pdf $(1)/$$name.pdf; \
		echo "wrote $(1)/$$name.pdf"; \
	done
endef

resume-pdf: resume  ## Build a PDF of every variant into build/resume (needs TeX Live)
	$(call build-pdfs,build/resume)

# Word: Pandoc turns each resume.<variant>.md into a .docx, styled by a reference
# document. Pandoc is a pinned release binary, downloaded and checked like the fonts;
# the reference document is generated (resume/word/reference.py), since binaries
# never go in git.
PANDOC_VERSION := 3.13
ifeq ($(shell uname -s),Darwin)
  ifeq ($(shell uname -m),arm64)
    PANDOC_ASSET := pandoc-$(PANDOC_VERSION)-arm64-macOS.zip
    PANDOC_SHA256 := f436d18bebbfb3ddd326597686dacab9fac5dc23bd060495020f35a28a43a7c3
  else
    PANDOC_ASSET := pandoc-$(PANDOC_VERSION)-x86_64-macOS.zip
    PANDOC_SHA256 := 1823278715aa405aaf21b59cde5dc789346aa8d4fc307467f9a9ee84162bedc5
  endif
else
  ifeq ($(shell uname -m),aarch64)
    PANDOC_ASSET := pandoc-$(PANDOC_VERSION)-linux-arm64.tar.gz
    PANDOC_SHA256 := 21e1d4ba09bd0b334590d204580fd4f6629b21be4ab6a32b26cc98982e9e0c68
  else
    PANDOC_ASSET := pandoc-$(PANDOC_VERSION)-linux-amd64.tar.gz
    PANDOC_SHA256 := 38aa33b325a2bc96a4c39409631212db55702f564250bcdbead180b3bf646389
  endif
endif
PANDOC := build/tools/pandoc-$(PANDOC_VERSION)/pandoc
PANDOC_ARCHIVE := build/downloads/$(PANDOC_ASSET)
WORD_REFERENCE := build/word/reference.docx

$(PANDOC_ARCHIVE):
	mkdir -p $(dir $@)
	curl -sSfL -o $@.part https://github.com/jgm/pandoc/releases/download/$(PANDOC_VERSION)/$(PANDOC_ASSET)
	echo "$(PANDOC_SHA256)  $@.part" | shasum -a 256 -c -
	mv $@.part $@

$(PANDOC): $(PANDOC_ARCHIVE)
	mkdir -p $(dir $@)
	case $< in \
		*.zip) unzip -o -q -j $< '*/bin/pandoc' -d $(dir $@) ;; \
		*) tar -xzf $< -C $(dir $@) --strip-components=2 --wildcards '*/bin/pandoc' ;; \
	esac
	touch $@

$(WORD_REFERENCE): $(PANDOC) resume/word/reference.py
	mkdir -p $(dir $@)
	$(PANDOC) -o build/word/pandoc-reference.docx --print-default-data-file reference.docx
	python3 -I resume/word/reference.py build/word/pandoc-reference.docx $@

define build-word
	for md in $(1)/resume.*.md; do \
		name=$$(basename $$md .md); \
		SOURCE_DATE_EPOCH=$(COMMIT_TIME) $(PANDOC) -f commonmark -t docx \
			--reference-doc=$(WORD_REFERENCE) -o $(1)/$$name.docx $$md || exit 1; \
		echo "wrote $(1)/$$name.docx"; \
	done
endef

resume-word: resume $(WORD_REFERENCE)  ## Build a Word file of every variant into build/resume
	$(call build-word,build/resume)

# What a parser reads: pdftotext (Poppler) extracts each PDF, and every passage must
# come back intact and in order (resume/check_pdf_text.py). CI installs Poppler; on a
# Mac without it this target fails with a note, and the rest of the build doesn't need it.
resume-text-check:  ## Check that each résumé PDF extracts as its text, in order (needs pdftotext; run after resume-pdf)
	python3 -I resume/check_pdf_text.py build/resume

# The files CI keeps from each build and attaches to a release, under the names the
# site's downloads save as: Jon-Reyna-Resume.pdf for the full résumé,
# Jon-Reyna-Resume-<variant>.pdf for the others. The .tex sources stay out.
RESUME_RELEASE := build/resume-release
resume-release:  ## Stage the rendered résumé files under their download names in build/resume-release (run after resume-pdf and resume-word)
	rm -rf $(RESUME_RELEASE)
	mkdir -p $(RESUME_RELEASE)
	for f in build/resume/resume.*.json build/resume/resume.*.pdf build/resume/resume.*.docx \
		build/resume/resume.*.txt build/resume/resume.*.md; do \
		test -f $$f || { echo "$$f is missing: run make resume-pdf resume-word first"; exit 1; }; \
		variant=$$(basename $$f | cut -d. -f2); ext=$${f##*.}; \
		if [ $$variant = full ]; then name=Jon-Reyna-Resume.$$ext; else name=Jon-Reyna-Resume-$$variant.$$ext; fi; \
		cp $$f $(RESUME_RELEASE)/$$name; \
	done
	@echo "staged $(RESUME_RELEASE): $$(ls $(RESUME_RELEASE) | wc -l | tr -d ' ') files"

# The site's résumé page and downloads come from the same render (public fields only).
# site/web/resume/ is git-ignored; Jaspr copies it into the build, so the files are
# served at /resume/resume.<variant>.<format>. The .tex sources stay out.
site-resume: resume-pdf resume-word  ## Copy the rendered résumé into site/web/resume for the site build
	rm -rf site/web/resume
	mkdir -p site/web/resume
	cp build/resume/resume.*.json build/resume/resume.*.pdf build/resume/resume.*.docx \
		build/resume/resume.*.txt build/resume/resume.*.md site/web/resume/

# Adds the phone number from resume/resume.private.yaml. Local only: never in CI.
resume-private: vitae  ## Render with the private overlay into build/resume-private
	@test -f resume/resume.private.yaml || { echo "copy resume/resume.private.example.yaml to resume/resume.private.yaml first"; exit 1; }
	rm -rf build/resume-private
	$(VITAE) render -i resume/resume.yaml -private resume/resume.private.yaml -o build/resume-private
	$(call build-pdfs,build/resume-private)
	$(MAKE) --no-print-directory $(WORD_REFERENCE)
	$(call build-word,build/resume-private)

clean:  ## Remove build output, the copied résumé, the Open Graph image and downloaded fonts
	rm -rf site/build dist .wrangler build site/web/resume $(MATH_FONT) $(OG_IMAGE)
