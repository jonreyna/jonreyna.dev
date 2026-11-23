package markup_test

import (
	"fmt"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/markup"
)

func ExampleParse() {

	spans, err := markup.Parse("Built *küber* in `Go`; see [the site](https://jonreyna.dev).")
	if err != nil {
		panic(err)
	}

	for _, s := range spans {
		fmt.Printf("%s %q\n", s.Kind, s.Text)
	}

	// Output:
	// text "Built "
	// emph "küber"
	// text " in "
	// code "Go"
	// text "; see "
	// link "the site"
	// text "."
}

func ExamplePlain() {
	fmt.Println(markup.Plain(markup.MustParse("drew on RFC~8342 and *etcd*")))
	// Output: drew on RFC 8342 and etcd
}

func ExampleHTML() {
	fmt.Println(markup.HTML(markup.MustParse("Go & C, [GitHub](https://github.com/jonreyna)")))
	// Output: Go &amp; C, <a href="https://github.com/jonreyna">GitHub</a>
}

func ExampleLaTeX() {
	fmt.Println(markup.LaTeX(markup.MustParse("99.8% of *15,000* a week; RFC~8342, [GitHub](https://github.com/jonreyna)")))
	// Output: 99.8\% of \emph{15,000} a week; RFC~8342, \href{https://github.com/jonreyna}{GitHub}
}

func ExampleMarkdown() {
	fmt.Println(markup.Markdown(markup.MustParse("*küber* in `Go`, $5 [site](https://jonreyna.dev)")))
	// Output: *küber* in `Go`, \$5 [site](https://jonreyna.dev)
}
