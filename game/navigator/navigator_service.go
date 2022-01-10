package navigator

import (
	"database/sql"
	"github.com/jtieri/habbgo/game/room"
	"sync"
)

var ns *navService
var nonce sync.Once

type navService struct {
	repo *NavRepo
	nav  *Navigator
	mux  *sync.Mutex
}

// NavigatorService will initialize the single instance of navService if it is the first time it is called and then it
// will return the instance.
func NavigatorService() *navService {
	nonce.Do(func() {
		ns = &navService{
			repo: nil,
			nav:  new(Navigator),
			mux:  &sync.Mutex{},
		}
	})

	return ns
}

// SetDBCon is called when a NavService struct is allocated initially so that it has access to the applications db.
func (ns *navService) SetDBCon(db *sql.DB) {
	ns.repo = NewNavRepo(db)
}

// BuildNavigator retrieves the room categories from the database and builds the in-game Navigator with them.
func (ns *navService) BuildNavigator() {
	ns.nav.Categories = ns.repo.Categories()
}

// CategoryById retrieves a navigator category given the int parameter id and returns it if there is a match.
func (ns *navService) CategoryById(id int) *Category {
	for _, cat := range ns.nav.Categories {
		if cat.ID == id {
			return &cat
		}
	}

	return nil
}

// CategoriesByParentId retrieves a slice of sub-categories given the int parameter pid and returns it if there is a match.
func (ns *navService) CategoriesByParentId(pid int) []*Category {
	var categories []*Category

	for _, cat := range ns.nav.Categories {
		if cat.ParentID == pid {
			categories = append(categories, &cat)
		}
	}

	return categories
}

func CurrentVisitors(cat *Category) int {
	visitors := 0

	for _, r := range room.RoomService().Rooms() {
		if r.Details.CatId == cat.ID {
			visitors += r.Details.CurrentVisitors
		}
	}

	return visitors
}

func MaxVisitors(cat *Category) int {
	visitors := 0

	for _, r := range room.RoomService().Rooms() {
		if r.Details.CatId == cat.ID {
			visitors += r.Details.MaxVisitors
		}
	}

	return visitors
}
