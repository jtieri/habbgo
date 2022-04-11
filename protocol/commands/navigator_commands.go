package commands

import (
	"github.com/jtieri/habbgo/game/navigator"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/room"
	"github.com/jtieri/habbgo/protocol/messages"
	"github.com/jtieri/habbgo/protocol/packets"
)

func Navigate(player *player.Player, packet *packets.IncomingPacket) {
	roomService := player.Services.RoomService()

	hideFullRooms := packet.ReadInt() == 1
	catId := packet.ReadInt()

	if catId >= room.PublicRoomOffset {
		r := roomService.RoomById(catId - room.PublicRoomOffset)
		if r != nil {
			catId = r.Details.CategoryID
		}
	}

	category := player.Services.NavigatorService().CategoryById(catId)

	// TODO also check that access rank isnt higher than players rank
	if category == nil {
		return
	}

	subCategories := player.Services.NavigatorService().CategoriesByParentId(category.ID)
	// sort categories by player count

	r := player.Services.RoomService().Rooms()
	currentVisitors := navigator.CurrentVisitors(category, r)
	maxVisitors := navigator.MaxVisitors(category, r)

	var rooms []*room.Room
	if category.IsPublic {
		for _, r := range roomService.ReplaceRooms(roomService.RoomsByPlayerId(0)) {
			if r.Details.CategoryID == category.ID && (!hideFullRooms) && r.Details.CurrentVisitors < r.Details.MaxVisitors {
				// if room is hidden or category id is not equal to the category id we are working with currently continue
				if r.Details.Hidden || r.Details.CategoryID != category.ID {
					continue
				}

				// if we are hiding full rooms in the navigator and the room is full continue
				if hideFullRooms && r.Details.CurrentVisitors >= r.Details.MaxVisitors {
					continue
				}

				rooms = append(rooms, r)
			}
		}
	} else {
		// TODO finish private room logic
	}

	//// ----------
	//fmt.Println("--------------------")
	//fmt.Println(category.Name)
	//for _, c := range subCategories {
	//	fmt.Println(c.Name)
	//}
	//fmt.Println("--------------------")
	//for _, r := range rooms {
	//	fmt.Println(r.Details.Name)
	//}

	// TODO sort rooms by player count before sending NAVNODEINFO

	player.Session.Send(messages.NAVNODEINFO, messages.NAVNODEINFO(player, category, hideFullRooms, subCategories, rooms, currentVisitors, maxVisitors))
}
