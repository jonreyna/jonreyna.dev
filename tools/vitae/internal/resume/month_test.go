package resume

import (
	"errors"
	"testing"
)

func TestParseMonth(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name    string
		in      string
		want    Month
		wantErr bool
	}{
		{name: "valid", in: "2023-10", want: Month{Year: 2023, Month: 10}},
		{name: "january", in: "2015-01", want: Month{Year: 2015, Month: 1}},
		{name: "month 13", in: "2023-13", wantErr: true},
		{name: "month 00", in: "2023-00", wantErr: true},
		{name: "full date", in: "2023-10-01", wantErr: true},
		{name: "slash", in: "2023/10", wantErr: true},
		{name: "letters", in: "20x3-10", wantErr: true},
		{name: "empty", in: "", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			got, err := ParseMonth(tc.in)
			if tc.wantErr {
				if !errors.Is(err, ErrInvalidMonth) {
					t.Fatalf("ParseMonth(%q) error = %v, want ErrInvalidMonth", tc.in, err)
				}

				return
			}

			if err != nil || got != tc.want {
				t.Fatalf("ParseMonth(%q) = %v, %v; want %v", tc.in, got, err, tc.want)
			}

			if got.String() != tc.in {
				t.Errorf("String() = %q, want %q", got.String(), tc.in)
			}
		})
	}
}

func TestMonthBefore(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name string
		a, b Month
		want bool
	}{
		{name: "earlier year", a: Month{2019, 12}, b: Month{2020, 1}, want: true},
		{name: "earlier month", a: Month{2020, 3}, b: Month{2020, 4}, want: true},
		{name: "same month", a: Month{2020, 4}, b: Month{2020, 4}, want: false},
		{name: "later", a: Month{2021, 1}, b: Month{2020, 12}, want: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			if got := tc.a.Before(tc.b); got != tc.want {
				t.Errorf("%v.Before(%v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestMonthJSON(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name string
		in   Month
		want string
	}{
		{name: "padded", in: Month{2015, 4}, want: `"2015-04"`},
		{name: "two digits", in: Month{2023, 10}, want: `"2023-10"`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			got, err := tc.in.MarshalJSON()
			if err != nil || string(got) != tc.want {
				t.Errorf("MarshalJSON() = %s, %v; want %s", got, err, tc.want)
			}
		})
	}
}
