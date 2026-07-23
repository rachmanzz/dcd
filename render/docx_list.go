package render

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gomutex/godocx/wml/ctypes"
)

func (d *DocxRenderer) AddList(items []ListItem, ordered bool, numFmt string, start int) error {
	return d.addListAtDepth(items, ordered, numFmt, 0, start)
}

func (d *DocxRenderer) addListAtDepth(items []ListItem, ordered bool, numFmt string, depth int, start int) error {
	if d.root == nil {
		if err := d.init(); err != nil {
			return err
		}
	}
	var style, numID string
	if ordered && numFmt != "" {
		if !d.numFmtInited {
			if err := d.injectNumFmts(); err != nil {
				return err
			}
		}
		d.listCount++
		numID = strconv.Itoa(100 + d.listCount)
		d.injectOLNumEntry(numID, numFmt, start)
		switch depth {
		case 0:
			style = "ListNumber"
		case 1:
			style = "ListNumber2"
		case 2:
			style = "ListNumber3"
		default:
			style = "ListParagraph"
		}
	} else if ordered {
		d.listCount++
		numID = strconv.Itoa(100 + d.listCount)
		d.injectOLNumEntry(numID, "decimal", start)
		switch depth {
		case 0:
			style = "ListNumber"
		case 1:
			style = "ListNumber2"
		case 2:
			style = "ListNumber3"
		default:
			style = "ListParagraph"
		}
	} else {
		switch {
		case depth == 0:
			style, numID = "ListBullet", "1"
		case depth == 1:
			style, numID = "ListBullet2", "2"
		case depth == 2:
			style, numID = "ListBullet3", "3"
		default:
			style, numID = "ListParagraph", "1"
		}
	}
	for _, item := range items {
		if len(item.Runs) > 0 {
			p := d.root.AddEmptyParagraph()
			p.Style(style)
			id, _ := strconv.Atoi(numID)
			p.Numbering(id, 0)
			pPr := p.GetCT().Property
			if pPr == nil {
				pPr = &ctypes.ParagraphProp{}
				p.GetCT().Property = pPr
			}
			d.applyIndent(pPr, item.Attrs, d.defaultStyle)
			for _, run := range item.Runs {
				if run.Tab {
					p.AddText("")
					ctRun := p.GetCT().Children[len(p.GetCT().Children)-1].Run
					ctRun.Children = []ctypes.RunChild{{Tab: &ctypes.Empty{}}}
					continue
				}
				if run.Break {
					p.AddText("")
					ctRun := p.GetCT().Children[len(p.GetCT().Children)-1].Run
					ctRun.Children = []ctypes.RunChild{{Break: &ctypes.Break{}}}
					continue
				}
				r := p.AddText(run.Text)
				ctRun := p.GetCT().Children[len(p.GetCT().Children)-1].Run

				applyRunProps(r, ctRun, run, d.defaultStyle, item.Attrs)
			}
		}
		if len(item.Items) > 0 {
			if err := d.addListAtDepth(item.Items, item.Ordered, item.NumFormat, depth+1, item.Start); err != nil {
				return err
			}
		}
	}
	return nil
}

func numFmtWVal(numFmt string) string {
	switch numFmt {
	case "A":
		return "upperLetter"
	case "a":
		return "lowerLetter"
	case "I":
		return "upperRoman"
	case "i":
		return "lowerRoman"
	default:
		return "decimal"
	}
}

func (d *DocxRenderer) genNsid() string {
	d.nsidCounter++
	return fmt.Sprintf("%08X", d.nsidCounter)
}

func (d *DocxRenderer) injectOLNumEntry(numID string, numFmt string, start int) {
	raw, ok := d.root.FileMap.Load("word/numbering.xml")
	if !ok {
		return
	}
	content := string(raw.([]byte))
	numIDInt, _ := strconv.Atoi(numID)
	abstractNumID := numIDInt
	numFmtVal := numFmtWVal(numFmt)
	nsid := d.genNsid()
	pStyle := "ListNumber"
	switch numFmt {
	case "a", "A":
		pStyle = "ListNumber"
	case "i", "I":
		pStyle = "ListNumber"
	default:
		pStyle = "ListNumber"
	}
	var numEntry string
	if start > 1 {
		numEntry = fmt.Sprintf(`
  <w:abstractNum w:abstractNumId="%d">
    <w:nsid w:val="%s"/>
    <w:multiLevelType w:val="singleLevel"/>
    <w:lvl w:ilvl="0">
      <w:start w:val="1"/>
      <w:numFmt w:val="%s"/>
      <w:pStyle w:val="%s"/>
      <w:lvlText w:val="%%1."/>
      <w:lvlJc w:val="left"/>
      <w:pPr>
        <w:tabs>
          <w:tab w:val="num" w:pos="360"/>
        </w:tabs>
        <w:ind w:left="360" w:hanging="360"/>
      </w:pPr>
    </w:lvl>
  </w:abstractNum>
  <w:num w:numId="%s">
    <w:abstractNumId w:val="%d"/>
    <w:lvlOverride w:ilvl="0">
      <w:startOverride w:val="%d"/>
    </w:lvlOverride>
  </w:num>`, abstractNumID, nsid, numFmtVal, pStyle, numID, abstractNumID, start)
	} else {
		numEntry = fmt.Sprintf(`
  <w:abstractNum w:abstractNumId="%d">
    <w:nsid w:val="%s"/>
    <w:multiLevelType w:val="singleLevel"/>
    <w:lvl w:ilvl="0">
      <w:start w:val="1"/>
      <w:numFmt w:val="%s"/>
      <w:pStyle w:val="%s"/>
      <w:lvlText w:val="%%1."/>
      <w:lvlJc w:val="left"/>
      <w:pPr>
        <w:tabs>
          <w:tab w:val="num" w:pos="360"/>
        </w:tabs>
        <w:ind w:left="360" w:hanging="360"/>
      </w:pPr>
    </w:lvl>
  </w:abstractNum>
  <w:num w:numId="%s">
    <w:abstractNumId w:val="%d"/>
  </w:num>`, abstractNumID, nsid, numFmtVal, pStyle, numID, abstractNumID)
	}
	endTag := "</w:numbering>"
	idx := strings.LastIndex(content, endTag)
	if idx < 0 {
		return
	}
	content = content[:idx] + numEntry + "\n" + content[idx:]
	d.root.FileMap.Store("word/numbering.xml", []byte(content))
}
