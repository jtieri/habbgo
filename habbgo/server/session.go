package server

import (
	"bufio"
	"bytes"
	"database/sql"
	"github.com/jtieri/HabbGo/habbgo/game/model/player"
	"github.com/jtieri/HabbGo/habbgo/game/service"
	"github.com/jtieri/HabbGo/habbgo/protocol/composers"
	"github.com/jtieri/HabbGo/habbgo/protocol/encoding"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
	"log"
	"net"
	"sync"
)

type Session struct {
	Connection net.Conn
	database   *sql.DB
	buffer     *buffer
	active     bool
	server     *Server
}

type buffer struct {
	mux  sync.Mutex
	buff *bufio.Writer
}

func NewSession(conn net.Conn, server *Server) *Session {
	s := &Session{
		Connection: conn,
		database:   server.Database,
		buffer:     &buffer{mux: sync.Mutex{}, buff: bufio.NewWriter(conn)},
		active:     true,
		server:     server,
	}
	return s
}

// Listen starts listening for incoming data from a Session and handles it appropriately.
func (session *Session) Listen() {
	p := player.New(session, service.New())
	p.Service.Prepare(p)
	reader := bufio.NewReader(session.Connection)

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

		if session.server.Config.Log.Incoming {
			log.Printf("Received packet [%v - %v] with contents: %v ",
				packet.Header, packet.HeaderId, packet.Payload.String())
		}

		go Handle(p, packet) // Handle packets coming in from p's Session
	}
}

func (session *Session) Send(packet *packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	_, err := session.buffer.buff.Write(packet.Payload.Bytes())
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.Connection.LocalAddr(), err)
	}

	err = session.buffer.buff.Flush()
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.Connection.LocalAddr(), err)
	}

	if session.server.Config.Log.Outgoing {
		log.Printf("Sent packet [%v - %v] with contents: %v ", packet.Header, packet.HeaderId, packet.String())
	}
}

func (session *Session) Queue(packet *packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	_, err := session.buffer.buff.Write(packet.Payload.Bytes())
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.Connection.LocalAddr(), err)
	}
}

func (session *Session) Flush(packet *packets.OutgoingPacket) {
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	err := session.buffer.buff.Flush()
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.Connection.LocalAddr(), err)
	}

	if session.server.Config.Log.Outgoing {
		log.Printf("Sent packet [%v - %v] with contents: %v ", packet.Header, packet.HeaderId, packet.String())
	}
}

func (session *Session) Database() *sql.DB {
	return session.database
}

// Close disconnects a Session from the server.
func (session *Session) Close() {
	log.Printf("Closing session for address: %v ", session.Connection.LocalAddr())
	session.server.RemoveSession(session)
	session.server = nil
	session.buffer = nil
	_ = session.Connection.Close()
	session.active = false
}
