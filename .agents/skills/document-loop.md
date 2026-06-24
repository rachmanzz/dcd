# Document Loop

Iterate over array data sources declared in `var`.

## Var Setup

The data source name must be listed in `var` (after the first name). See [document-body.md](document-body.md#var-usage).

```
[section 0]
name=example
var=info, entries
keys=title

--- BODY ---
<loop x from entries>
  {{x.field}}
</loop>
```

Here `entries` is the 2nd name in `var=info, entries` — an array data source for the loop.

## Tags

| Tag                                    | Description                            |
|----------------------------------------|----------------------------------------|
| `<loop x from name>...</loop>`         | Iterate array `name`, each item as `x` |
| `<loop:ol x from name>...</loop>`      | Iterate + wrap each in `<ol><li>`      |
| `<loop:ul x from name>...</loop>`      | Iterate + wrap each in `<ul><li>`      |
| `<loop:row x from name>...</loop:row>` | Iterate into table rows                |

## Basic Loop

```
<loop x from entries>
  <p>{{x.name}} — {{x.value}}</p>
</loop>
```

- `x` — loop variable alias (any name)
- `entries` — must match a name in `var` (2nd position or later)
- Inside: `{{x.field}}` accesses a field on each array element

## Loop with Ordered List

```
<loop:ol x from items>
  {{x.label}}
</loop:ol>
```

Renders as `<ol><li>value</li><li>value</li></ol>`.

## Loop with Unordered List

```
<loop:ul x from items>
  {{x.label}}
</loop:ul>
```

Renders as `<ul><li>value</li><li>value</li></ul>`.

## Loop into Table Rows

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

## Full Example

```
[section 0]
name=products
var=info, items
keys=title

--- BODY ---
<h1>{{info.title}}</h1>
<table border=1 width=100%>
  <loop:row x from items>
    <col>{{x.name}}</col>
    <col>${{x.price}}</col>
  </loop:row>
</table>
```
