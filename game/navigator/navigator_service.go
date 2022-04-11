package navigator

import (
	"database/sql"

	"github.com/jtieri/habbgo/game/room"
	"go.uber.org/zap"
)

type NavService struct {
	repo *NavRepo
	nav  *Navigator
	log  *zap.Logger
}

func NewNavigatorService(log *zap.Logger, db *sql.DB) *NavService {
	return &NavService{
		repo: NewNavRepo(db),
		nav:  new(Navigator),
		log:  log,
	}
}

// SetDBConnection is called when a NavService struct is allocated initially so that it has access to the applications db.
func (ns *NavService) SetDBConnection(db *sql.DB) {
	ns.repo = NewNavRepo(db)
}

// Build retrieves the room categories from the database and builds the in-game Navigator with them.
func (ns *NavService) Build() {
	ns.nav.Categories = ns.repo.Categories()
}

// CategoryById retrieves a navigator category given the int parameter id and returns it if there is a match.
func (ns *NavService) CategoryById(id int) *Category {
	for _, cat := range ns.nav.Categories {
		if cat.ID == id {
			return &cat
		}
	}

	return nil
}

// CategoriesByParentId retrieves a slice of sub-categories given the int parameter pid and returns it if there is a match.
func (ns *NavService) CategoriesByParentId(pid int) []Category {
	var categories []Category

	for _, cat := range ns.nav.Categories {
		if cat.ParentID == pid {
			categories = append(categories, cat)
		}
	}

	return categories
}

func CurrentVisitors(cat *Category, rooms []*room.Room) int {
	visitors := 0

	for _, r := range rooms {
		if r.Details.CategoryID == cat.ID {
			visitors += r.Details.CurrentVisitors
		}
	}

	return visitors
}

func MaxVisitors(cat *Category, rooms []*room.Room) int {
	visitors := 0

	for _, r := range rooms {
		if r.Details.CategoryID == cat.ID {
			visitors += r.Details.MaxVisitors
		}
	}

	return visitors
}
