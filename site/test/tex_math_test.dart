// Tests for lib/tex_math.dart: finding TeX math in Markdown, DESIGN.md's MathML
// rules, and (when PANDOC is set, as `make site-test` does) the conversion end to end.
import 'dart:io';

import 'package:jonreyna_dev/tex_math.dart';
import 'package:test/test.dart';

void main() {
  group('findMath', () {
    // (name, source, expected [tex, display, ownParagraph, trailing] per formula)
    final cases = <(String, String, List<(String, bool, bool, String)>)>[
      ('inline', r'If $y = f(u)$ then', [('y = f(u)', false, false, '')]),
      ('two inline', r'$a$ and $b_i$ here', [('a', false, false, ''), ('b_i', false, false, '')]),
      ('comma after inline stays with it', r'set $x$, then', [('x', false, false, ',')]),
      ('period after inline stays with it', r'for all $j$.', [('j', false, false, '.')]),
      ('prices are not math', r'It cost $5 and $10.', []),
      ('space after the opening dollar', r'a $ b$ c', []),
      ('digit after the closing dollar', r'$x$1 and', []),
      ('escaped dollar', r'costs \$5, and $x$ is', [('x', false, false, '')]),
      ('code span', r'run `echo $HOME` and $x$', [('x', false, false, '')]),
      ('double-backtick code span', r'``a ` $b$`` then $c$', [('c', false, false, '')]),
      ('fenced code', 'before\n\n```sh\necho \$a \$b\n```\n\nafter \$c\$', [('c', false, false, '')]),
      ('tilde fence', '~~~\n\$a\$\n~~~\n\$b\$', [('b', false, false, '')]),
      ('comment', r'<!-- $a$ --> $b$', [('b', false, false, '')]),
      ('display paragraph', 'Text.\n\n\$\$\n\\sum_i x_i\n\$\$\n\nMore.', [(r'\sum_i x_i', true, true, '')]),
      ('display one line', 'Text.\n\n\$\$x^2\$\$\n\nMore.', [('x^2', true, true, '')]),
      (
        'display inside HTML',
        '<div class="eq num" id="eq-1">\$\$x = 1\$\$<span class="no">(1)</span></div>',
        [('x = 1', true, false, '')],
      ),
      ('inline math stops at a blank line', 'a \$b\n\nc\$ d', []),
    ];
    for (final (name, source, want) in cases) {
      test(name, () {
        final got = [for (final s in findMath(source)) (s.tex, s.display, s.ownParagraph, s.trailing)];
        expect(got, want);
      });
    }

    test('unclosed display math is an error', () {
      expect(() => findMath('a \$\$x\n\nb'), throwsFormatException);
    });
  });

  group('polishMathML', () {
    final cases = <(String, String, String)>[
      (
        'drops the annotation, wrapper, namespace and inline display',
        '<math display="inline" xmlns="http://www.w3.org/1998/Math/MathML"><semantics><mi>x</mi>'
            '<annotation encoding="application/x-tex">x</annotation></semantics></math>',
        '<math><mi>x</mi></math>',
      ),
      (
        'keeps block display',
        '<math display="block" xmlns="http://www.w3.org/1998/Math/MathML"><mi>x</mi></math>',
        '<math display="block"><mi>x</mi></math>',
      ),
      (
        'parentheses get zero lspace and rspace',
        '<math><mi>f</mi><mo stretchy="false" form="prefix">(</mo><mi>u</mi><mo stretchy="false" form="postfix">)</mo></math>',
        '<math><mi>f</mi><mo form="prefix" stretchy="false" lspace="0" rspace="0">(</mo><mi>u</mi>'
            '<mo form="postfix" stretchy="false" lspace="0" rspace="0">)</mo></math>',
      ),
      (
        'an overbar is a flat U+00AF under accent="false"',
        '<math><msub><mover><mi>v</mi><mo accent="true">‾</mo></mover><mi>j</mi></msub></math>',
        '<math><msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>j</mi></msub></math>',
      ),
      (
        'an overbar over a group, inside another mover',
        '<mover><mover><mrow><mi>a</mi><mi>b</mi></mrow><mo accent="true">‾</mo></mover><mo>→</mo></mover>',
        '<mover><mover accent="false"><mrow><mi>a</mi><mi>b</mi></mrow><mo>&#xAF;</mo></mover><mo>→</mo></mover>',
      ),
      (
        'a division slash has zero lspace and rspace',
        '<mi>∂</mi><mi>/</mi><mi>∂</mi>',
        '<mi>∂</mi><mo lspace="0" rspace="0">/</mo><mi>∂</mi>',
      ),
      ('an ellipsis is an operator', '<mo>,</mo><mi>…</mi><mo>,</mo>', '<mo>,</mo><mo>…</mo><mo>,</mo>'),
      (
        'an upright name is one identifier',
        '<mi>i</mi><mo>∈</mo><mrow><mi mathvariant="normal">o</mi><mi mathvariant="normal">u</mi><mi mathvariant="normal">t</mi></mrow>',
        '<mi>i</mi><mo>∈</mo><mrow><mi>out</mi></mrow>',
      ),
      (
        'a single upright letter stays as it is',
        '<mi mathvariant="normal">d</mi><mi>x</mi>',
        '<mi mathvariant="normal">d</mi><mi>x</mi>',
      ),
    ];
    for (final (name, input, want) in cases) {
      test(name, () => expect(polishMathML(input), want));
    }
  });

  group('replaceMath', () {
    test('wraps a display paragraph and keeps trailing punctuation with inline math', () {
      const source = 'Set \$x\$, then\n\n\$\$y\$\$\n\nend';
      final spans = findMath(source);
      final got = replaceMath(source, spans, ['<math>X</math>', '<math display="block">Y</math>']);
      expect(got, 'Set <span class="nowrap"><math>X</math>,</span> then\n\n'
          '<div class="eq"><math display="block">Y</math></div>\n\nend');
    });
  });

  // End to end with the pinned Pandoc: the specimen's formulas, written in TeX, give
  // the MathML DESIGN.md asks for.
  group('Pandoc', () {
    final pandoc = Platform.environment['PANDOC'] ?? '';
    final skip = pandoc.isEmpty ? 'PANDOC is not set (make site-test sets it)' : null;
    final cases = <(String, String)>[
      (
        r'y = f(u)',
        '<math><mrow><mi>y</mi><mo>=</mo><mi>f</mi><mo form="prefix" stretchy="false" lspace="0" rspace="0">(</mo>'
            '<mi>u</mi><mo form="postfix" stretchy="false" lspace="0" rspace="0">)</mo></mrow></math>',
      ),
      (
        r'\bar{v}_j \gets \bar{v}_j + \bar{v}_i \, \partial v_i / \partial v_j',
        '<math><mrow><msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>j</mi></msub><mo>←</mo>'
            '<msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>j</mi></msub><mo>+</mo>'
            '<msub><mover accent="false"><mi>v</mi><mo>&#xAF;</mo></mover><mi>i</mi></msub>'
            '<mspace width="0.167em"></mspace><mi>∂</mi><msub><mi>v</mi><mi>i</mi></msub>'
            '<mo lspace="0" rspace="0">/</mo><mi>∂</mi><msub><mi>v</mi><mi>j</mi></msub></mrow></math>',
      ),
      (
        r'v_1, \ldots, v_n',
        '<math><mrow><msub><mi>v</mi><mn>1</mn></msub><mo>,</mo><mo>…</mo><mo>,</mo>'
            '<msub><mi>v</mi><mi>n</mi></msub></mrow></math>',
      ),
    ];
    for (final (tex, want) in cases) {
      test(tex, skip: skip, () async {
        final spans = findMath('\$$tex\$');
        expect(await convertWithPandoc(spans, pandoc, 'test'), [want]);
      });
    }

    test('TeX that Pandoc cannot parse fails', skip: skip, () async {
      final spans = findMath(r'$\frac{a$');
      await expectLater(convertWithPandoc(spans, pandoc, 'test'), throwsStateError);
    });
  });
}
