/// Runs the pinned Pandoc binary during the build, for TeX math (tex_math.dart) and
/// code highlighting (code_highlight.dart).
///
/// `make` downloads Pandoc and checks its SHA-256 (see the Makefile), and `make site`
/// and `make site-test` pass its path in the PANDOC environment variable.
library;

import 'dart:convert';
import 'dart:io';

/// The path of the Pandoc binary, from PANDOC. Throws if it isn't set, naming [what]
/// needs it.
String pandocPath(String what) {
  final pandoc = Platform.environment['PANDOC'];
  if (pandoc == null || pandoc.isEmpty) {
    throw StateError('$what needs Pandoc, but PANDOC is not set: build with `make site`');
  }
  return pandoc;
}

/// Runs Pandoc with [args] and [input] on stdin, and returns its stdout. Warnings
/// fail the run (`--fail-if-warnings`); a failure throws a [StateError] that starts
/// with [what], such as the page's path.
Future<String> runPandoc(String pandoc, List<String> args, String input, String what) async {
  final process = await Process.start(pandoc, [...args, '--fail-if-warnings']);
  process.stdin.write(input);
  await process.stdin.close();
  final (stdout, stderr, code) = await (
    process.stdout.transform(utf8.decoder).join(),
    process.stderr.transform(utf8.decoder).join(),
    process.exitCode,
  ).wait;
  if (code != 0) {
    throw StateError('$what: Pandoc failed (exit $code): $stderr');
  }
  return stdout;
}
