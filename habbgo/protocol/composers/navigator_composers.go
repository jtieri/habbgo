package composers

import "github.com/jtieri/HabbGo/habbgo/protocol/packets"

func ComposeNavNodeInfo() *packets.OutgoingPacket {
	p := packets.NewOutgoing(220) // Base64 Header C\

	return p
}
