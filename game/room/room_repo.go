package room

import (
	"database/sql"
	"log"
)

type RoomRepo struct {
	database *sql.DB
}

func NewRoomRepo(db *sql.DB) *RoomRepo {
	return &RoomRepo{database: db}
}

func (rr *RoomRepo) RoomsByPlayerId(id int) []*Room {
	stmt, err := rr.database.Prepare("SELECT * FROM rooms WHERE owner_id = $1")
	if err != nil {
		log.Printf("%v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		log.Printf("%v", err)
	}
	defer rows.Close()

	var rooms []*Room
	for rows.Next() {
		r := NewRoom()

		var tmpAccessType string
		err := rows.Scan(&r.Details.Id, &r.Details.CategoryID, &r.Details.Name, &r.Details.Description, &r.Details.OwnerId,
			&r.Model.ID, &r.Details.CCTs, &r.Details.Wallpaper, &r.Details.Floor, &r.Details.ShowOwner, &r.Details.Password,
			&tmpAccessType, &r.Details.SudoUsers, &r.Details.CurrentVisitors, &r.Details.MaxVisitors, &r.Details.Rating,
			&r.Details.Hidden, &r.Details.CreatedAt, &r.Details.UpdatedAt)
		if err != nil {
			log.Printf("%v", err)
		}

		r.Details.AccessType = AccessType(tmpAccessType)

		rooms = append(rooms, r)
	}

	return rooms
}

func (rr *RoomRepo) fillData(data *Details) {

}
