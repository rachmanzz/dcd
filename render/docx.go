package render

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/common/units"
	"github.com/gomutex/godocx/docx"
	"github.com/gomutex/godocx/wml/ctypes"
	"github.com/gomutex/godocx/wml/stypes"
)

type DocxRenderer struct {
	root          *docx.RootDoc
	headerRID     string
	footerRID     string
	headingStyles map[int]map[string]string
	tableStyles   map[string]map[string]string
	defaultStyle  map[string]string
	docTitle      string
	pageWidthMm   float64
	unit          string
}

func NewDocxRenderer() *DocxRenderer {
	return &DocxRenderer{
		headingStyles: make(map[int]map[string]string),
		tableStyles:   make(map[string]map[string]string),
	}
}

func (d *DocxRenderer) SetMetadata(props map[string]string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	title := props["title"]
	d.docTitle = title
	subject := props["subject"]
	author := props["author"]

	// Build core.xml
	coreXML := `<?xml version="1.0" encoding="UTF-8"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">`
	if title != "" {
		coreXML += fmt.Sprintf("<dc:title>%s</dc:title>", title)
	}
	if subject != "" {
		coreXML += fmt.Sprintf("<dc:subject>%s</dc:subject>", subject)
	}
	if author != "" {
		coreXML += fmt.Sprintf("<cp:creator>%s</cp:creator>", author)
	}
	coreXML += `</cp:coreProperties>`

	d.root.FileMap.Store("docProps/core.xml", []byte(coreXML))

	// Build app.xml with title
	appXML := `<?xml version="1.0" encoding="UTF-8"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes">`
	if title != "" {
		appXML += fmt.Sprintf("<Title>%s</Title>", title)
	}
	appXML += `</Properties>`

	d.root.FileMap.Store("docProps/app.xml", []byte(appXML))

	// Add content type overrides if not present
	d.root.ContentType.AddOverride("/docProps/core.xml", "application/vnd.openxmlformats-package.core-properties+xml")
	d.root.ContentType.AddOverride("/docProps/app.xml", "application/vnd.openxmlformats-officedocument.extended-properties+xml")

	return nil
}

func (d *DocxRenderer) SetDefaultStyle(props map[string]string) error {
	d.defaultStyle = props
	return nil
}

func (d *DocxRenderer) init() error {
	root, err := godocx.NewDocument()
	if err != nil {
		return err
	}
	d.root = root
	return nil
}

func (d *DocxRenderer) SetHeadingStyle(level int, props map[string]string) error {
	d.headingStyles[level] = props
	return nil
}

func (d *DocxRenderer) SetTableStyle(name string, props map[string]string) error {
	d.tableStyles[name] = props
	return nil
}

func (d *DocxRenderer) AddHeading(text string, level int, attrs map[string]string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	p, err := d.root.AddHeading(text, uint(level))
	if err != nil {
		return err
	}

	style := d.headingStyles[level]
	def := d.defaultStyle
	if style == nil && attrs == nil && def == nil {
		return nil
	}

	// Paragraph-level properties
	pPr := p.GetCT().Property
	if pPr == nil {
		pPr = &ctypes.ParagraphProp{}
		p.GetCT().Property = pPr
	}

	if s := chooseAttr(def, style, attrs, "align"); s != "" {
		switch s {
		case "center":
			pPr.Justification = ctypes.NewGenSingleStrVal(stypes.JustificationCenter)
		case "right":
			pPr.Justification = ctypes.NewGenSingleStrVal(stypes.JustificationRight)
		default:
			pPr.Justification = ctypes.NewGenSingleStrVal(stypes.JustificationLeft)
		}
	}

	if s := chooseAttr(def, style, attrs, "space-before"); s != "" {
		if pPr.Spacing == nil {
			pPr.Spacing = &ctypes.Spacing{}
		}
		v := uint64(atof(s) * 20)
		pPr.Spacing.Before = &v
	}

	if s := chooseAttr(def, style, attrs, "space-after"); s != "" {
		if pPr.Spacing == nil {
			pPr.Spacing = &ctypes.Spacing{}
		}
		v := uint64(atof(s) * 20)
		pPr.Spacing.After = &v
	}

		if s := chooseAttr(def, style, attrs, "border-bottom"); s != "" {
		if pPr.Border == nil {
			pPr.Border = &ctypes.ParaBorder{}
		}
		color := "auto"
		pPr.Border.Bottom = &ctypes.Border{
			Val:   stypes.BorderStyleSingle,
			Color: &color,
		}
	}

	// Run-level properties
	for _, child := range p.GetCT().Children {
		if child.Run == nil {
			continue
		}
		if child.Run.Property == nil {
			child.Run.Property = &ctypes.RunProperty{}
		}
		prop := child.Run.Property

		if s := chooseAttr(def, style, attrs, "font-color"); s != "" {
			prop.Color = ctypes.NewColor(s)
		}
		if s := chooseAttr(def, style, attrs, "font-size"); s != "" {
			prop.Size = ctypes.NewFontSize(uint64(atof(s) * 2))
		}
		if s := chooseAttr(def, style, attrs, "bold"); s != "" {
			prop.Bold = ctypes.OnOffFromBool(s == "true")
		}
		if s := chooseAttr(def, style, attrs, "italic"); s != "" {
			prop.Italic = ctypes.OnOffFromBool(s == "true")
		}
		if s := chooseAttr(def, style, attrs, "underline"); s != "" {
			prop.Underline = ctypes.NewGenSingleStrVal(stypes.UnderlineSingle)
		}
		if s := chooseAttr(def, style, attrs, "font-family"); s != "" {
			prop.Fonts = &ctypes.RunFonts{
				Ascii: s,
				HAnsi: s,
			}
		}
	}
	return nil
}

func (d *DocxRenderer) AddParagraph(runs []TextRun) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	p := d.root.AddEmptyParagraph()

	if d.defaultStyle != nil {
		pPr := p.GetCT().Property
		if pPr == nil {
			pPr = &ctypes.ParagraphProp{}
			p.GetCT().Property = pPr
		}
		if lh := d.defaultStyle["line-height"]; lh != "" {
			if pPr.Spacing == nil {
				pPr.Spacing = &ctypes.Spacing{}
			}
			v := int(atof(lh) * 240)
			pPr.Spacing.Line = &v
		}
	}

	for _, r := range runs {
		if r.Tab {
			p.AddText("")
			ctRun := p.GetCT().Children[len(p.GetCT().Children)-1].Run
			ctRun.Children = []ctypes.RunChild{{Tab: &ctypes.Empty{}}}
			continue
		}
		run := p.AddText(r.Text)
		ctRun := p.GetCT().Children[len(p.GetCT().Children)-1].Run

		if d.defaultStyle != nil {
			if fc := d.defaultStyle["font-color"]; fc != "" {
				run.Color(fc)
			}
			if fs := d.defaultStyle["font-size"]; fs != "" {
				run.Size(uint64(atof(fs)))
			}
			if ff := d.defaultStyle["font-family"]; ff != "" {
				if ctRun.Property == nil {
					ctRun.Property = &ctypes.RunProperty{}
				}
				ctRun.Property.Fonts = &ctypes.RunFonts{
					Ascii: ff,
					HAnsi: ff,
				}
			}
		}

		if r.Code {
			if ctRun.Property == nil {
				ctRun.Property = &ctypes.RunProperty{}
			}
			ctRun.Property.Fonts = &ctypes.RunFonts{
				Ascii: "Courier New",
				HAnsi: "Courier New",
			}
		}
		if r.Bold {
			run.Bold(true)
		}
		if r.Italic {
			run.Italic(true)
		}
		isLink := r.Link != ""
		if isLink && r.LinkAttrs["underline"] == "false" {
			// skip underline
		} else if isLink || r.Underline {
			run.Underline(stypes.UnderlineSingle)
		}
		if isLink {
			linkColor := "0055CC"
			if r.LinkAttrs != nil && r.LinkAttrs["color"] != "" {
				linkColor = strings.TrimPrefix(r.LinkAttrs["color"], "#")
			}
			run.Color(linkColor)
		}
	}
	return nil
}

func (d *DocxRenderer) AddLineBreak() error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	p := d.root.AddEmptyParagraph()
	p.AddText("")
	return nil
}

func (d *DocxRenderer) AddHorizontalRule(attrs map[string]string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}

	color := "CCCCCC"
	if attrs != nil && attrs["color"] != "" {
		color = strings.TrimPrefix(attrs["color"], "#")
	}

	p := d.root.AddEmptyParagraph()
	pPr := p.GetCT().Property
	if pPr == nil {
		pPr = &ctypes.ParagraphProp{}
		p.GetCT().Property = pPr
	}

	if pPr.Border == nil {
		pPr.Border = &ctypes.ParaBorder{}
	}

	bdr := &ctypes.Border{
		Val:   stypes.BorderStyleSingle,
		Color: &color,
	}
	pPr.Border.Bottom = bdr

	// Width as percentage of page width
	if attrs != nil {
		if w := attrs["width"]; strings.HasSuffix(w, "%") {
			if pct, err := strconv.ParseFloat(strings.TrimSuffix(w, "%"), 64); err == nil && d.pageWidthMm > 0 {
				indentMm := d.pageWidthMm * (100 - pct) / 200
				indentTwip := int(indentMm * 56.7)
				if pPr.Indent == nil {
					pPr.Indent = &ctypes.Indent{}
				}
				pPr.Indent.Left = &indentTwip
				pPr.Indent.Right = &indentTwip
			}
		}
	}

	return nil
}

func (d *DocxRenderer) AddPageBreak() error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	d.root.AddPageBreak()
	return nil
}

func (d *DocxRenderer) AddImage(src string, attrs map[string]string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	width := units.Inch(5)
	height := units.Inch(4)
	if attrs != nil {
		if w := attrs["width"]; w != "" {
			if strings.HasSuffix(w, "%") {
				if pct, err := strconv.ParseFloat(strings.TrimSuffix(w, "%"), 64); err == nil && d.pageWidthMm > 0 {
					width = units.Inch(d.pageWidthMm * pct / 100 / 25.4)
				}
			} else {
				width = units.Inch(atof(w))
			}
		}
		if h := attrs["height"]; h != "" {
			height = units.Inch(atof(h))
		}
	}
	pic, err := d.root.AddPicture(src, width, height)
	if err != nil {
		return err
	}

	if attrs != nil {
		// Align
		switch attrs["align"] {
		case "center":
			pic.Para.Justification(stypes.JustificationCenter)
		case "right":
			pic.Para.Justification(stypes.JustificationRight)
		}

		// Alt text
		if alt := attrs["alt"]; alt != "" {
			pic.Inline.DocProp.Description = alt
		}

		// Border
		if attrs["border"] != "" {
			pPr := pic.Para.GetCT().Property
			if pPr == nil {
				pPr = &ctypes.ParagraphProp{}
				pic.Para.GetCT().Property = pPr
			}
			if pPr.Border == nil {
				pPr.Border = &ctypes.ParaBorder{}
			}
			color := "auto"
			bdr := &ctypes.Border{Val: stypes.BorderStyleSingle, Color: &color}
			pPr.Border.Top = bdr
			pPr.Border.Bottom = bdr
			pPr.Border.Left = bdr
			pPr.Border.Right = bdr
		}

		// Shading
		if sh := attrs["shading"]; sh != "" {
			pPr := pic.Para.GetCT().Property
			if pPr == nil {
				pPr = &ctypes.ParagraphProp{}
				pic.Para.GetCT().Property = pPr
			}
			fill := sh
			pPr.Shading = &ctypes.Shading{
				Fill: &fill,
			}
		}
	}

	return nil
}

func (d *DocxRenderer) AddHyperlink(text, url string, attrs map[string]string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}

	rID := d.root.Document.IncRelationID()
	relID := "rId" + strconv.Itoa(rID)

	targetMode := "External"
	target := url
	if strings.HasPrefix(url, "#") {
		targetMode = "Internal"
		if d.root.Document.DocRels.Relationships == nil {
			_ = d.root.Document.IncRelationID()
		}
	}

	d.root.Document.DocRels.Relationships = append(d.root.Document.DocRels.Relationships, &docx.Relationship{
		ID:         relID,
		TargetMode: targetMode,
		Type:       "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink",
		Target:     target,
	})

	linkColor := "0055CC"
	hasUnderline := true
	if attrs != nil {
		if c := attrs["color"]; c != "" {
			linkColor = strings.TrimPrefix(c, "#")
		}
		if attrs["underline"] == "false" {
			hasUnderline = false
		}
	}

	runProp := &ctypes.RunProperty{
		Style: ctypes.NewRunStyle("a1"),
		Color: ctypes.NewColor(linkColor),
	}
	if hasUnderline {
		runProp.Underline = ctypes.NewGenSingleStrVal(stypes.UnderlineSingle)
	}

	run := &ctypes.Run{
		Children: []ctypes.RunChild{
			{Text: ctypes.TextFromString(text)},
		},
		Property: runProp,
	}

	hyperlink := &ctypes.Hyperlink{
		ID: relID,
		Children: []ctypes.ParagraphChild{
			{Run: run},
		},
	}

	p := d.root.AddEmptyParagraph()
	p.GetCT().Children = append(p.GetCT().Children, ctypes.ParagraphChild{
		Link: hyperlink,
	})

	return nil
}

func (d *DocxRenderer) AddList(items []ListItem, ordered bool) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	style := "List Bullet"
	if ordered {
		style = "List Number"
	}
	for _, item := range items {
		if len(item.Runs) > 0 {
			p := d.root.AddEmptyParagraph()
			p.Style(style)
			for _, run := range item.Runs {
				r := p.AddText(run.Text)
				ctRun := p.GetCT().Children[len(p.GetCT().Children)-1].Run

				if d.defaultStyle != nil {
					if fc := d.defaultStyle["font-color"]; fc != "" {
						r.Color(fc)
					}
					if fs := d.defaultStyle["font-size"]; fs != "" {
						r.Size(uint64(atof(fs)))
					}
					if ff := d.defaultStyle["font-family"]; ff != "" {
						if ctRun.Property == nil {
							ctRun.Property = &ctypes.RunProperty{}
						}
						ctRun.Property.Fonts = &ctypes.RunFonts{
							Ascii: ff,
							HAnsi: ff,
						}
					}
				}

				if run.Code {
					if ctRun.Property == nil {
						ctRun.Property = &ctypes.RunProperty{}
					}
					ctRun.Property.Fonts = &ctypes.RunFonts{
						Ascii: "Courier New",
						HAnsi: "Courier New",
					}
				}
				if run.Bold {
					r.Bold(true)
				}
				if run.Italic {
					r.Italic(true)
				}
				if run.Underline {
					r.Underline(stypes.UnderlineSingle)
				}
			}
		}
		if len(item.Items) > 0 {
			if err := d.AddList(item.Items, ordered); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *DocxRenderer) AddTable(rows []TableRow, attrs map[string]string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	table := d.root.AddTable()
	if attrs["border"] != "" {
		table.Style("TableGrid")
	} else {
		table.Style("LightList-Accent4")
	}
	for _, row := range rows {
		tblRow := table.AddRow()
		styleName := row.Props["style"]
		rowStyle := d.tableStyles[styleName]
		rowShading := row.Props["shading"]
		for _, cell := range row.Cells {
			c := tblRow.AddCell()
			p := c.AddParagraph("")

			// Build runs from cell.Runs
			for _, run := range cell.Runs {
				r := p.AddText(run.Text)
				ctRun := p.GetCT().Children[len(p.GetCT().Children)-1].Run

				if d.defaultStyle != nil {
					if fc := d.defaultStyle["font-color"]; fc != "" {
						r.Color(fc)
					}
					if fs := d.defaultStyle["font-size"]; fs != "" {
						r.Size(uint64(atof(fs)))
					}
					if ff := d.defaultStyle["font-family"]; ff != "" {
						if ctRun.Property == nil {
							ctRun.Property = &ctypes.RunProperty{}
						}
						ctRun.Property.Fonts = &ctypes.RunFonts{
							Ascii: ff,
							HAnsi: ff,
						}
					}
				}

				if run.Code {
					if ctRun.Property == nil {
						ctRun.Property = &ctypes.RunProperty{}
					}
					ctRun.Property.Fonts = &ctypes.RunFonts{
						Ascii: "Courier New",
						HAnsi: "Courier New",
					}
				}
				if run.Bold {
					r.Bold(true)
				}
				if run.Italic {
					r.Italic(true)
				}
				if run.Underline {
					r.Underline(stypes.UnderlineSingle)
				}

				// Apply named style font properties to each run
				if rowStyle != nil {
					if fc := rowStyle["font-color"]; fc != "" {
						r.Color(fc)
					}
					if rowStyle["font-weight"] == "bold" {
						r.Bold(true)
					}
				}
			}

			// Determine shading: row style > cell attr > row prop
			shading := cell.Attrs["shading"]
			if shading == "" && rowStyle != nil && rowStyle["shading"] != "" {
				shading = rowStyle["shading"]
			}
			if shading == "" && rowShading != "" {
				shading = rowShading
			}
			if shading != "" {
				pPr := p.GetCT().Property
				if pPr == nil {
					pPr = &ctypes.ParagraphProp{}
					p.GetCT().Property = pPr
				}
				fill := shading
				pPr.Shading = &ctypes.Shading{
					Fill: &fill,
				}
			}

			// Border-bottom from named row style
			if rowStyle != nil && rowStyle["border-bottom"] != "" {
				pPr := p.GetCT().Property
				if pPr == nil {
					pPr = &ctypes.ParagraphProp{}
					p.GetCT().Property = pPr
				}
				if pPr.Border == nil {
					pPr.Border = &ctypes.ParaBorder{}
				}
				color := "auto"
				pPr.Border.Bottom = &ctypes.Border{
					Val:   stypes.BorderStyleSingle,
					Color: &color,
				}
			}

			// Determine alignment: cell attr > row style > row prop
			align := cell.Attrs["align"]
			if align == "" && rowStyle != nil && rowStyle["align"] != "" {
				align = rowStyle["align"]
			}
			switch align {
			case "center":
				p.Justification(stypes.JustificationCenter)
			case "right":
				p.Justification(stypes.JustificationRight)
			}
		}
	}
	return nil
}

func (d *DocxRenderer) AddWrappedParagraph(text string, flags string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	p := d.root.AddParagraph(text)
	for _, f := range strings.Split(flags, "|") {
		switch f {
		case "c":
			p.Justification(stypes.JustificationCenter)
		case "b":
			if len(p.GetCT().Children) > 0 {
				if run := p.GetCT().Children[0].Run; run != nil {
					if run.Property == nil {
						run.Property = &ctypes.RunProperty{}
					}
					run.Property.Bold = ctypes.OnOffFromBool(true)
				}
			}
		case "i":
			if len(p.GetCT().Children) > 0 {
				if run := p.GetCT().Children[0].Run; run != nil {
					if run.Property == nil {
						run.Property = &ctypes.RunProperty{}
					}
					run.Property.Italic = ctypes.OnOffFromBool(true)
				}
			}
		case "u":
			if len(p.GetCT().Children) > 0 {
				if run := p.GetCT().Children[0].Run; run != nil {
					if run.Property == nil {
						run.Property = &ctypes.RunProperty{}
					}
					run.Property.Underline = ctypes.NewGenSingleStrVal(stypes.UnderlineSingle)
				}
			}
		}
	}
	return nil
}

func (d *DocxRenderer) SetPageStyle(props map[string]string) error {
	if props == nil {
		return nil
	}
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}

	d.unit = props["unit"]
	w, h := parsePageSize(props["layout"], props["orientation"], props["w"], props["h"])
	l, r, t, b := computeMargins(props)

	mmToTwip := func(mm float64) uint64 {
		return uint64(mm * 56.7)
	}
	mmToTwipInt := func(mm float64) int {
		return int(mm * 56.7)
	}

	wTwip := mmToTwip(w)
	hTwip := mmToTwip(h)

	pageSize := &ctypes.PageSize{
		Width:  &wTwip,
		Height: &hTwip,
	}
	if strings.ToLower(props["orientation"]) == "landscape" {
		pageSize.Orient = stypes.PageOrientLandscape
	} else {
		pageSize.Orient = stypes.PageOrientPortrait
	}

	pageMargin := &ctypes.PageMargin{
		Left:   intPtr(mmToTwipInt(l)),
		Right:  intPtr(mmToTwipInt(r)),
		Top:    intPtr(mmToTwipInt(t)),
		Bottom: intPtr(mmToTwipInt(b)),
	}

	sectPr := &ctypes.SectionProp{
		PageSize:   pageSize,
		PageMargin: pageMargin,
	}

	d.root.Document.Body.SectPr = sectPr
	d.pageWidthMm = w
	return nil
}

func (d *DocxRenderer) SetHeader(props map[string]string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}

	parts, cfg := d.collectHeaderFooterParts(props)
	if len(parts) == 0 {
		return nil
	}

	d.headerRID = fmt.Sprintf("rId%d", d.root.Document.IncRelationID())
	xml := d.buildHdrXML(parts, cfg)

	d.root.FileMap.Store("word/header1.xml", []byte(xml))

	d.root.Document.DocRels.Relationships = append(d.root.Document.DocRels.Relationships, &docx.Relationship{
		ID:     d.headerRID,
		Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/header",
		Target: "header1.xml",
	})

	d.root.ContentType.AddOverride("/word/header1.xml", "application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml")

	if d.root.Document.Body.SectPr == nil {
		d.root.Document.Body.SectPr = ctypes.NewSectionProper()
	}
	d.root.Document.Body.SectPr.HeaderReference = &ctypes.HeaderReference{
		Type: stypes.HdrFtrDefault,
		ID:   d.headerRID,
	}

	if cfg.margin != "" {
		if d.root.Document.Body.SectPr.PageMargin == nil {
			d.root.Document.Body.SectPr.PageMargin = &ctypes.PageMargin{}
		}
		scale := unitScale(d.unit)
		m := atof(cfg.margin) * scale
		headerTwip := int(m * 56.7)
		d.root.Document.Body.SectPr.PageMargin.Header = &headerTwip
	}

	if !cfg.firstPage {
		d.root.Document.Body.SectPr.TitlePg = ctypes.NewGenSingleStrVal(stypes.OnOff("true"))
	}
	return nil
}

func (d *DocxRenderer) SetFooter(props map[string]string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}

	parts, cfg := d.collectHeaderFooterParts(props)
	if len(parts) == 0 {
		return nil
	}

	d.footerRID = fmt.Sprintf("rId%d", d.root.Document.IncRelationID())
	xml := d.buildFtrXML(parts, cfg)

	d.root.FileMap.Store("word/footer1.xml", []byte(xml))

	d.root.Document.DocRels.Relationships = append(d.root.Document.DocRels.Relationships, &docx.Relationship{
		ID:     d.footerRID,
		Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer",
		Target: "footer1.xml",
	})

	d.root.ContentType.AddOverride("/word/footer1.xml", "application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml")

	if d.root.Document.Body.SectPr == nil {
		d.root.Document.Body.SectPr = ctypes.NewSectionProper()
	}
	d.root.Document.Body.SectPr.FooterReference = &ctypes.FooterReference{
		Type: stypes.HdrFtrDefault,
		ID:   d.footerRID,
	}

	if cfg.margin != "" {
		if d.root.Document.Body.SectPr.PageMargin == nil {
			d.root.Document.Body.SectPr.PageMargin = &ctypes.PageMargin{}
		}
		scale := unitScale(d.unit)
		m := atof(cfg.margin) * scale
		footerTwip := int(m * 56.7)
		d.root.Document.Body.SectPr.PageMargin.Footer = &footerTwip
	}

	if !cfg.firstPage {
		d.root.Document.Body.SectPr.TitlePg = ctypes.NewGenSingleStrVal(stypes.OnOff("true"))
	}
	return nil
}

type hdrSegment struct {
	isField bool
	isTab   bool
	content string
}

type hdrPart struct {
	segments []hdrSegment
	align    string
	tabStops []int
}

type hdrFooterCfg struct {
	fontFamily string
	fontSize   string
	fontColor  string
	border     string
	margin     string
	firstPage  bool
	mirror     bool
}

func splitComma(s string) []string {
	var parts []string
	var cur strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) && s[i+1] == ',' {
			cur.WriteByte(',')
			i++
		} else if s[i] == ',' {
			parts = append(parts, strings.TrimSpace(cur.String()))
			cur.Reset()
		} else {
			cur.WriteByte(s[i])
		}
	}
	if v := strings.TrimSpace(cur.String()); v != "" {
		parts = append(parts, v)
	}
	return parts
}

func (d *DocxRenderer) hdrTabPositions(n int) []int {
	pw := d.pageWidthMm
	if pw == 0 {
		pw = 210.0
	}
	pwTwip := int(pw * 56.7)
	l, r := 0, 0
	if d.root != nil && d.root.Document.Body.SectPr != nil && d.root.Document.Body.SectPr.PageMargin != nil {
		if d.root.Document.Body.SectPr.PageMargin.Left != nil {
			l = *d.root.Document.Body.SectPr.PageMargin.Left
		}
		if d.root.Document.Body.SectPr.PageMargin.Right != nil {
			r = *d.root.Document.Body.SectPr.PageMargin.Right
		}
	}
	usable := pwTwip - l - r
	if usable < 1 {
		usable = pwTwip
	}
	switch n {
	case 2:
		return []int{usable}
	case 3:
		return []int{usable / 2, usable}
	default:
		return nil
	}
}

func (d *DocxRenderer) collectHeaderFooterParts(props map[string]string) ([]hdrPart, hdrFooterCfg) {
	now := time.Now()
	var parts []hdrPart

	if jb, ok := props["justify_between"]; ok && jb != "" {
		cols := splitComma(jb)
		if len(cols) >= 2 && len(cols) <= 3 {
			var segs []hdrSegment
			for i, col := range cols {
				if i > 0 {
					segs = append(segs, hdrSegment{isTab: true})
				}
				segs = append(segs, d.parseHdrText(col, now)...)
			}
			parts = append(parts, hdrPart{
				segments: segs,
				tabStops: d.hdrTabPositions(len(cols)),
			})
		}
	} else {
		for _, align := range []string{"left", "center", "right"} {
			if t, ok := props[align]; ok && t != "" {
				parts = append(parts, hdrPart{
					segments: d.parseHdrText(t, now),
					align:    align,
				})
			}
		}
	}

	cfg := hdrFooterCfg{
		fontFamily: props["font-family"],
		fontSize:   props["font-size"],
		fontColor:  props["font-color"],
		border:     props["border"],
		margin:     props["margin"],
		firstPage:  props["first-page"] == "true",
		mirror:     props["mirror"] == "true",
	}
	return parts, cfg
}

func (d *DocxRenderer) parseHdrText(t string, now time.Time) []hdrSegment {
	var segs []hdrSegment
	for {
		start := strings.Index(t, "{{")
		if start < 0 {
			if t != "" {
				segs = append(segs, hdrSegment{content: t})
			}
			break
		}
		if start > 0 {
			segs = append(segs, hdrSegment{content: t[:start]})
		}
		end := strings.Index(t[start:], "}}")
		if end < 0 {
			segs = append(segs, hdrSegment{content: t})
			break
		}
		end = start + end + 2
		switch strings.TrimSpace(t[start+2 : end-2]) {
		case "page":
			segs = append(segs, hdrSegment{isField: true, content: " PAGE "})
		case "total":
			segs = append(segs, hdrSegment{isField: true, content: " NUMPAGES "})
		case "date":
			segs = append(segs, hdrSegment{content: now.Format("2006-01-02")})
		case "title":
			segs = append(segs, hdrSegment{content: d.docTitle})
		default:
			segs = append(segs, hdrSegment{content: t[start:end]})
		}
		t = t[end:]
	}
	return segs
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func (d *DocxRenderer) buildHdrXML(parts []hdrPart, cfg hdrFooterCfg) string {
	hasBorder := cfg.border == "top" || cfg.border == "bottom"
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString(`<w:hdr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">`)
	for _, part := range parts {
		sb.WriteString(`<w:p>`)
		sb.WriteString(`<w:pPr>`)
		switch part.align {
		case "center":
			sb.WriteString(`<w:jc w:val="center"/>`)
		case "right":
			sb.WriteString(`<w:jc w:val="right"/>`)
		}
		if len(part.tabStops) > 0 {
			sb.WriteString(`<w:tabs>`)
			for i, pos := range part.tabStops {
				val := "right"
				if len(part.tabStops) == 2 && i == 0 {
					val = "center"
				}
				sb.WriteString(fmt.Sprintf(`<w:tab w:val="%s" w:pos="%d"/>`, val, pos))
			}
			sb.WriteString(`</w:tabs>`)
		}
		if hasBorder {
			bdrPos := "bottom"
			if cfg.border == "top" {
				bdrPos = "top"
			}
			sb.WriteString(fmt.Sprintf(`<w:pBdr><w:%s w:val="single" w:color="auto" w:space="4"/></w:pBdr>`, bdrPos))
		}
		sb.WriteString(`</w:pPr>`)
		for _, seg := range part.segments {
			if seg.isTab {
				sb.WriteString(`<w:r><w:tab/></w:r>`)
			} else if seg.isField {
				sb.WriteString(fmt.Sprintf(`<w:fldSimple w:instr="%s"/>`, seg.content))
			} else {
				sb.WriteString(`<w:r>`)
				if d.hasRunProps(cfg) {
					sb.WriteString(d.runPropsXML(cfg))
				}
				sb.WriteString(fmt.Sprintf(`<w:t xml:space="preserve">%s</w:t>`, xmlEscape(seg.content)))
				sb.WriteString(`</w:r>`)
			}
		}
		sb.WriteString(`</w:p>`)
	}
	sb.WriteString(`</w:hdr>`)
	return sb.String()
}

func (d *DocxRenderer) buildFtrXML(parts []hdrPart, cfg hdrFooterCfg) string {
	hasBorder := cfg.border == "top" || cfg.border == "bottom"
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString(`<w:ftr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">`)
	for _, part := range parts {
		sb.WriteString(`<w:p>`)
		sb.WriteString(`<w:pPr>`)
		switch part.align {
		case "center":
			sb.WriteString(`<w:jc w:val="center"/>`)
		case "right":
			sb.WriteString(`<w:jc w:val="right"/>`)
		}
		if len(part.tabStops) > 0 {
			sb.WriteString(`<w:tabs>`)
			for i, pos := range part.tabStops {
				val := "right"
				if len(part.tabStops) == 2 && i == 0 {
					val = "center"
				}
				sb.WriteString(fmt.Sprintf(`<w:tab w:val="%s" w:pos="%d"/>`, val, pos))
			}
			sb.WriteString(`</w:tabs>`)
		}
		if hasBorder {
			bdrPos := "bottom"
			if cfg.border == "top" {
				bdrPos = "top"
			}
			sb.WriteString(fmt.Sprintf(`<w:pBdr><w:%s w:val="single" w:color="auto" w:space="4"/></w:pBdr>`, bdrPos))
		}
		sb.WriteString(`</w:pPr>`)
		for _, seg := range part.segments {
			if seg.isTab {
				sb.WriteString(`<w:r><w:tab/></w:r>`)
			} else if seg.isField {
				sb.WriteString(fmt.Sprintf(`<w:fldSimple w:instr="%s"/>`, seg.content))
			} else {
				sb.WriteString(`<w:r>`)
				if d.hasRunProps(cfg) {
					sb.WriteString(d.runPropsXML(cfg))
				}
				sb.WriteString(fmt.Sprintf(`<w:t xml:space="preserve">%s</w:t>`, xmlEscape(seg.content)))
				sb.WriteString(`</w:r>`)
			}
		}
		sb.WriteString(`</w:p>`)
	}
	sb.WriteString(`</w:ftr>`)
	return sb.String()
}

func (d *DocxRenderer) hasRunProps(cfg hdrFooterCfg) bool {
	return cfg.fontFamily != "" || cfg.fontSize != "" || cfg.fontColor != ""
}

func (d *DocxRenderer) runPropsXML(cfg hdrFooterCfg) string {
	var sb strings.Builder
	sb.WriteString(`<w:rPr>`)
	if cfg.fontFamily != "" {
		sb.WriteString(fmt.Sprintf(`<w:rFonts w:ascii="%s" w:eastAsia="%s" w:hAnsi="%s"/>`, cfg.fontFamily, cfg.fontFamily, cfg.fontFamily))
	}
	if cfg.fontSize != "" {
		if v, err := strconv.ParseFloat(cfg.fontSize, 64); err == nil {
			sz := int(v * 2)
			sb.WriteString(fmt.Sprintf(`<w:sz w:val="%d"/>`, sz))
		}
	}
	if cfg.fontColor != "" {
		sb.WriteString(fmt.Sprintf(`<w:color w:val="%s"/>`, strings.TrimPrefix(cfg.fontColor, "#")))
	}
	sb.WriteString(`</w:rPr>`)
	return sb.String()
}

func (d *DocxRenderer) Save(path string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	return d.root.SaveTo(path)
}

func intPtr(n int) *int { return &n }
