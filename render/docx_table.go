package render

import (
	"github.com/gomutex/godocx/wml/ctypes"
	"github.com/gomutex/godocx/wml/stypes"
)

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

			for _, run := range cell.Runs {
				r := p.AddText(run.Text)
				ctRun := p.GetCT().Children[len(p.GetCT().Children)-1].Run

				applyRunProps(r, ctRun, run, d.defaultStyle, nil)

				if rowStyle != nil {
					if fc := rowStyle["font-color"]; fc != "" {
						r.Color(fc)
					}
					if rowStyle["font-weight"] == "bold" {
						r.Bold(true)
					}
				}
				if cell.Attrs != nil {
					if fc := cell.Attrs["font-color"]; fc != "" {
						r.Color(fc)
					}
					if fs := cell.Attrs["font-size"]; fs != "" {
						r.Size(uint64(atof(fs)))
					}
				}
			}

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
