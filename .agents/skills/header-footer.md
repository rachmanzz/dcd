# Header & Footer

Header and footer for document pages.

## Usage

```
[header]
left={{title}}
right={{page}} / {{total}}

[footer]
center={{date}}
```

## Properties

| Property      | Description                            |
|---------------|----------------------------------------|
| `left`        | Left column content                    |
| `center`      | Center column content                  |
| `right`       | Right column content                   |
| `font-family` | Header/footer font override            |
| `font-size`   | Font size                              |
| `color`  | Text color                             |
| `border`      | `top`, `bottom`, `none`                |
| `margin`      | Distance from header/footer to content |
| `first-page`  | `true` / `false` — show on page 1     |
| `mirror`      | `true` / `false` — swap left↔right    |

## Variables

| Variable      | Description          |
|---------------|----------------------|
| `{{page}}`    | Page number          |
| `{{total}}`   | Total pages          |
| `{{title}}`   | Document title       |
| `{{date}}`    | Compilation date     |

## Full Example

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
