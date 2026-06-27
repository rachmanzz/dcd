# DCD — Document Compilation Definition

DCD compiles `.dcd` files into **DOCX** documents using a simple INI-style format with HTML-like template tags.

## How It Works

```
.dcd file → Parser → Compiler → Renderer → .docx
              ↑
          DataSet (variables)
```

1. **Parser** reads the `.dcd` file, splitting it into `[section]` blocks with properties and body content
2. **Compiler** processes each section: expands loops, resolves variables, applies formatting, then dispatches body tags
3. **Renderer** (DOCX) converts each tag into the output format

## Quick Start

```bash
go install github.com/rachmanzz/dcd@latest

dcd input.dcd output.docx
dcd --data variables.json input.dcd output.docx
```

## Project Structure

| Path | Description |
|------|-------------|
| `cmd/dcd/main.go` | CLI entry point |
| `parse/` | `.dcd` file parser |
| `render/` | Compiler, renderer (DOCX), style helpers |
| `data/` | Variable data set for `{{key}}` resolution |
| `.agents/skills/` | Detailed skill specifications |

## Document Structure

A `.dcd` file contains sections:

```ini
[style]
layout=A4
unit=mm
m=20
color=#000000

[title]
title=My Report
subject=Report Summary
author=John Doe

[header]
left=Confidential
right={{date}}

[section 0]
name=cover
keys=date_field
formats=[date_field:dd-MM-yyyy]

--- BODY ---
<h1>{{title}}</h1>
<p>Date: <b>{{date_field}}</b></p>
<p>This document uses <set:b|i>combined formatting</set:b|i></p>
```

## Key Features

### Inline Formatting
- Single tags: `<b>`, `<i>`, `<u>`, `<code>`
- Combined formatting: `<set:b|i>text</set:b|i>`
- Works in paragraphs, lists, and tables

### Dynamic Content
- Variables: `{{var.field}}`
- Loops: `<loop:row x from items>{{x.name}}</loop:row>`
- Loop styling: `<loop:row style.first=header x from data>`
- Dynamic styles: `<row style={{myStyle}}>`

### Table Styling
- Named styles: `[style:table header]`
- Properties: `bg`, `color`, `font-weight`, `align`
- Inline attributes: `<row bg=#f0f0f0>`

### Document Metadata
- Title, subject, author: `[title]` section
- Built-in variables: `{{title}}`, `{{date}}`, `{{page}}`, `{{total}}`
- Header/footer support with `left`, `center`, `right`, or `justify_between`

## Breaking Changes (v0.2.0)

If migrating from older versions:
- `font-color` → `color`
- `shading` → `bg`
- `[table-style]` → `[style:table]`

See `CHANGES.md` for migration guide.
