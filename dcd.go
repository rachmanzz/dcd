package dcd

import (
	"github.com/rachmanzz/dcd/data"
	"github.com/rachmanzz/dcd/parse"
	"github.com/rachmanzz/dcd/render"
)

type Doc = parse.Doc
type DataSet = data.DataSet
type DocxRenderer = render.DocxRenderer
type Compiler = render.Compiler
type TextRun = render.TextRun
type ListItem = render.ListItem
type TableCell = render.TableCell
type TableRow = render.TableRow

var (
	NewDataSet      = data.NewDataSet
	Parse           = parse.Parse
	NewDocxRenderer = render.NewDocxRenderer
	New             = render.New
)
