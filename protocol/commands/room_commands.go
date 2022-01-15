package commands

import (
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/room"
	"github.com/jtieri/habbgo/protocol/messages"
	"github.com/jtieri/habbgo/protocol/packets"
)

func GETINTERST(player *player.Player, packet *packets.IncomingPacket) {
	player.Session.Send(player.Details.Username, messages.INTERSTITIALDATA, messages.INTERSTITIALDATA())
}

// packet payload [#boolean:tTypeID, #integer:tRoomID, #integer:tDoorID]
func ROOM_DIRECTORY(player *player.Player, packet *packets.IncomingPacket) {
	isPublic := packet.ReadBool()
	if !isPublic {
		// send OPEN_CONNECTION
		return
	}

	roomID := packet.ReadInt()

	r := room.RoomService().RoomByID(roomID)

	if r == nil {
		return
	}

	// if room is habbo club only and player is not habbo club send CANTCONNECT with error CLUB_ONLY

	// if room is full and player does not have the fuse right for entering full rooms send CANTCONNECT with error FULL

	// enter room
}
