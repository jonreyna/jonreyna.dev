// Package plate reads and checks an atlas plate: the data behind a woven chart such
// as play/athens-and-jerusalem/athens-and-jerusalem.yaml.
//
// A [Plate] has lanes (the bands of the chart, such as "Philosophers"), threads (the
// ideas, each drawn in one color), people (the knots; a book can be one) and links.
// A [Link] says that one person took up, named, quoted or rejected another's idea on
// a thread, and it carries the citation for that claim: the passage, the translation
// it was copied from, and whether Jon has checked it.
//
// [Parse] decodes the YAML strictly, so a misspelled field is an error, and indexes
// the people by id. [Validate] then lists every problem it finds, rather than
// stopping at the first, so one run shows everything to fix. With requireChecked set
// it also reports every link Jon hasn't checked, which is how a production build
// refuses to publish an unchecked citation.
//
// Years are plain numbers, BC negative, and there is no year 0.
package plate
