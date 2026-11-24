import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';
import 'package:jaspr_content/jaspr_content.dart';

import '../dates.dart';
import '../essays.dart';

/// The list of published essays, newest first, for the home page and /writing.
///
/// ```markdown
/// <EssayList/>
/// <EssayList limit="3" heading="Writing" id="writing"/>
/// ```
///
/// An essay is any page with `kind: essay`; pages marked `noindex` or `draft` are left
/// out ([publishedEssays]). With no essays, the plain list says "Nothing published yet.";
/// given a [heading], it is a whole section with its own `h2`, left out entirely until
/// there is something to list. Needs `eagerlyLoadAllPages` on the ContentApp, or pages
/// built later are missing.
class EssayList extends StatelessComponent {
  const EssayList({this.limit, this.heading, this.id, super.key});

  /// How many essays to show; all of them when null.
  final int? limit;

  /// The section head drawn above the list, or none when null.
  final String? heading;

  /// The heading's id, the target of its § link and of links such as `/#writing`.
  final String? id;

  static CustomComponent component() => CustomComponent(
    pattern: 'EssayList',
    builder: (_, attrs, _) =>
        EssayList(limit: int.tryParse(attrs['limit'] ?? ''), heading: attrs['heading'], id: attrs['id']),
  );

  @override
  Component build(BuildContext context) {
    final essays = publishedEssays(pageEntries(context.pages));
    if (essays.isEmpty) {
      return heading == null ? p(classes: 'empty', [.text('Nothing published yet.')]) : Component.fragment([]);
    }
    final shown = limit == null ? essays : essays.take(limit!);
    final list = div(classes: 'entries', [
      for (final essay in shown)
        p(classes: 'entry', [
          a(classes: 't', href: essay.url, [.text('${essay.data['title']}')]),
          .text(', '),
          span(classes: 'm', [
            Component.element(
              tag: 'time',
              attributes: {'datetime': '${essay.data['date'] ?? ''}'},
              children: [.text(formatMonth('${essay.data['date'] ?? ''}'))],
            ),
            .text('.'),
          ]),
          if (essay.data['abstract'] ?? essay.data['description'] case final String summary) .text(' $summary'),
        ]),
      if (limit != null && essays.length > limit!)
        p(classes: 'more', [
          a(href: '/writing', [.text('All ${essays.length} essays')]),
        ]),
    ]);
    if (heading case final String text) {
      // The markup EssaySectionsExtension gives a plain page's h2, which it can't reach
      // here: the heading and its § link.
      final anchor = id ?? text.toLowerCase();
      return Component.fragment([
        h2(id: anchor, [
          .text(text),
          a(classes: 'anchor', href: '#$anchor', attributes: {'aria-label': 'Link to this section'}, [.text('§')]),
        ]),
        list,
      ]);
    }
    return list;
  }
}
