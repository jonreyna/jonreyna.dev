// Tests for lib/code_highlight.dart: finding code blocks in parsed pages, cleaning
// Pandoc's output, and (when PANDOC is set, as `make site-test` does) highlighting
// end to end.
import 'dart:io';

import 'package:jaspr_content/jaspr_content.dart';
import 'package:jonreyna_dev/code_highlight.dart';
import 'package:test/test.dart';

ElementNode block(String? language, String text) => ElementNode('pre', {}, [
  ElementNode('code', {if (language != null) 'class': 'language-$language'}, [TextNode(text)]),
]);

void main() {
  group('collectCodeBlocks', () {
    final cases = <(String, List<Node>, List<(String, String)>)>[
      ('a block with a language', [block('go', 'x := 1\n')], [('go', 'x := 1\n')]),
      ('a block without a language is left alone', [block(null, 'plain\n')], []),
      (
        'blocks inside other elements, in document order',
        [
          ElementNode('div', {}, [block('go', 'a\n'), ElementNode('section', {}, [block('sh', 'b\n')])]),
          block('python', 'c\n'),
        ],
        [('go', 'a\n'), ('sh', 'b\n'), ('python', 'c\n')],
      ),
      (
        'inline code is not a block',
        [ElementNode('p', {}, [ElementNode('code', {'class': 'language-go'}, [const TextNode('x')])])],
        [],
      ),
    ];
    for (final (name, nodes, want) in cases) {
      test(name, () {
        final got = <CodeBlock>[];
        collectCodeBlocks(nodes, got);
        expect([for (final b in got) (b.language, b.text)], want);
      });
    }
  });

  group('cleanHighlighted', () {
    final cases = <(String, String, String)>[
      (
        'drops the numbered span and empty anchor around each line',
        '<pre class="sourceCode go"><code class="sourceCode go">'
            '<span id="cb1-1"><a href="#cb1-1" aria-hidden="true" tabindex="-1"></a><span class="kw">type</span> T</span>\n'
            '<span id="cb1-2"><a href="#cb1-2" aria-hidden="true" tabindex="-1"></a><span class="co">// c</span></span>'
            '</code></pre>',
        '<span class="kw">type</span> T\n<span class="co">// c</span>',
      ),
      (
        'keeps an empty line',
        '<pre><code><span id="cb2-1"><a href="#cb2-1" aria-hidden="true" tabindex="-1"></a>a</span>\n'
            '<span id="cb2-2"><a href="#cb2-2" aria-hidden="true" tabindex="-1"></a></span>\n'
            '<span id="cb2-3"><a href="#cb2-3" aria-hidden="true" tabindex="-1"></a>b</span></code></pre>',
        'a\n\nb',
      ),
    ];
    for (final (name, input, want) in cases) {
      test(name, () => expect(cleanHighlighted(input), want));
    }
  });

  group('Pandoc', () {
    final pandoc = Platform.environment['PANDOC'] ?? '';
    final skip = pandoc.isEmpty ? 'PANDOC is not set (make site-test sets it)' : null;

    test('highlights Go: keywords, types, comments, escaping, and keeps tabs', skip: skip, () async {
      final blocks = <CodeBlock>[];
      collectCodeBlocks([block('go', 'type V struct {\n\tA float64 // a < b && "c"\n}\n')], blocks);
      expect(await highlightWithPandoc(blocks, pandoc, 'test'), [
        '<span class="kw">type</span> V <span class="kw">struct</span> <span class="op">{</span>\n'
            '\tA <span class="dt">float64</span> <span class="co">// a &lt; b &amp;&amp; &quot;c&quot;</span>\n'
            '<span class="op">}</span>',
      ]);
    });

    test('a language Pandoc does not know is left as written', skip: skip, () async {
      final blocks = <CodeBlock>[];
      collectCodeBlocks([block('nosuchlanguage', 'abc\n'), block('sh', 'echo hi\n')], blocks);
      final got = await highlightWithPandoc(blocks, pandoc, 'test');
      expect(got[0], isNull);
      expect(got[1], contains('echo'));
    });

    test('code containing a code fence', skip: skip, () async {
      final blocks = <CodeBlock>[];
      collectCodeBlocks([block('markdown', 'Use:\n\n```go\nx\n```\n')], blocks);
      final [got] = await highlightWithPandoc(blocks, pandoc, 'test');
      expect(got, allOf(contains('```go'), contains('x')));
    });
  });
}
