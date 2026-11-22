import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';
import 'package:jaspr_content/jaspr_content.dart';

import '../dates.dart';

/// The list of published essays, newest first, for the home page and /writing.
///
/// ```markdown
/// <EssayList limit="3"/>
/// ```
///
/// An essay is any page with `kind: essay`; pages marked `noindex` or `draft` are left
/// out. Needs `eagerlyLoadAllPages` on the ContentApp, or pages built later are missing.
class EssayList extends StatelessComponent {
  const EssayList({this.limit, super.key});

  /// How many essays to show; all of them when null.
  final int? limit;

  static CustomComponent component() => CustomComponent(
    pattern: 'EssayList',
    builder: (_, attrs, _) => EssayList(limit: int.tryParse(attrs['limit'] ?? '')),
  );

  @override
  Component build(BuildContext context) {
    final essays = [
      for (final page in context.pages)
        if (page.data.page case final data
            when data['kind'] == 'essay' && data['noindex'] != true && data['draft'] != true)
          (url: page.url, data: data),
    ]..sort((x, y) => '${y.data['date'] ?? ''}'.compareTo('${x.data['date'] ?? ''}'));

    if (essays.isEmpty) {
      return p(classes: 'empty', [.text('Nothing published yet.')]);
    }
    final shown = limit == null ? essays : essays.take(limit!);
    return div(classes: 'entries', [
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
  }
}
