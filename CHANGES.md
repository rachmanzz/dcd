# Breaking Changes - Property Rename & New Features

## Summary

This update implements a **breaking change** with no backward compatibility. All existing `.dcd` files must be updated to use the new property names and section format.

---

## 1. Property Renames (Global)

### Changed Properties

| Old Name | New Name | Applies To |
|----------|----------|------------|
| `font-color` | `color` | All styles, sections, inline attributes |
| `shading` | `bg` | All row/col/image attributes, table styles |

### Unchanged Properties

- `font-weight`
- `font-size`
- `font-family`
- `border-bottom`
- `align`

### Migration Examples

**Before:**
```ini
[style]
font-color=#000000

[table-style header]
shading=#4472C4
font-color=#ffffff
```

```html
<row shading=#f0f0f0>
  <col font-color=#000>Text</col>
</row>
```

**After:**
```ini
[style]
color=#000000

[style:table header]
bg=#4472C4
color=#ffffff
```

```html
<row bg=#f0f0f0>
  <col color=#000>Text</col>
</row>
```

---

## 2. Table Style Section Format

### Changed Format

**Before:**
```ini
[table-style header]
shading=#4472C4
font-color=#ffffff
```

**After:**
```ini
[style:table header]
bg=#4472C4
color=#ffffff
```

All `[table-style name]` sections must be renamed to `[style:table name]`.

---

## 3. New Feature: `style.first` for Loops

Apply a named style to the **first item only** in a loop.

### Syntax

Works with all loop variants: `loop:row`, `loop:ol`, `loop:ul`

```html
<loop:row style.first=header x from items>
  <col>{{x.name}}</col>
</loop:row>
```

Position flexible:
```html
<loop:row x from items style.first=header>
  <col>{{x.name}}</col>
</loop:row>
```

### Behavior

- **First iteration:** Gets the specified style (`<row style=header>`)
- **Remaining iterations:** Plain tags without style attribute

### Examples

**Table with styled header row:**
```html
<table border>
  <loop:row style.first=header x from data>
    <col>{{x.field1}}</col>
    <col>{{x.field2}}</col>
  </loop:row>
</table>
```

**Ordered list with styled first item:**
```html
<loop:ol style.first=highlight x from steps>
  {{x.text}}
</loop:ol>
```

---

## 4. New Feature: Dynamic `style={{var}}`

Resolve style names from variables at compile time.

### Syntax

**Static rows only** (not inside loops):

```html
<row style={{myStyleVar}}>
  <col>Data</col>
</row>
```

With data:
```json
{
  "myStyleVar": "header"
}
```

Rendered as:
```html
<row style=header>
  <col>Data</col>
</row>
```

### Scope

- ✅ Works with `<row>`
- ✅ Works with `<li>`
- ❌ NOT supported inside `<loop>` templates (use `style.first` instead)

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
    <col>Dynamically styled</col>
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

## 5. Implementation Details

### Property Normalization Layer

The system uses a normalization layer that converts new property names to internal keys:

- User writes: `color=`, `bg=`
- Internal code uses: `font-color`, `shading`

This minimizes code changes in renderers while enforcing new user-facing names.

### Normalization Points

1. **Section properties** (`parse/parse.go`) - Properties in `[section]` blocks
2. **Inline attributes** (`render/body.go`) - Attributes in tags like `<row bg=...>`
3. **Helper function** (`render/style.go`, `parse/parse.go`) - `normalizePropertyKey()`

---

## 6. Files Changed

### Code (3 files)
- `parse/parse.go` - Added property normalization in section parsing
- `render/style.go` - Added `normalizePropertyKey()` function
- `render/body.go` - Updated `parseAttrs()`, loop regex, `expandLoops()` for `style.first`
- `render/compiler.go` - Updated `applyTableStyles()` prefix, added `resolveRowStyles()`

### Documentation (11 files)
- `.agents/skills/document-table.md` - Updated examples, added new features
- `.agents/skills/document-heading.md` - Property rename
- `.agents/skills/document-image.md` - Property rename
- `.agents/skills/document-style.md` - Property rename
- `.agents/skills/header-footer.md` - Property rename
- `docs/style.md` - Property rename, section format
- `docs/tags.md` - Property rename

### Examples (4 files)
- `docs/examples/simple.dcd` - Updated properties
- `docs/examples/features.dcd` - Updated properties
- `docs/examples/report.dcd` - Updated properties and section format
- `docs/examples/invoice.dcd` - Updated properties and section format

---

## 7. Migration Checklist

For each `.dcd` file:

- [ ] Replace `font-color=` with `color=`
- [ ] Replace `shading=` with `bg=`
- [ ] Replace `[table-style name]` with `[style:table name]`
- [ ] Test compilation with new format

**Automated migration** (Linux/Mac):
```bash
sed -i 's/font-color=/color=/g; s/shading=/bg=/g; s/\[table-style /[style:table /g' yourfile.dcd
```

---

## 8. Testing

All existing examples successfully regenerated:
- ✅ simple.dcd
- ✅ features.dcd
- ✅ report.dcd (with JSON data)
- ✅ invoice.dcd (with JSON data)
- ✅ PDF outputs

New features tested:
- ✅ `style.first` with `loop:row`, `loop:ol`, `loop:ul`
- ✅ `style={{var}}` with static rows
- ✅ Property normalization (`color`, `bg`)
- ✅ Section format `[style:table name]`

---

## 9. Backward Compatibility

**NONE** - This is a breaking change. All old format files will fail or behave incorrectly.

Old format is **not supported**. Users must migrate all `.dcd` files.
