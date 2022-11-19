package room

import (
	"fmt"

	"github.com/jtieri/habbgo/game/item"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/scheduler"
	"github.com/jtieri/habbgo/game/service/query"
	"github.com/jtieri/habbgo/protocol/packets"
)

// RoomServiceProxy provides a high level API for communicating with a running instance of room.RoomService.
type RoomServiceProxy struct {
	service *ServiceChannels
}

func NewProxy(channels *ServiceChannels) *RoomServiceProxy {
	return &RoomServiceProxy{
		service: channels,
	}
}

func (p *RoomServiceProxy) Init() {

}

// AddRoom calls the underlying function on the RoomService by first sending a new Room over the channel,
// and then calling the function to trigger a read on the RoomService which reads the room and adds it to the cache.
func (p *RoomServiceProxy) AddRoom(room Room) {
	req := query.NewRequest(room.Details.Id, room)
	p.service.RoomAddChan <- req

	<-req.Response // wait for signal from service that the request has been processed.
}

// RemoveRoom ...
func (p *RoomServiceProxy) RemoveRoom(roomID int) {
	req := query.NewRequest(roomID, Room{})
	p.service.RoomDeleteChan <- req

	<-req.Response // wait for signal from service that the request has been processed.
}

// RoomsByPlayerID ...
func (p *RoomServiceProxy) RoomsByPlayerID(pid int) []Room {
	req := query.NewRequest(pid, []Room{})
	p.service.RoomByPlayerID <- req

	resp := <-req.Response // block until we get a response from the service.
	if resp == nil {
		return nil
	}

	return resp.Value
}

// RoomByID ...
func (p *RoomServiceProxy) RoomByID(roomID int) Room {
	req := query.NewRequest(roomID, Room{})
	p.service.RoomByRoomID <- req

	fmt.Println("Sent req")
	resp := <-req.Response
	fmt.Println("Read response")
	if resp == nil {
		return Room{}
	}

	fmt.Println("Returning from proxy call")
	return resp.Value
}

// Rooms calls the underlying function on the RoomService which will write a slice containing
// a copy of the currently cached rooms. This call reads that slice from the channel and returns it.
func (p *RoomServiceProxy) Rooms() []Room {
	req := query.NewRequest(0, []Room{})
	p.service.RoomsChan <- req

	resp := <-req.Response
	return resp.Value
}

// RoomCached calls the underlying function on the RoomService by first sending a query
// over the proper channel to the service. The service will then read the query and send a response
// over the channel.
func (p *RoomServiceProxy) RoomCached(roomID int) bool {
	req := query.NewRequest(roomID, false)
	p.service.RoomCachedChan <- req

	resp := <-req.Response
	return resp.Value
}

// RoomsCachedCount retrieves the number of rooms currently cached in the room service.
func (p *RoomServiceProxy) RoomsCachedCount() int {
	req := query.NewRequest(0, 0)
	p.service.RoomCountChan <- req

	resp := <-req.Response
	return resp.Value
}

// CheckRoomsQueried sends a slice of rooms to the room service and
// receives a new slice of rooms in response where any room from the original slice
// that is already initialized will be replaced with a copy of the initialized room.
// This is used after reading rooms from the service that were loaded from the database
// so that we are always using a rooms state from the servers in memory store vs. possibly
// stale data from the database.
func (p *RoomServiceProxy) CheckRoomsQueried(queryRooms []Room) []Room {
	req := query.NewRequest(0, queryRooms)
	p.service.RoomQueriedChan <- req

	resp := <-req.Response
	return resp.Value
}

// Players ...
func (p *RoomServiceProxy) Players(roomID int) []player.Player {
	req := query.NewRequest(roomID, []player.Player{})
	p.service.RoomPlayersChan <- req

	resp := <-req.Response
	if resp == nil {
		return nil
	}

	return resp.Value
}

// RoomAddPlayer Add a player to a room
func (p *RoomServiceProxy) RoomAddPlayer(roomID int, plyr player.Player) {
	req := query.NewRequest(roomID, plyr)
	p.service.RoomAddPlayerChan <- req

	<-req.Response // wait for signal from the service that req has been processed.
}

func (p *RoomServiceProxy) RoomRemovePlayer(roomID, playerID int) {
	req := query.NewRequest(roomID, playerID)
	p.service.RoomRemovePlayerChan <- req

	<-req.Response // wait for signal from the service that req has been processed.
}

// RoomPlayerCount Get number of players in a room
func (p *RoomServiceProxy) RoomPlayerCount(roomID int) int {
	req := query.NewRequest(roomID, 0)
	p.service.RoomPlayerCountChan <- req

	resp := <-req.Response
	if resp == nil {
		return 0
	}

	return resp.Value
}

// RoomAddItems Add items to a room
func (p *RoomServiceProxy) RoomAddItems(roomID int, items []item.Item) {
	req := query.NewRequest(roomID, items)
	p.service.RoomAddItemsChan <- req

	<-req.Response // wait for signal from the service that req has been processed.
}

// RoomFloorItems Get floor items in a room
func (p *RoomServiceProxy) RoomFloorItems(roomID int) []item.Item {
	req := query.NewRequest(roomID, []item.Item{})
	p.service.RoomFloorItemsChan <- req

	resp := <-req.Response
	if resp == nil {
		return nil
	}

	return resp.Value
}

// RoomPublicItems Get public items in a room
func (p *RoomServiceProxy) RoomPublicItems(roomID int) []item.Item {
	req := query.NewRequest(roomID, []item.Item{})
	p.service.RoomPublicItemsChan <- req

	resp := <-req.Response
	if resp == nil {
		return nil
	}

	return resp.Value
}

// RoomIsPublic Is a specific room a public room
func (p *RoomServiceProxy) RoomIsPublic(roomID int) bool {
	req := query.NewRequest(roomID, false)
	p.service.RoomPublicRoomChan <- req

	resp := <-req.Response
	if resp == nil {
		return false
	}

	return resp.Value
}

// UpdateRoom will send a request to update a rooms state to the room service.
func (p *RoomServiceProxy) UpdateRoom(room Room) {
	req := query.NewRequest(room.Details.Id, room)
	p.service.RoomUpdateChan <- req

	<-req.Response // wait for signal from the service that req has been processed.
}

// RoomIsHabboClub Is a specific room club only
func (p *RoomServiceProxy) RoomIsHabboClub(roomID int) bool {
	req := query.NewRequest(roomID, false)
	p.service.RoomIsHCChan <- req

	resp := <-req.Response
	if resp == nil {
		return false
	}

	return resp.Value
}

// RoomSendPacket Send a packet to every player in a room
func (p *RoomServiceProxy) RoomSendPacket(roomID int, packet packets.OutgoingPacket) {
	req := query.NewRequest(roomID, packet)
	p.service.RoomSendPacketChan <- req

	<-req.Response // wait for signal from the service that req has been processed.
}

// RoomScheduleJob schedule job
func (p *RoomServiceProxy) RoomScheduleJob(roomID int, job scheduler.Job) {
	req := query.NewRequest(roomID, job)
	p.service.RoomScheduleJobChan <- req

	<-req.Response // wait for signal from the service that req has been processed.
}

// RoomStopJob stop job
func (p *RoomServiceProxy) RoomStopJob(roomID int, jobName string) {
	req := query.NewRequest(roomID, jobName)
	p.service.RoomStopJobChan <- req

	<-req.Response // wait for signal from the service that req has been processed.
}

// EnterRoom Player enter room
func (p *RoomServiceProxy) EnterRoom(roomID int, ply PlayerEnter) {
	req := query.NewRequest(roomID, ply)
	fmt.Println("Sending enter room req to service")
	p.service.RoomPlayerEnterChan <- req
	fmt.Println("Sent enter room req to service")
	<-req.Response // wait for signal from the service that req has been processed.
}

// LeaveRoom ...
func (p *RoomServiceProxy) LeaveRoom(roomID int, ply player.Player) {
	// send values to room service over channel
	// call room service function to trigger a read on the service
	// room service should find that room and then update state to remove the player
	req := query.NewRequest(roomID, ply)
	p.service.RoomPlayerLeaveChan <- req

	<-req.Response
}

func (p *RoomServiceProxy) MovePlayer(roomID int, move PlayerMove) {
	req := query.NewRequest(roomID, move)
	p.service.RoomPlayerMove <- req

	<-req.Response
}

// TODO Possibly need cmds for getting a tile, tile height, tile state
// TODO Possibly need cmds for a rooms map size, for both x and y
// TODO Possibly need cmds for current and max visitors in a room
