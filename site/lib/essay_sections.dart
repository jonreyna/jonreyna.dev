import 'package:jaspr_content/jaspr_content.dart';

/// Numbers an essay's sections the way DESIGN.md asks ("Section head"), so authors
/// write a plain `## The rule` and get
///
/// ```html
/// <h2 id="the-rule"><span class="n">I.</span>The rule<a class="anchor" href="#the-rule"
///     aria-label="Link to section I">§</a></h2>
/// ```
///
/// Only top-level `h2`s in pages with `kind: essay` are numbered; headings with a class
/// (Exercises, Answers) are left alone. Other pages get the id and the § link but no numeral.
class EssaySectionsExtension implements PageExtension {
  const EssaySectionsExtension();

  @override
  Future<List<Node>> apply(Page page, List<Node> nodes) async {
    final essay = page.data.page['kind'] == 'essay';
    var count = 0;
    var headingIndex = 0;
    final seen = <String>{};
    return [
      for (final node in nodes)
        if (node is ElementNode && node.tag == 'h2' && node.attributes['class'] == null)
          _section(node, essay ? _roman(++count) : null, ++headingIndex, seen)
        else
          node,
    ];
  }

  ElementNode _section(ElementNode h2, String? numeral, int headingIndex, Set<String> seen) {
    final children = h2.children ?? const <Node>[];
    // An id written by hand on a plain page wins (`<h2 id="play">`); in essays the ids are
    // always made from the text, since the Markdown ones include the numeral's markup.
    final String id;
    if (h2.attributes['id'] case final existing? when numeral == null && existing.isNotEmpty) {
      id = existing;
      seen.add(id);
    } else {
      final slug = _slug(_text(children));
      final base = slug.isEmpty ? 'section-$headingIndex' : slug;
      var candidate = base;
      var suffix = 2;
      while (seen.contains(candidate)) {
        candidate = '$base-${suffix++}';
      }
      id = candidate;
      seen.add(id);
    }
    final label = numeral == null ? 'Link to this section' : 'Link to section $numeral';
    return ElementNode('h2', {...h2.attributes, 'id': id}, [
      if (numeral != null) ElementNode('span', {'class': 'n'}, [TextNode('$numeral.')]),
      ...children,
      ElementNode('a', {'class': 'anchor', 'href': '#$id', 'aria-label': label}, [const TextNode('§')]),
    ]);
  }

  static String _text(List<Node> nodes) => [
    for (final node in nodes)
      switch (node) {
        TextNode(:final text) => text,
        ElementNode(:final children) => _text(children ?? const []),
        _ => '',
      },
  ].join();

  /// "What the tutorials skip" → "what-the-tutorials-skip"; keeps letters such as ü.
  static String _slug(String text) => text
      .toLowerCase()
      .replaceAll(RegExp(r'&[a-z]+;|<[^>]*>'), ' ')
      .replaceAll(RegExp(r'[^\p{L}\p{N}]+', unicode: true), '-')
      .replaceAll(RegExp(r'^-+|-+$'), '');

  static String _roman(int n) {
    const values = [(10, 'X'), (9, 'IX'), (5, 'V'), (4, 'IV'), (1, 'I')];
    final out = StringBuffer();
    for (final (value, numeral) in values) {
      while (n >= value) {
        out.write(numeral);
        n -= value;
      }
    }
    return out.toString();
  }
}
