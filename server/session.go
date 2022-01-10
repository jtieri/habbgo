package server

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/jtieri/habbgo/game/player"
	logger "github.com/jtieri/habbgo/log"
	"github.com/jtieri/habbgo/protocol/encoding"
	"github.com/jtieri/habbgo/protocol/messages"
	"github.com/jtieri/habbgo/protocol/packets"
)

// Session represents a players underlying network session and connection to the server
type Session struct {
	connection net.Conn
	buffer     *buffer
	active     bool
	server     *Server
	router     *Router
}

type buffer struct {
	mux  sync.Mutex
	buff *bufio.Writer
}

// NewSession returns a pointer to a newly allocated Session struct, representing a players connection to the server.
func NewSession(conn net.Conn, server *Server) *Session {
	return &Session{
		connection: conn,
		buffer:     &buffer{mux: sync.Mutex{}, buff: bufio.NewWriter(conn)},
		active:     true,
		server:     server,
		router:     RegisterCommands(),
	}
}

// Listen starts listening for incoming data from a Session and handles it appropriately.
func (session *Session) Listen() {
	p := player.New(session, session.server.database)
	reader := bufio.NewReader(session.connection)

	session.Send(session.Address(), messages.HELLO, messages.HELLO()) // Send packet with Base64 header @@ to initialize connection with client.

	// Listen for incoming packets from a players session
	for {
		// Attempt to read three bytes; client->server packets in FUSEv0.2.0 begin with 3 byte Base64 encoded length.
		encodedLen := make([]byte, 3)

		for i := 0; i < 3; i++ {
			b, err := reader.ReadByte()
			if err != nil {
				log.Printf("There was an error while reading from %s's buffer. Err: %v \n", session.Address(), err)
				session.Close()
				return
			}
			encodedLen[i] = b
		}
		length := encoding.DecodeB64(encodedLen)

		// Check if data is junk before handling
		rawPacket := make([]byte, length)
		bytesRead, err := reader.Read(rawPacket)
		if err != nil {
			log.Printf("There was an error while reading from %s's buffer. Err: %v \n", session.Address(), err)
			continue
		} else if length == 0 {
			log.Println("Junk packet received, packet length is 0")
			continue
		} else if bytesRead < length {
			log.Printf("Packet length mismatch. Expected {%d bytes} but received {%d bytes} \n", length, bytesRead)
			continue
		}

		// Get Base64 encoded packet header
		payload := bytes.NewBuffer(rawPacket)
		rawHeader := make([]byte, 2)
		for i := 0; i < 2; i++ {
			rawHeader[i], _ = payload.ReadByte()
		}

		packet := packets.NewIncoming(rawHeader, payload)

		// Handle packets coming in from p's Session
		go Handle(p, packet, session.server.config.Debug)
	}
}

// Handle attempts to handle an incoming packet from a Players Session if it's registered in the Router
func Handle(p *player.Player, packet *packets.IncomingPacket, debug bool) {
	handler, found := p.Session.GetPacketCommand(packet.HeaderId)

	if found {
		if debug {
			// Check if the incoming packet is a handshake packet.
			// If the user is still logging in we don't have their username for logging so we use
			// their IP address for logging until handshake is complete
			// TODO this is kinda ugly, maybe remove this if logging is refactored
			switch packet.HeaderId {
			case 206, 202, 5, 6, 181, 4:
				logger.LogIncomingPacket(p.Session.Address(), handler, packet)
			default:
				logger.LogIncomingPacket(p.Details.Username, handler, packet)
			}
		}
		handler(p, packet)
	} else {
		if debug {
			logger.LogUnknownPacket(p.Details.Username, packet)
		}
	}

}

// Send finalizes an outgoing packet with 0x01 and then attempts to write the packet to a Session's buffer
// before flushing the buffer.
func (session *Session) Send(playerIdentifier string, caller interface{}, packet *packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	_, err := session.buffer.buff.Write(packet.Payload.Bytes())
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.Address(), err)
	}

	err = session.buffer.buff.Flush()
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.Address(), err)
	}

	if session.server.config.Debug {
		logger.LogOutgoingPacket(playerIdentifier, caller, packet)
	}
}

// Queue finalizes an outgoing packet with 0x01 and then attempts to write the packet to a Session's buffer.
func (session *Session) Queue(packet *packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	_, err := session.buffer.buff.Write(packet.Payload.Bytes())
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.Address(), err)
	}
}

// Flush attempts flush the packet to a Session's buffer.
func (session *Session) Flush(playerIdentifier string, caller interface{}, packet *packets.OutgoingPacket) {
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	err := session.buffer.buff.Flush()
	if err != nil {
		log.Printf("Error sending packet %v to session %v \n %v ", packet.Header, session.Address(), err)
	}

	if session.server.config.Debug {
		logger.LogOutgoingPacket(playerIdentifier, caller, packet)
	}
}

// GetPacketCommand attempts to retrieve a registered Command from the Router
func (session *Session) GetPacketCommand(headerId int) (func(*player.Player, *packets.IncomingPacket), bool) {
	return session.router.GetCommand(headerId)
}

// Address returns the IP address from a Session's connection
// Splits the address (e.g. 127.0.0.1:1234) and returns the IP part without the port
func (session *Session) Address() string {
	//
	return strings.Split(session.connection.RemoteAddr().String(), ":")[0]
}

// Close disconnects a Session from the server.
func (session *Session) Close() {
	if session.server.config.Debug {
		log.Printf("Closing session for address: %v ", session.Address())
	}

	session.server.RemoveSession(session)
	session.server = nil
	session.buffer = nil
	_ = session.connection.Close()
	session.active = false
}
