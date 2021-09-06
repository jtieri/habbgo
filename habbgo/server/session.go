package server

import (
	"bufio"
	"bytes"
	"github.com/jtieri/HabbGo/habbgo/game/player"
	"github.com/jtieri/HabbGo/habbgo/protocol/composers"
	"github.com/jtieri/HabbGo/habbgo/protocol/encoding"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
	"gorm.io/gorm"
	"log"
	"net"
	"sync"
)

type Session struct {
	connection net.Conn
	database   *gorm.DB
	buffer     *buffer
	active     bool
	server     *Server
}

type buffer struct {
	mux  sync.Mutex
	buff *bufio.Writer
}

// NewSession returns a pointer to a newly allocated Session struct, representing a players connection to the server.
func NewSession(conn net.Conn, server *Server) *Session {
	return &Session{
		connection: conn,
		database:   server.Database,
		buffer:     &buffer{mux: sync.Mutex{}, buff: bufio.NewWriter(conn)},
		active:     true,
		server:     server,
	}
}

// Listen starts listening for incoming data from a Session and handles it appropriately.
func (session *Session) Listen() {
	p := player.New(session)
	reader := bufio.NewReader(session.connection)

	session.Send(composers.ComposeHello()) // Send packet with Base64 header @@ to initialize connection with client.

	// Listen for incoming packets from a players session
	for {
		// Attempt to read three bytes; client->server packets in FUSEv0.2.0 begin with 3 byte Base64 encoded length.
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

		packet := packets.NewIncoming(rawHeader, payload)

		if session.server.Config.Server.Debug {
			log.Printf("Received packet [%v - %v] with contents: %v ",
				packet.Header, packet.HeaderId, packet.Payload.String())
		}

		go Handle(p, packet) // Handle packets coming in from p's Session
	}
}

// Send finalizes an outgoing packet with 0x01 and then attempts to write and flush the packet to a Session's buffer.
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

	if session.server.Config.Server.Debug {
		log.Printf("Sent packet [%v - %v] with contents: %v ", packet.Header, packet.HeaderId, packet.String())
	}
}

// Send finalizes an outgoing packet with 0x01 and then attempts to write the packet to a Session's buffer.
func (session *Session) Queue(packet *packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	_, err := session.buffer.buff.Write(packet.Payload.Bytes())
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.connection.LocalAddr(), err)
	}
}

// Send finalizes an outgoing packet with 0x01 and then attempts flush the packet to a Sessions's buffer.
func (session *Session) Flush(packet *packets.OutgoingPacket) {
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	err := session.buffer.buff.Flush()
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.connection.LocalAddr(), err)
	}

	if session.server.Config.Server.Debug {
		log.Printf("Sent packet [%v - %v] with contents: %v ", packet.Header, packet.HeaderId, packet.String())
	}
}

// Database returns a pointer to a Session's Conn access struct.
func (session *Session) Database() *gorm.DB {
	return session.database
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
