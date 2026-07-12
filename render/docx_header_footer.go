package render

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gomutex/godocx/docx"
	"github.com/gomutex/godocx/wml/ctypes"
	"github.com/gomutex/godocx/wml/stypes"
)

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
