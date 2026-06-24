# Style Configuration

## Page Style (`[style]`)

```ini
[style]
layout=A4
orientation=portrait
unit=mm
font-family="Times New Roman"
font-size=12
font-color=#000000
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
| `font-color` | `#000000` | Hex color |
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
font-color=#1F3864
bold=true
align=left
space-before=18
space-after=12
border-bottom=single
```

## Table Row Styles (`[table-style name]`)

```ini
[table-style header]
shading=#4472C4
font-color=#ffffff
font-weight=bold
align=center
border-bottom=single

[table-style alt]
shading=#D9E2F3
```

### Table Style Properties

| Property | Description |
|----------|-------------|
| `shading` | Row/cell background color |
| `align` | `center` / `right` |
| `font-color` | Text color |
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
font-color=#666666
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
| `font-family` | Font name |
| `font-size` | Font size in points |
| `font-color` | Hex color |
| `border` | `top` or `bottom` line |
| `margin` | Header/footer margin |

## Property Resolution Order

For any rendering property, the value is determined by:

1. **Local attribute** (on the tag itself, e.g. `align=center`)
2. **Named style** (heading style or table row style)
3. **Default style** (from `[style]` section)
