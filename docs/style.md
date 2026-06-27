# Style Configuration

## Page Style (`[style]`)

```ini
[style]
layout=A4
orientation=portrait
unit=mm
font-family="Times New Roman"
font-size=12
color=#000000
line-height=1.5
m=25.4
```

### Margin Shorthands

| Key | Description |
|-----|-------------|
| `m` | All margins uniform |
| `mx` | Horizontal (left + right) |
| `my` | Vertical (top + bottom) |
| `mt` | Margin top |
| `mb` | Margin bottom |
| `ml` | Margin left |
| `mr` | Margin right |
| `md` | All margins (overrides `m`) |

## Default Text Style

Properties in `[style]` serve as the default for all text:

| Property | Example | Description |
|----------|---------|-------------|
| `font-family` | `Arial` | Font name |
| `font-size` | `12` | Font size in points |
| `color` | `#000000` | Hex color |
| `line-height` | `1.5` | Line spacing multiplier |
| `bold` | `true` | Bold (headings) |
| `italic` | `true` | Italic (headings) |
| `underline` | `true` | Underline (headings) |
| `align` | `center` | Alignment (headings) |
| `space-before` | `12` | Space before in points (headings) |
| `space-after` | `6` | Space after in points (headings) |
| `border-bottom` | `single` | Bottom border (headings) |

## Heading Styles (`[style:heading-N]`)

Override defaults per heading level:

```ini
[style:heading-1]
font-family=Arial
font-size=24
color=#1F3864
bold=true
align=left
space-before=18
space-after=12
border-bottom=single
keep-next=true
keep-lines=true
```

Heading-specific properties:

| Property | Description |
|----------|-------------|
| `keep-next` | Keep heading on same page as next paragraph |
| `keep-lines` | Keep all heading lines on one page |
| `numbering` | Not available — requires library support for numbering definitions |

## Table Row Styles (`[style:table name]`)

```ini
[style:table header]
bg=#4472C4
color=#ffffff
font-weight=bold
align=center
border-bottom=single

[style:table alt]
bg=#D9E2F3
```

### Table Style Properties

| Property | Description |
|----------|-------------|
| `bg` | Row/cell background color |
| `align` | `center` / `right` |
| `color` | Text color |
| `font-weight` | `bold` |
| `border-bottom` | Bottom border on cells |

## Header / Footer (`[header]` / `[footer]`)

```ini
[header]
left=Draft Document
center=Quarterly Report
right={{date}}
font-family=Arial
font-size=9
color=#666666
border=bottom

[footer]
left=Page {{page}} of {{total}}
right=Confidential
border=top
```

### Header/Footer Properties

| Property | Description |
|----------|-------------|
| `left` | Left-aligned content |
| `center` | Center-aligned content |
| `right` | Right-aligned content |
| `justify_between` | 2 or 3 comma-separated items spread evenly via tab stops. Use `\,` for literal comma |
| `font-family` | Font name |
| `font-size` | Font size in points |
| `color` | Hex color |
| `border` | `top` or `bottom` line |
| `margin` | Header/footer margin |
| `first-page` | `true` / `false` — show on page 1 |
| `mirror` | `true` / `false` — swap left↔right for even pages |

### justify_between Example

```ini
[header]
justify_between={{title}}, {{page}} / {{total}}
font-size=9
color=#666666
border=bottom

[footer]
justify_between=Dept. A\, B\, and C, {{date}}, Page {{page}}
border=top
```

## Property Resolution Order

For any rendering property, the value is determined by:

1. **Local attribute** (on the tag itself, e.g. `align=center`)
2. **Named style** (heading style or table row style)
3. **Default style** (from `[style]` section)

## Style Enhancements (v0.2.0)

### Combined Inline Formatting

Beyond single-tag formatting (`<b>`, `<i>`, `<u>`), use `<set:flags>` for multiple styles:

```html
<p><set:b|i>Bold and Italic</set:b|i></p>
<p><set:b|u>Bold and Underline</set:b|u></p>
<p><set:i|code>Italic monospace</set:i|code></p>
```

**Works in:** Paragraphs, list items, table cells

### Dynamic Row Styling

Resolve style names from variables:

```html
<row style={{rowStyleVar}}>
  <col>Data</col>
</row>
```

### Loop Styling with style.first

Apply a style to the first iteration only:

```html
<loop:row style.first=header x from items>
  <col>{{x.name}}</col>
</loop:row>
```

**Result:** First row has `style=header`, remaining rows are plain.

### Inline Attributes

Both `<row>` and `<col>` support inline style attributes:

```html
<row bg=#f0f0f0>
  <col color=#000 align=center>Text</col>
</row>
```

## Migration Notes (v0.2.0)

### Property Renames

| Old (deprecated) | New | Applies To |
|-----------------|-----|------------|
| `font-color` | `color` | All sections, inline attributes |
| `shading` | `bg` | Table rows, cells, images |

### Section Format Change

`[table-style name]` → `[style:table name]`

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

### Automated Migration

```bash
sed -i 's/font-color=/color=/g; s/shading=/bg=/g; s/\[table-style /[style:table /g' file.dcd
```

See `CHANGES.md` for complete migration guide.
