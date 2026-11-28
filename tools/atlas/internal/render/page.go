package render

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/layout"
	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"
)

// Drawing constants, in SVG user units where they are lengths.
const (
	// laneLineChars is the most characters a line of a lane's name holds in the gutter.
	laneLineChars = 13
	// bandShades is how many fills the bands alternate between.
	bandShades = 2
	// laneTextDrop puts a lane's name below its band's top edge.
	laneTextDrop = 24
	// yearAbove and yearBelow place the century labels outside the bands.
	yearAbove = 11
	yearBelow = 21
	// knotWidth and knotHalf set the knot bar's width (drawn 5 wide).
	knotWidth = 5
	knotHalf  = 2.5
	// knotHitWidth and knotHitHalf set the invisible tap target behind each knot.
	knotHitWidth = 36
	knotHitHalf  = 18
	// knotHitPad extends the tap target slightly above the bar and below the label.
	knotHitPad = 5
	// stemMinGap is the smallest gap between a bar's foot and its band foot that
	// earns a drop stem down to the label.
	stemMinGap = 4
	// rejectHalf and rejectHeight set the perpendicular bar ending a rejection.
	rejectHalf   = 4.5
	rejectHeight = 9
	// half splits the chart into left and right halves for card placement.
	half = 0.5
)

// page is the markup being written for one plate.
type page struct {
	b        strings.Builder
	p        *plate.Plate
	l        *layout.Layout
	incoming map[string][]*plate.Link
	outgoing map[string][]*plate.Link
	knotX    map[string]float64
}

// newPage indexes the plate's links and knot x coordinates by person id.
func newPage(p *plate.Plate, l *layout.Layout) *page {

	pg := &page{
		p: p, l: l,
		incoming: map[string][]*plate.Link{}, outgoing: map[string][]*plate.Link{},
		knotX: make(map[string]float64, len(l.Knots)),
	}
	for i := range p.Links {
		link := &p.Links[i]
		pg.incoming[link.To] = append(pg.incoming[link.To], link)
		pg.outgoing[link.From] = append(pg.outgoing[link.From], link)
	}

	for _, k := range l.Knots {
		pg.knotX[k.Person] = k.X
	}

	return pg
}

// line writes s and a newline.
func (pg *page) line(s string) {
	pg.b.WriteString(s)
	pg.b.WriteByte('\n')
}

// linef writes a formatted line.
func (pg *page) linef(format string, args ...any) {
	fmt.Fprintf(&pg.b, format, args...)
	pg.b.WriteByte('\n')
}

// name returns a person's full name, set for HTML.
func (pg *page) name(id string) string { return smart(pg.p.Person(id).Name) }

// head writes the opening of <main>: the folio frame, top bar, title and subtitle.
func (pg *page) head() {

	pg.line(`<main id="main" class="atlas">`)
	pg.line(`<div class="folio-frame" aria-hidden="true"><i class="boss"></i><i class="boss"></i>` +
		`<i class="boss"></i><i class="boss"></i></div>`)
	pg.line(`<div class="topbar"><a class="home" href="/">jonreyna.dev</a></div>`)
	pg.line(`<header class="atlas-head">`)

	if pg.p.Kicker != "" {
		pg.linef(`<p class="kicker">%s</p>`, smart(pg.p.Kicker))
	}

	pg.linef(`<h1>%s</h1>`, smart(pg.p.Title))
	pg.linef(`<p class="sub">%s</p>`, smart(pg.p.Subtitle))
	pg.line(`<p class="how"><span class="pil" aria-hidden="true">¶</span>` +
		`Hover a thread or a name to follow it; select a name for its card and sources,` +
		` or <a href="#outline-body">read it all as an outline</a>.</p>`)
	pg.line(`</header>`)
}

// legend writes the threads, each a link that lights it, and the four kinds of link.
func (pg *page) legend() {

	pg.line(`<nav class="legend" aria-label="Threads">`)
	pg.line(`<ul class="threads">`)

	for _, t := range pg.p.Threads {
		pg.linef(`<li id="t-%[1]s" class="t-%[1]s"><a href="#t-%[1]s"><span class="swatch"></span>%[2]s</a>`+
			`<span class="summary">%[3]s</span></li>`, t.ID, smart(t.Name), smart(t.Summary))
	}

	pg.line(`</ul>`)
	pg.line(`<ul class="kinds">`)

	for _, k := range []struct{ class, svg, label string }{
		{"y-quotes", `<path d="M0 6H27"/><circle class="dot-solid" cx="27" cy="6" r="2.6"/>`, "quotes or names"},
		{"y-adapts", `<path d="M0 6H27"/><circle class="dot-hollow" cx="27" cy="6" r="2.8"/>`, "adapts"},
		{"y-rejects", `<path d="M0 6H29M29 1.5V10.5"/>`, "rejects"},
		{"y-argued", `<path d="M0 6H27"/><circle class="dot-argued" cx="27" cy="6" r="2.4"/>`, "scholars argue"},
	} {
		pg.linef(`<li><svg class="mark %s" viewBox="0 0 34 12" aria-hidden="true">%s</svg>%s</li>`,
			k.class, k.svg, k.label)
	}

	pg.line(`</ul>`)
	pg.line(`<p class="all"><a href="#chart">Show every thread</a></p>`)
	pg.line(`</nav>`)
}

// chart writes the figure: the sticky lane rail for narrow screens, the SVG and its caption.
func (pg *page) chart() {

	pg.line(`<figure class="chart" id="chart">`)
	pg.line(`<div class="frame">`)
	pg.line(`<div class="lane-rail" aria-hidden="true">`)

	for _, b := range pg.l.Bands {
		pg.linef(`<span class="lane-tag lane-%s">%s</span>`, b.Lane, smart(b.Name))
	}

	pg.line(`</div>`)
	pg.linef(`<svg viewBox="0 0 %s %s" role="group" aria-label="%s: the chart. Each name is a link to its card.">`,
		num(pg.l.Width), num(pg.l.Height), attr(pg.p.Title))
	pg.bands()
	pg.axis()
	pg.threads()
	pg.knots()
	pg.line(`</svg>`)
	pg.line(`</div>`)
	pg.linef(`<figcaption><span class="pil" aria-hidden="true">¶</span>%s`+
		` The time scale stretches where the writers crowd; the century marks show by how much.</figcaption>`,
		smart(pg.p.Scope))
	pg.line(`</figure>`)
}

// bands writes the lanes and their names.
func (pg *page) bands() {

	pg.line(`<g class="bands">`)

	for i, b := range pg.l.Bands {
		pg.linef(`<rect class="band band-%d" x="0" y="%s" width="%s" height="%s"/>`,
			i%bandShades, num(b.Top), num(pg.l.Width), num(b.Bottom-b.Top))

		lines := laneLines(b.Name)
		pg.linef(`<text class="lane" x="14" y="%s">`, num(b.Top+laneTextDrop))

		for j, s := range lines {
			dy := "0"
			if j > 0 {
				dy = "1.2em"
			}

			pg.linef(`<tspan x="14" dy="%s">%s</tspan>`, dy, smart(s))
		}

		pg.line(`</text>`)
	}

	pg.line(`</g>`)
}

// axis writes the century marks, with labels above and below the bands.
func (pg *page) axis() {

	top, bottom := pg.l.Bands[0].Top, pg.l.Bands[len(pg.l.Bands)-1].Bottom
	pg.line(`<g class="axis" aria-hidden="true">`)

	for _, t := range pg.l.Ticks {
		pg.linef(`<line class="tick" x1="%[1]s" y1="%[2]s" x2="%[1]s" y2="%[3]s"/>`, num(t.X), num(top), num(bottom))
		if t.Label != "" {
			pg.linef(`<text class="year" x="%[1]s" y="%[2]s">%[4]s</text><text class="year" x="%[1]s" y="%[3]s">%[4]s</text>`,
				num(t.X), num(top-yearAbove), num(bottom+yearBelow), t.Label)
		}
	}

	pg.line(`</g>`)
}

// threads writes every path, then a wide invisible copy of each on top to catch the
// pointer and carry its tooltip.
func (pg *page) threads() {

	pg.line(`<g class="threads" aria-hidden="true">`)

	for i := range pg.l.Segments {
		s := &pg.l.Segments[i]
		pg.linef(`<path class="%s" d="%s"/>`, pg.segmentClass(s, "s"), s.D)
	}

	for i := range pg.l.Segments {
		s := &pg.l.Segments[i]
		pg.linef(`<path class="%s" d="%s"><title>%s</title></path>`, pg.segmentClass(s, "hit"), s.D, pg.tooltip(s))
	}

	pg.line(`</g>`)
}

// segmentClass returns a path's classes: its kind, thread, firmness, link type and
// the knots it belongs to.
func (pg *page) segmentClass(s *layout.Segment, kind string) string {

	c := []string{kind, "t-" + s.Thread, "f-" + strconv.Itoa(s.Strength)}
	if s.Link != "" {
		c = append(c, "b", "y-"+string(s.Type))
	}

	for _, id := range s.People {
		c = append(c, "k-"+id)
	}

	return strings.Join(slices.Compact(c), " ")
}

// tooltip says what a path shows: the link for a branch, the thread and its source
// for a carrier.
func (pg *page) tooltip(s *layout.Segment) string {

	thread := smart(pg.p.Thread(s.Thread).Name)
	if s.Link == "" {
		return thread + ", from " + pg.name(s.People[0])
	}

	return pg.name(s.People[1]) + " " + s.Type.Verb() + " " + pg.name(s.People[0]) + " · " + thread
}

// knots writes each knot as a link to its card, then the marks where links arrive.
func (pg *page) knots() {

	pg.line(`<g class="knots">`)

	for _, k := range pg.l.Knots {
		person := pg.p.Person(k.Person)
		n := len(pg.incoming[k.Person]) + len(pg.outgoing[k.Person])
		pg.linef(`<a class="%s" href="#k-%s" aria-label="%s, %s, %d links">`,
			pg.knotClass(k.Person), k.Person, attr(person.Name), attr(person.Dates()), n)
		pg.linef(`<rect class="knot-hit" x="%s" y="%s" width="%d" height="%s"/>`,
			num(k.X-knotHitHalf), num(k.Top-knotHitPad), knotHitWidth, num(k.LabelY-k.Top+knotHitPad+knotHitPad))

		if k.Bottom < k.StemBottom-stemMinGap {
			pg.linef(`<line class="stem" x1="%[1]s" y1="%[2]s" x2="%[1]s" y2="%[3]s"/>`,
				num(k.X), num(k.Bottom), num(k.StemBottom))
		}

		pg.linef(`<rect class="bar" x="%s" y="%s" width="%d" height="%s" rx="%s"/>`,
			num(k.X-knotHalf), num(k.Top), knotWidth, num(k.Bottom-k.Top), num(knotHalf))
		pg.linef(`<text class="name" x="%s" y="%s">%s</text>`, num(k.X), num(k.LabelY), smart(k.Label))
		pg.line(`</a>`)
	}

	pg.ends()
	pg.line(`</g>`)
}

// ends writes the arrival marks on every knot's left edge, shaped by link type.
func (pg *page) ends() {

	pg.line(`<g class="ends" aria-hidden="true">`)

	for _, k := range pg.l.Knots {
		for _, e := range k.Ends {
			class := "end t-" + e.Thread + " y-" + string(e.Type) + " k-" + e.From + " k-" + k.Person

			switch e.Type {
			case plate.LinkRejects:
				pg.linef(`<path class="%s" d="M%s %sv%d"/>`, class, num(e.X-1), num(e.Y-rejectHalf), rejectHeight)
			case plate.LinkAdapts:
				pg.linef(`<circle class="%s dot-hollow" cx="%s" cy="%s" r="2.8"/>`, class, num(e.X), num(e.Y))
			case plate.LinkArgued:
				pg.linef(`<circle class="%s dot-argued" cx="%s" cy="%s" r="2.4"/>`, class, num(e.X), num(e.Y))
			case plate.LinkQuotes, plate.LinkNames:
				pg.linef(`<circle class="%s dot-solid" cx="%s" cy="%s" r="2.6"/>`, class, num(e.X), num(e.Y))
			}
		}
	}

	pg.line(`</g>`)
}

// knotClass returns a knot's classes: its own id, every thread that touches it, and
// every neighbor joined to it (so CSS focus keeps the neighborhood bright).
func (pg *page) knotClass(id string) string {

	threads := map[string]bool{}
	peers := map[string]bool{}

	for _, l := range pg.incoming[id] {
		threads[l.Thread] = true
		peers[l.From] = true
	}

	for _, l := range pg.outgoing[id] {
		threads[l.Thread] = true
		peers[l.To] = true
	}

	tids := make([]string, 0, len(threads))
	for t := range threads {
		tids = append(tids, "t-"+t)
	}

	pids := make([]string, 0, len(peers))
	for peer := range peers {
		pids = append(pids, "n-"+peer)
	}

	slices.Sort(tids)
	slices.Sort(pids)

	parts := make([]string, 0, 2+len(tids)+len(pids))
	parts = append(parts, "knot", "k-"+id)
	parts = append(parts, tids...)
	parts = append(parts, pids...)

	return strings.Join(parts, " ")
}

// laneLines breaks a lane's name into short lines at spaces.
func laneLines(name string) []string {

	var lines []string

	for w := range strings.FieldsSeq(name) {
		n := len(lines)
		if n > 0 && len(lines[n-1])+1+len(w) <= laneLineChars {
			lines[n-1] += " " + w

			continue
		}

		lines = append(lines, w)
	}

	return lines
}

// cards writes a card for every person, each shown when its address is opened.
func (pg *page) cards() {

	pg.line(`<section class="cards" aria-label="People">`)

	for i := range pg.p.People {
		pg.card(&pg.p.People[i])
	}

	pg.line(`</section>`)
}

// card writes one person's card: who they were, then the evidence for every link
// that arrives at them, then the links that leave them.
func (pg *page) card(person *plate.Person) {

	side := "card-right"
	if pg.knotX[person.ID] > pg.l.Width*half {
		side = "card-left"
	}

	pg.linef(`<article class="card %[2]s" id="k-%[1]s" aria-labelledby="k-%[1]s-name">`, person.ID, side)
	pg.line(`<a class="scrim" href="#chart" tabindex="-1" aria-label="Close card"></a>`)
	pg.line(`<div class="card-panel">`)
	pg.line(`<p class="close"><a href="#chart">Close</a></p>`)
	pg.linef(`<h2 id="k-%s-name">%s</h2>`, person.ID, smart(person.Name))
	pg.linef(`<p class="dates">%s</p>`, smart(person.Dates()))

	for _, s := range []string{person.Members, person.Dating, person.Card} {
		if s != "" {
			pg.linef(`<p>%s</p>`, smart(s))
		}
	}

	if in := pg.incoming[person.ID]; len(in) > 0 {
		pg.line(`<h3>Sources &amp; responses</h3>`)
		pg.line(`<ul class="links">`)

		for _, l := range in {
			pg.evidence(l, person.ID)
		}

		pg.line(`</ul>`)
	}

	if out := pg.outgoing[person.ID]; len(out) > 0 {
		pg.line(`<h3>Cited or answered by</h3>`)
		pg.line(`<ul class="links">`)

		for _, l := range out {
			pg.linef(`<li class="t-%s y-%s">%s <span class="ref">(<cite>%s</cite>)</span></li>`,
				l.Thread, l.Type, pg.claim(l, person.ID), smart(l.Cite))
		}

		pg.line(`</ul>`)
	}

	pg.linef(`<p class="more"><a href="https://en.wikipedia.org/wiki/%s">Wikipedia</a></p>`, attr(person.Wiki))
	pg.line(`</div>`)
	pg.line(`</article>`)
}

// claim writes a link as a sentence ("Justin names Heraclitus") with a link to the
// thread and to each person other than selfID.
func (pg *page) claim(l *plate.Link, selfID string) string {

	return fmt.Sprintf(`<a class="thread" href="#t-%s">%s</a> %s <b>%s</b> %s`,
		l.Thread, smart(pg.p.Thread(l.Thread).Name), pg.personRef(l.To, selfID), l.Type.Verb(), pg.personRef(l.From, selfID))
}

// personRef returns a person's name linked to their card, or plain text when it is
// the card already open.
func (pg *page) personRef(id, selfID string) string {

	if id == selfID {
		return pg.name(id)
	}

	return fmt.Sprintf(`<a href="#k-%s">%s</a>`, id, pg.name(id))
}

// evidence writes one link with everything that supports it.
func (pg *page) evidence(l *plate.Link, selfID string) {

	pg.linef(`<li class="t-%s y-%s" id="%s">`, l.Thread, l.Type, l.ID)
	pg.linef(`<p class="claim">%s</p>`, pg.claim(l, selfID))

	if l.Quote != "" {
		pg.linef(`<blockquote><p>%s</p></blockquote>`, smart(l.Quote))
	}

	if l.Paraphrase != "" {
		pg.linef(`<p class="paraphrase">%s</p>`, smart(l.Paraphrase))
	}

	ref := "<cite>" + smart(l.Cite) + "</cite>"
	if s, ok := pg.p.Sources[l.Source]; ok {
		ref += ", " + pg.source(s, l.URL)
	}

	pg.linef(`<p class="ref">%s</p>`, ref)

	if c := l.Compare; c != nil {
		text := "Compare <cite>" + smart(c.Cite) + "</cite>"
		if c.Quote != "" {
			text += ": “" + smart(c.Quote) + "”"
		}

		if c.URL != "" {
			text += ` (<a href="` + attr(c.URL) + `">read it</a>)`
		}

		pg.linef(`<p class="compare">%s</p>`, text)
	}

	for _, s := range l.Scholarship {
		pg.linef(`<p class="scholar">%s</p>`, smart(s))
	}

	if l.Note != "" {
		pg.linef(`<p class="note">%s</p>`, smart(l.Note))
	}

	if l.Checked == nil {
		pg.line(`<p class="unchecked">Not yet checked</p>`)
	}

	pg.line(`</li>`)
}

// source names a translation, linked to where the passage can be read.
func (pg *page) source(s plate.Source, link string) string {

	text := smart(s.Name)
	if link != "" {
		text = `<a href="` + attr(link) + `">` + text + `</a>`
	}

	if s.Via != "" {
		text += ", via " + smart(s.Via)
	}

	return text
}

// outline writes the whole plate as text: each thread with its links, then the
// sources.
func (pg *page) outline() {

	pg.line(`<details class="outline" id="outline">`)
	pg.line(`<summary>Read as an outline</summary>`)
	pg.line(`<div class="outline-body" id="outline-body">`)
	pg.linef(`<p>%s</p>`, smart(pg.p.Scope))

	for _, t := range pg.p.Threads {
		pg.linef(`<section class="t-%s">`, t.ID)
		pg.linef(`<h2>%s</h2>`, smart(t.Name))
		pg.linef(`<p>%s</p>`, smart(t.Summary))
		pg.line(`<ol>`)

		for i := range pg.p.Links {
			if l := &pg.p.Links[i]; l.Thread == t.ID {
				pg.linef(`<li>%s (<cite>%s</cite>): %s</li>`, pg.claim(l, ""), smart(l.Cite), quoted(l))
			}
		}

		pg.line(`</ol>`)
		pg.line(`</section>`)
	}

	pg.bibliography()
	pg.line(`</div>`)
	pg.line(`</details>`)
}

// bibliography lists the translations quoted and the scholarship cited.
func (pg *page) bibliography() {

	used := map[string]bool{}

	var scholarship []string

	for i := range pg.p.Links {
		l := &pg.p.Links[i]
		used[l.Source] = true
		if l.Compare != nil {
			used[l.Compare.Source] = true
		}

		scholarship = append(scholarship, l.Scholarship...)
	}

	var sources []string

	for key, s := range pg.p.Sources {
		if used[key] {
			sources = append(sources, s.Name+", via "+s.Via)
		}
	}

	byText := func(a, b string) int {
		return strings.Compare(strings.ReplaceAll(a, "*", ""), strings.ReplaceAll(b, "*", ""))
	}
	slices.SortFunc(sources, byText)
	slices.SortFunc(scholarship, byText)

	pg.line(`<section class="bibliography">`)
	pg.line(`<h2>Sources</h2>`)

	for _, group := range []struct {
		title string
		items []string
	}{{"Translations quoted", sources}, {"Scholarship cited", slices.Compact(scholarship)}} {
		pg.linef(`<h3>%s</h3>`, group.title)
		pg.line(`<ul>`)

		for _, s := range group.items {
			pg.linef(`<li>%s</li>`, smart(s))
		}

		pg.line(`</ul>`)
	}

	pg.line(`</section>`)
}

// quoted returns a link's quote in quotation marks, or its paraphrase.
func quoted(l *plate.Link) string {

	if l.Quote != "" {
		return "“" + smart(l.Quote) + "”"
	}

	return smart(l.Paraphrase)
}
