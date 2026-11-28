package plate_test

import (
	"fmt"

	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"
)

func ExampleParse() {

	p, err := plate.Parse([]byte(`
plate: tiny
title: Tiny
span: [-500, 100]
lanes: [{id: a, name: A}]
threads: [{id: t, name: T, color: gold}]
people:
  - {id: x, name: Xenophanes, lanes: [a], from: -570, to: -475, circa: true, card: One god.}
  - {id: c, name: Clement, lanes: [a], from: 150, to: 215, card: Quotes him.}
links:
  - {id: L1, thread: t, from: x, to: c, type: quotes, cite: Stromata V.14, paraphrase: Quotes him., located: paraphrase}
`))
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(p.Person("x").Dates())
	fmt.Println(plate.Validate(p, false))
	fmt.Println(plate.Validate(p, true))
	// Output:
	// c. 570–475 BC
	// []
	// [link L1: not checked]
}
