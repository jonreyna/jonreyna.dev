package resume

import (
	"fmt"
	"strconv"

	"go.yaml.in/yaml/v3"
)

// monthLayoutLen is the length of a month written "YYYY-MM".
const monthLayoutLen = len("2006-01")

// Month is a calendar month, written "YYYY-MM" in the file and in JSON.
type Month struct {
	// Year is the four-digit year.
	Year int
	// Month is the month of the year, 1 to 12.
	Month int
}

// ParseMonth reads a month written "YYYY-MM". The error wraps [ErrInvalidMonth].
func ParseMonth(s string) (Month, error) {
	if len(s) != monthLayoutLen || s[4] != '-' {
		return Month{}, fmt.Errorf("%w: %q", ErrInvalidMonth, s)
	}
	y, errY := strconv.Atoi(s[:4])
	m, errM := strconv.Atoi(s[5:])
	if errY != nil || errM != nil || m < 1 || m > 12 {
		return Month{}, fmt.Errorf("%w: %q", ErrInvalidMonth, s)
	}
	return Month{Year: y, Month: m}, nil
}

// String returns the month as "YYYY-MM".
func (m Month) String() string { return fmt.Sprintf("%04d-%02d", m.Year, m.Month) }

// Before reports whether m is earlier than o.
func (m Month) Before(o Month) bool {
	return m.Year < o.Year || (m.Year == o.Year && m.Month < o.Month)
}

// IsZero reports whether m is unset.
func (m Month) IsZero() bool { return m == Month{} }

// UnmarshalYAML decodes a "YYYY-MM" scalar.
func (m *Month) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.ScalarNode {
		return fmt.Errorf("line %d: %w: not a string", n.Line, ErrInvalidMonth)
	}
	v, err := ParseMonth(n.Value)
	if err != nil {
		return fmt.Errorf("line %d: %w", n.Line, err)
	}
	*m = v
	return nil
}

// MarshalJSON encodes the month as a "YYYY-MM" string.
func (m Month) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(m.String())), nil
}
