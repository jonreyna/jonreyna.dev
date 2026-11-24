package icon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image/png"
)

// The ICO format's fixed sizes, in bytes, and its largest image side.
const (
	icoHeaderSize = 6
	icoEntrySize  = 16
	icoMaxSide    = 256
	icoTypeIcon   = 1 // the header's type field for an icon (2 is a cursor)
	icoBitCount   = 32
)

// ErrICO reports images EncodeICO can't pack.
var ErrICO = errors.New("can't pack into an ICO")

// icoEntry is one image's directory entry in an ICO file.
type icoEntry struct {
	Width, Height uint8 // 0 means 256
	Colors        uint8 // 0: not a palette image
	Reserved      uint8
	Planes        uint16
	BitCount      uint16
	Size          uint32 // the image's length in bytes
	Offset        uint32 // where the image starts in the file
}

// EncodeICO packs PNG images into one ICO file, in the order given. Each must be a
// square from 1 to 256 pixels. Browsers read PNGs in ICO files, so the images are
// stored as they are.
func EncodeICO(pngs ...[]byte) ([]byte, error) {

	if len(pngs) == 0 {
		return nil, fmt.Errorf("%w: no images", ErrICO)
	}

	var out bytes.Buffer
	header := []uint16{0, icoTypeIcon, uint16(len(pngs))} //nolint:gosec // G115: at most a few images
	offset := icoHeaderSize + icoEntrySize*len(pngs)
	entries := make([]icoEntry, len(pngs))
	for i, data := range pngs {
		cfg, err := png.DecodeConfig(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("%w: image %d: %w", ErrICO, i+1, err)
		}

		if cfg.Width != cfg.Height || cfg.Width < 1 || cfg.Width > icoMaxSide {
			return nil, fmt.Errorf("%w: image %d is %d×%d, not a square of 1 to %d pixels",
				ErrICO, i+1, cfg.Width, cfg.Height, icoMaxSide)
		}

		side := uint8(cfg.Width & byteMask) // 256 wraps to 0, which the format reads as 256
		entries[i] = icoEntry{
			Width: side, Height: side, Planes: 1, BitCount: icoBitCount,
			Size: uint32(len(data)), Offset: uint32(offset), //nolint:gosec // G115: icons are kilobytes
		}
		offset += len(data)
	}

	for _, v := range []any{header, entries} {
		err := binary.Write(&out, binary.LittleEndian, v)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrICO, err)
		}
	}

	for _, data := range pngs {
		out.Write(data)
	}

	return out.Bytes(), nil
}
