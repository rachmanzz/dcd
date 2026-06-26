# New Features Summary

## 1. `<set:flags>` Tag for Combined Inline Formatting

### Overview
New syntax to apply multiple formatting styles simultaneously, replacing the need for nested tags.

### Syntax
```html
<set:flags>text</set:flags>
```

Where `flags` is a pipe-separated list of formatting options.

### Available Flags
- `b` - Bold
- `i` - Italic  
- `u` - Underline
- `code` - Monospace font

### Examples

**Single flag:**
```html
<p><set:b>Bold text</set:b></p>
<p><set:i>Italic text</set:i></p>
```

**Multiple flags:**
```html
<p><set:b|i>Bold and Italic</set:b|i></p>
<p><set:b|u>Bold and Underline</set:b|u></p>
<p><set:i|u>Italic and Underline</set:i|u></p>
<p><set:b|i|u>All three combined</set:b|i|u></p>
```

**With code:**
```html
<p><set:b|code>Bold monospace</set:b|code></p>
<p><set:i|code>Italic monospace</set:i|code></p>
```

### Usage Contexts
Works in all inline contexts:
- Inside `<p>` paragraphs
- Inside `<li>` list items  
- Inside `<col>` table cells

### Closing Tag
Can use either:
- `</set:flags>` (matching opening tag)
- `</set>` (simplified)

Both are equivalent.

### Backward Compatibility
✅ **Fully backward compatible**

Old syntax still works:
```html
<p><b>Bold</b> and <i>Italic</i></p>
```

Can mix old and new syntax:
```html
<p><b>old bold</b> with <set:b|i>new combined</set:b|i></p>
```

---

## 2. Global Property Renames (Breaking Change)

### Changed Properties

| Old Name | New Name | Scope |
|----------|----------|-------|
| `font-color` | `color` | All sections, inline attributes |
| `shading` | `bg` | Row/col attributes, table styles |

### Migration

**Before:**
```ini
[style]
font-color=#000000

[table-style header]
shading=#4472C4
font-color=#ffffff
```

**After:**
```ini
[style]
color=#000000

[style:table header]
bg=#4472C4
color=#ffffff
```

**Automated migration:**
```bash
sed -i 's/font-color=/color=/g; s/shading=/bg=/g; s/\[table-style /[style:table /g' file.dcd
```

---

## 3. Table Style Section Format Change (Breaking Change)

### Changed Format

**Before:**
```ini
[table-style header]
bg=#4472C4
color=#ffffff
```

**After:**
```ini
[style:table header]
bg=#4472C4
color=#ffffff
```

Consistent with heading styles: `[style:heading-1]`, `[style:heading-2]`, etc.

---

## 4. `style.first` for Loop Rows

Apply a named style to the first row/item only in a loop.

### Syntax
```html
<loop:row style.first=header x from items>
  <col>{{x.field}}</col>
</loop:row>
```

Position flexible:
```html
<loop:row x from items style.first=header>
  <col>{{x.field}}</col>
</loop:row>
```

### Supported Loop Types
- `<loop:row>` - Table rows
- `<loop:ol>` - Ordered lists
- `<loop:ul>` - Unordered lists

### Behavior
- **First iteration:** Gets the specified style
- **Remaining iterations:** Plain tags without style

### Example

```html
<table border>
  <loop:row style.first=header x from data>
    <col>{{x.name}}</col>
    <col>{{x.value}}</col>
  </loop:row>
</table>
```

With data:
```json
{
  "data": [
    {"name": "Name", "value": "Value"},
    {"name": "Item 1", "value": "100"},
    {"name": "Item 2", "value": "200"}
  ]
}
```

Result:
- First row has `style=header` (styled as header)
- Rows 2-3 are plain (data rows)

---

## 5. Dynamic `style={{var}}` for Rows

Resolve style names from variables at compile time.

### Syntax
```html
<row style={{styleVar}}>
  <col>Data</col>
</row>
```

### Scope
- ✅ Works with `<row>` and `<li>`
- ❌ Not supported inside `<loop>` templates (use `style.first` instead)

### Example

```ini
[style:table highlight]
bg=#ffff00
font-weight=bold

[section 0]
var=data
keys=rowStyle

--- BODY ---
<table border>
  <row style={{data.rowStyle}}>
    <col>Highlighted row</col>
  </row>
</table>
```

With data:
```json
{
  "data": {
    "rowStyle": "highlight"
  }
}
```

---

## 6. Loop Source Names with Dots

Fixed: Loop sources can now use dot notation for nested data.

### Before (didn't work)
```html
<loop:row x from invoice.items>
```

### After (works)
```html
<loop:row x from invoice.items>
  <col>{{x.id}}</col>
  <col>{{x.name}}</col>
</loop:row>
```

With data:
```json
{
  "invoice": {
    "items": [
      {"id": 1, "name": "Item A"},
      {"id": 2, "name": "Item B"}
    ]
  }
}
```

**Implementation:** Updated loop regex from `(\w+)` to `([\w.]+)` to allow dots.

---

## 7. Inline Formatting in Lists and Tables

Inline tags (`<b>`, `<i>`, `<u>`, `<code>`, `<set:>`) now work in:
- ✅ Paragraphs `<p>`
- ✅ List items `<li>`
- ✅ Table cells `<col>`

### Example

```html
<ul>
  <li><b>Bold item</b></li>
  <li><set:i|u>Italic underline item</set:i|u></li>
</ul>

<table border>
  <row>
    <col><set:b|u>Header</set:b|u></col>
  </row>
  <row>
    <col><i>Data</i></col>
  </row>
</table>
```

**Implementation:** Changed `ListItem` and `TableCell` from `Text string` to `Runs []TextRun`.

---

## 8. Documentation: `[title]` Section

New skill file documenting the `[title]` metadata section.

### Properties Supported

```ini
[title]
title=Document Title
subject=Document Subject
author=Author Name
```

### Built-in Variable

Use `{{title}}` in headers, footers, and body:

```ini
[title]
title=My Report

[header]
left={{title}}
```

### Documentation Location

`.agents/skills/document-metadata.md` - Complete guide with examples

---

---

## 9. Header/Footer `justify_between` (v0.2.1)

Evenly-spaced columns using OOXML tab stops.

### Syntax

2 or 3 comma-separated items. Use `\,` for literal comma.

```ini
[header]
justify_between={{title}}, {{page}} / {{total}}

[footer]
justify_between=Dept. A\, B\, and C, {{date}}, Page {{page}}
```

### Behavior

| Items | Tab Stops |
|---|---|
| 2 | Left + right |
| 3 | Left + center + right |

Tab positions auto-calculated from page width and margins.

### Props Support

`font-family`, `font-size`, `color`, `border`, `margin`, `first-page` all work with `justify_between`.

### Related Fixes

- `{{page}}` now renders only the PAGE field (not combined with NUMPAGES)
- XML escaping for `{{title}}` and all text content
- Segment-based OOXML generation (proper sibling elements, not nested in `<w:t>`)

---

## Summary of Files Changed

### Code (5 files)
- `parse/parse.go` - Property normalization
- `render/style.go` - normalizePropertyKey()
- `render/body.go` - <set:> tag, style.first, loop regex
- `render/compiler.go` - applyTableStyles(), resolveRowStyles()
- `render/types.go` - ListItem.Runs, TableCell.Runs

### Documentation (12 files)
- `.agents/skills/document-body.md` - <set:> examples
- `.agents/skills/document-table.md` - Updated examples, style.first
- `.agents/skills/document-metadata.md` - NEW: [title] documentation
- `.agents/skills/document-heading.md` - Property renames
- `.agents/skills/document-image.md` - Property renames
- `.agents/skills/document-style.md` - Property renames
- `.agents/skills/header-footer.md` - Property renames
- `docs/style.md` - Property renames, section format
- `docs/tags.md` - <set:> tag, property renames
- `docs/library.md` - Updated type definitions
- `docs/format.md` - Updated examples
- `docs/overview.md` - Updated quick start

### Examples (5 files)
- `docs/examples/simple.dcd` - Property renames
- `docs/examples/features.dcd` - Property renames
- `docs/examples/report.dcd` - Property renames, section format
- `docs/examples/invoice.dcd` - Property renames, section format
- `docs/examples/set-tag-demo.dcd` - NEW: <set:> demo
- `docs/examples/inline-test.dcd` - Inline formatting test

---

## Breaking Changes

⚠️ **Breaking changes require migration of existing .dcd files**

1. Property renames: `font-color` → `color`, `shading` → `bg`
2. Section format: `[table-style]` → `[style:table]`

All other changes are backward compatible or additive.

---

## Testing Summary

✅ All unit features tested
✅ All integration scenarios tested  
✅ DOCX output verified
✅ PDF output verified
✅ All existing examples compile
✅ Backward compatibility verified
✅ Build passes: `go build ./...`
✅ Vet passes: `go vet ./...`

