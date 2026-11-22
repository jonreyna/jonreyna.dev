// Server-only: this component reads files with dart:io while the site is pre-rendered,
// so never import it from main.client.dart.
import 'dart:convert';
import 'dart:io';

import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';

/// The résumé page: the full résumé as HTML, and a link to every rendered file.
///
/// ```markdown
/// <Resume/>
/// ```
///
/// It reads `web/resume/`, which `make site-resume` fills from `vitae render` (the
/// public fields only). The page shows `resume.full.json`; the downloads list whatever
/// files are there, one row per variant, so a format appears as soon as the pipeline
/// produces it and no link is ever dead. Text fields carry HTML that vitae escaped and
/// built from the résumé's own markup, so the page inserts it as is.
class ResumeView extends StatelessComponent {
  const ResumeView({this.dir = 'web/resume', super.key});

  /// The folder of rendered files, relative to `site/`.
  final String dir;

  /// The URL path the files are served under (Jaspr copies `web/` into the build).
  static const urlBase = '/resume';

  @override
  Component build(BuildContext context) {
    final files = ResumeFiles.read(dir);
    final doc = files.document('full');
    final basics = _map(doc['basics']);

    return div(classes: 'resume', [
      _contact(basics),
      _downloadLine(files.byVariant['full'] ?? const []),
      for (final section in _list(doc['sections']).cast<String>()) ..._section(section, doc),
      _heading('downloads', 'Downloads'),
      ul(classes: 'plain downloads', [
        for (final variant in files.variants)
          li([
            .text('${files.document(variant)['description'] ?? variant} '),
            span(classes: 'files', _fileLinks(files.byVariant[variant]!)),
          ]),
      ]),
    ]);
  }

  /// Label, location, email and profiles, as the home page's contact line.
  Component _contact(Map<String, Object?> b) {
    final location = b['remote'] == true ? '${b['location']} (remote)' : '${b['location']}';
    return p(
      classes: 'contact',
      _dotted([
        span(classes: 'nowrap', [.text('${b['label']}')]),
        span(classes: 'nowrap', [.text(location)]),
        a(href: 'mailto:${b['email']}', [.text('${b['email']}')]),
        for (final link in _list(b['links']).map(_map))
          if (link['kind'] != 'site') a(classes: 'bare', href: '${link['url']}', [.text('${link['text']}')]),
      ]),
    );
  }

  /// The full résumé's files, at the top so a recruiter finds the PDF at once.
  Component _downloadLine(List<String> names) => p(classes: 'contact', [.text('Download: '), ..._fileLinks(names)]);

  Iterable<Component> _section(String name, Map<String, Object?> doc) sync* {
    switch (name) {
      case 'summary':
        yield p(classes: 'summary', [_html(doc['summary'])]);
      case 'skills':
        yield _heading('skills', 'Skills');
        for (final g in _list(doc['skills']).map(_map)) {
          yield p(classes: 'run', [
            span(classes: 't', [.text('${g['group']}.')]),
            .text(' ${_list(g['items']).join(', ')}.'),
          ]);
        }
      case 'experience':
        yield _heading('experience', 'Experience');
        for (final role in _list(doc['experience']).map(_map)) {
          yield div(classes: 'role', [
            p(classes: 'run', [
              span(classes: 't', [.text('${role['org']}')]),
              .text(', ${role['title']}. '),
              span(classes: 'm', [.text('${role['location']}, ${role['dates']}.')]),
            ]),
            ul(classes: 'plain', [
              for (final h in _list(role['highlights']).map(_map))
                li([
                  if (h['lead'] case final String lead when lead.isNotEmpty) ...[
                    span(classes: 'lead', [.text('$lead.')]),
                    .text(' '),
                  ],
                  _html(h['text']),
                ]),
            ]),
          ]);
        }
      case 'earlier':
        yield _heading('earlier', 'Earlier');
        yield p([_html(doc['earlier'])]);
      case 'education':
        yield _heading('education', 'Education');
        for (final s in _list(doc['education']).map(_map)) {
          yield div(classes: 'role', [
            p(classes: 'run', [
              span(classes: 't', [.text('${s['org']}')]),
              .text(', ${s['degree'] ?? s['program']}. '),
              span(classes: 'm', [.text('${s['dates']}.')]),
            ]),
            ul(classes: 'plain', [
              if (_list(s['minors']).isNotEmpty) li([.text('Minors in ${_and(_list(s['minors']))}.')]),
              if (s['note'] != null) li([_html(s['note'])]),
              if (s['gpa'] != null || _list(s['honors']).isNotEmpty)
                li([
                  .text(
                    [
                      if (s['gpa'] != null) 'GPA ${s['gpa']}',
                      ..._list(s['honors']).map((h) => '$h'),
                    ].join('; '),
                  ),
                  .text('.'),
                ]),
              if (_list(s['courses']).isNotEmpty) li([.text('Coursework: ${_list(s['courses']).join(', ')}.')]),
            ]),
          ]);
        }
    }
  }

  /// A section head with an id and a § link, like the headings the build adds to
  /// Markdown pages (lib/essay_sections.dart).
  Component _heading(String id, String title) => h2(id: id, [
    .text(title),
    a(classes: 'anchor', href: '#$id', attributes: {'aria-label': 'Link to section $title'}, [.text('§')]),
  ]);

  List<Component> _fileLinks(List<String> names) => _dotted([
    for (final name in names) a(href: '$urlBase/$name', [.text(ResumeFiles.formatLabel(name))]),
  ]);

  /// Puts the contact line's red middle dot between [parts]. The space after each dot
  /// is where a narrow screen may break the line, since the parts themselves don't.
  static List<Component> _dotted(List<Component> parts) => [
    for (final (i, part) in parts.indexed) ...[
      if (i > 0) ...[
        span(classes: 'sep', attributes: {'aria-hidden': 'true'}, [.text('·')]),
        .text(' '),
      ],
      part,
    ],
  ];

  /// The HTML half of one of vitae's text fields ({"plain": …, "html": …}).
  static Component _html(Object? text) => RawText('${_map(text)['html'] ?? ''}');

  /// "a, b and c", the site's list style (no serial comma).
  static String _and(List<Object?> items) => switch (items.length) {
    0 => '',
    1 => '${items.first}',
    _ => '${items.take(items.length - 1).join(', ')} and ${items.last}',
  };

  static Map<String, Object?> _map(Object? v) => v is Map ? v.cast<String, Object?>() : const {};
  static List<Object?> _list(Object? v) => v is List ? v : const [];
}

/// The files `vitae render` wrote, named `resume.<variant>.<format>`, grouped by
/// variant.
class ResumeFiles {
  ResumeFiles._(this._dir, this.byVariant);

  /// Lists [dir]. Throws if it is missing or has no full résumé, so that a build
  /// without `make site-resume` fails instead of publishing an empty page.
  factory ResumeFiles.read(String dir) {
    final folder = Directory(dir);
    if (!folder.existsSync()) {
      throw StateError('$dir is missing: run `make site-resume` (or `make site`) first');
    }
    final byVariant = <String, List<String>>{};
    for (final entry in folder.listSync().whereType<File>()) {
      final name = entry.uri.pathSegments.last;
      final parts = name.split('.');
      if (parts.length != 3 || parts.first != 'resume') continue;
      byVariant.putIfAbsent(parts[1], () => []).add(name);
    }
    for (final names in byVariant.values) {
      names.sort((x, y) => _formatRank(x).compareTo(_formatRank(y)));
    }
    if (!byVariant.containsKey('full') || !byVariant['full']!.contains('resume.full.json')) {
      throw StateError('$dir has no resume.full.json: run `make site-resume` first');
    }
    return ResumeFiles._(dir, byVariant);
  }

  final String _dir;

  /// File names per variant, in the order of [formatLabels].
  final Map<String, List<String>> byVariant;

  /// Variants in display order: full, then onepage, then the rest by name.
  List<String> get variants {
    const first = ['full', 'onepage'];
    return [
      ...first.where(byVariant.containsKey),
      ...(byVariant.keys.where((v) => !first.contains(v)).toList()..sort()),
    ];
  }

  /// The decoded JSON document for [variant].
  Map<String, Object?> document(String variant) {
    final file = File('$_dir/resume.$variant.json');
    if (!file.existsSync()) return const {};
    return (jsonDecode(file.readAsStringSync()) as Map).cast<String, Object?>();
  }

  /// Link text for each format, in the order the downloads show them.
  static const formatLabels = {
    'pdf': 'PDF',
    'docx': 'Word',
    'txt': 'plain text',
    'md': 'Markdown',
    'json': 'JSON',
  };

  /// The link text for a file name's format.
  static String formatLabel(String name) {
    final ext = name.split('.').last;
    return formatLabels[ext] ?? ext.toUpperCase();
  }

  static int _formatRank(String name) {
    final i = formatLabels.keys.toList().indexOf(name.split('.').last);
    return i < 0 ? formatLabels.length : i;
  }
}
