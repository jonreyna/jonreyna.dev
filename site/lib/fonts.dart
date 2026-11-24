/// The self-hosted fonts' manifest, read during the build for the preload links.
///
/// `make fonts` subsets the fonts into web/fonts/ under content-hashed names
/// (tool/subset_fonts.py) and lists them in a manifest; `make site` passes its path in
/// the FONT_FACES environment variable.
library;

import 'dart:convert';
import 'dart:io';

List<String>? _preload;

/// The URLs of the fonts every page preloads, such as
/// `/fonts/alegreya-regular.68c44f2872.woff2`. Throws if FONT_FACES isn't set or its
/// manifest can't be read.
List<String> preloadFonts() => _preload ??= _readPreload();

List<String> _readPreload() {
  final path = Platform.environment['FONT_FACES'];
  if (path == null || path.isEmpty) {
    throw StateError('the layout needs the font manifest, but FONT_FACES is not set: build with `make site`');
  }
  final manifest = jsonDecode(File(path).readAsStringSync());
  if (manifest case {'preload': final List<dynamic> urls} when urls.isNotEmpty) {
    return [for (final url in urls) '$url'];
  }
  throw StateError('$path: no "preload" URLs');
}
