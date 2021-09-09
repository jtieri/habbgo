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
	rows, err := navRepo.database.Query("SELECT * FROM room_categories")
	if err != nil {
		log.Printf("%v", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		err = rows.Scan(&cat.Id, &cat.Pid, &cat.Node, &cat.Name, &cat.Public, &cat.Trading, &cat.MinRankAccess)
		if err != nil {
			log.Printf("%v", err)
		}

		categories = append(categories, cat)
	}

	return categories
}
