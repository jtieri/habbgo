package server

import (
	"bufio"
	"github.com/jtieri/HabbGo/habbgo/server/protocol/packets"
	"github.com/jtieri/HabbGo/habbgo/utils/encoding"
	"log"
	"net"
)

type Session struct {
	connection net.Conn
	buffer     bufio.Writer
	server     *server
}

// Listen starts listening for incoming data from a Session and handles it appropriately.
func (session *Session) Listen() {
	// TODO Create Player

	reader := bufio.NewReader(session.connection)

	// TODO Send HELLO packet to initialize connection

	// Main loop for listening for incoming packets from a players session
	for {
		// Attempt to read three bytes; client->server packets in FUSEv0.2.0 begin with 3 byte B64 encoded length.
		encodedLen := make([]byte, 3)
		for i := 0; i < 3; i++ {
			b, err := reader.ReadByte()

			if err != nil {
				// TODO handle errors reading packets
				session.Close()
				return
			}
			encodedLen[i] = b
		}
		length := encoding.DecodeB64(encodedLen)

		// Get Base64 encoded packet header
		rawHeader := make([]byte, 2)
		for i := 0; i < 2; i++ {
			rawHeader[i], _ = reader.ReadByte()
		}

		// Check if data is junk before handling
		var rawPacket []byte
		bytesRead, err := reader.Read(rawPacket)
		if length == 0 || err != nil || bytesRead < length {
			// TODO handle logging junk packets
			continue
		}

		// Create a struct for the packet, and pass it on to be handled.
		packet := packets.NewIncoming(rawHeader, rawPacket)
		log.Printf("Received packet %v with contents: %v ", packet.HeaderId, packet.Payload.String()) // TODO REMOVE THIS DUMB TEST
		// TODO handle packets coming in from player's Session
	}
}

func (session *Session) Send(packet *packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.Write(packet.Payload.Bytes())
	session.buffer.Flush()
}

func (session *Session) Queue(packet *packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.Write(packet.Payload.Bytes())
}

func (session *Session) Flush(packet *packets.OutgoingPacket) {
	session.buffer.Flush()
}

// Close disconnects a Session from the server.
func (session *Session) Close() {
	session.server.RemoveSession(session)
	session.server = nil
	session.connection.Close()
}
