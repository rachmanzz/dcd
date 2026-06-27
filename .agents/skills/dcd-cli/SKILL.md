---
name: dcd-cli
description: Complete guide to using the DCD command-line interface — installation, options, workflows, and troubleshooting
---

# CLI Usage

Complete guide to using the DCD command-line interface.

## Installation

```bash
# Install from source
go install github.com/rachmanzz/dcd/cmd/dcd@latest

# Or build locally
git clone https://github.com/rachmanzz/dcd.git
cd dcd
go build -o dcd ./cmd/dcd
```

## Basic Syntax

```bash
dcd [OPTIONS] <input.dcd> [output]
```

## Quick Start

```bash
# Generate DOCX (default)
dcd report.dcd

# Specify output path
dcd report.dcd my-report.docx



# With JSON data
dcd --data data.json invoice.dcd invoice.docx
```

## Command-Line Options

| Option | Short | Default | Description |
|--------|-------|---------|-------------|

| `--data` | `-d` | none | JSON file with variables |
| `--help` | `-h` | - | Show help message |
| `--version` | `-v` | - | Show version |

## Examples

### Simple Document

```bash
# Create report.dcd
cat > report.dcd << 'EOF'
[style]
layout=A4

[section 0]

--- BODY ---
<h1>My Report</h1>
<p>This is a simple report.</p>
EOF

# Generate DOCX
dcd report.dcd
# -> output.docx

# Generate with custom name
dcd report.dcd my-report.docx
```

### With Variables

**document.dcd:**
```ini
[section 0]
var=info
keys=title,author,date

--- BODY ---
<h1>{{info.title}}</h1>
<p>By: {{info.author}}</p>
<p>Date: {{info.date}}</p>
```

**data.json:**
```json
{
  "info": {
    "title": "Annual Report",
    "author": "John Doe",
    "date": "2025-01-15"
  }
}
```

**Generate:**
```bash
dcd --data data.json document.dcd report.docx
```

### Invoice with Loop

**invoice.dcd:**
```ini
[style:table header]
bg=#4472C4
color=#ffffff
font-weight=bold

[section 0]
var=invoice
keys=number,customer,items,total

--- BODY ---
<h1>Invoice #{{invoice.number}}</h1>
<p>Customer: {{invoice.customer}}</p>

<table border>
  <loop:row x from invoice.items style.first=header>
    <col>{{x.desc}}</col>
    <col align=right>${{x.amount}}</col>
  </loop:row>
  <row>
    <col align=right><b>Total:</b></col>
    <col align=right><b>${{invoice.total}}</b></col>
  </row>
</table>
```

**invoice.json:**
```json
{
  "invoice": {
    "number": "INV-001",
    "customer": "Acme Corp",
    "items": [
      {"desc": "Service A", "amount": 1000},
      {"desc": "Service B", "amount": 500}
    ],
    "total": 1500
  }
}
```

**Generate:**
```bash
dcd --data data.json invoice.dcd invoice.docx
```



## Workflows

### Batch Processing

```bash
# Process multiple files
for file in templates/*.dcd; do
  name=$(basename "$file" .dcd)
  dcd --data data.json "$file" "output/$name.docx"
done
```

### With Environment Variables

```bash
# Set data file via environment
export DATA_FILE=config.json

# Use in command
dcd --data "$DATA_FILE" template.dcd output.docx
```

### Pipeline Processing

```bash
# Generate from template
cat template.dcd | \
  sed "s/{{VERSION}}/v1.2.3/" > output.dcd

# Compile
dcd --data data.json output.dcd final.docx

# Clean up
rm output.dcd
```

### Automated Reports

```bash
#!/bin/bash
# generate-report.sh

DATE=$(date +%Y-%m-%d)
OUTPUT="report-$DATE.docx"

# Generate JSON data
cat > /tmp/data.json << EOF
{
  "report": {
    "date": "$DATE",
    "title": "Daily Report"
  }
}
EOF

# Generate report
dcd --data /tmp/data.json report.dcd "$OUTPUT"

echo "Generated: $OUTPUT"
```

## Troubleshooting

### Variable Not Resolved

**Problem:** `{{var.field}}` appears in output

**Solutions:**
1. Check `var=` is set in section
2. Verify `keys=` includes the field
3. Check JSON data structure matches
4. Ensure data file is loaded: `--data data.json`

**Debug:**
```bash
# Check JSON structure
cat data.json | jq .

# Verify .dcd file
cat document.dcd | grep -A 5 "\[section"
```

### Loop Not Expanding

**Problem:** `<loop:row x from items>` doesn't generate rows

**Solutions:**
1. Check data source exists in JSON
2. Use dot notation for nested: `from invoice.items`
3. Verify array is not empty
4. Check variable name matches JSON

**Debug:**
```bash
# Check array in JSON
cat data.json | jq '.invoice.items'

# Should return array
```

### Style Not Applied

**Problem:** `[style:table header]` doesn't work

**Solutions:**
1. Check section format (v0.2.0+): `[style:table name]`
2. Use `style=name` in `<row>` tag
3. Verify property names: `bg` not `shading`, `color` not `font-color`

**Migration (v0.1.x -> v0.2.0):**
```bash
sed -i 's/font-color=/color=/g; s/shading=/bg=/g; s/\[table-style /[style:table /g' *.dcd
```

### File Not Found

**Problem:** `parse: open report.dcd: no such file or directory`

**Solutions:**
1. Check file exists: `ls report.dcd`
2. Check current directory: `pwd`
3. Use absolute path: `dcd /full/path/to/report.dcd`

### Permission Denied

**Problem:** Cannot write output file

**Solutions:**
1. Check output directory exists
2. Check write permissions: `ls -la output/`
3. Create directory: `mkdir -p output/`

### JSON Parse Error

**Problem:** `invalid character` in JSON

**Solutions:**
1. Validate JSON: `cat data.json | jq .`
2. Check for trailing commas
3. Check quotes are correct (double quotes)
4. Use JSON validator online

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (parse error, file not found, invalid format, etc.) |

**Usage in scripts:**
```bash
if dcd report.dcd output.docx; then
  echo "Success!"
else
  echo "Failed with code $?"
  exit 1
fi
```

## Tips & Best Practices

### File Organization

```
project/
├── templates/
│   ├── invoice.dcd
│   └── report.dcd
├── data/
│   ├── invoice.json
│   └── report.json
└── output/
    ├── invoice.docx
    └── report.docx
```

### Naming Convention

```bash
# Template files
invoice-template.dcd
report-template.dcd

# Data files (same base name)
invoice-data.json
report-data.json

# Output files (with date)
invoice-2025-01-15.docx
report-2025-01-15.docx
```

### Version Control

```bash
# Track templates
git add templates/*.dcd

# Ignore generated files
echo "*.docx" >> .gitignore

echo "output/" >> .gitignore

# Track sample data
git add data/sample-*.json
```

### Testing

```bash
# Quick test with sample data
dcd --data data/sample.json template.dcd /tmp/test.docx

# Verify output
open /tmp/test.docx

# Clean up
rm /tmp/test.docx
```

## Advanced Usage

### Custom Output Names

```bash
# With timestamp
dcd report.dcd "report-$(date +%Y%m%d-%H%M%S).docx"

# With version
VERSION="v1.2.3"
dcd manual.dcd "manual-$VERSION.docx"

# From variable in data
CUSTOMER=$(jq -r '.customer.name' data.json)
dcd invoice.dcd "invoice-$CUSTOMER.docx"
```

### Parallel Processing

```bash
# Process multiple templates in parallel
find templates/ -name "*.dcd" | \
  parallel dcd --data data.json {} output/{/.}.docx
```

### Error Handling

```bash
#!/bin/bash
set -e  # Exit on error

trap 'echo "Error on line $LINENO"' ERR

dcd --data data.json template.dcd output.docx

echo "Success!"
```

## Integration Examples

### Make Integration

```makefile
# Makefile
all: report.docx invoice.docx

%.docx: templates/%.dcd data/%.json
	dcd --data data/$*.json $< output/$@

clean:
	rm -f output/*.docx

.PHONY: all clean
```

### CI/CD Integration

```yaml
# .github/workflows/generate-docs.yml
name: Generate Documents

on: [push]

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: '1.21'

      - name: Install DCD
        run: go install github.com/rachmanzz/dcd/cmd/dcd@latest

      - name: Generate Documents
        run: |
          dcd --data data.json template.dcd output.docx

      - name: Upload Artifacts
        uses: actions/upload-artifact@v2
        with:
          name: documents
          path: output/*.docx
```

## See Also

- `dcd-documents` — Document template syntax reference
- `golang-programming` — Go library API
- `dcd-guide` — Project overview and patterns
