package handlers

import (
	"github.com/jtieri/HabbGo/habbgo/game/model"
	"github.com/jtieri/HabbGo/habbgo/protocol/composers"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
)

func HandleGetInfo(player *model.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeUserObj(player))
}

func HandleGetCredits(player *model.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeCreditBalance(player.Details.Credits))
}

func HandleGetAvailableBadges(player *model.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeAvailableBadges(player))
}

func HandleGetSoundSetting(player *model.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeSoundSetting(player.Details.SoundEnabled))
}

func HandleTestLatency(player *model.Player, packet *packets.IncomingPacket) {
	l := packet.ReadInt()
	player.Session.Send(composers.ComposeLatency(l))
}
