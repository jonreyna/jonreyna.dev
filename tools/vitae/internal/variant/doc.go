// Package variant applies a [resume.Variant] to a résumé, producing a [Resolved]
// résumé with only the sections, roles, highlights and courses the variant keeps.
//
// It's the only code that understands variants: every renderer, and the website,
// reads its output instead of filtering again.
//
// Selection works like this, per role:
//
//  1. Roles that ended before Since are dropped; at most MaxRoles remain, newest first.
//  2. Highlights carrying an Exclude tag are dropped.
//  3. Each remaining highlight scores its Weight plus BoostBy (default [DefaultBoost])
//     for every Boost tag it carries. Scores under MinWeight are dropped.
//  4. If more than MaxHighlights remain, the highest scores stay; ties go to the
//     earlier highlight.
//  5. The survivors keep their order from the file. Variants select; they never reorder
//     or reword.
//
// Courses are kept when they carry a ShowCourses tag, or all of them when ShowCourses
// is empty.
package variant
