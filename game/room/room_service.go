package room

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/jtieri/habbgo/collections"
	"github.com/jtieri/habbgo/game/item"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/scheduler"
	"github.com/jtieri/habbgo/game/service/query"
	"github.com/jtieri/habbgo/protocol/packets"
	"go.uber.org/zap"
)

// PublicRoomOffset is used as an offset in the navigator.Navigator to easily dichotomize public and private rooms.
const PublicRoomOffset = 1000

const channelBufferSize = 100

// RoomService manages the global server state associated with Rooms and provides facilities for handling Rooms.
type RoomService struct {
	ctx    context.Context
	cancel context.CancelFunc

	repo      RoomRepo                     // Database access object
	scheduler *scheduler.GameScheduler     // Game scheduler for scheduling room tasks
	roomCache collections.Cache[int, Room] // Currently initialized rooms
	channels  *ServiceChannels             // Channels used for reading/writing data to a players session
	running   bool

	log *zap.Logger
}

// NewRoomService returns a new RoomService struct.
func NewRoomService(ctx context.Context, log *zap.Logger, db *sql.DB, scheduler *scheduler.GameScheduler, cancel context.CancelFunc) *RoomService {
	return &RoomService{
		ctx:    ctx,
		cancel: cancel,

		repo:      NewRoomRepo(db),
		scheduler: scheduler,
		roomCache: collections.NewCache(make(map[int]Room)),
		channels:  newServiceChannel(),
		running:   false,

		log: log,
	}
}

// ServiceChannels is a wrapper type for all the channels needed to send and receive requests from/to the service.
type ServiceChannels struct {
	// Channels for reading/writing to the database and RoomService cache
	RoomAddChan     chan *query.Request[int, Room]   // Channel used to receive rooms to add or update
	RoomDeleteChan  chan *query.Request[int, Room]   // Channel used to receive room ID's for rooms to delete
	RoomByPlayerID  chan *query.Request[int, []Room] // Channel used to query for rooms by player id
	RoomByRoomID    chan *query.Request[int, Room]   // Channel used to query for rooms by room id
	RoomsChan       chan *query.Request[int, []Room] // Channel used to return a copy of all cached rooms
	RoomCachedChan  chan *query.Request[int, bool]   // Channel used to query for existance of a room in cache
	RoomCountChan   chan *query.Request[int, int]    // Channel used to query for number of rooms in the cache
	RoomQueriedChan chan *query.Request[int, []Room] // check if any rooms in slice is already initialized, send new slice back

	// Channels for reading/writing to a specific Room's state
	RoomPlayersChan      chan *query.Request[int, []player.Player]
	RoomAddPlayerChan    chan *query.Request[int, player.Player]
	RoomRemovePlayerChan chan *query.Request[int, int]
	RoomPlayerCountChan  chan *query.Request[int, int]
	RoomAddItemsChan     chan *query.Request[int, []item.Item]
	RoomFloorItemsChan   chan *query.Request[int, []item.Item]
	RoomPublicItemsChan  chan *query.Request[int, []item.Item]
	RoomPublicRoomChan   chan *query.Request[int, bool]
	RoomUpdateChan       chan *query.Request[int, Room]

	RoomIsHCChan        chan *query.Request[int, bool]
	RoomSendPacketChan  chan *query.Request[int, packets.OutgoingPacket]
	RoomScheduleJobChan chan *query.Request[int, scheduler.Job]
	RoomStopJobChan     chan *query.Request[int, string]
	RoomPlayerEnterChan chan *query.Request[int, PlayerEnter]
	RoomPlayerLeaveChan chan *query.Request[int, player.Player]
	RoomPlayerMove      chan *query.Request[int, PlayerMove]
}

// newServiceChannel creates a new ServiceChannels object with all of its channels initialized.
func newServiceChannel() *ServiceChannels {
	return &ServiceChannels{
		RoomAddChan:          make(chan *query.Request[int, Room], channelBufferSize),
		RoomDeleteChan:       make(chan *query.Request[int, Room], channelBufferSize),
		RoomByPlayerID:       make(chan *query.Request[int, []Room], channelBufferSize),
		RoomByRoomID:         make(chan *query.Request[int, Room], channelBufferSize),
		RoomsChan:            make(chan *query.Request[int, []Room], channelBufferSize),
		RoomCachedChan:       make(chan *query.Request[int, bool], channelBufferSize),
		RoomCountChan:        make(chan *query.Request[int, int], channelBufferSize),
		RoomQueriedChan:      make(chan *query.Request[int, []Room], channelBufferSize),
		RoomPlayersChan:      make(chan *query.Request[int, []player.Player], channelBufferSize),
		RoomAddPlayerChan:    make(chan *query.Request[int, player.Player], channelBufferSize),
		RoomRemovePlayerChan: make(chan *query.Request[int, int], channelBufferSize),
		RoomPlayerCountChan:  make(chan *query.Request[int, int], channelBufferSize),
		RoomAddItemsChan:     make(chan *query.Request[int, []item.Item], channelBufferSize),
		RoomFloorItemsChan:   make(chan *query.Request[int, []item.Item], channelBufferSize),
		RoomPublicItemsChan:  make(chan *query.Request[int, []item.Item], channelBufferSize),
		RoomPublicRoomChan:   make(chan *query.Request[int, bool], channelBufferSize),
		RoomUpdateChan:       make(chan *query.Request[int, Room], channelBufferSize),
		RoomIsHCChan:         make(chan *query.Request[int, bool], channelBufferSize),
		RoomSendPacketChan:   make(chan *query.Request[int, packets.OutgoingPacket], channelBufferSize),
		RoomScheduleJobChan:  make(chan *query.Request[int, scheduler.Job], channelBufferSize),
		RoomStopJobChan:      make(chan *query.Request[int, string], channelBufferSize),
		RoomPlayerEnterChan:  make(chan *query.Request[int, PlayerEnter], channelBufferSize),
		RoomPlayerLeaveChan:  make(chan *query.Request[int, player.Player], channelBufferSize),
		RoomPlayerMove:       make(chan *query.Request[int, PlayerMove], channelBufferSize),
	}
}

// Start will load the appropriate data and perform the necessary actions to setup the RoomService on startup.
func (rs *RoomService) Start() {
	// Load and cache the public rooms on startup.
	//publicRooms, err := rs.repo.LoadPublicRooms()
	//if err != nil {
	//	panic(err)
	//}
	//
	//rs.log.Debug(
	//	"Loaded public rooms from database",
	//	zap.Int("public_rooms_loaded", len(publicRooms)),
	//)
	//
	//for _, r := range publicRooms {
	//	r.scheduler = rs.scheduler
	//	rs.roomCache.Set(r.Details.Id, r)
	//}

	rs.running = true
	wg := &sync.WaitGroup{}
	for {
		// When the service starts we need to spin up a new goroutine that handles reading/writing
		// for one specific channel. This will allow us to concurrently handle requests from each channel at once.
		for _, handle := range rs.handlers() {
			go handle(wg)
			wg.Add(1)
		}

		// Block here until the context is cancelled and all the worker goroutines die.
		wg.Wait()

		// TODO finish gracefully closing out a rooms state.
		// Main game context has been cancelled, server is shutting down.
		// likely need to save player state to db and send hotel closing packet.
		// save room/item state and then exit.
		for _, r := range rs.roomCache.Items() {
			r.Cleanup()
		}

		rs.running = false
		return
	}
}

func (rs *RoomService) Channels() *ServiceChannels {
	return rs.channels
}

// handlers returns a slice of function callbacks which are used to start the
// various channel listeners in their own goroutines at startup.
func (rs *RoomService) handlers() []func(*sync.WaitGroup) {
	return []func(wg *sync.WaitGroup){
		rs.handleAddRoom,
		rs.handleRemoveRoom,
		rs.handleRoomByPlayerID,
		rs.handleRoomByRoomID,
		rs.handleRoomsCopy,
		rs.handleRoomCached,
		rs.handleRoomCount,
		rs.handleCheckRoomsQueried,
		rs.handleRoomPlayers,
		rs.handleRoomAddPlayer,
		rs.handleRoomRemovePlayer,
		rs.handleRoomPlayerCount,
		rs.handleRoomAddItems,
		rs.handleRoomFloorItems,
		rs.handleRoomPublicItems,
		rs.handleRoomPublic,
		rs.handleUpdateRoom,
		rs.handleRoomIsHC,
		rs.handleRoomSendPacket,
		rs.handleRoomScheduleJob,
		rs.handleRoomStopJob,
		rs.handlePlayerEnterRoom,
		rs.handlePlayerLeaveRoom,
		rs.handleRoomMovePlayer,
	}
}

// handleAddRoom listens for requests to add a new room to the cache.
func (rs *RoomService) handleAddRoom(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomAddChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomAddChan:
			rs.AddRoom(req)
		}
	}
}

// handleRemoveRoom listens for requests to remove a room from the cache.
func (rs *RoomService) handleRemoveRoom(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomDeleteChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomDeleteChan:
			rs.RemoveRoom(req)
		}
	}
}

// handleRoomByPlayerID listens for requests to query for a set of rooms owned by a specific player.
func (rs *RoomService) handleRoomByPlayerID(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomByPlayerID)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomByPlayerID:
			rs.RoomsByPlayerID(req)
		}
	}
}

// handleRoomByRoomID listens for requests to query for a room based off its ID.
func (rs *RoomService) handleRoomByRoomID(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomByRoomID)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomByRoomID:
			rs.RoomByID(req)
		}
	}
}

// handleRoomsCopy listens for requests for a copy of the rooms currently in the cache.
func (rs *RoomService) handleRoomsCopy(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomsChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomsChan:
			rs.Rooms(req)
		}
	}
}

// handleRoomCached listens for requests to query if a room is currently cached or not.
func (rs *RoomService) handleRoomCached(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomCachedChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomCachedChan:
			rs.RoomCached(req)
		}
	}
}

// handleRoomCount listens for requests for the current size of the room cache.
func (rs *RoomService) handleRoomCount(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomCountChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomCountChan:
			rs.RoomsCachedCount(req)
		}
	}
}

// handleCheckRoomsQueried listens for requests to determine if any room in a slice of rooms is already initialized.
// A new slice is built and sent back over the channel containing the already initialized rooms if they exist, along
// with the original rooms sent in the request if there was no room already initialized.
func (rs *RoomService) handleCheckRoomsQueried(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomQueriedChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomQueriedChan:
			rs.CheckRoomsQueried(req)
		}
	}
}

// handleRoomPlayers listens for requests to query for a copy of the current players in a specified room.
func (rs *RoomService) handleRoomPlayers(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomPlayersChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomPlayersChan:
			rs.RoomPlayers(req)
		}
	}
}

// handleRoomAddPlayer listens for requests to add a player to a specified rooms player cache.
func (rs *RoomService) handleRoomAddPlayer(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomAddPlayerChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomAddPlayerChan:
			rs.RoomAddPlayer(req)
		}
	}
}

// handleRoomRemovePlayer listens for requests to remove a player from the player cache of a specified room.
func (rs *RoomService) handleRoomRemovePlayer(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomRemovePlayerChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomRemovePlayerChan:
			rs.RoomRemovePlayer(req)
		}
	}
}

// handleRoomPlayerCount listens for requests for the count of players in a specified room.
func (rs *RoomService) handleRoomPlayerCount(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomPlayerCountChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomPlayerCountChan:
			rs.RoomPlayerCount(req)
		}
	}
}

// handleRoomAddItems listens for requests to add a slice of items to a specified rooms item cache.
func (rs *RoomService) handleRoomAddItems(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomAddItemsChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomAddItemsChan:
			rs.RoomAddItems(req)
		}
	}
}

// handleRoomFloorItems listens for requests to query for the floor items in a specified room.
func (rs *RoomService) handleRoomFloorItems(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomFloorItemsChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomFloorItemsChan:
			rs.RoomFloorItems(req)
		}
	}
}

// handleRoomPublicItems listens for requests to query for the public room items for a specified room.
func (rs *RoomService) handleRoomPublicItems(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomPublicItemsChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomPublicItemsChan:
			rs.RoomPublicItems(req)
		}
	}
}

// handleRoomPublic listens for requests to determine if a specified room is a public room or not.
func (rs *RoomService) handleRoomPublic(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomPublicRoomChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomPublicRoomChan:
			rs.RoomIsPublic(req)
		}
	}
}

// handleUpdateRoom listens for requests to update the state of a cached room.
func (rs *RoomService) handleUpdateRoom(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomUpdateChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomUpdateChan:
			rs.UpdateRoom(req)
		}
	}
}

func (rs *RoomService) handleRoomIsHC(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomIsHCChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomIsHCChan:
			rs.RoomIsHC(req)
		}
	}
}

func (rs *RoomService) handleRoomSendPacket(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomSendPacketChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomSendPacketChan:
			rs.RoomSendPacket(req)
		}
	}
}

func (rs *RoomService) handleRoomScheduleJob(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomScheduleJobChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomScheduleJobChan:
			rs.RoomScheduleJob(req)
		}
	}
}

func (rs *RoomService) handleRoomStopJob(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomStopJobChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomStopJobChan:
			rs.RoomStopJob(req)
		}
	}
}

func (rs *RoomService) handlePlayerEnterRoom(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomPlayerEnterChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomPlayerEnterChan:
			rs.RoomPlayerEnter(req)
		}
	}
}

func (rs *RoomService) handlePlayerLeaveRoom(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomPlayerLeaveChan)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomPlayerLeaveChan:
			rs.RoomPlayerLeave(req)
		}
	}
}

func (rs *RoomService) handleRoomMovePlayer(wg *sync.WaitGroup) {
	defer close(rs.channels.RoomPlayerMove)
	defer wg.Done()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case req := <-rs.channels.RoomPlayerMove:
			rs.RoomMovePlayer(req)
		}
	}
}

// Rooms returns a slice of the currently cached Rooms.
func (rs *RoomService) Rooms(req *query.Request[int, []Room]) {
	req.Query.Value = rs.roomCache.Items()
	req.Response <- req.Query
}

// RoomCached returns true if the Room is in the cache of loaded Rooms or false otherwise.
func (rs *RoomService) RoomCached(req *query.Request[int, bool]) {
	req.Query.Value = rs.roomCache.Has(req.Query.Key)
	req.Response <- req.Query
}

// AddRoom adds a Room to the cache of loaded Rooms.
func (rs *RoomService) AddRoom(req *query.Request[int, Room]) {
	rs.roomCache.SetIfAbsent(req.Query.Key, req.Query.Value)
	req.Response <- nil // this signals to the calling Proxy that the request has been processed
}

// RemoveRoom removes a Room from the cache of loaded Rooms.
func (rs *RoomService) RemoveRoom(req *query.Request[int, Room]) {
	rs.roomCache.Remove(req.Query.Key)
	req.Response <- nil // this signals to the calling Proxy that the request has been processed
}

// RoomByID will check if the room associated with the specified ID is already cached and if so returns it,
// otherwise the room is loaded from the database, added to the cache and returned.
func (rs *RoomService) RoomByID(req *query.Request[int, Room]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if ok {
		req.Query.Value = room
		req.Response <- req.Query
		return
	}

	// if room is not cached load from database, add to cache and return.
	r, err := rs.repo.RoomByID(req.Query.Key)
	if err != nil {
		rs.log.Warn(
			"Failed to retrieve room by ID from database.",
			zap.Int("room_id", req.Query.Key),
			zap.Error(err),
		)
		req.Response <- nil
		return
	}
	r.scheduler = rs.scheduler

	rs.roomCache.Set(r.Details.Id, r)

	req.Query.Value = r
	req.Response <- req.Query
}

// RoomsByPlayerID returns the rooms stored in the database that are owned by the player with the specified ID.
func (rs *RoomService) RoomsByPlayerID(req *query.Request[int, []Room]) {
	rooms, err := rs.repo.RoomsByPlayerId(req.Query.Key)
	if err != nil {
		rs.log.Warn(
			"Failed to retrieve rooms by player ID from database.",
			zap.Int("player_id", req.Query.Key),
			zap.Error(err),
		)
		req.Response <- nil
	}

	req.Query.Value = rooms
	req.Response <- req.Query
}

// CheckRoomsQueried reads a slice of rooms from a channel and checks if any of these rooms are already initialized.
// A new slice of rooms will be composed, containing any replaced rooms along with the rooms not replaced, and sent
// back over the channel.
func (rs *RoomService) CheckRoomsQueried(req *query.Request[int, []Room]) {
	rooms := make([]Room, len(req.Query.Value))

	count := 0
	for _, room := range req.Query.Value {
		r, ok := rs.roomCache.Get(room.Details.Id)
		if ok {
			rooms[count] = r
		} else {
			rooms[count] = room
		}

		count++
	}

	req.Query.Value = rooms
	req.Response <- req.Query
}

func (rs *RoomService) RoomsCachedCount(req *query.Request[int, int]) {
	req.Query.Value = rs.roomCache.Count()
	req.Response <- req.Query
}

func (rs *RoomService) RoomPlayers(req *query.Request[int, []player.Player]) {
	r, ok := rs.roomCache.Get(req.Query.Key)
	if ok {
		req.Query.Value = r.playerCache.Items()
		req.Response <- req.Query
	}

	req.Response <- nil
}

func (rs *RoomService) RoomAddPlayer(req *query.Request[int, player.Player]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if !ok {
		// We should never be adding a player to a room that does not exist.
		req.Response <- nil
		return
	}

	room.playerCache.Set(req.Query.Value.Details.Id, req.Query.Value)
	req.Response <- nil
}

func (rs *RoomService) RoomRemovePlayer(req *query.Request[int, int]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if !ok {
		// We should never be removing a player from a room that does not exist.
		req.Response <- nil
		return
	}

	room.playerCache.Remove(req.Query.Value)
	req.Response <- nil
}

func (rs *RoomService) RoomPlayerCount(req *query.Request[int, int]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if !ok {
		req.Response <- nil
	}

	req.Query.Value = room.playerCache.Count()
	req.Response <- req.Query
}

func (rs *RoomService) RoomAddItems(req *query.Request[int, []item.Item]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if !ok {
		// We should never be trying to add items to a room that does not exist
		req.Response <- nil
		return
	}

	room.AddAllItems(req.Query.Value)
	req.Response <- nil
}

func (rs *RoomService) RoomFloorItems(req *query.Request[int, []item.Item]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if !ok {
		// We should never be trying to read items from a room that does not exist
		req.Response <- nil
		return
	}

	req.Query.Value = room.floorItems()
	req.Response <- req.Query
}

func (rs *RoomService) RoomPublicItems(req *query.Request[int, []item.Item]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if !ok {
		// We should never be trying to read items from a room that does not exist
		req.Response <- nil
		return
	}

	req.Query.Value = room.publicRoomItems()
	req.Response <- req.Query
}

func (rs *RoomService) RoomIsPublic(req *query.Request[int, bool]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if !ok {
		// We should never be trying to determine room type for a room that does not exist
		req.Response <- nil
		return
	}

	req.Query.Value = room.publicRoom()
	req.Response <- req.Query
}

func (rs *RoomService) UpdateRoom(req *query.Request[int, Room]) {
	_, ok := rs.roomCache.Get(req.Query.Key)
	if ok {
		// TODO update the rooms state appropriately
		rs.roomCache.Set(req.Query.Key, req.Query.Value)
	}

	req.Response <- nil
}

func (rs *RoomService) RoomIsHC(req *query.Request[int, bool]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if ok {
		req.Query.Value = room.habboClubOnly()
		req.Response <- req.Query
		return
	}

	req.Response <- nil
}

func (rs *RoomService) RoomSendPacket(req *query.Request[int, packets.OutgoingPacket]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if ok {
		room.Send(req.Query.Value, req.Query.Value)
		req.Response <- nil
		return
	}

	req.Response <- nil
}

func (rs *RoomService) RoomScheduleJob(req *query.Request[int, scheduler.Job]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if ok {
		room.scheduleJob(req.Query.Value)
		req.Response <- nil
		return
	}

	req.Response <- nil
}

func (rs *RoomService) RoomStopJob(req *query.Request[int, string]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if ok {
		room.stopJob(req.Query.Value)
		req.Response <- nil
		return
	}

	req.Response <- nil
}

func (rs *RoomService) RoomPlayerEnter(req *query.Request[int, PlayerEnter]) {
	fmt.Println("Looking for room to enter in service cache")
	room, ok := rs.roomCache.Get(req.Query.Key)
	if ok {
		fmt.Println("Room was in service cache")
		err := room.enterRoom(req.Query.Value.Player, req.Query.Value.Jobs...)
		if err != nil {
			rs.log.Error("Player failed to enter room",
				zap.String("player_name", req.Query.Value.Player.Details.Username),
				zap.String("room_name", room.Details.Name),
				zap.Error(err),
			)
		}
	}

	fmt.Println("Sending response to proxy from service")

	req.Response <- nil
}

func (rs *RoomService) RoomMovePlayer(req *query.Request[int, PlayerMove]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if !ok {
		req.Response <- nil
		return
	}

	room.moveTo(req.Query.Value)
	req.Response <- nil
}

func (rs *RoomService) RoomPlayerLeave(req *query.Request[int, player.Player]) {
	room, ok := rs.roomCache.Get(req.Query.Key)
	if !ok {
		req.Response <- nil
	}

	err := room.leaveRoom(req.Query.Value)
	if err != nil {
		rs.log.Error("Player failed to leave room",
			zap.String("player_name", req.Query.Value.Details.Username),
			zap.String("room_name", room.Details.Name),
			zap.Error(err),
		)
	}

	req.Response <- nil
}

// PublicName returns a more readable version of the room name for certain public rooms.
func PublicName(room *Room) string {
	if room.publicRoom() {
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
