---
name: dcd-guide
description: Project overview and development guide for AI agents working with DCD — codebase structure, common tasks, testing patterns, and best practices
---

# For AI Agents - DCD Project Guide

This guide helps AI assistants (Claude, Gemini, OpenCode, etc.) understand and work effectively with the DCD (Document Compilation Description) project.

## Project Overview

**DCD** is a document compilation system that converts `.dcd` template files into DOCX documents.

**Core Concept:**
```
.dcd template + JSON data -> Parser -> Compiler -> Renderer -> .docx
```

**Language:** Go
**Output Formats:** DOCX
**Primary Use:** Automated document generation (invoices, reports, contracts, etc.)

## Project Structure

```
document-compilation/
├── cmd/dcd/          # CLI application
├── parse/            # .dcd file parser
├── render/           # Document renderer (DOCX)
│   ├── compiler.go   # Template compiler
│   ├── docx.go       # DOCX renderer

│   ├── body.go       # Body content parser
│   ├── style.go      # Style utilities
│   └── types.go      # Data types
├── data/             # Variable resolution
├── docs/             # User documentation
├── .agents/skills/   # AI agent skill files (you are here)
└── examples/         # Example .dcd files
```

## Key Concepts

### 1. DCD File Format

A `.dcd` file uses INI-style sections with HTML-like tags:

```ini
[style]
layout=A4
color=#000000

[style:table header]
bg=#4472C4
color=#ffffff
font-weight=bold

[section 0]
var=invoice
keys=number,items,total

--- BODY ---
<h1>Invoice #{{invoice.number}}</h1>
<table border>
  <loop:row x from invoice.items style.first=header>
    <col>{{x.desc}}</col>
    <col align=right>${{x.amount}}</col>
  </loop:row>
</table>
<p>Total: ${{invoice.total}}</p>
```

### 2. Property Names (v0.2.0)

**Current (v0.2.0):**
- `color` - Text color
- `bg` - Background color
- `[style:table name]` - Table style section

**Deprecated (v0.1.x):**
- ~~`font-color`~~ -> use `color`
- ~~`shading`~~ -> use `bg`
- ~~`[table-style name]`~~ -> use `[style:table name]`

### 3. Variable Resolution

Variables use `{{path.to.field}}` syntax:

```json
{
  "invoice": {
    "number": "INV-001",
    "items": [
      {"desc": "Service A", "amount": 1000}
    ],
    "total": 1000
  }
}
```

Access: `{{invoice.number}}`, `{{invoice.items.0.desc}}`

> **Variable Registration Rule:** Every `{{...}}` variable must be registered in the section's `keys` or `var`. Unregistered variables are treated as literal strings in the output. For array object fields, use dotted path notation (e.g. `items.price` in `keys`).

### 4. Key Features (v0.2.0)

**Combined Inline Formatting:**
```html
<p><set:b|i>Bold and Italic</set:b|i></p>
```

**Loop with First-Row Styling:**
```html
<loop:row x from items style.first=header>
  <col>{{x.name}}</col>
</loop:row>
```

**Dynamic Styles:**
```html
<row style={{myStyle}}>
  <col>Data</col>
</row>
```

**Loop Dot Notation:**
```html
<loop:row x from invoice.items>
  <col>{{x.id}}</col>
</loop:row>
```

## Common Tasks for AI Agents

### Task 1: Code Changes

When modifying code:

1. **Read first:** Always read the file before editing
2. **Understand context:** Check related files
3. **Follow patterns:** Match existing code style
4. **Test:** Build and run tests after changes

**Key files for features:**
- `render/body.go` - Body content parsing (tags, loops, inline)
- `render/compiler.go` - Template compilation
- `parse/parse.go` - .dcd file parsing
- `render/style.go` - Style utilities
- `render/types.go` - Data structures

### Task 2: Documentation Updates

When updating docs:

1. **Check all locations:** Update all relevant files
2. **Maintain consistency:** Keep terminology consistent
3. **Update examples:** Ensure examples work
4. **Cross-reference:** Update related docs

**Documentation locations:**
- `docs/*.md` - User documentation
- `.agents/skills/*.md` - AI agent skills
- `README.md` - Project overview
- `CHANGES.md` - Breaking changes

### Task 3: Adding New Features

**Standard workflow:**

1. **Plan:** Understand requirements, check existing code
2. **Implement:** Update code files
3. **Test:** Write/run tests
4. **Document:** Update all relevant documentation
5. **Examples:** Add/update example files
6. **Migrate:** Update existing examples if breaking change

**Checklist:**
- [ ] Code implemented
- [ ] Tests passing
- [ ] User docs updated
- [ ] Skill docs updated
- [ ] Examples updated/added
- [ ] Breaking changes documented
- [ ] Migration guide provided (if breaking)

### Task 4: Bug Fixes

**Standard workflow:**

1. **Reproduce:** Understand the issue
2. **Locate:** Find the problematic code
3. **Fix:** Make minimal changes
4. **Test:** Verify fix works
5. **Regression test:** Ensure no new issues
6. **Document:** Add to changelog if user-facing

### Task 5: Migration (v0.1.x -> v0.2.0)

**Property renames:**
```bash
sed -i 's/font-color=/color=/g; s/shading=/bg=/g; s/\[table-style /[style:table /g' *.dcd
```

**Code changes required:**
- Update all `.dcd` files
- Update documentation
- Update examples
- Test everything

## Understanding the Codebase

### Parser (`parse/parse.go`)

**Purpose:** Parse `.dcd` files into structured data

**Key functions:**
- `Parse(path string) (*Doc, error)` - Parse a .dcd file
- Splits file into sections
- Extracts properties and body content

**Data structure:**
```go
type Doc struct {
    Sections []Section
}

type Section struct {
    Name  string
    Props map[string]string
    Body  string
}
```

### Compiler (`render/compiler.go`)

**Purpose:** Process sections and coordinate rendering

**Key functions:**
- `New(doc *Doc, ds *DataSet, r Renderer) *Compiler`
- `Run(outputPath string) error`
- `renderSection(sec Section) error`
- `expandLoops(body string) string`

**Processing order:**
1. Apply metadata (`[title]`)
2. Apply header/footer
3. Apply styles
4. Apply table styles
5. Render sections in order

### Body Parser (`render/body.go`)

**Purpose:** Parse body content (HTML-like tags)

**Key functions:**
- `renderBody(body string) error` - Main entry point
- `splitInline(s string) []inlinePart` - Parse inline tags
- `expandLoops(body string) string` - Expand loop tags
- `parseAttrs(s string) map[string]string` - Parse attributes

**Handles:**
- Headings: `<h1>`, `<h2>`, etc.
- Paragraphs: `<p>`
- Tables: `<table>`, `<row>`, `<col>`
- Lists: `<ul>`, `<ol>`, `<li>`
- Loops: `<loop:row>`, `<loop:ol>`, `<loop:ul>`
- Inline: `<b>`, `<i>`, `<u>`, `<code>`, `<set:flags>`
- Images: `<img>`
- Links: `<a>`

### Renderer (`render/docx.go`)

**Purpose:** Generate actual document files

**Interface:**
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
    AddTable(rows []TableRow, attrs map[string]string) error
    AddList(items []ListItem, ordered bool) error
    AddImage(src string, attrs map[string]string) error
    AddLineBreak() error
    AddPageBreak() error

    Save(path string) error
}
```

### Data Resolution (`data/dataset.go`)

**Purpose:** Variable resolution from JSON data

**Key functions:**
- `NewDataSet(data any) *DataSet`
- `Set(key string, value any)`
- `Get(key string) (any, bool)`
- `Resolve(template string) string` - Resolve `{{var}}`

## Common Patterns

### Pattern 1: Adding a New Tag

**Example:** Adding `<highlight>` tag

1. **Add regex:**
```go
// In render/body.go
var highlightRe = regexp.MustCompile(`<highlight>(.*?)</highlight>`)
```

2. **Parse in renderBody:**
```go
if strings.HasPrefix(line, "<highlight>") {
    content := highlightRe.FindStringSubmatch(line)
    if len(content) > 1 {
        c.renderHighlight(content[1])
        continue
    }
}
```

3. **Implement renderer method:**
```go
func (c *Compiler) renderHighlight(text string) error {
    runs := inlineToRuns(text)
    for i := range runs {
        runs[i].Highlight = true
    }
    return c.r.AddParagraph(runs)
}
```

4. **Update TextRun type:**
```go
type TextRun struct {
    Text      string
    Bold      bool
    Italic    bool
    Underline bool
    Code      bool
    Highlight bool  // NEW
    Link      string
    LinkAttrs map[string]string
}
```

5. **Update renderers:**
```go
// In render/docx.go - handle runs[i].Highlight in AddParagraph
```

6. **Document:**
   - Update `docs/tags.md`
   - Update `dcd-documents` skill
   - Add example

7. **Test:**
```bash
go build ./...
./dcd test.dcd test.docx
```

### Pattern 2: Adding a Section Property

**Example:** Adding `indent` property

1. **Parser handles automatically** (all properties parsed)

2. **Use in renderer:**
```go
indent := sec.Props["indent"]
if indent != "" {
    // Apply indentation
}
```

3. **Document:**
   - Update `docs/style.md`
   - Add to property tables
   - Add examples

### Pattern 3: Adding Loop Variant

**Example:** Adding `<loop:div>`

1. **Update loop regex:**
```go
var loopRe = regexp.MustCompile(`<loop(?::(\w+))?\s+(?:style\.first=(\w+)\s+)?(\w+)\s+from\s+([\w.]+)(?:\s+style\.first=(\w+))?>`)
```

2. **Handle in expandLoops:**
```go
loopType := match[1] // "div"
switch loopType {
case "row":
    // ... existing
case "div":  // NEW
    output += "<div>" + rendered + "</div>"
}
```

3. **Document and test**

## Testing Strategy

### Manual Testing

```bash
# Build
go build -o dcd ./cmd/dcd

# Test with example
./dcd docs/examples/simple.dcd test.docx

# Test with data
./dcd --data docs/examples/invoice.json docs/examples/invoice.dcd test-invoice.docx

```

### Automated Testing

```bash
# Run test script
./test-examples.sh

# Build and vet
go build ./...
go vet ./...
```

### Verification Checklist

After making changes:

- [ ] `go build ./...` passes
- [ ] `go vet ./...` passes
- [ ] `./test-examples.sh` passes
- [ ] Manual test: `./dcd test.dcd test.docx` works
- [ ] Open generated .docx file and verify
- [ ] Documentation updated
- [ ] Examples updated if needed

## Breaking Changes Protocol

When introducing breaking changes:

1. **Document in CHANGES.md**
   - Clear description
   - Migration path
   - Automated script if possible

2. **Update all examples**
   - Migrate to new format
   - Test all examples

3. **Update all documentation**
   - User docs
   - Skill docs
   - API reference

4. **Version bump**
   - Major version for breaking changes
   - Update README with version

5. **Migration guide**
   - Before/after examples
   - Script or manual steps
   - Troubleshooting

## Common Pitfalls

### 1. Property Name Confusion

**Wrong:**
```ini
[style]
font-color=#000000
```

**Correct:**
```ini
[style]
color=#000000
```

### 2. Section Format

**Wrong:**
```ini
[table-style header]
```

**Correct:**
```ini
[style:table header]
```

### 3. Variable Access

**Wrong:**
```html
<p>{{items}}</p>  <!-- Array not expanded -->
```

**Correct:**
```html
<loop x from items>
  <p>{{x.name}}</p>
</loop>
```

### 4. Nested Lists

**Not Supported:**
```html
<ul>
  <li>Item
    <ul>
      <li>Nested</li>
    </ul>
  </li>
</ul>
```

**Use Flat Lists:**
```html
<ul>
  <li>Item</li>
  <li>Sub-item A</li>
  <li>Sub-item B</li>
</ul>
```

### 5. Inline Tags Scope

**Wrong:**
```html
<h1><b>Bold heading</b></h1>  <!-- Inline in block -->
```

**Correct:**
```html
<h1>Bold heading</h1>  <!-- Use style properties -->

<!-- OR in paragraphs/lists/tables -->
<p><b>Bold text</b></p>
<li><b>Bold item</b></li>
<col><b>Bold cell</b></col>
```

## Quick Reference

### File Locations

**For code changes:**
- Tags/parsing: `render/body.go`
- Compilation: `render/compiler.go`
- Styles: `render/style.go`
- Types: `render/types.go`
- Parser: `parse/parse.go`

**For documentation:**
- User docs: `docs/*.md`
- Skills: `.agents/skills/*.md`
- Changes: `CHANGES.md`
- Features: `NEW-FEATURES.md`
- Limitations: `KNOWN-LIMITATIONS.md`

**For testing:**
- Examples: `docs/examples/*.dcd`
- Test script: `./test-examples.sh`
- Build: `go build ./...`

### Important Regex Patterns

```go
// Inline tags
bRe := regexp.MustCompile(`<b>(.*?)</b>`)
setRe := regexp.MustCompile(`<set:([^>]+)>(.*?)</set(?::[^>]+)?>`)

// Loops
loopRe := regexp.MustCompile(`<loop(?::(\w+))?\s+(?:style\.first=(\w+)\s+)?(\w+)\s+from\s+([\w.]+)(?:\s+style\.first=(\w+))?>`)

// Variables
varRe := regexp.MustCompile(`\{\{([^}]+)\}\}`)
```

### Property Normalization

User input -> Internal representation:
- `color` -> `font-color`
- `bg` -> `shading`

Function: `normalizePropertyKey()` in `parse/parse.go` and `render/style.go`

## AI Agent Best Practices

### When Helping Users

1. **Check version:** Ask what version they're using
2. **Check syntax:** Ensure they use v0.2.0 syntax
3. **Verify examples:** Test examples before sharing
4. **Reference docs:** Point to specific documentation
5. **Provide complete examples:** Full working code

### When Modifying Code

1. **Read first:** Always read files before editing
2. **Understand context:** Check related code
3. **Follow patterns:** Match existing style
4. **Test thoroughly:** Run all tests
5. **Document:** Update relevant docs

### When Diagnosing Issues

1. **Reproduce:** Try to reproduce the issue
2. **Check version:** Verify version compatibility
3. **Check syntax:** Look for common mistakes
4. **Test minimal:** Create minimal test case
5. **Verify fix:** Test the solution

## Resources

**Quick Links:**
- Main README: `../README.md`
- User docs: `../docs/`
- Examples: `../docs/examples/`
- Skill files: `./*.md` (this directory)

**Essential Skills:**
- `dcd-cli` — CLI usage
- `golang-programming` — Go library API
- `dcd-documents` — Document template syntax

**Version Info:**
- Current: v0.2.0
- Breaking changes: `../CHANGES.md`
- New features: `../NEW-FEATURES.md`
- Limitations: `../KNOWN-LIMITATIONS.md`

## Summary

**For AI Agents working with DCD:**

1. **Understand the format:** INI-style sections + HTML-like tags
2. **Know v0.2.0 changes:** Property renames, new features
3. **Follow patterns:** Match existing code style
4. **Test everything:** Build, examples, manual verification
5. **Document thoroughly:** Update all relevant files
6. **Be version-aware:** Check user's version, provide correct syntax

**This project is well-documented. When in doubt, check the docs!**

---

**Version:** v0.2.0
**Last Updated:** Session complete
**Status:** Production ready

## See Also

- `dcd-cli` — CLI usage and options
- `dcd-documents` — Document template syntax reference
- `golang-programming` — Go library API
