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

func (rr *RoomRepo) RoomsByPlayerId(playerID int) []*Room {
	stmt, err := rr.database.Prepare("SELECT * FROM rooms WHERE owner_id = $1")
	if err != nil {
		log.Printf("%v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(playerID)
	if err != nil {
		log.Printf("%v", err)
	}
	defer rows.Close()

	var rooms []*Room
	for rows.Next() {
		r := NewRoom()

		var tmpAccessType string
		err := rows.Scan(&r.Details.ID, &r.Details.CategoryID, &r.Details.Name, &r.Details.Description, &r.Details.OwnerId,
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

func (rr *RoomRepo) RoomByID(roomID int) *Room {
	stmt, err := rr.database.Prepare("SELECT * FROM rooms WHERE id = $1")
	if err != nil {
		log.Printf("%v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(roomID)
	if err != nil {
		log.Printf("%v", err)
	}
	defer rows.Close()

	var room *Room
	if rows.Next() {
		r := NewRoom()

		err = fillData(rows, r)
		if err != nil {
			log.Printf("%v", err)
		}
	}

	return room
}

func fillData(rows *sql.Rows, room *Room) error {
	var tmpAccessType string
	err := rows.Scan(&room.Details.ID, &room.Details.CategoryID, &room.Details.Name, &room.Details.Description, &room.Details.OwnerId,
		&room.Model.ID, &room.Details.CCTs, &room.Details.Wallpaper, &room.Details.Floor, &room.Details.ShowOwner, &room.Details.Password,
		&tmpAccessType, &room.Details.SudoUsers, &room.Details.CurrentVisitors, &room.Details.MaxVisitors, &room.Details.Rating,
		&room.Details.Hidden, &room.Details.CreatedAt, &room.Details.UpdatedAt)
	if err != nil {
		return err
	}
	room.Details.AccessType = AccessType(tmpAccessType)
	return nil
}
