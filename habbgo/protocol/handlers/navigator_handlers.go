package handlers

import (
	"github.com/jtieri/HabbGo/habbgo/game/model"
	"github.com/jtieri/HabbGo/habbgo/game/service"
	"github.com/jtieri/HabbGo/habbgo/protocol/composers"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
)

func HandleNavigate(player *model.Player, packet *packets.IncomingPacket) {
	roomService := service.RoomService()

	nodeMask := packet.ReadInt()==1
	catId := packet.ReadInt()

	if catId >= service.PublicRoomOffset {
		room := roomService.RoomById(catId - service.PublicRoomOffset)
		if room != nil {
			catId = room.Details.CatId
		}
	}

	category := service.NavigatorService().CategoryById(catId)

	// TODO also check that access rank isnt higher than players rank
	if category == nil {
		return
	}

	subCategories := service.NavigatorService().CategoriesByParentId(category.Id)
	// sort categories by player count

	currentVisitors := service.CurrentVisitors(category)
	maxVisitors := service.MaxVisitors(category)

	var rooms []*model.Room
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
