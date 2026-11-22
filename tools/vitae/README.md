# vitae

Turns `resume/resume.yaml` into every résumé format and variant. It's the only code that understands variants: the site reads its JSON output instead of filtering again.

```
vitae validate [-i FILE] [-private FILE] [-strict]
vitae render   [-i FILE] [-private FILE] [-strict] [-variant NAME|all] [-format NAME|all] [-o DIR|-]
vitae list     [-i FILE] variants|tags
vitae version
```

Exit codes: 0 success, 1 validation failure, 2 usage error. From the repository root, `make resume-check` tests and validates, and `make resume` renders into `build/resume/`.

## How it works

Each package has a `doc.go`; `go doc ./internal/<package>` is the reference. In brief:

- **`internal/resume`** holds the YAML as Go types, one file per type. `Load` rejects unknown fields, so a typo fails loudly.
- **`internal/lint`** checks structure and content. Errors always fail; warnings fail with `-strict`, which CI uses. The rules include:
  - roles newest first, and only the first one current;
  - every tag declared and used;
  - highlights under 650 characters;
  - curly quotes and real dashes;
  - no phone number in the public file;
  - an education entry must say whether the degree was conferred, and if it wasn't, only its note may mention the degree.
- **`internal/variant`** applies a variant. Each highlight scores its weight plus `boost_by` (default 20) per boosted tag. Excluded tags and scores under `min_weight` drop out, and the best `max_highlights` survive in their original order.
- **`internal/markup`** handles the markup allowed in text fields: `*emphasis*`, `` `code` ``, `[text](https://url)`, and `~` for a non-breaking space.
- **`internal/render`** writes the outputs. Each format is a `Renderer` (one method), looked up with `render.NewRenderer`. JSON so far. LaTeX, plain text and Markdown come next, then PDF (latexmk) and Word (Pandoc) from `make`.

Output is deterministic, so `testdata/golden/` holds exact expected files for `testdata/resume.yaml`. After an intended change, run `go test ./... -update` and review the diff.

## Code rules

Each public type has its own file, and every package, type and field has a doc comment. Tests are table-driven, and examples run as tests. golangci-lint v2.15.0 runs with the strict config in `.golangci.yml`. `make go-lint` checks formatting and lint, and `make resume-check` runs it along with the tests.

## Private fields

The phone number lives in `resume/resume.private.yaml`, which is git-ignored (start from `resume.private.example.yaml`). `-private` merges it; `make resume-private` renders with it into `build/resume-private/`. Build private versions on your own machine only. CI never sees the file.
