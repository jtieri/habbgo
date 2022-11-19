package player

import (
	"context"
	"database/sql"
	"sync"

	"github.com/jtieri/habbgo/collections"
	"github.com/jtieri/habbgo/game/scheduler"
	"github.com/jtieri/habbgo/game/service/query"
	"go.uber.org/zap"
)

/*
When the server accepts a new connection, it creates a new session and starts a goroutine that listens for data from
that session that comes in over the wire. The session's listener creates a player and adds it to the PlayerService and
keeps a copy of the PlayerID which will be passed into the handle() call for a callback when a valid packet is read.
The handle() call is started in a new goroutine, so it will query the player from the player service at the start of
each call.

This will allow us to remove the reference to player in a session and makes sure that each time a packet handler
is called we query for the player. This way we can make sure that there is a global accessible reference to the player
which can always be queried and updated from other services, jobs, etc.
*/

const channelBufferSize = 100

type PlayerService struct {
	ctx    context.Context
	cancel context.CancelFunc

	repo        PlayerRepo
	scheduler   *scheduler.GameScheduler          // Game scheduler for scheduling
	playerCache collections.Cache[string, Player] // Currently initialized players
	channels    *ServiceChannels                  // Channels used for reading/writing data to a players session
	running     bool

	log *zap.Logger
}

func NewPlayerService(ctx context.Context, log *zap.Logger, db *sql.DB, scheduler *scheduler.GameScheduler, cancel context.CancelFunc) *PlayerService {
	return &PlayerService{
		ctx:    ctx,
		cancel: cancel,

		repo:        NewPlayerRepo(db),
		scheduler:   scheduler,
		playerCache: collections.NewCache(make(map[string]Player)),
		channels:    newServiceChannel(),
		running:     false,

		log: log,
	}
}

type ServiceChannels struct {
	AddPlayer    chan query.Request[string, Player]
	RemovePlayer chan query.Request[string, Player]
	UpdatePlayer chan query.Request[string, Player]
	HasPlayer    chan query.Request[string, Player]
	GetPlayer    chan query.Request[string, Player]
}

func newServiceChannel() *ServiceChannels {
	return &ServiceChannels{
		AddPlayer:    make(chan query.Request[string, Player], channelBufferSize),
		RemovePlayer: make(chan query.Request[string, Player], channelBufferSize),
		UpdatePlayer: make(chan query.Request[string, Player], channelBufferSize),
		HasPlayer:    make(chan query.Request[string, Player], channelBufferSize),
		GetPlayer:    make(chan query.Request[string, Player], channelBufferSize),
	}
}

func (ps *PlayerService) Start() {
	ps.running = true
	wg := &sync.WaitGroup{}
	for {
		// When the service starts we need to spin up a new goroutine that handles reading/writing
		// for one specific channel. This will allow us to concurrently handle requests from each channel at once.
		for _, handle := range ps.handlers() {
			go handle(wg)
			wg.Add(1)
		}

		// Block here until the context is cancelled and all the worker goroutines die.
		wg.Wait()

		// TODO finish gracefully closing out a players state
		// Main game context has been cancelled, server is shutting down.
		for _, p := range ps.playerCache.Items() {
			p.Cleanup()
		}

		ps.running = false

		// listen for incoming tasks
		// for each task coming in,
		return
	}
}

func (ps *PlayerService) Channels() *ServiceChannels {
	return ps.channels
}

func (ps *PlayerService) handlers() []func(group *sync.WaitGroup) {
	return []func(wg *sync.WaitGroup){
		ps.handleAddPlayer,
		ps.handleRemovePlayer,
		ps.handleUpdatePlayer,
		ps.handleGetPlayer,
	}
}

func (ps *PlayerService) handleAddPlayer(wg *sync.WaitGroup) {
	defer close(ps.channels.AddPlayer)
	defer wg.Done()

	for {
		select {
		case <-ps.ctx.Done():
			return
		case req := <-ps.channels.AddPlayer:
			ps.addPlayer(req)
		}
	}
}

func (ps *PlayerService) handleRemovePlayer(wg *sync.WaitGroup) {
	defer close(ps.channels.RemovePlayer)
	defer wg.Done()

	for {
		select {
		case <-ps.ctx.Done():
			return
		case req := <-ps.channels.RemovePlayer:
			ps.removePlayer(req)
		}
	}
}

func (ps *PlayerService) handleUpdatePlayer(wg *sync.WaitGroup) {
	defer close(ps.channels.UpdatePlayer)
	defer wg.Done()

	for {
		select {
		case <-ps.ctx.Done():
			return
		case req := <-ps.channels.UpdatePlayer:
			ps.updatePlayer(req)
		}
	}
}

func (ps *PlayerService) handleGetPlayer(wg *sync.WaitGroup) {
	defer close(ps.channels.GetPlayer)
	defer wg.Done()

	for {
		select {
		case <-ps.ctx.Done():
			return
		case req := <-ps.channels.GetPlayer:
			ps.getPlayer(req)
		}
	}
}

type Task interface {
	Execute()
}

type AddPlayerTask struct {
	key   string
	value Player
}

func (t AddPlayerTask) Execute() {

}

func (ps *PlayerService) addPlayer(req query.Request[string, Player]) {
	ps.playerCache.SetIfAbsent(req.Query.Key, req.Query.Value)
	req.Response <- nil // Write to the channel to signal op complete.
}

func (ps *PlayerService) removePlayer(req query.Request[string, Player]) {
	ps.playerCache.Remove(req.Query.Key)
	req.Response <- nil // Write to the channel to signal op complete.
}

func (ps *PlayerService) updatePlayer(req query.Request[string, Player]) {
	ps.playerCache.Set(req.Query.Key, req.Query.Value)
	req.Response <- nil // Write to the channel to signal op complete.
}

func (ps *PlayerService) getPlayer(req query.Request[string, Player]) {
	p, ok := ps.playerCache.Get(req.Query.Key)
	if ok {
		req.Query.Value = p
		req.Response <- req.Query // Write to the channel to signal op complete.
		return
	}
	req.Response <- nil // Write to the channel to signal op complete.
}
