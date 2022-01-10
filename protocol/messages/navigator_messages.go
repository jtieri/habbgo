package messages

import (
	"github.com/jtieri/habbgo/game/navigator"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/room"
	"github.com/jtieri/habbgo/protocol/packets"
	"strconv"
	"strings"
)

func NAVNODEINFO(player *player.Player, cat *navigator.Category, nodeMask bool, subcats []*navigator.Category,
	rooms []*room.Room, currentVisitors int, maxVisitors int) *packets.OutgoingPacket {
	p := packets.NewOutgoing(220) // Base64 Header C\

	p.WriteBool(nodeMask) // hideCategory
	p.WriteInt(cat.ID)

	if cat.IsPublic {
		p.WriteInt(0)
	} else {
		p.WriteInt(2)
	}

	p.WriteString(cat.Name)
	p.WriteInt(currentVisitors)
	p.WriteInt(maxVisitors)
	p.WriteInt(cat.ParentID)

	if !cat.IsPublic {
		p.WriteInt(len(rooms))
	}

	for _, r := range rooms {
		if r.Details.Owner_Id == 0 { // if r is public
			desc := r.Details.Desc

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
			p.WriteInt(r.Details.CatId)                      // writeInt catId
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
			if player.Details.Username == r.Details.Owner_Name {
				p.WriteString(r.Details.Owner_Name)
			} else {
				p.WriteString("-")
			}

			p.WriteString(room.AccessType(r.Details.AccessType))
			p.WriteInt(r.Details.CurrentVisitors)
			p.WriteInt(r.Details.MaxVisitors)
			p.WriteString(r.Details.Desc)
		}
	}

	// iterate over sub-categories
	for _, subcat := range subcats {
		if subcat.MinRankAccess > 1 {
			continue
		}

		p.WriteInt(subcat.ID)
		p.WriteInt(0)
		p.WriteString(subcat.Name)
		p.WriteInt(navigator.CurrentVisitors(subcat)) // writeInt currentVisitors
		p.WriteInt(navigator.MaxVisitors(subcat))     // writeInt maxVisitors
		p.WriteInt(cat.ID)
	}

	return p
}
