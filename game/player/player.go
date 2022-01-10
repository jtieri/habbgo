package player

import (
	"database/sql"
	"github.com/jtieri/habbgo/protocol/packets"
	"strings"
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
	PlayerRank   Rank
	CurrentBadge string
	DisplayBadge bool
	SoundEnabled bool
}

type Rank int

const (
	None Rank = iota
	Normal
	CommunityManager
	Guide
	Hobba
	SuperHobba
	Moderator
	Administrator
)

func (r Rank) String() string {
	switch r {
	case 0:
		return "none"
	case 1:
		return "normal"
	case 2:
		return "community manager"
	case 3:
		return "guide"
	case 4:
		return "hobba"
	case 5:
		return "super hobba"
	case 6:
		return "moderator"
	case 7:
		return "administrator"
	default:
		return "normal"
	}
}
func PlayerRank(rankString string) Rank {
	switch strings.ToLower(rankString) {
	case "none":
		return None
	case "normal":
		return Normal
	case "community manager":
		return CommunityManager
	case "guide":
		return Guide
	case "hobba":
		return Hobba
	case "super hobba":
		return SuperHobba
	case "moderator":
		return Moderator
	case "administrator":
		return Administrator
	default:
		return Normal
	}
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
