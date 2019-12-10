package player

import (
	"github.com/jtieri/HabbGo/habbgo/server/protocol/packets"
)

type Player struct {
	Session Network
	PlayerDetails *PlayerDetails
}

type PlayerDetails struct {
	Username string
	Motto string
	Sex rune
	Figure string
}

type Network interface {
	Listen()
	Send(packet *packets.OutgoingPacket)
	Queue(packet *packets.OutgoingPacket)
	Flush(packet *packets.OutgoingPacket)
	Close()
}