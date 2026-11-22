package variant_test

import (
	"fmt"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// A variant keeps the best-scoring highlights and leaves them in their original order.
func ExampleApply() {

	r := &resume.Resume{
		Experience: []resume.Role{{
			ID: "acme",
			Highlights: []resume.Highlight{
				{Text: "Led the Go rewrite", Weight: 70, Tags: []string{"go"}},
				{Text: "Ran the on-call rotation", Weight: 80},
				{Text: "Built the ML feature store", Weight: 60, Tags: []string{"ai"}},
			},
		}},
		Variants: map[string]resume.Variant{
			// "ai" adds 20 to its highlight's weight: 60 becomes 80, tying the on-call
			// highlight; the tie goes to the earlier one, and both beat 70.
			"ai-infra": {Sections: []string{resume.SectionExperience}, Boost: []string{"ai"}, MaxHighlights: 2},
		},
	}
	out, err := variant.Apply(r, "ai-infra")
	if err != nil {
		panic(err)
	}

	for _, h := range out.Experience[0].Highlights {
		fmt.Println(h.Text)
	}

	// Output:
	// Ran the on-call rotation
	// Built the ML feature store
}
