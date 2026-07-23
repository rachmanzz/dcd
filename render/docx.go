package render

import (
	"fmt"
	"strings"

	"github.com/gomutex/godocx"
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
	numFmtInited  bool
	listCount     int
	nsidCounter   int
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

	appXML := `<?xml version="1.0" encoding="UTF-8"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes">`
	if title != "" {
		appXML += fmt.Sprintf("<Title>%s</Title>", title)
	}
	appXML += `</Properties>`

	d.root.FileMap.Store("docProps/app.xml", []byte(appXML))

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
	w, h := parsePageSize(props["layout"], props["orientation"], props["w"], props["h"], props["unit"])
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

func (d *DocxRenderer) Save(path string) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	return d.root.SaveTo(path)
}

func intPtr(n int) *int { return &n }

func (d *DocxRenderer) injectNumFmts() error {
	if d.numFmtInited {
		return nil
	}
	raw, ok := d.root.FileMap.Load("word/numbering.xml")
	if !ok {
		return nil
	}
	content := string(raw.([]byte))
	endTag := "</w:numbering>"
	idx := strings.LastIndex(content, endTag)
	if idx < 0 {
		return nil
	}

	indents := []string{"360", "720", "1080"}
	formats := []struct {
		name string
		base int
	}{
		{"upperLetter", 10},
		{"lowerLetter", 13},
		{"upperRoman", 16},
		{"lowerRoman", 19},
	}
	var sb strings.Builder
	for _, f := range formats {
		for lvl := 0; lvl < 3; lvl++ {
			absID := f.base + lvl
			numID := f.base + lvl
			nsid := d.genNsid()
			sb.WriteString(fmt.Sprintf(`
  <w:abstractNum w:abstractNumId="%d">
    <w:nsid w:val="%s"/>
    <w:multiLevelType w:val="singleLevel"/>
    <w:lvl w:ilvl="0">
      <w:start w:val="1"/>
      <w:numFmt w:val="%s"/>
      <w:pStyle w:val="ListNumber"/>
      <w:lvlText w:val="%%1."/>
      <w:lvlJc w:val="left"/>
      <w:pPr>
        <w:tabs>
          <w:tab w:val="num" w:pos="%s"/>
        </w:tabs>
        <w:ind w:left="%s" w:hanging="360"/>
      </w:pPr>
    </w:lvl>
  </w:abstractNum>
  <w:num w:numId="%d">
    <w:abstractNumId w:val="%d"/>
  </w:num>`, absID, nsid, f.name, indents[lvl], indents[lvl], numID, absID))
		}
	}
	content = content[:idx] + sb.String() + "\n" + content[idx:]
	d.root.FileMap.Store("word/numbering.xml", []byte(content))
	d.numFmtInited = true
	return nil
}
