package database

import (
	"database/sql"
	"github.com/jtieri/HabbGo/habbgo/game/model"
	"log"
)

type RoomRepo struct {
	database *sql.DB
}

func NewRoomRepo(db *sql.DB) *RoomRepo {
	return &RoomRepo{database: db}
}

func (rr *RoomRepo) RoomsByPlayerId(id int) []*model.Room {
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

	var rooms []*model.Room
	for rows.Next() {
		room := new(model.Room)
		room.Details = new(model.Data)

		err := rows.Scan(&room.Details.Id, &room.Details.CatId, &room.Details.Name, &room.Details.Desc, &room.Details.CCTs,
			&room.Details.Wallpaper, &room.Details.Floor, &room.Details.Landscape, &room.Details.Owner_Id, &room.Details.Owner_Name,
			&room.Details.ShowOwner, &room.Details.SudoUsers, &room.Details.AccessType, &room.Details.Password,
			&room.Details.CurrentVisitors, &room.Details.MaxVisitors, &room.Details.Rating)
		if err != nil {
			log.Printf("%v", err)
		}

		rooms = append(rooms, room)
	}

	return rooms
}

func (rr *RoomRepo) fillData(data *model.Data) {

}