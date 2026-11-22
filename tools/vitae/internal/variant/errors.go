package variant

import "errors"

// ErrUnknown reports a variant name the résumé doesn't define.
var ErrUnknown = errors.New("unknown variant")
