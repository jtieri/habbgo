package server

import (
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/jtieri/HabbGo/habbgo/app"
)

type Server struct {
	activeSessions []*Session
}

// New returns a pointer to a newly allocated server struct.
func New() *Server {
	return &Server{}
}

// Start will setup the game server, start listening for incoming connections, and handle connections appropriately.
func (server *Server) Start() {
	listener, err := net.Listen("tcp", app.HabbGo().Config.Server.Host+":"+strconv.Itoa(app.HabbGo().Config.Server.Port))
	if err != nil {
		log.Fatalf("There was an issue starting the game server on port %v.", app.HabbGo().Config.Server.Port) // TODO properly handle errors
	}
	log.Printf("Successfully started the game server at %v", listener.Addr().String())
	defer listener.Close()

	// Main loop for handling connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error trying to handle incoming connection.") // TODO properly handle errors
			_ = conn.Close()
			continue
		}

		// Check that there aren't multiple sessions for a given IP address
		// TODO kick a session to make room for the new one
		if server.sessionsFromSameAddr(conn) < app.HabbGo().Config.Server.MaxConns {
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

func (server *Server) RemoveSession(session *Session) {
	mux := sync.Mutex{}
	for i, activeSession := range server.activeSessions {
		if activeSession.connection.LocalAddr().String() == session.connection.LocalAddr().String() {
			mux.Lock()
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
