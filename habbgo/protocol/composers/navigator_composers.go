package composers

import (
	"github.com/jtieri/HabbGo/habbgo/game/model"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
)

func ComposeNavNodeInfo(player *model.Player, cat *model.Category, nodeMask bool) *packets.OutgoingPacket {
	p := packets.NewOutgoing(220) // Base64 Header C\

	// player
	// category
	// rooms
	// nodemask
	// subCategories
	// currentVisitors
	// maxVisitors
	// accessRank

	p.WriteBool(nodeMask)   // writeBool NodeMask
	p.WriteInt(cat.Id)      // writeInt CategoryId
	p.WriteBool(cat.Public) // writeBool publicSpace
	p.WriteString(cat.Name) // writeString CategoryName
	// writeInt categoryCurrentVisitors
	// writeInt categoryMaxVisitors
	// writeInt parentCategoryId

	// if category is for public rooms writeInt numberOfRooms
	if cat.Public {

	}

	// iterate over rooms
	// if room is public
	// writeInt roomId
	// writeInt 1
	// writeString roomName
	// writeInt currentVisitors
	// writeInt maxVisitors
	// writeInt catId
	// writeString roomDesc
	// writeInt roomId
	// writeInt door
	// writeString roomCCTs
	// writeInt 0
	// writeInt 1

	// if room is not private
	// writeInt roomId
	// writeString roomName
	// writeString ownerName or -
	// writeString accessType
	// writeInt currentVisitors
	// writeInt maxVisitors
	// writeString roomDesc

	// iterate over sub-categories
	// if user can access sub-cats
	// writeInt subCategoryId
	// writeInt 0
	// writeString subCatName
	// writeInt currentVisitors
	// writeInt maxVisitors
	// writeInt categoryId

	return p
}
