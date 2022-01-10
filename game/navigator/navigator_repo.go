package navigator

import (
	"database/sql"
	"log"
)

type NavRepo struct {
	database *sql.DB
}

// NewNavRepo returns a new instance of NavRepo for use in the navigator service.
func NewNavRepo(db *sql.DB) *NavRepo {
	return &NavRepo{database: db}
}

// Categories retrieves the navigator categories found in database table room_categories and returns them as a slice of
// Category structs.
func (navRepo *NavRepo) Categories() []Category {
	rows, err := navRepo.database.Query("SELECT id, parent_id, is_node, name, is_public, is_trading, min_rank_access, min_rank_setflatcat FROM room_categories")
	if err != nil {
		log.Printf("%v", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		err = rows.Scan(&cat.ID, &cat.ParentID, &cat.IsNode, &cat.Name, &cat.IsPublic, &cat.IsTrading, &cat.MinRankAccess, &cat.MinRankSetFlat)
		if err != nil {
			log.Printf("%v", err)
		}

		categories = append(categories, cat)
	}

	return categories
}
