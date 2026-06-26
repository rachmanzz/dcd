package render

type TextRun struct {
	Text          string
	Bold          bool
	Italic        bool
	Underline     bool
	Code          bool
	Tab           bool
	Link          string
	LinkAttrs     map[string]string
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
	Runs  []TextRun
	Items []ListItem
}

type Renderer interface {
	AddHeading(text string, level int, attrs map[string]string) error
	AddParagraph(runs []TextRun) error
	AddLineBreak() error
	AddHorizontalRule(attrs map[string]string) error
	AddPageBreak() error
	AddImage(src string, attrs map[string]string) error
	AddHyperlink(text, url string, attrs map[string]string) error
	AddWrappedParagraph(text string, flags string) error
	AddList(items []ListItem, ordered bool) error
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
