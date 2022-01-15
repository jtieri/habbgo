package commands

import (
	"github.com/jtieri/habbgo/game/navigator"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/room"
	"github.com/jtieri/habbgo/protocol/messages"
	"github.com/jtieri/habbgo/protocol/packets"
)

func NAVIGATE(player *player.Player, packet *packets.IncomingPacket) {
	roomService := room.RoomService()

	hideFullRooms := packet.ReadInt() == 1
	categoryID := packet.ReadInt()

	if categoryID >= room.PublicRoomOffset {
		r := roomService.RoomByID(categoryID - room.PublicRoomOffset)
		if r != nil {
			categoryID = r.Details.CategoryID
		}
	}

	category := navigator.NavigatorService().CategoryById(categoryID)

	// TODO also check that access rank isnt higher than players rank
	if category == nil || category.MinRankAccess > player.Details.PlayerRank {
		return
	}

	subCategories := navigator.NavigatorService().CategoriesByParentId(category.ID)
	// TODO sort categories by player count

	var rooms []*room.Room
	if category.IsPublic {
		for _, r := range roomService.ReplaceRooms(roomService.RoomsByPlayerID(0)) {
			// if room is hidden or category id is not equal to the category id we are working with currently skip
			if r.Details.Hidden || r.Details.CategoryID != category.ID {
				continue
			}

			// if we are hiding full rooms in the navigator and the room is full skip
			if hideFullRooms && r.Details.CurrentVisitors >= r.Details.MaxVisitors {
				continue
			}

			rooms = append(rooms, r)
		}
	} else {
		// TODO finish private room logic
	}

	// TODO sort rooms by player count before sending NAVNODEINFO
	player.Session.Send(player.Details.Username, messages.NAVNODEINFO, messages.NAVNODEINFO(player, category,
		hideFullRooms, subCategories, rooms, navigator.CurrentVisitors(category), navigator.MaxVisitors(category)))
}
