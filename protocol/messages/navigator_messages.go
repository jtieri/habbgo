package messages

import (
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

	p.WriteBool(hideFullRooms)
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
		// if room is public
		if r.Details.OwnerId == 0 {
			desc := r.Details.Description

			var door int
			if strings.Contains(desc, "/") {
				data := strings.Split(desc, "/")
				desc = data[0]
				door, _ = strconv.Atoi(data[1])
			}

			p.WriteInt(r.Details.ID + room.PublicRoomOffset)
			p.WriteInt(1)
			p.WriteString(r.Details.Name)
			p.WriteInt(r.Details.CurrentVisitors)
			p.WriteInt(r.Details.MaxVisitors)
			p.WriteInt(r.Details.CategoryID)
			p.WriteString(desc)
			p.WriteInt(r.Details.ID)
			p.WriteInt(door)
			p.WriteString(r.Details.CCTs)
			p.WriteInt(0)
			p.WriteInt(1)
		} else {
			p.WriteInt(r.Details.ID)
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
		p.WriteInt(subcat.ID)
		p.WriteInt(0)
		p.WriteString(subcat.Name)
		p.WriteInt(navigator.CurrentVisitors(&subcat))
		p.WriteInt(navigator.MaxVisitors(&subcat))
		p.WriteInt(parentCat.ID)
	}

	return p
}
