# Body Tags

All tags go in the body section (after `---`).

## Block Tags

| Tag | Description |
|-----|-------------|
| `<h1>text</h1>` – `<h6>text</h6>` | Heading level 1–6 |
| `<p>text</p>` | Paragraph with inline formatting |
| `<w:flags>text</w>` | Wrapped paragraph (`c` = center, `b` = bold) |
| `<br>` | Line break |
| `<hr attrs>` | Horizontal rule |
| `<pb>` / `<page-break>` | Page break |

### Horizontal Rule Attributes

| Attribute | Example | Description |
|-----------|---------|-------------|
| `color` | `color=#333` | Line color |
| `width` | `width=50%` or `width=100mm` | Width |
| `thick` | `thick=2pt` | Thickness |

## Inline Tags (inside `<p>`)

| Tag | Description |
|-----|-------------|
| `<b>text</b>` | Bold |
| `<i>text</i>` | Italic |
| `<u>text</u>` | Underline |
| `<a=url>text</a>` | Hyperlink |

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
| `shading` | `shading=#f0f0f0` | Background shading |

## Table Tags

```html
<table border>
  <row style=header shading=#e0e0e0>
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
| `shading` | `<row>` | Row background color |
| `align` | `<col>` | `center` / `right` |
| `shading` | `<col>` | Cell background color |

### Named Table Row Styles

```ini
[table-style header]
shading=#4472C4
font-color=#ffffff
font-weight=bold
align=center
border-bottom=single
```

## List Tags

```html
<ul>
  <li>Item 1</li>
  <li>Item 2
    <ol>
      <li>Sub-item</li>
    </ol>
  </li>
</ul>
```

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
