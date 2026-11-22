package filter

import (
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/model"
)

func hl(text string, weight int, tags ...string) model.Highlight {
	return model.Highlight{Text: text, Weight: weight, Tags: tags}
}

func texts(hs []model.Highlight) []string {
	var out []string
	for _, h := range hs {
		out = append(out, h.Text)
	}
	return out
}

func TestHighlights(t *testing.T) {
	all := []model.Highlight{hl("a", 50, "go"), hl("b", 90), hl("c", 60, "ai"), hl("d", 95, "old")}
	tests := []struct {
		name string
		v    model.Variant
		want []string
	}{
		{"everything, in order", model.Variant{}, []string{"a", "b", "c", "d"}},
		{"exclude", model.Variant{Exclude: []string{"old"}}, []string{"a", "b", "c"}},
		{"min weight", model.Variant{MinWeight: 60}, []string{"b", "c", "d"}},
		{"top two keep file order", model.Variant{MaxHighlights: 2}, []string{"b", "d"}},
		{"boost lifts a highlight", model.Variant{MaxHighlights: 2, Boost: []string{"ai"}, BoostBy: 40}, []string{"c", "d"}},
		{"default boost is 20", model.Variant{MaxHighlights: 1, Boost: []string{"go"}, Exclude: []string{"old"}}, []string{"b"}},
		// a scores 50+40 = 90, tying b; the earlier one wins.
		{"ties go to the earlier one", model.Variant{MaxHighlights: 1, Boost: []string{"go"}, BoostBy: 40, Exclude: []string{"old"}}, []string{"a"}},
	}
	for _, tt := range tests {
		got := texts(highlights(all, tt.v))
		if len(got) != len(tt.want) {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
				break
			}
		}
	}
}

func TestSinceAndMaxRoles(t *testing.T) {
	m := func(y, mo int) *model.Month { return &model.Month{Year: y, Month: mo} }
	roles_ := []model.Role{
		{ID: "now", Start: *m(2024, 1)},
		{ID: "mid", Start: *m(2015, 4), End: m(2023, 12)},
		{ID: "old", Start: *m(2013, 9), End: m(2015, 1)},
	}
	ids := func(rs []model.Role) (out []string) {
		for _, r := range rs {
			out = append(out, r.ID)
		}
		return
	}
	if got := ids(roles(roles_, model.Variant{Since: m(2015, 4)})); len(got) != 2 || got[1] != "mid" {
		t.Errorf("since: got %v", got)
	}
	if got := ids(roles(roles_, model.Variant{Since: m(2015, 1)})); len(got) != 3 {
		t.Errorf("a role ending in the since month stays: got %v", got)
	}
	if got := ids(roles(roles_, model.Variant{MaxRoles: 1})); len(got) != 1 || got[0] != "now" {
		t.Errorf("max roles: got %v", got)
	}
}
