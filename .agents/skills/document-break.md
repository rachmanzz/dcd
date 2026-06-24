# Document Break

Page break and section break.

## Page Break

```
--- BODY ---
<p>page 1</p>
<pb>
<p>page 2</p>
```

| Tag             | Description       |
|-----------------|-------------------|
| `<pb>`          | Page break        |
| `<page-break>`  | Alias for `<pb>`  |

## Section Break

```
[section 0]
name=cover
var=info
keys=title, author

--- BODY ---
<h1>{{info.title}}</h1>
<p>{{info.author}}</p>

[section:next-page 1]
orientation=landscape
start-page=1

--- BODY ---
<p>new section, landscape, page number reset to 1</p>
```

| Syntax                           | Description                           |
|----------------------------------|---------------------------------------|
| `[section:next-page N]`          | Section break + page break            |

`N` = section sequence number.

### Properties

| Property      | Description                         |
|---------------|-------------------------------------|
| `layout`      | Page size: `A4`, `Letter`, `Legal`, `A3`, `A5`, `B5`, `custom` |
| `orientation` | `portrait` / `landscape`            |
| `start-page`  | Starting page number (default 1) — recognized, not yet implemented |
| `title`       | Section title (for navigation)      |
