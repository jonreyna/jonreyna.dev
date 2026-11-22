package render

import (
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

func TestFunc(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name string
		f    Func
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
