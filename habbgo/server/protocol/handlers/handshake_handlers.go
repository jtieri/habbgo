package handlers

import (
	"github.com/jtieri/HabbGo/habbgo/game/model/player"
	"github.com/jtieri/HabbGo/habbgo/server/protocol/composers"
	"github.com/jtieri/HabbGo/habbgo/server/protocol/packets"
)

func HandleInitCrypto(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeCryptoParams())
}

func HandleGenerateKey(player *player.Player, packet *packets.IncomingPacket) {
	// TODO send
	player.Session.Send(composers.ComposeEndCrypto())
}

func HandleGetSessionParams(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeSessionParams())
}

func HandleSSO(player *player.Player, packet *packets.IncomingPacket) {

}

func HandleTryLogin(player *player.Player, packet *packets.IncomingPacket) {

}
