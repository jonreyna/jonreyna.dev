import 'package:jaspr_content/jaspr_content.dart';

/// Unwraps a link nested inside another link, keeping its contents.
///
/// The Markdown parser's GitHub extensions turn an email address into a link even
/// inside a link's text, so `[jon@jonreyna.dev](mailto:jon@jonreyna.dev)` comes out as
/// `<a href="mailto:…"><a href="mailto:…">jon@jonreyna.dev</a></a>`. A link inside a
/// link is invalid HTML: browsers split it, leaving the outer link empty, which screen
/// readers announce with no name. The outer link wins, as the author wrote it.
class NestedLinksExtension implements PageExtension {
  const NestedLinksExtension();

  @override
  Future<List<Node>> apply(Page page, List<Node> nodes) async => unwrapNestedLinks(nodes);
}

/// Returns [nodes] with every `a` inside another `a` replaced by its children.
List<Node> unwrapNestedLinks(List<Node> nodes, {bool insideLink = false}) => [
  for (final node in nodes)
    if (node is ElementNode && node.tag == 'a' && insideLink)
      ...unwrapNestedLinks(node.children ?? const [], insideLink: true)
    else if (node is ElementNode && node.children != null)
      ElementNode(
        node.tag,
        node.attributes,
        unwrapNestedLinks(node.children!, insideLink: insideLink || node.tag == 'a'),
      )
    else
      node,
];
