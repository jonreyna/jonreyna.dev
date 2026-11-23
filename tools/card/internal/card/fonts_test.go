package card

import (
	"strings"
	"testing"

	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/goregular"
)

func TestParseFonts(t *testing.T) {

	t.Parallel()
	notAFont := []byte("not a font")
	tests := []struct {
		name      string
		smallCaps []byte
		italic    []byte
		wantErr   string // empty for success
	}{
		{"both fonts", goregular.TTF, goitalic.TTF, ""},
		{"a bad small-capitals font", notAFont, goitalic.TTF, "parse small-capitals font"},
		{"a bad italic font", goregular.TTF, notAFont, "parse italic font"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			fonts, err := ParseFonts(tc.smallCaps, tc.italic)
			switch {
			case tc.wantErr == "" && (err != nil || fonts.SmallCaps == nil || fonts.Italic == nil):
				t.Errorf("ParseFonts = %v, %v; want both fonts", fonts, err)
			case tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)):
				t.Errorf("ParseFonts error %v, want one containing %q", err, tc.wantErr)
			}
		})
	}
}
