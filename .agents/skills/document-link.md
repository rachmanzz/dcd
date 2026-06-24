# Document Link

Internal and external hyperlinks.

## Usage

From data section:

```
<section 0>
var=source
keys=url, label

--- BODY ---
<a={{source.url}}>{{source.label}}</a>
```

Static:

```
<a=https://example.com>visit website</a>
```

Inline:

```
<p>click <a={{source.url}} target=_blank>here</a> for more info</p>
```

## Properties

| Property    | Example         | Description          |
|-------------|-----------------|----------------------|
| `target`    | `_blank`        | Open in new tab      |
| `color`     | `#0055cc`       | Link color           |
| `underline` | `true`          | Underline            |

## Bookmark

```
<a=#chapter1>see Chapter 1</a>
```
