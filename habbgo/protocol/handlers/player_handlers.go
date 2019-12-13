package handlers

import (
	"github.com/jtieri/HabbGo/habbgo/game/model/player"
	"github.com/jtieri/HabbGo/habbgo/protocol/composers"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
)

func HandleGetInfo(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeUserObj(player))
}

func HandleGetCredits(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeCreditBalance(player.Details.Credits))
}

func HandleGetAvailableBadges(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeAvailableBadges())
}
