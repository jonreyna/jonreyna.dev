// Tests for lib/nested_links.dart: unwrapping a link the Markdown parser nested in another.
import 'package:jaspr_content/jaspr_content.dart';
import 'package:jonreyna_dev/nested_links.dart';
import 'package:test/test.dart';

/// Renders nodes as compact HTML, so expectations read like the page source.
String html(List<Node> nodes) => [
  for (final node in nodes)
    switch (node) {
      TextNode(:final text) => text,
      ElementNode(:final tag, :final attributes, :final children) =>
        '<$tag${attributes.entries.map((e) => ' ${e.key}="${e.value}"').join()}>'
            '${html(children ?? const [])}</$tag>',
      _ => '',
    },
].join();

ElementNode a(String href, List<Node> children) => ElementNode('a', {'href': href}, children);

void main() {
  final cases = <(String, List<Node>, String)>[
    (
      'an email link the parser linked twice',
      [
        ElementNode('p', {}, [
          const TextNode('Write to '),
          a('mailto:jon@jonreyna.dev', [
            a('mailto:jon@jonreyna.dev', [const TextNode('jon@jonreyna.dev')]),
          ]),
        ]),
      ],
      '<p>Write to <a href="mailto:jon@jonreyna.dev">jon@jonreyna.dev</a></p>',
    ),
    (
      'a link nested deeper keeps the markup around it',
      [
        a('/outer', [
          ElementNode('em', {}, [
            a('/inner', [const TextNode('text')]),
          ]),
        ]),
      ],
      '<a href="/outer"><em>text</em></a>',
    ),
    (
      'links side by side are left alone',
      [
        ElementNode('p', {}, [
          a('/a', [const TextNode('A')]),
          const TextNode(' and '),
          a('/b', [const TextNode('B')]),
        ]),
      ],
      '<p><a href="/a">A</a> and <a href="/b">B</a></p>',
    ),
    (
      'text without links is unchanged',
      [
        ElementNode('p', {'class': 'x'}, [const TextNode('plain')]),
      ],
      '<p class="x">plain</p>',
    ),
  ];
  for (final (name, input, want) in cases) {
    test(name, () => expect(html(unwrapNestedLinks(input)), want));
  }
}
