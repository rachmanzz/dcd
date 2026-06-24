# Document Body

Document content template with structured data.

## Usage

Multiple sections, each with its own `--- BODY ---`:

```
[section 0]
name=userinfo
var=info, entries
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

## Section

| Property   | Description                            |
|------------|----------------------------------------|
| `name`       | Section identifier                     |
| `var`        | Comma-separated variable names, each is an **object/map**. Pattern: `var=nameA, nameB, ...` — **first** `nameA` is prefix for `{{nameA.key}}` via `keys`. **Subsequent** `nameB` are data source names used by `<loop x from nameB>`. See [Var Usage](#var-usage) below. |
| `keys`       | Key list, comma or `[...]`. Used when `var` is not needed (standalone). Required when `var` is absent. |
| `formats`    | Per-key format: `[key:format]`. Defines the output format of a key. The key must be listed in `keys`. When formatting a var object key (e.g. `info.name`), `keys` must include `info.name`. |
| `layout`     | Page size: `A4`, `Letter`, `Legal`, `A3`, `A5`, `B5`, `custom` (overrides `[style]`) |
| `orientation`| `portrait` / `landscape` (overrides `[style]`) |
| `start-page` | Starting page number (default 1)              |

## Block Tags (outside `<p>`)

| Tag                              | Description                     |
|----------------------------------|---------------------------------|
| `<w:c>...</w:c>`                 | Center                          |
| `<w:b>...</w:b>`                 | Bold                            |
| `<w:i>...</w:i>`                 | Italic                          |
| `<w:u>...</w:u>`                 | Underline                       |
| `<w:c|b>...</w:c|b>`             | Center + Bold                   |
| `<w:b|i>...</w:b|i>`             | Bold + Italic                   |
| `<w:b|i|u>...</w:b|i|u>`         | Bold + Italic + Underline       |
| `<p>`                            | Paragraph                       |
| `<br>`                           | Line break                      |
| `<loop x from var>...</loop>`     | Iterate array `var`, each item as `x` |
| `<loop:ol x from var>...</loop>`  | Iterate + wrap `<ol><li>`       |
| `<loop:ul x from var>...</loop>`  | Iterate + wrap `<ul><li>`       |

## Inline Tags (inside `<p>`, `<li>`, `<col>`)

| Tag              | Description             |
|------------------|-------------------------|
| `<b>...</b>`     | Bold                    |
| `<i>...</i>`     | Italic                  |
| `<u>...</u>`     | Underline               |
| `<code>...</code>`| Monospace / code font  |
| `<set:flags>...</set:flags>` | Combined formatting |

### Combined Formatting with `<set:>`

Apply multiple formatting flags simultaneously:

```
<p><set:b|i>Bold and Italic</set:b|i></p>
<p><set:b|u>Bold and Underline</set:b|u></p>
<p><set:i|code>Italic monospace</set:i|code></p>
<p><set:b|i|u>Bold, Italic, and Underline</set:b|i|u></p>
```

**Available flags:** `b` (bold), `i` (italic), `u` (underline), `code` (monospace)

**Closing tag:** Can be `</set:flags>` (matching) or `</set>` (simplified)

## Var Usage

```
var=info, entries
```

| Position | Name       | Source of data           | Access in body                      |
|----------|------------|--------------------------|-------------------------------------|
| 1st      | `info`     | Resolved via `keys`      | `{{info.username}}`                 |
| 2nd+     | `entries`  | Array data source        | `<loop x from entries>{{x.name}}</loop>` |

- **First name** (`info`): variable prefix. Fields listed in `keys`. Accessed as `{{info.key}}`.
- **Additional names** (`entries`, ...): data sources for loops. Accessed via `<loop x from entries>`, then `{{x.field}}` per item.

## Variables

`{{var.key}}` — e.g. `{{info.username}}`, `{{info.date_field}}`.

## Format Specifiers

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
