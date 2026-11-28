package layout

// timeScale maps years to x by straight lines between its points, which run forward
// in both year and x.
type timeScale struct {
	years []float64
	xs    []float64
}

// add appends a point, skipping one at a year already placed.
func (s *timeScale) add(year, x float64) {

	if n := len(s.years); n > 0 && year <= s.years[n-1] {
		return
	}

	s.years = append(s.years, year)
	s.xs = append(s.xs, x)
}

// x returns the x for a year, clamped to the scale's ends.
func (s *timeScale) x(year float64) float64 {

	n := len(s.years)
	if n == 0 {
		return 0
	}

	if year <= s.years[0] {
		return s.xs[0]
	}

	for i := 1; i < n; i++ {
		if year <= s.years[i] {
			f := (year - s.years[i-1]) / (s.years[i] - s.years[i-1])

			return s.xs[i-1] + f*(s.xs[i]-s.xs[i-1])
		}
	}

	return s.xs[n-1]
}
