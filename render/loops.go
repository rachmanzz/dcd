package render

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	loopRe        = regexp.MustCompile(`(?s)<loop(?::(\w+))?(?:\s+style\.first=(\w+))?\s+(\w+)\s+from\s+([\w.]+)(?:\s+type=(\w+))?(?:\s+style\.first=(\w+))?(?:\s+indexType=(\w+))?>(.*?)</loop(?::(\w+))?>`)
	loopSourceRe  = regexp.MustCompile(`<loop(?::\w+)?\s+\w+\s+from\s+([\w.]+)(?:\s+type=(\w+))?`)
	objectVarRe   = regexp.MustCompile(`\{\{(\w+)\.`)
	indexPatternRe  = regexp.MustCompile(`\{index\+(\d+)\}`)
	indexBareRe     = regexp.MustCompile(`\{index\}`)
	loopIndexRe     = regexp.MustCompile(`\{\{loop\.index\}\}`)
)

func (c *Compiler) expandLoops(body string) string {
	return loopRe.ReplaceAllStringFunc(body, func(match string) string {
		m := loopRe.FindStringSubmatch(match)
		if len(m) < 10 {
			return match
		}
		// Validate closing tag variant matches opening variant
		if m[1] != m[9] {
			return match
		}
		variant := m[1]
		styleFirst := m[2] // style.first before varName
		loopType := m[5]   // type= (a/A/i/I)
		if styleFirst == "" {
			styleFirst = m[6] // style.first after sourceName
		}
		indexType := m[7] // indexType= (a/A/i/I/1)
		varName := m[3]
		sourceName := m[4]
		tmpl := m[8]

		raw, ok := c.ds.Get(sourceName)
		if !ok {
			return match
		}

		arr, ok := raw.([]any)
		if !ok {
			return match
		}

		var items []string
		for i, item := range arr {
			expanded := expandLoopTemplate(tmpl, varName, item, i, indexType)
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
			hasLi := strings.Contains(tmpl, "<li")
			for i, item := range items {
				trimmed := strings.TrimSpace(item)
				if hasLi || strings.HasPrefix(trimmed, "<li") {
					sb.WriteString(item)
					sb.WriteString("\n")
				} else {
					if i == 0 && styleFirst != "" {
						sb.WriteString(fmt.Sprintf(`<li style=%s>`, styleFirst))
					} else {
						sb.WriteString("<li>")
					}
					sb.WriteString(item)
					sb.WriteString("</li>\n")
				}
			}
			sb.WriteString("</ol>")
			return sb.String()
		case "ul":
			var sb strings.Builder
			sb.WriteString("<ul>\n")
			hasLi := strings.Contains(tmpl, "<li")
			for i, item := range items {
				trimmed := strings.TrimSpace(item)
				if hasLi || strings.HasPrefix(trimmed, "<li") {
					sb.WriteString(item)
					sb.WriteString("\n")
				} else {
					if i == 0 && styleFirst != "" {
						sb.WriteString(fmt.Sprintf(`<li style=%s>`, styleFirst))
					} else {
						sb.WriteString("<li>")
					}
					sb.WriteString(item)
					sb.WriteString("</li>\n")
				}
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

func expandLoopTemplate(tmpl, varName string, item any, index int, indexType string) string {
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

		if expr == "loop.index" || expr == "index" {
			result.WriteString(formatIndex(index, indexType))
		} else if expr == varName {
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
	s := result.String()
	s = expandIndexPattern(s, index)
	s = expandIndexBare(s, index, indexType)
	return s
}

func expandIndexPattern(s string, index int) string {
	return indexPatternRe.ReplaceAllStringFunc(s, func(match string) string {
		m := indexPatternRe.FindStringSubmatch(match)
		offset, _ := strconv.Atoi(m[1])
		return strconv.Itoa(index + offset)
	})
}

func expandIndexBare(s string, index int, indexType string) string {
	return indexBareRe.ReplaceAllStringFunc(s, func(match string) string {
		return formatIndex(index, indexType)
	})
}

func formatIndex(index int, indexType string) string {
	switch indexType {
	case "a":
		return string(rune('a' + index))
	case "A":
		return string(rune('A' + index))
	case "i":
		return toRomanLower(index + 1)
	case "I":
		return toRomanUpper(index + 1)
	default:
		return strconv.Itoa(index + 1)
	}
}

func toRomanUpper(n int) string {
	vals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	syms := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	var result strings.Builder
	for i := 0; i < len(vals); i++ {
		for n >= vals[i] {
			result.WriteString(syms[i])
			n -= vals[i]
		}
	}
	return result.String()
}

func toRomanLower(n int) string {
	return strings.ToLower(toRomanUpper(n))
}

func resolveField(item any, key string) (any, bool) {
	m, ok := item.(map[string]any)
	if ok {
		val, found := m[key]
		return val, found
	}
	return nil, false
}