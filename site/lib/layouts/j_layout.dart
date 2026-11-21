import 'package:jaspr/dom.dart';
import 'package:jaspr/server.dart';
import 'package:jaspr_content/jaspr_content.dart';

/// Design J ("Red and black, with algorithms") as a jaspr_content layout.
/// See ../../../site-specimens/DESIGN.md for the rules this follows.
class JLayout extends PageLayoutBase {
  const JLayout();

  @override
  Pattern get name => 'j';

  static const _months = [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December',
  ];

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
    if (data['image'] case final String image) {
      yield meta(attributes: {'property': 'og:image'}, content: image);
    }
    if (data['noindex'] == true) {
      yield meta(name: 'robots', content: 'noindex');
    }
    // Google Fonts for now; self-hosted subset WOFF2 files replace this before launch.
    yield link(rel: 'preconnect', href: 'https://fonts.googleapis.com');
    yield link(rel: 'preconnect', href: 'https://fonts.gstatic.com', attributes: {'crossorigin': ''});
    yield link(
      rel: 'stylesheet',
      href: 'https://fonts.googleapis.com/css2?family=Alegreya:ital,wght@0,400;0,500;1,400'
          '&family=Alegreya+SC:wght@500&family=Courier+Prime:wght@400;700&display=swap',
    );
    yield link(rel: 'stylesheet', href: '/tokens.css');
    yield link(rel: 'stylesheet', href: '/site.css');
    yield link(rel: 'icon', href: '/favicon.svg', type: 'image/svg+xml');
    yield meta(name: 'theme-color', content: '#f6f0e1', attributes: {'media': '(prefers-color-scheme: light)'});
    yield meta(name: 'theme-color', content: '#1a1612', attributes: {'media': '(prefers-color-scheme: dark)'});
    yield link(rel: 'canonical', href: 'https://jonreyna.dev${page.url}');
    yield link(rel: 'alternate', type: 'application/atom+xml', href: '/feed.xml', attributes: {'title': 'Jon Reyna'});
    yield meta(attributes: {'property': 'og:url', 'content': 'https://jonreyna.dev${page.url}'});
    // Reading pages ship no JavaScript; a page with an island opts in (DESIGN.md, principle 4).
    if (data['islands'] == true) {
      yield script(src: '/main.client.dart.js', defer: true);
    }
  }

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
          nav(attributes: {'aria-label': 'Site'}, [
            a(href: '/#work', [.text('work')]),
            a(href: '/#writing', [.text('writing')]),
            a(href: '/#play', [.text('play')]),
            a(href: '/resume', [.text('résumé')]),
          ]),
        ]),
        main_(id: 'main', classes: data['kind'] == 'essay' ? 'essay' : null, [
          if (data['kind'] == 'essay') _essayHead(data),
          if (data['kind'] == 'page' && data['title'] is String) h1(classes: 'page-title', [.text(data['title'] as String)]),
          child,
        ]),
        footer(classes: 'site-foot', [
          a(href: '/about', [.text('about')]),
          a(href: 'mailto:jon@jonreyna.dev', [.text('jon@jonreyna.dev')]),
          a(href: 'https://github.com/jonreyna', [.text('GitHub')]),
          a(href: 'https://www.linkedin.com/in/jonreyna', [.text('LinkedIn')]),
          a(href: '/feed.xml', [.text('feed')]),
        ]),
      ]),
    ]);
  }

  Component _essayHead(Map<String, Object?> data) {
    final date = '${data['date'] ?? ''}';
    final parts = date.split('-');
    final shown = parts.length >= 2 ? '${_months[int.parse(parts[1]) - 1]} ${parts[0]}' : date;
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
