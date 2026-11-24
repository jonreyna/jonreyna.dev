package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image/png"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/image/font/sfnt"

	"github.com/jonreyna/jonreyna.dev/tools/card/internal/icon"
)

// Exit codes (see doc.go).
const (
	exitFailure = 1
	exitUsage   = 2
)

// The icons' sizes in pixels and how tall the letter stands on each, as a fraction of
// the side. A favicon is tiny, so its letter nearly fills the tile; phones round the
// touch icon's corners, so its letter keeps clear of them.
const (
	tabHeight   = 0.8
	touchHeight = 0.62
	touchSize   = 180
)

// icoSizes are the favicon.ico's images: a browser tab's, and a high-density tab's.
func icoSizes() []int { return []int{16, 32} }

// filePerm is the files' mode: public images for the website, readable by all.
const filePerm = 0o644

// errUsage marks errors in the command line, which exit with exitUsage.
var errUsage = errors.New("usage")

// run parses args, draws the icons and writes them, reporting each file on stdout.
func run(args []string, stdout, stderr io.Writer) error {

	fs := flag.NewFlagSet("icon", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fontPath := fs.String("font", "", "the `font` the letter is set in")
	letter := fs.String("letter", "", "the one `character` to draw")
	svgPath := fs.String("svg", "", "the SVG favicon `file` to write")
	icoPath := fs.String("ico", "", "the favicon.ico `file` to write")
	touchPath := fs.String("touch", "", "the apple-touch-icon PNG `file` to write")

	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("%w: %w", errUsage, err)
	}

	for _, f := range []struct{ name, value string }{
		{"font", *fontPath}, {"letter", *letter}, {"svg", *svgPath}, {"ico", *icoPath}, {"touch", *touchPath},
	} {
		if f.value == "" {
			return fmt.Errorf("%w: -%s is required", errUsage, f.name)
		}
	}

	r, n := utf8.DecodeRuneInString(*letter)
	if n != len(*letter) || r == utf8.RuneError {
		return fmt.Errorf("%w: -letter must be one character, not %q", errUsage, *letter)
	}

	data, err := os.ReadFile(*fontPath)
	if err != nil {
		return fmt.Errorf("read font: %w", err)
	}

	f, err := sfnt.Parse(data)
	if err != nil {
		return fmt.Errorf("parse font %s: %w", *fontPath, err)
	}

	tab := icon.Icon{Letter: r, Font: f, Height: tabHeight}
	touch := icon.Icon{Letter: r, Font: f, Height: touchHeight}
	for _, w := range []func() (string, error){
		func() (string, error) { return writeSVG(&tab, *svgPath) },
		func() (string, error) { return writeICO(&tab, *icoPath) },
		func() (string, error) { return writePNG(&touch, touchSize, *touchPath) },
	} {
		line, err := w()
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(stdout, line)
		if err != nil {
			return fmt.Errorf("report: %w", err)
		}
	}

	return nil
}

// writeSVG writes ic's SVG favicon to path and returns the line reporting it.
func writeSVG(ic *icon.Icon, path string) (string, error) {

	svg, err := ic.SVG(icon.Light(), icon.Dark())
	if err != nil {
		return "", fmt.Errorf("draw %s: %w", path, err)
	}

	return "wrote " + path, write(path, svg)
}

// writeICO writes ic as a favicon.ico of icoSizes to path and returns the line
// reporting it.
func writeICO(ic *icon.Icon, path string) (string, error) {

	pngs := make([][]byte, 0, len(icoSizes()))
	dims := make([]string, 0, len(icoSizes()))
	for _, size := range icoSizes() {
		data, err := encodePNG(ic, size)
		if err != nil {
			return "", fmt.Errorf("draw %s: %w", path, err)
		}

		pngs = append(pngs, data)
		dims = append(dims, fmt.Sprintf("%d×%d", size, size))
	}

	ico, err := icon.EncodeICO(pngs...)
	if err != nil {
		return "", fmt.Errorf("pack %s: %w", path, err)
	}

	return fmt.Sprintf("wrote %s (%s)", path, strings.Join(dims, ", ")), write(path, ico)
}

// writePNG writes ic size pixels square to path and returns the line reporting it.
func writePNG(ic *icon.Icon, size int, path string) (string, error) {

	data, err := encodePNG(ic, size)
	if err != nil {
		return "", fmt.Errorf("draw %s: %w", path, err)
	}

	return fmt.Sprintf("wrote %s (%d×%d)", path, size, size), write(path, data)
}

// encodePNG draws ic size pixels square in the light palette, as a PNG.
func encodePNG(ic *icon.Icon, size int) ([]byte, error) {

	img, err := ic.Render(size, icon.Light())
	if err != nil {
		return nil, fmt.Errorf("render: %w", err)
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	return buf.Bytes(), nil
}

// write writes data to path.
func write(path string, data []byte) error {

	err := os.WriteFile(path, data, filePerm)
	if err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	return nil
}
