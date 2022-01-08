package server

import (
	"github.com/jtieri/habbgo/app"
	"github.com/jtieri/habbgo/game/player"
	logger "github.com/jtieri/habbgo/log"
	"github.com/jtieri/habbgo/protocol/packets"
)

func Handle(p *player.Player, packet *packets.IncomingPacket) {
	handler, found := p.Session.GetPacketHandler(packet.HeaderId)

	if found {
		if app.Habbgo().Config.Server.Debug {
			logger.LogIncomingPacket(p.Session.Address(), handler, packet)
		}
		handler(p, packet)
	} else {
		if app.Habbgo().Config.Server.Debug {
			logger.LogUnknownPacket(p.Session.Address(), packet)
		}
	}

}
