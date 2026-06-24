# CLI Reference

## Usage

```bash
dcd [--format FORMAT] <input.dcd> [output]
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
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (parse error, missing input, unknown format) |
