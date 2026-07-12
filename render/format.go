package render

import (
	"regexp"
	"strings"
	"time"
)

var specConv = map[string]string{
	"dd":   "02",
	"MM":   "01",
	"yyyy": "2006",
	"HH":   "15",
	"mm":   "04",
	"ss":   "05",
}

// specRe matches any custom format specifier at word boundaries.
var specRe = func() *regexp.Regexp {
	keys := []string{"dd", "MM", "yyyy", "HH", "mm", "ss"}
	pat := `\b(?:` + strings.Join(keys, `|`) + `)\b`
	return regexp.MustCompile(pat)
}()

var fmtRe = regexp.MustCompile(`\[([^:]+):([^\]]+)\]`)

// parseFormats parses "formats" property value like:
//
//	[date_field:dd-MM-yyyy], [time_field:HH\:m]
//
// into map[key]format.
func parseFormats(s string) map[string]string {
	m := make(map[string]string)
	for _, match := range fmtRe.FindAllStringSubmatch(s, -1) {
		if len(match) == 3 {
			key := strings.TrimSpace(match[1])
			fmtStr := strings.TrimSpace(match[2])
			if key != "" && fmtStr != "" {
				m[key] = fmtStr
			}
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

// convertFormat translates custom specifiers to Go time format.
//
//	dd → 02   (day)
//	MM → 01   (month)
//	yyyy → 2006 (year)
//	HH → 15   (hour)
//	mm → 04   (minute)
//	ss → 05   (second)
//
// Escaped \: is unescaped → literal : (to avoid breaking the [key:format] parser).
// Non-matching text is passed through as-is (supports regex patterns like \d, \w).
func convertFormat(fmtStr string) string {
	result := specRe.ReplaceAllStringFunc(fmtStr, func(match string) string {
		if goFmt, ok := specConv[match]; ok {
			return goFmt
		}
		return match
	})
	result = strings.ReplaceAll(result, "\\:", ":")
	return result
}

// applyFormat applies a time format string to a value.
// Custom specifiers (dd, MM, yyyy, etc.) are converted to Go format automatically.
func applyFormat(val string, fmtStr string) string {
	fmtStr = convertFormat(fmtStr)
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"02-01-2006",
		"01/02/2006",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, val); err == nil {
			return t.Format(fmtStr)
		}
	}
	return val
}