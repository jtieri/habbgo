package handlers

import (
	"github.com/jtieri/HabbGo/habbgo/database"
	"github.com/jtieri/HabbGo/habbgo/game/model"
	"github.com/jtieri/HabbGo/habbgo/game/service"
	"github.com/jtieri/HabbGo/habbgo/protocol/composers"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
)

func HandleInitCrypto(player *model.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeCryptoParams())
}

func HandleGenerateKey(player *model.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeEndCrypto())
}

func HandleGetSessionParams(player *model.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeSessionParams())
}

func HandleSSO(player *model.Player, packet *packets.IncomingPacket) {
	token := packet.ReadString()

	// TODO if player login with token is success login, otherwise send LOCALISED ERROR & disconnect from server
	if token == "" {
		service.Login(player)
	} else {

	}
}

func HandleTryLogin(player *model.Player, packet *packets.IncomingPacket) {
	username := packet.ReadString()
	password := packet.ReadString()

	if database.Login(player, username, password) {
		service.Login(player)
	} else {
		// TODO send LOCALISED ERROR
	}
}
