package service

import (
	"database/sql"
	"github.com/jtieri/HabbGo/habbgo/database"
	"github.com/jtieri/HabbGo/habbgo/game/model"
	"sync"
)

var rs *roomService
var ronce sync.Once

type roomService struct {
	repo  *database.RoomRepo
	rooms []*model.Room
}

func RoomService() *roomService {
	ronce.Do(func() {
		rs = &roomService{
			repo:  nil,
			rooms: make([]*model.Room, 50),
		}
	})

	return rs
}

func (rs *roomService) SetDBConn(db *sql.DB) {
	rs.repo = database.NewRoomRepo(db)
}
