package database

import "database/sql"

type RoomRepo struct {
	database *sql.DB
}

func NewRoomRepo(db *sql.DB) *RoomRepo {
	return &RoomRepo{database: db}
}
