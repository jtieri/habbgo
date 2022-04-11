package navigator

import (
	"github.com/jtieri/habbgo/game/ranks"
)

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
	MinRankAccess  ranks.Rank
	MinRankSetFlat int
}
