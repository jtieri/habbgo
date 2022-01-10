package messages

import (
	"fmt"
	"github.com/jtieri/habbgo/game/navigator"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/room"
	"github.com/jtieri/habbgo/protocol/packets"
	"strconv"
	"strings"
)

func NAVNODEINFO(player *player.Player, parentCat *navigator.Category, hideFullRooms bool, subcats []navigator.Category,
	rooms []*room.Room, currentVisitors int, maxVisitors int) *packets.OutgoingPacket {
	p := packets.NewOutgoing(220) // Base64 Header C\

	p.WriteBool(hideFullRooms) // hideCategory
	p.WriteInt(parentCat.ID)

	if parentCat.IsPublic {
		p.WriteInt(0)
	} else {
		p.WriteInt(2)
	}

	p.WriteString(parentCat.Name)
	p.WriteInt(currentVisitors)
	p.WriteInt(maxVisitors)
	p.WriteInt(parentCat.ParentID)

	if !parentCat.IsPublic {
		p.WriteInt(len(rooms))
	}

	for _, r := range rooms {
		if r.Details.OwnerId == 0 { // if r is public
			desc := r.Details.Description

			var door int
			if strings.Contains(desc, "/") {
				data := strings.Split(desc, "/")
				desc = data[0]
				door, _ = strconv.Atoi(data[1])
			}

			p.WriteInt(r.Details.Id + room.PublicRoomOffset) // writeInt roomId
			p.WriteInt(1)                                    // writeInt 1
			p.WriteString(r.Details.Name)                    // writeString roomName
			p.WriteInt(r.Details.CurrentVisitors)            // writeInt currentVisitors
			p.WriteInt(r.Details.MaxVisitors)                // writeInt maxVisitors
			p.WriteInt(r.Details.CategoryID)                 // writeInt catId
			p.WriteString(desc)                              // writeString roomDesc
			p.WriteInt(r.Details.Id)                         // writeInt roomId
			p.WriteInt(door)                                 // writeInt door
			p.WriteString(r.Details.CCTs)                    // writeString roomCCTs
			p.WriteInt(0)                                    // writeInt 0
			p.WriteInt(1)                                    // writeInt 1
		} else {
			p.WriteInt(r.Details.Id)
			p.WriteString(r.Details.Name)

			// TODO check that player is owner of r, that r is showing owner name, or that player has right SEE_ALL_ROOMOWNERS
			if player.Details.Username == r.Details.OwnerName || r.Details.ShowOwner {
				p.WriteString(r.Details.OwnerName)
			} else {
				p.WriteString("-")
			}

			p.WriteString(r.Details.AccessType.String())
			p.WriteInt(r.Details.CurrentVisitors)
			p.WriteInt(r.Details.MaxVisitors)
			p.WriteString(r.Details.Description)
		}
	}

	// iterate over sub-categories
	for _, subcat := range subcats {
		if subcat.MinRankAccess > player.Details.PlayerRank {
			continue
		}

		fmt.Println(subcat.Name)
		p.WriteInt(subcat.ID)
		p.WriteInt(0)
		p.WriteString(subcat.Name)
		p.WriteInt(navigator.CurrentVisitors(&subcat)) // writeInt currentVisitors
		p.WriteInt(navigator.MaxVisitors(&subcat))     // writeInt maxVisitors
		p.WriteInt(parentCat.ID)
	}

	return p
}
