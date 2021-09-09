package handlers

import (
	"github.com/jtieri/HabbGo/habbgo/game/player"
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

func HandleSSO(p *player.Player, packet *packets.IncomingPacket) {
	token := packet.ReadString()

	// TODO if p login with token is success login, otherwise send LOCALISED ERROR & disconnect from server
	if token == "" {
		player.Login(p)
	} else {

	}
}

func HandleTryLogin(p *player.Player, packet *packets.IncomingPacket) {
	username := packet.ReadString()
	password := packet.ReadString()

	if player.LoginDB(p, username, password) {
		player.Login(p)
		p.Session.Send(composers.ComposeLoginOk())
	} else {
		// TODO send LOCALISED ERROR
	}
}
