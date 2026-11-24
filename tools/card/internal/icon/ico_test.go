package icon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/png"
	"testing"
)

// newTestPNG encodes a blank w×h image as a PNG.
func newTestPNG(t *testing.T, w, h int) []byte {

	t.Helper()
	var buf bytes.Buffer
	err := png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, w, h)))
	if err != nil {
		t.Fatalf("encode PNG: %v", err)
	}

	return buf.Bytes()
}

func TestEncodeICO(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name  string
		sides []int
		want  []uint8 // each entry's width byte
	}{
		{"one image", []int{32}, []uint8{32}},
		{"two images in order", []int{16, 32}, []uint8{16, 32}},
		{"256 is stored as 0", []int{256}, []uint8{0}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			pngs := make([][]byte, len(tc.sides))
			for i, s := range tc.sides {
				pngs[i] = newTestPNG(t, s, s)
			}

			ico, err := EncodeICO(pngs...)
			if err != nil {
				t.Fatalf("EncodeICO: %v", err)
			}

			var header [3]uint16
			r := bytes.NewReader(ico)
			if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
				t.Fatalf("read header: %v", err)
			}

			if header != [3]uint16{0, icoTypeIcon, uint16(len(tc.sides))} { //nolint:gosec // G115: a few test images
				t.Fatalf("header %v, want an icon of %d images", header, len(tc.sides))
			}

			entries := make([]icoEntry, len(tc.sides))
			if err := binary.Read(r, binary.LittleEndian, entries); err != nil {
				t.Fatalf("read entries: %v", err)
			}

			for i, e := range entries {
				if e.Width != tc.want[i] || e.Height != tc.want[i] {
					t.Errorf("entry %d is %d×%d, want %d", i, e.Width, e.Height, tc.want[i])
				}

				img := ico[e.Offset : e.Offset+e.Size]
				if !bytes.Equal(img, pngs[i]) {
					t.Errorf("entry %d doesn't point at image %d's PNG", i, i)
				}
			}
		})
	}
}

func TestEncodeICOErrors(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name string
		pngs [][]byte
	}{
		{"no images", nil},
		{"not a PNG", [][]byte{[]byte("GIF89a")}},
		{"not square", [][]byte{newTestPNG(t, 16, 32)}},
		{"too big", [][]byte{newTestPNG(t, 257, 257)}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			if _, err := EncodeICO(tc.pngs...); !errors.Is(err, ErrICO) {
				t.Errorf("EncodeICO error %v, want ErrICO", err)
			}
		})
	}
}
