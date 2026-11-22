package markup

import "errors"

// Errors from [Parse]. Each is wrapped with the position of the problem.
var (
	// ErrTrailingBackslash reports a backslash with nothing after it.
	ErrTrailingBackslash = errors.New("trailing backslash")
	// ErrUnclosed reports emphasis or code that is never closed.
	ErrUnclosed = errors.New("unclosed markup")
	// ErrEmpty reports emphasis or code with nothing inside.
	ErrEmpty = errors.New("empty markup")
	// ErrBadLink reports a [ that doesn't start a well-formed [text](url) link, or a
	// URL that isn't https:// or mailto:.
	ErrBadLink = errors.New("malformed link")
)
