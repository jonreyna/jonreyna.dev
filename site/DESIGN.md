# DESIGN.md: jonreyna.dev

**Direction J, "Red and black, with algorithms."** Chosen 2026-11-20 from ten specimens and revised the same day after a design review (see "Review log" at the end). The token values live in `web/tokens.css` and the interim favicon is `web/favicon.svg`. The page `/specimen` (`content/specimen.md`, not linked and not indexed) exercises every component below; check it after any change to the styles.

## The idea in one paragraph

The site looks like an old English hardback, printed in red and black on pale parchment, the way two-color title pages once were. Inside, the technical essays take their apparatus from Knuth's *The Art of Computer Programming*: numbered equations, algorithms set step by step, and rated exercises with answers at the back. It should read as literary and as rigorous, never as costume. Keep the restraint: no runes, maps, page borders, paper textures, script fonts or imitation of any particular book's art.

## Principles

1. **Red does the work color would do.** Red marks structure: headings, initials, step labels, links, ratings. Everything else is ink. Never add a third hue.
2. **Typography, not decoration.** Hierarchy comes from small capitals, italics, size and red. The only ornaments are the masthead's rules and one asterism (⁂) before the notes.
3. **One column, set like a book.** About 69 characters per line, indented paragraphs with no space between them, hyphenation on, and line breaking that avoids lone words, as TeX does.
4. **Reading pages ship no JavaScript.** Math is MathML rendered at build time. Interactive pieces load only when asked for.
5. **Dark mode is the same book in leather and copper.** The same roles, the same names, new values.
6. **A recruiter finds the résumé in five seconds.** The home page says who you are, links the résumé and gives the email before any list.

## Color

| Token | Light | Dark | Role | Contrast on `--bg` (light / dark) |
|---|---|---|---|---|
| `--bg` | `#f6f0e1` | `#1a1612` | Page ground | — |
| `--text` | `#2a2119` | `#e9e0cd` | Body ink, the masthead name, table rules | 13.9 / 13.7 |
| `--muted` | `#6d5e4a` | `#a8997f` | Dates, bylines, abstracts, captions, code comments, equation numbers, the contact line's plain text | 5.5 / 6.5 |
| `--red` | `#9a2a17` | `#e47d63` | Section heads, the drop cap, links, nav separators, algorithm names and step labels, ratings, ▸, ⁂, masthead rules | 6.8 / 6.4 |
| `--red-hover` | `#6a1a0c` | `#f3ab98` | Link hover and active | 10.5 / 9.5 |
| `--rule` | `#d9ccb2` | `#3e352b` | Quiet dividers (dotted page labels). Never used for text. | 1.4 / 1.5 |
| `--code-bg` | `#efe6d1` | `#231e18` | Code blocks. Red and muted text stay above 5:1 on it. | — |
| `--link-line` | red at 45% | (follows red) | The underline under links in running prose | — |
| `--selection` | red at 22% | (follows red) | Selected text | — |
| `--highlight` | red at 12% | (follows red) | The footnote, step, exercise or heading a link just jumped to | — |
| `--code-shade` | ink at 16% | (follows ink) | The fade at the edge of a code block that scrolls | — |

Red and ink differ by only about 2:1, so color alone can't mark a link inside a sentence. **Links in running prose are underlined** with `--link-line`. Links that stand alone (the nav, lists, the contact line, the footer, exercise numbers) need no underline.

Text sizes and colors all pass WCAG AA. No success, warning or error colors exist yet, because the site has no forms. If it ever needs them, define them here first.

## Typography

| Family | Token | Weights | Use | License |
|---|---|---|---|---|
| Alegreya | `--font-text` | 400, 400 italic, 500 | Body text, italics, essay titles, work names | OFL |
| Alegreya SC | `--font-sc` | 500 only | The name, nav, section heads, algorithm names, table heads, the Exercises heading | OFL |
| Courier Prime | `--font-code` | 400, 700 | Code: keywords bold, comments italic in `--muted` | OFL |
| Libertinus Math | `--font-math` | — | MathML | OFL |

**Fallback while fonts load.** `tokens.css` defines "Alegreya Fallback": Georgia resized to Alegreya's width and line metrics (`size-adjust: 91.6%`, `ascent-override: 111.4%`, `descent-override: 38.2%`). Measured in Chrome, the sample paragraphs take exactly the same height in the fallback as in Alegreya, against 7% taller in plain Georgia, so text doesn't jump when the font arrives.

**Font budget.** Self-host everything as WOFF2, subset to the characters the site uses (Latin, the math and arrow symbols in the essays, ▸ ▮ ⁂ ↩ §). Preload only Alegreya regular. Aim for about 150kB on a first visit. Libertinus Math is 396kB as WOFF2 (release 7.051, which `make fonts` downloads and checks; no font file is committed); subset it in CI (`pyftsubset`, keeping the MATH table) to under 120kB. Browsers fetch it only on pages that contain math. No Google Fonts requests on the live site. All four families are under the SIL Open Font License: subsetting is allowed (none of the names used here is a Reserved Font Name), and each family's `OFL.txt` ships in `web/fonts/` beside its files.

**Scale** (rem; the root is 21px on desktop and 18px at 640px or narrower):

| Token | Size | Use | Leading |
|---|---|---|---|
| `--size-title` | 1.9 | Essay title (weight 500, balanced lines) | 1.2 |
| `--size-name` | 1.65 | The masthead name (small capitals, tracking 0.1em) | 1.25 |
| `--size-h2` | 1.1 | Section heads (small capitals, tracking 0.06em, red) | 1.3 |
| `--size-body` | 1 | Body | 1.5 |
| `--size-small` | 0.94 | Exercises and answers | 1.5 |
| `--size-note` | 0.86 | Footnotes, captions, tables | 1.5 |
| `--size-label` | 0.75 | Small-capital labels | 1.5 |
| `--size-code` | 0.74 | Code blocks | 1.5 |

**Features:** old-style figures in text (`oldstyle-nums`), tabular old-style figures in tables, kerning and common ligatures on, discretionary ligatures off. Use real small capitals from Alegreya SC, never synthesized ones.

**Line breaking:** every paragraph uses `text-wrap: pretty` (the browser's nearest thing to TeX's whole-paragraph line breaking, so last lines aren't a lone word) and `hyphenate-limit-chars: 7 3 3` (only words of seven letters or more are hyphenated, keeping at least three letters on each side). Headings use `text-wrap: balance`. Text is ragged-right; justified text was considered and not adopted.

**The drop cap** (`--size-initial`, red, about three lines) opens an essay's first section only, like a chapter in a book. Never on the home page and never on every section.

## Layout and spacing

- **Column:** `--measure` = 30rem (630px on desktop). Measured in Chrome: 69 characters per line on average, 75 at most; about 47 on a 375px phone. Centered, with `--gutter` (16px) side margins on phones. No sidebars. The page never scrolls sideways; wide tables and code scroll inside their own box.
- **Paragraphs:** no space between them; each one after the first is indented by `--indent` (1.4em). A paragraph that follows a heading, list, figure or algorithm is not indented.
- **Space scale:** `--space-1` to `--space-6` = 0.3, 0.6, 1, 1.4, 2, 3 rem. Use 2rem before a section head and 3rem between page regions.
- **Shape:** no border radius and no shadows, anywhere.

## Components

**Skip link.** The first thing in the page: "Skip to content," hidden until it receives keyboard focus, then shown at the top left in small capitals with a red outline.

**Masthead.** Centered, like a title page: a red double rule above (`--rule-double`), the name in small capitals, an italic tagline in `--muted`, then the nav in small capitals with red middle dots between items. A single red rule below. On the home page the name is the page's `<h1>`; on every other page it's a `<p class="name">` and the page title is the `<h1>`. The nav has four items, **work · writing · play · résumé**, which fit on one line on a 375px phone (the dots disappear below 640px). They go to `/work`, `/writing`, the home page's Play section (until Play has a page of its own, Phase 6) and `/resume`. The item for the current section is red and carries `aria-current="page"`; essays count as part of Writing. About lives in the footer.

**Section head (`h2`).** Red small capitals. In essays, chapters are numbered with roman numerals (`I.`, `II.`, `III.`) and carry an id. The build does both (`lib/essay_sections.dart`): authors write a plain `## The rule` and get the numeral, the id `the-rule` and the § link. Other pages get the id and the § link without a numeral; a hand-written `<h2 id>` keeps its id. A § link to the section appears beside the heading on hover or keyboard focus, and at half strength on touch screens.

**Home page.** The opening paragraph (no drop cap), then the **contact line**: "Résumé (PDF) · jon@jonreyna.dev · open to staff and principal roles (remote)," in italics, with the links in red, the plain words in `--muted` and red middle dots between. Drop or change the last phrase when you're no longer looking. Then three sections: *Selected work* (four entries, then a link to `/work`), *Writing* (the three newest essays, from `<EssayList limit="3"/>`) and *Play*. A work entry is run in: the name in weight 500, then employer and years in italic `--muted`, then one or two sentences. Lists use a red en dash as the bullet.

**Work page** (`/work`). A short opening paragraph, then one section per employer, newest first, headed "Equinix, 2023–2026". Each entry is run in like the home page's, at more length: two to four sentences on what it was, what was hard and what it still does.

**Writing list** (`/writing`, and on the home page). The same run-in form: the essay's title as a link in weight 500, its month and year in italic `--muted`, then its abstract (or its description). Newest first. Essays live in `content/writing/` and appear once they have `kind: essay`, a `date` and no `draft` or `noindex`. With none, the list says "Nothing published yet."

**Résumé page** (`/resume`, `<Resume/>` in `lib/components/resume_view.dart`). Built from `resume.full.json`, which `vitae` renders from `resume/resume.yaml` in the same build. First the contact line (label, location, email, LinkedIn and GitHub, with the profile URLs as their text) and a "Download:" line with the full résumé's files, so a recruiter has the PDF in five seconds (principle 6). Then the summary, unindented, and the sections in the order the JSON lists them. Roles and schools are run in like the home page's work entries: the employer in weight 500, the title in roman, the place and dates in italic `--muted`. Highlights are a red-dash list; each opens with its lead in italics. Last, *Downloads*: one line per variant, labeled with its description, with a link per format. The page links only files the build actually produced, so no link is dead, and the build fails if `resume.full.json` is missing.

**Print.** One stylesheet for every page, written for the résumé: black ink on white, 10.5pt, no masthead rules, nav, motto, footer, § links or page title, and no Downloads section. A link to another site prints its URL after it, unless its text is already the URL (class `bare`). Headings stay with what follows; list items and entries don't split across pages.

**Essay foot.** A centered italic line after the essay, before the footer: "Written by Jon Reyna. All writing."

**Essay head.** Centered: the title (`<h1>`), "by Jon Reyna" in italics, and the month and year in red small capitals inside `<time datetime>`. Then an optional abstract: centered, italic, `--muted`, at most about 50 words, inset 1.5rem on each side (no inset on phones).

**Equation.** A display equation is centered. When the text refers to it, give it an id (`eq-1`) and a number at the right margin, `(1)`, in `--muted`; references read "equation (1)" or "(1)" and link to it.

**Algorithm.** Set the way Knuth sets them:

```html
<div class="algorithm" id="algorithm-b" role="group" aria-label="Algorithm B">
  <p class="head"><span class="name">Algorithm B</span> (<i>Backpropagation</i>).
     One sentence on the input and what the algorithm computes.</p>
  <ol>
    <li id="step-b1"><span class="step">B1.</span>[Initialize.] Set …</li>
    <li id="step-b2"><span class="step">B2.</span>[Propagate.] For each …</li>
    <li id="step-b3"><span class="step">B3.</span>[Decrease <i>i</i>.] … return to <a href="#step-b2">B2</a>;
        otherwise the algorithm terminates. <span class="qed" aria-label="End of algorithm">▮</span></li>
  </ol>
</div>
```

The name is red small capitals; the step labels are red, weight 500, hung in a `--hang` (2.7rem) column; the step names stay roman inside brackets; assignment is `←`; the end mark `▮` is red. Every mention of a step elsewhere ("step B2") links to it. In Markdown this will be a fenced block that the build turns into this markup.

**Exercises.** An `<h2 class="exercises">`, "Exercises," centered in red small capitals with wide tracking, then a numbered list at `--size-small`. Each exercise has an id (`exercise-1`); its number links to its answer. Each starts with a rating in red italics inside upright brackets. Knuth's scale: `[00]` immediate, `[10]` a minute, `[20]` up to twenty minutes, `[30]` a few hours, `[40]` a long project, `[50]` an open problem; a prefix `M` means mathematical and `HM` higher mathematics. A red `▸` in the margin marks a recommended exercise.

**Answers.** After the notes, an `<h2 class="exercises">`, "Answers to exercises," with one item per exercise (`answer-1`). Each answer's number links back to its exercise. Answers are short; a long one links to its own page.

**Code.** Courier Prime on `--code-bg`, no border, no radius, tab width 4, scrolling sideways inside its box. Keywords bold, comments italic in `--muted`. **Keep lines to 60 characters or fewer**, so listings fit without scrolling on tablets and desktops. On phones they still scroll, and a soft fade (`--code-shade`) appears at whichever edge has more code beyond it; the fade is pure CSS. Inline code is 0.86em with no background and never breaks across lines, so keep it short.

**Table.** Booktabs style: a 1.5px ink rule above and below, a 0.75px rule under the head row, no vertical rules. Head cells in small capitals. The caption sits above, in italics, with "Table 1." in red small capitals. Give the figure an id (`table-1`).

**Notes.** A centered red ⁂, then the numbered notes at `--size-note` with red numbers. Each note ends with a ↩ link back to the text, kept on the same line as the note's last word. Reference marks in the text are superscript links with no underline, padded to a target about 24px square without moving the text.

**Jump highlight.** Whatever a link jumps to (a note, step, exercise, answer or section head) gets a `--highlight` background so the reader's eye lands on it. Every element with an id has a 1rem scroll margin.

**Footer.** A single red rule, then **about** · email · GitHub · LinkedIn · feed, centered, in italics.

## Typesetting rules for the build

The Markdown processor applies these, so authors never have to:

- **Ties** (TeX's `~`): a non-breaking space in "Algorithm B," "step B2," "Table 1," "equation (1)," "Exercise 2," "by 1," and between a number and its unit or noun where a break would strand it ("2,317 store," "24 teams," "2 × model size").
- **Punctuation stays with math:** a comma or period after inline math is wrapped with it so it can't start a line. The last word before inline math at the end of a sentence stays with it too ("for all *j*.").
- **Dashes and signs:** en dash for ranges (2019–2020), em dash without spaces for breaks, true minus (−) and multiplication (×) signs, curly quotes and apostrophes.
- **Back links** (↩) are tied to the last word of their note.
- **Math is written in TeX** and becomes MathML when the site is built (`lib/tex_math.dart`, with the pinned Pandoc; no script on the page):
  - `$…$` inline, `$$…$$` display. A `$$…$$` paragraph becomes `<div class="eq">`; a numbered equation is one line, `<div class="eq num" id="eq-1">$$…$$<span class="no">(1)</span></div>`. `\$` is a dollar sign, and a `$` before a digit or after a space isn't math, so prices need no escaping.
  - Upright names: `\mathrm{out}` or `\operatorname{out}`. Overbars: `\bar{v}`. Arrows: `\gets`.
  - TeX that Pandoc can't read stops the build, naming the page and the formula.
- **MathML rules,** applied to Pandoc's output automatically and covered by `make site-test`:
  - Parentheses get `stretchy="false"` with zero `lspace` and `rspace`, or Chrome spaces them like operators.
  - Overbars (adjoints and similar) use `<mover accent="false">` with U+00AF, or the bar collapses to a sub-pixel line.
  - A slash used as division inside a formula gets zero `lspace` and `rspace`.
  - An ellipsis is an operator; an upright name is one `<mi>`.
  - A period, comma or other punctuation right after inline math is wrapped with it (the punctuation rule above). Keeping the word before it on the line is still by hand: `<span class="nowrap">all&nbsp;$j$.</span>`.

## Page template

Every page carries:

- `<html lang="en">`, a skip link, exactly one `<h1>`, and `<main id="main">`.
- `<title>`: "Jon Reyna" on the home page; "Essay title · Jon Reyna" elsewhere.
- `<meta name="description">`, at most about 155 characters, written per page.
- `<link rel="canonical">` with the page's full https://jonreyna.dev URL.
- `<link rel="alternate" type="application/atom+xml">` pointing at the feed.
- Favicon: `favicon.svg`, which switches to leather and copper in dark mode, plus a 180px PNG `apple-touch-icon`.
- `<meta name="theme-color">` twice: `#f6f0e1` for light and `#1a1612` for dark.
- Open Graph: `og:type`, `og:title`, `og:description`, `og:url` and `og:image` (1200×630).
- Dates in `<time datetime>`.
- Essays add `article:published_time` from `date` and, when the front matter has `updated`, `article:modified_time`. A page with `noindex: true` or `draft: true` gets `<meta name="robots" content="noindex">`.

## Feed and sitemap

`make dist` writes both from the staged pages (`tool/feed_sitemap.py`), so they always match what is published:

- **`/feed.xml`**: Atom, one entry per published essay (`kind: essay`, a `date`, not `noindex` or `draft`), newest first. Each entry carries the essay's title, description as the summary, and its full text as rendered (section numbers, MathML, footnotes), with site links made absolute. Dates: `date` (`2027-03` means the 1st, `2027-03-14` midnight UTC) and `updated`.
- **`/sitemap.xml`**: every page except the 404 page and noindex pages, at its canonical URL. Only essays have a `lastmod`.
- **`/robots.txt`** (static, in `web/`): allows everything and names the sitemap. Until launch the `X-Robots-Tag: noindex` header in `_headers` keeps the whole site out of search results.

`make verify` checks that every URL in the feed and sitemap is a published page and that the sitemap misses none.

## Imagery and icons

No photos in the design and no icon set. Diagrams are inline SVG drawn in `currentColor`, with red reserved for the one thing the diagram is about. Interactive Play pieces follow the same two colors.

## Accessibility

- Contrast as listed above, AA or better throughout.
- Links in prose are underlined; every link and button shows a 2px red focus outline; the skip link works from the keyboard.
- Touch targets are at least 24px, including footnote marks.
- Heading order on every page: one `<h1>`, then `<h2>`. Exercises and Answers are `<h2>`s.
- Algorithms are ordered lists inside a labeled group; exercise ratings stay as text so screen readers read them; § and ↩ links carry labels ("Link to section I," "Back to the text").
- Text resizes freely: everything is in rem or em except the 16px gutter.

## Still open

- The final favicon: a red "J" in Alegreya SC, outlined to a path at build time. (`favicon.svg` uses Georgia for now.)
- The Open Graph image (1200×630): the masthead as a title page, on parchment.

## Review log

- **2026-11-20, design review** at 375, 768 and 1440px, light and dark, keyboard only. Applied:
  - Lines were 75 characters on average and up to 81; the measure went from 33rem to 30rem and the desktop root from 20px to 21px.
  - The name is the home page's `<h1>`.
  - The page template checklist was added.
  - No drop cap on the home page; it had left a 35px gap after "I."
  - The contact line was added under the home intro.
  - The nav went to four items; it had wrapped on phones.
  - The 60-character code rule and the edge fade were added; hidden code on a phone went from 212px to 116px.
  - Footnote targets grew from 4–6px wide to about 24px.
  - Fonts: SC is down to one weight, the fallback metrics are in place, and the subsetting budget is set.
  - TeX-style line breaking, ties, section and step links, the jump highlight, the Answers section, selection color, the skip link and theme colors.
- **Considered and not adopted:** justified text (a matter of taste; possible later for wide screens); lining figures for metrics; a narrower code font (every monospace font is about as wide).
