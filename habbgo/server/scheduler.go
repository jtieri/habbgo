package server

import (
	"github.com/jtieri/HabbGo/habbgo/game/player"
	logger "github.com/jtieri/HabbGo/habbgo/log"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
)

func Handle(p *player.Player, packet *packets.IncomingPacket) {
	handler, found := p.Session.GetPacketHandler(packet.HeaderId)

	if found {
		if GetConfig().Server.Debug {
			logger.LogIncomingPacket(p.Session.Address(), handler, packet)
		}
		handler(p, packet)
	} else {
		if GetConfig().Server.Debug {
			logger.LogUnknownPacket(p.Session.Address(), packet)
		}
	}

}
