# Document Image

## Usage

From data section:

```
[section 0]
name=gallery
var=source
keys=img, caption

--- BODY ---
<img={{source.img}} width=80% align=center>
<p><i>{{source.caption}}</i></p>
```

Static path:

```
<img=./assets/photo.jpg width=400>
```

## Properties

| Property   | Example        | Description                 |
|------------|----------------|-----------------------------|
| `width`    | `100%`, `400`  | Width (px or %)             |
| `height`   | `300`          | Height (px)                 |
| `align`    | `center`       | `left`, `center`, `right`   |
| `alt`      | "photo"        | Alternative text            |
| `border`   | `1`            | Border width                |
| `bg`  | `#f0f0f0`      | Background container        |
