package room

import (
	"database/sql"
	"strings"

	"go.uber.org/zap"
)

const PublicRoomOffset = 1000

type RoomService struct {
	repo  *RoomRepo
	rooms map[int]*Room

	log *zap.Logger
}

func NewRoomService(log *zap.Logger, db *sql.DB) *RoomService {
	return &RoomService{
		repo:  NewRoomRepo(db),
		rooms: nil,
		log:   log,
	}
}

func (rs *RoomService) Build() {

}

func (rs *RoomService) Rooms() []*Room {
	var rooms []*Room
	for _, room := range rs.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

func (rs *RoomService) RoomById(id int) *Room {
	if room, ok := rs.rooms[id]; ok {
		return room
	}
	return nil
}

func (rs *RoomService) RoomsByPlayerId(id int) []*Room {
	return rs.repo.RoomsByPlayerId(id)
}

func (rs *RoomService) RoomByModelName(name string) *Room {
	return &Room{}
}

func (rs *RoomService) ReplaceRooms(queryRooms []*Room) []*Room {
	var rooms []*Room

	for _, room := range queryRooms {
		if _, ok := rs.rooms[room.Details.Id]; ok {
			rooms = append(rooms, rs.RoomById(room.Details.Id))
		} else {
			rooms = append(rooms, room)
		}
	}

	return rooms
}

func (rs *RoomService) PublicRoom(room *Room) bool {
	if room.Details.OwnerId == 0 {
		return true
	}
	return false
}

func (rs *RoomService) PublicName(room *Room) string {
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

func (rs *RoomService) CurrentVisitors() int {
	var visitors int

	return visitors
}

func (rs *RoomService) MaxVisitors() int {
	var visitors int

	return visitors
}

func (rs *RoomService) LoadChildRooms(room *Room) {
	if room.Model.Name == "gate_park" {
		room.Details.ChildRooms = append(room.Details.ChildRooms)
	}
}
