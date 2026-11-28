"""Gives the staged stylesheets content-hashed names, so browsers can keep them a year.

Run by `make dist` as `fingerprint_css.py DIST`. Moves each stylesheet in DIST's root
(fonts.css, tokens.css, site.css) to DIST/css/NAME.HASH.css, where HASH is the first ten
hex digits of its SHA-256 (the scheme the fonts use), and points every page's
<link href="/NAME.css"> at the new name. A changed stylesheet gets a new URL, so _headers
marks /css/* immutable and a deploy reaches returning visitors on their next page load;
under a fixed name with an hour's cache, they kept the old styles for up to an hour.

Fails (exit 1) if a page links a stylesheet that isn't staged, if a staged stylesheet
isn't linked by any page, or if a stylesheet has a relative url(), which would break
once the file moves to /css/.
"""

import hashlib
import pathlib
import re
import sys

DIST = pathlib.Path(sys.argv[1] if len(sys.argv) > 1 else "dist")

# A stylesheet <link> pointing at a file in the site's root, such as href="/site.css".
LINK = re.compile(r'(<link\b[^>]*\bhref=")/([\w-]+\.css)(")')

problems = []
renamed = {}
for css in sorted(DIST.glob("*.css")):
    text = css.read_text(encoding="utf-8")
    # A quoted data: URL can hold url() references of its own (an SVG filter's
    # url(%23n)); they point inside it, so leave them out of the check.
    outer = re.sub(r"url\(\s*([\"'])data:.*?\1\s*\)", "url(data:)", text)
    for url in re.findall(r"url\(\s*['\"]?([^'\")]+)", outer):
        if not url.startswith(("/", "data:", "#")):
            problems.append(f"{css.name}: relative url({url}) would break under /css/")
    digest = hashlib.sha256(css.read_bytes()).hexdigest()[:10]
    renamed[css.name] = f"/css/{css.stem}.{digest}.css"

pages = sorted(DIST.rglob("*.html"))
linked = set()
rewritten = {}
for page in pages:
    html = page.read_text(encoding="utf-8")

    def swap(m, page=page):
        name = m.group(2)
        if name not in renamed:
            problems.append(f"{page.relative_to(DIST)}: links /{name}, which isn't staged")
            return m.group(0)
        linked.add(name)
        return m.group(1) + renamed[name] + m.group(3)

    new = LINK.sub(swap, html)
    if new != html:
        rewritten[page] = new
for name in sorted(set(renamed) - linked):
    problems.append(f"{name}: staged but no page links it")
if problems:
    print("\n".join(problems))
    sys.exit(1)

(DIST / "css").mkdir(exist_ok=True)
for name, url in renamed.items():
    (DIST / name).rename(DIST / url.lstrip("/"))
for page, html in rewritten.items():
    page.write_text(html, encoding="utf-8")
print(f"fingerprinted {len(renamed)} stylesheets in {len(rewritten)} pages: "
      + ", ".join(url.removeprefix("/css/") for url in renamed.values()))
