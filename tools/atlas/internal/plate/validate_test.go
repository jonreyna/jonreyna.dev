package plate

import (
	"os"
	"strings"
	"testing"
)

// minimal returns the shared test plate, parsed.
func minimal(t *testing.T) *Plate {

	t.Helper()

	data, err := os.ReadFile("../../testdata/minimal.yaml")
	if err != nil {
		t.Fatal(err)
	}

	p, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	return p
}

func TestValidate(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name           string
		change         func(p *Plate)
		requireChecked bool
		want           string // a substring of one issue; empty means no issues
	}{
		{name: "minimal plate passes", change: func(*Plate) {}},
		{name: "unchecked link fails a production build", change: func(*Plate) {}, requireChecked: true,
			want: "link L4: not checked"},
		{name: "unknown thread", change: func(p *Plate) { p.Links[0].Thread = "nope" }, want: `unknown thread "nope"`},
		{name: "unknown person", change: func(p *Plate) { p.Links[0].To = "zed" }, want: "unknown person"},
		{name: "link backward in time", change: func(p *Plate) { p.Links[1].From, p.Links[1].To = "dee", "ana" },
			want: "runs backward in time"},
		{name: "argued without scholarship", change: func(p *Plate) { p.Links[3].Scholarship = nil },
			want: "must name its scholarship"},
		{name: "passage without url", change: func(p *Plate) { p.Links[0].URL = "" }, want: "needs quote, source and url"},
		{name: "unsafe url", change: func(p *Plate) { p.Links[0].URL = "javascript:alert(1)" },
			want: "must be an http or https address"},
		{name: "unknown source", change: func(p *Plate) { p.Links[0].Source = "nope" }, want: `unknown source "nope"`},
		{name: "quote too long", change: func(p *Plate) { p.Links[0].Quote = strings.Repeat("word ", 61) },
			want: "quote has 61 words"},
		{name: "lanes skipped", change: func(p *Plate) { p.People[1].Lanes = []string{"new", "old"} },
			want: "top to bottom"},
		{name: "year zero", change: func(p *Plate) { p.People[0].To = 0 }, want: "no year 0"},
		{name: "unused thread", change: func(p *Plate) { p.Threads = append(p.Threads, Thread{ID: "x", Name: "X", Color: "c"}) },
			want: "thread x: no link uses it"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			p := minimal(t)
			tc.change(p)

			issues := Validate(p, tc.requireChecked)
			if tc.want == "" {
				if len(issues) > 0 {
					t.Fatalf("Validate() = %v, want no issues", issues)
				}

				return
			}

			for _, i := range issues {
				if strings.Contains(i.String(), tc.want) {
					return
				}
			}

			t.Fatalf("Validate() = %v, want an issue containing %q", issues, tc.want)
		})
	}
}

func TestParse(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name string
		yaml string
		want string
	}{
		{name: "empty", yaml: "", want: "empty plate"},
		{name: "unknown field", yaml: "plate: x\nshade: red\n", want: "field shade not found"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			_, err := Parse([]byte(tc.yaml))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("Parse() error = %v, want it to contain %q", err, tc.want)
			}
		})
	}
}
