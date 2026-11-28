package render

import (
	"math"
	"strconv"
	"strings"
)

// tenths is what num rounds to.
const tenths = 10

// num writes a coordinate with one decimal place at most.
func num(v float64) string {

	s := strconv.FormatFloat(math.Round(v*tenths)/tenths, 'f', 1, 64)

	return strings.TrimSuffix(s, ".0")
}
