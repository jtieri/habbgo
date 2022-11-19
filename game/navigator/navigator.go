package navigator

import (
	"github.com/jtieri/habbgo/collections"
	"github.com/jtieri/habbgo/game/ranks"
)

type Navigator struct {
	categoryCache collections.Cache[int, Category]
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

func newNavigator() Navigator {
	return Navigator{
		categoryCache: collections.NewCache(make(map[int]Category)),
	}
}

func (n *Navigator) categoryByID(id int) (Category, bool) {
	return n.categoryCache.Get(id)
}

func (n *Navigator) categoryByParentID(pid int) []Category {
	cats := make([]Category, 0)

	for _, cat := range n.categoryCache.Items() {
		if cat.ParentID == pid {
			cats = append(cats, cat)
		}
	}
	return cats
}

// categories returns a slice of Category's that is a clone of the Navigator's map of Category's.
func (n *Navigator) categories() []Category {
	return n.categoryCache.Items()
}
