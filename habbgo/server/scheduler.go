package server

import (
	"github.com/jtieri/HabbGo/habbgo/app"
	"github.com/jtieri/HabbGo/habbgo/game/player"
	logger "github.com/jtieri/HabbGo/habbgo/log"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
)

func Handle(p *player.Player, packet *packets.IncomingPacket) {
	handler, found := p.Session.GetPacketHandler(packet.HeaderId)

	if found {
		if app.HabbGo().Config.Server.Debug {
			logger.LogIncomingPacket(p.Session.Address(), handler, packet)
		}
		handler(p, packet)
	} else {
		if app.HabbGo().Config.Server.Debug {
			logger.LogUnknownPacket(p.Session.Address(), packet)
		}
	}

}
