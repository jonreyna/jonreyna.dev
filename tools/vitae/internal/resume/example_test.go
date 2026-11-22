package resume_test

import (
	"fmt"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

func ExampleParseMonth() {

	m, err := resume.ParseMonth("2023-10")
	fmt.Println(m.Year, m.Month, err)

	_, err = resume.ParseMonth("October 2023")
	fmt.Println(err)
	// Output:
	// 2023 10 <nil>
	// not a YYYY-MM month: "October 2023"
}

func ExampleMonth_Before() {

	start := resume.Month{Year: 2020, Month: 4}
	end := resume.Month{Year: 2023, Month: 10}
	fmt.Println(start.Before(end), end.Before(start))
	// Output: true false
}
