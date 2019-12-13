package composers

import (
	"github.com/jtieri/HabbGo/habbgo/game/model/player"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
	"strconv"
)

func ComposeUserObj(player *player.Player) *packets.OutgoingPacket {
	p := packets.NewOutgoing(5) // Base64 Header @E

	p.WriteString(strconv.Itoa(player.Details.Id)) // writeString userId
	p.WriteString(player.Details.Username)             // writeString name
	p.WriteString(player.Details.Figure)               // writeString figure
	p.WriteString(string(player.Details.Sex))          // writeString sex
	p.WriteString(player.Details.Motto)                // writeString motto
	p.WriteInt(player.Details.Tickets)                 // writeInt ph_tickets
	p.WriteString(player.Details.PoolFigure)           // writeString ph_figure
	p.WriteInt(player.Details.Film)                    // writeInt photo_film
	//p.WriteInt(directMail)

	return p
}

func ComposeCreditBalance(credits int) *packets.OutgoingPacket {
	p := packets.NewOutgoing(6) // Base64 Header @F
	p.WriteString(strconv.Itoa(credits) + ".0")
	return p
}

func ComposeAvailableBadges() *packets.OutgoingPacket {
	p := packets.NewOutgoing(229) // Base64 Header

	// writeInt num of badges

	// loop and writeString each badge id

	// writeInt chosenBadge
	// writeInt visible
	return p
}
