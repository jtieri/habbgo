package service

import (
	"database/sql"
	"github.com/jtieri/HabbGo/habbgo/database"
	"github.com/jtieri/HabbGo/habbgo/game/model"
	"strings"
	"sync"
)

var rs *roomService
var ronce sync.Once
const PublicRoomOffset = 1000

type roomService struct {
	repo  *database.RoomRepo
	rooms map[int]*model.Room
}

func RoomService() *roomService {
	ronce.Do(func() {
		rs = &roomService{
			repo:  nil,
			rooms: make(map[int]*model.Room, 50),
		}
	})

	return rs
}

func (rs *roomService) SetDBConn(db *sql.DB) {
	rs.repo = database.NewRoomRepo(db)
}

func (rs *roomService) Rooms() []*model.Room {
	var rooms []*model.Room
	for _, room := range rs.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

func (rs *roomService) RoomById(id int) *model.Room {
	if room, ok := rs.rooms[id]; ok {
		return room
	}
	return nil
}

func (rs *roomService) RoomsByPlayerId(id int) []*model.Room {
	return rs.repo.RoomsByPlayerId(id)
}

func (rs *roomService) RoomByModelName(name string) *model.Room {
	return &model.Room{}
}

func (rs *roomService) ReplaceRooms(queryRooms []*model.Room) []*model.Room {
	var rooms []*model.Room

	for _, room := range queryRooms {
		if _, ok := rs.rooms[room.Details.Id]; ok {
			rooms = append(rooms, rs.RoomById(room.Details.Id))
		} else {
			rooms = append(rooms, room)
		}
	}

	return rooms
}

func AccessType(accessId int) string {
	switch accessId {
	case 1:
		return "closed"
	case 2:
		return "password"
	default:
		return "open"
	}
}

func (rs *roomService) PublicRoom(room *model.Room) bool {
	if room.Details.Owner_Id == 0 {
		return true
	} else {
		return false
	}
}

func (rs *roomService) PublicName(room *model.Room) string {
	 if rs.PublicRoom(room) {
	 	if strings.HasPrefix(room.Details.Name, "Upper Hallways") {
			return "Upper Hallways"
		}

		if strings.HasPrefix(room.Details.Name, "Lower Hallways") {
			 return "Lower Hallways"
		}

		if strings.HasPrefix(room.Details.Name, "Club Massiva") {
			return "Club Massiva"
		}

		if strings.HasPrefix(room.Details.Name, "The Chromide Club") {
			 return "The Chromide Club"
		}

		if room.Details.CCTs == "hh_room_gamehall,hh_games" {
			return "Cunning Fox Gamehall"
		}
	 }

	 return room.Details.Name
}

func (rs *roomService) CurrentVisitors() int {
	var visitors int

	return visitors
}

func (rs *roomService) MaxVisitors() int {
	var visitors int

	return visitors
}

func (rs *roomService) LoadChildRooms(room *model.Room) {
	if room.Model.Name == "gate_park" {
		room.Details.ChildRooms = append(room.Details.ChildRooms,  )
	}
}