import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';
import 'package:jaspr_content/jaspr_content.dart';

/// Renders content/_data/resume.json, which jaspr_content loads at build time.
/// In the real site this file is `resume.full.json`, produced by `vitae` (Phase 4).
class ResumeView extends StatelessComponent {
  const ResumeView({super.key});

  @override
  Component build(BuildContext context) {
    final resume = context.page.data['resume'];
    if (resume is! Map) {
      return p([.text('No résumé data found.')]);
    }
    final work = (resume['work'] as List?) ?? const [];
    return div(classes: 'resume', [
      if (resume['summary'] case final String summary) p([.text(summary)]),
      h2([.text('Experience')]),
      for (final role in work.cast<Map>())
        div(classes: 'role', [
          p(classes: 'role-head', [
            .text('${role['company']}, ${role['title']}'),
            span(classes: 'm', [.text(', ${role['dates']}.')]),
          ]),
          ul([
            for (final h in (role['highlights'] as List? ?? const [])) li([.text('$h')]),
          ]),
        ]),
      p(classes: 'contact', [
        a(href: '/resume.pdf', [.text('PDF')]),
        span(classes: 'sep', attributes: {'aria-hidden': 'true'}, [.text('·')]),
        a(href: '/resume.txt', [.text('plain text')]),
        span(classes: 'sep', attributes: {'aria-hidden': 'true'}, [.text('·')]),
        a(href: '/resume.json', [.text('JSON')]),
      ]),
    ]);
  }
}
