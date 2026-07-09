# Implementation Plan

## Completed

| Feature | Status | Files |
|---------|--------|-------|
| Inline tags: `<b>`, `<i>`, `<u>`, `<s>`, `<code>`, `<mark>`, `<sub>`, `<sup>` | ✅ | `body.go`, `docx.go` |
| Underline variants (`single`, `double`, `dotted`, `dash`, `wavy`) | ✅ | `style.go`, `docx.go`, `body.go` |
| `<set:flags>` combined formatting (closing tag must match) | ✅ | `body.go` |
| `<w:flags>` block formatting (`c`, `r`, `j`, `b`, `i`, `s`, `u`) + `underline` attr | ✅ | `docx.go` |
| Style properties: `caps`, `small-caps`, `letter-spacing`, underline variants | ✅ | `docx.go` |
| Nested lists (≤3 levels, mixed `ul`/`ol`) | ✅ | `types.go`, `body.go`, `docx.go` |
| Dead code removal (`nestedListRe`, `parseNestedListItems`) | ✅ | `body.go` |
| Mark/Sub/Sup in `addListAtDepth` | ✅ | `docx.go` |
| Code/Mark/Sub/Sup in `AddWrappedParagraph` | ✅ | `docx.go` |
| Strike in `addListAtDepth` | ✅ | `docx.go` |
| Keep-next / Keep-lines in heading | ✅ | `docx.go` |
| Header/footer mirror (`HdrFtrEven` type) | ✅ | `docx.go` |

---

## Gap Fix: Mark/Sub/Sup in AddList

**Status:** ✅ Done

**Why:** `addListAtDepth` (`docx.go:562`) handles `Code`, `Bold`, `Italic`, `Underline`, `Strike` but not `Mark`, `Sub`, `Sup`. These work in `AddParagraph` and `AddTable` but were missing in lists.

**Changes:**
- `render/docx.go` — Add `Mark`, `Sub`, `Sup` blocks in `addListAtDepth` loop (after line 638), matching the pattern from `AddParagraph` (lines 314–335)

**Pattern:**
```go
if run.Mark {
    if ctRun.Property == nil {
        ctRun.Property = &ctypes.RunProperty{}
    }
    color := "yellow"
    if run.MarkColor != "" {
        color = run.MarkColor
    }
    ctRun.Property.Highlight = ctypes.NewCTString(color)
}
if run.Sub {
    if ctRun.Property == nil {
        ctRun.Property = &ctypes.RunProperty{}
    }
    ctRun.Property.VertAlign = ctypes.NewGenSingleStrVal(stypes.VerticalAlignRunSubscript)
}
if run.Sup {
    if ctRun.Property == nil {
        ctRun.Property = &ctypes.RunProperty{}
    }
    ctRun.Property.VertAlign = ctypes.NewGenSingleStrVal(stypes.VerticalAlignRunSuperscript)
}
```

**Effort:** Trivial (copy-paste from `AddParagraph`).

---

## Gap Fix: AddWrappedParagraph missing tags

**Status:** ✅ Done

**Why:** `AddWrappedParagraph` (`docx.go:813`) handles `c`, `r`, `j`, `b`, `i`, `s`, `u` flags but not `code`, `mark`, `sub`, `sup`. Block tags `<w:code>`, `<w:mark>`, `<w:sub>`, `<w:sup>` didn't work.

**Changes:**
- `render/docx.go` — Add `case "code"`, `case "mark"`, `case "sub"`, `case "sup"` in the switch at line 838

**Pattern** (same as other cases in that switch — operate on the first run's Property):
```go
case "code":
    if run := p.GetCT().Children[0].Run; run != nil {
        if run.Property == nil {
            run.Property = &ctypes.RunProperty{}
        }
        run.Property.Fonts = &ctypes.RunFonts{Ascii: "Courier New", HAnsi: "Courier New"}
    }
```

**Note:** `mark` needs `attrs["color"]` for highlight color, `sub`/`sup` set `VertAlign`.

**Effort:** Trivial.

---

## Audit Finding: Strike in addListAtDepth

**Status:** ✅ Done

**Why:** Audit found `if run.Strike { ... }` block is missing from `addListAtDepth` (`docx.go`). All other functions (`AddParagraph`, `AddTable`, `AddWrappedParagraph`, `AddHeading`) handle Strike. This is a consistency gap.

**Changes:**
- `render/docx.go` — Add Strike block in `addListAtDepth` loop (between `Underline` and `Mark` blocks)

**Pattern:**
```go
if run.Strike {
    if ctRun.Property == nil {
        ctRun.Property = &ctypes.RunProperty{}
    }
    ctRun.Property.Strike = ctypes.OnOffFromBool(true)
}
```

**Effort:** Trivial.

---

## Phase 3: Table Colspan / Rowspan

**Status:** ❌ Cancelled — beyond library API/CT scope

**Why:** `docx.Cell.ct` is unexported. `CellProperty.GridSpan`/`VMerge` exist in the CT struct but cannot be accessed. Would require library patch to export `Cell.GetCT()`.

---

## Phase 4: Table Column Width

**Status:** ❌ Cancelled — beyond library API/CT scope

**Why:** Same `Cell.ct` issue as colspan + `w:tblGrid` needs raw table properties XML injection.

---
## Phase 5: Keep-next / Keep-lines

**Status:** ✅ Done

**Why:** Available on `ParagraphProp` (`KeepNext *OnOff`, `KeepLines *OnOff`). No parser support yet.

**Changes:**
- `render/docx.go` — In `AddHeading`, added `keep-next` and `keep-lines` style property checks after border-bottom block

**Effort:** Trivial.

---

## Phase 6: Outline Numbering

**Status:** ❌ Cancelled — beyond library API/CT scope

**Why:** Numbering definitions (`w:abstractNum`/`w:num`/`w:lvl`) have no Go structs in the library. `Paragraph.Numbering(id, level)` can _apply_ numbering but numbering definitions must exist in `numbering.xml`. Creating them requires raw XML injection or library patch.

---

## Phase 7: `<hr>` Thickness

**Status:** ❌ Library limitation — won't implement

**Why:** `ctypes.Border` struct missing `Sz` field. Library fork needed (see KNOWN-LIMITATIONS.md).

**Effort:** — (blocked by library).

---

## Phase 8: Header/Footer `mirror` Property

**Status:** ✅ Done

**Why:** `mirror=true` for different odd/even headers uses `HdrFtrEven` reference type.

**Limitation:** Library's `SectionProp` only supports one `HeaderReference`/`FooterReference`. When mirror is enabled, `HdrFtrEven` type is used instead of `HdrFtrDefault`. Full odd+even support requires library patch to support multiple header/footer references per `SectionProp`.

**Changes:**
- `render/docx.go` — In `SetHeader`/`SetFooter`, when `cfg.mirror` is true, set `HdrFtrEven` as reference type instead of `HdrFtrDefault`

**Effort:** Trivial.

---

## Summary Table

| # | Feature | Effort | Status |
|---|---------|--------|--------|
| 1 | Mark/Sub/Sup in AddList | Trivial | ✅ Done |
| 2 | Code/Mark/Sub/Sup in AddWrappedParagraph | Trivial | ✅ Done |
| 3 | Strike in addListAtDepth | Trivial | ✅ Done |
| 4 | Keep-next / Keep-lines | Trivial | ✅ Done |
| 5 | Header/footer mirror | Trivial | ✅ Done |
| 6 | Table colspan/rowspan | — | ❌ Beyond lib scope |
| 7 | Table column width | — | ❌ Beyond lib scope |
| 8 | Outline numbering | — | ❌ Beyond lib scope |
| 9 | `<hr>` thickness | — | ❌ Beyond lib scope |

---

## Sync dcdmaker Rules — Validation & Documentation Gaps

### Status

| # | Item | SKILL.md | Compiler | Priority |
|---|------|----------|----------|----------|
| 1 | Built-in vars `{{page}}`, `{{total}}` | 🔴 Missing | 🔴 Not resolved (template literal) | High |
| 2 | `<ol type=a/A/i/I>` | 🔴 Missing | ⚠️ Need check parser regex | Medium |
| 3 | Loop action before attributes | 🔴 Missing | ✅ Parser accepts both orders | Low |
| 4 | List Loop Prohibition | 🔴 Missing | ✅ Parse error anyway | Low |
| 5 | Section limits (3 var, 15 keys) | 🔴 Missing | ❌ Not enforced | Low |
| 6 | Strict Usage (var/key must be used) | 🔴 Missing | ❌ Not enforced | High |
| 7 | `<w:*>` nesting prohibition | 🔴 Missing | ⚠️ Parser accepts nested | Medium |
| 8 | Heading nesting prohibition | 🔴 Missing | ⚠️ Parser accepts nested | Medium |
| 9 | Tag balancing global | 🔴 Missing | ❌ Not validated | Low |
| 10 | Section `name=` required | 🔴 Missing | ❌ Not enforced | High |

### Planned Changes

**SKILL.md only (no compiler change):**
- #1: Add `{{page}}`, `{{total}}` to built-in vars table
- #2: Add `<ol type=a>` note or "not supported"
- #3, #4, #7, #8, #9: Add documentation rules

**Compiler enforcement:**
- #5: Section limits (warning level, not error)
- #6: Strict Usage — validate each var/key appears in body
- #10: Section `name=` required — error if missing
