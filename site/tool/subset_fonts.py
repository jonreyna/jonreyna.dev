"""Cuts the site's fonts down to the characters it uses and writes them as WOFF2.

Run by `make fonts` with the pinned fontTools in build/venv (site/tool/requirements.txt):

    subset_fonts.py SOURCE_DIR OUT_DIR CSS MANIFEST

SOURCE_DIR holds the source fonts that make downloaded and checked. For each face in
FACES this writes OUT_DIR/<name>.<hash>.woff2, where <hash> starts the file's SHA-256,
so a changed font gets a new URL and the year-long immutable cache (site/web/_headers)
never serves a stale one. It removes WOFF2 files there that it didn't write, writes the
@font-face rules to CSS, and writes MANIFEST: the preload URLs, and each face's file,
size and the code points it covers, which site/tool/verify_dist.py checks the built
pages against.

Fails (exit 1), writing nothing, if the text faces or the math face go over their byte
budgets (DESIGN.md, "Font budget").

Same sources, same output: fontTools keeps the source's timestamps, and WOFF2's Brotli
compression is deterministic.
"""

import dataclasses
import hashlib
import io
import json
import logging
import pathlib
import sys

from fontTools import subset
from fontTools.ttLib import TTFont
from fontTools.varLib import instancer

# Body text, small capitals: Latin with its extensions, Greek, punctuation, the euro,
# trademark and minus signs, and arrows. Alegreya has no ⁂ ↩ ▮ ▸; verify_dist.py lets
# the system's symbol font draw those.
TEXT = "U+0020-007E,U+00A0-017F,U+0370-03FF,U+2000-206F,U+20AC,U+2122,U+2190-2199,U+2212"
# Code: ASCII, Latin-1, punctuation and arrows.
CODE = "U+0020-007E,U+00A0-00FF,U+2000-206F,U+2190-2199"
# MathML: Latin, Greek, operators, arrows, delimiters and their extension pieces, and the
# math alphanumerics for bold, italic, bold italic, script, double-struck and Greek
# (Chrome draws a one-letter <mi> in math italic). Fraktur, sans-serif and monospace
# letters are left out to stay under budget; add their ranges if an essay needs them.
MATH = ("U+0020-007E,U+00A0-00FF,U+0300-036F,U+0391-03A9,U+03B1-03C9,U+03D0-03F6,"
        "U+2000-206F,U+2100-214F,U+2190-21FF,U+2200-22FF,U+2308-230B,U+2329-232A,"
        "U+239B-23B5,U+23DC-23E1,U+25A1,U+25B3,U+25BD,U+25C7,U+25CB,U+27C0-27FF,"
        "U+2980-29FF,U+2A00-2AFF,U+1D400-1D4CF,U+1D538-1D56B,U+1D6A8-1D71B,U+1D7CE-1D7E1")

# The OpenType features the CSS asks for beyond fontTools' defaults (kerning,
# ligatures, marks, fractions): old-style, lining, proportional and tabular figures.
NUMERALS = ["onum", "lnum", "pnum", "tnum"]

# Bytes per group (DESIGN.md): the text faces load on every page; the math face loads
# only on pages with math.
BUDGETS = {"text": 150_000, "math": 120_000}


@dataclasses.dataclass(frozen=True)
class Face:
    """One @font-face: a source font, pinned to a weight if it is variable."""

    name: str  # the output file's stem
    source: str  # file name in SOURCE_DIR
    family: str  # the CSS font-family
    weight: int
    style: str  # "normal" or "italic"
    group: str  # "text", "code" or "math", for the budgets and verify_dist.py
    unicodes: str
    features: tuple = ()  # extra OpenType features; ("*",) keeps them all
    wght: int | None = None  # the weight to pin a variable source at
    preload: bool = False  # every page uses it, so pages fetch it with the HTML


FACES = [
    Face("alegreya-regular", "Alegreya.ttf", "Alegreya", 400, "normal", "text", TEXT,
         tuple(NUMERALS), wght=400, preload=True),
    Face("alegreya-medium", "Alegreya.ttf", "Alegreya", 500, "normal", "text", TEXT,
         tuple(NUMERALS), wght=500),
    Face("alegreya-italic", "Alegreya-Italic.ttf", "Alegreya", 400, "italic", "text", TEXT,
         tuple(NUMERALS), wght=400, preload=True),
    Face("alegreya-sc-medium", "AlegreyaSC-Medium.ttf", "Alegreya SC", 500, "normal", "text",
         TEXT, tuple(NUMERALS), preload=True),
    Face("courier-prime-regular", "CourierPrime-Regular.ttf", "Courier Prime", 400, "normal",
         "code", CODE),
    Face("courier-prime-bold", "CourierPrime-Bold.ttf", "Courier Prime", 700, "normal", "code",
         CODE),
    Face("courier-prime-italic", "CourierPrime-Italic.ttf", "Courier Prime", 400, "italic",
         "code", CODE),
    # A math font has few features, and MathML needs its script-style forms (ssty).
    Face("libertinus-math", "LibertinusMath-Regular.otf", "Libertinus Math", 400, "normal",
         "math", MATH, ("*",)),
]


def subset_face(face, source_dir):
    """Returns the face as WOFF2 bytes and the code points it covers."""
    # Keep the source's head.modified; by default fontTools stamps the time of saving.
    font = TTFont(source_dir / face.source, recalcTimestamp=False)
    if face.wght is not None:
        font = instancer.instantiateVariableFont(font, {"wght": face.wght})
    options = subset.Options()
    options.hinting = False  # browsers on every platform render these sizes well unhinted
    options.desubroutinize = True  # smaller after WOFF2 compression (CFF fonts)
    options.name_IDs = [0, 1, 2, 3, 4, 5, 6, 13, 14]  # keep the copyright and license
    if face.features == ("*",):
        options.layout_features = ["*"]
    else:
        options.layout_features = options.layout_features + list(face.features)
    subsetter = subset.Subsetter(options)
    subsetter.populate(unicodes=subset.parse_unicodes(face.unicodes))
    subsetter.subset(font)
    font.flavor = "woff2"
    out = io.BytesIO()
    font.save(out)
    return out.getvalue(), sorted(font.getBestCmap())


def ranges(codepoints):
    """Collapses sorted code points into [first, last] pairs."""
    pairs = []
    for cp in codepoints:
        if pairs and cp == pairs[-1][1] + 1:
            pairs[-1][1] = cp
        else:
            pairs.append([cp, cp])
    return pairs


def css_rule(face, url):
    """Returns the @font-face rule for face, served from url."""
    return (f'@font-face {{ font-family: "{face.family}"; src: url({url}) format("woff2"); '
            f"font-weight: {face.weight}; font-style: {face.style}; font-display: swap; }}")


def main(argv):
    if len(argv) != 5:
        print(__doc__.split("\n\n")[1], file=sys.stderr)
        return 2
    source_dir, out_dir, css_path, manifest_path = map(pathlib.Path, argv[1:])
    logging.getLogger("fontTools").setLevel(logging.ERROR)  # "dropped table" notes

    built = []
    for face in FACES:
        data, codepoints = subset_face(face, source_dir)
        file = f"{face.name}.{hashlib.sha256(data).hexdigest()[:10]}.woff2"
        built.append((face, file, data, codepoints))

    totals = {}
    for face, _, data, _ in built:
        totals[face.group] = totals.get(face.group, 0) + len(data)
    over = [f"{g} faces are {totals[g]:,} bytes, over the {b:,}-byte budget"
            for g, b in BUDGETS.items() if totals.get(g, 0) > b]
    if over:
        print("subset_fonts: " + "; ".join(over), file=sys.stderr)
        return 1

    out_dir.mkdir(parents=True, exist_ok=True)
    keep = {file for _, file, _, _ in built}
    for old in out_dir.glob("*.woff2"):
        if old.name not in keep:
            old.unlink()
    for _, file, data, _ in built:
        (out_dir / file).write_bytes(data)

    url_dir = "/" + out_dir.name
    css_path.write_text(
        "/* Generated by site/tool/subset_fonts.py (make fonts); do not edit. */\n"
        + "".join(css_rule(face, f"{url_dir}/{file}") + "\n" for face, file, _, _ in built),
        encoding="utf-8")

    manifest = {
        "preload": [f"{url_dir}/{file}" for face, file, _, _ in built if face.preload],
        "faces": [{"file": f"{url_dir}/{file}", "family": face.family, "weight": face.weight,
                   "style": face.style, "group": face.group, "bytes": len(data),
                   "codepoints": ranges(codepoints)}
                  for face, file, data, codepoints in built],
    }
    manifest_path.parent.mkdir(parents=True, exist_ok=True)
    manifest_path.write_text(json.dumps(manifest, indent=1) + "\n", encoding="utf-8")

    summary = ", ".join(f"{g} {t / 1000:.0f} kB" for g, t in sorted(totals.items()))
    print(f"subset_fonts: wrote {len(built)} faces to {out_dir} ({summary}) and {css_path}")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
