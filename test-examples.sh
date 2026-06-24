#!/bin/bash
# Test all examples in docs/examples/

set -e

echo "Building dcd..."
go build -o dcd ./cmd/dcd

echo ""
echo "Testing examples..."
echo ""

# Simple (no data)
echo "→ simple.dcd → simple.docx"
./dcd docs/examples/simple.dcd docs/examples/simple.docx

# Features (no data)
echo "→ features.dcd → features.docx"
./dcd docs/examples/features.dcd docs/examples/features.docx

# Report (with JSON)
echo "→ report.dcd + report.json → report.docx"
./dcd --data docs/examples/report.json docs/examples/report.dcd docs/examples/report.docx

# Invoice (with JSON)
echo "→ invoice.dcd + invoice.json → invoice.docx"
./dcd --data docs/examples/invoice.json docs/examples/invoice.dcd docs/examples/invoice.docx

# PDF examples
echo "→ simple.dcd → simple.pdf"
./dcd --format pdf docs/examples/simple.dcd docs/examples/simple.pdf

echo "→ report.dcd + report.json → report.pdf"
./dcd --format pdf --data docs/examples/report.json docs/examples/report.dcd docs/examples/report.pdf

echo ""
echo "✓ All examples generated successfully"
