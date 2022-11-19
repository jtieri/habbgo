package server

import (
	"bufio"
	"bytes"
	"context"
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

// listen starts listening for incoming data from a Session's connection and handles it appropriately as
// per the FUSEv0.2.0 protocol.
func (session *Session) listen(ctx context.Context) {
	p := player.New(
		ctx,
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
		encodedLen, err := readPacketLength(reader)
		if err != nil {
			// If the network connection is closed, it's because the server closed the Session
			// which means we don't need to log again or call session.Close
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}

			session.log.Error(
				"Error reading encoded packet length from session",
				zap.String("session_address", session.address()),
				zap.Error(err),
			)
			session.Close()
			return
		}

		// Decode packet length & check if data is junk before handling.
		packetLen := encoding.DecodeB64(encodedLen)
		if !validate(session, packetLen, reader.Size()) {
			// R we find junk data in the connections buffer
			reader.Reset(session.connection)
			continue
		}

		// Build a packet object from the remaining bytes
		packetBz := make([]byte, packetLen)

		if _, err = reader.Read(packetBz); err != nil {
			session.log.Error(
				"Error reading packet data from session",
				zap.String("session_address", session.address()),
				zap.Error(err),
			)
			session.Close()
			return
		}

		packet, err := packetFromBytes(packetBz)
		if err != nil {
			session.log.Error(
				"Error reading packet data from session",
				zap.String("session_address", session.address()),
				zap.Error(err),
			)
			session.Close()
			return
		}

		// handle packets coming in from the Player's Session.
		ps := p.Services.PlayerService().(*player.PlayerServiceProxy)
		cachedPlayer := ps.GetPlayer(p)

		if cachedPlayer.LoggedIn() {
			go session.handle(cachedPlayer, packet)
		} else {
			go session.handle(p, packet)
		}
	}
}

// packetFromBytes attempts to build a packets.IncomingPacket from a slice of bytes.
func packetFromBytes(packetBytes []byte) (packets.IncomingPacket, error) {
	payload := bytes.NewBuffer(packetBytes)
	rawHeader := make([]byte, 2)

	for i := 0; i < 2; i++ {
		b, err := payload.ReadByte()
		if err != nil {
			return packets.IncomingPacket{}, err
		}
		rawHeader[i] = b
	}

	return packets.NewIncoming(rawHeader, payload), nil
}

// readPacketLength attempts to read the 3 byte Base64 encoded packet length from the buffered reader.
func readPacketLength(reader *bufio.Reader) ([]byte, error) {
	encodedLen := make([]byte, 3)

	for i := 0; i < 3; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		encodedLen[i] = b
	}

	return encodedLen, nil
}

// validate checks if an incoming packet is valid.
// It's argument is a 3 byte Base64 encoded length.
func validate(session *Session, packetLen, bytesRead int) bool {

	switch {
	case packetLen == 0:
		session.log.Info("Junk packet received")
		return false
	case bytesRead < packetLen:
		session.log.Info(
			"Packet length mismatch",
			zap.Int("expected_length", packetLen),
			zap.Int("got_length", bytesRead),
		)
		return false
	}

	return true
}

// handle attempts to handle an incoming packet from a player.Player's Session.
// If the packet is not registered in the Router with an appropriate handler,
// the packet is ignored.
func (session *Session) handle(p player.Player, packet packets.IncomingPacket) {
	handler, found := session.router.GetCommand(packet.HeaderId)

	if !found {
		session.log.Debug(
			"Incoming Packet",
			zap.String("player_name", p.Details.Username),
			zap.String("packet_header", packet.Header),
			zap.Int("header_id", packet.HeaderId),
			zap.String("payload", packet.Payload.String()),
		)
		return
	}

	// Avoid using reflection unless the server is in debug mode.
	var handlerName string
	if session.server.config.debug {
		handlerName = getPacketHandlerName(handler)
	}

	// If the user is still logging in we don't have their username for logging so,
	// we check to see if we should log it or not.
	// TODO possibly just remove this and never log usernames
	//switch {
	//case p.Details.Username != "":
	//	session.log.Debug(
	//		"Incoming Packet",
	//		zap.String("player_name", p.Details.Username),
	//		zap.String("packet_name", handlerName),
	//		zap.String("packet_header", packet.Header),
	//		zap.Int("header_id", packet.HeaderId),
	//		zap.String("payload", packet.Payload.String()),
	//	)
	//default:
	//	session.log.Debug(
	//		"Incoming Packet",
	//		zap.String("packet_name", handlerName),
	//		zap.String("packet_header", packet.Header),
	//		zap.Int("header_id", packet.HeaderId),
	//		zap.String("payload", packet.Payload.String()),
	//	)
	//}
	session.log.Debug(
		"Incoming Packet",
		zap.String("packet_name", handlerName),
		zap.String("packet_header", packet.Header),
		zap.Int("header_id", packet.HeaderId),
		zap.String("payload", packet.Payload.String()),
	)

	handler(p, packet)
}

// Send finalizes an outgoing packet with 0x01 and then attempts to write the packet to a Session's buffer
// before flushing the buffer.
func (session *Session) Send(caller interface{}, packet packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	_, err := session.buffer.buff.Write(packet.Payload.Bytes())
	if err != nil {
		session.log.Warn(
			"Error writing packet to session buffer",
			zap.String("packet_name", getPacketHandlerName(caller)),
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
		session.log.Warn(
			"Error sending packet to session",
			zap.String("packet_name", getPacketHandlerName(caller)),
			zap.String("packet_header", packet.Header),
			zap.Int("header_id", packet.HeaderId),
			zap.String("payload", packet.Payload.String()),
			zap.Error(err),
		)
		session.Close()
		return
	}

	session.log.Debug(
		"Outgoing Packet",
		zap.String("packet_name", getPacketHandlerName(caller)),
		zap.String("packet_header", packet.Header),
		zap.Int("header_id", packet.HeaderId),
		zap.String("payload", packet.Payload.String()),
	)
}

// Queue finalizes an outgoing packet with 0x01 and then attempts to write the packet to a Session's buffer.
func (session *Session) Queue(caller interface{}, packet packets.OutgoingPacket) {
	packet.Finish()
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	_, err := session.buffer.buff.Write(packet.Payload.Bytes())
	if err != nil {
		session.log.Warn(
			"Error writing packet to session buffer",
			zap.String("packet_name", getPacketHandlerName(caller)),
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
func (session *Session) Flush(caller interface{}, packet packets.OutgoingPacket) {
	session.buffer.mux.Lock()
	defer session.buffer.mux.Unlock()

	err := session.buffer.buff.Flush()
	if err != nil {
		session.log.Warn(
			"Error sending packet to session",
			zap.String("packet_name", getPacketHandlerName(caller)),
			zap.String("packet_header", packet.Header),
			zap.Int("header_id", packet.HeaderId),
			zap.String("payload", packet.Payload.String()),
			zap.Error(err),
		)
		session.Close()
		return
	}

	session.log.Debug(
		"Outgoing Packet",
		zap.String("packet_name", getPacketHandlerName(caller)),
		zap.String("packet_header", packet.Header),
		zap.Int("header_id", packet.HeaderId),
		zap.String("payload", packet.Payload.String()),
	)
}

// Close disconnects a Session from the server.
func (session *Session) Close() {
	defer session.connection.Close()

	session.log.Debug(
		"Closing session",
		zap.String("session_addr", session.address()),
	)

	session.server.RemoveSession(session)

	session.active = false
}

// address returns the IP address from a Session's connection
// Splits the address (e.g. 127.0.0.1:1234) and returns the IP part without the port
func (session *Session) address() string {
	return strings.Split(session.connection.RemoteAddr().String(), ":")[0]
}

// getPacketHandlerName is a hacky way to get the name of the incoming/outgoing packet function call,
// this is useful in debugging so you can quickly analyze the flow of packets.
func getPacketHandlerName(message interface{}) string {
	handler := runtime.FuncForPC(reflect.ValueOf(message).Pointer()).Name()
	sp := strings.Split(handler, "/") // e.g. github.com/jtieri/habbgo/protocol/handlers.GenerateKey
	s2 := sp[len(sp)-1]               // e.g. handlers.GenerateKey
	return strings.Split(s2, ".")[1]  // e.g. GenerateKey
}
