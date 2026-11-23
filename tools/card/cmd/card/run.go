package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image/png"
	"io"
	"os"

	"github.com/jonreyna/jonreyna.dev/tools/card/internal/card"
)

// Exit codes (see doc.go).
const (
	exitFailure = 1
	exitUsage   = 2
)

// filePerm is the PNG's mode: a public image for the website, readable by all.
const filePerm = 0o644

// errUsage marks errors in the command line, which exit with exitUsage.
var errUsage = errors.New("usage")

// run parses args, renders the card and writes the PNG, reporting on stdout.
func run(args []string, stdout, stderr io.Writer) error {

	fs := flag.NewFlagSet("card", flag.ContinueOnError)
	fs.SetOutput(stderr)
	smallCaps := fs.String("smallcaps", "", "small-capitals `font` for the name and address")
	italic := fs.String("italic", "", "italic `font` for the tagline")
	c := card.Card{}
	fs.StringVar(&c.Name, "name", "", "the name, the largest line")
	fs.StringVar(&c.Motto, "motto", "", "the tagline under the name")
	fs.StringVar(&c.Site, "site", "", "the address on the last line")
	out := fs.String("o", "", "the PNG `file` to write")

	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("%w: %w", errUsage, err)
	}

	for _, f := range []struct{ name, value string }{
		{"smallcaps", *smallCaps}, {"italic", *italic}, {"name", c.Name},
		{"motto", c.Motto}, {"site", c.Site}, {"o", *out},
	} {
		if f.value == "" {
			return fmt.Errorf("%w: -%s is required", errUsage, f.name)
		}
	}

	fonts, err := readFonts(*smallCaps, *italic)
	if err != nil {
		return err
	}

	img, err := c.Render(fonts)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return fmt.Errorf("encode %s: %w", *out, err)
	}

	err = os.WriteFile(*out, buf.Bytes(), filePerm)
	if err != nil {
		return fmt.Errorf("write %s: %w", *out, err)
	}

	_, err = fmt.Fprintf(stdout, "wrote %s (%d×%d)\n", *out, card.Width, card.Height)
	if err != nil {
		return fmt.Errorf("report: %w", err)
	}

	return nil
}

// readFonts reads and parses the two font files.
func readFonts(smallCapsPath, italicPath string) (*card.Fonts, error) {

	//nolint:gosec // G304: the font paths are the command's own arguments
	smallCaps, err := os.ReadFile(smallCapsPath)
	if err != nil {
		return nil, fmt.Errorf("read small-capitals font: %w", err)
	}

	//nolint:gosec // G304: the font paths are the command's own arguments
	italic, err := os.ReadFile(italicPath)
	if err != nil {
		return nil, fmt.Errorf("read italic font: %w", err)
	}

	fonts, err := card.ParseFonts(smallCaps, italic)
	if err != nil {
		return nil, fmt.Errorf("%s, %s: %w", smallCapsPath, italicPath, err)
	}

	return fonts, nil
}
