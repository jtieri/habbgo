package player

import (
	"container/list"
	"context"
	"database/sql"
	"time"

	"github.com/jtieri/habbgo/game/pathfinder/position"
	"github.com/jtieri/habbgo/game/ranks"
	"github.com/jtieri/habbgo/game/room/actions"
	"github.com/jtieri/habbgo/game/types"
	"github.com/jtieri/habbgo/protocol/packets"
	"go.uber.org/zap"
)

type Player struct {
	Ctx context.Context

	state    *PlayerState
	Session  Session
	Details  *Details
	Repo     PlayerRepo
	Services types.ServiceProxies
	loggedIn bool

	log *zap.Logger
}

// PlayerState represents a player.Player's state in a room.Room.
type PlayerState struct {
	InstanceID int
	StateType  string // TODO in the future this will be used to distinguish between different types e.g. Players, pets, bots, Items

	RoomID    int
	ModelName string

	Position     position.Position
	NextPosition position.Position
	EndPosition  position.Position
	Path         list.List

	needsUpdate bool
	canMove     bool
	moving      bool
	inRoom      bool
	kicked      bool

	Actions map[string]actions.Action
}

type Details struct {
	Id           int
	Username     string
	Figure       string
	Sex          string
	Motto        string
	ConsoleMotto string
	Tickets      int
	PoolFigure   string
	Film         int
	Credits      int
	LastOnline   time.Time
	Badges       []string
	PlayerRank   ranks.Rank
	CurrentBadge string
	DisplayBadge bool
	SoundEnabled bool
}

type Session interface {
	Send(caller interface{}, packet packets.OutgoingPacket)
	Queue(caller interface{}, packet packets.OutgoingPacket)
	Flush(caller interface{}, packet packets.OutgoingPacket)
	Close()
}

func New(ctx context.Context, log *zap.Logger, session Session, database *sql.DB, s types.ServiceProxies) Player {
	return Player{
		Ctx:      ctx,
		state:    NewPlayerState(),
		Session:  session,
		Repo:     NewPlayerRepo(database),
		Details:  new(Details),
		Services: s,
		loggedIn: false,
		log:      log,
	}
}

func NewPlayerState() *PlayerState {
	return &PlayerState{
		Actions:     make(map[string]actions.Action),
		inRoom:      false,
		canMove:     false,
		needsUpdate: false,
		kicked:      false,
		moving:      false,
	}
}

func (p *Player) Cleanup() {
	// TODO finish this
}

func (p *Player) Login() {
	p.loggedIn = true

	// Set player logged in & ping ready for latency test
	// Possibly add player to a list of online players? Health endpoint with server stats?
	// Save current time to Conn for players last online time

	// Check if player is banned & if so send USER_BANNED
	// Log IP address to Conn

	p.Repo.LoadBadges(p)

	// If Config has alerts enabled, send player ALERT

	// Check if player gets club gift & update club status
}

func (p *Player) Register(username, figure, gender, email, birthday, createdAt, password string, salt []byte) {
	err := p.Repo.Register(username, figure, gender, email, birthday, createdAt, password, salt)
	if err != nil {
		p.log.Warn("Failed to register player",
			zap.String("username", username),
			zap.Error(err),
		)
	}
}

func (p *Player) NeedsUpdate() bool {
	return p.state.needsUpdate
}

func (p *Player) SetUpdate(update bool) {
	p.state.needsUpdate = update
}

func (p *Player) State() *PlayerState {
	return p.state
}

func (p *Player) InRoom() bool {
	return p.state.inRoom
}

func (p *Player) SetInRoom(inRoom bool) {
	p.state.inRoom = inRoom
}

func (p *Player) SetPosition(pos position.Position) {
	p.state.Position = pos
}

func (p *Player) CanMove() bool {
	return p.state.canMove
}

func (p *Player) LoggedIn() bool {
	return p.loggedIn
}
