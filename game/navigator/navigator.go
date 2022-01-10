package navigator

import "github.com/jtieri/habbgo/game/player"

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
	MinRankAccess  player.Rank
	MinRankSetFlat int
}
