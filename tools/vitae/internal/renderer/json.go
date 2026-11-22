package renderer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// JSON renders r as its [Document]: two-space indentation, fields in a fixed order and
// no HTML escaping, so the output is stable byte for byte.
func JSON(r *variant.Resolved) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(NewDocument(r)); err != nil {
		return nil, fmt.Errorf("encode json: %w", err)
	}
	return buf.Bytes(), nil
}
