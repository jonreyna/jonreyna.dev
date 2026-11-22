package lint

// Severity is how serious an [Issue] is.
type Severity int

// The severities.
const (
	// Warning is a problem that fails validation only in strict mode.
	Warning Severity = iota
	// Error is a problem that always fails validation.
	Error
)

// String returns "warning" or "error".
func (s Severity) String() string {

	if s == Error {
		return "error"
	}

	return "warning"
}
