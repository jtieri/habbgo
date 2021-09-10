package player

import (
	"database/sql"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
)

type Player struct {
	Session Session
	Details *Details
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
	LastOnline   string
	Badges       []string
	CurrentBadge string
	DisplayBadge bool
	SoundEnabled int
}

type Session interface {
	Listen()
	Send(packet *packets.OutgoingPacket)
	Queue(packet *packets.OutgoingPacket)
	Flush(packet *packets.OutgoingPacket)
	Database() *sql.DB
	Address() string
	GetPacketHandler(headerId int) (func(*Player, *packets.IncomingPacket), bool)
	Close()
}

func New(session Session) *Player {
	return &Player{
		Session: session,
		Details: &Details{},
	}
}
