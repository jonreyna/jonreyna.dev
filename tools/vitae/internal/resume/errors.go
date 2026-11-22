package resume

import "errors"

// ErrInvalidMonth reports a month that isn't written "YYYY-MM" with a month from 01 to 12.
var ErrInvalidMonth = errors.New("not a YYYY-MM month")
