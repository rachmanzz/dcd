# Document Heading

Heading `<h1>`–`<h6>` with global style in `[style:heading-N]`.

## Usage

```
[style]
layout=A4
unit=inch
m=1

[style:heading-1]
font-family="Arial"
font-size=24
font-color=#2b5797
bold=true
space-before=18
space-after=12
border-bottom=1pt

[style:heading-2]
font-family="Arial"
font-size=18
font-color=#444444
bold=true
space-before=12
space-after=6

[style:heading-3]
font-family="Arial"
font-size=14
font-color=#444444
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
<h1 font-color=red font-size=28>Chapter with local style</h1>
```

## Style Properties

| Property      | Description             |
|---------------|-------------------------|
| `font-family` | Heading font            |
| `font-size`   | Font size (pt)          |
| `font-color`  | Text color              |
| `bold`        | `true` / `false`        |
| `italic`      | `true` / `false`        |
| `underline`   | `true` / `false`        |
| `align`       | `left`, `center`, `right` |
| `space-before`| Space before (pt)       |
| `space-after` | Space after (pt)        |
| `border-bottom` | Bottom border line    |

## Precedence

1. Local attribute on tag `<h1 font-color=red>`
2. `[style:heading-N]` global
3. `[style]` font default
