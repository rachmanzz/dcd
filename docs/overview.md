# DCD — Document Compilation Definition

DCD compiles `.dcd` files into **DOCX** or **PDF** documents using a simple INI-style format with HTML-like template tags.

## How It Works

```
.dcd file → Parser → Compiler → Renderer → .docx / .pdf
              ↑
          DataSet (variables)
```

1. **Parser** reads the `.dcd` file, splitting it into `[section]` blocks with properties and body content
2. **Compiler** processes each section: expands loops, resolves variables, applies formatting, then dispatches body tags
3. **Renderer** (DOCX or PDF) converts each tag into the output format

## Quick Start

```bash
go install github.com/rachmanzz/dcd@latest

dcd input.dcd output.docx
dcd --format pdf input.dcd output.pdf
```

## Project Structure

| Path | Description |
|------|-------------|
| `cmd/dcd/main.go` | CLI entry point |
| `parse/` | `.dcd` file parser |
| `render/` | Compiler, renderers (DOCX/PDF), style helpers |
| `data/` | Variable data set for `{{key}}` resolution |
| `.agents/skills/` | Detailed skill specifications |

## Document Structure

A `.dcd` file contains sections:

```ini
[style]
layout=A4
unit=inch
m=1

[title]
subject=Report
author=John

[header]
left=Confidential
right={{date}}

[section 0]
name=cover
formats=[date_field:02-01-2006]

--- BODY ---
<h1>Title</h1>
<p>Date: <b>{{date_field}}</b></p>
```
