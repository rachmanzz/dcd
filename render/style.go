package render

import (
	"regexp"
	"strings"
	"time"

	"github.com/gomutex/godocx/wml/stypes"
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

// normalizePropertyKey maps user-facing property names to internal keys
func normalizePropertyKey(key string) string {
	switch key {
	case "color":
		return "font-color"
	case "bg":
		return "shading"
	default:
		return key
	}
}

var fmtRe = regexp.MustCompile(`\[([^:]+):([^\]]+)\]`)

func parsePageSize(layout, orient, customW, customH string) (width, height float64) {
	switch strings.ToLower(layout) {
	case "letter":
		width, height = 215.9, 279.4
	case "legal":
		width, height = 215.9, 355.6
	case "a3":
		width, height = 297, 420
	case "a5":
		width, height = 148, 210
	case "b5":
		width, height = 176, 250
	case "custom":
		width = atof(customW)
		height = atof(customH)
		if width == 0 {
			width = 210
		}
		if height == 0 {
			height = 297
		}
	default:
		width, height = 210, 297
	}
	if strings.ToLower(orient) == "landscape" {
		width, height = height, width
	}
	return
}

func unitScale(unit string) float64 {
	switch strings.ToLower(unit) {
	case "inch", "in":
		return 25.4
	case "cm":
		return 10
	case "pt":
		return 0.3528
	case "pica":
		return 4.2333
	default:
		return 1
	}
}

func computeMargins(props map[string]string) (left, right, top, bottom float64) {
	scale := unitScale(props["unit"])
	if m := props["m"]; m != "" {
		v := atof(m) * scale
		left, right, top, bottom = v, v, v, v
	}
	if mx := props["mx"]; mx != "" {
		v := atof(mx) * scale
		left, right = v, v
	}
	if my := props["my"]; my != "" {
		v := atof(my) * scale
		top, bottom = v, v
	}
	if md := props["md"]; md != "" {
		v := atof(md) * scale
		left, right, top, bottom = v, v, v, v
	}
	if mt := props["mt"]; mt != "" {
		top = atof(mt) * scale
	}
	if mb := props["mb"]; mb != "" {
		bottom = atof(mb) * scale
	}
	if ml := props["ml"]; ml != "" {
		left = atof(ml) * scale
	}
	if mr := props["mr"]; mr != "" {
		right = atof(mr) * scale
	}
	return
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

func atof(s string) float64 {
	n := 0.0
	dec := false
	div := 1.0
	for _, c := range s {
		if c == '.' {
			dec = true
			continue
		}
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + float64(c-'0')
		if dec {
			div *= 10
		}
	}
	return n / div
}

func chooseAttr(def, style, attrs map[string]string, key string) string {
	if attrs != nil {
		if v, ok := attrs[key]; ok && v != "" {
			return v
		}
	}
	if style != nil {
		if v, ok := style[key]; ok && v != "" {
			return v
		}
	}
	if def != nil {
		if v, ok := def[key]; ok && v != "" {
			return v
		}
	}
	return ""
}

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

func underlineFromString(s string) stypes.Underline {
	switch s {
	case "double":
		return stypes.UnderlineDouble
	case "dotted":
		return stypes.UnderlineDotted
	case "dash":
		return stypes.UnderlineDash
	case "wavy":
		return stypes.UnderlineWavy
	default:
		return stypes.UnderlineSingle
	}
}
