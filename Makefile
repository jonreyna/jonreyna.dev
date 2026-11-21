# Same targets locally and in CI.
#
# Jaspr 0.23.5 needs Dart 3.13 (see site/pubspec.yaml). On a Mac that is Homebrew's
# dart@3.13.7, and `jaspr` refuses Homebrew's bin/dart symlink, so the real SDK goes
# first on PATH. CI sets up Dart itself and DART_HOME stays empty there.
DART_HOME ?= $(firstword $(wildcard /usr/local/opt/dart@3.13.7/libexec /opt/homebrew/opt/dart@3.13.7/libexec))
export PATH := $(if $(DART_HOME),$(DART_HOME)/bin:)$(PATH):$(HOME)/.pub-cache/bin

.PHONY: help setup site dist serve check deploy clean

help:  ## List targets
	@grep -E '^[a-z]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-8s %s\n", $$1, $$2}'

setup:  ## Fetch Dart packages and the Jaspr CLI
	cd site && dart pub get
	dart pub global activate jaspr_cli 0.23.5

site:  ## Build the static site into site/build/jaspr
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

deploy: dist  ## Upload dist/ to Cloudflare (needs `wrangler login` or CLOUDFLARE_API_TOKEN)
	npx --yes wrangler@4 deploy

clean:  ## Remove build output
	rm -rf site/build dist .wrangler
