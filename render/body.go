package render

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	hRe          = regexp.MustCompile(`^<h(\d)(\s+[^>]*)?>(.+)</h(\d)>$`)
	pRe          = regexp.MustCompile(`^<p(\s+[^>]*)?>(.*)</p>$`)
	wRe          = regexp.MustCompile(`^<w:([^\s>]+)(\s+[^>]*)?>(.*)</w:([^\s>]+)>$`)
	imgRe        = regexp.MustCompile(`^<img=(\S+)\s*(.*?)>$`)
	linkRe       = regexp.MustCompile(`<a=(\S+?)(\s+[^>]*)?>([^<]+)</a>`)
	loopRe       = regexp.MustCompile(`(?s)<loop(?::(\w+))?(?:\s+style\.first=(\w+))?\s+(\w+)\s+from\s+([\w.]+)(?:\s+type=(\w+))?(?:\s+style\.first=(\w+))?>(.*?)</loop(?::(\w+))?>`)
	loopSourceRe = regexp.MustCompile(`<loop(?::\w+)?\s+\w+\s+from\s+([\w.]+)(?:\s+type=(\w+))?`)
	objectVarRe  = regexp.MustCompile(`\{\{(\w+)\.`)
	bRe          = regexp.MustCompile(`<b>(.*?)</b>`)
	iRe          = regexp.MustCompile(`<i>(.*?)</i>`)
	uRe          = regexp.MustCompile(`<u(?:=(\w+))?>([^<]+)</u>`)
	sRe          = regexp.MustCompile(`<s>(.*?)</s>`)
	codeRe       = regexp.MustCompile(`<code>(.*?)</code>`)
	markRe       = regexp.MustCompile(`<mark(?:\s+color=(\S+))?>(.*?)</mark>`)
	subRe        = regexp.MustCompile(`<sub>(.*?)</sub>`)
	supRe        = regexp.MustCompile(`<sup>(.*?)</sup>`)
	setRe        = regexp.MustCompile(`<set:([^\s>]+)(\s+[^>]+)?>(.*?)</set(?::([^>]+))?>`)
	brRe         = regexp.MustCompile(`^<br/?>$`)
	tabRe        = regexp.MustCompile(`^<tab(\s+size=(\d+))?\s*/?>$`)
	hrRe         = regexp.MustCompile(`^<hr(\s+[^>]*)?/?>$`)
	pageBreakRe  = regexp.MustCompile(`^<pb/?>$|^<page-break/?>$`)
)

type inlinePart struct {
	tag            string
	text           string
	url            string
	linkAttrs      map[string]string
	markColor      string
	underlineStyle string
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
		case strings.HasPrefix(line, "<ul") || strings.HasPrefix(line, "<ol"):
			ordered := strings.HasPrefix(line, "<ol")
			var numFmt string
			if ordered && !strings.HasPrefix(line, "<ol>") {
				numFmt = parseListType(line)
			}
			items, next, err := c.collectListItems(lines, i)
			if err != nil {
				return err
			}
			if err := c.r.AddList(items, ordered, numFmt); err != nil {
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

		case (strings.HasPrefix(line, "<p>") || strings.HasPrefix(line, "<p ")) && !strings.HasSuffix(line, "</p>"):
			closeIdx := strings.LastIndex(line, "</p>")
			if closeIdx < 0 {
				content, attrs, next, err := c.collectPTag(lines, i, line)
				if err != nil {
					return err
				}
				i = next
				if err := c.renderParagraph(content, attrs); err != nil {
					return err
				}
				break
			}
			if err := c.parseLine(line); err != nil {
				return err
			}

		case strings.HasPrefix(line, "<w:") && !strings.Contains(line, "</w:"):
			flags, attrs, content, next, err := c.collectWTag(lines, i, line)
			if err != nil {
				return err
			}
			i = next
			if err := c.renderWrappedContent(content, flags, attrs); err != nil {
				return err
			}

		default:
			if err := c.parseLine(line); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Compiler) collectPTag(lines []string, start int, firstLine string) (string, map[string]string, int, error) {
	gtIdx := strings.IndexByte(firstLine, '>')
	if gtIdx < 0 {
		return "", nil, start, nil
	}
	attrs := parseAttrs(tagAttrs(firstLine, gtIdx))
	var buf strings.Builder
	after := firstLine[gtIdx+1:]
	if after != "" {
		buf.WriteString(after)
	}
	i := start
	for i < len(lines) {
		tr := strings.TrimSpace(lines[i])
		i++
		if strings.HasSuffix(tr, "</p>") {
			idx := strings.LastIndex(tr, "</p>")
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
	content := strings.TrimSpace(buf.String())
	if strings.Contains(content, "<h") {
		return "", nil, 0, fmt.Errorf("<p>: heading tags inside <p> are not allowed")
	}
	return content, attrs, i, nil
}

func tagAttrs(tagLine string, gtIdx int) string {
	spIdx := strings.IndexByte(tagLine, ' ')
	if spIdx < 0 || spIdx > gtIdx {
		return ""
	}
	return tagLine[spIdx:gtIdx]
}

func (c *Compiler) collectWTag(lines []string, start int, firstLine string) (string, map[string]string, string, int, error) {
	gtIdx := strings.IndexByte(firstLine, '>')
	if gtIdx < 0 {
		return "", nil, "", start, nil
	}
	flags := firstLine[3:gtIdx]
	if spIdx := strings.IndexByte(flags, ' '); spIdx >= 0 {
		flags = flags[:spIdx]
	}
	attrs := parseAttrs(tagAttrs(firstLine, gtIdx))
	var buf strings.Builder
	after := firstLine[gtIdx+1:]
	if after != "" {
		buf.WriteString(after)
	}
	i := start
	closeTag := "</w:" + flags + ">"
	for i < len(lines) {
		tr := strings.TrimSpace(lines[i])
		i++
		if strings.Contains(tr, closeTag) {
			idx := strings.Index(tr, closeTag)
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
	content := strings.TrimSpace(buf.String())
	if strings.Contains(content, "<w:") {
		return "", nil, "", 0, fmt.Errorf("<w:%s>: nested <w:> tags are not allowed", flags)
	}
	if strings.Contains(content, "<h") {
		return "", nil, "", 0, fmt.Errorf("<w:%s>: heading tags inside <w:> are not allowed", flags)
	}
	return flags, attrs, content, i, nil
}

func (c *Compiler) renderWrappedContent(content, flags string, attrs map[string]string) error {
	runs, err := inlineToRuns(content)
	if err != nil {
		return err
	}
	for _, f := range strings.Split(flags, "|") {
		switch f {
		case "c":
			attrs = mergeAttrs(attrs, map[string]string{"align": "center"})
		case "r":
			attrs = mergeAttrs(attrs, map[string]string{"align": "right"})
		case "j":
			attrs = mergeAttrs(attrs, map[string]string{"align": "justify"})
		case "b":
			for i := range runs {
				if !runs[i].Bold {
					runs[i].Bold = true
				}
			}
		case "i":
			for i := range runs {
				if !runs[i].Italic {
					runs[i].Italic = true
				}
			}
		case "u":
			for i := range runs {
				if !runs[i].Underline {
					runs[i].Underline = true
				}
				if runs[i].UnderlineStyle == "" && attrs["underline"] != "" {
					runs[i].UnderlineStyle = attrs["underline"]
				}
			}
		case "s":
			for i := range runs {
				if !runs[i].Strike {
					runs[i].Strike = true
				}
			}
		case "code":
			for i := range runs {
				if !runs[i].Code {
					runs[i].Code = true
				}
			}
		}
	}
	return c.r.AddParagraph(runs, attrs)
}

// For renderWrapped, we need to parse the <w:> tag from the collected content
// Let me redo this differently

func (c *Compiler) expandLoops(body string) string {
	return loopRe.ReplaceAllStringFunc(body, func(match string) string {
		m := loopRe.FindStringSubmatch(match)
		if len(m) < 9 {
			return match
		}
		// Validate closing tag variant matches opening variant
		if m[1] != m[8] {
			return match
		}
		variant := m[1]
		styleFirst := m[2] // style.first before varName
		loopType := m[5]   // type= (a/A/i/I)
		if styleFirst == "" {
			styleFirst = m[6] // style.first after sourceName
		}
		varName := m[3]
		sourceName := m[4]
		tmpl := m[7]

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
			if loopType != "" {
				sb.WriteString(fmt.Sprintf("<ol type=%s>\n", loopType))
			} else {
				sb.WriteString("<ol>\n")
			}
			for i, item := range items {
				if i == 0 && styleFirst != "" {
					sb.WriteString(fmt.Sprintf(`<li style=%s>`, styleFirst))
				} else {
					sb.WriteString("<li>")
				}
				sb.WriteString(item)
				sb.WriteString("</li>\n")
			}
			sb.WriteString("</ol>")
			return sb.String()
		case "ul":
			var sb strings.Builder
			sb.WriteString("<ul>\n")
			for i, item := range items {
				if i == 0 && styleFirst != "" {
					sb.WriteString(fmt.Sprintf(`<li style=%s>`, styleFirst))
				} else {
					sb.WriteString("<li>")
				}
				sb.WriteString(item)
				sb.WriteString("</li>\n")
			}
			sb.WriteString("</ul>")
			return sb.String()
		case "row":
			var sb strings.Builder
			for i, item := range items {
				if i == 0 && styleFirst != "" {
					sb.WriteString(fmt.Sprintf(`<row style=%s>`, styleFirst))
				} else {
					sb.WriteString("<row>")
				}
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
			if strings.Contains(m[3], "<h") {
				return fmt.Errorf("heading nesting is not allowed")
			}
			return c.r.AddHeading(m[3], atoi(m[1]), attrs)
		}

	case strings.HasPrefix(line, "<p>") || strings.HasPrefix(line, "<p "):
		m := pRe.FindStringSubmatch(line)
		if len(m) == 3 {
			if strings.Contains(m[2], "<h") {
				return fmt.Errorf("heading tags inside <p> are not allowed")
			}
			return c.renderParagraph(m[2], parseAttrs(m[1]))
		}

	case strings.HasPrefix(line, "<w:"):
		m := wRe.FindStringSubmatch(line)
		if len(m) == 5 && m[1] == m[4] {
			if strings.Contains(m[3], "<w:") {
				return fmt.Errorf("<w:%s>: nested <w:> tags are not allowed", m[1])
			}
			if strings.Contains(m[3], "<h") {
				return fmt.Errorf("<w:%s>: heading tags inside <w:> are not allowed", m[1])
			}
			return c.renderWrappedContent(m[3], m[1], parseAttrs(m[2]))
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

func (c *Compiler) renderParagraph(content string, attrs map[string]string) error {
	runs, err := inlineToRuns(content)
	if err != nil {
		return err
	}
	return c.r.AddParagraph(runs, attrs)
}

var (
	openTagRe  = regexp.MustCompile(`<(b|i|u|s|code|mark|sub|sup)(?:[=>\s][^><]*)?>`)
	openARe    = regexp.MustCompile(`<a=[^><]*>`)
	openSetRe  = regexp.MustCompile(`<set:(\w+(?:\|\w+)*)(?:\s+[^><]*)?>`)
	closeTagRe = regexp.MustCompile(`</(b|i|u|s|code|mark|sub|sup|a)>`)
	closeSetRe = regexp.MustCompile(`</set(?::(\w+(?:\|\w+)*))?>`)
)

type tagMatch struct {
	m    string
	name string
	idx  int
	end  int
	kind byte // 'o' = opening, 'c' = closing
}

func validateTagBalance(s string) error {
	var stack []string
	pos := 0
	for pos < len(s) {
		next := strings.IndexByte(s[pos:], '<')
		if next < 0 {
			break
		}
		pos += next

		var best *tagMatch
		check := func(re *regexp.Regexp, kind byte, nameFn func(string) string) {
			loc := re.FindStringIndex(s[pos:])
			if len(loc) < 2 {
				return
			}
			m := s[pos+loc[0] : pos+loc[1]]
			name := nameFn(m)
			if best == nil || loc[0] < best.idx {
				best = &tagMatch{m: m, name: name, idx: loc[0], end: loc[1], kind: kind}
			}
		}

		check(closeTagRe, 'c', func(m string) string { return closeTagRe.FindStringSubmatch(m)[1] })
		check(closeSetRe, 'c', func(m string) string {
			sm := closeSetRe.FindStringSubmatch(m)
			if sm[1] == "" {
				return "set:"
			}
			return "set:" + sm[1]
		})
		check(openTagRe, 'o', func(m string) string { return openTagRe.FindStringSubmatch(m)[1] })
		check(openARe, 'o', func(m string) string { return "a" })
		check(openSetRe, 'o', func(m string) string { return "set:" + openSetRe.FindStringSubmatch(m)[1] })

		if best == nil {
			// Unknown tag — skip past it
			if end := strings.IndexByte(s[pos+1:], '>'); end >= 0 {
				pos += end + 2
			} else {
				pos++
			}
			continue
		}

		switch best.kind {
		case 'c':
			if len(stack) == 0 {
				return fmt.Errorf("unexpected closing tag <%s>", best.m)
			}
			top := stack[len(stack)-1]
			if top != best.name {
				if !(best.name == "set:" && strings.HasPrefix(top, "set:")) {
					return fmt.Errorf("tag balancing error: expected </%s> but found <%s>", top, best.m)
				}
			}
			stack = stack[:len(stack)-1]
		case 'o':
			stack = append(stack, best.name)
		}
		pos += best.end
	}
	if len(stack) > 0 {
		return fmt.Errorf("tag balancing error: unclosed tag <%s>", stack[len(stack)-1])
	}
	return nil
}

func inlineToRuns(content string) ([]TextRun, error) {
	if err := validateTagBalance(content); err != nil {
		return nil, err
	}
	parts := splitInline(content)
	runs := make([]TextRun, 0, len(parts))
	for _, part := range parts {
		switch part.tag {
		case "a":
			runs = append(runs, TextRun{Text: part.text, Link: part.url, LinkAttrs: part.linkAttrs})
		case "tab":
			n := 1
			if part.text != "" {
				if v, err := strconv.Atoi(part.text); err == nil && v > 0 {
					n = v
				}
			}
			for i := 0; i < n; i++ {
				runs = append(runs, TextRun{Tab: true})
			}
		default:
			// Check if part.tag contains "|" or has attrs (set:flags format)
			if strings.Contains(part.tag, "|") || strings.Contains(part.tag, " ") {
				tokens := strings.Fields(part.tag)
				flags := strings.Split(tokens[0], "|")
				run := TextRun{Text: part.text}
				for _, flag := range flags {
					switch strings.TrimSpace(flag) {
					case "b":
						run.Bold = true
					case "i":
						run.Italic = true
					case "u":
						run.Underline = true
					case "s":
						run.Strike = true
					case "code":
						run.Code = true
					}
				}
				if len(tokens) > 1 {
					setAttrs := parseAttrs(strings.Join(tokens[1:], " "))
					if u := setAttrs["underline"]; u != "" {
						run.UnderlineStyle = u
					}
				}
				runs = append(runs, run)
			} else {
				// Single tag (existing behavior)
				runs = append(runs, TextRun{
					Text:           part.text,
					Bold:           part.tag == "b",
					Italic:         part.tag == "i",
					Underline:      part.tag == "u",
					UnderlineStyle: part.underlineStyle,
					Strike:         part.tag == "s",
					Code:           part.tag == "code",
					Mark:           part.tag == "mark",
					MarkColor:      part.markColor,
					Sub:            part.tag == "sub",
					Sup:            part.tag == "sup",
				})
			}
		}
	}
	return runs, nil
}

func splitInline(s string) []inlinePart {
	var parts []inlinePart
	pos := 0

	for pos < len(s) {
		type match struct {
			tag            string
			text           string
			url            string
			linkAttrs      map[string]string
			markColor      string
			underlineStyle string
			skip           int
			idx            int
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
		check("s", sRe)
		check("code", codeRe)
		check("sub", subRe)
		check("sup", supRe)

		checkUnderline := func() {
			loc := uRe.FindStringSubmatchIndex(s[pos:])
			if len(loc) < 6 {
				return
			}
			idx := loc[0]
			uStyle := ""
			if loc[2] >= 0 {
				uStyle = s[pos+loc[2] : pos+loc[3]]
			}
			text := s[pos+loc[4] : pos+loc[5]]
			if best == nil || idx < best.idx {
				best = &match{tag: "u", text: text, underlineStyle: uStyle, skip: loc[1] - loc[0], idx: idx}
			}
		}
		checkUnderline()

		checkMark := func() {
			loc := markRe.FindStringSubmatchIndex(s[pos:])
			if len(loc) < 6 {
				return
			}
			idx := loc[0]
			color := ""
			if loc[2] >= 0 {
				color = s[pos+loc[2] : pos+loc[3]]
			}
			text := s[pos+loc[4] : pos+loc[5]]
			if best == nil || idx < best.idx {
				best = &match{tag: "mark", text: text, markColor: color, skip: loc[1] - loc[0], idx: idx}
			}
		}
		checkMark()

		checkTab := func() {
			loc := tabRe.FindStringSubmatchIndex(s[pos:])
			if len(loc) < 6 || loc[0] < 0 {
				return
			}
			idx := loc[0]
			size := "1"
			if loc[4] >= 0 {
				size = s[pos+loc[4] : pos+loc[5]]
			}
			if best == nil || idx < best.idx {
				best = &match{tag: "tab", text: size, skip: loc[1] - loc[0], idx: idx}
			}
		}
		checkTab()

		// Check for <set:flags> tag
		checkSet := func() {
			loc := setRe.FindStringSubmatchIndex(s[pos:])
			if len(loc) < 10 {
				return
			}
			idx := loc[0]
			openFlags := s[pos+loc[2] : pos+loc[3]]
			closeFlags := ""
			if loc[8] >= 0 {
				closeFlags = s[pos+loc[8] : pos+loc[9]]
			}
			if openFlags != closeFlags && closeFlags != "" {
				return
			}
			attrs := ""
			if loc[4] >= 0 {
				attrs = s[pos+loc[4] : pos+loc[5]]
			}
			text := s[pos+loc[6] : pos+loc[7]]
			if attrs != "" {
				openFlags = openFlags + " " + attrs
			}
			if best == nil || idx < best.idx {
				best = &match{tag: openFlags, text: text, skip: loc[1] - loc[0], idx: idx}
			}
		}
		checkSet()

		if best == nil {
			parts = append(parts, inlinePart{text: s[pos:]})
			break
		}

		if best.idx > 0 {
			parts = append(parts, inlinePart{text: s[pos : pos+best.idx]})
		}
		parts = append(parts, inlinePart{tag: best.tag, text: best.text, url: best.url, linkAttrs: best.linkAttrs, markColor: best.markColor, underlineStyle: best.underlineStyle})

		pos += best.idx + best.skip
	}

	return parts
}

func mergeAttrs(base, over map[string]string) map[string]string {
	if len(over) == 0 {
		return base
	}
	if base == nil {
		base = make(map[string]string)
	}
	for k, v := range over {
		if _, ok := base[k]; !ok {
			base[k] = v
		}
	}
	return base
}

func parseListType(line string) string {
	idx := strings.Index(line, "type=")
	if idx < 0 {
		return ""
	}
	rest := line[idx+5:]
	end := strings.IndexAny(rest, " >")
	if end < 0 {
		return ""
	}
	return rest[:end]
}

func parseAttrs(s string) map[string]string {
	m := make(map[string]string)
	for _, token := range strings.Fields(s) {
		if k, v, ok := strings.Cut(token, "="); ok {
			k = strings.TrimSpace(k)
			v = strings.TrimSpace(v)
			if k != "" {
				m[normalizePropertyKey(k)] = v
			}
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}
