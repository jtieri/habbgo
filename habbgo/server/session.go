package server

import (
	"bufio"
	"bytes"
	"github.com/jtieri/HabbGo/habbgo/server/protocol/packets"
	"github.com/jtieri/HabbGo/habbgo/utils/encoding"
	"log"
	"net"
	"sync"
)

type Session struct {
	connection net.Conn
	buffer     *buffer
	active     bool
	server     *server
}

type buffer struct {
	mux  sync.Mutex
	buff *bufio.Writer
}

// Listen starts listening for incoming data from a Session and handles it appropriately.
func (session *Session) Listen() {
	// TODO create player and add to list of online players
	reader := bufio.NewReader(session.connection)

	// TODO Send HELLO packet to initialize connection

	// Main loop for listening for incoming packets from a players session
	for {
		// Attempt to read three bytes; client->server packets in FUSEv0.2.0 begin with 3 byte B64 encoded length.
		encodedLen := make([]byte, 3)
		for i := 0; i < 3; i++ {
			b, err := reader.ReadByte()

			if err != nil {
				// TODO handle errors parsing packets
				session.Close()
				return
			}
			encodedLen[i] = b
		}
		length := encoding.DecodeB64(encodedLen)

		// Check if data is junk before handling
		rawPacket := make([]byte, length)
		bytesRead, err := reader.Read(rawPacket)
		if length == 0 || err != nil || bytesRead < length {
			log.Println("Junk packet received.") // TODO handle logging junk packets
			continue
		}

		// Get Base64 encoded packet header
		payload := bytes.NewBuffer(rawPacket)
		rawHeader := make([]byte, 2)
		for i := 0; i < 2; i++ {
			rawHeader[i], _ = payload.ReadByte()
		}

		// Create a struct for the packet, and pass it on to be handled.
		packet := packets.NewIncoming(rawHeader, payload)
		log.Printf("Received packet [{%v} - %v] with contents: %v ", packet.Header, packet.HeaderId, packet.Payload.String())
		// TODO handle packets coming in from player's Session
	}
}

func (session *Session) Send(packet *packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	_, err := session.buffer.buff.Write(packet.Payload.Bytes())
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.connection.LocalAddr(), err)
	}

	err = session.buffer.buff.Flush()
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.connection.LocalAddr(), err)
	}
}

func (session *Session) Queue(packet *packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	_, err := session.buffer.buff.Write(packet.Payload.Bytes())
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.connection.LocalAddr(), err)
	}
}

func (session *Session) Flush(packet *packets.OutgoingPacket) {
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	err := session.buffer.buff.Flush()
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.connection.LocalAddr(), err)
	}
}

// Close disconnects a Session from the server.
func (session *Session) Close() {
	log.Printf("Closing session for address: %v ", session.connection.LocalAddr())
	session.server.RemoveSession(session)
	session.server = nil
	session.buffer = nil
	_ = session.connection.Close()
	session.active = false
}
