import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';
import 'package:jaspr_content/jaspr_content.dart';

/// Markdown components for Knuth-style algorithms and exercises (DESIGN.md, "Components").
///
/// ```markdown
/// <Algorithm id="algorithm-b" letter="B" title="Backpropagation">
/// One sentence on the input and what it computes.
///
/// <Steps>
/// <Step n="B1" label="Initialize.">Set …</Step>
/// <Step n="B3" label="Decrease *i*." end>…</Step>
/// </Steps>
/// </Algorithm>
/// ```
List<CustomComponent> knuthComponents() => [
  CustomComponent(
    pattern: 'Algorithm',
    builder: (_, attrs, child) {
      final letter = attrs['letter'] ?? '';
      return div(
        classes: 'algorithm',
        id: attrs['id'],
        attributes: {'role': 'group', 'aria-label': 'Algorithm $letter'},
        [
          // The lead runs into the first paragraph of the child (see site.css).
          span(classes: 'lead', [
            span(classes: 'name', [.text('Algorithm $letter')]),
            .text(' ('),
            i([.text(attrs['title'] ?? '')]),
            .text('). '),
          ]),
          ?child,
        ],
      );
    },
  ),
  CustomComponent(pattern: 'Steps', builder: (_, _, child) => ol([?child])),
  CustomComponent(
    pattern: 'Step',
    builder: (_, attrs, child) {
      final n = attrs['n'] ?? '';
      return li(id: 'step-${n.toLowerCase()}', [
        span(classes: 'step', [.text('$n.')]),
        .text('['),
        ..._italics(attrs['label'] ?? ''),
        .text('] '),
        ?child,
        if (attrs.containsKey('end')) ...[
          .text(' '),
          span(classes: 'qed', attributes: {'aria-label': 'End of algorithm'}, [.text('▮')]),
        ],
      ]);
    },
  ),
  CustomComponent(
    pattern: 'Exercises',
    builder: (_, _, child) => Component.fragment([
      h2(classes: 'exercises', id: 'exercises', [.text('Exercises')]),
      ol(classes: 'ex', [?child]),
    ]),
  ),
  CustomComponent(
    pattern: 'Exercise',
    builder: (_, attrs, child) {
      final n = attrs['n'] ?? '';
      return li(id: 'exercise-$n', [
        if (attrs.containsKey('recommended'))
          span(classes: 'rec', attributes: {'aria-label': 'Recommended'}, [.text('▸')]),
        a(classes: 'no', href: '#answer-$n', attributes: {'aria-label': 'Exercise $n, see the answer'}, [.text('$n.')]),
        span([
          span(classes: 'rating', [
            span(classes: 'b', [.text('[')]),
            .text(attrs['rating'] ?? ''),
            span(classes: 'b', [.text(']')]),
          ]),
          ?child,
        ]),
      ]);
    },
  ),
  CustomComponent(
    pattern: 'Answers',
    builder: (_, _, child) => Component.fragment([
      h2(classes: 'exercises', id: 'answers', [.text('Answers to exercises')]),
      ol(classes: 'ex answers', [?child]),
    ]),
  ),
  CustomComponent(
    pattern: 'Answer',
    builder: (_, attrs, child) {
      final n = attrs['n'] ?? '';
      return li(id: 'answer-$n', [
        a(classes: 'no', href: '#exercise-$n', attributes: {'aria-label': 'Back to exercise $n'}, [.text('$n.')]),
        span([?child]),
      ]);
    },
  ),
];

/// Turns `*x*` in a step label into italics; labels are attributes, not Markdown.
List<Component> _italics(String label) {
  final out = <Component>[];
  final parts = label.split('*');
  for (var k = 0; k < parts.length; k++) {
    if (parts[k].isEmpty) continue;
    out.add(k.isOdd ? i([.text(parts[k])]) : .text(parts[k]));
  }
  return out;
}
