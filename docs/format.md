# Document Format

## Sections

A `.dcd` file is split into sections delimited by `[section-name]` headers.

### Section Types

| Section | Purpose |
|---------|---------|
| `[style]` | Page layout and default text style |
| `[style:heading-N]` | Heading style (N = 1–6) |
| `[style:table name]` | Named table row style |
| `[title]` | Document metadata |
| `[header]` | Page header |
| `[footer]` | Page footer |
| `[section N]` | Content section with index N |
| `[section:next-page N]` | Content section with page break |

### Section Format

```ini
[section 0]
name=cover
var=data
keys=title, author
formats=[date_field:02-01-2006]
layout=A4
orientation=portrait

--- BODY ---
Body content here...
```

Properties use `=` or `:` as separator:

```ini
font-family=Arial
border-bottom:2pt
```

The body marker is exactly `---` on its own line.

### Section Properties

| Property | Applies To | Description |
|----------|-----------|-------------|
| `name` | `[section N]` | Section identifier |
| `var` | `[section N]` | Variable prefix for `{{var.key}}` |
| `keys` | `[section N]` | Comma-separated key names |
| `formats` | `[section N]` | Date format mappings: `[key:GoDateFormat]` |
| `layout` | `[section N]` | Page size (overrides `[style]`) |
| `orientation` | `[section N]` | `portrait` / `landscape` |

## Variables

`{{path.to.key}}` is resolved at compile time:

```ini
[section 0]
var=data
keys=title, author

--- BODY ---
<p>Title: <b>{{data.title}}</b></p>
<p>Author: <i>{{data.author}}</i></p>
```

Inside `<loop x from source>`, use `{{x.field}}`.

### Built-in Variables

| Variable | Resolves To |
|----------|------------|
| `{{date}}` | Current date (2006-01-02) |
| `{{page}}` | Page number field |
| `{{total}}` | Total pages field (DOCX only) |
| `{{title}}` | Document title from `[title]` |

## Page Sizes

| Preset | Width × Height (mm) |
|--------|-------------------|
| `letter` | 215.9 × 279.4 |
| `legal` | 215.9 × 355.6 |
| `a3` | 297 × 420 |
| `a4` | 210 × 297 |
| `a5` | 148 × 210 |
| `b5` | 176 × 250 |
| `custom` | From `w`/`h` props |

## Units

| Unit | Aliases | Conversion to mm |
|------|---------|-----------------|
| mm | — | 1 |
| cm | — | 10 |
| inch | `in` | 25.4 |
| pt | — | 0.3528 |
| pica | — | 4.2333 |

## New Features (v0.2.0)

### Combined Inline Formatting

Use `<set:flags>` for multiple formatting styles:

```html
<p><set:b|i>Bold and Italic</set:b|i></p>
<p><set:b|i|u>All three styles</set:b|i|u></p>
```

**Available flags:** `b`, `i`, `u`, `code`

### Loop Styling

Apply style to first item only:

```html
<loop:row style.first=header x from items>
  <col>{{x.name}}</col>
</loop:row>
```

**Position flexible:** Can be before or after `x from items`

### Dynamic Styles

Resolve style names from variables:

```html
<row style={{myStyleVar}}>
  <col>Data</col>
</row>
```

### Loop Dot Notation

Access nested data in loops:

```html
<loop:row x from invoice.items>
  <col>{{x.id}}</col>
  <col>{{x.name}}</col>
</loop:row>
```

## Property Names (v0.2.0)

### Updated Names

| Old Name (deprecated) | New Name | Usage |
|----------------------|----------|-------|
| `font-color` | `color` | Text color |
| `shading` | `bg` | Background color |

**Migration:** Replace old names with new names in all `.dcd` files.

**Example:**
```ini
# Old
[style]
font-color=#000000

# New
[style]
color=#000000
```

## Limitations

- **Nested lists:** Not supported. Use flat lists only.
- **PDF inline formatting:** Limited support in lists/tables (use DOCX for rich formatting).

See `KNOWN-LIMITATIONS.md` for details.
