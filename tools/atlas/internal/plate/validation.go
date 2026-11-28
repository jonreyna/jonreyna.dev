package plate

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// maxQuoteWords is the longest quote a card may carry: a line or two, well inside
// fair use even where a translation isn't in the public domain.
const maxQuoteWords = 60

// slug matches ids that are safe in anchors and CSS class names.
var slug = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

// validation is one run of [Validate]: the plate, what has been seen so far, and the
// issues found.
type validation struct {
	p         *Plate
	issues    []Issue
	threadIDs map[string]bool
	linkIDs   map[string]bool
	used      map[string]bool
	linked    map[string]bool
}

// newValidation starts a run over p.
func newValidation(p *Plate) *validation {

	return &validation{
		p: p, threadIDs: map[string]bool{}, linkIDs: map[string]bool{},
		used: map[string]bool{}, linked: map[string]bool{},
	}
}

// errorf records an issue at where.
func (v *validation) errorf(where, format string, args ...any) {
	v.issues = append(v.issues, Issue{Where: where, Message: fmt.Sprintf(format, args...)})
}

// plate checks the plate's own fields.
func (v *validation) plate() {

	if !slug.MatchString(v.p.ID) {
		v.errorf("plate", "id %q must be lowercase letters, digits and hyphens", v.p.ID)
	}

	if v.p.Title == "" {
		v.errorf("plate", "title is required")
	}

	if v.p.Span[0] >= v.p.Span[1] {
		v.errorf("plate", "span %v must run forward in time", v.p.Span)
	}

	for key, s := range v.p.Sources {
		if s.Name == "" {
			v.errorf("source "+key, "name is required")
		}
	}
}

// lanes checks that there are lanes and that their ids are unique.
func (v *validation) lanes() {

	if len(v.p.Lanes) == 0 {
		v.errorf("plate", "at least one lane is required")
	}

	seen := map[string]bool{}
	for _, l := range v.p.Lanes {
		if seen[l.ID] || !slug.MatchString(l.ID) {
			v.errorf("lane "+l.ID, "id must be unique and lowercase")
		}

		seen[l.ID] = true
	}
}

// threads checks each thread and records its id.
func (v *validation) threads() {

	for _, t := range v.p.Threads {
		where := "thread " + t.ID
		if v.threadIDs[t.ID] || !slug.MatchString(t.ID) {
			v.errorf(where, "id must be unique and lowercase")
		}

		v.threadIDs[t.ID] = true
		if t.Name == "" || t.Color == "" {
			v.errorf(where, "name and color are required")
		}
	}
}

// person checks one person's id, years and lanes.
func (v *validation) person(p *Person) {

	where := "person " + p.ID
	if !slug.MatchString(p.ID) || v.p.Person(p.ID) != p {
		v.errorf(where, "id must be unique and lowercase")
	}

	if p.Name == "" || p.Card == "" {
		v.errorf(where, "name and card are required")
	}

	if p.From == 0 || p.To == 0 || p.From > p.To {
		v.errorf(where, "years %d to %d must run forward, with no year 0", p.From, p.To)
	}

	if p.At != 0 && (p.At < v.p.Span[0] || p.At > v.p.Span[1]) {
		v.errorf(where, "at %d is outside the span", p.At)
	}

	v.personLanes(where, p.Lanes)
}

// personLanes checks that a knot's lanes exist and sit next to each other.
func (v *validation) personLanes(where string, lanes []string) {

	if len(lanes) == 0 {
		v.errorf(where, "at least one lane is required")

		return
	}

	for i, id := range lanes {
		at := v.p.LaneIndex(id)
		if at < 0 {
			v.errorf(where, "unknown lane %q", id)

			return
		}

		if i > 0 && at != v.p.LaneIndex(lanes[i-1])+1 {
			v.errorf(where, "lanes must be listed top to bottom with none skipped")
		}
	}
}

// link checks one link's ends, type, evidence and check.
func (v *validation) link(l *Link, requireChecked bool) {

	where := "link " + l.ID
	if l.ID == "" || v.linkIDs[l.ID] {
		v.errorf(where, "id must be present and unique")
	}

	v.linkIDs[l.ID] = true
	v.used[l.Thread] = true

	if !v.threadIDs[l.Thread] {
		v.errorf(where, "unknown thread %q", l.Thread)
	}

	if !l.Type.Valid() {
		v.errorf(where, "unknown type %q", l.Type)
	}

	v.ends(where, l)
	v.evidence(where, l)

	if requireChecked && (l.Checked == nil || l.Checked.On == "") {
		v.errorf(where, "not checked")
	}
}

// ends checks that a link joins two different known people, earlier to later.
func (v *validation) ends(where string, l *Link) {

	from, to := v.p.Person(l.From), v.p.Person(l.To)
	if from == nil || to == nil {
		v.errorf(where, "unknown person in %q → %q", l.From, l.To)

		return
	}

	v.linked[l.From], v.linked[l.To] = true, true
	if l.From == l.To {
		v.errorf(where, "joins %s to itself", l.From)
	}

	if from.From > to.To {
		v.errorf(where, "runs backward in time: %s begins after %s ends", l.From, l.To)
	}
}

// evidence checks a link's citation, evidence, quote length and addresses.
func (v *validation) evidence(where string, l *Link) {

	if l.Cite == "" {
		v.errorf(where, "cite is required")
	}

	v.located(where, l)

	if l.Type == LinkArgued && len(l.Scholarship) == 0 {
		v.errorf(where, "an argued link must name its scholarship")
	}

	if n := len(strings.Fields(l.Quote)); n > maxQuoteWords {
		v.errorf(where, "quote has %d words; the most is %d", n, maxQuoteWords)
	}

	v.reference(where, l.Source, l.URL)
	if l.Compare != nil {
		v.reference(where+" compare", l.Compare.Source, l.Compare.URL)
	}
}

// located checks that a link carries what its Located value promises.
func (v *validation) located(where string, l *Link) {

	switch l.Located {
	case LocatedPassage:
		if l.Quote == "" || l.URL == "" || l.Source == "" {
			v.errorf(where, "a passage needs quote, source and url")
		}
	case LocatedParaphrase:
		if l.Paraphrase == "" {
			v.errorf(where, "a paraphrase needs paraphrase")
		}
	case LocatedScholarship:
		if len(l.Scholarship) == 0 || l.Paraphrase == "" {
			v.errorf(where, "scholarship needs scholarship and paraphrase")
		}
	default:
		v.errorf(where, "unknown located %q", l.Located)
	}
}

// reference checks a source key and a URL, either of which may be empty.
func (v *validation) reference(where, source, link string) {

	if _, ok := v.p.Sources[source]; source != "" && !ok {
		v.errorf(where, "unknown source %q", source)
	}

	if link == "" {
		return
	}

	u, err := url.Parse(link)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
		v.errorf(where, "url %q must be an http or https address", link)
	}
}

// unused reports threads no link uses and people no link joins.
func (v *validation) unused() {

	for _, t := range v.p.Threads {
		if !v.used[t.ID] {
			v.errorf("thread "+t.ID, "no link uses it")
		}
	}

	for i := range v.p.People {
		if id := v.p.People[i].ID; !v.linked[id] {
			v.errorf("person "+id, "no link joins it")
		}
	}
}
