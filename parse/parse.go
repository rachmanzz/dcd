package parse

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func Parse(path string) (*Doc, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	defer f.Close()

	doc := &Doc{}
	var cur *Section
	var inBody bool

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		trimmed := strings.TrimSpace(line)

		switch {
		case trimmed == "":
			// skip empty

		case strings.HasPrefix(trimmed, "---"):
			inBody = true

		case trimmed[0] == '[' && trimmed[len(trimmed)-1] == ']':
			cur = &Section{Props: make(map[string]string)}
			cur.Name = trimmed[1 : len(trimmed)-1]
			cur.N = parseSectionN(cur.Name)
			doc.Sections = append(doc.Sections, *cur)
			inBody = false

		default:
			switch {
			case inBody && cur != nil:
				doc.Sections[len(doc.Sections)-1].Body += line + "\n"

			case cur != nil:
				if k, v, ok := strings.Cut(trimmed, "="); ok {
					k = normalizePropertyKey(strings.TrimSpace(k))
					doc.Sections[len(doc.Sections)-1].Props[k] = strings.TrimSpace(v)
				}
			}
		}
	}

	return doc, scan.Err()
}

func parseSectionN(name string) int {
	if strings.HasPrefix(name, "section:next-page ") {
		parts := strings.Fields(name)
		if len(parts) >= 2 {
			if i, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
				return i
			}
		}
	}
	if strings.HasPrefix(name, "section ") {
		parts := strings.Fields(name)
		if len(parts) >= 2 {
			if i, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
				return i
			}
		}
	}
	return 0
}
