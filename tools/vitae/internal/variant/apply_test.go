package variant

import (
	"errors"
	"slices"
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

func hl(text string, weight int, tags ...string) resume.Highlight {
	return resume.Highlight{Text: text, Weight: weight, Tags: tags}
}

func month(y, m int) *resume.Month { return &resume.Month{Year: y, Month: m} }

func TestHighlights(t *testing.T) {
	t.Parallel()
	all := []resume.Highlight{hl("a", 50, "go"), hl("b", 90), hl("c", 60, "ai"), hl("d", 95, "old")}
	tests := []struct {
		name string
		v    resume.Variant
		want []string
	}{
		{name: "everything, in order", v: resume.Variant{}, want: []string{"a", "b", "c", "d"}},
		{name: "exclude", v: resume.Variant{Exclude: []string{"old"}}, want: []string{"a", "b", "c"}},
		{name: "min weight", v: resume.Variant{MinWeight: 60}, want: []string{"b", "c", "d"}},
		{name: "top two keep file order", v: resume.Variant{MaxHighlights: 2}, want: []string{"b", "d"}},
		{name: "boost lifts a highlight", v: resume.Variant{MaxHighlights: 2, Boost: []string{"ai"}, BoostBy: 40}, want: []string{"c", "d"}},
		{name: "default boost is 20", v: resume.Variant{MaxHighlights: 1, Boost: []string{"go"}, Exclude: []string{"old"}}, want: []string{"b"}},
		// a scores 50+40 = 90, tying b; the earlier one wins.
		{name: "ties go to the earlier one", v: resume.Variant{MaxHighlights: 1, Boost: []string{"go"}, BoostBy: 40, Exclude: []string{"old"}}, want: []string{"a"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var got []string
			for _, h := range highlights(all, &tc.v) {
				got = append(got, h.Text)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRoles(t *testing.T) {
	t.Parallel()
	all := []resume.Role{
		{ID: "now", Start: *month(2024, 1)},
		{ID: "mid", Start: *month(2015, 4), End: month(2023, 12)},
		{ID: "old", Start: *month(2013, 9), End: month(2015, 1)},
	}
	tests := []struct {
		name string
		v    resume.Variant
		want []string
	}{
		{name: "no limits", v: resume.Variant{}, want: []string{"now", "mid", "old"}},
		{name: "since drops roles that ended before", v: resume.Variant{Since: month(2015, 4)}, want: []string{"now", "mid"}},
		{name: "a role ending in the since month stays", v: resume.Variant{Since: month(2015, 1)}, want: []string{"now", "mid", "old"}},
		{name: "max roles keeps the newest", v: resume.Variant{MaxRoles: 1}, want: []string{"now"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var got []string
			for _, r := range roles(all, &tc.v) {
				got = append(got, r.ID)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCourses(t *testing.T) {
	t.Parallel()
	all := []resume.Course{{Name: "linear algebra", Tags: []string{"math"}}, {Name: "pottery"}}
	tests := []struct {
		name string
		show []string
		want []string
	}{
		{name: "empty shows all", show: nil, want: []string{"linear algebra", "pottery"}},
		{name: "by tag", show: []string{"math"}, want: []string{"linear algebra"}},
		{name: "no match", show: []string{"art"}, want: nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var got []string
			for _, c := range courses(all, tc.show) {
				got = append(got, c.Name)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestApply(t *testing.T) {
	t.Parallel()
	r := &resume.Resume{
		Basics:  resume.Basics{Name: "Ada", Label: "Engineer"},
		Summary: map[string]string{"default": "long", "short": "short"},
		Earlier: "odd jobs",
		Variants: map[string]resume.Variant{
			"full":  {Sections: []string{resume.SectionSummary, resume.SectionEarlier}},
			"brief": {Sections: []string{resume.SectionSummary}, Summary: "short", Label: "Lead"},
		},
	}
	tests := []struct {
		name                   string
		variant                string
		wantSummary, wantLabel string
		wantEarlier            string
		wantErr                error
	}{
		{name: "default summary", variant: "full", wantSummary: "long", wantLabel: "Engineer", wantEarlier: "odd jobs"},
		{name: "chosen summary, label and sections", variant: "brief", wantSummary: "short", wantLabel: "Lead"},
		{name: "unknown variant", variant: "nope", wantErr: ErrUnknown},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := Apply(r, tc.variant)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got.Summary != tc.wantSummary || got.Basics.Label != tc.wantLabel || got.Earlier != tc.wantEarlier {
				t.Errorf("got summary %q, label %q, earlier %q", got.Summary, got.Basics.Label, got.Earlier)
			}
			if r.Basics.Label != "Engineer" {
				t.Errorf("Apply changed the résumé's label to %q", r.Basics.Label)
			}
		})
	}
}
