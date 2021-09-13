package handlers

import (
	"github.com/jtieri/HabbGo/habbgo/game/player"
	"github.com/jtieri/HabbGo/habbgo/protocol/composers"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
)

func InitCrypto(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeCryptoParams())
}

func GenerateKey(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeAvailableSets())
	player.Session.Send(composers.ComposeEndCrypto())
}

func GetSessionParams(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(composers.ComposeSessionParams())
}

func VersionCheck(player *player.Player, packet *packets.IncomingPacket) {

}

func UniqueID(player *player.Player, packet *packets.IncomingPacket) {

}

func SSO(p *player.Player, packet *packets.IncomingPacket) {
	token := packet.ReadString()

	// TODO if p login with token is success login, otherwise send LOCALISED ERROR & disconnect from server
	if token == "" {
		player.Login(p)
	} else {

	}
}

func TryLogin(p *player.Player, packet *packets.IncomingPacket) {
	username := packet.ReadString()
	password := packet.ReadString()

	if player.LoginDB(p, username, password) {
		player.Login(p)
		p.Session.Send(composers.ComposeLoginOk())
	} else {
		// TODO send LOCALISED ERROR
	}
}
