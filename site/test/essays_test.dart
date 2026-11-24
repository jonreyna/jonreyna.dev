// Tests for lib/essays.dart: which pages count as published essays, and their order.
import 'package:jonreyna_dev/essays.dart';
import 'package:test/test.dart';

PageEntry page(String url, Map<String, Object?> data) => (url: url, data: data);

void main() {
  final cases = <(String, List<PageEntry>, List<String>)>[
    ('no pages, no essays', [], []),
    (
      'plain pages are not essays',
      [
        page('/', {'kind': 'page'}),
        page('/about', {'title': 'About'}),
      ],
      [],
    ),
    (
      'drafts and noindex essays are left out',
      [
        page('/writing/a', {'kind': 'essay', 'date': '2027-04-24', 'draft': true}),
        page('/writing/b', {'kind': 'essay', 'date': '2027-04-25', 'noindex': true}),
        page('/writing/c', {'kind': 'essay', 'date': '2027-04-26'}),
      ],
      ['/writing/c'],
    ),
    (
      'newest first, undated last',
      [
        page('/writing/old', {'kind': 'essay', 'date': '2027-04-24'}),
        page('/writing/undated', {'kind': 'essay'}),
        page('/writing/new', {'kind': 'essay', 'date': '2027-06-01'}),
      ],
      ['/writing/new', '/writing/old', '/writing/undated'],
    ),
  ];
  for (final (name, pages, want) in cases) {
    test(name, () => expect([for (final e in publishedEssays(pages)) e.url], want));
  }
}
