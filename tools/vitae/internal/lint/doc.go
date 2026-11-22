// Package lint checks a résumé's structure and content before anything is rendered.
//
// [Check] returns every [Issue] it finds, in file order. An issue is an [Error], which
// always fails validation, or a [Warning], which fails only in strict mode ([HasErrors]
// with strict set; CI runs strict). The rules:
//
//   - Structure: schema version, required fields, unique role IDs, valid dates, roles
//     newest first with only the first one current, https links.
//   - Tags: every tag used is declared under tags, and every declared tag is used.
//   - Content: inline markup parses; highlights stay under [MaxHighlight] characters;
//     curly quotes and real dashes, not straight quotes or hyphens.
//   - Privacy: the public file has no phone number; it belongs in the private overlay.
//   - Education: an entry states whether its degree was conferred, and if it wasn't,
//     only its note may mention a degree.
//   - Variants: "full" exists; sections, summaries and tags they name exist; and each
//     selects at least one highlight.
package lint
