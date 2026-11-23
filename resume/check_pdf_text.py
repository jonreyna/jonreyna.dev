"""Checks that each résumé PDF reads back as the text it was built from.

Usage: python3 -I check_pdf_text.py DIR

Applicant-tracking systems don't see the PDF; they see the text a tool like pdftotext
pulls out of it. This script extracts every DIR/resume.<variant>.pdf with pdftotext
(from Poppler) and compares it with resume.<variant>.json beside it, the same render's
structured copy. It fails when:

- a passage is missing or garbled: the name, label and email; the summary; each
  section heading; every role's company, and every highlight word for word; the skill
  groups and their items; the note on earlier work; every school;
- the passages come out of order (a parser would misread the structure);
- the text has ligature characters (an "ffi" glyph extracts as one character that
  parsers drop) or replacement characters (a glyph with no Unicode mapping).

Whitespace is collapsed before comparing, so line and page breaks don't matter, but a
space inside a word or a missing one between words does. Only the standard library is
used.
"""

import json
import pathlib
import re
import shutil
import subprocess
import sys

# Section headings as jrresume.cls prints them. The summary and the note on earlier
# work have none.
HEADINGS = {"skills": "Skills", "experience": "Experience", "education": "Education"}

# Characters a parser can't use: Unicode's ligature block and the replacement character.
BAD_CHARS = re.compile("[ﬀ-ﬆ�]")

# pdfLaTeX's hyphen may map to U+2010 (hyphen) or U+2011 (non-breaking hyphen); both
# read as an ASCII hyphen. A no-break space reads as a space.
TRANSLATE = str.maketrans({"‐": "-", "‑": "-", " ": " "})


def normalize(text):
    """Collapses whitespace and the hyphen and space variants."""
    return " ".join(text.translate(TRANSLATE).split())


def plain(field):
    """Returns the plain half of one of vitae's text fields ({"plain": …, "html": …})."""
    return field.get("plain", "") if isinstance(field, dict) else str(field or "")


def passages(doc):
    """Lists (what, text) pairs in the order the PDF prints them."""
    basics = doc.get("basics", {})
    out = [("name", basics.get("name")), ("label", basics.get("label")), ("email", basics.get("email"))]
    for section in doc.get("sections", []):
        if section in HEADINGS:
            out.append((f"{section} heading", HEADINGS[section]))
        if section == "summary":
            out.append(("summary", plain(doc.get("summary"))))
        elif section == "skills":
            for group in doc.get("skills", []):
                out.append(("skill group", group["group"]))
                out.extend(("skill", item) for item in group.get("items", []))
        elif section == "experience":
            for role in doc.get("experience", []):
                out.append(("company", role["org"]))
                for h in role.get("highlights", []):
                    if h.get("lead"):
                        out.append((f"{role['org']} lead", h["lead"]))
                    out.append((f"{role['org']} highlight", plain(h["text"])))
        elif section == "earlier":
            out.append(("earlier", plain(doc.get("earlier"))))
        elif section == "education":
            out.extend(("school", school["org"]) for school in doc.get("education", []))
    return [(what, normalize(text)) for what, text in out if text]


def check(doc, extracted):
    """Returns the problems found comparing doc (the JSON) with the extracted text."""
    problems = []
    for char in sorted(set(BAD_CHARS.findall(extracted))):
        problems.append(f"has character U+{ord(char):04X}, which parsers can't read")
    text = normalize(extracted)
    pos = 0
    for what, expected in passages(doc):
        found = text.find(expected, pos)
        if found >= 0:
            pos = found + len(expected)
            continue
        if expected in text:
            problems.append(f"{what} is out of order: {expected[:70]!r}")
        else:
            problems.append(f"{what} is missing or garbled: {departure(expected, text, pos)}")
    return problems, len(passages(doc))


def departure(expected, text, pos):
    """Describes where the extracted text first departs from an expected passage."""
    near = text.find(expected[:24], pos)
    if near < 0:
        return f"{expected[:70]!r} not found"
    got = text[near : near + len(expected)]
    i = next((i for i, (a, b) in enumerate(zip(expected, got)) if a != b), len(got))
    start = max(0, i - 30)
    return f"\n  expected …{expected[start : i + 30]!r}\n  got      …{got[start : i + 30]!r}"


def extract(pdf):
    """Returns the PDF's text as pdftotext reads it."""
    result = subprocess.run(
        ["pdftotext", "-enc", "UTF-8", "-nopgbrk", str(pdf), "-"],
        check=True,
        capture_output=True,
        text=True,
    )
    return result.stdout


def main():
    if len(sys.argv) != 2:
        sys.exit("usage: check_pdf_text.py DIR")
    if shutil.which("pdftotext") is None:
        sys.exit("check_pdf_text.py: pdftotext is missing; install Poppler (apt install poppler-utils)")
    folder = pathlib.Path(sys.argv[1])
    pdfs = sorted(folder.glob("resume.*.pdf"))
    if not pdfs:
        sys.exit(f"check_pdf_text.py: no resume.*.pdf in {folder}; run make resume-pdf first")
    failed = False
    for pdf in pdfs:
        doc = json.loads(pdf.with_suffix(".json").read_text(encoding="utf-8"))
        problems, count = check(doc, extract(pdf))
        if problems:
            failed = True
            print(f"{pdf}: {len(problems)} problems")
            for problem in problems:
                print(f"  {problem}")
        else:
            print(f"checked {pdf}: {count} passages extracted intact and in order")
    if failed:
        sys.exit(1)


if __name__ == "__main__":
    main()
