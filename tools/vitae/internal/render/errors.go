package render

import "errors"

// ErrUnknownFormat reports a format name that [NewRenderer] doesn't know.
var ErrUnknownFormat = errors.New("unknown format")
