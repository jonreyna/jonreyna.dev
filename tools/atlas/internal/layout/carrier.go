package layout

// carrier is one thread's line leaving one knot: it runs right in a lane's track
// from x0 to x1, and each of its targets branches off it.
type carrier struct {
	thread  string
	source  int
	lane    int
	x0, x1  float64
	targets []carrierTarget
	track   int
	y       float64
}

// active reports whether the carrier runs past x.
func (c *carrier) active(x float64) bool { return c.x0 < x && x < c.x1 }

// reaches reports whether spot is one of the carrier's targets.
func (c *carrier) reaches(spot int) bool {

	for _, t := range c.targets {
		if t.spot == spot {
			return true
		}
	}

	return false
}
