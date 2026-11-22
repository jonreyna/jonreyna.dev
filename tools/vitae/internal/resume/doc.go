// Package resume holds the résumé as written in resume.yaml, and loads it.
//
// [Resume] is the whole file: the header ([Basics]), keyed summaries, declared tags,
// roles ([Role] with [Highlight]s), skills ([SkillGroup]), education ([School] with
// [Course]s) and the [Variant]s that select from them. [Load] decodes a file and
// rejects unknown fields, so a misspelled key fails instead of vanishing. [ApplyPrivate]
// merges the git-ignored overlay ([Private]) that holds the phone number.
//
// This package only decodes. Package lint validates a [Resume], package variant applies
// a [Variant], and package render writes the result.
//
// Dates in the file are months written "YYYY-MM" ([Month]); education uses plain years.
// A role whose End is nil is current.
package resume
