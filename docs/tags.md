# Body Tags

All tags go in the body section (after `---`).

## Block Tags

| Tag | Description |
|-----|-------------|
| `<h1>text</h1>` – `<h6>text</h6>` | Heading level 1–6 |
| `<p>text</p>` | Paragraph with inline formatting |
| `<w:flags>text</w:flags>` | Wrapped paragraph (`c` = center, `b` = bold, `i` = italic, `u` = underline) |
| `<br>` | Line break |
| `<hr attrs>` | Horizontal rule |
| `<pb>` / `<page-break>` | Page break |

### Horizontal Rule Attributes

| Attribute | Example | Description |
|-----------|---------|-------------|
| `color` | `color=#333` | Line color |
| `width` | `width=50%` or `width=100mm` | Width |
| `thick` | `thick=2pt` | Thickness |

## Inline Tags

Inline formatting works inside `<p>`, `<li>`, and `<col>` tags.

| Tag | Description |
|-----|-------------|
| `<b>text</b>` | Bold |
| `<i>text</i>` | Italic |
| `<u>text</u>` | Underline |
| `<code>text</code>` | Monospace / code font |
| `<set:flags>text</set>` | Combined formatting |
| `<a=url>text</a>` | Hyperlink |

### Combined Formatting

Use `<set:flags>` to apply multiple formatting styles:

```html
<p><set:b|i>Bold and Italic</set:b|i></p>
<p><set:b|u>Bold and Underline</set:b|u></p>
<p><set:i|code>Italic monospace</set:i|code></p>
<p><set:b|i|u>All three</set:b|i|u></p>
```

**Available flags:** `b`, `i`, `u`, `code`

### Examples

```html
<p>Paragraph with <b>bold</b> and <i>italic</i></p>
<p>Combined: <set:b|i>bold-italic</set:b|i></p>
<li>List item with <code>code</code> and <u>underline</u></li>
<col>Table cell with <set:b|u>bold-underline</set:b|u> text</col>
```

### Hyperlink Attributes

| Attribute | Example | Description |
|-----------|---------|-------------|
| `color` | `color=blue` | Link color |
| `underline` | `underline=false` | Show/hide underline |

## Image Tag

```html
<img=path/to/image.png attrs>
```

| Attribute | Example | Description |
|-----------|---------|-------------|
| `width` | `width=50%` or `width=100mm` | Width (percentage of page width, or absolute) |
| `height` | `height=50mm` | Height |
| `align` | `align=center` | `center` or `right` |
| `alt` | `alt=Diagram` | Alt text |
| `border` | `border=1pt` | Border thickness |
| `bg` | `bg=#f0f0f0` | Background shading |

## Table Tags

```html
<table border>
  <row style=header bg=#e0e0e0>
    <col>Name</col>
    <col align=center>Age</col>
  </row>
  <row>
    <col>John</col>
    <col>30</col>
  </row>
</table>
```

### Table Attributes

| Attribute | Applies To | Description |
|-----------|-----------|-------------|
| `border` | `<table>` | Enables table grid borders |
| `style` | `<row>` | Named table row style |
| `bg` | `<row>` | Row background color |
| `align` | `<col>` | `center` / `right` |
| `bg` | `<col>` | Cell background color |

### Named Table Row Styles

```ini
[style:table header]
bg=#4472C4
color=#ffffff
font-weight=bold
align=center
border-bottom=single
```

## List Tags

```html
<ul>
  <li>Item 1</li>
  <li>Item 2</li>
  <li>Item 3</li>
</ul>

<ol>
  <li>Step 1</li>
  <li>Step 2</li>
  <li>Step 3</li>
</ol>
```

**Note:** Nested lists are **not supported**. Nested list tags will be stripped from the output.

## Loop Tags

```html
<loop x from items>
  <p>{{x.name}}: {{x.value}}</p>
</loop>
```

| Variant | Description |
|---------|-------------|
| `<loop x from source>...</loop>` | Plain iteration |
| `<loop:ol x from source>...</loop>` | Ordered list iteration |
| `<loop:ul x from source>...</loop>` | Unordered list iteration |
| `<loop:row x from source>...</loop>` | Table row iteration |

Inside loops, `{{x}}` is the item itself, `{{x.field}}` accesses item fields.
