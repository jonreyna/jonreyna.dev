package markup

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name, in, plain, html string
	}{
		{name: "plain text", in: "plain text", plain: "plain text", html: "plain text"},
		{name: "emphasis", in: "a *b* c", plain: "a b c", html: "a <em>b</em> c"},
		{name: "code", in: "run `go test`", plain: "run go test", html: "run <code>go test</code>"},
		{name: "link and escaping", in: "[site](https://x.dev) & <b>", plain: "site & <b>", html: `<a href="https://x.dev">site</a> &amp; &lt;b&gt;`},
		{name: "mailto link", in: "[mail](mailto:a@b.c)", plain: "mail", html: `<a href="mailto:a@b.c">mail</a>`},
		{name: "tie", in: "RFC~8342", plain: "RFC 8342", html: "RFC 8342"},
		{name: "escapes", in: `\*not\* \~ \[x]`, plain: "*not* ~ [x]", html: "*not* ~ [x]"},
		{name: "tie in emphasis", in: "*a~b*", plain: "a b", html: "<em>a b</em>"},
		{name: "escaped star in emphasis", in: `*a\*b*`, plain: "a*b", html: "<em>a*b</em>"},
		{name: "tilde in code stays", in: "`~/bin`", plain: "~/bin", html: "<code>~/bin</code>"},
		{name: "empty", in: "", plain: "", html: ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			spans, err := Parse(tc.in)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tc.in, err)
			}

			if got := Plain(spans); got != tc.plain {
				t.Errorf("Plain = %q, want %q", got, tc.plain)
			}

			if got := HTML(spans); got != tc.html {
				t.Errorf("HTML = %q, want %q", got, tc.html)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name string
		in   string
		want error
	}{
		{name: "unclosed emphasis", in: "unclosed *emphasis", want: ErrUnclosed},
		{name: "unclosed code", in: "unclosed `code", want: ErrUnclosed},
		{name: "empty emphasis", in: "empty ** here", want: ErrEmpty},
		{name: "bracket without url", in: "[no url] here", want: ErrBadLink},
		{name: "http url", in: "[bad](http://insecure.example)", want: ErrBadLink},
		{name: "unclosed url", in: "[open](https://x.dev", want: ErrBadLink},
		{name: "trailing backslash", in: `trailing \`, want: ErrTrailingBackslash},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			if _, err := Parse(tc.in); !errors.Is(err, tc.want) {
				t.Errorf("Parse(%q) error = %v, want %v", tc.in, err, tc.want)
			}
		})
	}
}

func TestKindString(t *testing.T) {

	t.Parallel()
	tests := []struct {
		kind Kind
		want string
	}{
		{Text, "text"}, {Emph, "emph"}, {Code, "code"}, {Link, "link"}, {Kind(99), "unknown"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {

			t.Parallel()
			if got := tc.kind.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestMustParsePanics(t *testing.T) {

	t.Parallel()
	tests := []struct {
		name      string
		in        string
		wantPanic bool
	}{
		{name: "valid", in: "*ok*"},
		{name: "invalid", in: "*open", wantPanic: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()
			defer func() {
				if got := recover() != nil; got != tc.wantPanic {
					t.Errorf("panicked = %v, want %v", got, tc.wantPanic)
				}
			}()

			MustParse(tc.in)
		})
	}
}
