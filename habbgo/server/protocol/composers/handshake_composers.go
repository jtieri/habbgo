package composers

import "github.com/jtieri/HabbGo/habbgo/server/protocol/packets"

func ComposeHello() *packets.OutgoingPacket {
	return packets.NewOutgoing(0) // Base64 Header @@
}

func ComposeCryptoParams() *packets.OutgoingPacket {
	packet := packets.NewOutgoing(277) // Base64 Header DU
	packet.WriteInt(0) // Toggles server->client encryption; 0=off | non-zero=on
	return packet
}

func ComposeEndCrypto() *packets.OutgoingPacket {
	packet := packets.NewOutgoing(278) // Base 64 Header DV
	return packet
}

func ComposeSessionParams() *packets.OutgoingPacket {
	packet := packets.NewOutgoing(257) // Base64 Header DA

	// Map the param types to their values
	// writeInt number of params
	// iterate over the params & writeInt(paramiD) then writeInt for numeric values and writeString for strings

	return packet
}