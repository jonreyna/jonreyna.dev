package render

import (
	"errors"
	"testing"
)

func TestFor(t *testing.T) {

	t.Parallel()
	tests := []struct {
		format  string
		wantErr error
	}{
		{format: "json"},
		{format: "tex"},
		{format: "txt"},
		{format: "md"},
		{format: "pdf", wantErr: ErrUnknownFormat},
		{format: "", wantErr: ErrUnknownFormat},
	}
	for _, tc := range tests {
		t.Run(tc.format, func(t *testing.T) {

			t.Parallel()
			r, err := For(tc.format)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("For(%q) error = %v, want %v", tc.format, err, tc.wantErr)
			}

			if tc.wantErr == nil && r == nil {
				t.Fatalf("For(%q) returned no renderer", tc.format)
			}
		})
	}
}
