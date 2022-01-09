package commands

import (
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/protocol/messages"
	"github.com/jtieri/habbgo/protocol/packets"
)

func INIT_CRYPTO(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(player.Session.Address(), messages.CRYPTOPARAMETERS())
}

func GENERATEKEY(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(player.Session.Address(), messages.AVAILABLESETS())
	player.Session.Send(player.Session.Address(), messages.ENDCRYPTO())
	//player.Session.Send(composers.SECRETKEY())
}

func GET_SESSION_PARAMETERS(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(player.Session.Address(), messages.SESSIONPARAMETERS())
}

func VERSIONCHECK(player *player.Player, packet *packets.IncomingPacket) {

}

func UNIQUEID(player *player.Player, packet *packets.IncomingPacket) {

}

func SECRETKEY(player *player.Player, packets *packets.IncomingPacket) {
	player.Session.Send(player.Session.Address(), messages.ENDCRYPTO())
}

func SSO(p *player.Player, packet *packets.IncomingPacket) {
	token := packet.ReadString()

	// TODO if p login with token is success login, otherwise send LOCALISED ERROR & disconnect from server
	if token == "" {
		p.Login()
	} else {

	}
}

func TRY_LOGIN(p *player.Player, packet *packets.IncomingPacket) {
	username := packet.ReadString()
	password := packet.ReadString()

	if player.LoginDB(p, username, password) {
		p.Login()
		p.Session.Send(p.Session.Address(), messages.LOGINOK())
	} else {
		// TODO send LOCALISED ERROR
	}
}
