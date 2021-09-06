package composers

import (
	"github.com/jtieri/HabbGo/habbgo/game/player"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
	"strconv"
)

func ComposeUserObj(player *player.Player) *packets.OutgoingPacket {
	p := packets.NewOutgoing(5) // Base64 Header @E

	p.WriteString(strconv.Itoa(player.Details.Id))
	p.WriteString(player.Details.Username)
	p.WriteString(player.Details.Figure)
	p.WriteString(player.Details.Sex)
	p.WriteString(player.Details.Motto)
	p.WriteInt(player.Details.Tickets)
	p.WriteString(player.Details.PoolFigure)
	p.WriteInt(player.Details.Film)
	//p.WriteInt(directMail)

	return p
}

func ComposeCreditBalance(credits int) *packets.OutgoingPacket {
	p := packets.NewOutgoing(6) // Base64 Header @F
	p.WriteString(strconv.Itoa(credits) + ".0")
	return p
}

func ComposeAvailableBadges(player *player.Player) *packets.OutgoingPacket {
	p := packets.NewOutgoing(229) // Base64 Header

	p.WriteInt(len(player.Details.Badges))

	var bSlot int
	for i, b := range player.Details.Badges {
		p.WriteString(b)

		if b == player.Details.CurrentBadge {
			bSlot = i
		}
	}

	p.WriteInt(bSlot)
	p.WriteBool(player.Details.DisplayBadge)
	return p
}

func ComposeSoundSetting(ss int) *packets.OutgoingPacket {
	p := packets.NewOutgoing(308) // Base 64 Header Dt
	p.WriteInt(ss)
	return p
}

func ComposeLatency(l int) *packets.OutgoingPacket {
	p := packets.NewOutgoing(354) // Base 64 Header Eb
	p.WriteInt(l)
	return p
}
