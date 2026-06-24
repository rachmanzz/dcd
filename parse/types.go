package parse

type Doc struct {
	Sections []Section
}

type Section struct {
	N     int
	Name  string
	Props map[string]string
	Body  string
}
