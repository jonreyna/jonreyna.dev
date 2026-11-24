package icon_test

import (
	"fmt"
	"strings"

	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/sfnt"

	"github.com/jonreyna/jonreyna.dev/tools/card/internal/icon"
)

// The build passes Alegreya SC Medium; Go Regular stands in here.
func ExampleIcon_SVG() {

	f, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		fmt.Println(err)
		return
	}

	ic := icon.Icon{Letter: 'J', Font: f, Height: 0.8}
	svg, err := ic.SVG(icon.Light(), icon.Dark())
	if err != nil {
		fmt.Println(err)
		return
	}

	head, rest, _ := strings.Cut(string(svg), "<style>")
	fmt.Println(head)
	fmt.Println(strings.Contains(rest, `<path class="l" d="M`))
	// Output:
	// <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
	// true
}

func ExampleIcon_Render() {

	f, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		fmt.Println(err)
		return
	}

	ic := icon.Icon{Letter: 'J', Font: f, Height: 0.62}
	img, err := ic.Render(180, icon.Light())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(img.Bounds(), img.RGBAAt(0, 0) == icon.Light().Ground)
	// Output: (0,0)-(180,180) true
}
