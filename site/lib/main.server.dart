/// Server entry point: pre-renders every page in content/ at build time.
library;

import 'package:jaspr/server.dart';
import 'package:jaspr_content/jaspr_content.dart';
import 'package:jaspr_content/theme.dart';

import 'code_highlight.dart';
import 'components/essay_list.dart';
import 'components/knuth.dart';
import 'components/resume_view.dart';
import 'essay_sections.dart';
import 'layouts/j_layout.dart';
import 'main.server.options.dart';
import 'tex_math.dart';

void main() {
  // Jaspr adds main.client.dart.js to every page whenever a client entry point exists.
  // Leaving clientId out stops that; JLayout adds the script only to pages that set
  // `islands: true`, so reading pages ship no JavaScript (DESIGN.md, principle 4).
  final generated = defaultServerOptions;
  Jaspr.initializeApp(
    options: ServerOptions(clients: generated.clients, stylesId: generated.stylesId, styles: generated.styles),
  );

  runApp(
    ContentApp(
      // The writing index and the home page list other pages, so load them all first.
      eagerlyLoadAllPages: true,
      // TeX math ($…$, $$…$$) becomes MathML before the Markdown is parsed.
      templateEngine: const TexMath(),
      parsers: [MarkdownParser()],
      // Code blocks are highlighted while the site is built (Pandoc, as for the math).
      extensions: [const EssaySectionsExtension(), const CodeHighlightExtension()],
      components: [
        ...knuthComponents(),
        EssayList.component(),
        CustomComponent(pattern: 'Resume', builder: (_, _, _) => const ResumeView()),
      ],
      layouts: [const JLayout()],
      // Design J supplies all styles; turn off jaspr_content's own theme.
      theme: const ContentTheme.none(),
    ),
  );
}
