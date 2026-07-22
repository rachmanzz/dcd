package render

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rachmanzz/dcd/internal/property"
)

var (
	hRe         = regexp.MustCompile(`^<h(\d)(\s+[^>]*)?>(.+)</h(\d)>$`)
	pRe         = regexp.MustCompile(`^<p(\s+[^>]*)?>(.*)</p>$`)
	wRe         = regexp.MustCompile(`^<w:([^\s>]+)(\s+[^>]*)?>(.*)</w:([^\s>]+)>$`)
	imgRe       = regexp.MustCompile(`^<img=(\S+)\s*(.*?)>$`)
	linkRe      = regexp.MustCompile(`<a=(\S+?)(\s+[^>]*)?>([^<]+)</a>`)
	brRe        = regexp.MustCompile(`^<br/?>$`)
	hrRe        = regexp.MustCompile(`^<hr(\s+[^>]*)?/?>$`)
	pageBreakRe = regexp.MustCompile(`^<pb/?>$|^<page-break/?>$`)
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
	// <tab>, <tab/>, <tab size=N>, and <br> are allowed inside <w:*> tags
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
			// <tab>, <tab/>, <tab size=N>, and <br> are allowed inside <w:*> tags
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
				m[property.NormalizeKey(k)] = v
			}
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}