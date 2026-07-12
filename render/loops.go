package render

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	loopRe       = regexp.MustCompile(`(?s)<loop(?::(\w+))?(?:\s+style\.first=(\w+))?\s+(\w+)\s+from\s+([\w.]+)(?:\s+type=(\w+))?(?:\s+style\.first=(\w+))?>(.*?)</loop(?::(\w+))?>`)
	loopSourceRe = regexp.MustCompile(`<loop(?::\w+)?\s+\w+\s+from\s+([\w.]+)(?:\s+type=(\w+))?`)
	objectVarRe  = regexp.MustCompile(`\{\{(\w+)\.`)
)

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