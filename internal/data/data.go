package data

type TypeProblem int

const (
	Single TypeProblem = iota
	Multi
)

type Item struct {
	title, desc string
}

func NewItem(title, desc string) Item {
	return Item{title: title, desc: desc}
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }

func (i Item) FilterValue() string { return i.title }
