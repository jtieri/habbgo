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
	stmt, err := rr.database.Prepare("SELECT r.id, r.cat_id, r.name, r.`desc`, r.ccts, r.wallpaper, r.floor, r.landscape, r.owner_id, r.owner_name, r.show_owner, r.sudo_users, r.access_type, r.password, r.current_visitors, r.max_visitors, r.rating FROM rooms r WHERE r.owner_id = ?")
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
		r := new(Room)
		r.Details = new(Data)

		err := rows.Scan(&r.Details.Id, &r.Details.CatId, &r.Details.Name, &r.Details.Desc, &r.Details.CCTs,
			&r.Details.Wallpaper, &r.Details.Floor, &r.Details.Landscape, &r.Details.Owner_Id, &r.Details.Owner_Name,
			&r.Details.ShowOwner, &r.Details.SudoUsers, &r.Details.AccessType, &r.Details.Password,
			&r.Details.CurrentVisitors, &r.Details.MaxVisitors, &r.Details.Rating)
		if err != nil {
			log.Printf("%v", err)
		}

		rooms = append(rooms, r)
	}

	return rooms
}

func (rr *RoomRepo) fillData(data *Data) {

}
