package render

import (
	"fmt"
	"strings"

	"github.com/rachmanzz/dcd/parse"
)

type varEntry struct {
	Name    string
	IsArray bool
}

func parseVarDecl(raw string) []varEntry {
	var entries []varEntry
	for _, part := range strings.Split(raw, ",") {
		name := strings.TrimSpace(part)
		if name == "" {
			continue
		}
		isArray := strings.HasPrefix(name, "[]")
		if isArray {
			name = strings.TrimPrefix(name, "[]")
		}
		entries = append(entries, varEntry{Name: name, IsArray: isArray})
	}
	return entries
}

// validateSectionProps checks section properties against skill rules.
// Only [section N] / [section:next-page N] sections are validated.
func (c *Compiler) validateSectionProps(sec parse.Section) error {
	name := sec.Name
	if !strings.HasPrefix(name, "section") {
		return nil
	}

	// Check 10: name= is REQUIRED
	if sec.Props["name"] == "" {
		return fmt.Errorf("name= is required in section %q", name)
	}

	hasVar := sec.Props["var"] != ""
	hasKeys := false
	for k := range sec.Props {
		if k == "keys" {
			hasKeys = true
			break
		}
	}

	if !hasVar && !hasKeys {
		return nil
	}

	if hasVar {
		vars := parseVarDecl(sec.Props["var"])

		// Check 5: Section limits — warning only
		if len(vars) > 5 {
			fmt.Printf("warning: section %q has %d var entries (max 5)\n", sec.Name, len(vars))
		}

		// Collect array names and object names from var=
		arrayNames := make(map[string]bool)
		objectNames := make(map[string]bool)
		for _, v := range vars {
			if v.IsArray {
				arrayNames[v.Name] = true
			} else {
				objectNames[v.Name] = true
			}
		}

		// Check 1: Array ([]) used as object prefix {{name.key}}
		for _, m := range objectVarRe.FindAllStringSubmatch(sec.Body, -1) {
			if len(m) >= 2 {
				prefix := m[1]
				if arrayNames[prefix] {
					return fmt.Errorf("array var %q (declared with []) used as object prefix {{%s.key}}", prefix, prefix)
				}
			}
		}

		// Check 2 & 3: Loop sources — must be declared with [] in var=
		for _, m := range loopSourceRe.FindAllStringSubmatch(sec.Body, -1) {
			if len(m) < 2 {
				continue
			}
			source := m[1]
			// Skip dotted paths (e.g. invoice.items) — resolved as data path
			if strings.Contains(source, ".") {
				continue
			}
			if objectNames[source] {
				return fmt.Errorf("object var %q (without []) used as loop source", source)
			}
			if !arrayNames[source] && !objectNames[source] {
				return fmt.Errorf("loop source %q not declared in var=", source)
			}
		}

		// Check 4: [] var declared but never used as loop source
		for _, v := range vars {
			if v.IsArray {
				matched := false
				for _, m := range loopSourceRe.FindAllStringSubmatch(sec.Body, -1) {
					if len(m) >= 2 && m[1] == v.Name {
						matched = true
						break
					}
				}
				if !matched {
					return fmt.Errorf("array var %q declared with [] but never used as loop source", v.Name)
				}
			}
		}

		// Check 6: Strict Usage — object vars must appear as {{name.key}}
		for _, v := range vars {
			if !v.IsArray {
				matched := false
				for _, m := range objectVarRe.FindAllStringSubmatch(sec.Body, -1) {
					if len(m) >= 2 && m[1] == v.Name {
						matched = true
						break
					}
				}
				if !matched {
					return fmt.Errorf("object var %q declared but never used as {{%s.key}} in body", v.Name, v.Name)
				}
			}
		}
	}

	// Check 5: Section limits — keys warning
	keyList := strings.Split(sec.Props["keys"], ",")
	if len(keyList) > 15 {
		fmt.Printf("warning: section %q has %d keys entries (max 15)\n", sec.Name, len(keyList))
	}

	// Check 5a (CONDITIONAL DOT-NOTATION): dotted paths in keys= must have format
	for _, k := range keyList {
		k = strings.TrimSpace(k)
		if k != "" && strings.Contains(k, ".") {
			fmts := sec.Props["formats"]
			if fmts == "" || !strings.Contains(fmts, k) {
				return fmt.Errorf("dotted key %q in keys= must have corresponding format in formats=", k)
			}
		}
	}

	if fmts := sec.Props["formats"]; fmts != "" {
		fmtMap := parseFormats(fmts)
		keySet := make(map[string]bool)
		for _, k := range keyList {
			keySet[strings.TrimSpace(k)] = true
		}
		for key := range fmtMap {
			leaf := key
			if lastDot := strings.LastIndex(key, "."); lastDot >= 0 {
				leaf = key[lastDot+1:]
			}
			if !keySet[leaf] && !keySet[key] {
				return fmt.Errorf("format key %q not found in keys", key)
			}
		}
	}

	return nil
}