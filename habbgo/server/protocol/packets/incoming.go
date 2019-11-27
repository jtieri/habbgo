package packets

import (
	"bytes"
	"github.com/jtieri/HabbGo/habbgo/utils/encoding"
)

type IncomingPacket struct {
	Header string
	HeaderId uint
	Payload *bytes.Buffer
}

func (packet *IncomingPacket) ReadB64() int {
	data := make([]byte, 2)
	data[0], _ = packet.Payload.ReadByte()
	data[1], _ = packet.Payload.ReadByte()
	return encoding.DecodeB64(data)
}

func (packet *IncomingPacket) ReadBytes(i int) []byte {
	data := packet.Payload.Next(i)
	return data
}

func (packet *IncomingPacket) ReadInt() int {
	data := packet.Bytes()
	length := int(data[0] >> 3 & 7)
	value := encoding.DecodeVl64(data)
	packet.ReadBytes(length)
	return value
}

func (packet *IncomingPacket) ReadBool() bool {
	return packet.ReadInt() == 1
}

func (packet *IncomingPacket) Bytes() []byte {
	return packet.Payload.Bytes()
}

// VerifyIncoming verifies that an incoming packet is not some garbage or scripted value.
// It will return true if the packet is acceptable, and false otherwise.
func VerifyIncoming(rawPacket []byte) bool {
	if len(rawPacket) > 5 {
		buffer := bytes.Buffer{}
		buffer.Write(rawPacket)

		rawLen := make([]byte, 3)
		for i := 0; i < 3; i++ {
			rawLen[i], _ = buffer.ReadByte()
		}
		length := encoding.DecodeB64(rawLen)

		if length == 0 || buffer.Len() < length {
			return false
		}

		return true
	}

	return false
}