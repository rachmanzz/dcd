---
name: dcd-documents
description: Complete reference for the DCD DSL — sections, variables, body tags, styles, headings, tables, lists, loops, images, links, breaks, metadata, and header/footer
---

# DCD Documents

## 1. Style Configuration

Page layout and margin configuration for DCD documents.

```
[style]
layout=A4
unit=inch
orientation=portrait
font-family="Times New Roman"
font-size=12
color=#000000
line-height=1.5
m=1
```

### Layout

| Value    | Description                |
|----------|----------------------------|
| `A4`     | 210 × 297 mm               |
| `letter` | 8.5 × 11 in                |
| `legal`  | 8.5 × 14 in                |
| `A3`     | 297 × 420 mm               |
| `A5`     | 148 × 210 mm               |
| `B5`     | 176 × 250 mm               |
| `custom` | Requires explicit w / h    |

### Unit

`inch`, `cm`, `mm`, `pt`, `pica`

### Orientation

| Value       | Description                         |
|-------------|-------------------------------------|
| `portrait`  | Default. Taller than wide.          |
| `landscape` | Wider than tall. Swap width/height. |

```
[style]
layout=A4
unit=inch
orientation=landscape
```

### Font

| Property      | Description       | Example                   |
|---------------|-------------------|---------------------------|
| `font-family` | Font family name  | "Times New Roman", Arial  |
| `font-size`   | Base font size    | 12pt                      |
| `color`  | Text color        | #000000, black            |
| `line-height` | Line spacing      | 1.5                       |

```
[style]
layout=A4
unit=inch
font-family="Times New Roman"
font-size=12
color=#000000
line-height=1.5
```

### Paragraph

Global default for paragraph indentation.

| Property  | Description                      | Example |
|-----------|----------------------------------|---------|
| `indent`  | Left indent (in document unit)   | 1       |
| `hanging` | Hanging indent (in document unit)| 0.5     |

```
[style]
indent=0.5
hanging=0.25
```

Inline `<p indent=X>` / `<li indent=X>` overrides this default.

### Margins

All margin examples below assume:

```
[style]
layout=A4
unit=inch
orientation=portrait
font-family="Times New Roman"
font-size=12
line-height=1.5
```

**Uniform:**
```
m=1
```

**Axis:** `mx` = left & right, `my` = top & bottom.
```
mx=1
my=1
```

**Individual:** `mt` top, `mb` bottom, `ml` left, `mr` right.
```
mt=1
mb=1
ml=1
mr=1
```

**Default + Bottom:** `md` = margin default (all sides), `mb` = bottom (override).
```
md=1
mb=1
```

**Precedence (low → high):**
1. `m`
2. `mx` / `my`
3. `md`
4. `mt` / `mb` / `ml` / `mr`

## 2. Sections & Variables

Document content template with structured data.

```
[section 0]
name=userinfo
var=info, []entries
keys=username, date_field, time_field
formats=[date_field:dd-MM-yyyy], [time_field:HH\:m]

--- BODY ---
<w:c|b>Center Bold</w:c|b>
<p>your username is <b>{{info.username}}</b> created on <i>{{info.date_field}}</i> at <u>{{info.time_field}}</u></p>
<loop x from entries>
   {{x.name}} lives at {{x.address}}
</loop>

[section 1]
name=address
var=addr
keys=street, city, zip

--- BODY ---
<p>{{addr.street}}, {{addr.city}} - {{addr.zip}}</p>

[section 2]
name=simple
keys=title, message

--- BODY ---
<p>{{title}}: {{message}}</p>
```

### Section Properties

| Property   | Description                            |
|------------|----------------------------------------|
| `name`       | Section identifier — **REQUIRED** in every `[section N]`. Must be declared before `var=` and `keys=`. |
| `var`        | Comma-separated variable names. **Objects:** plain name (e.g. `info`). **Arrays (loop sources):** prefix with `[]` (e.g. `[]entries`). Pattern: `var=info, []entries` — **first** `info` is prefix for `{{info.key}}` via `keys`. **Subsequent** `[]entries` is a data source for `<loop x from entries>`. |
| `keys`       | Comma-separated field names for data variable resolution. For primary var: field names. For array fields requiring formatting: `source.field` (e.g. `items.date_field`). **CONDITIONAL DOT-NOTATION RULE:** Dotted paths (object/array fields) MUST NOT be registered in `keys=` UNLESS they are explicitly formatted in `formats=`. Optional — sections without `var`/`keys` pass `{{...}}` through as literals. |
| `formats`    | Per-key format: `[key:format]` or `[source.field:format]`. Defines the output format of a key. **EXCLUSIVE REGISTRATION RULE:** Any key or dotted path targeted in `formats=` MUST be explicitly listed in `keys=`. For array fields in loops, use `[source.field:format]` (e.g. `[items.date_field:dd-MM-yyyy`). |

> Properties use `=` separator (e.g. `name=example`).

### Var Usage

- **Section Limits:** Aim for ≤ 5 `var` entries and ≤ 15 `keys` entries per section.
- **Splitting Rule:** If you exceed these limits, split into a new logical section (e.g. `[section 1]`, `[section 2]`).

```
var=info, []entries
```

| Position | Name         | Source of data           | Access in body                                |
|----------|--------------|--------------------------|-----------------------------------------------|
| 1st      | `info`       | Resolved via `keys`      | `{{info.username}}`                           |
| 2nd+     | `[]entries`  | Array data source        | `<loop x from entries>{{x.name}}</loop>`      |

- **First name** (`info`): variable prefix. Fields listed in `keys`. Accessed as `{{info.key}}`.
- **Additional names** (`[]entries`, ...): data sources for loops. Accessed via `<loop x from entries>`, then `{{x.field}}` per item.

### Variables

`{{var.key}}` — e.g. `{{info.username}}`, `{{info.date_field}}`.

Built-in variables (resolved automatically, no registration needed):
- `{{date}}` — current compilation date
- `{{title}}` — document title (from `[title]` section)
- `{{page}}` — page number (works in header/footer; passed as literal in body)
- `{{total}}` — total pages (works in header/footer; passed as literal in body)

### Format Specifiers

Format is defined as `[key:format]` in the `formats` property.

| Specifier | Description     |
|-----------|-----------------|
| `dd`      | Day (01–31)     |
| `MM`      | Month (01–12)   |
| `yyyy`    | Year (4 digit)  |
| `HH`      | Hour (00–23)    |
| `mm`      | Minute (00–59)  |
| `ss`      | Second (00–59)  |

Example: `[date_field:yyyy-MM-dd]` → `2026-06-24`

Besides the specifiers above, format also supports regex patterns like `\d`, `\w`, or other regex.

### Format for Array Fields

For fields inside array objects (used in `<loop x from source>`) that need formatting, use **dotted path** in both `keys` and `formats`. Fields that do NOT need formatting must NOT appear in `keys=`:

```ini
[section 0]
var=info, []items
keys=title, items.date_field
formats=[items.date_field:dd-MM-yyyy]

--- BODY ---
<h1>{{info.title}}</h1>
<loop x from items>
  <p>{{x.name}} — {{x.date_field}}</p>   ← formatted via dotted path match
</loop>
```

After loop expansion, `{{x.date_field}}` becomes `{{items.0.date_field}}` — the format engine matches it against `items.date_field` by stripping the array index.

### Variable Registration Rule

`{{...}}` variables that reference data fields must be registered in:
- **`keys`** — field names or dotted paths
- **`var`** — data source names

**Strict Usage:** Every variable in `var=` and every key in `keys=` MUST be used at least once in `--- BODY ---`. Do NOT declare unused variables or keys.

Sections without `var` or `keys` are allowed: unresolvable `{{...}}` variables pass through as **literal strings** (e.g. `{{unknown}}` appears as-is). Built-in variables (`{{date}}`, `{{title}}`, `{{page}}`, `{{total}}`) are resolved automatically regardless of registration.

### Block Tags (outside `<p>`)

> **TAG BALANCING:** Every opened tag must be closed exactly. `<loop:ol>` MUST close with `</loop:ol>`, NOT `</loop>`.
>
> **NO `<w:*>` NESTING:** `<w:*>` tags MUST NOT contain other `<w:*>` tags. `<w:c><w:b>text</w:b></w:c>` is an error. Use OR logic in a single tag: `<w:c|b>text</w:c|b>`.
>
> **TAB INSIDE `<w:*>`:** `<tab>`, `<tab/>`, and `<tab size=N>` ARE allowed inside `<w:*>` tags. `<br>` is also allowed.
>
> **NO HEADING NESTING:** Heading tags `<h1>`–`<h6>` MUST NOT appear inside `<p>`, `<w:*>`, or other heading tags. `<p><h2>text</h2></p>` and `<h1><h2>text</h2></h1>` are errors.

| Tag                              | Description                     |
|----------------------------------|---------------------------------|
| `<w:c>...</w:c>`                 | Center                          |
| `<w:b>...</w:b>`                 | Bold                            |
| `<w:i>...</w:i>`                 | Italic                          |
| `<w:u>...</w:u>`                 | Underline                       |
| `<w:s>...</w:s>`                 | Strikethrough                   |
| `<w:c\|b>...</w:c\|b>`           | Center + Bold                   |
| `<w:b\|i>...</w:b\|i>`           | Bold + Italic                   |
| `<w:b\|i\|u>...</w:b\|i\|u>`     | Bold + Italic + Underline       |
| `<w:c\|s>...</w:c\|s>`           | Center + Strikethrough           |
| `<p>`                            | Paragraph                       |
| `<br>`                           | Line break                      |
| `<tab>`                          | Tab character (inside `<p>`)    |
| `<tab size=N>`                   | Tab with N spaces               |
| `<loop x from var>...</loop>`     | Iterate array `var`, each item as `x` |
| `<loop:ol x from var>...</loop:ol>` | Iterate + wrap `<ol><li>`       |
| `<loop:ul x from var>...</loop:ul>` | Iterate + wrap `<ul><li>`       |

> Note: `\|` inside table cells is markdown escape for `|` — the actual tag is `<w:b|i>` etc.

### Paragraph Properties

`<p>` and `<li>` tags accept paragraph-level formatting attributes:

| Property  | Example       | Description                            |
|-----------|---------------|----------------------------------------|
| `indent`  | `indent=0.5`  | Left indent (in document unit)         |
| `hanging` | `hanging=0.25`| Hanging indent (removed from first line) |

Applies to the entire paragraph. `hanging` shifts the first line left relative to `indent`.

```
<p indent=1 hanging=0.5>
  First line starts 0.5 from left margin,
  rest of paragraph indented 1 from margin.
</p>

<li indent=0.5 hanging=0.25>list item with custom indent</li>
```

**Precedence (high → low):**
1. Inline attribute on `<p>` or `<li>`
2. `[style] indent` / `[style] hanging` global

### Inline Tags (inside `<p>`, `<li>`, `<col>`)

| Tag              | Description             |
|------------------|-------------------------|
| `<b>...</b>`     | Bold                    |
| `<i>...</i>`     | Italic                  |
| `<u>...</u>`     | Underline               |
| `<s>...</s>`     | Strikethrough           |
| `<code>...</code>`| Monospace / code font  |
| `<mark>...</mark>`| Highlight (default yellow, optional: `<mark color=green>`) |
| `<sub>...</sub>`  | Subscript               |
| `<sup>...</sup>`  | Superscript             |
| `<set:flags>...</set:flags>` | Combined formatting |

Combined formatting with `<set:>`:

```
<p><set:b|i>Bold and Italic</set:b|i></p>
<p><set:b|u>Bold and Underline</set:b|u></p>
<p><set:i|code>Italic monospace</set:i|code></p>
<p><set:b|i|u>Bold, Italic, and Underline</set:b|i|u></p>
<p><set:s|b>Strikethrough and Bold</set:s|b></p>
```

**Available flags:** `b` (bold), `i` (italic), `u` (underline), `s` (strikethrough), `code` (monospace)

**Closing tag:** Must match opening flags: `<set:u>text</set:u>`, `<set:b|i>text</set:b|i>`

**Attributes:** Pass additional formatting via attributes on `<set:flags>`:

```
<p><set:u underline=double>double underline</set:u></p>
<p><set:b|u underline=dash>bold with dashed underline</set:b|u></p>
```

| Attribute    | Values                    | Description        |
|--------------|---------------------------|--------------------|
| `underline`  | `single`, `double`, `dotted`, `dash`, `wavy` | Underline style |

### Block Tags with Attributes

Block `<w:>` tags also accept attributes for additional formatting:

```
<w:u underline=double>double underline paragraph</w:u>
<w:u underline=dash>dashed underline paragraph</w:u>
```

| Tag                     | Attribute               | Values                    | Description        |
|-------------------------|-------------------------|---------------------------|--------------------|
| `<w:u>`                 | `underline`             | `single`, `double`, `dotted`, `dash`, `wavy` | Underline style |

### Tab Inside W Block Tags

`<tab>`, `<tab/>`, and `<tab size=N>` are allowed inside `<w:*>` tags:

```
<w:c|b>Name:<tab>John Doe</w:c|b>
<w:c>City:<tab size=4>Jakarta</w:c>
<w:b>Phone:<tab/>+62-812-3456-7890</w:b>
```

This enables formatted key-value layouts with tab stops inside centered, bold, or other styled blocks.

## 3. Headings

Heading `<h1>`–`<h6>` with global style in `[style:heading-N]`.

**RESTRICTION:** `<h1>` through `<h6>` MUST contain ONLY plain text and `{{vars}}`. Nested tags (`<b>`, `<i>`, `<u>`, `<code>`, etc.) are STRICTLY FORBIDDEN inside headings.

```
[style]
layout=A4
unit=inch
m=1

[style:heading-1]
font-family="Arial"
font-size=24
color=#2b5797
bold=true
space-before=18
space-after=12
border-bottom=1pt

[style:heading-2]
font-family="Arial"
font-size=18
color=#444444
bold=true
space-before=12
space-after=6

[style:heading-3]
font-family="Arial"
font-size=14
color=#444444
bold=false
space-before=6
space-after=3
```

Body:

```
--- BODY ---
<h1>Chapter 1: Introduction</h1>
<p>lorem ipsum...</p>
<h2>1.1 Background</h2>
<p>lorem ipsum...</p>
<h3>1.1.1 Sub Section</h3>
<p>lorem ipsum...</p>
```

Local override (higher priority):

```
<h1 color=red font-size=28>Chapter with local style</h1>
```

### Style Properties

| Property        | Description                      |
|-----------------|----------------------------------|
| `font-family`   | Heading font                     |
| `font-size`     | Font size (pt)                   |
| `color`         | Text color                       |
| `bold`          | `true` / `false`                 |
| `italic`        | `true` / `false`                 |
| `strike`        | `true` / `false`                 |
| `underline`     | `true`, `single`, `double`, `dotted`, `dash`, `wavy` |
| `caps`          | `true` / `false` — all capitals |
| `small-caps`    | `true` / `false` — small capitals |
| `letter-spacing`| Letter spacing (pt)              |
| `align`         | `left`, `center`, `right`        |
| `space-before`  | Space before (pt)                |
| `space-after`   | Space after (pt)                 |
| `border-bottom` | Bottom border line               |

### Precedence

1. Local attribute on tag `<h1 color=red>`
2. `[style:heading-N]` global
3. `[style]` font default

## 4. Tables

### Dynamic Table

```
<table border=1 width=100%>
<loop:row x from headers>
   <col>{{x}}</col>
</loop:row>
<loop:row x from entries>
   <col>{{x.field1}}</col>
   <col>{{x.field2}}</col>
</loop:row>
</table>
```

### Static Table

```
<table border=1>
  <row bg=#f0f0f0>
    <col align=center width=30%>Name</col>
    <col align=center width=30%>City</col>
    <col align=center width=40%>Age</col>
  </row>
  <row>
    <col align=left>John</col>
    <col align=left>Jakarta</col>
    <col align=center>25</col>
  </row>
</table>
```

### Tags

| Tag                              | Description                  |
|----------------------------------|------------------------------|
| `<table>...</table>`             | Table wrapper                |
| `<row>...</row>`                 | Row                          |
| `<col>...</col>`                 | Cell                         |
| `<loop:row x from var>...</loop:row>` | Loop data into rows    |

### Table Properties

| Property  | Example   | Description          |
|-----------|-----------|----------------------|
| `border`  | `1`       | Border width         |
| `width`   | `100%`    | Table width¹         |

¹ Not yet implemented (roadmap item).

### Row Properties

| Property  | Example       | Description          |
|-----------|---------------|----------------------|
| `bg`      | `#f0f0f0`     | Row background       |
| `style`   | `header`      | Named table-style    |

### Col Properties

| Property  | Example       | Description          |
|-----------|---------------|----------------------|
| `align`   | `center`      | Text alignment       |
| `width`   | `30%`         | Column width¹        |
| `bg`      | `#e0e0e0`     | Cell background      |
| `colspan` | `2`           | Merge columns¹       |
| `rowspan` | `2`           | Merge rows¹          |

¹ `docx.Cell.ct` is unexported — `GridSpan`/`VMerge`/`CellWidth` cannot be set. Library patch required.

### Named Table Style

```
[style:table header]
bg=#2b5797
color=white
font-weight=bold
align=center
border-bottom=2pt

[style:table alt]
bg=#f5f5f5
```

Usage:

```
<table border=1>
  <row style=header>
    <col>Name</col>
    <col>City</col>
  </row>
  <row style=alt>
    <col>John</col>
    <col>Jakarta</col>
  </row>
</table>
```

### Loop with style.first

Apply style to first row only:

```
<table border=1>
  <loop:row x from items style.first=header>
    <col>{{x.name}}</col>
    <col>{{x.value}}</col>
  </loop:row>
</table>
```

### Dynamic Row Style

Use variable for style name:

```
<row style={{myStyle}}>
  <col>Data</col>
</row>
```

## 5. Lists

Standalone lists (not from loop).

```
<ul>
  <li>item one</li>
  <li>item two</li>
  <li>item three</li>
</ul>

<ol>
  <li>first</li>
  <li>second</li>
  <li>third</li>
</ol>
```

Nested:

```
<ul>
  <li>fruit
    <ul>
      <li>apple</li>
      <li>mango</li>
    </ul>
  </li>
  <li>vegetable</li>
</ul>
```

### Tags

| Tag       | Description                                    |
|-----------|------------------------------------------------|
| `<ol>`    | Ordered list (supports `type=a/A/1/i/I`)       |
| `<ul>`    | Unordered list                                 |
| `<li>`    | List item (supports `indent`/`hanging`, see [Paragraph Properties](#paragraph-properties)) |

Supported `type` values for `<ol>`: `1` (default, numbers), `a` (lowercase letters), `A` (uppercase letters), `i` (lowercase roman), `I` (uppercase roman).

### Horizontal Rule

```
<hr>
```

| Tag       | Description              |
|-----------|--------------------------|
| `<hr>`    | Horizontal rule          |

Properties:

| Property | Example   | Description     |
|----------|-----------|-----------------|
| `width`  | `50%`     | Line width      |
| `color`  | `#cccccc` | Line color      |
| `thick`  | `2`       | Thickness (pt)¹ |

¹ Not yet implemented (roadmap item).

## 6. Loops

Iterate over array data sources declared in `var`.

The data source name must be listed with `[]` prefix in `var=`. See [Var Usage](#var-usage).

### Critical Loop Constraints

1. **Strict Syntax Order:** The iteration action (`x from source`) MUST come BEFORE any attributes (`type=...`). `<loop:ol x from items type=A>` ✅, `<loop:ol type=A x from items>` ❌.
2. **Source Matching:** The array source MUST be listed with a `[]` prefix in `var=`.
3. **Variable Access:** Inside the loop, access fields using the alias (e.g. `{{x.field}}`).
4. **Closing Tag Rule:** The closing tag MUST EXACTLY MATCH the opening variant prefix, but MUST OMIT the action and attributes. Opening: `<loop:ol x from items type=A>` ➔ Closing: `</loop:ol>` (NOT `</loop>` and NOT `</loop:ol type=A>`).
5. **List Loop Prohibition:** NEVER wrap a standard `<loop>` inside static `<ol>` or `<ul>` tags. Use the native `<loop:ol>` or `<loop:ul>` tags instead.
6. **Silent Index:** Use `{index+N}` (single braces) to insert a 1-based counter. Index starts at 0, so `{index+1}` = 1, 2, 3...

```
[section 0]
name=example
var=info, []entries
keys=title

--- BODY ---
<loop x from entries>
  {index+1}. {{x.field}}
</loop>
```

Here `entries` is an array source declared as `[]entries` in `var=info, []entries`.

### Tags

| Tag                                              | Description                            |
|--------------------------------------------------|----------------------------------------|
| `<loop x from name>...</loop>`                   | Iterate array `name`, each item as `x` |
| `<loop:ol x from name type=1>...</loop:ol>`      | Iterate + wrap `<ol><li>`              |
| `<loop:ul x from name>...</loop:ul>`             | Iterate + wrap `<ul><li>`              |
| `<loop:row x from name>...</loop:row>`           | Iterate into table rows                |

> Closing tag MUST EXACTLY MATCH the opening variant: `<loop:ol>` closes with `</loop:ol>` (NOT `</loop>`).

### Basic Loop

```
<loop x from entries>
  <p>{{x.name}} — {{x.value}}</p>
</loop>
```

- `x` — loop variable alias (any name)
- `entries` — must match a name with `[]` prefix in `var=` (e.g. `var=info, []entries`)
- Inside: `{{x.field}}` accesses a field on each array element

### Loop with Ordered List

```
<loop:ol x from items type=A>
  {{x.label}}
</loop:ol>
```

Renders as `<ol type=A><li>value</li><li>value</li></ol>`. Default `type` is `1` (numeric). Supported: `1`, `a`, `A`, `i`, `I`.

### Loop with Unordered List

```
<loop:ul x from items>
  {{x.label}}
</loop:ul>
```

Renders as `<ul><li>value</li><li>value</li></ul>`.

### Loop Index Counter

Use `{index+N}` (single braces) inside any loop variant to insert an auto-incrementing counter. Index starts at **0**, so `{index+1}` produces 1, 2, 3...

| Pattern     | Result                |
|-------------|-----------------------|
| `{index+1}` | 1, 2, 3, 4, ...      |
| `{index+0}` | 0, 1, 2, 3, ...      |
| `{index+10}`| 10, 11, 12, 13, ...  |

Works in `<loop>`, `<loop:ol>`, `<loop:ul>`, and `<loop:row>`:

```
<loop:ol x from items>
  {index+1}. {{x.label}}
</loop:ol>
```

Renders as `<ol><li>1. Item A</li><li>2. Item B</li><li>3. Item C</li></ol>`.

Multiple indexes in one template:

```
<loop x from entries>
  <p>Entry #{index+1} of {{x.total}}: {{x.name}}</p>
</loop>
```

### Loop into Table Rows

```
<table border=1>
<loop:row x from headers>
  <col>{{x}}</col>
</loop:row>
<loop:row x from entries>
  <col>{{x.field1}}</col>
  <col>{{x.field2}}</col>
</loop:row>
</table>
```

- First `loop:row` iterates `headers` — each item is a cell value (`{{x}}`)
- Second `loop:row` iterates `entries` — each item is an object (`{{x.field}}`)
- Each iteration produces a `<row>` with `<col>` cells

### Full Example

```
[section 0]
name=products
var=info, []items
keys=title, items.date, items.price
formats=[items.date:dd-MM-yyyy], [items.price:#,##0.00]

--- BODY ---
<h1>{{info.title}}</h1>
<table border=1 width=100%>
  <loop:row x from items>
    <col>{{x.name}}</col>
    <col align=right>{{x.price}}</col>
    <col>{{x.date}}</col>
  </loop:row>
</table>
```

Fields from array objects (`items.date`, `items.price`) use dotted path notation in `keys` and `formats`. See [Format for Array Fields](#format-for-array-fields).

## 7. Images

From data section:

```
[section 0]
name=gallery
var=source
keys=img, caption

--- BODY ---
<img={{source.img}} width=80% align=center>
<p><i>{{source.caption}}</i></p>
```

Static path:

```
<img=./assets/photo.jpg width=400>
```

### Properties

| Property   | Example        | Description                 |
|------------|----------------|-----------------------------|
| `width`    | `100%`, `400`  | Width (px or %)             |
| `height`   | `300`          | Height (px)                 |
| `align`    | `center`       | `left`, `center`, `right`   |
| `alt`      | "photo"        | Alternative text            |
| `border`   | `1`            | Border width                |
| `bg`  | `#f0f0f0`      | Background container        |

## 8. Links

Internal and external hyperlinks.

From data section:

```
<section 0>
var=source
keys=url, label

--- BODY ---
<a={{source.url}}>{{source.label}}</a>
```

Static:

```
<a=https://example.com>visit website</a>
```

Inline:

```
<p>click <a={{source.url}} target=_blank>here</a> for more info</p>
```

### Properties

| Property    | Example         | Description          |
|-------------|-----------------|----------------------|
| `target`    | `_blank`        | Open in new tab (DOCX always opens external links in new window) |
| `color`     | `#0055cc`       | Link color           |
| `underline` | `true`          | Underline            |

> **Limitation:** Hyperlinks render as blue underlined text but are **not clickable** in the DOCX output. The `godocx v0.1.5` library's `Hyperlink` struct lacks proper OOXML serialization (`<Children>` wrapper instead of raw `<w:r>`), which Word ignores. This affects all `<a=>` usage (inline and standalone). See [`KNOWN-LIMITATIONS.md`](/KNOWN-LIMITATIONS.md) for details.

### Bookmark

```
<a=#chapter1>see Chapter 1</a>
```

## 9. Page & Section Breaks

### Page Break

```
--- BODY ---
<p>page 1</p>
<pb>
<p>page 2</p>
```

| Tag              | Description       |
|------------------|-------------------|
| `<pb>`           | Page break        |
| `<page-break>`   | Alias for `<pb>`  |
| `<tab>`          | Tab character     |
| `<tab size=N>`   | Tab with N spaces |

`<tab>` can appear inside `<p>` paragraphs to insert a tab stop. The optional `size=N` attribute sets
an explicit number of spaces (defaults to one standard tab).

```
<p>Name:<tab>John Doe</p>
<p>Age:<tab size=4>25</p>
```

### Section Break

```
[section 0]
name=cover
var=info
keys=title, author

--- BODY ---
<h1>{{info.title}}</h1>
<p>{{info.author}}</p>

[section:next-page 1]

--- BODY ---
<p>new section after page break</p>
```

| Syntax                           | Description                           |
|----------------------------------|---------------------------------------|
| `[section:next-page N]`          | Section break + page break            |

`N` = section sequence number.

## 10. Metadata

Set document properties like title, subject, and author using the `[title]` section.

```
[title]
title=Document Title
subject=Document Subject
author=Author Name
```

### Properties

| Property  | Description                          | Example                    |
|-----------|--------------------------------------|----------------------------|
| `title`   | Document title                       | Annual Report 2025         |
| `subject` | Document subject/description         | Financial Summary          |
| `author`  | Document author/creator              | Finance Team               |

These properties are written to:
- **DOCX:** Document properties (`docProps/core.xml`)


### Built-in Variable: `{{title}}`

The `title` property can be referenced in headers and footers using the `{{title}}` variable.

```
[title]
title=My Report

[header]
left={{title}}
right={{date}}

[footer]
center={{title}} - Page {{page}}
```

### Full Example

```
[style]
layout=A4

[title]
title=Quarterly Business Review
subject=Q4 2024 Performance Report
author=Executive Team

[header]
left={{title}}
right={{date}}
font-size=8
color=#666666

[section 0]

--- BODY ---
<h1>{{title}}</h1>
<p>Prepared by: Finance Department</p>
```

The document will have:
- Title property set to "Quarterly Business Review"
- Subject set to "Q4 2024 Performance Report"
- Author set to "Executive Team"
- Header showing the title and current date
- Body displaying the title as heading

### Notes

- All properties are optional
- Properties are visible in document properties dialog (DOCX)
- The `{{title}}` variable only works in headers/footers and body content
- Use `{{date}}` for current date, `{{page}}` for page numbers

## 11. Header & Footer

Header and footer for document pages.

```
[header]
left={{title}}
right={{page}} / {{total}}

[footer]
center={{date}}
```

### Properties

| Property      | Description                            |
|---------------|----------------------------------------|
| `left`        | Left column content                    |
| `center`      | Center column content                  |
| `right`       | Right column content                   |
| `justify_between` | 2 or 3 comma-separated items spread evenly via tab stops. Use `\,` for literal comma |
| `font-family` | Header/footer font override            |
| `font-size`   | Font size                              |
| `color`  | Text color                             |
| `border`      | `top`, `bottom`, `none`                |
| `margin`      | Distance from header/footer to content |
| `first-page`  | `true` / `false` — show on page 1     |
| `mirror`      | `true` / `false` — swap left↔right    |

### justify_between

Replaces `left`/`center`/`right` with evenly-spaced columns using OOXML tab stops.

```
[header]
justify_between={{title}}, {{page}} / {{total}}

[footer]
justify_between=Dept. A\, B\, and C, {{date}}, Page {{page}}
```

| Items | Behavior |
|---|---|
| 2 items | Left-aligned + right-aligned |
| 3 items | Left + center + right |

**Comma escaping:** Use `\,` for a literal comma inside a column value (e.g. `Dept. A\, B\, and C`).

Works with all header/footer variables and font styling properties.

### Variables

| Variable      | Description          |
|---------------|----------------------|
| `{{page}}`    | Page number          |
| `{{total}}`   | Total pages          |
| `{{title}}`   | Document title       |
| `{{date}}`    | Compilation date     |

### Full Example

```
[style]
layout=A4
unit=inch
m=1

[header]
left={{title}}
right={{page}} / {{total}}
font-size=10
color=#999999
border=bottom
margin=0.3

[footer]
center={{date}}
font-size=9
color=#666666
border=top
margin=0.2
first-page=false
```

### justify_between Example

```
[style]
layout=A4
unit=inch
m=1

[header]
justify_between={{title}}, {{page}} / {{total}}
font-size=10
color=#999999
border=bottom
margin=0.3

[footer]
justify_between=Dept. A\, B\, and C, {{date}}, Page {{page}}
font-size=9
color=#666666
border=top
margin=0.2
```

## See Also

- `dcd-cli` — CLI usage and options
- `golang-programming` — Go library API
- `dcd-guide` — Project overview and patterns
