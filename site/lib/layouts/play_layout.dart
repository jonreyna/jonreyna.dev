import 'dart:io';

import 'package:jaspr/dom.dart';
import 'package:jaspr/server.dart';
import 'package:jaspr_content/jaspr_content.dart';

import 'j_layout.dart';

/// A Play piece that fills the window, such as Arbor Ecclesiae: the site's head (title,
/// description, Open Graph, fonts, favicons), then the piece's own markup, styles and
/// scripts in place of the masthead, the site's styles and the footer. The piece
/// carries its own link back to the site. See play/README.md.
///
/// A page picks it with `layout: play` and names its piece, a directory under play/
/// in the repository, with `piece:`. Its scripts, in load order, are listed in
/// `scripts:` as files in that directory:
///
/// ```yaml
/// layout: play
/// piece: arbor-ecclesiae
/// scripts: [d3.min.js, arbor-ecclesiae.js]
/// ```
///
/// `make play` stages the piece's files where the build serves them.
///
/// A piece built from data, such as an atlas drawn by tools/atlas, sets
/// `generated: true`: its markup is then read from build/play/PIECE/PIECE.html, which
/// `make play` writes, instead of from play/PIECE/.
class PlayLayout extends JLayout {
  const PlayLayout();

  @override
  Pattern get name => 'play';

  @override
  List<String> stylesheets(Page page) => ['/fonts.css', '/play-${_piece(page)}.css'];

  @override
  List<String> scripts(Page page) => [
    for (final file in page.data.page['scripts'] as List? ?? const []) '/play/${_piece(page)}/$file',
  ];

  @override
  Component buildBody(Page page, Component child) =>
      RawText(_markup(_piece(page), generated: page.data.page['generated'] == true));

  /// The page's piece, from its `piece:` front matter.
  static String _piece(Page page) => switch (page.data.page['piece']) {
    final String piece when RegExp(r'^[a-z0-9-]+$').hasMatch(piece) => piece,
    final other => throw StateError('${page.url}: `layout: play` needs `piece:` naming a directory in play/, got $other'),
  };

  /// The piece's markup, play/PIECE/PIECE.html (or build/play/PIECE/PIECE.html for a
  /// generated piece), read once per build. The build runs in site/, so both are one
  /// level up.
  static String _markup(String piece, {required bool generated}) => _cache[piece] ??= () {
    final file = File(generated ? '../build/play/$piece/$piece.html' : '../play/$piece/$piece.html');
    if (!file.existsSync()) {
      throw StateError('${file.path} is missing; run `make play` to build the piece');
    }
    return file.readAsStringSync();
  }();

  static final _cache = <String, String>{};
}
