package navigator

type Navigator struct {
	Categories []Category
}

type Category struct {
	ID             int
	ParentID       int
	Name           string
	IsNode         bool
	IsPublic       bool
	IsTrading      bool
	MinRankAccess  int
	MinRankSetFlat int
}
