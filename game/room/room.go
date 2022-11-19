package room

import (
	"context"
	"fmt"
	"time"

	"github.com/jtieri/habbgo/collections"
	"github.com/jtieri/habbgo/game/item"
	"github.com/jtieri/habbgo/game/navigator"
	"github.com/jtieri/habbgo/game/pathfinder/position"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/scheduler"
	"github.com/jtieri/habbgo/num"
	"github.com/jtieri/habbgo/protocol/packets"
	"golang.org/x/sync/errgroup"
)

const (
	WallpaperProperty = "wallpaper"
	FloorProperty     = "floor"
)

// Room represents an in-game room.
type Room struct {
	ctx    context.Context
	cancel context.CancelFunc

	Details     Details
	Model       Model
	mapping     Map
	Initialized bool

	playerCache collections.Cache[int, player.Player]
	itemCache   collections.Cache[int, item.Item]

	scheduler   *scheduler.GameScheduler
	runningJobs collections.Cache[string, scheduler.Job]

	// Ready is set to true when a Room is built using the constructor.
	// this allows us to check for empty Room values by checking if Ready == false.
	Ready bool
}

// Details are the metadata describing a Room.
type Details struct {
	Id              int
	CategoryID      int
	Name            string
	Description     string
	CCTs            string
	Wallpaper       int
	Floor           int
	Landscape       float32
	OwnerId         int
	OwnerName       string
	ShowOwner       bool
	SudoUsers       bool
	Hidden          bool
	AccessType      Access
	Password        string
	CurrentVisitors int
	MaxVisitors     int
	Rating          int
	ChildRooms      []*Room
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewRoom returns a new Room struct.
func NewRoom() Room {
	return Room{
		Details:     Details{},
		Model:       Model{},
		mapping:     Map{},
		Initialized: false,
		playerCache: collections.NewCache(make(map[int]player.Player)),
		itemCache:   collections.NewCache(make(map[int]item.Item)),
		runningJobs: collections.NewCache(make(map[string]scheduler.Job)),
		Ready:       true,
	}
}

// Initialize is called the first time the room is loaded when a player enters.
// The function will run in a goroutine that stays alive until it's context is cancelled.
// The context should only be getting cancelled if the room is removed from the cache of active rooms,
// or if the server is being shutdown.
func (r *Room) Initialize() {
	// TODO perform logic for a room upon initialization

	for {
		select {
		case <-r.ctx.Done():
			// TODO rooms get cleaned up when there are no Players inside and they are removed from the
			// cache of rooms in the room service. The context will be cancelled and we should clean up the room here.

			r.StopRoomJobs()
		default:
			_ = r // TODO remove this, just using this to silence linter

			// if there is anything else to do then do that here.
			// we could constantly check for some trigger or change in state that does S O M E T H I N G
		}
	}
}

// Cleanup ...
// TODO cleanup should completely wipe a rooms state/running jobs and prepare the room
// for being removed from the RoomService cache. We should check that if the call to Cleanup
// comes with players in the room that their state is persisted properly and that their state is
// cleaned up properly.
func (r *Room) Cleanup() {
	// TODO reset item states?

	r.Initialized = false
	r.ClearPlayers()

	// TODO clear Items, Players, votes, etc.
}

// Players returns a channel containing all the (key, value) pairs from Room's cache
// of Players.
func (r *Room) Players() []player.Player {
	return r.playerCache.Items()
}

// AddAllItems adds all the new Items into the Room's list of active Items.
func (r *Room) AddAllItems(newItems []item.Item) {
	for _, i := range newItems {
		r.itemCache.SetIfAbsent(i.ID, i)
	}
	fmt.Printf("Added %d items to room \n", r.itemCache.Count())
}

// AddPlayer adds a player.Player to the map of Players currently in the room.
func (r *Room) AddPlayer(instanceID int, player player.Player) {
	r.playerCache.SetIfAbsent(instanceID, player)
}

// removePlayer removes a player.Player from the map of Players currently in the room.
func (r *Room) removePlayer(playerID int) {
	r.playerCache.Remove(playerID)
}

// PlayerCount returns the number of Players in the map of active Players.
func (r *Room) PlayerCount() int {
	return r.playerCache.Count()
}

// ClearPlayers removes all the Players from the map of activer Players.
func (r *Room) ClearPlayers() {
	r.playerCache.Clear()
}

// Tile returns the tile in the room's map at the specified (x, y) coordinate.
func (r *Room) Tile(x, y int) (Tile, error) {
	if x >= r.mapping.sizeX || y >= r.mapping.sizeY {
		return Tile{}, ErrIndexOutOfBounds
	}

	if r.tileState(x, y) == Inaccessible {
		return Tile{}, ErrTileInaccessible
	}

	tile := r.mapping.tiles[x][y]

	return tile, nil
}

// tileHeight returns the tile height for the Tile at the given (x, y) coordinate.
// If the (x, y) coordinate is invalid then 0 is returned.
func (r *Room) tileHeight(x, y int) float64 {
	if x >= r.mapping.sizeX || y >= r.mapping.sizeY {
		return 0
	}

	return r.mapping.tiles[x][y].Height
}

// TileState returns the TileState for the Tile at the given (x, y) coordinate.
// If the (x, y) coordinate is invalid then Inaccessible is returned.
func (r *Room) tileState(x, y int) TileState {
	if x < 0 || y < 0 || x >= r.mapping.sizeX || y >= r.mapping.sizeY {
		return Inaccessible
	}

	return r.mapping.tiles[x][y].State
}

// floorItems returns a slice containing all the regular floor furniture Items that are in the Room.
func (r Room) floorItems() []item.Item {
	var items []item.Item

	for _, i := range r.itemCache.Items() {
		if i.Definition.ContainsBehavior(item.PublicSpaceObject) || i.Definition.ContainsBehavior(item.WallItem) {
			continue
		}

		items = append(items, i)
	}
	return items
}

// publicRoomItems returns a slice containing all the public room furniture Items that are in the Room.
func (r Room) publicRoomItems() []item.Item {
	var items []item.Item

	for _, i := range r.itemCache.Items() {
		if i.Definition.ContainsBehavior(item.Invisible) ||
			i.Definition.ContainsBehavior(item.PrivateFurniture) ||
			!i.Definition.ContainsBehavior(item.PublicSpaceObject) {
			continue
		}

		items = append(items, i)
	}

	fmt.Printf("There are %d public items in this room \n", len(items))
	return items
}

// publicRoom returns true if the Room is a public room.
func (r *Room) publicRoom() bool {
	// We use the owner ID of 0 in the database for public rooms.
	return r.Details.OwnerId == 0
}

// MapSizeX returns the size of the room Map's X dimension.
func (r Room) MapSizeX() int {
	return r.mapping.sizeX
}

// MapSizeY returns the size of the room Map's Y dimension.
func (r Room) MapSizeY() int {
	return r.mapping.sizeY
}

// Send sends a packet to every player currently in the room.
func (r *Room) Send(caller interface{}, packet packets.OutgoingPacket) {
	for _, p := range r.playerCache.Items() {
		p.Session.Send(caller, packet)
	}
}

// StartRoomJobs schedules the appropriate jobs that should be running once a room is initialized.
// This function should only be called once when the room is initialized.
//func (r *Room) StartRoomJobs() {
//	ctx, cancel := context.WithCancel(context.Background())
//	r.scheduler.scheduleJob(jobs.NewPlayerJob(ctx, cancel, r))
//
//	// TODO handle rollers when implemented
//}

// StopRoomJobs iterates over all the running jobs and calls their Stop function.
// After stopping all running jobs the map will be cleared.
func (r *Room) StopRoomJobs() {
	var eg errgroup.Group

	for _, job := range r.runningJobs.Items() {
		eg.Go(func() error {
			job.Stop()
			for !job.Running() {
			}
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return
	}
	r.runningJobs.Clear()
}

// scheduleJob sends a new job to the scheduler.GameScheduler to be executed.
// If there is already a job running with the same name then this call is a no-op.
func (r *Room) scheduleJob(j scheduler.Job) {
	// We don't want two instances of the same job running at once.
	if r.runningJobs.Has(j.Name()) {
		return
	}

	r.scheduler.ScheduleJob(j)
	r.runningJobs.Set(j.Name(), j)
}

// hasJob returns true if there is a running job for the specified key name.
func (r *Room) hasJob(name string) bool {
	return r.runningJobs.Has(name)
}

// stopJob calls the Stop function of a running job if there is a match for the specified key name.
// If there is no match for the key name then this call is a no-op.
func (r *Room) stopJob(name string) {
	j, ok := r.runningJobs.Get(name)
	if ok {
		j.Stop()
		r.runningJobs.Remove(name)
	}
}

// getJob returns a running job from the map of running jobs using the specified key name.
func (r *Room) getJob(name string) scheduler.Job {
	j, ok := r.runningJobs.Get(name)
	if ok {
		return j
	}
	return nil
}

func (r *Room) habboClubOnly() bool {
	return false
}

type PlayerEnter struct {
	Player player.Player
	Jobs   []scheduler.Job
}

func NewPlayerEnter(p player.Player, jobs []scheduler.Job) PlayerEnter {
	return PlayerEnter{
		Player: p,
		Jobs:   jobs,
	}
}

func (r *Room) enterRoom(p player.Player, jobs ...scheduler.Job) error {

	// TODO make sure p is out of previous room and take appropriate actions updating state
	if p.InRoom() {
		panic("Player was already in room and entered a new room")
	}

	fmt.Println("Setting p state")
	p.State().RoomID = r.Details.Id
	p.SetInRoom(true)
	p.State().InstanceID = int(num.RandomInt(1000000)) // TODO set instance id to a random integer id, each entity in the room needs a unique ID

	fmt.Println("Setting p position")
	doorPos := r.Model.Door.Position()
	p.SetPosition(doorPos)

	ps := p.Services.PlayerService().(*player.PlayerServiceProxy)
	ps.UpdatePlayer(p)

	fmt.Println("Adding p to room")
	r.AddPlayer(p.State().InstanceID, p)

	fmt.Println("Loading room")
	r.loadRoom(p, jobs...)

	fmt.Println("Leaving enterRoom call")
	// TODO set visitor count in room
	return nil
}

func (r *Room) loadRoom(player player.Player, jobs ...scheduler.Job) {
	// The room is already initialized
	if r.Initialized {
		fmt.Println("Room already initialized")
		return
	}

	// TODO sanitize room state on initialization to ensure its starting from a clean state
	// e.g.

	fmt.Println("Adding public room items to public room")
	itemService := player.Services.ItemService().(*item.ItemServiceProxy)
	if r.publicRoom() {
		items := itemService.PublicItems(r.Details.Id, r.Model.Name)
		r.AddAllItems(items)
	}
	// TODO else if public room load the rooms Items from the database

	// TODO possible reset room item states? Is this how Habbo worked?
	// e.g. you turn on some holoboys, leave room, come back and they are off

	// TODO regenerate room collision map

	fmt.Println("Room initialized")
	r.Initialized = true

	fmt.Println("Scheduling room jobs")
	for _, job := range jobs {
		r.scheduleJob(job)
	}
}

func (r *Room) leaveRoom(p player.Player) error {
	// TODO finish logic for sanitizing a Players state when they leave a room
	r.removePlayer(p.Details.Id)

	// Remove player from the cache of Players on a certain tile on the room map.
	tile, err := r.Tile(p.State().Position.X, p.State().Position.Y)
	if err != nil {
		return err
	}

	tile.RemovePlayer(p.Details.Id)

	// TODO update rooms player count and persist to database
	// TODO send logout or hotel view packets
	// TODO update room status for messenger in the future

	// TODO ideally we want to run this in a job with a delay so that if the room needs to be
	// initialized right after this player leaves we can cancel room cleanup vs. re-initializing the room.
	// We only want to run this job if there are no Players in the room.
	if r.PlayerCount() == 0 {
		// TODO reset item states?

		r.Initialized = false
		r.StopRoomJobs()

		// TODO clear Items, Players, votes, etc.
		r.ClearPlayers()
	}

	//p, ok := r.playerCache.Get(p.Details.Id)
	//if !ok {
	//	return fmt.Errorf("failed to get player with ID %d from cache of room %s", p.Details.Id, r.Details.Name)
	//}

	// Remove player from Room's cache of players.
	r.playerCache.Remove(p.Details.Id)

	// TODO update rooms player count and persist to database
	// TODO send logout or hotel view packets to player
	// TODO update room status for messenger in the future

	// TODO ideally we want to run this in a job with a delay so that if the room needs to be
	// initialized right after this player leaves we can cancel room cleanup vs. re-initializing the room.
	// We only want to run this job if there are no Players in the room.

	if r.PlayerCount() == 0 {
		r.Cleanup()
	}

	// Reset players state
	// TODO finish logic for sanitizing a Players state when they leave a room

	p.State().RoomID = 0
	p.SetInRoom(false)

	ps := p.Services.PlayerService().(*player.PlayerServiceProxy)
	ps.UpdatePlayer(p)

	return nil
}

type PlayerMove struct {
	Player         player.Player
	XCoord, YCoord int
}

func NewPlayerMove(p player.Player, xCoord, yCoord int) PlayerMove {
	return PlayerMove{
		Player: p,
		XCoord: xCoord,
		YCoord: yCoord,
	}
}

func (r *Room) moveTo(move PlayerMove) {
	p := move.Player
	ps := p.State()

	if (ps.NextPosition != position.Position{}) {
		oldPos := ps.Position

		_ = oldPos

		ps.Position.X = ps.NextPosition.X
		ps.Position.Y = ps.NextPosition.Y

		tile, err := r.Tile(ps.Position.X, ps.Position.Y)
		if err != nil {
			return
		}

		height := tile.WalkingHeight()
		oldHeight := ps.Position.Z

		if height != oldHeight {
			ps.Position.Z = height
			p.SetUpdate(true)

			ps := p.Services.PlayerService().(*player.PlayerServiceProxy)
			ps.UpdatePlayer(p)
		}

		pCurrentItem := tile.TopItem
		if err != nil {
			return
		}

		// TODO handle item trigger for what happens when a player walks

		endPos := position.Position{X: move.XCoord, Y: move.YCoord}

		if ValidTile(endPos, *r) {
			// TODO finish back off here for walking
			_ = pCurrentItem
			_ = tile
		}

	}

	// TODO this function is called from the goroutine spun up in session.go during the call to handle(player, packet)
	// which means that if the players state is mutated here it needs to be propagated to the cached room in RoomService
}

func CurrentVisitors(cat navigator.Category, rooms []Room) int {
	visitors := 0

	for _, r := range rooms {
		if r.Details.CategoryID == cat.ID {
			visitors += r.Details.CurrentVisitors
		}
	}

	return visitors
}

func MaxVisitors(cat navigator.Category, rooms []Room) int {
	visitors := 0

	for _, r := range rooms {
		if r.Details.CategoryID == cat.ID {
			visitors += r.Details.MaxVisitors
		}
	}

	return visitors
}
