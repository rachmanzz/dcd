package render

import "strings"

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