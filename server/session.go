package server

import (
	"bufio"
	"bytes"
	"net"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/protocol/encoding"
	"github.com/jtieri/habbgo/protocol/messages"
	"github.com/jtieri/habbgo/protocol/packets"
	"go.uber.org/zap"
)

// Session represents a player.Player's underlying network session and connection to the server.
type Session struct {
	connection net.Conn
	buffer     *buffer
	active     bool
	server     *Server
	router     *Router
	log        *zap.Logger
}

// buffer is the buffered Writer used to write data to a Session's connection.
type buffer struct {
	mux  sync.Mutex
	buff *bufio.Writer
}

// NewSession returns a pointer to a newly allocated Session struct.
func NewSession(log *zap.Logger, conn net.Conn, server *Server) *Session {
	return &Session{
		connection: conn,
		buffer:     &buffer{mux: sync.Mutex{}, buff: bufio.NewWriter(conn)},
		active:     true,
		server:     server,
		router:     RegisterCommands(),
		log:        log,
	}
}

// Listen starts listening for incoming data from a Session's connection and handles it appropriately as
// per the FUSEv0.2.0 protocol.
func (session *Session) Listen() {
	p := player.New(
		session.log.With(),
		session,
		session.server.database,
		session.server.services,
	)
	reader := bufio.NewReader(session.connection)

	// Send packet with Base64 header @@ to initialize connection with client.
	session.Send(messages.HELLO, messages.HELLO())

	// Listen for incoming packets from a player's session.
	for {
		// Attempt to read three bytes,
		// client->server packets in FUSEv0.2.0 begin with 3 byte Base64 encoded packet length.
		encodedLen := make([]byte, 3)

		for i := 0; i < 3; i++ {
			b, err := reader.ReadByte()
			if err != nil {
				// If the network connection is closed, it's because the server closed the Session
				// which means we don't need to log again or call session.Close
				if strings.Contains(err.Error(), "use of closed network connection") {
					return
				}

				session.log.Warn("Error reading encoded packet length from session",
					zap.String("session_address", session.Address()),
					zap.Error(err),
				)
				session.Close()
				return
			}
			encodedLen[i] = b
		}
		packetLen := encoding.DecodeB64(encodedLen)

		// Check if data is junk before handling.
		rawPacket := make([]byte, packetLen)
		bytesRead, err := reader.Read(rawPacket)
		switch {
		case err != nil:
			session.log.Warn("Error reading packet data from session",
				zap.String("session_address", session.Address()),
				zap.Error(err),
			)
			session.Close()
			return
		case packetLen == 0:
			session.log.Info("Junk packet received")
			continue
		case bytesRead < packetLen:
			session.log.Info("Packet length mismatch",
				zap.Int("expected_length", packetLen),
				zap.Int("got_length", bytesRead),
			)
			continue
		}

		// Get Base64 encoded packet header.
		payload := bytes.NewBuffer(rawPacket)
		rawHeader := make([]byte, 2)
		for i := 0; i < 2; i++ {
			rawHeader[i], _ = payload.ReadByte()
		}

		packet := packets.NewIncoming(rawHeader, payload)

		// Handle packets coming in from the Player's Session.
		go session.Handle(p, packet)
	}
}

// Handle attempts to handle an incoming packet from a player.Player's Session.
// If the packet is not registered in the Router with an appropriate handler,
// the packet is ignored.
func (session *Session) Handle(p *player.Player, packet *packets.IncomingPacket) {
	handler, found := p.Session.GetPacketCommand(packet.HeaderId)

	if found {
		// Avoid using reflection unless the server is in debug mode.
		var handlerName string
		if session.server.config.debug {
			handlerName = GetPacketHandlerName(handler)
		}

		// If the user is still logging in we don't have their username for logging so,
		// we check to see if we should log it or not.
		// TODO possibly just remove this and never log usernames
		switch {
		case p.Details.Username != "":
			session.log.Debug("Incoming Packet",
				zap.String("player_name", p.Details.Username),
				zap.String("packet_name", handlerName),
				zap.String("packet_header", packet.Header),
				zap.Int("header_id", packet.HeaderId),
				zap.String("payload", packet.Payload.String()),
			)
		default:
			session.log.Debug("Incoming Packet",
				zap.String("packet_name", handlerName),
				zap.String("packet_header", packet.Header),
				zap.Int("header_id", packet.HeaderId),
				zap.String("payload", packet.Payload.String()),
			)
		}

		handler(p, packet)
	} else {
		session.log.Debug("Incoming Packet",
			zap.String("player_name", p.Details.Username),
			zap.String("packet_header", packet.Header),
			zap.Int("header_id", packet.HeaderId),
			zap.String("payload", packet.Payload.String()),
		)
	}

}

// Send finalizes an outgoing packet with 0x01 and then attempts to write the packet to a Session's buffer
// before flushing the buffer.
func (session *Session) Send(caller interface{}, packet *packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	_, err := session.buffer.buff.Write(packet.Payload.Bytes())
	if err != nil {
		session.log.Warn("Error writing packet to session buffer",
			zap.String("packet_name", GetPacketHandlerName(caller)),
			zap.String("packet_header", packet.Header),
			zap.Int("header_id", packet.HeaderId),
			zap.String("payload", packet.Payload.String()),
			zap.Error(err),
		)
		session.Close()
		return
	}

	err = session.buffer.buff.Flush()
	if err != nil {
		session.log.Warn("Error sending packet to session",
			zap.String("packet_name", GetPacketHandlerName(caller)),
			zap.String("packet_header", packet.Header),
			zap.Int("header_id", packet.HeaderId),
			zap.String("payload", packet.Payload.String()),
			zap.Error(err),
		)
		session.Close()
		return
	}

	session.log.Debug("Outgoing Packet",
		zap.String("packet_name", GetPacketHandlerName(caller)),
		zap.String("packet_header", packet.Header),
		zap.Int("header_id", packet.HeaderId),
		zap.String("payload", packet.Payload.String()),
	)
}

// Queue finalizes an outgoing packet with 0x01 and then attempts to write the packet to a Session's buffer.
func (session *Session) Queue(packet *packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	_, err := session.buffer.buff.Write(packet.Payload.Bytes())
	if err != nil {
		session.log.Warn("Error writing packet to session buffer",
			zap.String("packet_name", GetPacketHandlerName(packet)),
			zap.String("packet_header", packet.Header),
			zap.Int("header_id", packet.HeaderId),
			zap.String("payload", packet.Payload.String()),
			zap.Error(err),
		)
		session.Close()
		return
	}
}

// Flush attempts flush the packet to a Session's buffer.
func (session *Session) Flush(caller interface{}, packet *packets.OutgoingPacket) {
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	err := session.buffer.buff.Flush()
	if err != nil {
		session.log.Warn("Error sending packet to session",
			zap.String("packet_name", GetPacketHandlerName(caller)),
			zap.String("packet_header", packet.Header),
			zap.Int("header_id", packet.HeaderId),
			zap.String("payload", packet.Payload.String()),
			zap.Error(err),
		)
		session.Close()
		return
	}

	session.log.Debug("Outgoing Packet",
		zap.String("packet_name", GetPacketHandlerName(caller)),
		zap.String("packet_header", packet.Header),
		zap.Int("header_id", packet.HeaderId),
		zap.String("payload", packet.Payload.String()),
	)
}

// GetPacketCommand attempts to retrieve a registered Command from the Router.
// If there is no Command registered in the router with header ID equal to the function parameter,
// false will be returned.
func (session *Session) GetPacketCommand(headerId int) (func(*player.Player, *packets.IncomingPacket), bool) {
	return session.router.GetCommand(headerId)
}

// Address returns the IP address from a Session's connection
// Splits the address (e.g. 127.0.0.1:1234) and returns the IP part without the port
func (session *Session) Address() string {
	return strings.Split(session.connection.RemoteAddr().String(), ":")[0]
}

// Close disconnects a Session from the server.
func (session *Session) Close() {
	defer session.connection.Close()

	session.log.Debug("Closing session",
		zap.String("session_addr", session.Address()),
	)

	session.server.RemoveSession(session)
	session.active = false
}

// GetPacketHandlerName is a hacky way to get the name of the incoming/outgoing packet function call,
// this is useful in debugging so you can quickly analyze the flow of packets.
func GetPacketHandlerName(message interface{}) string {
	handler := runtime.FuncForPC(reflect.ValueOf(message).Pointer()).Name()
	sp := strings.Split(handler, "/") // e.g. github.com/jtieri/habbgo/protocol/handlers.GenerateKey
	s2 := sp[len(sp)-1]               // e.g. handlers.GenerateKey
	return strings.Split(s2, ".")[1]  // e.g. GenerateKey
}
