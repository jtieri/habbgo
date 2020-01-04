package handlers

import (
	"github.com/jtieri/HabbGo/habbgo/game/model"
	"github.com/jtieri/HabbGo/habbgo/game/service"
	"github.com/jtieri/HabbGo/habbgo/protocol/composers"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
)

func HandleNavigate(player *model.Player, packet *packets.IncomingPacket) {
	nodeMask := packet.ReadBool()
	catId := packet.ReadInt()

	// get category using catID
	category := service.NavigatorService().CategoryById(catId)

	if category == nil {
		return
	}

	// if minrank for cat is > playerRank then return without sending response

	// get sub categories of category
	// sort categories by player count

	// get category currentvisitors
	// get category maxvisitors

	player.Session.Send(composers.ComposeNavNodeInfo(player, category, nodeMask))
}
