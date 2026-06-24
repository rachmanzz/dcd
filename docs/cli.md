# CLI Reference

## Usage

```bash
dcd [--format FORMAT] [--data FILE] <input.dcd> [output]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `input.dcd` | Input document file (required) |
| `output` | Output file path (optional, defaults to `output.docx` or `output.pdf`) |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `docx` | Output format: `docx` or `pdf` |
| `--data` | `""` | JSON file with variables |

## Examples

```bash
# Default DOCX output
dcd report.dcd

# Specify output path
dcd report.dcd report.docx

# PDF output
dcd --format pdf report.dcd report.pdf

# PDF with default name
dcd --format pdf report.dcd
# → output.pdf

# With JSON variables
dcd --data report.json report.dcd report.docx
```

Uses `{{info.title}}`, `{{info.author}}`, etc. from the JSON file.
```json
{
  "info": { "title": "Quarterly Report", "author": "Finance Team" },
  "items": [
    { "name": "Revenue", "value": 2100000 }
  ]
}
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (parse error, missing input, unknown format) |

## Common Workflows

### Generate Invoice from Data

```bash
dcd --data invoice.json invoice.dcd invoice.docx
```

**invoice.json:**
```json
{
  "invoice": {
    "number": "INV-2025-001",
    "items": [
      {"id": 1, "desc": "Service A", "qty": 10, "price": 100, "total": 1000}
    ]
  }
}
```

**invoice.dcd:**
```ini
[style:table header]
bg=#4472C4
color=#ffffff

[section 0]
var=invoice
keys=number,items

--- BODY ---
<h1>Invoice {{invoice.number}}</h1>
<table border>
  <loop:row style.first=header x from invoice.items>
    <col>{{x.desc}}</col>
    <col align=right>${{x.total}}</col>
  </loop:row>
</table>
```

### Batch Processing

```bash
# Process multiple files
for file in reports/*.dcd; do
  dcd --data data.json "$file" "output/$(basename "$file" .dcd).docx"
done

# Generate both DOCX and PDF
dcd report.dcd report.docx
dcd --format pdf report.dcd report.pdf
```

### Environment Variables

```bash
# Use environment for data path
export DATA_FILE=config.json
dcd --data "$DATA_FILE" template.dcd output.docx
```

## Troubleshooting

### Variable Not Resolved

If `{{var.field}}` appears in output:
1. Check `var=` is set in section
2. Check `keys=` includes the field
3. Verify JSON data structure matches
4. Ensure data file is loaded with `--data`

### Loop Not Expanding

If `<loop:row x from items>` doesn't work:
1. Check data source exists in JSON
2. Use dot notation for nested data: `from invoice.items`
3. Verify array is not empty

### Style Not Applied

If `[style:table header]` doesn't work:
1. Check section name format (v0.2.0+): `[style:table name]`
2. Use `style=name` in `<row>` tag
3. Verify property names: `bg` not `shading`, `color` not `font-color`

## Version Compatibility

### v0.2.0 (Current)

Breaking changes:
- Property renames: `font-color` → `color`, `shading` → `bg`
- Section format: `[table-style]` → `[style:table]`

New features:
- `<set:flags>` tag for combined formatting
- `style.first` for loops
- `style={{var}}` for dynamic styles
- Loop dot notation support

See `CHANGES.md` for migration guide.
