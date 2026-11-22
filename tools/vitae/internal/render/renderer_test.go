package render

import (
	"errors"
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

func TestFor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		format  string
		wantErr error
	}{
		{format: "json"},
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

func TestRendererFunc(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		f    RendererFunc
		want string
	}{
		{name: "passes the résumé through", f: func(r *variant.Resolved) ([]byte, error) { return []byte(r.Variant), nil }, want: "full"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := tc.f.Render(&variant.Resolved{Variant: "full"})
			if err != nil || string(got) != tc.want {
				t.Errorf("Render = %q, %v; want %q", got, err, tc.want)
			}
		})
	}
}
