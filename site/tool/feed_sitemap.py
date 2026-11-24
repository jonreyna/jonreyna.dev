"""Writes the Atom feed and the sitemap for the staged site in DIR.

Usage: python3 -I feed_sitemap.py DIR

Run by `make dist` after the pages are staged, so both files describe exactly what
gets published. It reads the rendered pages rather than the Markdown, so the feed
carries each essay as the site shows it: section numbers, MathML and footnotes
included.

- DIR/feed.xml (Atom): one entry per published essay, newest first, with the essay's
  full text. An essay is a page with og:type "article"; a page marked noindex (which
  JLayout also sets for drafts) is left out. Dates come from article:published_time
  and article:modified_time, from the essay's `date` and `updated` front matter.
- DIR/sitemap.xml: every page except the 404 page and noindex pages, at its canonical
  URL. Essays carry a lastmod; other pages don't, since the build can't know when
  their content last changed, and a wrong date is worse than none.

The output is deterministic. With no essays, the feed's updated date is the last
commit's time, from SOURCE_DATE_EPOCH. Only the standard library is used.
"""

import datetime
import html
import os
import pathlib
import re
import sys
import xml.etree.ElementTree as ET

SITE = "https://jonreyna.dev"
TITLE = "Jon Reyna"
SUBTITLE = "Essays on AI infrastructure, distributed systems and the mathematics underneath them."
AUTHOR = "Jon Reyna"
ATOM = "http://www.w3.org/2005/Atom"
SITEMAP = "http://www.sitemaps.org/schemas/sitemap/0.9"


def tags(page, tag):
    """Yields the attributes of each <tag> in the page, in any order, as a dict."""
    for m in re.finditer(rf"<{tag}\b([^>]*)>", page):
        yield {k: html.unescape(v) for k, v in re.findall(r'([\w:-]+)="([^"]*)"', m.group(1))}


def meta(page, attr, name):
    """Returns the content of the page's <meta attr="name"> tag, or None."""
    return next((t.get("content") for t in tags(page, "meta") if t.get(attr) == name), None)


def is_noindex(page):
    """Returns whether the page's robots meta tag includes noindex."""
    robots = meta(page, "name", "robots") or ""
    return "noindex" in {d.strip().lower() for d in robots.split(",")}


def canonical(page, rel):
    """Returns the page's canonical URL; every page must have one."""
    href = next((t.get("href") for t in tags(page, "link") if t.get("rel") == "canonical"), None)
    if not href:
        sys.exit(f"feed_sitemap.py: {rel} has no canonical link")
    return href


def timestamp(value, what):
    """Turns a front-matter date ("2027-03", "2027-03-14" or a full RFC 3339 time)
    into an RFC 3339 time. A month means its first day; a day means midnight UTC."""
    if re.fullmatch(r"\d{4}-\d{2}", value):
        value += "-01"
    if re.fullmatch(r"\d{4}-\d{2}-\d{2}", value):
        value += "T00:00:00Z"
    try:
        datetime.datetime.fromisoformat(value.replace("Z", "+00:00"))
    except ValueError:
        sys.exit(f"feed_sitemap.py: {what}: {value!r} is not a date such as 2027-03 or 2027-03-14")
    return value


def essay_body(page, url):
    """Returns the essay's text as HTML: the article without its head (title, byline,
    date) and foot, with site links made absolute so they work in a feed reader."""
    m = re.search(r"<article>(.*)</article>", page, re.S)
    if not m:
        sys.exit(f"feed_sitemap.py: {url} is an essay but has no <article>")
    body = m.group(1)
    body = re.sub(r'<header class="essay-head">.*?</header>', "", body, count=1, flags=re.S)
    body = re.sub(r'<footer class="essay-foot">.*?</footer>', "", body, count=1, flags=re.S)
    body = re.sub(r'\b(href|src)="/(?!/)', rf'\1="{SITE}/', body)
    body = re.sub(r'\bhref="#', f'href="{url}#', body)
    return body.strip()


def essay_title(page, url):
    """Returns the essay's title as its heading shows it."""
    m = re.search(r'<header class="essay-head">\s*<h1>(.*?)</h1>', page, re.S)
    if not m:
        sys.exit(f"feed_sitemap.py: {url} is an essay but has no title heading")
    return html.unescape(re.sub(r"<[^>]+>", "", m.group(1))).strip()


def read_pages(dist):
    """Lists the pages to publish in the feed and sitemap, sorted by URL."""
    pages = []
    for path in sorted(dist.rglob("*.html")):
        rel = path.relative_to(dist).as_posix()
        page = path.read_text(encoding="utf-8")
        if rel == "404.html" or is_noindex(page):
            continue
        url = canonical(page, rel)
        entry = {"url": url}
        if meta(page, "property", "og:type") == "article":
            published = meta(page, "property", "article:published_time")
            if not published:
                sys.exit(f"feed_sitemap.py: {url} is an essay but has no date")
            entry["published"] = timestamp(published, f"{url} date")
            modified = meta(page, "property", "article:modified_time")
            entry["updated"] = timestamp(modified, f"{url} updated") if modified else entry["published"]
            entry["title"] = essay_title(page, url)
            entry["summary"] = meta(page, "name", "description")
            entry["content"] = essay_body(page, url)
        pages.append(entry)
    return sorted(pages, key=lambda p: p["url"])


def write_feed(dist, essays):
    """Writes the Atom feed of the essays, newest first."""
    essays = sorted(essays, key=lambda e: (e["published"], e["url"]), reverse=True)
    if essays:
        updated = max(e["updated"] for e in essays)
    else:
        epoch = int(os.environ.get("SOURCE_DATE_EPOCH", "0"))
        updated = datetime.datetime.fromtimestamp(epoch, datetime.UTC).strftime("%Y-%m-%dT%H:%M:%SZ")

    ET.register_namespace("", ATOM)

    def sub(parent, tag, text=None, **attrs):
        el = ET.SubElement(parent, f"{{{ATOM}}}{tag}", attrs)
        if text is not None:
            el.text = text
        return el

    feed = ET.Element(f"{{{ATOM}}}feed", {"xml:lang": "en"})
    sub(feed, "id", f"{SITE}/")
    sub(feed, "title", TITLE)
    sub(feed, "subtitle", SUBTITLE)
    sub(feed, "updated", updated)
    sub(sub(feed, "author"), "name", AUTHOR)
    sub(feed, "link", rel="self", type="application/atom+xml", href=f"{SITE}/feed.xml")
    sub(feed, "link", rel="alternate", type="text/html", href=f"{SITE}/writing")
    sub(feed, "icon", f"{SITE}/favicon.svg")
    sub(feed, "rights", "© Jon Reyna. Licensed under Creative Commons Attribution 4.0.")
    for e in essays:
        entry = sub(feed, "entry")
        sub(entry, "id", e["url"])
        sub(entry, "title", e["title"])
        sub(entry, "link", rel="alternate", type="text/html", href=e["url"])
        sub(entry, "published", e["published"])
        sub(entry, "updated", e["updated"])
        if e["summary"]:
            sub(entry, "summary", e["summary"])
        sub(entry, "content", e["content"], type="html")
    ET.indent(feed)
    (dist / "feed.xml").write_bytes(ET.tostring(feed, xml_declaration=True, encoding="UTF-8") + b"\n")
    return len(essays)


def write_sitemap(dist, pages):
    """Writes the sitemap: every published page, essays with their last change."""
    ET.register_namespace("", SITEMAP)
    urlset = ET.Element(f"{{{SITEMAP}}}urlset")
    for p in pages:
        url = ET.SubElement(urlset, f"{{{SITEMAP}}}url")
        ET.SubElement(url, f"{{{SITEMAP}}}loc").text = p["url"]
        if "updated" in p:
            ET.SubElement(url, f"{{{SITEMAP}}}lastmod").text = p["updated"]
    ET.indent(urlset)
    (dist / "sitemap.xml").write_bytes(ET.tostring(urlset, xml_declaration=True, encoding="UTF-8") + b"\n")
    return len(pages)


def main():
    if len(sys.argv) != 2:
        sys.exit("usage: feed_sitemap.py DIR")
    dist = pathlib.Path(sys.argv[1])
    if not (dist / "index.html").is_file():
        sys.exit(f"feed_sitemap.py: {dist} has no index.html; stage the site first")
    pages = read_pages(dist)
    essays = write_feed(dist, [p for p in pages if "content" in p])
    urls = write_sitemap(dist, pages)
    print(f"wrote {dist}/feed.xml: {essays} essays")
    print(f"wrote {dist}/sitemap.xml: {urls} pages")


if __name__ == "__main__":
    main()
