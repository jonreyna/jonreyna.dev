package plate

import "testing"

func TestPersonDates(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name   string
		person Person
		want   string
	}{
		{name: "BC", person: Person{From: -428, To: -348, Circa: true}, want: "c. 428–348 BC"},
		{name: "across the era", person: Person{From: -20, To: 50, Circa: true}, want: "c. 20 BC – AD 50"},
		{name: "AD", person: Person{From: 354, To: 430}, want: "354–430"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			if got := tc.person.Dates(); got != tc.want {
				t.Fatalf("Dates() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestPersonAnchor(t *testing.T) {

	t.Parallel()

	span := [2]int{-800, 450}
	tests := []struct {
		name   string
		person Person
		want   float64
	}{
		{name: "middle of the years", person: Person{From: 100, To: 200}, want: 150},
		{name: "clipped to the span", person: Person{From: -1000, To: -600}, want: -700},
		{name: "placed by At", person: Person{From: -1000, To: -600, At: -750}, want: -750},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			if got := tc.person.Anchor(span); got != tc.want {
				t.Fatalf("Anchor() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestLinkType(t *testing.T) {

	t.Parallel()

	tests := []struct {
		typ      LinkType
		valid    bool
		strength int
		verb     string
	}{
		{LinkQuotes, true, 3, "quotes"},
		{LinkRejects, true, 3, "rejects"},
		{LinkAdapts, true, 2, "adapts"},
		{LinkArgued, true, 1, "may draw on"},
		{"cites", false, 0, "cites"},
	}

	for _, tc := range tests {
		t.Run(string(tc.typ), func(t *testing.T) {

			t.Parallel()

			if tc.typ.Valid() != tc.valid || tc.typ.Strength() != tc.strength || tc.typ.Verb() != tc.verb {
				t.Fatalf("%q: Valid %v, Strength %d, Verb %q; want %v, %d, %q",
					tc.typ, tc.typ.Valid(), tc.typ.Strength(), tc.typ.Verb(), tc.valid, tc.strength, tc.verb)
			}
		})
	}
}
