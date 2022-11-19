package server

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jtieri/habbgo/collections"
	"github.com/jtieri/habbgo/game/service"
	"go.uber.org/zap"
)

// Server is the main game server. It maintains a slice of active sessions connected to the server
// as well as a reference to all the available game services.
type Server struct {
	config   Config
	database *sql.DB

	sessions collections.Cache[string, *Session]
	services *service.Proxies

	log *zap.Logger
}

// Config is the game server configuration settings.
// NOTE: to avoid circular dependencies we avoid cmd.Config and use a local reference to the game server config.
type Config struct {
	Host              string
	Port              int
	MaxConnsPerPlayer int
	debug             bool
}

// New returns a pointer to a newly allocated Server struct.
func New(
	log *zap.Logger,
	database *sql.DB,
	host string,
	port, maxConnsPerPlayer int,
	debug bool,
	services *service.Proxies,
) *Server {
	return &Server{
		config: Config{
			Host:              host,
			Port:              port,
			MaxConnsPerPlayer: maxConnsPerPlayer,
			debug:             debug,
		},
		database: database,
		sessions: collections.NewCache(make(map[string]*Session)),
		services: services,
		log:      log,
	}
}

// Start will start the Server's main loop which listens for incoming TCP connections.
func (server *Server) Start(ctx context.Context) chan error {
	errorChan := make(chan error, 1)
	go server.HandleConnections(ctx, errorChan)
	return errorChan
}

// HandleConnections listens for new incoming connections and creates a new session
// for valid requests.
func (server *Server) HandleConnections(ctx context.Context, errorChan chan error) {
	defer close(errorChan)

	address := fmt.Sprintf("%s:%d", server.config.Host, server.config.Port)

	localAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		errorChan <- err
		return
	}

	// Use ListenTCP vs. net.Listen so that we can set a deadline on the listener,
	listener, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		errorChan <- err
		return
	}
	defer listener.Close()

	server.log.Info(
		"Successfully started the game server",
		zap.String("server_address", listener.Addr().String()),
	)

	// Main loop for handling connections.
	for {
		select {
		case <-ctx.Done():
			// The context was cancelled so return and the deferred call to server.Stop() will clean up.
			return
		default:
			// Set a deadline so that we don't stay blocking forever during listener.Accept()
			// This allows us to gracefully shutdown if the context is cancelled.
			if err := listener.SetDeadline(time.Now().Add(time.Second)); err != nil {
				continue
			}

			// Block and listen for incoming connections.
			conn, err := listener.Accept()
			if err != nil {
				if os.IsTimeout(err) {
					continue
				}

				server.log.Warn(
					"Error trying to handle incoming connection",
					zap.Error(err),
				)

				continue
			}

			// Check that there aren't multiple sessions for a given IP address.
			// TODO kick a session to make room for the new one instead of not letting new session connect
			if server.sessionsFromSameAddr(conn.LocalAddr().String()) >= server.config.MaxConnsPerPlayer {
				server.log.Info(
					"Too many sessions for address",
					zap.String("address", conn.LocalAddr().String()),
				)
				_ = conn.Close()
				continue

			}

			session := NewSession(
				server.log.With(zap.String("session", conn.LocalAddr().String())),
				conn,
				server,
			)

			server.log.Info(
				"New session created",
				zap.String("address", conn.LocalAddr().String()),
				zap.Int("num_sessions_for_usr", server.sessionsFromSameAddr(conn.LocalAddr().String())),
			)

			server.sessions.Set(session.address(), session)

			go session.listen(ctx)
		}

	}
}

// RemoveSession removes a Session from the slice of active Sessions and adjusts the slice so that there are no gaps.
func (server *Server) RemoveSession(session *Session) {
	server.sessions.Remove(session.address())
	server.log.Debug(
		"Active sessions updated",
		zap.Int("num_active_sessions", server.sessions.Count()),
	)
}

// Stop terminates all active sessions and shuts down the game server.
func (server *Server) Stop() {
	defer server.database.Close()

	// TODO send shutdown signals to the running ActorServices
	// Need to wait till we know each service is shutdown before we close all the sessions.
	// Otherwise the Players session will be terminated before the PlayerActor is
	// finished clearing its task queue which could result in data loss.

	sessions := server.sessions.Items()
	for _, session := range sessions {
		session.Close()
		session = nil
	}

	server.log.Info("Shutting down game server")
}

// sessionsFromSameAddr returns the number of active Sessions connected to the server for one IP address.
func (server *Server) sessionsFromSameAddr(address string) int {
	count := 0

	for _, session := range server.sessions.Items() {
		if address == session.connection.LocalAddr().String() {
			count++
		}
	}

	return count
}
