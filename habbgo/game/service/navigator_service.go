package service

import (
	"database/sql"
	"github.com/jtieri/HabbGo/habbgo/database"
	"github.com/jtieri/HabbGo/habbgo/game/model"
)

type NavService struct {
	np  *database.NavRepo
	nav *model.Navigator
}

func NewNavService(db *sql.DB) *NavService {
	return &NavService{
		np:  database.NewNavRepo(db),
		nav: &model.Navigator{Categories: nil},
	}
}

func (ns NavService) BuildNavigator() {
	ns.nav.Categories = ns.np.Categories()
}

func (ns NavService) CategoryById(id int) *model.Category {
	for _, cat := range ns.nav.Categories {
		if cat.Id == id {
			return &cat
		}
	}

	return nil
}
