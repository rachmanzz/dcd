# Known Limitations

## 1. Nested Lists Not Supported

**Status:** Not implemented

**Description:**
Lists within lists (nested lists) are not currently supported.

**Not working:**
```html
<ul>
  <li>First item</li>
  <li>Second item
    <ul>
      <li>Nested item A</li>
      <li>Nested item B</li>
    </ul>
  </li>
  <li>Third item</li>
</ul>
```

**Behavior:**
- Nested `<ul>` or `<ol>` tags are stripped from the content
- Only the outer list is rendered
- Inner list items are ignored

**Workaround:**
Use separate flat lists:

```html
<ul>
  <li>First item</li>
  <li>Second item</li>
  <li>Third item</li>
</ul>

<p>Sub-items:</p>
<ul>
  <li>Nested item A</li>
  <li>Nested item B</li>
</ul>
```

**Future:**
This feature may be implemented in a future version.

---

## 2. Nested Inline Tags

**Status:** Alternative solution provided

**Description:**
Original nested inline tags like `<b><i>text</i></b>` are not supported due to parsing complexity.

**Solution:**
Use the `<set:flags>` tag instead:

```html
<!-- Instead of: <b><i>text</i></b> -->
<p><set:b|i>Bold and Italic</set:b|i></p>
```

This provides the same functionality with cleaner syntax.

---

## 3. PDF Inline Formatting

**Status:** Limitation of underlying library

**Description:**
PDF output does not support rich inline formatting in lists and tables due to limitations in the `gopdf` library.

**Behavior:**
- DOCX: Full inline formatting support (bold, italic, underline, code, combinations)
- PDF: Plain text only in lists and tables (formatting stripped)

**Affected:**
- `<li>` content with inline tags
- `<col>` content with inline tags

**Workaround:**
For documents requiring rich formatting in lists/tables, use DOCX output format.

---

## Summary

| Feature | Status | Impact | Workaround Available |
|---------|--------|--------|---------------------|
| Nested lists | Not supported | Medium | Yes - use flat lists |
| Nested inline tags | Alternative provided | Low | Yes - use `<set:>` |
| PDF inline formatting | Library limitation | Low | Use DOCX format |
