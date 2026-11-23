"""Builds the Word reference document that styles the résumé's .docx files.

Usage: python3 -I reference.py PANDOC_DEFAULT.docx OUTPUT.docx

Pandoc styles a .docx after a reference document. Rather than commit one (it is a
binary, and binaries never go in git), make asks Pandoc for its default
(`pandoc --print-default-data-file reference.docx`) and this script restyles it to
follow the PDF's design (resume/latex/jrresume.cls) as far as Word allows:

- Calibri, which every copy of Word has, instead of Lato, which most don't;
- 10.5pt near-black text, one deep blue for the name and section heads;
- section heads with a hairline rule under them, as in the PDF;
- US Letter with the PDF's margins (0.7in sides, 0.55in top and bottom);
- tight spacing, so the full résumé stays about two pages.

Only the standard library is used. The output is deterministic: entries keep their
order and get a fixed timestamp.
"""

import re
import sys
import xml.etree.ElementTree as ET
import zipfile

W = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"
FONT = "Calibri"
INK = "222222"
ACCENT = "1F4E79"
RULE = "A7B6C6"

# Sizes are in half-points, spacing in twentieths of a point.
STYLES = {
    # style id: (run properties, paragraph properties)
    "BodyText": ({}, {"spacing": {"before": "0", "after": "80"}}),
    "FirstParagraph": ({}, {"spacing": {"before": "0", "after": "80"}}),
    "Compact": ({}, {"spacing": {"before": "0", "after": "40"}}),
    "Heading1": (
        {"sz": "52", "b": True, "color": ACCENT},
        {"spacing": {"before": "0", "after": "40"}},
    ),
    "Heading2": (
        {"sz": "26", "b": True, "color": ACCENT},
        {"spacing": {"before": "220", "after": "80"}, "border": True},
    ),
    "Heading3": (
        {"sz": "23", "b": True, "color": INK},
        {"spacing": {"before": "160", "after": "0"}},
    ),
    "Hyperlink": ({"color": ACCENT}, {}),
}


def q(tag):
    """Returns the qualified name of a WordprocessingML tag."""
    return f"{{{W}}}{tag}"


def child(parent, tag):
    """Returns parent's child with the given tag, adding it if missing."""
    found = parent.find(q(tag))
    if found is None:
        found = ET.SubElement(parent, q(tag))
    return found


def set_fonts(rpr, color=INK, size=None):
    """Sets explicit fonts (replacing theme fonts), color and size on run properties."""
    fonts = child(rpr, "rFonts")
    fonts.attrib.clear()
    for side in ("ascii", "hAnsi", "cs", "eastAsia"):
        fonts.set(q(side), FONT)
    col = child(rpr, "color")
    col.attrib.clear()
    col.set(q("val"), color)
    if size:
        for tag in ("sz", "szCs"):
            child(rpr, tag).set(q("val"), size)


def restyle(styles_xml):
    """Returns styles.xml restyled for the résumé."""
    ET.register_namespace("w", W)
    ET.register_namespace(
        "r", "http://schemas.openxmlformats.org/officeDocument/2006/relationships"
    )
    root = ET.fromstring(styles_xml)

    defaults = root.find(q("docDefaults"))
    set_fonts(child(child(defaults, "rPrDefault"), "rPr"), size="21")
    spacing = child(child(child(defaults, "pPrDefault"), "pPr"), "spacing")
    spacing.attrib.clear()
    spacing.set(q("after"), "0")

    for style in root.findall(q("style")):
        sid = style.get(q("styleId"))
        if sid not in STYLES:
            continue
        run, para = STYLES[sid]
        rpr = child(style, "rPr")
        if sid == "Hyperlink":
            child(rpr, "color").attrib.clear()
            child(rpr, "color").set(q("val"), run["color"])
            continue
        set_fonts(rpr, color=run.get("color", INK), size=run.get("sz"))
        if run.get("b"):
            child(rpr, "b")
            child(rpr, "bCs")
        if para:
            ppr = child(style, "pPr")
            sp = child(ppr, "spacing")
            sp.attrib.clear()
            for key, value in para["spacing"].items():
                sp.set(q(key), value)
            if para.get("border"):
                bottom = child(child(ppr, "pBdr"), "bottom")
                for key, value in (("val", "single"), ("sz", "4"), ("space", "1"), ("color", RULE)):
                    bottom.set(q(key), value)
        # Word expects pPr before rPr inside a style.
        if style.find(q("pPr")) is not None:
            ppr = style.find(q("pPr"))
            style.remove(ppr)
            style.remove(rpr)
            style.append(ppr)
            style.append(rpr)

    return ET.tostring(root, xml_declaration=True, encoding="UTF-8")


def page_setup(document_xml):
    """Adds US Letter and the PDF's margins to the document's section properties."""
    text = document_xml.decode("utf-8")
    setup = (
        '<w:pgSz w:w="12240" w:h="15840" />'
        '<w:pgMar w:top="792" w:right="1008" w:bottom="792" w:left="1008"'
        ' w:header="0" w:footer="0" w:gutter="0" />'
    )
    text, count = re.subn(r"<w:sectPr>", "<w:sectPr>" + setup, text, count=1)
    if count != 1:
        sys.exit("reference.py: no <w:sectPr> in word/document.xml")
    return text.encode("utf-8")


def main():
    if len(sys.argv) != 3:
        sys.exit("usage: reference.py PANDOC_DEFAULT.docx OUTPUT.docx")
    source, target = sys.argv[1], sys.argv[2]
    with zipfile.ZipFile(source) as zin, zipfile.ZipFile(target, "w", zipfile.ZIP_DEFLATED) as zout:
        for item in zin.infolist():
            data = zin.read(item.filename)
            if item.filename == "word/styles.xml":
                data = restyle(data)
            elif item.filename == "word/document.xml":
                data = page_setup(data)
            info = zipfile.ZipInfo(item.filename, date_time=(1980, 1, 1, 0, 0, 0))
            info.compress_type = zipfile.ZIP_DEFLATED
            zout.writestr(info, data)
    print(f"wrote {target}")


if __name__ == "__main__":
    main()
