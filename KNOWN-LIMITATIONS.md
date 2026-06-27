# Known Limitations

This document maps each feature to its implementation layer in the `gomutex/godocx v0.1.5` library:

| Layer | What it means | Used for |
|-------|---------------|----------|
| **Library API** | Direct method call on `docx.*` types | Init, paragraphs, runs, tables, images, numbering apply, save |
| **CT struct field** | Access via `GetCT()` + set `ctypes.*` fields | Formatting properties not exposed by library methods |
| **Beyond scope** | Requires library fork, XML injection, or struct access the library doesn't expose | Table span/width, outline numbering defs, border width |

---

## Implemented (via Library API + CT)

### Inline Formatting

| Tag | CT Field | Method |
|-----|----------|--------|
| `<b>` bold | `RunProperty.Bold` | `r.Bold(true)` |
| `<i>` italic | `RunProperty.Italic` | `r.Italic(true)` |
| `<u>` underline (with variants) | `RunProperty.Underline` | `r.Underline(stype)` |
| `<s>` strikethrough | `RunProperty.Strike` | CT field set |
| `<code>` monospace | `RunProperty.Fonts` | CT field set |
| `<mark>` highlight | `RunProperty.Highlight` | CT field set |
| `<sub>` subscript | `RunProperty.VertAlign` | CT field set |
| `<sup>` superscript | `RunProperty.VertAlign` | CT field set |
| `<set:flags>` combined | Any of the above | CT field set (parser handles `\|` syntax) |

### Style Properties

| Property | CT Field | Notes |
|----------|----------|-------|
| `bold`, `italic`, `strike`, `underline` | `RunProperty.*` | Applied per-run |
| `caps=true` | `RunProperty.Caps` | All-caps |
| `small-caps=true` | `RunProperty.SmallCaps` | Small capitals |
| `letter-spacing=N` | `RunProperty.Spacing` | Character spacing in twips |
| `underline=double/dotted/dash/wavy` | `stypes.Underline*` | Underline variants |
| `keep-next=true` | `ParagraphProp.KeepNext` | Keep with next paragraph |
| `keep-lines=true` | `ParagraphProp.KeepLines` | Keep lines together |

### Structural

| Feature | API/CT | Notes |
|---------|--------|-------|
| Headings 1-6 | `d.root.AddHeading(text, level)` | Library applies heading styles |
| Paragraphs | `d.root.AddEmptyParagraph()` / `AddParagraph()` | + CT for alignment, spacing, shading |
| Lists (≤3 levels) | `p.Style(name)` + `p.Numbering(id, level)` | Uses template's numbering defs |
| Tables | `d.root.AddTable()` / row / cell | Border/shading/font via CT |
| Image | `d.root.AddPicture(path, w, h)` | Full drawingML pipeline |
| Hyperlink | CT Hyperlink + relationship | Manual CT construction |
| Page break | `d.root.AddPageBreak()` | |
| Horizontal rule | CT `pPr.Border.Bottom` | **No thickness control** (see below) |
| Page size/margins | `SectPr.PageSize` / `PageMargin` | Full CT struct |
| Header/footer (single) | `FileMap` raw XML | Content built as XML string |
| Header/footer `mirror` | `HdrFtrEven` type | Single reference only |
| Header/footer `first-page` | `SectPr.TitlePg` | CT field |
| Core metadata | `FileMap` raw XML | Library has no metadata struct |

---

## Not Possible (Beyond Library API/CT Scope)

| Feature | Root Cause | What Would Be Needed |
|---------|-----------|---------------------|
| Table colspan/rowspan | `docx.Cell.ct` is unexported — can't set `CellProperty.GridSpan` / `VMerge` | Library patch: export `Cell.GetCT()` or add span methods |
| Table column width | Same `Cell.ct` issue + `w:tblGrid` needs raw XML in table properties | Library patch: export `Table.GetCT()` + `Cell.GetCT()` |
| Outline numbering (1, 1.1, 1.1.1) | `w:abstractNum`/`w:num`/`w:lvl` have no Go structs. Library provides `Paragraph.Numbering(id, level)` to _apply_ numbering, but no API to _create_ numbering definitions. | Library patch: numbering definition builder OR raw XML injection into `numbering.xml` |
| `<hr thick=N>` | `ctypes.Border` struct missing `Sz` (border width) field. `MarshalXML` doesn't emit `w:sz`. | Library patch: add `Sz *uint64` to `Border` struct |
| Header/footer odd+even separately | `SectionProp` has only one `HeaderReference`/`FooterReference` slot. Can only use `HdrFtrEven` **instead of** `HdrFtrDefault`, not both. | Library patch: support multiple references |

### Why Not Raw XML Injection?

For **outline numbering** and **table features**, raw XML injection into `FileMap` (the approach used for header/footer content and core metadata) is theoretically possible but was deliberately excluded:

- **Numbering XML**: Requires injection into `word/numbering.xml` before `</w:numbering>`. Feasible via string replace, but fragile on library upgrades and adds non-obvious coupling.
- **Table Cell CT**: The cell's `tcPr` lives inside document.xml's paragraph hierarchy — not a separate file. Can't inject at FileMap level; would need to navigate the live XML tree inside `RootDoc.Document.Body.Children`, which is possible but tightly coupled to library internals.

Decision: these features are deferred until the library provides proper API support.

---

## Won't Fix — Alternative Exists

| Feature | Alternative |
|---------|-------------|
| Nested inline tags (`<b><i>text</i></b>`) | Use `<set:b\|i>text</set:b\|i>` |
