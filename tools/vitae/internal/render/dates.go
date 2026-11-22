package render

import (
	"strconv"
	"time"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

// Dates formats a role's span as the résumé prints it: "Oct 2023 – Present" for a
// current role (end is nil), "Aug 2019 – Apr 2020" otherwise.
func Dates(start resume.Month, end *resume.Month) string {
	if end == nil {
		return month(start) + " – Present"
	}
	return month(start) + " – " + month(*end)
}

// month formats m as "Oct 2023".
func month(m resume.Month) string {
	return time.Month(m.Month).String()[:3] + " " + strconv.Itoa(m.Year)
}
