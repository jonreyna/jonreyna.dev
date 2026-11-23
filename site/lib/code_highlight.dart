/// Syntax highlighting for code blocks, done while the site is built (DESIGN.md,
/// principle 4: reading pages ship no JavaScript).
///
/// A fenced block with a language (```` ```go ````) is highlighted by Pandoc
/// (skylighting, about 140 languages), one call per page, with the binary `make`
/// pins. Pandoc marks each token with a class, such as `kw` for a keyword or `co` for a
/// comment, and site.css styles the few that DESIGN.md asks for: keywords bold,
/// comments italic in `--muted`. Everything else stays in the text color.
///
/// A block without a language, or in one Pandoc doesn't know, is left as written.
library;

import 'package:jaspr_content/jaspr_content.dart';

import 'pandoc.dart';

/// Highlights every code block with a language, in place of its plain text.
class CodeHighlightExtension implements PageExtension {
  const CodeHighlightExtension();

  @override
  Future<List<Node>> apply(Page page, List<Node> nodes) async {
    final blocks = <CodeBlock>[];
    collectCodeBlocks(nodes, blocks);
    if (blocks.isEmpty) return nodes;
    final html = await highlightWithPandoc(blocks, pandocPath('${page.path} has code to highlight, which'), page.path);
    final replacements = <ElementNode, ElementNode>{
      for (final (k, b) in blocks.indexed)
        if (html[k] case final h?)
          b.pre: ElementNode('pre', b.pre.attributes, [
            ElementNode('code', b.code.attributes, [TextNode(h, raw: true)]),
          ]),
    };
    return _replace(nodes, replacements);
  }

  static List<Node> _replace(List<Node> nodes, Map<ElementNode, ElementNode> replacements) => [
    for (final node in nodes)
      if (node is ElementNode && replacements.containsKey(node))
        replacements[node]!
      else if (node is ElementNode && node.children != null)
        ElementNode(node.tag, node.attributes, _replace(node.children!, replacements))
      else
        node,
  ];
}

/// A code block found in a parsed page: `<pre><code class="language-…">`.
class CodeBlock {
  const CodeBlock(this.pre, this.code, this.language, this.text);

  /// The `pre` element, which the highlighted block replaces.
  final ElementNode pre;

  /// Its `code` element, whose attributes the highlighted block keeps.
  final ElementNode code;

  /// The language from the fence, such as `go`.
  final String language;

  /// The code, as plain text.
  final String text;
}

/// Adds every code block with a language in [nodes] to [out], in document order.
void collectCodeBlocks(List<Node> nodes, List<CodeBlock> out) {
  for (final node in nodes) {
    if (node is! ElementNode) continue;
    final children = node.children ?? const [];
    if (children case [final ElementNode code] when node.tag == 'pre' && code.tag == 'code') {
      final language = RegExp(r'(?:^|\s)language-(\S+)').firstMatch(code.attributes['class'] ?? '')?.group(1);
      if (language != null) out.add(CodeBlock(node, code, language, code.innerText));
      continue;
    }
    collectCodeBlocks(children, out);
  }
}

/// Highlights [blocks] with one Pandoc call and returns, for each, the highlighted
/// HTML to put inside its `code` element, or null when Pandoc doesn't know the
/// language.
Future<List<String?>> highlightWithPandoc(List<CodeBlock> blocks, String pandoc, String path) async {
  final input = [
    for (final b in blocks) _fenced(b),
  ].join('\n\n');
  final stdout = await runPandoc(
    pandoc,
    ['--from=markdown', '--to=html', '--preserve-tabs', '--wrap=none', '--syntax-highlighting=default'],
    input,
    '$path: highlighting code',
  );
  final pres = RegExp(r'<pre\b[^>]*>.*?</pre>', dotAll: true).allMatches(stdout).map((m) => m[0]!).toList();
  if (pres.length != blocks.length) {
    throw StateError('$path: Pandoc returned ${pres.length} code blocks for ${blocks.length}');
  }
  return [for (final pre in pres) pre.contains('class="sourceCode') ? cleanHighlighted(pre) : null];
}

/// A block as a Pandoc code fence, longer than any run of backticks in the code.
String _fenced(CodeBlock b) {
  final longest = RegExp('`+').allMatches(b.text).fold(0, (n, m) => m[0]!.length > n ? m[0]!.length : n);
  final fence = '`' * (longest < 3 ? 3 : longest + 1);
  final text = b.text.endsWith('\n') ? b.text : '${b.text}\n';
  return '$fence${b.language}\n$text$fence';
}

/// The inside of Pandoc's `<code>` element, without the numbered span and empty
/// anchor link Pandoc wraps around each line (the site doesn't link to lines).
String cleanHighlighted(String pre) {
  final inner = RegExp(r'<code\b[^>]*>(.*)</code>', dotAll: true).firstMatch(pre)!.group(1)!;
  final line = RegExp(r'^<span id="cb\d+-\d+"><a href="#cb\d+-\d+" aria-hidden="true" tabindex="-1"></a>(.*)</span>$');
  return [
    for (final l in inner.split('\n')) line.firstMatch(l)?.group(1) ?? l,
  ].join('\n');
}
