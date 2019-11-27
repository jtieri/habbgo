package server

import (
	"bufio"
	"github.com/jtieri/HabbGo/habbgo/server/protocol/packets"
	"log"
	"net"
)

type Session struct {
	connection net.Conn
	buffer bufio.Writer
	server     *server
}

// Listen starts listening for incoming data from a Session and handles it appropriately.
func (session *Session) Listen() {
	reader := bufio.NewReader(session.connection)

	// Main loop for listening for incoming packets from a players session
	for {
		var rawPacket []byte
		_, err := reader.Read(rawPacket)
		if err != nil {
			// TODO handle errors reading packets
		}
		log.Println(rawPacket) // TODO REMOVE THIS DUMB TEST

		// Decode raw rawPacket & if acceptable
		if packets.VerifyIncoming(rawPacket) {
			// Create an incoming packet struct with remaining bytes in rawPacket
			// TODO handle packets coming in from player's Sessions
		}
	}
}

func (session *Session) Send() {
	// Write
	// Flush
}

func (session *Session) Queue() {
	// Write
}

func (session *Session) Flush() {
	// Flush
}

// Close disconnects a Session from the server.
func (session *Session) Close() {
	session.connection.Close()
}
