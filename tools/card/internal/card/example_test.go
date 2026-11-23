package card_test

import (
	"fmt"

	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/jonreyna/jonreyna.dev/tools/card/internal/card"
)

// The build passes Alegreya SC Medium and Alegreya Italic; the Go fonts stand in here.
func ExampleCard_Render() {

	fonts, err := card.ParseFonts(goregular.TTF, goitalic.TTF)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := card.Card{Name: "Jon Reyna", Motto: "distributed systems & AI infrastructure", Site: "jonreyna.dev"}
	img, err := c.Render(fonts)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(img.Bounds())
	// Output: (0,0)-(1200,630)
}
