"""Checks the staged site in dist/ before it is published.

Run by `make verify` as `verify_dist.py DIST FONT_FACES`, where FONT_FACES is the
font manifest site/tool/subset_fonts.py wrote. Fails (exit 1) if any page:
  - loads a script other than Jaspr's island bundle (DESIGN.md, principle 4),
  - loads the island bundle without containing an island,
  - lacks a <title>, a lang attribute, or a meta description (the 404 page is exempt
    from the description), or an og:image,
  - emits a <base> tag (which breaks in-page #fragment links on subpages),
  - doesn't load the fonts stylesheet or preload each of the manifest's preload fonts,
  - links a stylesheet outside /css/, or one whose name's hash doesn't match its
    content (site/tool/fingerprint_css.py), so the year-long cache could serve it stale,
  - nests a link inside another link,
  - doesn't link the three favicons,
  - shows a character the face that draws it lacks: text in <math> needs the math
    face, text in <code> or <pre> every code face, and other body text every text face,
or if a file the deploy needs is missing (including the résumé page and its JSON, and
every font in the manifest), if any staged file still mentions Google Fonts, or if the
feed or sitemap is malformed or lists a URL that isn't a published page (one staged and
not marked noindex), or the sitemap misses a published page.
"""

import hashlib
import html.parser
import json
import pathlib
import re
import sys
import xml.etree.ElementTree as ET

DIST = pathlib.Path(sys.argv[1] if len(sys.argv) > 1 else "dist")
FONT_FACES = pathlib.Path(sys.argv[2] if len(sys.argv) > 2 else "build/fonts/faces.json")
ISLAND_JS = "/main.client.dart.js"


# Symbols no pinned font has: ⁂ (asterism before notes), ↩ (footnote back-links),
# ▮ (end of algorithm) and ▸ (recommended exercise). Browsers draw them from the
# system's symbol font, as they did under Google Fonts.
FALLBACK = {"\u2042", "\u21a9", "\u25ae", "\u25b8"}


class TextByFace(html.parser.HTMLParser):
    """Collects a page's visible characters by the font group that draws them, and
    counts links nested inside other links."""

    HIDDEN = {"head", "script", "style", "template"}  # text no web font draws
    GROUPS = (("math", {"math"}), ("code", {"code", "pre"}))  # checked in this order

    def __init__(self):
        super().__init__()
        self.depth = {tag: 0 for tag in self.HIDDEN | {"math", "code", "pre"}}
        self.chars = {"text": set(), "code": set(), "math": set()}
        self.links = 0  # how many <a> elements are open
        self.nested_links = 0  # <a> elements opened inside another

    def handle_starttag(self, tag, attrs):
        if tag in self.depth:
            self.depth[tag] += 1
        if tag == "a":
            self.nested_links += self.links > 0
            self.links += 1

    def handle_endtag(self, tag):
        if self.depth.get(tag):
            self.depth[tag] -= 1
        if tag == "a" and self.links:
            self.links -= 1

    def handle_data(self, data):
        if any(self.depth[tag] for tag in self.HIDDEN):
            return
        group = next((g for g, tags in self.GROUPS if any(self.depth[t] for t in tags)), "text")
        self.chars[group].update(set(data) - set("\n\r\t"))


def load_fonts():
    """Returns the font manifest and, per group, the code points every face in it has."""
    manifest = json.loads(FONT_FACES.read_text(encoding="utf-8"))
    covered = {}
    for face in manifest["faces"]:
        cps = {chr(cp) for first, last in face["codepoints"] for cp in range(first, last + 1)}
        group = face["group"]
        covered[group] = covered[group] & cps if group in covered else cps
    return manifest, covered


try:
    fonts, covered = load_fonts()
except (OSError, ValueError, KeyError) as e:
    sys.exit(f"verify: can't read the font manifest {FONT_FACES} (run make fonts): {e}")

problems = []
linked_css = set()
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
    if re.search(r"<base\b", html):
        problems.append(f"{rel}: has a <base> tag, which breaks in-page #fragment links on subpages")
    if not re.search(r"<title>[^<]+</title>", html):
        problems.append(f"{rel}: no <title>")
    if not re.search(r"<html[^>]*\blang=", html):
        problems.append(f"{rel}: no lang attribute")
    if rel.name != "404.html" and '<meta name="description"' not in html:
        problems.append(f"{rel}: no meta description")
    if 'property="og:image"' not in html:
        problems.append(f"{rel}: no og:image")
    links = [dict(re.findall(r'([\w-]+)="([^"]*)"', tag)) for tag in re.findall(r"<link\b[^>]*>", html)]
    sheets = [a.get("href", "") for a in links if a.get("rel") == "stylesheet"]
    if not any(re.fullmatch(r"/css/fonts\.[0-9a-f]{10}\.css", href) for href in sheets):
        problems.append(f"{rel}: doesn't load the fonts stylesheet")
    for href in sheets:
        linked_css.add(href)
        m = re.fullmatch(r"/css/[\w-]+\.([0-9a-f]{10})\.css", href)
        if not m:
            problems.append(f"{rel}: stylesheet {href} isn't fingerprinted")
        elif not (DIST / href.lstrip("/")).is_file():
            problems.append(f"{rel}: links {href}, which isn't staged")
        elif hashlib.sha256((DIST / href.lstrip("/")).read_bytes()).hexdigest()[:10] != m.group(1):
            problems.append(f"{rel}: {href} doesn't match its content's hash")
    for url in fonts["preload"]:
        if not any(a.get("rel") == "preload" and a.get("href") == url and a.get("as") == "font" for a in links):
            problems.append(f"{rel}: doesn't preload {url}")
    for want in ({"rel": "icon", "href": "/favicon.ico"}, {"rel": "icon", "href": "/favicon.svg"},
                 {"rel": "apple-touch-icon", "href": "/apple-touch-icon.png"}):
        if not any(want.items() <= a.items() for a in links):
            problems.append(f"{rel}: doesn't link {want['href']}")
    text = TextByFace()
    text.feed(html)
    if text.nested_links:
        problems.append(f"{rel}: {text.nested_links} links nested inside another link")
    for group, chars in sorted(text.chars.items()):
        if missing := sorted(chars - covered[group] - FALLBACK):
            shown = " ".join(f"U+{ord(c):04X} {c!r}" for c in missing[:8])
            problems.append(f"{rel}: characters the {group} fonts lack: {shown}" + (" …" if len(missing) > 8 else ""))

# Every stylesheet is fingerprinted under /css/ and linked; a stray one would be cached
# without a hash, or shipped for nothing.
for css in sorted(DIST.rglob("*.css")):
    url = "/" + css.relative_to(DIST).as_posix()
    if url not in linked_css:
        problems.append(f"{url.lstrip('/')}: staged stylesheet no page links (or not under /css/)")

# The Open Graph image (tools/card) must be the 1200×630 PNG the pages announce.
card = DIST / "og.png"
if card.is_file():
    head = card.read_bytes()[:24]
    size = (int.from_bytes(head[16:20], "big"), int.from_bytes(head[20:24], "big"))
    if head[:8] != b"\x89PNG\r\n\x1a\n" or size != (1200, 630):
        problems.append(f"og.png: not a 1200×630 PNG (got {size})")

# The fonts are self-hosted: every face in the manifest is staged, and nothing still
# points at Google Fonts (the pages, the CSS or the security policy in _headers).
for face in fonts["faces"]:
    if not (DIST / face["file"].lstrip("/")).is_file():
        problems.append(f"missing {face['file'].lstrip('/')}")
for path in sorted(DIST.rglob("*")):
    if path.suffix in {".html", ".css", ".xml", ".txt", ""} and path.is_file():
        if re.search(r"fonts\.(googleapis|gstatic)\.com", path.read_text(encoding="utf-8", errors="replace")):
            problems.append(f"{path.relative_to(DIST)}: mentions Google Fonts")

# The favicons (make favicon): an SVG drawn from the font's outline, so it names no
# font; an ICO; and the 180×180 touch icon.
svg = DIST / "favicon.svg"
if svg.is_file():
    text = svg.read_text(encoding="utf-8")
    if not text.startswith("<svg ") or "<path " not in text or "<text" in text or "font-family" in text:
        problems.append("favicon.svg: not the outlined letter (make favicon)")
ico = DIST / "favicon.ico"
if ico.is_file() and ico.read_bytes()[:4] != b"\x00\x00\x01\x00":
    problems.append("favicon.ico: not an ICO file")
touch = DIST / "apple-touch-icon.png"
if touch.is_file():
    head = touch.read_bytes()[:24]
    size = (int.from_bytes(head[16:20], "big"), int.from_bytes(head[20:24], "big"))
    if head[:8] != b"\x89PNG\r\n\x1a\n" or size != (180, 180):
        problems.append(f"apple-touch-icon.png: not a 180×180 PNG (got {size})")

for required in ("index.html", "favicon.svg", "favicon.ico", "apple-touch-icon.png", "404.html", "_headers", "robots.txt", "og.png", "feed.xml", "sitemap.xml", "resume/index.html", "resume/resume.full.json", "resume/resume.full.pdf",
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
print(f"verify: {len(list(DIST.rglob('*.html')))} pages OK; sitemap {len(sitemap)} pages, feed {len(feed)} essays; "
      f"{len(fonts['faces'])} fonts")
