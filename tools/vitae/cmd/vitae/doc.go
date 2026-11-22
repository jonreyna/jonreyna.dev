// Command vitae validates resume.yaml and renders it, one file per variant and format.
//
// Usage:
//
//	vitae validate [-i FILE] [-private FILE] [-strict]
//	vitae render   [-i FILE] [-private FILE] [-strict] [-variant NAME|all] [-format NAME|all] [-o DIR|-]
//	vitae list     [-i FILE] variants|tags
//	vitae version
//
// validate prints every issue lint finds and fails on errors (on warnings too with
// -strict). render validates the same way, refuses to render a résumé with errors,
// then writes resume.<variant>.<format> into DIR for each variant and format, printing
// "wrote PATH" after each file. With -o - it writes a single file to standard output.
// list prints the variants with their descriptions, or the tags with how many
// highlights use each. version prints the commit vitae was built from.
//
// Flags:
//
//	-i FILE         the résumé (default resume/resume.yaml)
//	-private FILE   merge this private overlay (the phone number); local builds only
//	-strict         treat warnings as errors
//	-variant NAME   the variant to render, or all (default all)
//	-format NAME    the format to render, or all (default all)
//	-o DIR          the output directory, or - for standard output (default .)
//
// Exit codes: 0 success, 1 validation or I/O failure, 2 usage error.
package main
