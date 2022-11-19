package commands

import (
	"github.com/jtieri/habbgo/game/navigator"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/room"
	"github.com/jtieri/habbgo/protocol/messages"
	"github.com/jtieri/habbgo/protocol/packets"
)

func NAVIGATE(player player.Player, packet packets.IncomingPacket) {
	roomService := player.Services.RoomService().(*room.RoomServiceProxy)

	navService := player.Services.NavigatorService().(*navigator.NavigatorServiceProxy)

	hideFullRooms := packet.ReadInt() == 1
	catId := packet.ReadInt()

	if catId >= room.PublicRoomOffset {
		r := roomService.RoomByID(catId - room.PublicRoomOffset)
		if r.Ready {
			catId = r.Details.CategoryID
		}
	}

	category := navService.CategoryByID(catId)
	if (category == navigator.Category{}) {
		return
	}

	// TODO also check that access rank isnt higher than players rank

	subCategories := navService.CategoriesByParentID(category.ID)

	r := roomService.Rooms()
	currentVisitors := room.CurrentVisitors(category, r)
	maxVisitors := room.MaxVisitors(category, r)

	var rooms []room.Room
	if category.IsPublic {
		playerRooms := roomService.RoomsByPlayerID(0)
		if playerRooms == nil {
			return
		}

		replacedRooms := roomService.CheckRoomsQueried(playerRooms)
		for _, r := range replacedRooms {
			if r.Details.Hidden {
				continue
			}

			if r.Details.CategoryID != category.ID {
				continue
			}

			if hideFullRooms && r.Details.CurrentVisitors >= r.Details.MaxVisitors {
				continue
			}

			rooms = append(rooms, r)
		}
	}

	// TODO finish private room logic

	// TODO sort rooms by player count before sending NAVNODEINFO

	player.Session.Send(messages.NAVNODEINFO, messages.NAVNODEINFO(player, category, hideFullRooms, subCategories, rooms, currentVisitors, maxVisitors))
}

// GETUSERFLATCATS is sent from the client requesting the navigator.Navigator private room categories that
// should be visible for the specified user.
func GETUSERFLATCATS(player player.Player, packet packets.IncomingPacket) {
	var privateRoomCategories []navigator.Category

	navService := player.Services.NavigatorService().(*navigator.NavigatorServiceProxy)

	// We only want to send category information for private rooms that
	// should be visible by the player, so don't add categories that are
	// set with a minimum rank to access that is greater than the players rank.
	for _, category := range navService.Categories() {
		if category.IsPublic && player.Details.PlayerRank < category.MinRankAccess {
			continue
		}
		privateRoomCategories = append(privateRoomCategories, category)
	}

	player.Session.Send(messages.USERFLATCATS, messages.USERFLATCATS(privateRoomCategories))
}
