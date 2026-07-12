package render

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Inline markup regexes
	bRe       = regexp.MustCompile(`<b>(.*?)</b>`)
	iRe       = regexp.MustCompile(`<i>(.*?)</i>`)
	uRe       = regexp.MustCompile(`<u(?:=(\w+))?>([^<]+)</u>`)
	sRe       = regexp.MustCompile(`<s>(.*?)</s>`)
	codeRe    = regexp.MustCompile(`<code>(.*?)</code>`)
	markRe    = regexp.MustCompile(`<mark(?:\s+color=(\S+))?>(.*?)</mark>`)
	subRe     = regexp.MustCompile(`<sub>(.*?)</sub>`)
	supRe     = regexp.MustCompile(`<sup>(.*?)</sup>`)
	setRe     = regexp.MustCompile(`<set:([^\s>]+)(\s+[^>]+)?>(.*?)</set(?::([^>]+))?>`)
	tabRe     = regexp.MustCompile(`^<tab(\s+size=(\d+))?\s*/?>$`)
	openTagRe = regexp.MustCompile(`<(b|i|u|s|code|mark|sub|sup)(?:[=>\s][^><]*)?>`)
	openARe   = regexp.MustCompile(`<a=[^><]*>`)
	openSetRe = regexp.MustCompile(`<set:(\w+(?:\|\w+)*)(?:\s+[^><]*)?>`)
	closeTagRe  = regexp.MustCompile(`</(b|i|u|s|code|mark|sub|sup|a)>`)
	closeSetRe  = regexp.MustCompile(`</set(?::(\w+(?:\|\w+)*))?>`)
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