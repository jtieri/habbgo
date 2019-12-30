package database

import (
	"database/sql"
	"github.com/jtieri/HabbGo/habbgo/game/model"
	"log"
)

type NavRepo struct {
	database *sql.DB
}

func NewNavRepo(db *sql.DB) *NavRepo {
	return &NavRepo{database: db}
}

func (navRepo *NavRepo) Categories() []model.Category {
	rows, err := navRepo.database.Query("SELECT * FROM nav_categories")
	if err != nil {
		log.Printf("%v", err)
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var cat model.Category
		err = rows.Scan(&cat.Id, &cat.Pid, &cat.Node, &cat.Name, &cat.Public, &cat.Trading, &cat.MinRank)
		if err != nil {
			log.Printf("%v", err)
		}

		categories = append(categories, cat)
	}

	return categories
}
