package server

import (
	"database/sql"
	"fmt"
	"github.com/jtieri/habbgo/config"
	"log"
	"net"
	"os"
	"sync"
)

// Server represents the main game server
type Server struct {
	activeSessions []*Session
	config         *config.Config
	database       *sql.DB
}

// New returns a pointer to a newly allocated server struct.
func New(config *config.Config, database *sql.DB) *Server {
	return &Server{
		config:   config,
		database: database,
	}
}

// Start will set up the game server, start listening for incoming connections, and handle connections appropriately.
func (server *Server) Start() {
	address := fmt.Sprintf("%s:%d", server.config.ServerHost, server.config.ServerPort)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("There was an issue starting the game server at %v. Err: %v", address, err)
	}
	log.Printf("Successfully started the game server at %v", listener.Addr().String())
	defer listener.Close()

	server.HandleConnections(listener)
}

// HandleConnections listens for new connections and creates a new session for valid requests from the listener
func (server *Server) HandleConnections(listener net.Listener) {
	// Main loop for handling connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error trying to handle incoming connection. Err: %v", err)
			_ = conn.Close()
			continue
		}

		// Check that there aren't multiple sessions for a given IP address
		// TODO kick a session to make room for the new one
		if server.sessionsFromSameAddr(conn) < server.config.MaxConnsPerPlayer {
			session := NewSession(conn, server)
			log.Printf("New session created for address: %v", conn.LocalAddr().String())
			server.activeSessions = append(server.activeSessions, session)

			go session.Listen()
		} else {
			log.Printf("Too many concurrent connections from address %v \n", conn.LocalAddr().String())
			_ = conn.Close()
		}
	}
}

// RemoveSession removes a session from the slice of active connections and adjusts the slice so that there are no gaps
func (server *Server) RemoveSession(session *Session) {
	mux := sync.Mutex{}
	for i, activeSession := range server.activeSessions {
		if activeSession.connection.LocalAddr().String() == session.connection.LocalAddr().String() {
			mux.Lock()
			// This re-adjusts the slice of active connections so there are no gaps in the slice
			// i.e. there is an active session at every index in the slice
			server.activeSessions[i] = server.activeSessions[len(server.activeSessions)-1]
			server.activeSessions[len(server.activeSessions)-1] = nil
			server.activeSessions = server.activeSessions[:len(server.activeSessions)-1]
			mux.Unlock()

			log.Printf("There are now %v sessions connected to the server. ", len(server.activeSessions))
			break
		}
	}
}

// Stop terminates all active sessions and shuts down the game server.
func (server *Server) Stop() {
	for _, session := range server.activeSessions {
		session.Close()
	}

	log.Println("Shutting down the game server...")
	os.Exit(0)
}

// sessionsFromSameAddr returns the number of connections to the server coming from one IP.
func (server *Server) sessionsFromSameAddr(conn net.Conn) int {
	count := 0
	for _, session := range server.activeSessions {
		if conn.LocalAddr().String() == session.connection.LocalAddr().String() {
			count++
		}
	}

	return count
}
