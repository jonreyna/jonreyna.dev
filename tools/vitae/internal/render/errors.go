package render

import "errors"

// ErrUnknownFormat reports a format name that [For] doesn't know.
var ErrUnknownFormat = errors.New("unknown format")
