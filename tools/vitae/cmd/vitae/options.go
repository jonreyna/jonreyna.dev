package main

// options are the flags every command accepts; each command uses the ones it needs.
type options struct {
	input, private, variant, format, outDir string
	strict                                  bool
}
