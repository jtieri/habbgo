package handlers

import (
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/protocol/composers"
	"github.com/jtieri/habbgo/protocol/packets"
)

func GetInfo(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeUserObj(player))
}

func GetCredits(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeCreditBalance(player.Details.Credits))
}

func GetAvailableBadges(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeAvailableBadges(player))
}

func GetSoundSetting(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeSoundSetting(player.Details.SoundEnabled))
}

func TestLatency(player *player.Player, packet *packets.IncomingPacket) {
	l := packet.ReadInt()
	player.Session.Send(composers.ComposeLatency(l))
}
