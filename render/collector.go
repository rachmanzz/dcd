package render

import (
	"regexp"
	"strings"
)

var pTagOpenRe = regexp.MustCompile(`(?s)<p(?:\s[^>]*)?>`)

func (c *Compiler) collectListItems(lines []string, start int) ([]ListItem, int, error) {
	var items []ListItem
	i := start
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		i++
		if line == "</ul>" || line == "</ol>" {
			break
		}
		if strings.HasPrefix(line, "<li") {
			item, next, err := c.collectLi(lines, i, line)
			if err != nil {
				return nil, 0, err
			}
			items = append(items, item)
			i = next
		}
	}
	return items, i, nil
}

func (c *Compiler) collectLi(lines []string, start int, startLine string) (ListItem, int, error) {
	gtIdx := strings.IndexByte(startLine, '>')
	if gtIdx < 0 {
		return ListItem{}, start, nil
	}

	attrs := parseAttrs(startLine[3:gtIdx])

	if strings.HasSuffix(startLine, "</li>") {
		text := startLine[gtIdx+1 : len(startLine)-5]
		runs, err := inlineToRuns(text)
		if err != nil {
			return ListItem{}, 0, err
		}
		return ListItem{Runs: runs, Attrs: attrs}, start, nil
	}

	textAfter := startLine[gtIdx+1:]

	var buf strings.Builder
	if textAfter != "" {
		buf.WriteString(textAfter)
		buf.WriteString("\n")
	}

	var subItems []ListItem
	var subOrdered bool
	var subNumFmt string
	i := start
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		i++

		if line == "<ul>" || strings.HasPrefix(line, "<ol") {
			subOrdered = strings.HasPrefix(line, "<ol")
			subNumFmt = ""
			if subOrdered && line != "<ol>" {
				subNumFmt = parseListType(line)
			}
			var err error
			subItems, i, err = c.collectListItems(lines, i)
			if err != nil {
				return ListItem{}, 0, err
			}
			continue
		}

		if strings.HasSuffix(line, "</li>") {
			if line != "</li>" {
				buf.WriteString(line[:len(line)-5])
			}
			break
		}
		buf.WriteString(line)
		buf.WriteString("\n")
	}

	raw := strings.TrimSpace(buf.String())

	raw = pTagOpenRe.ReplaceAllString(raw, "")
	raw = strings.ReplaceAll(raw, "</p>", "")

	runs, err := inlineToRuns(raw)
	if err != nil {
		return ListItem{}, 0, err
	}
	item := ListItem{Runs: runs, Attrs: attrs}
	if len(subItems) > 0 {
		item.Items = subItems
		item.Ordered = subOrdered
		item.NumFormat = subNumFmt
	}
	return item, i, nil
}

func (c *Compiler) collectTableRows(lines []string, start int) ([]TableRow, int, error) {
	var rows []TableRow
	i := start
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		i++
		if line == "</table>" {
			break
		}
		if strings.HasPrefix(line, "<row") {
			end := strings.IndexByte(line, '>')
			props := parseAttrs(line[4:end])
			cells, next, err := c.collectRowCells(lines, i)
			if err != nil {
				return nil, 0, err
			}
			rows = append(rows, TableRow{Cells: cells, Props: props})
			i = next
		}
	}
	return rows, i, nil
}

func (c *Compiler) collectRowCells(lines []string, start int) ([]TableCell, int, error) {
	var cells []TableCell
	for i := start; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "</row>" {
			return cells, i + 1, nil
		}
		if strings.HasPrefix(line, "<col") {
			gtIdx := strings.Index(line, ">")
			if gtIdx < 0 {
				continue
			}
			if strings.HasSuffix(line, "</col>") {
				attrs := parseAttrs(line[4:gtIdx])
				text := line[gtIdx+1 : len(line)-6]
				runs, err := inlineToRuns(text)
				if err != nil {
					return nil, 0, err
				}
				cells = append(cells, TableCell{Runs: runs, Attrs: attrs})
			} else {
				// Multi-line <col> — collect until </col>
				firstAttrs := parseAttrs(line[4:gtIdx])
				var buf strings.Builder
				after := line[gtIdx+1:]
				if after != "" {
					buf.WriteString(after)
				}
				i++
				for i < len(lines) {
					tr := strings.TrimSpace(lines[i])
					i++
					if strings.HasSuffix(tr, "</col>") {
						idx := strings.LastIndex(tr, "</col>")
						if idx > 0 {
							if buf.Len() > 0 {
								buf.WriteString("\n")
							}
							buf.WriteString(tr[:idx])
						}
						break
					}
					if buf.Len() > 0 {
						buf.WriteString("\n")
					}
					buf.WriteString(tr)
				}
				runs2, err := inlineToRuns(strings.TrimSpace(buf.String()))
				if err != nil {
					return nil, 0, err
				}
				cells = append(cells, TableCell{Runs: runs2, Attrs: firstAttrs})
			}
		}
	}
	return cells, len(lines), nil
}