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

--- BODY ---
<p>new section after page break</p>
```

| Syntax                           | Description                           |
|----------------------------------|---------------------------------------|
| `[section:next-page N]`          | Section break + page break            |

`N` = section sequence number.
