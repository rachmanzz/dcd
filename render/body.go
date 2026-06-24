package render

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	hRe         = regexp.MustCompile(`^<h(\d)(\s+[^>]*)?>(.+)</h(\d)>$`)
	pRe         = regexp.MustCompile(`^<p>(.*)</p>$`)
	wRe         = regexp.MustCompile(`^<w:([^>]+)>(.*)</w>$`)
	imgRe       = regexp.MustCompile(`^<img=(\S+)\s*(.*?)>$`)
	linkRe      = regexp.MustCompile(`<a=(\S+?)(\s+[^>]*)?>([^<]+)</a>`)
	loopRe      = regexp.MustCompile(`(?s)<loop(?::(\w+))?\s+(\w+)\s+from\s+(\w+)>(.*?)</loop(?::\w+)?>`)
	bRe         = regexp.MustCompile(`<b>(.*?)</b>`)
	iRe         = regexp.MustCompile(`<i>(.*?)</i>`)
	uRe         = regexp.MustCompile(`<u>(.*?)</u>`)
	brRe        = regexp.MustCompile(`^<br>$`)
	hrRe        = regexp.MustCompile(`^<hr(\s+[^>]*)?>$`)
	pageBreakRe = regexp.MustCompile(`^<pb>$|^<page-break>$`)
	nestedListRe = regexp.MustCompile(`(?s)<(ul|ol)(?:\s+[^>]*)?>(.*?)</(?:ul|ol)>`)
)

type inlinePart struct {
	tag       string
	text      string
	url       string
	linkAttrs map[string]string
}

func (c *Compiler) renderBody(body string) error {
	lines := strings.Split(body, "\n")
	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		i++
		if line == "" {
			continue
		}

		switch {
		case line == "<ul>" || line == "<ol>":
			ordered := line == "<ol>"
			items, next, err := c.collectListItems(lines, i)
			if err != nil {
				return err
			}
			if err := c.r.AddList(items, ordered); err != nil {
				return err
			}
			i = next

		case strings.HasPrefix(line, "<table"):
			tableEnd := strings.IndexByte(line, '>')
			tableAttrs := parseAttrs(line[6:tableEnd])
			rows, next, err := c.collectTableRows(lines, i)
			if err != nil {
				return err
			}
			if err := c.r.AddTable(rows, tableAttrs); err != nil {
				return err
			}
			i = next

		default:
			if err := c.parseLine(line); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Compiler) expandLoops(body string) string {
	return loopRe.ReplaceAllStringFunc(body, func(match string) string {
		m := loopRe.FindStringSubmatch(match)
		if len(m) < 5 {
			return match
		}
		variant := m[1]
		varName := m[2]
		sourceName := m[3]
		tmpl := m[4]

		raw, ok := c.ds.Get(sourceName)
		if !ok {
			return match
		}

		arr, ok := raw.([]any)
		if !ok {
			return match
		}

		var items []string
		for _, item := range arr {
			expanded := expandLoopTemplate(tmpl, varName, item)
			expanded = c.ds.Resolve(expanded)
			items = append(items, expanded)
		}

		switch variant {
		case "ol":
			var sb strings.Builder
			for _, item := range items {
				sb.WriteString("<li>")
				sb.WriteString(item)
				sb.WriteString("</li>\n")
			}
			return "<ol>\n" + sb.String() + "</ol>"
		case "ul":
			var sb strings.Builder
			for _, item := range items {
				sb.WriteString("<li>")
				sb.WriteString(item)
				sb.WriteString("</li>\n")
			}
			return "<ul>\n" + sb.String() + "</ul>"
		case "row":
			var sb strings.Builder
			for _, item := range items {
				sb.WriteString("<row>")
				sb.WriteString(item)
				sb.WriteString("</row>\n")
			}
			return sb.String()
		default:
			var sb strings.Builder
			for _, item := range items {
				sb.WriteString(item)
				sb.WriteString("\n")
			}
			return sb.String()
		}
	})
}

func expandLoopTemplate(tmpl, varName string, item any) string {
	var result strings.Builder
	pos := 0
	for pos < len(tmpl) {
		start := strings.Index(tmpl[pos:], "{{")
		if start == -1 {
			result.WriteString(tmpl[pos:])
			break
		}
		start += pos
		end := strings.Index(tmpl[start:], "}}")
		if end == -1 {
			result.WriteString(tmpl[pos:])
			break
		}
		end += start + 2
		result.WriteString(tmpl[pos:start])
		expr := tmpl[start+2 : end-2]

		if expr == varName {
			result.WriteString(fmt.Sprintf("%v", item))
		} else if strings.HasPrefix(expr, varName+".") {
			key := expr[len(varName)+1:]
			val, ok := resolveField(item, key)
			if ok {
				result.WriteString(fmt.Sprintf("%v", val))
			} else {
				result.WriteString(tmpl[start:end])
			}
		} else {
			result.WriteString(tmpl[start:end])
		}
		pos = end
	}
	return result.String()
}

func resolveField(item any, key string) (any, bool) {
	m, ok := item.(map[string]any)
	if ok {
		val, found := m[key]
		return val, found
	}
	return nil, false
}

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

	if strings.HasSuffix(startLine, "</li>") {
		return ListItem{Text: startLine[gtIdx+1 : len(startLine)-5]}, start, nil
	}

	textAfter := startLine[gtIdx+1:]

	var buf strings.Builder
	if textAfter != "" {
		buf.WriteString(textAfter)
		buf.WriteString("\n")
	}

	i := start
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		i++
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

	var children []ListItem
	for {
		loc := nestedListRe.FindStringSubmatchIndex(raw)
		if len(loc) < 6 || loc[0] != 0 {
			break
		}
		if loc[0] > 0 {
			break
		}
		nestedType := raw[loc[2]:loc[3]]
		inner := raw[loc[4]:loc[5]]

		subItems, err := parseNestedListItems(inner, nestedType)
		if err != nil {
			return ListItem{}, i, err
		}
		children = append(children, ListItem{Items: subItems})
		raw = raw[loc[1]:]
	}

	text := raw
	if len(children) > 0 {
		for _, loc := range nestedListRe.FindAllStringIndex(raw, -1) {
			text = strings.TrimSpace(raw[:loc[0]])
			break
		}
	}

	item := ListItem{Text: text}
	if len(children) > 0 {
		item.Items = children
	}
	return item, i, nil
}

func parseNestedListItems(raw string, listType string) ([]ListItem, error) {
	var items []ListItem
	// Split on <li>...</li> boundaries
	for {
		liStart := strings.Index(raw, "<li")
		if liStart < 0 {
			break
		}
		gtIdx := strings.IndexByte(raw[liStart:], '>')
		if gtIdx < 0 {
			break
		}
		gtIdx += liStart
		contentStart := gtIdx + 1

		liEnd := strings.Index(raw[contentStart:], "</li>")
		if liEnd < 0 {
			break
		}
		liEnd += contentStart

		itemText := strings.TrimSpace(raw[contentStart:liEnd])

		// Check for nested list inside this item
		item := ListItem{}
		if nlLoc := nestedListRe.FindStringIndex(itemText); nlLoc != nil {
			item.Text = strings.TrimSpace(itemText[:nlLoc[0]])
			innerRaw := itemText[nlLoc[0]:]
			if nlLoc2 := nestedListRe.FindStringSubmatchIndex(innerRaw); len(nlLoc2) >= 6 {
				innerType := innerRaw[nlLoc2[2]:nlLoc2[3]]
				inner := innerRaw[nlLoc2[4]:nlLoc2[5]]
				subItems, _ := parseNestedListItems(inner, innerType)
				item.Items = subItems
			}
		} else {
			item.Text = itemText
		}

		items = append(items, item)
		raw = raw[liEnd+5:]
	}
	return items, nil
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
			if !strings.HasSuffix(line, "</col>") {
				continue
			}
			attrs := parseAttrs(line[4:gtIdx])
			text := line[gtIdx+1 : len(line)-6]
			cells = append(cells, TableCell{Text: text, Attrs: attrs})
		}
	}
	return cells, len(lines), nil
}

func (c *Compiler) parseLine(line string) error {
	switch {
	case brRe.MatchString(line):
		return c.r.AddLineBreak()

	case hrRe.MatchString(line):
		m := hrRe.FindStringSubmatch(line)
		attrs := parseAttrs(m[1])
		return c.r.AddHorizontalRule(attrs)

	case pageBreakRe.MatchString(line):
		return c.r.AddPageBreak()

	case strings.HasPrefix(line, "<h"):
		m := hRe.FindStringSubmatch(line)
		if len(m) == 5 && m[1] == m[4] {
			attrs := parseAttrs(m[2])
			return c.r.AddHeading(m[3], atoi(m[1]), attrs)
		}

	case strings.HasPrefix(line, "<p>"):
		m := pRe.FindStringSubmatch(line)
		if len(m) == 2 {
			return c.renderParagraph(m[1])
		}

	case strings.HasPrefix(line, "<w:"):
		m := wRe.FindStringSubmatch(line)
		if len(m) == 3 {
			return c.r.AddWrappedParagraph(m[2], m[1])
		}

	case strings.HasPrefix(line, "<img="):
		m := imgRe.FindStringSubmatch(line)
		if len(m) == 3 {
			attrs := parseAttrs(m[2])
			return c.r.AddImage(m[1], attrs)
		}

	case strings.HasPrefix(line, "<a=") && strings.HasSuffix(line, "</a>"):
		m := linkRe.FindStringSubmatch(line)
		if len(m) >= 3 {
			attrs := parseAttrs(m[2])
			return c.r.AddHyperlink(m[len(m)-1], m[1], attrs)
		}
	}

	return nil
}

func (c *Compiler) renderParagraph(content string) error {
	parts := splitInline(content)
	runs := make([]TextRun, 0, len(parts))
	for _, part := range parts {
		switch part.tag {
		case "a":
			runs = append(runs, TextRun{Text: part.text, Link: part.url, LinkAttrs: part.linkAttrs})
		default:
			runs = append(runs, TextRun{
				Text:      part.text,
				Bold:      part.tag == "b",
				Italic:    part.tag == "i",
				Underline: part.tag == "u",
			})
		}
	}
	return c.r.AddParagraph(runs)
}

func splitInline(s string) []inlinePart {
	var parts []inlinePart
	pos := 0

	for pos < len(s) {
		type match struct {
			tag       string
			text      string
			url       string
			linkAttrs map[string]string
			skip int
			idx  int
		}
		var best *match

		checkLink := func() {
			loc := linkRe.FindStringSubmatchIndex(s[pos:])
			if len(loc) < 8 {
				return
			}
			idx := loc[0]
			url := s[pos+loc[2] : pos+loc[3]]
			text := s[pos+loc[6] : pos+loc[7]]
			var linkAttrs map[string]string
			if loc[5] > loc[4] {
				linkAttrs = parseAttrs(s[pos+loc[4] : pos+loc[5]])
			}
			if best == nil || idx < best.idx {
				best = &match{tag: "a", text: text, url: url, linkAttrs: linkAttrs, skip: loc[1] - loc[0], idx: idx}
			}
		}

		check := func(tag string, re *regexp.Regexp) {
			loc := re.FindStringSubmatchIndex(s[pos:])
			if len(loc) < 4 {
				return
			}
			idx := loc[0]
			text := s[pos+loc[2] : pos+loc[3]]
			if best == nil || idx < best.idx {
				best = &match{tag: tag, text: text, skip: loc[1] - loc[0], idx: idx}
			}
		}

		checkLink()
		check("b", bRe)
		check("i", iRe)
		check("u", uRe)

		if best == nil {
			parts = append(parts, inlinePart{text: s[pos:]})
			break
		}

		if best.idx > 0 {
			parts = append(parts, inlinePart{text: s[pos : pos+best.idx]})
		}
		parts = append(parts, inlinePart{tag: best.tag, text: best.text, url: best.url, linkAttrs: best.linkAttrs})

		pos += best.idx + best.skip
	}

	return parts
}

func parseAttrs(s string) map[string]string {
	m := make(map[string]string)
	for _, token := range strings.Fields(s) {
		if k, v, ok := strings.Cut(token, "="); ok {
			k = strings.TrimSpace(k)
			v = strings.TrimSpace(v)
			if k != "" {
				m[k] = v
			}
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}
