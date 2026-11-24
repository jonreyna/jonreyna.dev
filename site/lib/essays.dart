import 'package:jaspr_content/jaspr_content.dart';

/// A page's URL and its front matter.
typedef PageEntry = ({String url, Map<String, Object?> data});

/// Returns the published essays among [pages], newest first by `date`: pages with
/// `kind: essay` that aren't marked `noindex` or `draft`.
///
/// The essay lists, the Writing nav link, the feed link and the Writing page's indexing
/// all follow this, so the site shows Writing only once something is published.
List<PageEntry> publishedEssays(Iterable<PageEntry> pages) => [
  for (final page in pages)
    if (page.data['kind'] == 'essay' && page.data['noindex'] != true && page.data['draft'] != true) page,
]..sort((x, y) => '${y.data['date'] ?? ''}'.compareTo('${x.data['date'] ?? ''}'));

/// The pages being built, as [PageEntry]s for [publishedEssays]. Needs
/// `eagerlyLoadAllPages` on the ContentApp, or pages built later are missing.
Iterable<PageEntry> pageEntries(List<Page> pages) => [for (final page in pages) (url: page.url, data: page.data.page)];
