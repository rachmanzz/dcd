# Document Table

## Dynamic Table

```
<table border=1 width=100%>
<loop:row x from headers>
   <col>{{x}}</col>
</loop:row>
<loop:row x from entries>
   <col>{{x.field1}}</col>
   <col>{{x.field2}}</col>
</loop:row>
</table>
```

## Static Table

```
<table border=1>
  <row bg=#f0f0f0>
    <col align=center width=30%>Name</col>
    <col align=center width=30%>City</col>
    <col align=center width=40%>Age</col>
  </row>
  <row>
    <col align=left>John</col>
    <col align=left>Jakarta</col>
    <col align=center>25</col>
  </row>
</table>
```

## Tags

| Tag                              | Description                  |
|----------------------------------|------------------------------|
| `<table>...</table>`             | Table wrapper                |
| `<row>...</row>`                 | Row                          |
| `<col>...</col>`                 | Cell                         |
| `<loop:row x from var>...</loop:row>` | Loop data into rows    |

## Table Properties

| Property  | Example   | Description          |
|-----------|-----------|----------------------|
| `border`  | `1`       | Border width         |
| `width`   | `100%`    | Table width          |

## Row Properties

| Property  | Example       | Description          |
|-----------|---------------|----------------------|
| `bg`      | `#f0f0f0`     | Row background       |
| `style`   | `header`      | Named table-style    |

## Col Properties

| Property  | Example       | Description          |
|-----------|---------------|----------------------|
| `align`   | `center`      | Text alignment       |
| `width`   | `30%`         | Column width         |
| `bg`      | `#e0e0e0`     | Cell background      |
| `colspan` | `2`           | Merge columns        |
| `rowspan` | `2`           | Merge rows           |

## Named Table Style

```
[style:table header]
bg=#2b5797
color=white
font-weight=bold
align=center
border-bottom=2pt

[style:table alt]
bg=#f5f5f5
```

Usage:

```
<table border=1>
  <row style=header>
    <col>Name</col>
    <col>City</col>
  </row>
  <row style=alt>
    <col>John</col>
    <col>Jakarta</col>
  </row>
</table>
```

## Loop with style.first

Apply style to first row only:

```
<table border=1>
  <loop:row style.first=header x from items>
    <col>{{x.name}}</col>
    <col>{{x.value}}</col>
  </loop:row>
</table>
```

## Dynamic Row Style

Use variable for style name:

```
<row style={{myStyle}}>
  <col>Data</col>
</row>
```
