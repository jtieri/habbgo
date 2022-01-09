package commands

import (
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/protocol/messages"
	"github.com/jtieri/habbgo/protocol/packets"
)

func GET_INFO(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(player.Details.Username, messages.USEROBJ(player))
}

func GET_CREDITS(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(player.Details.Username, messages.CREDITBALANCE(player.Details.Credits))
}

func GETAVAILABLEBADGES(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(player.Details.Username, messages.AVAILABLEBADGES(player))
}

func GET_SOUND_SETTING(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(player.Details.Username, messages.SOUNDSETTING(player.Details.SoundEnabled))
}

func TestLatency(player *player.Player, packet *packets.IncomingPacket) {
	l := packet.ReadInt()
	player.Session.Send(player.Details.Username, messages.Latency(l))
}