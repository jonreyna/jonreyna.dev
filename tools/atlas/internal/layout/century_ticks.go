package layout

import (
	"math"
	"strconv"
)

// The least room between two labeled ticks, and the years between ticks.
const (
	tickSpacing = 52.0
	century     = 100
)

// centuryTicks returns a tick at every century in the span, labeling those with room.
// Year 0 doesn't exist, so the tick between 100 BC and AD 100 is AD 1.
func centuryTicks(span [2]int, s *timeScale) []Tick {

	var ticks []Tick

	lastLabel := math.Inf(-1)
	for y := int(math.Ceil(float64(span[0])/century)) * century; y <= span[1]; y += century {
		year := y
		if year == 0 {
			year = 1
		}

		t := Tick{Year: year, X: s.x(float64(year))}
		if t.X-lastLabel >= tickSpacing {
			t.Label = yearLabel(year)
			lastLabel = t.X
		}

		ticks = append(ticks, t)
	}

	return ticks
}

// yearLabel writes a year as "300 BC" or "AD 200".
func yearLabel(year int) string {

	if year < 0 {
		return strconv.Itoa(-year) + " BC"
	}

	return "AD " + strconv.Itoa(year)
}
