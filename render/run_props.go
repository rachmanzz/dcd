package render

import (
	"strings"

	"github.com/gomutex/godocx/docx"
	"github.com/gomutex/godocx/wml/ctypes"
	"github.com/gomutex/godocx/wml/stypes"
)

func applyRunProps(r *docx.Run, ctRun *ctypes.Run, run TextRun, defaultStyle, attrs map[string]string) {
	if defaultStyle != nil {
		if fc := defaultStyle["font-color"]; fc != "" {
			r.Color(fc)
		}
		if fs := defaultStyle["font-size"]; fs != "" {
			r.Size(uint64(atof(fs)))
		}
		if ff := defaultStyle["font-family"]; ff != "" {
			if ctRun.Property == nil {
				ctRun.Property = &ctypes.RunProperty{}
			}
			ctRun.Property.Fonts = &ctypes.RunFonts{
				Ascii: ff,
				HAnsi: ff,
			}
		}
	}
	if attrs != nil {
		if fc := attrs["font-color"]; fc != "" {
			r.Color(fc)
		}
		if fs := attrs["font-size"]; fs != "" {
			r.Size(uint64(atof(fs)))
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
	if run.Strike {
		r.Strike(true)
	}
	if run.Bold {
		r.Bold(true)
	}
	if run.Italic {
		r.Italic(true)
	}
	isLink := run.Link != ""
	if isLink && run.LinkAttrs["underline"] == "false" {
		// skip underline
	} else if isLink || run.Underline {
		r.Underline(underlineFromString(run.UnderlineStyle))
	}
	if run.Mark {
		color := "yellow"
		if run.MarkColor != "" {
			color = run.MarkColor
		}
		r.Highlight(color)
	}
	if run.Sub {
		if ctRun.Property == nil {
			ctRun.Property = &ctypes.RunProperty{}
		}
		ctRun.Property.VertAlign = ctypes.NewGenSingleStrVal(stypes.VerticalAlignRunSubscript)
	}
	if run.Sup {
		if ctRun.Property == nil {
			ctRun.Property = &ctypes.RunProperty{}
		}
		ctRun.Property.VertAlign = ctypes.NewGenSingleStrVal(stypes.VerticalAlignRunSuperscript)
	}
	if isLink {
		linkColor := "0055CC"
		if run.LinkAttrs != nil && run.LinkAttrs["color"] != "" {
			linkColor = strings.TrimPrefix(run.LinkAttrs["color"], "#")
		}
		r.Color(linkColor)
	}
}
