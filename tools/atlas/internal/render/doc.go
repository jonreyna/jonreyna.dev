// Package render writes a laid-out plate as the markup and stylesheet of a Play piece
// (see play/README.md): everything inside the page's <body>, and the CSS rules that
// depend on the plate's threads and people.
//
// [HTML] writes one <main> holding the title, the legend of threads and link types,
// the chart as inline SVG, a card for every person, and an outline with a
// bibliography that gives the same content as text. [CSS] writes each thread's
// color and the rules that dim the other threads while a thread or a knot is
// hovered, focused or opened.
//
// The page needs no JavaScript. Hovering uses CSS :has(); opening a card or a thread
// uses :target, so every view has its own address (#k-justin, #t-logos). The site's
// security policy forbids inline styles, so the markup carries classes only and the
// colors come from the stylesheet: a thread's class t-ID sets --c to the pigment its
// color names.
//
// Prose from the plate is escaped, and its straight quotes and apostrophes are set
// as curly ones.
package render
