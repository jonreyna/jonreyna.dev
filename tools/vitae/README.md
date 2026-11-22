# vitae

Turns `resume/resume.yaml` into every résumé format and variant. It's the only code that understands variants: the site reads its JSON output instead of filtering again.

```
vitae validate [-i FILE] [--private FILE] [--strict]
vitae render   [-i FILE] [--private FILE] [--variant NAME|all] [--format json|all] [-o DIR|-]
vitae list     [-i FILE] variants|tags
vitae version
```

Exit codes: 0 success, 1 validation failure, 2 usage error. From the repository root, `make resume-check` tests and validates, and `make resume` renders into `build/resume/`.

## How it works

- **`internal/model`**: the YAML as Go types. Unknown fields are errors, so a typo fails loudly.
- **`internal/lint`**: structure and content checks. Errors always fail; warnings fail with `--strict`, which CI uses. Among the rules:
  - roles newest first, and only the first one current;
  - every tag declared;
  - highlights under 650 characters;
  - curly quotes and real dashes;
  - no phone number in the public file;
  - an education entry must say whether the degree was conferred, and if it wasn't, only its note may mention the degree.
- **`internal/filter`**: applies a variant. Each highlight scores its weight plus `boost_by` (default 20) per boosted tag. Excluded tags and scores under `min_weight` drop out, and the best `max_highlights` survive in their original order.
- **`internal/inline`**: the markup allowed in text fields (`*emphasis*`, `` `code` ``, `[text](https://url)`, `~` for a non-breaking space), written out per format.
- **`internal/render`**: the outputs. JSON so far. LaTeX, plain text and Markdown come next, then PDF (latexmk) and Word (Pandoc) from `make`.

Output is deterministic, so `testdata/golden/` holds exact expected files for `testdata/resume.yaml`. After an intended change, run `go test ./... -update` and review the diff.

## Private fields

The phone number lives in `resume/resume.private.yaml`, which is git-ignored (start from `resume.private.example.yaml`). `--private` merges it; `make resume-private` renders with it into `build/resume-private/`. Build private versions on your own machine only. CI never sees the file.
