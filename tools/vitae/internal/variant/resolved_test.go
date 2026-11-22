package variant

import "testing"

func TestResolvedHas(t *testing.T) {
	t.Parallel()
	r := &Resolved{Sections: []string{"summary", "experience"}}
	tests := []struct {
		section string
		want    bool
	}{
		{"summary", true}, {"experience", true}, {"skills", false}, {"", false},
	}
	for _, tc := range tests {
		t.Run(tc.section, func(t *testing.T) {
			t.Parallel()
			if got := r.Has(tc.section); got != tc.want {
				t.Errorf("Has(%q) = %v, want %v", tc.section, got, tc.want)
			}
		})
	}
}
