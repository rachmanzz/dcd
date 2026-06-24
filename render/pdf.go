package render

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/pdf"
	"github.com/gpdf-dev/gpdf/template"
)

type hdrFooterEntry struct {
	text  string
	align string
}

type PdfRenderer struct {
	doc           *template.Document
	page          *template.PageBuilder
	options       []template.Option
	headerParts   []hdrFooterEntry
	footerParts   []hdrFooterEntry
	headingStyles map[int]map[string]string
	tableStyles   map[string]map[string]string
	defaultStyle  map[string]string
	docTitle      string
}

func NewPdfRenderer() *PdfRenderer {
	return &PdfRenderer{
		headingStyles: make(map[int]map[string]string),
		tableStyles:   make(map[string]map[string]string),
	}
}

func (p *PdfRenderer) SetMetadata(props map[string]string) error {
	if title := props["title"]; title != "" {
		p.docTitle = title
	}
	return nil
}

func (p *PdfRenderer) SetDefaultStyle(props map[string]string) error {
	p.defaultStyle = props
	return nil
}

func (p *PdfRenderer) init() error {
	if p.doc != nil {
		return nil
	}
	if p.options == nil {
		p.options = []template.Option{
			template.WithPageSize(document.A4),
			template.WithMargins(document.UniformEdges(document.Mm(20))),
		}
	}
	p.doc = template.New(p.options...)

	if len(p.headerParts) > 0 {
		parts := p.headerParts
		p.doc.Header(func(pb *template.PageBuilder) {
			for _, part := range parts {
				pb.AutoRow(func(r *template.RowBuilder) {
					r.Col(12, func(c *template.ColBuilder) {
						var opts []template.TextOption
						switch part.align {
						case "center":
							opts = append(opts, template.AlignCenter())
						case "right":
							opts = append(opts, template.AlignRight())
						}
						text := part.text
						text = resolveBuiltins(text)
						c.Text(text, opts...)
					})
				})
			}
		})
	}
	if len(p.footerParts) > 0 {
		parts := p.footerParts
		p.doc.Footer(func(pb *template.PageBuilder) {
			for _, part := range parts {
				pb.AutoRow(func(r *template.RowBuilder) {
					r.Col(12, func(c *template.ColBuilder) {
						var opts []template.TextOption
						switch part.align {
						case "center":
							opts = append(opts, template.AlignCenter())
						case "right":
							opts = append(opts, template.AlignRight())
						}
						text := part.text
						text = resolveBuiltins(text)
						c.Text(text, opts...)
					})
				})
			}
		})
	}
	return nil
}

func (p *PdfRenderer) SetHeadingStyle(level int, props map[string]string) error {
	p.headingStyles[level] = props
	return nil
}

func (p *PdfRenderer) SetTableStyle(name string, props map[string]string) error {
	p.tableStyles[name] = props
	return nil
}

func (p *PdfRenderer) getPage() *template.PageBuilder {
	if p.page == nil {
		p.page = p.doc.AddPage()
	}
	return p.page
}

func (p *PdfRenderer) AddHeading(text string, level int, attrs map[string]string) error {
	if err := p.init(); err != nil {
		return err
	}

	style := p.headingStyles[level]
	def := p.defaultStyle

	sizes := []float64{24, 20, 16, 14, 12, 11}
	size := 16.0
	if level >= 1 && level <= 6 {
		size = sizes[level-1]
	}
	if s := chooseAttr(def, style, attrs, "font-size"); s != "" {
		size = atof(s)
	}

	page := p.getPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			var opts []template.TextOption
			opts = append(opts, template.FontSize(size))

			s := chooseAttr(def, style, attrs, "font-color")
			if s != "" {
				if clr := parseHexColor(s); clr != nil {
					opts = append(opts, template.TextColor(*clr))
				}
			}

			if s := chooseAttr(def, style, attrs, "bold"); s == "true" {
				opts = append(opts, template.Bold())
			}
			if s := chooseAttr(def, style, attrs, "italic"); s == "true" {
				opts = append(opts, template.Italic())
			}
			if s := chooseAttr(def, style, attrs, "underline"); s == "true" {
				opts = append(opts, template.Underline())
			}

			switch chooseAttr(def, style, attrs, "align") {
			case "center":
				opts = append(opts, template.AlignCenter())
			case "right":
				opts = append(opts, template.AlignRight())
			}

			c.Text(text, opts...)

			if s := chooseAttr(def, style, attrs, "space-before"); s != "" {
				c.Spacer(document.Mm(atof(s) * 0.3528))
			}

			if s := chooseAttr(def, style, attrs, "space-after"); s != "" {
				c.Spacer(document.Mm(atof(s) * 0.3528))
			} else {
				c.Spacer(document.Mm(3))
			}

			if chooseAttr(def, style, attrs, "border-bottom") != "" {
				c.Line()
				c.Spacer(document.Mm(2))
			}
		})
	})
	return nil
}

func parseHexColor(s string) *pdf.Color {
	if !strings.HasPrefix(s, "#") || len(s) != 7 {
		return nil
	}
	hex, err := strconv.ParseUint(s[1:], 16, 32)
	if err != nil {
		return nil
	}
	c := pdf.RGBHex(uint32(hex))
	return &c
}

func (p *PdfRenderer) defaultTextOpts() []template.TextOption {
	var opts []template.TextOption
	if p.defaultStyle != nil {
		if ff := p.defaultStyle["font-family"]; ff != "" {
			opts = append(opts, template.FontFamily(ff))
		}
		if fs := p.defaultStyle["font-size"]; fs != "" {
			opts = append(opts, template.FontSize(atof(fs)))
		}
		if fc := p.defaultStyle["font-color"]; fc != "" {
			if clr := parseHexColor(fc); clr != nil {
				opts = append(opts, template.TextColor(*clr))
			}
		}
	}
	return opts
}

func (p *PdfRenderer) AddParagraph(runs []TextRun) error {
	if err := p.init(); err != nil {
		return err
	}
	if len(runs) == 0 {
		return nil
	}

	page := p.getPage()
	defaultOpts := p.defaultTextOpts()

	allPlain := true
	for _, run := range runs {
		if run.Bold || run.Italic || run.Underline || run.Code || run.Link != "" {
			allPlain = false
			break
		}
	}
	if allPlain {
		page.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text(runs[0].Text, defaultOpts...)
				c.Spacer(document.Mm(2))
			})
		})
		return nil
	}

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.RichText(func(rt *template.RichTextBuilder) {
				for _, run := range runs {
					var opts []template.TextOption
					opts = append(opts, defaultOpts...)
				if run.Code {
					opts = append(opts, template.FontFamily("Courier New"))
				}
				if run.Bold {
					opts = append(opts, template.Bold())
				}
				if run.Italic {
					opts = append(opts, template.Italic())
				}
				isLink := run.Link != ""
				if isLink && run.LinkAttrs["underline"] == "false" {
					// skip underline
				} else if isLink || run.Underline {
					opts = append(opts, template.Underline())
				}
				if isLink {
					linkColor := uint32(0x0055CC)
					if run.LinkAttrs != nil && run.LinkAttrs["color"] != "" {
						hex := strings.TrimPrefix(run.LinkAttrs["color"], "#")
						if v, err := strconv.ParseUint(hex, 16, 32); err == nil {
							linkColor = uint32(v)
						}
					}
					opts = append(opts, template.TextColor(pdf.RGBHex(linkColor)))
				}
					rt.Span(run.Text, opts...)
				}
			})
			c.Spacer(document.Mm(2))
		})
	})
	return nil
}

func (p *PdfRenderer) AddLineBreak() error {
	if err := p.init(); err != nil {
		return err
	}
	page := p.getPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(3))
		})
	})
	return nil
}

func (p *PdfRenderer) AddHorizontalRule(attrs map[string]string) error {
	if err := p.init(); err != nil {
		return err
	}
	page := p.getPage()
	page.AutoRow(func(r *template.RowBuilder) {
		fullWidth := true
		colUnits := 12
		if attrs != nil {
			if w := attrs["width"]; strings.HasSuffix(w, "%") {
				if pct, err := strconv.ParseFloat(strings.TrimSuffix(w, "%"), 64); err == nil && pct > 0 && pct < 100 {
					colUnits = int(12 * pct / 100)
					if colUnits < 1 {
						colUnits = 1
					}
					offset := (12 - colUnits) / 2
					if offset > 0 {
						r.Col(offset, func(c *template.ColBuilder) {})
					}
					fullWidth = false
				}
			}
		}
		r.Col(colUnits, func(c *template.ColBuilder) {
			c.Line()
			c.Spacer(document.Mm(3))
		})
		if !fullWidth {
			remaining := 12 - colUnits - (12-colUnits)/2
			if remaining > 0 {
				r.Col(remaining, func(c *template.ColBuilder) {})
			}
		}
	})
	return nil
}

func (p *PdfRenderer) AddPageBreak() error {
	if err := p.init(); err != nil {
		return err
	}
	p.page = p.doc.AddPage()
	return nil
}

func (p *PdfRenderer) AddImage(src string, attrs map[string]string) error {
	if err := p.init(); err != nil {
		return err
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	page := p.getPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			var opts []template.ImageOption
			if attrs != nil {
				if w := attrs["width"]; w != "" {
					if strings.HasSuffix(w, "%") {
						if pct, err := strconv.ParseFloat(strings.TrimSuffix(w, "%"), 64); err == nil {
							opts = append(opts, template.FitWidth(document.Mm(float64(160*pct/100))))
						}
					} else {
						opts = append(opts, template.FitWidth(document.Mm(atof(w))))
					}
				}
				switch attrs["align"] {
				case "center":
					opts = append(opts, template.WithAlign(document.AlignCenter))
				case "right":
					opts = append(opts, template.WithAlign(document.AlignRight))
				}
			}
			if len(opts) == 0 {
				opts = append(opts, template.FitWidth(document.Mm(150)))
			}
			c.Image(data, opts...)
		})
	})
	return nil
}

func (p *PdfRenderer) AddHyperlink(text, url string, attrs map[string]string) error {
	if err := p.init(); err != nil {
		return err
	}
	page := p.getPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			linkColor := uint32(0x0055CC)
			hasUnderline := true
			if attrs != nil {
				if hex := strings.TrimPrefix(attrs["color"], "#"); hex != "" {
					if v, err := strconv.ParseUint(hex, 16, 32); err == nil {
						linkColor = uint32(v)
					}
				}
				if attrs["underline"] == "false" {
					hasUnderline = false
				}
			}
			var opts []template.TextOption
			if hasUnderline {
				opts = append(opts, template.Underline())
			}
			opts = append(opts, template.TextColor(pdf.RGBHex(linkColor)))
			c.Text(text, opts...)
			c.Spacer(document.Mm(2))
		})
	})
	_ = url
	return nil
}

func (p *PdfRenderer) AddList(items []ListItem, ordered bool) error {
	if err := p.init(); err != nil {
		return err
	}
	
	// Flatten to plain text strings for gpdf List/OrderedList
	// These APIs don't support rich text formatting per item yet
	var texts []string
	var flatten func(items []ListItem)
	flatten = func(items []ListItem) {
		for _, item := range items {
			if len(item.Runs) > 0 {
				var buf strings.Builder
				for _, run := range item.Runs {
					buf.WriteString(run.Text)
				}
				texts = append(texts, buf.String())
			}
			if len(item.Items) > 0 {
				flatten(item.Items)
			}
		}
	}
	flatten(items)

	page := p.getPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			if ordered {
				c.OrderedList(texts)
			} else {
				c.List(texts)
			}
			c.Spacer(document.Mm(3))
		})
	})
	return nil
}

func (p *PdfRenderer) AddTable(rows []TableRow, attrs map[string]string) error {
	if err := p.init(); err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	page := p.getPage()
	
	// Convert runs to plain text for gpdf table API
	runsToText := func(runs []TextRun) string {
		var buf strings.Builder
		for _, run := range runs {
			buf.WriteString(run.Text)
		}
		return buf.String()
	}
	
	var header []string
	var data [][]string
	if len(rows) == 1 {
		data = make([][]string, 1)
		for _, cell := range rows[0].Cells {
			data[0] = append(data[0], runsToText(cell.Runs))
		}
	} else {
		for i, r := range rows {
			styleName := r.Props["style"]
			style := p.tableStyles[styleName]
			cells := make([]string, len(r.Cells))
			for j, cell := range r.Cells {
				cells[j] = runsToText(cell.Runs)
			}
			if style != nil && i == 0 {
				header = cells
				continue
			}
			data = append(data, cells)
		}
		if header == nil {
			header = make([]string, len(rows[0].Cells))
			for j, cell := range rows[0].Cells {
				header[j] = runsToText(cell.Runs)
			}
			for _, r := range rows[1:] {
				cells := make([]string, len(r.Cells))
				for j, cell := range r.Cells {
					cells[j] = runsToText(cell.Runs)
				}
				data = append(data, cells)
			}
		}
	}
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Table(header, data)
			c.Spacer(document.Mm(3))
		})
	})
	return nil
}

func (p *PdfRenderer) AddWrappedParagraph(text string, flags string) error {
	if err := p.init(); err != nil {
		return err
	}

	page := p.getPage()
	var opts []template.TextOption
	for _, f := range strings.Split(flags, "|") {
		switch f {
		case "c":
			opts = append(opts, template.AlignCenter())
		case "b":
			opts = append(opts, template.Bold())
		}
	}

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text(text, opts...)
			c.Spacer(document.Mm(2))
		})
	})
	return nil
}

func (p *PdfRenderer) SetPageStyle(props map[string]string) error {
	if props == nil {
		return nil
	}

	wMm, hMm := parsePageSize(props["layout"], props["orientation"], props["w"], props["h"])
	pageSize := document.Size{
		Width:  wMm * 72 / 25.4,
		Height: hMm * 72 / 25.4,
	}

	l, r, t, b := computeMargins(props)
	margins := document.Edges{
		Top:    document.Mm(t),
		Bottom: document.Mm(b),
		Left:   document.Mm(l),
		Right:  document.Mm(r),
	}

	p.options = []template.Option{
		template.WithPageSize(pageSize),
		template.WithMargins(margins),
	}
	return nil
}

func (p *PdfRenderer) SetHeader(props map[string]string) error {
	now := time.Now()
	for _, align := range []string{"left", "center", "right"} {
		if t, ok := props[align]; ok && t != "" {
			t = strings.ReplaceAll(t, "{{date}}", now.Format("2006-01-02"))
			t = strings.ReplaceAll(t, "{{page}}", "{#page}")
			t = strings.ReplaceAll(t, "{{title}}", p.docTitle)
			t = strings.ReplaceAll(t, "{{total}}", "{#total}")
			p.headerParts = append(p.headerParts, hdrFooterEntry{text: t, align: align})
		}
	}
	return nil
}

func (p *PdfRenderer) SetFooter(props map[string]string) error {
	now := time.Now()
	for _, align := range []string{"left", "center", "right"} {
		if t, ok := props[align]; ok && t != "" {
			t = strings.ReplaceAll(t, "{{date}}", now.Format("2006-01-02"))
			t = strings.ReplaceAll(t, "{{page}}", "{#page}")
			t = strings.ReplaceAll(t, "{{title}}", p.docTitle)
			t = strings.ReplaceAll(t, "{{total}}", "{#total}")
			p.footerParts = append(p.footerParts, hdrFooterEntry{text: t, align: align})
		}
	}
	return nil
}

func (p *PdfRenderer) Save(path string) error {
	if err := p.init(); err != nil {
		return err
	}
	data, err := p.doc.Generate()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
