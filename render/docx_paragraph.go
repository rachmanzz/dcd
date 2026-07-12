package render

import (
	"strconv"
	"strings"

	"github.com/gomutex/godocx/common/units"
	"github.com/gomutex/godocx/docx"
	"github.com/gomutex/godocx/wml/ctypes"
	"github.com/gomutex/godocx/wml/stypes"
)

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

	if s := chooseAttr(def, style, attrs, "line-height"); s != "" {
		if pPr.Spacing == nil {
			pPr.Spacing = &ctypes.Spacing{}
		}
		v := int(atof(s) * 240)
		pPr.Spacing.Line = &v
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

	if s := chooseAttr(def, style, attrs, "keep-next"); s == "true" {
		pPr.KeepNext = ctypes.OnOffFromBool(true)
	}
	if s := chooseAttr(def, style, attrs, "keep-lines"); s == "true" {
		pPr.KeepLines = ctypes.OnOffFromBool(true)
	}

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
		if s := chooseAttr(def, style, attrs, "strike"); s != "" {
			prop.Strike = ctypes.OnOffFromBool(s == "true")
		}
		if s := chooseAttr(def, style, attrs, "underline"); s != "" {
			prop.Underline = ctypes.NewGenSingleStrVal(underlineFromString(s))
		}
		if s := chooseAttr(def, style, attrs, "caps"); s != "" {
			prop.Caps = ctypes.OnOffFromBool(s == "true")
		}
		if s := chooseAttr(def, style, attrs, "small-caps"); s != "" {
			prop.SmallCaps = ctypes.OnOffFromBool(s == "true")
		}
		if s := chooseAttr(def, style, attrs, "letter-spacing"); s != "" {
			prop.Spacing = &ctypes.DecimalNum{Val: int(atof(s) * 20)}
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

func (d *DocxRenderer) applyIndent(pPr *ctypes.ParagraphProp, attrs, defaults map[string]string) {
	indent := ""
	hanging := ""
	if attrs != nil {
		indent = attrs["indent"]
		hanging = attrs["hanging"]
	}
	if indent == "" && defaults != nil {
		indent = defaults["indent"]
	}
	if hanging == "" && defaults != nil {
		hanging = defaults["hanging"]
	}
	if indent == "" && hanging == "" {
		return
	}
	scale := unitScale(d.unit)
	if pPr.Indent == nil {
		pPr.Indent = &ctypes.Indent{}
	}
	if indent != "" {
		v := int(atof(indent) * scale * 56.7)
		pPr.Indent.Left = &v
	}
	if hanging != "" {
		v := uint64(atof(hanging) * scale * 56.7)
		pPr.Indent.Hanging = &v
	}
}

func (d *DocxRenderer) AddParagraph(runs []TextRun, attrs map[string]string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	p := d.root.AddEmptyParagraph()

	pPr := p.GetCT().Property
	if pPr == nil {
		pPr = &ctypes.ParagraphProp{}
		p.GetCT().Property = pPr
	}

	if attrs != nil {
		if s := attrs["align"]; s != "" {
			switch s {
			case "center":
				pPr.Justification = ctypes.NewGenSingleStrVal(stypes.JustificationCenter)
			case "right":
				pPr.Justification = ctypes.NewGenSingleStrVal(stypes.JustificationRight)
			case "justify":
				pPr.Justification = ctypes.NewGenSingleStrVal(stypes.JustificationBoth)
			default:
				pPr.Justification = ctypes.NewGenSingleStrVal(stypes.JustificationLeft)
			}
		}
	}

	d.applyIndent(pPr, attrs, d.defaultStyle)

	for _, r := range runs {
		if r.Tab {
			p.AddText("")
			ctRun := p.GetCT().Children[len(p.GetCT().Children)-1].Run
			ctRun.Children = []ctypes.RunChild{{Tab: &ctypes.Empty{}}}
			continue
		}
		run := p.AddText(r.Text)
		ctRun := p.GetCT().Children[len(p.GetCT().Children)-1].Run

		applyRunProps(run, ctRun, r, d.defaultStyle, attrs)
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

func (d *DocxRenderer) AddSectionBreak(sectionType string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	if d.root.Document.Body.SectPr == nil {
		d.root.Document.Body.SectPr = ctypes.NewSectionProper()
	}
	switch sectionType {
	case "continuous":
		d.root.Document.Body.SectPr.Type = ctypes.NewGenSingleStrVal(stypes.SectionMarkNextContinuous)
	case "next-page":
		d.root.Document.Body.SectPr.Type = ctypes.NewGenSingleStrVal(stypes.SectionMarkNextPage)
	case "even-page":
		d.root.Document.Body.SectPr.Type = ctypes.NewGenSingleStrVal(stypes.SectionMarkEvenPage)
	case "odd-page":
		d.root.Document.Body.SectPr.Type = ctypes.NewGenSingleStrVal(stypes.SectionMarkOddPage)
	}
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
		switch attrs["align"] {
		case "center":
			pic.Para.Justification(stypes.JustificationCenter)
		case "right":
			pic.Para.Justification(stypes.JustificationRight)
		}

		if alt := attrs["alt"]; alt != "" {
			pic.Inline.DocProp.Description = alt
		}

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
		c := attrs["font-color"]
		if c == "" {
			c = attrs["color"]
		}
		if c != "" {
			linkColor = strings.TrimPrefix(c, "#")
		}
		if attrs["underline"] == "false" {
			hasUnderline = false
		}
		_ = attrs["target"]
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

func (d *DocxRenderer) AddWrappedParagraph(text string, flags string, attrs map[string]string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	p := d.root.AddParagraph(text)

	if attrs != nil {
		if len(p.GetCT().Children) > 0 {
			if run := p.GetCT().Children[0].Run; run != nil {
				if run.Property == nil {
					run.Property = &ctypes.RunProperty{}
				}
				if fc := attrs["font-color"]; fc != "" {
					run.Property.Color = ctypes.NewColor(fc)
				}
				if fs := attrs["font-size"]; fs != "" {
					run.Property.Size = ctypes.NewFontSize(uint64(atof(fs)) * 2)
				}
			}
		}
	}

	for _, f := range strings.Split(flags, "|") {
		switch f {
		case "c":
			p.Justification(stypes.JustificationCenter)
		case "r":
			p.Justification(stypes.JustificationRight)
		case "j":
			p.Justification(stypes.JustificationBoth)
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
		case "s":
			if len(p.GetCT().Children) > 0 {
				if run := p.GetCT().Children[0].Run; run != nil {
					if run.Property == nil {
						run.Property = &ctypes.RunProperty{}
					}
					run.Property.Strike = ctypes.OnOffFromBool(true)
				}
			}
		case "u":
			if len(p.GetCT().Children) > 0 {
				if run := p.GetCT().Children[0].Run; run != nil {
					if run.Property == nil {
						run.Property = &ctypes.RunProperty{}
					}
					uType := "single"
					if attrs != nil && attrs["underline"] != "" {
						uType = attrs["underline"]
					}
					run.Property.Underline = ctypes.NewGenSingleStrVal(underlineFromString(uType))
				}
			}
		case "code":
			if len(p.GetCT().Children) > 0 {
				if run := p.GetCT().Children[0].Run; run != nil {
					if run.Property == nil {
						run.Property = &ctypes.RunProperty{}
					}
					run.Property.Fonts = &ctypes.RunFonts{
						Ascii: "Courier New",
						HAnsi: "Courier New",
					}
				}
			}
		case "mark":
			if len(p.GetCT().Children) > 0 {
				if run := p.GetCT().Children[0].Run; run != nil {
					if run.Property == nil {
						run.Property = &ctypes.RunProperty{}
					}
					color := "yellow"
					if attrs != nil && attrs["color"] != "" {
						color = attrs["color"]
					}
					run.Property.Highlight = ctypes.NewCTString(color)
				}
			}
		case "sub":
			if len(p.GetCT().Children) > 0 {
				if run := p.GetCT().Children[0].Run; run != nil {
					if run.Property == nil {
						run.Property = &ctypes.RunProperty{}
					}
					run.Property.VertAlign = ctypes.NewGenSingleStrVal(stypes.VerticalAlignRunSubscript)
				}
			}
		case "sup":
			if len(p.GetCT().Children) > 0 {
				if run := p.GetCT().Children[0].Run; run != nil {
					if run.Property == nil {
						run.Property = &ctypes.RunProperty{}
					}
					run.Property.VertAlign = ctypes.NewGenSingleStrVal(stypes.VerticalAlignRunSuperscript)
				}
			}
		}
	}
	return nil
}
