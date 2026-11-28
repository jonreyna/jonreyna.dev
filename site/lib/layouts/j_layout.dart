import 'package:jaspr/dom.dart';
import 'package:jaspr/server.dart';
import 'package:jaspr_content/jaspr_content.dart';

import '../dates.dart';
import '../essays.dart';
import '../fonts.dart';

/// Design J ("Red and black, with algorithms") as a jaspr_content layout.
/// See DESIGN.md (in site/) for the rules this follows.
class JLayout extends PageLayoutBase {
  const JLayout();

  @override
  Pattern get name => 'j';

  // PageLayoutBase.buildHead hard-codes "Title | Base", so this replaces it rather than
  // calling it. It covers the parts of the base that the site uses: title, description,
  // Open Graph title, description and image. (Not covered: `keywords` and `meta`.)
  @override
  // ignore: must_call_super
  Iterable<Component> buildHead(Page page) sync* {
    final data = page.data.page;
    final site = page.data.site;
    // "Title · Jon Reyna", or just the name on pages without a title (DESIGN.md).
    final base = '${site['titleBase'] ?? 'Jon Reyna'}';
    final title = data['title'] is String ? '${data['title']} · $base' : base;
    yield Component.element(tag: 'title', children: [.text(title)]);
    yield meta(attributes: {'property': 'og:title'}, content: title);
    yield meta(attributes: {'property': 'og:type'}, content: data['kind'] == 'essay' ? 'article' : 'website');
    if (data['description'] case final String description) {
      yield meta(name: 'description', content: description);
      yield meta(attributes: {'property': 'og:description'}, content: description);
    }
    // The feed and sitemap (tool/feed_sitemap.py) read an essay's dates from these.
    if (data['kind'] == 'essay' && data['date'] is String) {
      yield meta(attributes: {'property': 'article:published_time'}, content: data['date'] as String);
    }
    if (data['kind'] == 'essay' && data['updated'] is String) {
      yield meta(attributes: {'property': 'article:modified_time'}, content: data['updated'] as String);
    }
    // The picture beside a shared link: the page's own `image` if it sets one, or the
    // site's card (`make og-image`, tools/card), 1200×630.
    yield meta(attributes: {'property': 'og:site_name'}, content: base);
    if (data['image'] case final String image) {
      final url = image.startsWith('/') ? 'https://jonreyna.dev$image' : image;
      yield meta(attributes: {'property': 'og:image'}, content: url);
    } else {
      yield meta(attributes: {'property': 'og:image'}, content: 'https://jonreyna.dev/og.png');
      yield meta(attributes: {'property': 'og:image:width'}, content: '1200');
      yield meta(attributes: {'property': 'og:image:height'}, content: '630');
      yield meta(
        attributes: {'property': 'og:image:alt'},
        content: 'Jon Reyna: distributed systems & AI infrastructure. jonreyna.dev',
      );
    }
    yield meta(name: 'twitter:card', content: 'summary_large_image');
    // A draft is built (so it can be previewed) but, like a noindex page, stays out of
    // search engines, the essay lists, the feed and the sitemap.
    if (data['noindex'] == true || data['draft'] == true) {
      yield meta(name: 'robots', content: 'noindex');
    }
    // The essay index (`essayIndex: true`, /writing) stays out of search engines, the
    // sitemap and the nav until an essay is published, and the feed link waits too.
    yield Builder(
      builder: (context) {
        final hasEssays = publishedEssays(pageEntries(context.pages)).isNotEmpty;
        return Component.fragment([
          if (!hasEssays && data['essayIndex'] == true && data['noindex'] != true && data['draft'] != true)
            meta(name: 'robots', content: 'noindex'),
          if (hasEssays)
            link(rel: 'alternate', type: 'application/atom+xml', href: '/feed.xml', attributes: {'title': 'Jon Reyna'}),
        ]);
      },
    );
    // Self-hosted, subset fonts (make fonts). Every page sets text in Alegreya regular,
    // the italic (the tagline) and Alegreya SC (the name and nav), so those three are
    // fetched with the HTML rather than after the CSS; arriving late, they reflow the
    // page (DESIGN.md, font budget).
    for (final font in preloadFonts()) {
      yield link(rel: 'preload', href: font, attributes: {'as': 'font', 'type': 'font/woff2', 'crossorigin': ''});
    }
    // make dist renames each of these for its content, as /css/NAME.HASH.css, and
    // rewrites these links to match (tool/fingerprint_css.py).
    for (final href in stylesheets(page)) {
      yield link(rel: 'stylesheet', href: href);
    }
    // The SVG for browsers that take one; favicon.ico for the rest and for tools that
    // ask for /favicon.ico; the PNG for phones' home screens (make favicon).
    yield link(rel: 'icon', href: '/favicon.ico', attributes: {'sizes': '32x32'});
    yield link(rel: 'icon', href: '/favicon.svg', type: 'image/svg+xml');
    yield link(rel: 'apple-touch-icon', href: '/apple-touch-icon.png');
    yield meta(name: 'theme-color', content: '#f6f0e1', attributes: {'media': '(prefers-color-scheme: light)'});
    yield meta(name: 'theme-color', content: '#1a1612', attributes: {'media': '(prefers-color-scheme: dark)'});
    yield link(rel: 'canonical', href: 'https://jonreyna.dev${page.url}');
    yield meta(attributes: {'property': 'og:url', 'content': 'https://jonreyna.dev${page.url}'});
    for (final src in scripts(page)) {
      yield script(src: src, defer: true);
    }
  }

  /// The page's stylesheets, each a file in web/'s root: the fonts, the design
  /// tokens and the site's styles.
  List<String> stylesheets(Page page) => const ['/fonts.css', '/tokens.css', '/site.css'];

  /// The page's scripts. Reading pages ship no JavaScript; a page with an island opts
  /// in with `islands: true` (DESIGN.md, principle 4).
  List<String> scripts(Page page) => [if (page.data.page['islands'] == true) '/main.client.dart.js'];

  @override
  Component buildBody(Page page, Component child) {
    final data = page.data.page;
    final isHome = page.url == '/' || page.url.isEmpty;
    final name = a(href: '/', [.text('Jon Reyna')]);

    return Component.fragment([
      a(classes: 'skip', href: '#main', [.text('Skip to content')]),
      div(classes: 'wrap', [
        header(classes: 'masthead', [
          if (isHome) h1(classes: 'name', [name]) else p(classes: 'name', [name]),
          p(classes: 'motto', [.text('distributed systems & AI infrastructure')]),
          Builder(
            builder: (context) {
              final hasEssays = publishedEssays(pageEntries(context.pages)).isNotEmpty;
              return nav(
                attributes: {'aria-label': 'Site'},
                [
                  for (final (href, label) in _nav)
                    if (hasEssays || href != '/writing')
                      a(
                        href: href,
                        attributes: {if (_inSection(page.url, href)) 'aria-current': 'page'},
                        [.text(label)],
                      ),
                ],
              );
            },
          ),
        ]),
        main_(id: 'main', classes: data['kind'] == 'essay' ? 'essay' : null, [
          if (data['kind'] == 'essay')
            article([_essayHead(data), child, _essayFoot()])
          else ...[
            if (data['kind'] == 'page' && data['title'] is String)
              h1(classes: 'page-title', [.text(data['title'] as String)]),
            child,
          ],
        ]),
        footer(classes: 'site-foot', [
          a(href: '/about', [.text('about')]),
          a(href: 'mailto:jon@jonreyna.dev', [.text('jon@jonreyna.dev')]),
          a(href: 'https://github.com/jonreyna', [.text('GitHub')]),
          a(href: 'https://www.linkedin.com/in/jonreyna', [.text('LinkedIn')]),
          Builder(
            builder: (context) => publishedEssays(pageEntries(context.pages)).isEmpty
                ? Component.fragment([])
                : a(href: '/feed.xml', [.text('feed')]),
          ),
        ]),
      ]),
    ]);
  }

  /// The nav items (DESIGN.md, "Masthead"). Writing shows once an essay is published.
  /// Play goes to its section on the home page until it has a page of its own.
  static const _nav = [('/work', 'work'), ('/writing', 'writing'), ('/#play', 'play'), ('/resume', 'résumé')];

  /// Whether [url] is the nav item's page or one below it (an essay under /writing).
  static bool _inSection(String url, String href) => !href.contains('#') && (url == href || url.startsWith('$href/'));

  /// Ends every essay: who wrote it and the way back to the list.
  Component _essayFoot() => footer(classes: 'essay-foot', [
    p([
      .text('Written by Jon Reyna. '),
      a(href: '/writing', [.text('All writing')]),
      .text('.'),
    ]),
  ]);

  Component _essayHead(Map<String, Object?> data) {
    final date = '${data['date'] ?? ''}';
    final shown = formatMonth(date);
    return Component.fragment([
      header(classes: 'essay-head', [
        h1([.text('${data['title']}')]),
        p(classes: 'by', [.text('by Jon Reyna')]),
        p(classes: 'date', [
          Component.element(tag: 'time', attributes: {'datetime': date}, children: [.text(shown)]),
        ]),
      ]),
      if (data['abstract'] case final String abstract) p(classes: 'abstract', [.text(abstract)]),
    ]);
  }
}
