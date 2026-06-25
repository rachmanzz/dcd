---
name: golang-programming
description: Complete guide to using DCD as a Go library — installation, API reference, data types, renderer interface, and advanced patterns
---

# Programmatic Usage

Complete guide to using DCD as a Go library in your applications.

## Installation

```bash
go get github.com/rachmanzz/dcd
```

## Import Packages

```go
import (
    "github.com/rachmanzz/dcd/data"
    "github.com/rachmanzz/dcd/parse"
    "github.com/rachmanzz/dcd/render"
)
```

## Quick Start

### Basic Document Generation

```go
package main

import (
    "log"

    "github.com/rachmanzz/dcd/data"
    "github.com/rachmanzz/dcd/parse"
    "github.com/rachmanzz/dcd/render"
)

func main() {
    // Parse .dcd file
    doc, err := parse.Parse("template.dcd")
    if err != nil {
        log.Fatal(err)
    }

    // Create data set
    ds := data.NewDataSet(map[string]any{
        "title": "My Report",
        "author": "John Doe",
    })

    // Create renderer
    renderer := render.NewDocxRenderer()

    // Compile and generate
    compiler := render.New(doc, ds, renderer)
    if err := compiler.Run("output.docx"); err != nil {
        log.Fatal(err)
    }

    log.Println("Generated: output.docx")
}
```

## Without .dcd Files

Generate documents programmatically without template files:

```go
package main

import (
    "github.com/rachmanzz/dcd/render"
)

func main() {
    r := render.NewDocxRenderer()

    // Set page style
    r.SetPageStyle(map[string]string{
        "layout": "A4",
        "m":      "25.4", // margins in mm
    })

    // Set default text style
    r.SetDefaultStyle(map[string]string{
        "font-family": "Arial",
        "font-size":   "12",
        "color":       "#000000",
    })

    // Add heading
    r.AddHeading("My Report", 1, nil)

    // Add paragraph
    r.AddParagraph([]render.TextRun{
        {Text: "This is "},
        {Text: "bold", Bold: true},
        {Text: " and "},
        {Text: "italic", Italic: true},
        {Text: " text."},
    })

    // Save
    if err := r.Save("output.docx"); err != nil {
        log.Fatal(err)
    }
}
```

## Working with Data

### DataSet Creation

```go
// From map
ds := data.NewDataSet(map[string]any{
    "name":  "John",
    "age":   30,
    "items": []any{
        map[string]any{"id": 1, "title": "Item A"},
        map[string]any{"id": 2, "title": "Item B"},
    },
})

// From struct
type Report struct {
    Title  string
    Author string
    Items  []Item
}

type Item struct {
    ID    int
    Title string
}

report := Report{
    Title:  "Annual Report",
    Author: "Finance Team",
    Items: []Item{
        {ID: 1, Title: "Item A"},
        {ID: 2, Title: "Item B"},
    },
}

ds := data.NewDataSet(report)
```

### DataSet Operations

```go
// Set values
ds.Set("name", "Alice")
ds.Set("users", []any{
    map[string]any{"name": "Bob", "age": 25},
    map[string]any{"name": "Carol", "age": 30},
})

// Get values
name, ok := ds.Get("name")
if ok {
    fmt.Println(name) // "Alice"
}

// Resolve template strings
result := ds.Resolve("Hello {{name}}!") // "Hello Alice!"
result = ds.Resolve("User: {{users.0.name}}") // "User: Bob"
```

## Renderer Interface

All renderers implement this interface:

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

### TextRun

For rich text formatting:

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

// Example
runs := []render.TextRun{
    {Text: "Normal "},
    {Text: "bold", Bold: true},
    {Text: " and "},
    {Text: "italic", Italic: true},
    {Text: " and "},
    {Text: "both", Bold: true, Italic: true},
}

r.AddParagraph(runs)
```

### ListItem

For lists with rich formatting:

```go
type ListItem struct {
    Runs  []TextRun
    Items []ListItem  // For nested lists (not supported in v0.2.0)
}

// Example
items := []render.ListItem{
    {Runs: []render.TextRun{{Text: "Plain item"}}},
    {Runs: []render.TextRun{
        {Text: "Item with "},
        {Text: "bold", Bold: true},
        {Text: " text"},
    }},
    {Runs: []render.TextRun{{Text: "Code item", Code: true}}},
}

r.AddList(items, false) // false = unordered list
```

### TableCell & TableRow

For tables with rich formatting:

```go
type TableCell struct {
    Runs  []TextRun
    Attrs map[string]string
}

type TableRow struct {
    Cells []TableCell
    Props map[string]string
}

// Example
rows := []render.TableRow{
    {
        Cells: []render.TableCell{
            {
                Runs: []render.TextRun{{Text: "Header 1", Bold: true}},
                Attrs: map[string]string{"align": "center"},
            },
            {
                Runs: []render.TextRun{{Text: "Header 2", Bold: true}},
            },
        },
        Props: map[string]string{"style": "header"},
    },
    {
        Cells: []render.TableCell{
            {Runs: []render.TextRun{{Text: "Data 1"}}},
            {Runs: []render.TextRun{{Text: "Data 2", Italic: true}}},
        },
    },
}

r.AddTable(rows, map[string]string{"border": "1"})
```

## Complete Examples

### Invoice Generator

```go
package main

import (
    "log"
    "time"

    "github.com/rachmanzz/dcd/render"
)

type Invoice struct {
    Number   string
    Date     string
    Customer string
    Items    []InvoiceItem
    Total    float64
}

type InvoiceItem struct {
    Description string
    Quantity    int
    Price       float64
    Total       float64
}

func GenerateInvoice(inv Invoice, outputPath string) error {
    r := render.NewDocxRenderer()

    // Page setup
    r.SetPageStyle(map[string]string{
        "layout": "A4",
        "m":      "20",
    })

    // Header style
    r.SetHeadingStyle(1, map[string]string{
        "font-size": "24",
        "color":     "#1F3864",
        "bold":      "true",
    })

    // Table header style
    r.SetTableStyle("header", map[string]string{
        "bg":          "#4472C4",
        "color":       "#ffffff",
        "font-weight": "bold",
        "align":       "center",
    })

    // Header
    r.AddHeading("INVOICE", 1, nil)

    // Invoice details
    r.AddParagraph([]render.TextRun{
        {Text: "Invoice #: ", Bold: true},
        {Text: inv.Number},
    })
    r.AddParagraph([]render.TextRun{
        {Text: "Date: ", Bold: true},
        {Text: inv.Date},
    })
    r.AddParagraph([]render.TextRun{
        {Text: "Customer: ", Bold: true},
        {Text: inv.Customer},
    })

    r.AddLineBreak()

    // Items table
    rows := []render.TableRow{
        {
            Cells: []render.TableCell{
                {Runs: []render.TextRun{{Text: "Description"}}},
                {Runs: []render.TextRun{{Text: "Qty"}}, Attrs: map[string]string{"align": "center"}},
                {Runs: []render.TextRun{{Text: "Price"}}, Attrs: map[string]string{"align": "right"}},
                {Runs: []render.TextRun{{Text: "Total"}}, Attrs: map[string]string{"align": "right"}},
            },
            Props: map[string]string{"style": "header"},
        },
    }

    for _, item := range inv.Items {
        rows = append(rows, render.TableRow{
            Cells: []render.TableCell{
                {Runs: []render.TextRun{{Text: item.Description}}},
                {Runs: []render.TextRun{{Text: fmt.Sprintf("%d", item.Quantity)}}, Attrs: map[string]string{"align": "center"}},
                {Runs: []render.TextRun{{Text: fmt.Sprintf("$%.2f", item.Price)}}, Attrs: map[string]string{"align": "right"}},
                {Runs: []render.TextRun{{Text: fmt.Sprintf("$%.2f", item.Total)}}, Attrs: map[string]string{"align": "right"}},
            },
        })
    }

    rows = append(rows, render.TableRow{
        Cells: []render.TableCell{
            {Runs: []render.TextRun{{Text: ""}}},
            {Runs: []render.TextRun{{Text: ""}}},
            {Runs: []render.TextRun{{Text: "TOTAL", Bold: true}}, Attrs: map[string]string{"align": "right"}},
            {Runs: []render.TextRun{{Text: fmt.Sprintf("$%.2f", inv.Total), Bold: true}}, Attrs: map[string]string{"align": "right"}},
        },
    })

    r.AddTable(rows, map[string]string{"border": "1"})

    return r.Save(outputPath)
}

func main() {
    invoice := Invoice{
        Number:   "INV-2025-001",
        Date:     time.Now().Format("2006-01-02"),
        Customer: "Acme Corporation",
        Items: []InvoiceItem{
            {Description: "Consulting Services", Quantity: 10, Price: 150, Total: 1500},
            {Description: "Support Services", Quantity: 5, Price: 100, Total: 500},
        },
        Total: 2000,
    }

    if err := GenerateInvoice(invoice, "invoice.docx"); err != nil {
        log.Fatal(err)
    }

    log.Println("Invoice generated: invoice.docx")
}
```

### Report Generator with Templates

```go
package main

import (
    "encoding/json"
    "log"
    "os"

    "github.com/rachmanzz/dcd/data"
    "github.com/rachmanzz/dcd/parse"
    "github.com/rachmanzz/dcd/render"
)

type Report struct {
    Title    string    `json:"title"`
    Author   string   `json:"author"`
    Date     string   `json:"date"`
    Sections []Section `json:"sections"`
}

type Section struct {
    Heading string `json:"heading"`
    Content string `json:"content"`
}

func GenerateReport(templatePath, dataPath, outputPath string) error {
    doc, err := parse.Parse(templatePath)
    if err != nil {
        return err
    }

    dataFile, err := os.ReadFile(dataPath)
    if err != nil {
        return err
    }

    var report Report
    if err := json.Unmarshal(dataFile, &report); err != nil {
        return err
    }

    ds := data.NewDataSet(map[string]any{
        "report": report,
    })

    renderer := render.NewDocxRenderer()
    compiler := render.New(doc, ds, renderer)
    return compiler.Run(outputPath)
}

func main() {
    if err := GenerateReport(
        "templates/report.dcd",
        "data/report.json",
        "output/report.docx",
    ); err != nil {
        log.Fatal(err)
    }

    log.Println("Report generated successfully")
}
```

### Multi-Format Generator

```go
package main

import (
    "log"

    "github.com/rachmanzz/dcd/data"
    "github.com/rachmanzz/dcd/parse"
    "github.com/rachmanzz/dcd/render"
)

func GenerateDocument(templatePath string, dataset *data.DataSet, outputPath string, format string) error {
    doc, err := parse.Parse(templatePath)
    if err != nil {
        return err
    }

    var renderer render.Renderer
    switch format {
    case "pdf":
        renderer = render.NewPdfRenderer()
    case "docx":
        renderer = render.NewDocxRenderer()
    default:
        return fmt.Errorf("unsupported format: %s", format)
    }

    compiler := render.New(doc, dataset, renderer)
    return compiler.Run(outputPath)
}

func main() {
    ds := data.NewDataSet(map[string]any{
        "title":  "My Document",
        "author": "John Doe",
    })

    if err := GenerateDocument("template.dcd", ds, "output.docx", "docx"); err != nil {
        log.Fatal(err)
    }

    if err := GenerateDocument("template.dcd", ds, "output.pdf", "pdf"); err != nil {
        log.Fatal(err)
    }

    log.Println("Documents generated successfully")
}
```

## Advanced Patterns

### Custom Data Loader

```go
type DataLoader interface {
    Load() (map[string]any, error)
}

type JSONLoader struct {
    Path string
}

func (l *JSONLoader) Load() (map[string]any, error) {
    data, err := os.ReadFile(l.Path)
    if err != nil {
        return nil, err
    }

    var result map[string]any
    if err := json.Unmarshal(data, &result); err != nil {
        return nil, err
    }

    return result, nil
}

type DatabaseLoader struct {
    DB    *sql.DB
    Query string
}

func (l *DatabaseLoader) Load() (map[string]any, error) {
    rows, err := l.DB.Query(l.Query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []map[string]any
    // ... scan rows ...

    return map[string]any{"items": items}, nil
}
```

### Document Builder Pattern

```go
type DocumentBuilder struct {
    renderer render.Renderer
}

func NewDocumentBuilder() *DocumentBuilder {
    return &DocumentBuilder{
        renderer: render.NewDocxRenderer(),
    }
}

func (b *DocumentBuilder) SetPageLayout(layout string) *DocumentBuilder {
    b.renderer.SetPageStyle(map[string]string{"layout": layout})
    return b
}

func (b *DocumentBuilder) AddTitle(text string) *DocumentBuilder {
    b.renderer.AddHeading(text, 1, nil)
    return b
}

func (b *DocumentBuilder) AddSection(heading, content string) *DocumentBuilder {
    b.renderer.AddHeading(heading, 2, nil)
    b.renderer.AddParagraph([]render.TextRun{{Text: content}})
    return b
}

func (b *DocumentBuilder) Save(path string) error {
    return b.renderer.Save(path)
}

// Usage
func main() {
    err := NewDocumentBuilder().
        SetPageLayout("A4").
        AddTitle("My Report").
        AddSection("Introduction", "This is the introduction.").
        AddSection("Conclusion", "This is the conclusion.").
        Save("report.docx")

    if err != nil {
        log.Fatal(err)
    }
}
```

### Template Registry

```go
type TemplateRegistry struct {
    templates map[string]*parse.Doc
}

func NewTemplateRegistry() *TemplateRegistry {
    return &TemplateRegistry{
        templates: make(map[string]*parse.Doc),
    }
}

func (r *TemplateRegistry) Register(name, path string) error {
    doc, err := parse.Parse(path)
    if err != nil {
        return err
    }
    r.templates[name] = doc
    return nil
}

func (r *TemplateRegistry) Generate(name string, ds *data.DataSet, outputPath string) error {
    doc, ok := r.templates[name]
    if !ok {
        return fmt.Errorf("template not found: %s", name)
    }

    renderer := render.NewDocxRenderer()
    compiler := render.New(doc, ds, renderer)
    return compiler.Run(outputPath)
}

// Usage
func main() {
    registry := NewTemplateRegistry()
    registry.Register("invoice", "templates/invoice.dcd")
    registry.Register("report", "templates/report.dcd")

    ds := data.NewDataSet(map[string]any{"number": "INV-001"})

    if err := registry.Generate("invoice", ds, "invoice.docx"); err != nil {
        log.Fatal(err)
    }
}
```

## Error Handling

### Best Practices

```go
func GenerateDocument(templatePath, dataPath, outputPath string) error {
    doc, err := parse.Parse(templatePath)
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }

    data, err := loadData(dataPath)
    if err != nil {
        return fmt.Errorf("failed to load data: %w", err)
    }

    ds := data.NewDataSet(data)
    renderer := render.NewDocxRenderer()
    compiler := render.New(doc, ds, renderer)

    if err := compiler.Run(outputPath); err != nil {
        return fmt.Errorf("failed to generate document: %w", err)
    }

    return nil
}

func main() {
    if err := GenerateDocument("template.dcd", "data.json", "output.docx"); err != nil {
        log.Printf("Error: %v", err)
        os.Exit(1)
    }
}
```

## Performance Tips

### Reuse Parsers

```go
// Bad: Parse on every request
func handleRequest(w http.ResponseWriter, r *http.Request) {
    doc, _ := parse.Parse("template.dcd") // Slow!
    // ... generate ...
}

// Good: Parse once, reuse
var templateDoc *parse.Doc

func init() {
    var err error
    templateDoc, err = parse.Parse("template.dcd")
    if err != nil {
        log.Fatal(err)
    }
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Use cached templateDoc
    // ... generate ...
}
```

### Concurrent Generation

```go
func GenerateMultiple(templates []Template, output string) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(templates))

    for i, tmpl := range templates {
        wg.Add(1)
        go func(index int, t Template) {
            defer wg.Done()

            outputPath := fmt.Sprintf("%s/doc-%d.docx", output, index)
            if err := generateOne(t, outputPath); err != nil {
                errChan <- err
            }
        }(i, tmpl)
    }

    wg.Wait()
    close(errChan)

    for err := range errChan {
        if err != nil {
            return err
        }
    }

    return nil
}
```

## Testing

### Unit Testing

```go
package myapp

import (
    "testing"

    "github.com/rachmanzz/dcd/data"
    "github.com/rachmanzz/dcd/render"
)

func TestDocumentGeneration(t *testing.T) {
    ds := data.NewDataSet(map[string]any{
        "title": "Test Report",
    })

    r := render.NewDocxRenderer()

    if err := r.Save("/tmp/test.docx"); err != nil {
        t.Fatalf("Failed to generate document: %v", err)
    }

    if _, err := os.Stat("/tmp/test.docx"); os.IsNotExist(err) {
        t.Error("Output file was not created")
    }

    os.Remove("/tmp/test.docx")
}
```

### Integration Testing

```go
func TestEndToEnd(t *testing.T) {
    templatePath := "/tmp/test-template.dcd"
    dataPath := "/tmp/test-data.json"
    outputPath := "/tmp/test-output.docx"

    template := `[section 0]
var=data
keys=title

--- BODY ---
<h1>{{data.title}}</h1>`
    os.WriteFile(templatePath, []byte(template), 0644)

    data := `{"data": {"title": "Test Document"}}`
    os.WriteFile(dataPath, []byte(data), 0644)

    err := GenerateDocument(templatePath, dataPath, outputPath)
    if err != nil {
        t.Fatalf("Failed: %v", err)
    }

    if _, err := os.Stat(outputPath); os.IsNotExist(err) {
        t.Error("Output not created")
    }

    os.Remove(templatePath)
    os.Remove(dataPath)
    os.Remove(outputPath)
}
```

## See Also

- `dcd-cli` — CLI usage and options
- `dcd-documents` — Document template syntax reference
- `dcd-guide` — Project overview and patterns
