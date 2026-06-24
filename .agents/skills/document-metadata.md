# Document Metadata

Set document properties like title, subject, and author using the `[title]` section.

## Section Format

```
[title]
title=Document Title
subject=Document Subject
author=Author Name
```

## Properties

| Property  | Description                          | Example                    |
|-----------|--------------------------------------|----------------------------|
| `title`   | Document title                       | Annual Report 2025         |
| `subject` | Document subject/description         | Financial Summary          |
| `author`  | Document author/creator              | Finance Team               |

## Usage

These properties are written to:
- **DOCX:** Document properties (`docProps/core.xml`)
- **PDF:** Document metadata (Title, Subject, Author fields)

## Built-in Variable: `{{title}}`

The `title` property can be referenced in headers and footers using the `{{title}}` variable.

```
[title]
title=My Report

[header]
left={{title}}
right={{date}}

[footer]
center={{title}} - Page {{page}}
```

## Example

```
[style]
layout=A4

[title]
title=Quarterly Business Review
subject=Q4 2024 Performance Report
author=Executive Team

[header]
left={{title}}
right={{date}}
font-size=8
color=#666666

[section 0]

--- BODY ---
<h1>{{title}}</h1>
<p>Prepared by: Finance Department</p>
```

The document will have:
- Title property set to "Quarterly Business Review"
- Subject set to "Q4 2024 Performance Report"
- Author set to "Executive Team"
- Header showing the title and current date
- Body displaying the title as heading

## Notes

- All properties are optional
- Properties are visible in document properties dialog (DOCX) or PDF metadata
- The `{{title}}` variable only works in headers/footers and body content
- Use `{{date}}` for current date, `{{page}}` for page numbers
