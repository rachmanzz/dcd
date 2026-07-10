package render

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rachmanzz/dcd/data"
	"github.com/rachmanzz/dcd/parse"
)

var varRe = regexp.MustCompile(`\{\{[^}]+\}\}`)

type Compiler struct {
	doc *parse.Doc
	ds  *data.DataSet
	r   Renderer
}

func New(doc *parse.Doc, ds *data.DataSet, r Renderer) *Compiler {
	return &Compiler{doc: doc, ds: ds, r: r}
}

func (c *Compiler) Run(output string) error {
	if err := c.r.SetPageStyle(c.collectSection("style")); err != nil {
		return fmt.Errorf("render page style: %w", err)
	}
	if err := c.r.SetDefaultStyle(c.collectSection("style")); err != nil {
		return fmt.Errorf("set default style: %w", err)
	}
	if err := c.applyHeadingStyles(); err != nil {
		return fmt.Errorf("apply heading styles: %w", err)
	}
	if err := c.applyTableStyles(); err != nil {
		return fmt.Errorf("apply table styles: %w", err)
	}
	if err := c.applyMetadata(); err != nil {
		return err
	}
	if err := c.applyHeaderFooter(); err != nil {
		return err
	}

	for _, sec := range c.doc.Sections {
		if sec.Body == "" {
			continue
		}
		if err := c.renderSection(sec); err != nil {
			return fmt.Errorf("render section %q: %w", sec.Name, err)
		}
	}

	return c.r.Save(output)
}

func (c *Compiler) collectSection(name string) map[string]string {
	for _, sec := range c.doc.Sections {
		if sec.Name == name {
			return sec.Props
		}
	}
	return nil
}

func (c *Compiler) applyHeadingStyles() error {
	for _, sec := range c.doc.Sections {
		if strings.HasPrefix(sec.Name, "style:heading-") {
			levelStr := strings.TrimPrefix(sec.Name, "style:heading-")
			level, err := strconv.Atoi(levelStr)
			if err != nil || level < 1 || level > 6 {
				continue
			}
			if err := c.r.SetHeadingStyle(level, sec.Props); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Compiler) applyTableStyles() error {
	for _, sec := range c.doc.Sections {
		if strings.HasPrefix(sec.Name, "style:table ") {
			name := strings.TrimPrefix(sec.Name, "style:table ")
			if name == "" {
				continue
			}
			if err := c.r.SetTableStyle(name, sec.Props); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Compiler) applyMetadata() error {
	if props := c.collectSection("title"); props != nil {
		if err := c.r.SetMetadata(props); err != nil {
			return fmt.Errorf("set metadata: %w", err)
		}
	}
	return nil
}

func (c *Compiler) applyHeaderFooter() error {
	if props := c.collectSection("header"); props != nil {
		if err := c.r.SetHeader(props); err != nil {
			return fmt.Errorf("set header: %w", err)
		}
	}
	if props := c.collectSection("footer"); props != nil {
		if err := c.r.SetFooter(props); err != nil {
			return fmt.Errorf("set footer: %w", err)
		}
	}
	return nil
}

func (c *Compiler) renderSection(sec parse.Section) error {
	if err := c.validateSectionProps(sec); err != nil {
		return fmt.Errorf("section %q: %w", sec.Name, err)
	}

	switch {
	case strings.HasPrefix(sec.Name, "section:next-page"):
		if err := c.r.AddPageBreak(); err != nil {
			return err
		}
		if err := c.r.AddSectionBreak("next-page"); err != nil {
			return err
		}
	case strings.HasPrefix(sec.Name, "section:continuous"):
		if err := c.r.AddSectionBreak("continuous"); err != nil {
			return err
		}
	}

	body := c.expandLoops(sec.Body)
	body = c.resolveSectionBuiltins(body)
	body = c.applyFormats(body, sec.Props["formats"])
	body = c.ds.Resolve(body)
	body = resolveRowStyles(body) // Resolve style={{var}} after variable resolution
	if body == "" {
		return nil
	}
	return c.renderBody(body)
}

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

func (c *Compiler) applyFormats(body, formats string) string {
	if formats == "" {
		return body
	}
	fmtMap := parseFormats(formats)
	if fmtMap == nil {
		return body
	}
	return varRe.ReplaceAllStringFunc(body, func(match string) string {
		path := match[2 : len(match)-2]
		path = strings.TrimSpace(path)
		parts := strings.Split(path, ".")
		key := parts[len(parts)-1]
		fmtStr, ok := fmtMap[key]
		// For dotted format keys (e.g. items.date_field), match array paths
		// like items.0.date_field by stripping the index.
		if !ok && len(parts) >= 3 {
			sourceField := parts[0] + "." + parts[len(parts)-1]
			fmtStr, ok = fmtMap[sourceField]
		}
		if !ok {
			return match
		}
		if val, okVal := c.ds.Get(path); okVal {
			return applyFormat(fmt.Sprintf("%v", val), fmtStr)
		}
		return match
	})
}

func resolveBuiltins(s string) string {
	now := time.Now()
	date := now.Format("2006-01-02")
	s = strings.ReplaceAll(s, "{{date}}", date)
	return s
}

func (c *Compiler) resolveSectionBuiltins(s string) string {
	s = resolveBuiltins(s)
	if strings.Contains(s, "{{title}}") {
		if p := c.collectSection("title"); p != nil {
			if t := p["title"]; t != "" {
				s = strings.ReplaceAll(s, "{{title}}", t)
			}
		}
	}
	return s
}

func resolveRowStyles(body string) string {
	// Resolve style={{var}} in row/li tags after variables are already resolved
	// Pattern matches: <row style={{...}}> or <li style={{...}}>
	re := regexp.MustCompile(`<(row|li)\s+style=\{\{([^}]+)\}\}`)
	return re.ReplaceAllString(body, `<$1 style=$2`)
}
