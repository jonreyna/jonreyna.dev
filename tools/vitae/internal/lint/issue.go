package lint

import "fmt"

// Issue is one problem found by [Check].
type Issue struct {
	// Severity says whether the issue fails validation outside strict mode.
	Severity Severity
	// Where locates the problem by path, such as "experience[equinix].highlights[2]".
	Where string
	// Msg describes the problem and, where it helps, the fix.
	Msg string
}

// String formats the issue as "severity: where: message".
func (i Issue) String() string { return fmt.Sprintf("%s: %s: %s", i.Severity, i.Where, i.Msg) }
