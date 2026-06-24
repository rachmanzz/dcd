# DCD — Document Compilation Definition

**DCD** is a specification and tool for compiling `.dcd` files into Microsoft Word (`.docx`) and PDF documents. It uses a simple INI-style format with sections and HTML-like template tags.

## Install

```bash
go install github.com/rachmanzz/dcd@latest
```

## Quick Start

Create `report.dcd`:

```ini
[style]
layout=A4
unit=inch
m=1

[section 0]
name=report
var=info
keys=title, author, date

--- BODY ---
<w:c|b>Test Document</w:c>
<p>Title: <b>{{info.title}}</b></p>
<p>Author: <i>{{info.author}}</i></p>
```

Run:

```bash
dcd report.dcd report.docx
dcd --data report.json report.dcd report.docx
```

## Sections

| Section                              | Description                      |
|--------------------------------------|----------------------------------|
| `[section N]`                        | Section with index N             |
| `[section:next-page N]`              | Section break + page break       |
| `[style]`                            | Page layout configuration        |
| `[style:heading-N]`                  | Heading 1-6 global style         |
| `[table-style name]`                 | Named table style                |
| `[header]` / `[footer]`              | Header and footer                |

### Section Properties

| Property    | Description                      |
|-------------|----------------------------------|
| `name`      | Section identifier               |
| `var`       | Variable prefix                  |
| `keys`      | Key list                         |
| `formats`   | Per-key format                   |

## Body Tags

### Block Tags

| Tag                              | Description                 |
|----------------------------------|-----------------------------|
| `<h1>`–`<h6>`                    | Heading                     |
| `<p>...</p>`                     | Paragraph                   |
| `<w:c>...</w:c>`                 | Center block                |
| `<w:b>...</w:b>`                 | Bold block                  |
| `<w:i>...</w:i>`                 | Italic block                |
| `<w:u>...</w:u>`                 | Underline block             |
| `<w:c|b>...</w:c|b>`             | Center + Bold               |
| `<w:b|i>...</w:b|i>`             | Bold + Italic               |
| `<w:b|i|u>...</w:b|i|u>`         | Bold + Italic + Underline   |
| `<pb>` / `<page-break>`          | Page break                  |
| `<br>`                           | Line break                  |
| `<hr>`                           | Horizontal rule             |

### Inline Tags (inside `<p>`)

| Tag               | Description   |
|-------------------|---------------|
| `<b>...</b>`      | Bold          |
| `<i>...</i>`      | Italic        |
| `<u>...</u>`      | Underline     |
| `<code>...</code>`| Monospace font|

### Loop Tags

| Tag                              | Description                  |
|----------------------------------|------------------------------|
| `<loop x from var>...</loop>`    | Iterate array                |
| `<loop:ol x from var>...</loop>` | Iterate + ordered list       |
| `<loop:ul x from var>...</loop>` | Iterate + unordered list     |
| `<loop:row x from var>...</loop>`| Iterate into table rows      |

### Table Tags

| Tag          | Description  |
|--------------|--------------|
| `<table>`    | Table        |
| `<row>`      | Table row    |
| `<col>`      | Table cell   |

### Other Tags

| Tag              | Description          |
|------------------|----------------------|
| `<img=path>`     | Image                |
| `<a=url>text</a>`| Hyperlink            |
| `<ul>` / `<ol>`  | List wrapper         |
| `<li>`           | List item            |

## Variables

`{{var.key}}` — resolved at compile time.

```
{{info.username}}
{{info.date_field}}
{{x.name}}          ← inside <loop x from var>
```

## Style Configuration

```ini
[style]
layout=A4
unit=inch
orientation=portrait
font-family="Times New Roman"
font-size=12
font-color=#000000
line-height=1.5
m=1
```

### Margin Shorthands

| Syntax  | Description                     |
|---------|---------------------------------|
| `m`     | Uniform margin (all sides)      |
| `mx`    | Horizontal (left & right)       |
| `my`    | Vertical (top & bottom)         |
| `mt`    | Margin top                      |
| `mb`    | Margin bottom                   |
| `ml`    | Margin left                     |
| `mr`    | Margin right                    |
| `md`    | Margin default (all sides)      |

## Skills

Detailed specifications for each feature are available in [`.agents/skills/`](.agents/skills/):

- [document-style.md](.agents/skills/document-style.md)
- [document-body.md](.agents/skills/document-body.md)
- [document-heading.md](.agents/skills/document-heading.md)
- [document-table.md](.agents/skills/document-table.md)
- [document-image.md](.agents/skills/document-image.md)
- [document-list.md](.agents/skills/document-list.md)
- [document-loop.md](.agents/skills/document-loop.md)
- [document-link.md](.agents/skills/document-link.md)
- [document-break.md](.agents/skills/document-break.md)
- [header-footer.md](.agents/skills/header-footer.md)

## Library

```go
import "github.com/rachmanzz/dcd/data"
import "github.com/rachmanzz/dcd/parse"
import "github.com/rachmanzz/dcd/render"

doc, _ := parse.Parse("input.dcd")
r := render.NewDocxRenderer()
render.New(doc, data.NewDataSet(nil), r).Run("output.docx")
```

## License

MIT
