package player

import (
	"database/sql"
	"strings"
	"time"

	"github.com/jtieri/habbgo/game/navigator"
	"github.com/jtieri/habbgo/game/ranks"
	"github.com/jtieri/habbgo/game/room"
	"github.com/jtieri/habbgo/protocol/packets"
	"go.uber.org/zap"
)

type Player struct {
	Session Session
	Details *Details

	Database *sql.DB
	Services ServiceManager

	log *zap.Logger
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
	Listen()
	Send(caller interface{}, packet *packets.OutgoingPacket)
	Queue(packet *packets.OutgoingPacket)
	Flush(caller interface{}, packet *packets.OutgoingPacket)
	Address() string
	GetPacketCommand(headerId int) (func(*Player, *packets.IncomingPacket), bool)
	Close()
}

type Service interface {
	Build()
}

type ServiceManager interface {
	RoomService() *room.RoomService
	PlayerService() *PlayerService
	NavigatorService() *navigator.NavService
}

func New(log *zap.Logger, session Session, database *sql.DB, s ServiceManager) *Player {
	return &Player{
		Session:  session,
		Database: database,
		Details:  &Details{},
		Services: s,
		log:      log,
	}
}

func PlayerRank(rankString string) ranks.Rank {
	switch strings.ToLower(rankString) {
	case "none":
		return ranks.None
	case "normal":
		return ranks.Normal
	case "community manager":
		return ranks.CommunityManager
	case "guide":
		return ranks.Guide
	case "hobba":
		return ranks.Hobba
	case "super hobba":
		return ranks.SuperHobba
	case "moderator":
		return ranks.Moderator
	case "administrator":
		return ranks.Administrator
	default:
		return ranks.Normal
	}
}

func (p *Player) Login() {
	// Set player logged in & ping ready for latency test
	// Possibly add player to a list of online players? Health endpoint with server stats?
	// Save current time to Conn for players last online time

	// Check if player is banned & if so send USER_BANNED
	// Log IP address to Conn

	LoadBadges(p)

	// If Config has alerts enabled, send player ALERT

	// Check if player gets club gift & update club status
}

func (p *Player) Register(username, figure, gender, email, birthday, createdAt, password string, salt []byte) {
	err := Register(p, username, figure, gender, email, birthday, createdAt, password, salt)
	if err != nil {
		p.log.Warn("Failed to register player",
			zap.String("username", username),
			zap.Error(err),
		)
	}
}
