package handlers

import (
	"github.com/jtieri/HabbGo/habbgo/game/navigator"
	"github.com/jtieri/HabbGo/habbgo/game/player"
	"github.com/jtieri/HabbGo/habbgo/game/room"
	"github.com/jtieri/HabbGo/habbgo/protocol/composers"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
)

func Navigate(player *player.Player, packet *packets.IncomingPacket) {
	roomService := room.RoomService()

	nodeMask := packet.ReadInt() == 1
	catId := packet.ReadInt()

	if catId >= room.PublicRoomOffset {
		r := roomService.RoomById(catId - room.PublicRoomOffset)
		if r != nil {
			catId = r.Details.CatId
		}
	}

	category := navigator.NavigatorService().CategoryById(catId)

	// TODO also check that access rank isnt higher than players rank
	if category == nil {
		return
	}

	subCategories := navigator.NavigatorService().CategoriesByParentId(category.Id)
	// sort categories by player count

	currentVisitors := navigator.CurrentVisitors(category)
	maxVisitors := navigator.MaxVisitors(category)

	var rooms []*room.Room
	if category.Public {
		for _, room := range roomService.ReplaceRooms(roomService.RoomsByPlayerId(0)) {
			if room.Details.CatId == category.Id && (!nodeMask) && room.Details.CurrentVisitors < room.Details.MaxVisitors {
				rooms = append(rooms, room)
			}
		}
	} else {
		// TODO finish private room logic
	}

	// TODO sort rooms by player count before sending NavNodeInfo

	player.Session.Send(composers.ComposeNavNodeInfo(player, category, nodeMask, subCategories, rooms, currentVisitors, maxVisitors))
}
