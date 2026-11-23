"""Checks the staged site in dist/ before it is published.

Run by `make verify`. Fails (exit 1) if any page:
  - loads a script other than Jaspr's island bundle (DESIGN.md, principle 4),
  - loads the island bundle without containing an island,
  - lacks a <title>, a lang attribute, or a meta description (the 404 page is exempt
    from the description),
or if a file the deploy needs is missing (including the résumé page and its JSON).
"""

import pathlib
import re
import sys

DIST = pathlib.Path(sys.argv[1] if len(sys.argv) > 1 else "dist")
ISLAND_JS = "/main.client.dart.js"

problems = []
for page in sorted(DIST.rglob("*.html")):
    rel = page.relative_to(DIST)
    html = page.read_text(encoding="utf-8")
    scripts = re.findall(r"<script\b[^>]*>", html)
    for tag in scripts:
        if ISLAND_JS not in tag:
            problems.append(f"{rel}: script other than the island bundle: {tag}")
    # Jaspr marks each pre-rendered island with a comment such as <!--@kalman_demo-->.
    if any(ISLAND_JS in tag for tag in scripts) and "<!--@" not in html:
        problems.append(f"{rel}: loads {ISLAND_JS} but has no island")
    if not re.search(r"<title>[^<]+</title>", html):
        problems.append(f"{rel}: no <title>")
    if not re.search(r"<html[^>]*\blang=", html):
        problems.append(f"{rel}: no lang attribute")
    if rel.name != "404.html" and '<meta name="description"' not in html:
        problems.append(f"{rel}: no meta description")

for required in ("index.html", "404.html", "_headers", "resume/index.html", "resume/resume.full.json", "resume/resume.full.pdf",
                 "resume/resume.full.docx", "resume/resume.full.txt", "resume/resume.full.md"):
    if not (DIST / required).is_file():
        problems.append(f"missing {required}")

if problems:
    print("\n".join(problems))
    sys.exit(1)
print(f"verify: {len(list(DIST.rglob('*.html')))} pages OK")
