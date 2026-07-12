package render

import "github.com/gomutex/godocx/wml/stypes"

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