package database

import (
	"gorm.io/gorm"
)

type NavRepo struct {
	database gorm.DB
}

// NewNavRepo returns a new instance of NavRepo for use in the navigator service.
func NewNavRepo(db gorm.DB) *NavRepo {
	return &NavRepo{database: db}
}

// Categories retrieves the navigator categories found in database table room_categories and returns them as a slice of
// Category structs.
//func (navRepo *NavRepo) Categories() []model.Category {
//	rows, err := navRepo.database.Query("SELECT * FROM room_categories")
//	if err != nil {
//		log.Printf("%v", err)
//	}
//	defer rows.Close()
//
//	var categories []model.Category
//	for rows.Next() {
//		var cat model.Category
//		err = rows.Scan(&cat.Id, &cat.Pid, &cat.Node, &cat.Name, &cat.Public, &cat.Trading, &cat.MinRankAccess)
//		if err != nil {
//			log.Printf("%v", err)
//		}
//
//		categories = append(categories, cat)
//	}
//
//	return categories
//}
