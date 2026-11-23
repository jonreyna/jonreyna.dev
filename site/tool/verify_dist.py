"""Checks the staged site in dist/ before it is published.

Run by `make verify`. Fails (exit 1) if any page:
  - loads a script other than Jaspr's island bundle (DESIGN.md, principle 4),
  - loads the island bundle without containing an island,
  - lacks a <title>, a lang attribute, or a meta description (the 404 page is exempt
    from the description),
or if a file the deploy needs is missing (including the résumé page and its JSON), or
if the feed or sitemap is malformed or lists a URL that isn't a published page (one
staged and not marked noindex), or the sitemap misses a published page.
"""

import pathlib
import re
import sys
import xml.etree.ElementTree as ET

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

for required in ("index.html", "404.html", "_headers", "robots.txt", "feed.xml", "sitemap.xml", "resume/index.html", "resume/resume.full.json", "resume/resume.full.pdf",
                 "resume/resume.full.docx", "resume/resume.full.txt", "resume/resume.full.md"):
    if not (DIST / required).is_file():
        problems.append(f"missing {required}")

# Published pages by canonical URL: staged, not the 404 page, not marked noindex.
SITE = "https://jonreyna.dev"
published = set()
for page in DIST.rglob("*.html"):
    html = page.read_text(encoding="utf-8")
    if page.name == "404.html" or re.search(r'<meta name="robots" content="noindex"', html):
        continue
    url = page.relative_to(DIST).as_posix().removesuffix("index.html").removesuffix("/")
    published.add(f"{SITE}/{url}" if url else f"{SITE}/")


def urls(name, path):
    """Returns the URLs an XML file lists at path (an ElementTree path), or None."""
    try:
        root = ET.parse(DIST / name).getroot()
    except (OSError, ET.ParseError) as e:
        problems.append(f"{name}: {e}")
        return None
    return [el.text or el.get("href") for el in root.iterfind(path)]


sitemap = urls("sitemap.xml", "{http://www.sitemaps.org/schemas/sitemap/0.9}url/{http://www.sitemaps.org/schemas/sitemap/0.9}loc")
if sitemap is not None:
    for url in sorted(set(sitemap) - published):
        problems.append(f"sitemap.xml: lists {url}, which is not a published page")
    for url in sorted(published - set(sitemap)):
        problems.append(f"sitemap.xml: misses {url}")
feed = urls("feed.xml", "{http://www.w3.org/2005/Atom}entry/{http://www.w3.org/2005/Atom}link[@rel='alternate']")
if feed is not None:
    for url in feed:
        if url not in published:
            problems.append(f"feed.xml: lists {url}, which is not a published page")

if problems:
    print("\n".join(problems))
    sys.exit(1)
print(f"verify: {len(list(DIST.rglob('*.html')))} pages OK; sitemap {len(sitemap)} pages, feed {len(feed)} essays")
