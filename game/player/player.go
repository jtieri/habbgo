package player

import (
	"database/sql"
	"github.com/jtieri/habbgo/protocol/packets"
	"time"
)

type Player struct {
	Session  Session
	Details  *Details
	Database *sql.DB
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
	CurrentBadge string
	DisplayBadge bool
	SoundEnabled bool
}

type Session interface {
	Listen()
	Send(playerIdentifier string, caller interface{}, packet *packets.OutgoingPacket)
	Queue(packet *packets.OutgoingPacket)
	Flush(playerIdentifier string, caller interface{}, packet *packets.OutgoingPacket)
	Address() string
	GetPacketCommand(headerId int) (func(*Player, *packets.IncomingPacket), bool)
	Close()
}

func New(session Session, database *sql.DB) *Player {
	return &Player{
		Session:  session,
		Database: database,
		Details:  &Details{},
	}
}
