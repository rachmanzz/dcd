# Document Style

Page layout and margin configuration for DCD documents.

## Usage

```
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

## Layout

| Value    | Description                |
|----------|----------------------------|
| `A4`     | 210 × 297 mm               |
| `letter` | 8.5 × 11 in                |
| `legal`  | 8.5 × 14 in                |
| `A3`     | 297 × 420 mm               |
| `A5`     | 148 × 210 mm               |
| `B5`     | 176 × 250 mm               |
| `custom` | Requires explicit w / h    |

## Unit

`inch`, `cm`, `mm`, `pt`, `pica`

## Orientation

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

## Font

| Property      | Description       | Example                   |
|---------------|-------------------|---------------------------|
| `font-family` | Font family name  | "Times New Roman", Arial  |
| `font-size`   | Base font size    | 12pt                      |
| `font-color`  | Text color        | #000000, black            |
| `line-height` | Line spacing      | 1.5                       |

```
[style]
layout=A4
unit=inch
font-family="Times New Roman"
font-size=12
font-color=#000000
line-height=1.5
```

## Margins

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

### Uniform

```
[style]
layout=A4
unit=inch
orientation=portrait
font-family="Times New Roman"
font-size=12
line-height=1.5
m=1
```

### Axis

```
[style]
layout=A4
unit=inch
orientation=portrait
font-family="Times New Roman"
font-size=12
line-height=1.5
mx=1
my=1
```

`mx` = left & right, `my` = top & bottom.

### Individual

```
[style]
layout=A4
unit=inch
orientation=portrait
font-family="Times New Roman"
font-size=12
line-height=1.5
mt=1
mb=1
ml=1
mr=1
```

`mt` top, `mb` bottom, `ml` left, `mr` right.

### Default + Bottom

```
[style]
layout=A4
unit=inch
orientation=portrait
font-family="Times New Roman"
font-size=12
line-height=1.5
md=1
mb=1
```

`md` = margin default (all sides), `mb` = bottom (override).

### Precedence (low → high)

1. `m`
2. `mx` / `my`
3. `md`
4. `mt` / `mb` / `ml` / `mr`
