package render

type TextRun struct {
	Text           string
	Bold           bool
	Italic         bool
	Underline      bool
	UnderlineStyle string
	Code           bool
	Strike         bool
	Mark           bool
	MarkColor      string
	Sub            bool
	Sup            bool
	Tab            bool
	Link           string
	LinkAttrs      map[string]string
}

type TableCell struct {
	Runs  []TextRun
	Attrs map[string]string
}

type TableRow struct {
	Cells []TableCell
	Props map[string]string
}

type ListItem struct {
	Runs      []TextRun
	Items     []ListItem
	Attrs     map[string]string
	Ordered   bool   // sub-list type — true for <ol>, false for <ul>
	NumFormat string // "" = decimal, "A" = upperLetter, "a" = lowerLetter, "I" = upperRoman, "i" = lowerRoman
}

type Renderer interface {
	AddHeading(text string, level int, attrs map[string]string) error
	AddParagraph(runs []TextRun, attrs map[string]string) error
	AddLineBreak() error
	AddHorizontalRule(attrs map[string]string) error
	AddPageBreak() error
	AddSectionBreak(sectionType string) error
	AddImage(src string, attrs map[string]string) error
	AddHyperlink(text, url string, attrs map[string]string) error
	AddWrappedParagraph(text string, flags string, attrs map[string]string) error
	AddList(items []ListItem, ordered bool, numFmt string) error
	AddTable(rows []TableRow, attrs map[string]string) error
	SetPageStyle(props map[string]string) error
	SetHeader(props map[string]string) error
	SetFooter(props map[string]string) error
	SetHeadingStyle(level int, props map[string]string) error
	SetTableStyle(name string, props map[string]string) error
	SetDefaultStyle(props map[string]string) error
	SetMetadata(props map[string]string) error
	Save(path string) error
}
