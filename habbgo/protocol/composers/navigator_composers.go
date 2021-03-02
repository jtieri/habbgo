package composers

import (
	"github.com/jtieri/HabbGo/habbgo/game/model"
	"github.com/jtieri/HabbGo/habbgo/game/service"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
	"strconv"
	"strings"
)

func ComposeNavNodeInfo(player *model.Player, cat *model.Category, nodeMask bool, subcats []*model.Category, rooms []*model.Room, currentVisitors int, maxVisitors int) *packets.OutgoingPacket {
	p := packets.NewOutgoing(220) // Base64 Header C\

	p.WriteBool(nodeMask)   // hideCategory
	p.WriteInt(cat.Id)

	if cat.Public {
		p.WriteInt(0)
	} else {
		p.WriteInt(2)
	}

	p.WriteString(cat.Name)
	p.WriteInt(currentVisitors)
	p.WriteInt(maxVisitors)
	p.WriteInt(cat.Pid)

	if !cat.Public {
		p.WriteInt(len(rooms))
	}

	for _, room := range rooms {
		if room.Details.Owner_Id == 0 { // if room is public
			desc := room.Details.Desc

			var door int
			if strings.Contains(desc, "/") {
				data := strings.Split(desc, "/")
				desc = data[0]
				door, _ = strconv.Atoi(data[1])
			}

			p.WriteInt(room.Details.Id + service.PublicRoomOffset)// writeInt roomId
			p.WriteInt(1) // writeInt 1
			p.WriteString(room.Details.Name)// writeString roomName
			p.WriteInt(room.Details.CurrentVisitors) // writeInt currentVisitors
			p.WriteInt(room.Details.MaxVisitors) // writeInt maxVisitors
			p.WriteInt(room.Details.CatId) // writeInt catId
			p.WriteString(desc)// writeString roomDesc
			p.WriteInt(room.Details.Id) // writeInt roomId
			p.WriteInt(door) // writeInt door
			p.WriteString(room.Details.CCTs)// writeString roomCCTs
			p.WriteInt(0) // writeInt 0
			p.WriteInt(1) // writeInt 1
		} else {
			p.WriteInt(room.Details.Id)
			p.WriteString(room.Details.Name)

			// TODO check that player is owner of room, that room is showing owner name, or that player has right SEE_ALL_ROOMOWNERS
			if player.Details.Username == room.Details.Owner_Name {
				p.WriteString(room.Details.Owner_Name)
			} else {
				p.WriteString("-")
			}

			p.WriteString(service.AccessType(room.Details.AccessType))
			p.WriteInt(room.Details.CurrentVisitors)
			p.WriteInt(room.Details.MaxVisitors)
			p.WriteString(room.Details.Desc)
		}
	}

	// iterate over sub-categories
	for _, subcat := range subcats {
		if subcat.MinRankAccess > 1 {
			continue
		}

		p.WriteInt(subcat.Id)
		p.WriteInt(0)
		p.WriteString(subcat.Name)
		p.WriteInt(service.CurrentVisitors(subcat))// writeInt currentVisitors
		p.WriteInt(service.MaxVisitors(subcat))// writeInt maxVisitors
		p.WriteInt(cat.Id)
	}

	return p
}
