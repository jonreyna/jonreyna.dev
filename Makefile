# Same targets locally and in CI.
#
# Jaspr 0.23.5 needs Dart 3.13 (see site/pubspec.yaml). On a Mac that is Homebrew's
# dart@3.13.7, and `jaspr` refuses Homebrew's bin/dart symlink, so the real SDK goes
# first on PATH. CI sets up Dart itself and DART_HOME stays empty there.
DART_HOME ?= $(firstword $(wildcard /usr/local/opt/dart@3.13.7/libexec /opt/homebrew/opt/dart@3.13.7/libexec))
export PATH := $(if $(DART_HOME),$(DART_HOME)/bin:)$(PATH):$(HOME)/.pub-cache/bin

# Wrangler, pinned exactly; bump it on purpose.
WRANGLER ?= npx --yes wrangler@4.158.0

.PHONY: help setup no-binaries fonts site dist serve check verify deploy preview clean vitae go-lint go-test resume-check resume resume-private

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

site: fonts  ## Build the static site into site/build/jaspr
	cd site && jaspr build

# Copy only the site's own files: the build also leaves packages/, .dart_tool/ and a manifest.
dist: site  ## Stage the deployable files in dist/
	rm -rf dist
	rsync -a --exclude='packages/' --exclude='.dart_tool/' --exclude='.build.manifest' site/build/jaspr/ dist/
	cp site/web/_headers dist/_headers
	# Cloudflare's 404-page handling looks for 404.html; Jaspr writes 404/index.html.
	mv dist/404/index.html dist/404.html && rmdir dist/404
	@echo "dist: $$(find dist -type f | wc -l | tr -d ' ') files, $$(du -sh dist | cut -f1)"

serve: dist  ## Serve dist/ on http://localhost:8767
	python3 -m http.server 8767 --directory dist

check:  ## Static analysis
	cd site && dart analyze

verify: dist-exists  ## Check the staged pages (scripts, titles, descriptions)
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

# Go modules under tools/, each linted and tested the same way.
GO_MODULES := tools/vitae

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

# Adds the phone number from resume/resume.private.yaml. Local only: never in CI.
resume-private: vitae  ## Render with the private overlay into build/resume-private
	@test -f resume/resume.private.yaml || { echo "copy resume/resume.private.example.yaml to resume/resume.private.yaml first"; exit 1; }
	rm -rf build/resume-private
	$(VITAE) render -i resume/resume.yaml -private resume/resume.private.yaml -o build/resume-private

clean:  ## Remove build output and downloaded fonts
	rm -rf site/build dist .wrangler build $(MATH_FONT)
