package layout

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"
)

// minimal returns the shared test plate, parsed.
func minimal(t *testing.T) *plate.Plate {

	t.Helper()

	data, err := os.ReadFile("../../testdata/minimal.yaml")
	if err != nil {
		t.Fatal(err)
	}

	p, err := plate.Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	return p
}

func TestNew(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name    string
		options func(o *Options)
		wantErr string
	}{
		{name: "defaults", options: func(*Options) {}},
		{name: "no annealing", options: func(o *Options) { o.Iterations = 0 }},
		{name: "too narrow for the labels", options: func(o *Options) { o.Width = 200 }, wantErr: "knots need a width"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			o := Defaults()
			tc.options(&o)

			l, err := New(minimal(t), &o)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("New() error = %v, want %q", err, tc.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatal(err)
			}

			checkLayout(t, l)
		})
	}
}

// checkLayout checks what any layout of the minimal plate must satisfy.
func checkLayout(t *testing.T, l *Layout) {

	t.Helper()

	if len(l.Bands) != 2 || len(l.Knots) != 4 {
		t.Fatalf("got %d bands and %d knots, want 2 and 4", len(l.Bands), len(l.Knots))
	}

	for i := 1; i < len(l.Knots); i++ {
		if l.Knots[i].X <= l.Knots[i-1].X {
			t.Errorf("knot %s at x %.1f is not right of %s at %.1f, though later",
				l.Knots[i].Person, l.Knots[i].X, l.Knots[i-1].Person, l.Knots[i-1].X)
		}
	}

	ben := l.Knots[1]
	if ben.Top >= l.Bands[0].Bottom || ben.Bottom <= l.Bands[1].Top {
		t.Errorf("ben spans %.1f–%.1f, want it across both bands", ben.Top, ben.Bottom)
	}

	if n := len(l.Knots[3].Ends); n != 2 {
		t.Errorf("dee has %d end marks, want 2", n)
	}
}

// TestNewDeterministic checks that the annealing is reproducible: one scenario, run
// twice, so a table would only obscure it.
func TestNewDeterministic(t *testing.T) {

	t.Parallel()

	o := Defaults()

	a, err := New(minimal(t), &o)
	if err != nil {
		t.Fatal(err)
	}

	b, err := New(minimal(t), &o)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(a, b) {
		t.Fatal("two layouts of the same plate differ")
	}
}

func TestTimeScale(t *testing.T) {

	t.Parallel()

	s := &timeScale{}
	s.add(-500, 100)
	s.add(0, 200)
	s.add(0, 999) // a second point at the same year is skipped
	s.add(400, 600)

	tests := []struct {
		year, want float64
	}{
		{-600, 100}, {-500, 100}, {-250, 150}, {0, 200}, {200, 400}, {400, 600}, {500, 600},
	}

	for _, tc := range tests {
		if got := s.x(tc.year); got != tc.want {
			t.Errorf("x(%v) = %v, want %v", tc.year, got, tc.want)
		}
	}
}

func TestCenturyTicks(t *testing.T) {

	t.Parallel()

	s := &timeScale{}
	s.add(-200, 0)
	s.add(200, 400)

	ticks := centuryTicks([2]int{-200, 200}, s)
	got := make([]string, 0, len(ticks))

	for _, tk := range ticks {
		got = append(got, tk.Label)
	}

	want := []string{"200 BC", "100 BC", "AD 1", "AD 100", "AD 200"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("labels = %q, want %q", got, want)
	}
}
