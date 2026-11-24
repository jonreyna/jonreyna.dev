/// TeX math in Markdown, converted to MathML while the site is built (DESIGN.md,
/// principle 4: reading pages ship no JavaScript).
///
/// Write math as in Pandoc's Markdown:
///
/// - `$…$` inline. The opening `$` must be followed by a non-space and the closing
///   one preceded by a non-space and not followed by a digit, so "$5 and $10" stays
///   text. `\$` is a literal dollar sign.
/// - `$$…$$` display. On a paragraph of its own it becomes `<div class="eq">`; inside
///   HTML, such as a numbered equation's `<div class="eq num">`, it stays as written.
///
/// Code fences, code spans and HTML comments are left alone.
///
/// The conversion is Pandoc's (texmath), one call per page, with the binary that
/// `make` downloads and pins; `make site` passes its path in the PANDOC environment
/// variable. Any Pandoc warning, such as TeX it can't parse, fails the build. The
/// MathML is then adjusted to DESIGN.md's typesetting rules (see [polishMathML]),
/// and a period or comma right after inline math is kept on its line.
library;

import 'package:jaspr_content/jaspr_content.dart';

import 'pandoc.dart';

/// Converts each page's TeX math to MathML before the Markdown is parsed, so the
/// Markdown parser sees only raw MathML (which it passes through) and never reads a
/// TeX subscript as emphasis.
class TexMath implements TemplateEngine {
  const TexMath();

  @override
  Future<void> render(Page page, List<Page> pages) async {
    final spans = findMath(page.content);
    if (spans.isEmpty) return;
    final pandoc = pandocPath('${page.path} has TeX math, which');
    final mathml = await convertWithPandoc(spans, pandoc, page.path);
    page.apply(content: replaceMath(page.content, spans, mathml));
  }
}

/// One formula found in a page's Markdown.
class MathSpan {
  const MathSpan({
    required this.start,
    required this.end,
    required this.tex,
    required this.display,
    this.ownParagraph = false,
    this.trailing = '',
  });

  /// Offsets of the formula in the source, delimiters and [trailing] included.
  final int start, end;

  /// The TeX between the delimiters, trimmed.
  final String tex;

  /// Whether it was written `$$…$$`.
  final bool display;

  /// Whether a display formula stands alone as a paragraph.
  final bool ownParagraph;

  /// Punctuation right after inline math, kept on the same line as it.
  final String trailing;
}

/// Finds the TeX math in [source], in order.
List<MathSpan> findMath(String source) {
  final spans = <MathSpan>[];
  var i = 0;
  var lineStart = true;
  while (i < source.length) {
    final c = source[i];
    // Code fences: ``` or ~~~ at the start of a line, closed (as in CommonMark) by a
    // run of the same character at least as long.
    if (lineStart) {
      final fence = RegExp(r' {0,3}(`{3,}|~{3,})').matchAsPrefix(source, i);
      if (fence != null) {
        final marker = fence.group(1)!;
        final close = RegExp(
          '^ {0,3}${RegExp.escape(marker[0])}{${marker.length},}[ \t]*\$',
          multiLine: true,
        );
        final afterOpen = source.indexOf('\n', i);
        final m = afterOpen < 0 ? null : close.firstMatch(source.substring(afterOpen + 1));
        i = m == null ? source.length : afterOpen + 1 + m.end;
        continue;
      }
    }
    lineStart = c == '\n';
    if (c == '\\') {
      i += 2; // an escaped character, such as \$, is never a delimiter
      continue;
    }
    if (c == '`') {
      // A code span ends at the next run of the same number of backticks.
      final run = RegExp('`+').matchAsPrefix(source, i)!.group(0)!;
      final close = source.indexOf(RegExp('(?<!`)${RegExp.escape(run)}(?!`)'), i + run.length);
      i = close < 0 ? i + run.length : close + run.length;
      continue;
    }
    if (source.startsWith('<!--', i)) {
      final close = source.indexOf('-->', i + 4);
      i = close < 0 ? source.length : close + 3;
      continue;
    }
    if (source.startsWith(r'$$', i)) {
      final close = _closingDoubleDollar(source, i + 2);
      if (close == null || source.substring(i + 2, close).contains(RegExp(r'\n[ \t]*\n'))) {
        throw FormatException(r'unclosed $$ (display math may not contain a blank line)', source, i);
      }
      final end = close + 2;
      final before = i == 0 ? -1 : source.lastIndexOf('\n', i - 1);
      final after = source.indexOf('\n', end);
      final ownParagraph = source.substring(before + 1, i).trim().isEmpty &&
          source.substring(end, after < 0 ? source.length : after).trim().isEmpty &&
          _blankOrEdge(source, before, back: true) &&
          _blankOrEdge(source, after, back: false);
      spans.add(MathSpan(
        start: i,
        end: end,
        tex: source.substring(i + 2, close).trim(),
        display: true,
        ownParagraph: ownParagraph,
      ));
      i = end;
      continue;
    }
    if (c == r'$' && i + 1 < source.length && !_isSpace(source[i + 1])) {
      final close = _closingDollar(source, i + 1);
      if (close != null) {
        var end = close + 1;
        var trailing = '';
        if (end < source.length && '.,;:!?)'.contains(source[end])) {
          trailing = source[end];
          end++;
        }
        spans.add(MathSpan(
          start: i,
          end: end,
          tex: source.substring(i + 1, close),
          display: false,
          trailing: trailing,
        ));
        i = end;
        continue;
      }
    }
    i++;
  }
  return spans;
}

/// The closing `$` of inline math opened just before [from], or null if it has none.
///
/// As in Pandoc, the formula ends at the first unescaped `$` outside braces: a `$`
/// inside a group, as in `$\text{if $x$ then}$`, belongs to the formula. That `$`
/// closes it only with a non-space before it and no digit after it; otherwise the
/// opening `$` is just a dollar sign, so "It cost $5 and $10" stays text. A blank
/// line also ends the search.
int? _closingDollar(String source, int from) {
  var depth = 0;
  for (var j = from; j < source.length; j++) {
    final c = source[j];
    if (c == '\\') {
      j++;
      continue;
    }
    if (c == '\n' && RegExp(r'\n[ \t]*\n').matchAsPrefix(source, j) != null) return null;
    if (c == '{') depth++;
    if (c == '}' && depth > 0) depth--;
    if (c == r'$' && depth == 0) {
      final digitAfter = j + 1 < source.length && RegExp(r'\d').hasMatch(source[j + 1]);
      return j > from && !_isSpace(source[j - 1]) && !digitAfter ? j : null;
    }
  }
  return null;
}

/// The closing `$$` of display math whose TeX starts at [from], or null. Escaped
/// characters are skipped, so `$$\$$$` holds a dollar sign.
int? _closingDoubleDollar(String source, int from) {
  for (var j = from; j < source.length - 1; j++) {
    if (source[j] == '\\') {
      j++;
      continue;
    }
    if (source.startsWith(r'$$', j)) return j;
  }
  return null;
}

bool _isSpace(String c) => c == ' ' || c == '\t' || c == '\n';

/// Whether the line next to the one ending (or starting) at [newline] is blank, or
/// there is none.
bool _blankOrEdge(String source, int newline, {required bool back}) {
  if (newline < 0) return true;
  if (back) {
    if (newline == 0) return true;
    final prev = source.lastIndexOf('\n', newline - 1);
    return source.substring(prev + 1, newline).trim().isEmpty;
  }
  final next = source.indexOf('\n', newline + 1);
  return source.substring(newline + 1, next < 0 ? source.length : next).trim().isEmpty;
}

/// Converts [spans] to MathML with one Pandoc call. Each formula is its own
/// paragraph in Pandoc's input, so the output holds one `<math>` per formula, in
/// order.
Future<List<String>> convertWithPandoc(List<MathSpan> spans, String pandoc, String path) async {
  final input = spans.map((s) => s.display ? '\$\$${s.tex}\$\$' : '\$${s.tex.trim()}\$').join('\n\n');
  final stdout = await runPandoc(
    pandoc,
    ['--from=markdown', '--to=html', '--math-method=mathml', '--wrap=none'],
    input,
    '$path: converting TeX math',
  );
  final math = RegExp(r'<math\b.*?</math>', dotAll: true).allMatches(stdout).map((m) => m[0]!).toList();
  if (math.length != spans.length) {
    // Pandoc leaves TeX it can't read as text, still between its dollar signs.
    final unread = RegExp(r'<p>(\$.*?\$)</p>', dotAll: true).allMatches(stdout).map((m) => m[1]).join(', ');
    throw StateError('$path: Pandoc read ${math.length} of ${spans.length} formulas as math; check the TeX: $unread');
  }
  return [for (final m in math) polishMathML(m)];
}

/// Adjusts Pandoc's MathML to DESIGN.md's typesetting rules:
///
/// - the TeX annotation, `semantics` wrapper, namespace and `display="inline"` go
///   (the page is HTML, and the source is in the repository);
/// - fences that don't stretch, such as parentheses, get zero `lspace` and `rspace`,
///   or Chrome spaces them like operators;
/// - an overbar is `<mover accent="false">` over U+00AF, or Chrome draws it as a
///   sub-pixel line;
/// - an upright name such as `\mathrm{out}` is one `<mi>`, not one per letter;
/// - a slash used as division is an operator with zero `lspace` and `rspace`;
/// - an ellipsis is an operator, spaced like the commas around it.
String polishMathML(String mathml) {
  var m = mathml
      .replaceAll(RegExp(r'<annotation\b.*?</annotation>', dotAll: true), '')
      .replaceAll(RegExp(r'</?semantics>'), '')
      .replaceAll(' xmlns="http://www.w3.org/1998/Math/MathML"', '')
      .replaceAll(' display="inline"', '');
  m = m.replaceAllMapped(
    RegExp(r'<mo stretchy="false" form="(prefix|postfix)">'),
    (x) => '<mo form="${x[1]}" stretchy="false" lspace="0" rspace="0">',
  );
  // \mathrm{out} arrives as one upright <mi> per letter; a single multi-letter <mi> is
  // upright already and spaced as one name.
  m = m.replaceAllMapped(
    RegExp(r'(?:<mi mathvariant="normal">[^<]</mi>){2,}'),
    (x) => '<mi>${x[0]!.replaceAll(RegExp(r'</?mi[^>]*>'), '')}</mi>',
  );
  m = m.replaceAll('<mi>/</mi>', '<mo lspace="0" rspace="0">/</mo>');
  m = m.replaceAll('<mi>…</mi>', '<mo>…</mo>');
  return _flatOverbars(m);
}

/// Rewrites each `<mover>` whose last child is an accent bar (U+203E or U+00AF) as
/// `<mover accent="false">` over U+00AF.
String _flatOverbars(String m) {
  const barSuffix = '<mo accent="true">‾</mo>';
  final bar = RegExp(r'<mo accent="true">[‾¯]</mo></mover>');
  final out = StringBuffer();
  final opens = <({int offset, String openTag})>[]; // offsets in [out] of each open <mover>'s tag
  final tag = RegExp(r'<mover\b[^>]*>|</mover>');
  var pos = 0;
  for (final t in tag.allMatches(m)) {
    out.write(m.substring(pos, t.start));
    if (t[0] != '</mover>') {
      opens.add((offset: out.length, openTag: t[0]!));
      out.write(t[0]!);
      pos = t.end;
      continue;
    }
    if (opens.isEmpty) {
      out.write('</mover>');
      pos = t.end;
      continue;
    }
    final open = opens.removeLast();
    final barStart = t.start - barSuffix.length;
    final barMatch = barStart >= 0 && open.openTag == '<mover>' ? bar.matchAsPrefix(m, barStart) : null;
    if (barMatch != null) {
      // Drop the bar written so far, write the flat one, and mark the opening tag.
      final text = out.toString();
      final cut = text.length - barSuffix.length;
      out
        ..clear()
        ..write(text.substring(0, open.offset))
        ..write('<mover accent="false">')
        ..write(text.substring(open.offset + '<mover>'.length, cut))
        ..write('<mo>&#xAF;</mo></mover>');
    } else {
      out.write('</mover>');
    }
    pos = t.end;
  }
  out.write(m.substring(pos));
  return out.toString();
}

/// Replaces each span in [source] with its MathML.
String replaceMath(String source, List<MathSpan> spans, List<String> mathml) {
  final out = StringBuffer();
  var pos = 0;
  for (final (k, s) in spans.indexed) {
    out.write(source.substring(pos, s.start));
    final math = mathml[k];
    if (s.display && s.ownParagraph) {
      out.write('<div class="eq">$math</div>');
    } else if (s.trailing.isNotEmpty) {
      out.write('<span class="nowrap">$math${s.trailing}</span>');
    } else {
      out.write(math);
    }
    pos = s.end;
  }
  out.write(source.substring(pos));
  return out.toString();
}
