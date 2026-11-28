package render

import (
	"fmt"
	"io"

	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/layout"
	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"
)

// HTML writes the piece's markup: one <main id="main"> with the page's only <h1>.
func HTML(w io.Writer, p *plate.Plate, l *layout.Layout) error {

	pg := newPage(p, l)
	pg.head()
	pg.legend()
	pg.chart()
	pg.cards()
	pg.outline()
	pg.line(`</main>`)

	_, err := io.WriteString(w, pg.b.String())
	if err != nil {
		return fmt.Errorf("write markup: %w", err)
	}

	return nil
}
