package inline

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		in, plain, html string
	}{
		{"plain text", "plain text", "plain text"},
		{"a *b* c", "a b c", "a <em>b</em> c"},
		{"run `go test`", "run go test", "run <code>go test</code>"},
		{"[site](https://x.dev) & <b>", "site & <b>", `<a href="https://x.dev">site</a> &amp; &lt;b&gt;`},
		{"RFC~8342", "RFC 8342", "RFC 8342"},
		{`\*not\* \~ \[x]`, "*not* ~ [x]", "*not* ~ [x]"},
		{"*a~b*", "a b", "<em>a b</em>"},
	}
	for _, tt := range tests {
		spans, err := Parse(tt.in)
		if err != nil {
			t.Errorf("Parse(%q): %v", tt.in, err)
			continue
		}
		if got := Plain(spans); got != tt.plain {
			t.Errorf("Plain(%q) = %q, want %q", tt.in, got, tt.plain)
		}
		if got := HTML(spans); got != tt.html {
			t.Errorf("HTML(%q) = %q, want %q", tt.in, got, tt.html)
		}
	}
}

func TestParseErrors(t *testing.T) {
	for _, in := range []string{
		"unclosed *emphasis",
		"unclosed `code",
		"empty ** here",
		"[no url] here",
		"[bad](http://insecure.example)",
		"[open](https://x.dev",
		`trailing \`,
	} {
		if _, err := Parse(in); err == nil {
			t.Errorf("Parse(%q) succeeded, want an error", in)
		}
	}
}
