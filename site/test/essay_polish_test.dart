// Tests for lib/essay_polish.dart: footnote fleuron, backref tie and label, and
// table scroll wrappers.
import 'package:jaspr_content/jaspr_content.dart';
import 'package:jonreyna_dev/essay_polish.dart';
import 'package:test/test.dart';

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

void main() {
  group('polishNodes', () {
    test('adds fleuron, ties backref with nbsp, and adds aria-label to footnotes', () {
      final input = <Node>[
        ElementNode('p', {}, [const TextNode('Body.')]),
        ElementNode('section', {'class': 'footnotes'}, [
          ElementNode('ol', {}, [
            ElementNode('li', {'id': 'fn-1'}, [
              ElementNode('p', {}, [
                const TextNode('Note text. '),
                ElementNode('a', {'href': '#fnref-1', 'class': 'footnote-backref'}, [const TextNode('↩')]),
              ]),
            ]),
          ]),
        ]),
      ];

      expect(
        html(polishNodes(input)),
        '<p>Body.</p>'
        '<p class="fleuron" aria-hidden="true">⁂</p>'
        '<section class="footnotes"><ol><li id="fn-1"><p>Note text.\u00A0'
        '<a href="#fnref-1" class="footnote-backref" aria-label="Back to the text">↩</a>'
        '</p></li></ol></section>',
      );
    });

    test('does not duplicate an existing fleuron or overwrite custom aria-label', () {
      final input = <Node>[
        ElementNode('p', {'class': 'fleuron', 'aria-hidden': 'true'}, [const TextNode('⁂')]),
        ElementNode('section', {'class': 'footnotes'}, [
          ElementNode('ol', {}, [
            ElementNode('li', {'id': 'fn-1'}, [
              ElementNode('p', {}, [
                ElementNode('em', {}, [const TextNode('Work')]),
                const TextNode(' '),
                ElementNode('a', {'href': '#fnref-1', 'class': 'footnote-backref', 'aria-label': 'Return'}, [
                  const TextNode('↩'),
                ]),
              ]),
            ]),
          ]),
        ]),
      ];

      expect(
        html(polishNodes(input)),
        '<p class="fleuron" aria-hidden="true">⁂</p>'
        '<section class="footnotes"><ol><li id="fn-1"><p><em>Work</em>\u00A0'
        '<a href="#fnref-1" class="footnote-backref" aria-label="Return">↩</a>'
        '</p></li></ol></section>',
      );
    });

    test('wraps standalone tables in .table-scroll without double-wrapping', () {
      final input = <Node>[
        ElementNode('table', {}, [
          ElementNode('tr', {}, [ElementNode('td', {}, [const TextNode('a')])]),
        ]),
        ElementNode('div', {'class': 'table-scroll'}, [
          ElementNode('table', {}, [
            ElementNode('tr', {}, [ElementNode('td', {}, [const TextNode('b')])]),
          ]),
        ]),
      ];

      expect(
        html(polishNodes(input)),
        '<div class="table-scroll"><table><tr><td>a</td></tr></table></div>'
        '<div class="table-scroll"><table><tr><td>b</td></tr></table></div>',
      );
    });
  });
}
