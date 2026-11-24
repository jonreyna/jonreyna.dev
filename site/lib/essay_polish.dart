import 'package:jaspr_content/jaspr_content.dart';

/// Polishes essay footnotes and tables to match DESIGN.md:
///
/// - inserts the centered asterism (`<p class="fleuron" aria-hidden="true">⁂</p>`)
///   before `<section class="footnotes">` when not already present;
/// - ties each footnote back-link (`↩`) to its preceding word with a non-breaking space
///   and gives it `aria-label="Back to the text"`;
/// - wraps standalone `<table>` elements in `<div class="table-scroll">` so wide tables
///   scroll inside their own box rather than widening the page.
class EssayPolishExtension implements PageExtension {
  const EssayPolishExtension();

  @override
  Future<List<Node>> apply(Page page, List<Node> nodes) async => polishNodes(nodes);
}

/// Returns [nodes] with DESIGN.md's footnote and table adjustments applied.
List<Node> polishNodes(List<Node> nodes, {bool inFootnotes = false, bool inTableScroll = false}) {
  final out = <Node>[];
  for (final node in nodes) {
    if (node is! ElementNode) {
      out.add(node);
      continue;
    }
    if (node.tag == 'section' && _hasClass(node, 'footnotes')) {
      if (!_precededByFleuron(out)) {
        out.add(ElementNode('p', const {'class': 'fleuron', 'aria-hidden': 'true'}, const [TextNode('⁂')]));
      }
      out.add(
        ElementNode(
          node.tag,
          node.attributes,
          node.children == null ? null : polishNodes(node.children!, inFootnotes: true, inTableScroll: inTableScroll),
        ),
      );
      continue;
    }
    if (inFootnotes && node.tag == 'a' && _hasClass(node, 'footnote-backref')) {
      if (out.isNotEmpty && out.last is TextNode) {
        final prev = out.last as TextNode;
        if (prev.text.endsWith(' ')) {
          out[out.length - 1] = TextNode('${prev.text.substring(0, prev.text.length - 1)}\u00A0', raw: prev.raw);
        }
      }
      final attrs = node.attributes.containsKey('aria-label')
          ? node.attributes
          : {...node.attributes, 'aria-label': 'Back to the text'};
      out.add(
        ElementNode(
          node.tag,
          attrs,
          node.children == null
              ? null
              : polishNodes(node.children!, inFootnotes: inFootnotes, inTableScroll: inTableScroll),
        ),
      );
      continue;
    }
    final nextInTableScroll = inTableScroll || _hasClass(node, 'table-scroll');
    final polished = ElementNode(
      node.tag,
      node.attributes,
      node.children == null
          ? null
          : polishNodes(node.children!, inFootnotes: inFootnotes, inTableScroll: nextInTableScroll),
    );
    if (node.tag == 'table' && !inTableScroll) {
      out.add(ElementNode('div', const {'class': 'table-scroll'}, [polished]));
    } else {
      out.add(polished);
    }
  }
  return out;
}

bool _hasClass(ElementNode node, String token) =>
    (node.attributes['class'] ?? '').split(RegExp(r'\s+')).contains(token);

bool _precededByFleuron(List<Node> nodes) {
  for (var i = nodes.length - 1; i >= 0; i--) {
    final node = nodes[i];
    if (node is TextNode && node.text.trim().isEmpty) continue;
    return node is ElementNode && _hasClass(node, 'fleuron');
  }
  return false;
}
