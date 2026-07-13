package render

import (
	"strconv"

	"github.com/gomutex/godocx/wml/ctypes"
)

func (d *DocxRenderer) AddList(items []ListItem, ordered bool, numFmt string) error {
	return d.addListAtDepth(items, ordered, numFmt, 0)
}

func (d *DocxRenderer) addListAtDepth(items []ListItem, ordered bool, numFmt string, depth int) error {
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
		baseID := numFmtBaseID(numFmt)
		if baseID > 0 {
			numID = strconv.Itoa(baseID + depth)
		}
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
		case ordered && depth == 0:
			style, numID = "ListNumber", "5"
		case ordered && depth == 1:
			style, numID = "ListNumber2", "6"
		case ordered && depth == 2:
			style, numID = "ListNumber3", "7"
		case !ordered && depth == 0:
			style, numID = "ListBullet", "1"
		case !ordered && depth == 1:
			style, numID = "ListBullet2", "2"
		case !ordered && depth == 2:
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
			if err := d.addListAtDepth(item.Items, item.Ordered, item.NumFormat, depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}
