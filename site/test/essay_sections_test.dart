// Tests for lib/essay_sections.dart: section numbering, slug generation, hand-written
// ids, deduplication, and non-empty fallback ids.
import 'package:jaspr_content/jaspr_content.dart';
import 'package:jonreyna_dev/essay_sections.dart';
import 'package:test/test.dart';

Page makePage({bool essay = false}) => Page(
  path: 'test.md',
  url: '/test',
  content: '',
  initialData: {
    'page': {if (essay) 'kind': 'essay'},
  },
  config: PageConfig(),
  loader: MemoryLoader(pages: const []),
);

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

ElementNode h2(String text, {Map<String, String> attrs = const {}}) =>
    ElementNode('h2', attrs, [if (text.isNotEmpty) TextNode(text)]);

void main() {
  const ext = EssaySectionsExtension();
  final cases = <(String, bool, List<Node>, String)>[
    (
      'essays get roman numerals and section anchors',
      true,
      [h2('First'), h2('Second'), h2('Third'), h2('Fourth')],
      '<h2 id="first"><span class="n">I.</span>First<a class="anchor" href="#first" aria-label="Link to section I">§</a></h2>'
          '<h2 id="second"><span class="n">II.</span>Second<a class="anchor" href="#second" aria-label="Link to section II">§</a></h2>'
          '<h2 id="third"><span class="n">III.</span>Third<a class="anchor" href="#third" aria-label="Link to section III">§</a></h2>'
          '<h2 id="fourth"><span class="n">IV.</span>Fourth<a class="anchor" href="#fourth" aria-label="Link to section IV">§</a></h2>',
    ),
    (
      'plain pages get ids and anchors without numerals',
      false,
      [h2('Selected work')],
      '<h2 id="selected-work">Selected work<a class="anchor" href="#selected-work" aria-label="Link to this section">§</a></h2>',
    ),
    (
      'plain pages preserve hand-written ids',
      false,
      [h2('Play', attrs: {'id': 'custom-play'})],
      '<h2 id="custom-play">Play<a class="anchor" href="#custom-play" aria-label="Link to this section">§</a></h2>',
    ),
    (
      'headings with a class are skipped',
      true,
      [h2('Exercises', attrs: {'class': 'exercises'}), h2('The rule')],
      '<h2 class="exercises">Exercises</h2>'
          '<h2 id="the-rule"><span class="n">I.</span>The rule<a class="anchor" href="#the-rule" aria-label="Link to section I">§</a></h2>',
    ),
    (
      'duplicate heading titles get unique ids',
      true,
      [h2('Notes'), h2('Notes'), h2('Notes')],
      '<h2 id="notes"><span class="n">I.</span>Notes<a class="anchor" href="#notes" aria-label="Link to section I">§</a></h2>'
          '<h2 id="notes-2"><span class="n">II.</span>Notes<a class="anchor" href="#notes-2" aria-label="Link to section II">§</a></h2>'
          '<h2 id="notes-3"><span class="n">III.</span>Notes<a class="anchor" href="#notes-3" aria-label="Link to section III">§</a></h2>',
    ),
    (
      'symbol-only headings fall back to a non-empty id',
      true,
      [
        ElementNode('h2', {}, [ElementNode('math', {}, [ElementNode('mi', {}, [const TextNode('x')])])]),
        h2('—'),
      ],
      '<h2 id="x"><span class="n">I.</span><math><mi>x</mi></math><a class="anchor" href="#x" aria-label="Link to section I">§</a></h2>'
          '<h2 id="section-2"><span class="n">II.</span>—<a class="anchor" href="#section-2" aria-label="Link to section II">§</a></h2>',
    ),
  ];

  for (final (name, essay, input, want) in cases) {
    test(name, () async {
      final got = await ext.apply(makePage(essay: essay), input);
      expect(html(got), want);
    });
  }
}
