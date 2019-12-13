package handlers

import (
	"github.com/jtieri/HabbGo/habbgo/database"
	"github.com/jtieri/HabbGo/habbgo/game/model/player"
	"github.com/jtieri/HabbGo/habbgo/protocol/composers"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
)

func HandleInitCrypto(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeCryptoParams())
}

func HandleGenerateKey(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeEndCrypto())
}

func HandleGetSessionParams(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeSessionParams())
}

func HandleSSO(player *player.Player, packet *packets.IncomingPacket) {
	token := packet.ReadString()

	// TODO if player login with token is success login, otherwise send LOCALISED ERROR & disconnect from server
	if token == "" {
		player.Service.Login()
	} else {

	}
}

func HandleTryLogin(player *player.Player, packet *packets.IncomingPacket) {
	username := packet.ReadString()
	password := packet.ReadString()

	if database.Login(player, username, password) {
		player.Service.Login()
	} else {
		// TODO send LOCALISED ERROR
	}
}
