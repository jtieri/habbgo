package room

import (
	"database/sql"
	"strings"
	"sync"
)

var (
	rs   *roomService
	once sync.Once
)

const PublicRoomOffset = 1000

type roomService struct {
	repo  *RoomRepo
	rooms map[int]*Room
}

func RoomService() *roomService {
	once.Do(func() {
		rs = &roomService{
			repo:  nil,
			rooms: make(map[int]*Room, 67),
		}
	})
	return rs
}

// SetDBConn sets the database pool handle for the room service
func (rs *roomService) SetDBConn(db *sql.DB) {
	rs.repo = NewRoomRepo(db)
}

// Rooms returns a slice containing all the loaded Rooms currently in the cache
func (rs *roomService) Rooms() []*Room {
	var rooms []*Room
	for _, room := range rs.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

// RoomIsCached returns true if the Room is in the cache of loaded Rooms or false otherwise
func (rs *roomService) RoomIsCached(roomID int) bool {
	_, cached := rs.rooms[roomID]
	return cached
}

// AddRoom adds a Room to the cache of loaded Rooms
func (rs *roomService) AddRoom(room *Room) {
	if room == nil || rs.RoomIsCached(room.Details.ID) {
		return
	}
	rs.rooms[room.Details.ID] = room
}

// RemoveRoom removes a Room from the cache of loaded Rooms
func (rs *roomService) RemoveRoom(roomID int) {
	delete(rs.rooms, roomID)
}

// RoomByID attempts to return a Room with the given room ID if it exists.
func (rs *roomService) RoomByID(roomID int) *Room {
	// if room is cached return it
	if room, ok := rs.rooms[roomID]; ok {
		return room
	}

	// if room is not cached load from database, add to cache and return
	r := rs.repo.RoomByID(roomID)
	if r != nil {
		rs.rooms[r.Details.ID] = r
	}

	return r
}

// RoomsByPlayerID returns a slice of Rooms that are owned by the given Player
func (rs *roomService) RoomsByPlayerID(playerID int) []*Room {
	return rs.repo.RoomsByPlayerId(playerID)
}

func (rs *roomService) RoomByModelName(name string) *Room {
	return &Room{}
}

func (rs *roomService) ReplaceRooms(queryRooms []*Room) []*Room {
	var rooms []*Room

	for _, room := range queryRooms {
		if _, ok := rs.rooms[room.Details.ID]; ok {
			rooms = append(rooms, rs.RoomByID(room.Details.ID))
		} else {
			rooms = append(rooms, room)
		}
	}

	return rooms
}

func (rs *roomService) PublicRoom(room *Room) bool {
	if room.Details.OwnerId == 0 {
		return true
	}
	return false
}

func (rs *roomService) PublicName(room *Room) string {
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

func (rs *roomService) LoadChildRooms(room *Room) {
	if room.Model.Name == "gate_park" {
		room.Details.ChildRooms = append(room.Details.ChildRooms)
	}
}
