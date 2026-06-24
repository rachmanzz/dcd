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
    Link      string
    LinkAttrs map[string]string
}

type TableCell struct {
    Text  string
    Attrs map[string]string
}

type TableRow struct {
    Cells []TableCell
    Props map[string]string
}

type ListItem struct {
    Text  string
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
        {Text: "Name", Attrs: map[string]string{"align": "center"}},
        {Text: "Value"},
    }, Props: map[string]string{"style": "header"}},
    {Cells: []render.TableCell{
        {Text: "Alpha"},
        {Text: "42"},
    }},
}, map[string]string{"border": "true"})

r.Save("output.docx")
```
