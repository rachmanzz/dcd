# Refactoring Plan: SRP & KISS Compliance

## Problem Summary

| File | Lines | Issue |
|------|------:|-------|
| `render/docx.go` | 1,584 | 5 responsibility domains, run-styling duplicated 3x (~90 lines x 3) |
| `render/body.go` | 1,049 | 5 responsibility domains, tokenizer 148 lines, loop expander 111 lines |
| `render/compiler.go` | 369 | `validateSectionProps` 139 lines with 6+ interleaved rules |
| `render/style.go` | 243 | 4 unrelated concern groups |
| `parse/parse.go` + `render/style.go` | - | `normalizePropertyKey` duplicated across packages |

---

## Phase 1: Extract Shared Helpers (DRY Fix)

**Goal:** Eliminate code duplication before splitting files.

### 1a. Extract `applyRunProps` helper

The run-property application block (~90 lines: bold, italic, underline, code font, strike, mark, sub, sup) is **duplicated 3 times** in `render/docx.go`:
- Lines 299-380 (`AddParagraph`)
- Lines 686-759 (`addListAtDepth`)
- Lines 796-878 (`AddTable`)

**Action:** Create `render/run_props.go` with a shared helper:

```go
// render/run_props.go (~60 lines)
func applyRunProps(ctRun *ctypes.Run, run TextRun, defaultStyle, attrs map[string]string) {
    // Apply font-color, font-size, font-family from defaultStyle
    // Apply font-color, font-size from attrs
    // Apply code font, strike, bold, italic, underline, mark, sub, sup
}
```

Replace the 3 duplicated blocks with calls to `applyRunProps(...)`.

**Estimated:** -180 lines duplication, +60 lines helper = **net -120 lines**

### 1b. Consolidate `normalizePropertyKey`

Identical function defined in:
- `parse/parse.go:12`
- `render/style.go:28`

**Action:** Create `internal/property/property.go` with canonical definition, import from both packages. Or keep in `render/` and create a thin wrapper in `parse/`.

---

## Phase 2: Split `render/docx.go` (1,584 -> ~5 files)

**Goal:** Each file has one clear responsibility.

### 2a. `render/run_props.go` (from Phase 1a)
- `applyRunProps` helper
- ~60 lines

### 2b. `render/docx_paragraph.go` (~350 lines)
Migrate:
- `AddHeading` (lines 104-226)
- `applyIndent` (lines 228-256)
- `AddParagraph` (lines 258-383, ~40 lines after extract)
- `AddWrappedParagraph` (lines 934-1052)
- `AddLineBreak` (lines 385-394)
- `AddPageBreak` (lines 443-451)
- `AddSectionBreak` (lines 453-473)
- `AddHorizontalRule` (lines 396-441)

### 2c. `render/docx_list.go` (~80 lines)
Migrate:
- `AddList` (lines 623-625)
- `addListAtDepth` (lines 627-769, ~70 lines after extract)

### 2d. `render/docx_table.go` (~90 lines)
Migrate:
- `AddTable` (lines 771-932, ~80 lines after extract)

### 2e. `render/docx_header_footer.go` (~270 lines)
Migrate:
- `SetHeader` (lines 1105-1155)
- `SetFooter` (lines 1157-1204)
- `hdrSegment`, `hdrPart`, `hdrFooterCfg` structs (lines 1206-1226)
- `splitComma` (lines 1228-1246)
- `hdrTabPositions` (lines 1248-1275)
- `collectHeaderFooterParts` (lines 1277-1317)
- `parseHdrText` (lines 1319-1353)
- `xmlEscape` (lines 1355-1360)
- `buildHdrXML` (lines 1362-1413)
- `buildFtrXML` (lines 1415-1466)
- `hasRunProps` (lines 1468-1470)
- `runPropsXML` (lines 1472-1489)

### 2f. `render/docx.go` (remainder, ~400 lines)
- `DocxRenderer` struct + `NewDocxRenderer` + `init`
- `SetMetadata`, `SetDefaultStyle`, `SetHeadingStyle`, `SetTableStyle`
- `SetPageStyle`
- `AddImage`, `AddHyperlink`
- `Save`
- `intPtr`, `numFmtBaseID`, `numFmtName`, `injectNumFmts`

---

## Phase 3: Split `render/body.go` (1,049 -> ~4 files)

### 3a. `render/inline.go` (~320 lines) -- Inline Markup Tokenizer
Migrate:
- Regex variables: `bRe`, `iRe`, `uRe`, `sRe`, `codeRe`, `markRe`, `subRe`, `supRe`, `setRe`, `brRe`, `tabRe`, `hrRe`, `pageBreakRe`, `linkRe`
- `inlinePart` struct
- `splitInline` (lines 857-1004, 148 lines)
- `inlineToRuns` (lines 789-855, 67 lines)
- `validateTagBalance` (lines 721-787, 67 lines)

### 3b. `render/loops.go` (~200 lines) -- Template Loop Expansion
Migrate:
- Regex variables: `loopRe`, `loopSourceRe`, `objectVarRe`
- `expandLoops` (lines 261-371, 111 lines)
- `expandLoopTemplate` (lines 373-408)
- `resolveField` (lines 410-417)
- `extractLiFromTemplate` (lines 424-446)
- `mergeLiStyleFirst` (lines 450-459)

### 3c. `render/collector.go` (~180 lines) -- Structural Element Collection
Migrate:
- `collectListItems` (lines 461-480)
- `collectLi` (lines 482-552)
- `collectTableRows` (lines 554-575)
- `collectRowCells` (lines 577-633)
- Regex variables used by collector: `hRe`, `pRe`, `wRe`, `imgRe`

### 3d. `render/body.go` (remainder, ~250 lines) -- Line Dispatch
- `renderBody` (lines 43-115)
- `collectPTag` (lines 117-152)
- `tagAttrs` (lines 154-160)
- `collectWTag` (lines 162-205)
- `renderWrappedContent` (lines 207-256)
- `parseLine` (lines 635-695)
- `renderParagraph` (lines 697-703)
- Utilities: `mergeAttrs`, `parseListType`, `parseAttrs`

---

## Phase 4: Extract Validation from `render/compiler.go`

### 4a. `render/validate.go` (~180 lines)
Migrate and refactor:
- `validateSectionProps` (lines 174-313) -- decompose into sub-validators:
  - `validateVarDecl` -- checks 1, 2, 3, 4, 6
  - `validateKeyLimits` -- check 5
  - `validateFormatKeys` -- check 5a
- `parseVarDecl` (lines 156-170)
- `varEntry` struct (lines 151-154)

### 4b. `render/compiler.go` (remainder, ~190 lines)
- Orchestration: `New`, `Run`, `collectSection`, `renderSection`
- Style application: `applyHeadingStyles`, `applyTableStyles`, `applyMetadata`, `applyHeaderFooter`
- Format/builtin: `applyFormats`, `resolveBuiltins`, `resolveSectionBuiltins`, `resolveRowStyles`

---

## Phase 5: Reorganize `render/style.go` (243 -> ~3 files)

### 5a. `render/page.go` (~75 lines) -- Page Geometry
- `parsePageSize`
- `unitScale`
- `computeMargins`

### 5b. `render/format.go` (~50 lines) -- Date/Value Formatting
- `specConv`, `specRe`, `fmtRe` variables
- `parseFormats`
- `convertFormat`
- `applyFormat`

### 5c. `render/style.go` (remainder, ~70 lines) -- Style Helpers
- `normalizePropertyKey` (or move to shared package)
- `chooseAttr`
- `atoi`, `atof` (consider replacing with `strconv`)
- `underlineFromString`

---

## Expected Result

| Source File | Before | After | New Files |
|-------------|-------:|------:|-----------|
| `render/docx.go` | 1,584 | ~400 | +4 files (paragraph, list, table, header_footer) |
| `render/body.go` | 1,049 | ~250 | +3 files (inline, loops, collector) |
| `render/compiler.go` | 369 | ~190 | +1 file (validate) |
| `render/style.go` | 243 | ~70 | +2 files (page, format) |
| **Total** | **3,245** | **~910** | **+11 new files** |

Plus: ~120 lines reduction from `applyRunProps` dedup.

---

## Execution Order & Risk Assessment

| Phase | Risk | Mitigation |
|-------|------|------------|
| 1 (DRY helpers) | Medium -- changes rendering logic | Manual test: generate docx, compare output |
| 2 (split docx.go) | Low -- move code only | `go build ./...` + `go vet ./...` |
| 3 (split body.go) | Low -- move code only | `go build ./...` + `go vet ./...` |
| 4 (extract validate) | Low -- move + small refactor | Test with existing .dcd files |
| 5 (reorganize style) | Low -- move code only | `go build ./...` + `go vet ./...` |

**After each phase:** Run `go build ./...`, `go vet ./...`, and `./test-examples.sh`.

---

## Verification Checklist

- [ ] `go build ./...` passes
- [ ] `go vet ./...` passes
- [ ] `./test-examples.sh` passes
- [ ] Generated .docx files are identical before/after each phase
- [ ] No circular imports introduced
- [ ] All exported functions/types remain accessible
