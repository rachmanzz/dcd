# Library API

Import DCD as a Go library:

```go
import (
    "github.com/rachmanzz/dcd/data"
    "github.com/rachmanzz/dcd/parse"
    "github.com/rachmanzz/dcd/render"
)
```

## Basic Usage

```go
// Parse a .dcd file
doc, err := parse.Parse("input.dcd")

// Create a data set with variables
ds := data.NewDataSet(map[string]any{
    "name":   "John",
    "items":  []any{map[string]any{"title": "A", "val": 1}},
})

// Choose a renderer
r := render.NewDocxRenderer()  // or render.NewPdfRenderer()

// Compile and save
err = render.New(doc, ds, r).Run("output.docx")
```

## DataSet

```go
// From map
ds := data.NewDataSet(map[string]any{"key": "value"})

// From struct
ds := data.NewDataSet(MyStruct{...})

// From map[string]string
ds := data.NewDataSet(map[string]string{"key": "value"})

// Set values programmatically
ds.Set("name", "Alice")
ds.Set("items", []any{map[string]any{"id": 1, "label": "Item 1"}})

// Resolve {{variables}} in a template string
result := ds.Resolve("Hello {{name}}")  // → "Hello Alice"
```

Variables are accessed with dot notation: `{{info.title}}`, `{{x.name}}`.

## Renderer Interface

```go
type Renderer interface {
    SetPageStyle(props map[string]string) error
    SetDefaultStyle(props map[string]string) error
    SetHeadingStyle(level int, props map[string]string) error
    SetTableStyle(name string, props map[string]string) error
    SetHeader(props map[string]string) error
    SetFooter(props map[string]string) error
    SetMetadata(props map[string]string) error
    AddHeading(text string, level int, attrs map[string]string) error
    AddParagraph(runs []TextRun) error
    AddWrappedParagraph(text string, flags string) error
    AddLineBreak() error
    AddHorizontalRule(attrs map[string]string) error
    AddPageBreak() error
    AddImage(src string, attrs map[string]string) error
    AddHyperlink(text, url string, attrs map[string]string) error
    AddList(items []ListItem, ordered bool) error
    AddTable(rows []TableRow, attrs map[string]string) error
    Save(path string) error
}
```

## Data Types

```go
type TextRun struct {
    Text      string
    Bold      bool
    Italic    bool
    Underline bool
    Code      bool
    Link      string
    LinkAttrs map[string]string
}

type TableCell struct {
    Runs  []TextRun
    Attrs map[string]string
}

type TableRow struct {
    Cells []TableCell
    Props map[string]string
}

type ListItem struct {
    Runs  []TextRun
    Items []ListItem
}
```

## Programmatic Usage (Without `.dcd` Files)

```go
r := render.NewDocxRenderer()
r.SetPageStyle(map[string]string{"layout": "A4", "m": "25.4"})
r.SetDefaultStyle(map[string]string{"font-family": "Arial", "font-size": "12"})

r.AddHeading("Chapter 1", 1, nil)
r.AddParagraph([]render.TextRun{
    {Text: "Hello "},
    {Text: "world", Bold: true},
})

// Table
r.AddTable([]render.TableRow{
    {Cells: []render.TableCell{
        {Runs: []render.TextRun{{Text: "Name", Bold: true}}, Attrs: map[string]string{"align": "center"}},
        {Runs: []render.TextRun{{Text: "Value"}}},
    }, Props: map[string]string{"style": "header"}},
    {Cells: []render.TableCell{
        {Runs: []render.TextRun{{Text: "Alpha"}}},
        {Runs: []render.TextRun{{Text: "42"}}},
    }},
}, map[string]string{"border": "true"})

r.Save("output.docx")
```

## Changes in v0.2.0

### Updated Data Types

**ListItem and TableCell now use `Runs` instead of `Text`:**

```go
// Old (v0.1.x)
type ListItem struct {
    Text  string
    Items []ListItem
}

type TableCell struct {
    Text  string
    Attrs map[string]string
}

// New (v0.2.0+)
type ListItem struct {
    Runs  []TextRun  // Changed: supports rich formatting
    Items []ListItem
}

type TableCell struct {
    Runs  []TextRun  // Changed: supports rich formatting
    Attrs map[string]string
}
```

**Migration:**

```go
// Old
cell := TableCell{Text: "Hello"}

// New
cell := TableCell{Runs: []TextRun{{Text: "Hello"}}}

// With formatting
cell := TableCell{Runs: []TextRun{
    {Text: "Bold", Bold: true},
    {Text: " and "},
    {Text: "Italic", Italic: true},
}}
```

### Property Name Changes

When setting styles programmatically, use new property names:

```go
// Old
r.SetDefaultStyle(map[string]string{
    "font-color": "#000000",
})
r.SetTableStyle("header", map[string]string{
    "shading": "#4472C4",
})

// New
r.SetDefaultStyle(map[string]string{
    "color": "#000000",  // Changed
})
r.SetTableStyle("header", map[string]string{
    "bg": "#4472C4",  // Changed
})
```

**Note:** The library automatically normalizes property names, so both old and new names work internally. However, using new names is recommended for consistency.

## Example: Rich Text in Tables

```go
r := render.NewDocxRenderer()

// Table with rich formatting
r.AddTable([]render.TableRow{
    {
        Cells: []render.TableCell{
            {Runs: []render.TextRun{
                {Text: "Product", Bold: true, Underline: true},
            }},
            {Runs: []render.TextRun{
                {Text: "Price", Bold: true},
            }},
        },
        Props: map[string]string{"style": "header"},
    },
    {
        Cells: []render.TableCell{
            {Runs: []render.TextRun{
                {Text: "Item A", Italic: true},
            }},
            {Runs: []render.TextRun{
                {Text: "$100", Code: true},
            }},
        },
    },
}, map[string]string{"border": "1"})

r.Save("output.docx")
```

## Example: Lists with Formatting

```go
r.AddList([]render.ListItem{
    {Runs: []render.TextRun{{Text: "Plain item"}}},
    {Runs: []render.TextRun{
        {Text: "Item with ", Bold: false},
        {Text: "bold", Bold: true},
        {Text: " text", Bold: false},
    }},
    {Runs: []render.TextRun{{Text: "Code item", Code: true}}},
}, false) // false = unordered list
```

## See Also

- `parse/` — Parser implementation
- `render/` — Renderer implementations (DOCX, PDF)
- `data/` — DataSet for variable resolution
- `.agents/skills/` — Detailed feature documentation
