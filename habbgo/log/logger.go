package log

import (
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
	"log"
)

func PrintOutgoingPacket(p *packets.OutgoingPacket) {

}

func PrintIncomingPacket(p *packets.IncomingPacket) {
	log.Printf("Received packet [%v - %v] with contents: %v ", p.Header, p.HeaderId, p.Payload.String())
}

func PrintUnkownPacket(p *packets.IncomingPacket) {

}
