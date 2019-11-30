package server

import (
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

type server struct {
	host           string
	port           int16
	activeSessions []*Session
}

// New returns a pointer to a newly allocated server struct.
func New(port int16, host string) *server {
	server := &server{port: port, host: host}
	return server
}

// Start will setup the game server to listen for incoming connections and handle them appropriately.
func (server *server) Start() {
	listener, err := net.Listen("tcp", server.host+":"+strconv.Itoa(int(server.port)))
	if err != nil {
		log.Fatalf("There was an issue starting the game server on port %v.", server.port)
	}
	log.Printf("Successfully started the game server at %v", listener.Addr().String())
	defer listener.Close()

	// Main loop for handling connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error trying to handle new connection.")
			conn.Close()
			continue
		}

		// Check that there aren't multiple sessions for a given IP address
		if !(server.sessionsFromSameAddr(conn) >= 1) {
			session := &Session{
				connection: conn,
				server:     server,
			}

			log.Printf("New connection from: %v", conn.LocalAddr().String())
			server.activeSessions = append(server.activeSessions, session)
			go session.Listen()
		} else {
			log.Printf("Too many concurrent connections from address %v \n", conn.LocalAddr().String())
			conn.Close()
		}
	}
}

func (server *server) RemoveSession(session *Session) {
	mux := &sync.Mutex{}
	for i, activeSession := range server.activeSessions {
		if activeSession.connection.LocalAddr().String() == session.connection.LocalAddr().String() {
			mux.Lock()
			server.activeSessions[i] = server.activeSessions[len(server.activeSessions)-1]
			server.activeSessions[len(server.activeSessions)-1] = nil
			server.activeSessions = server.activeSessions[:len(server.activeSessions)-1]
			mux.Unlock()

			log.Printf("There are now %v sessions connected to the server. ", len(server.activeSessions))
		}
	}
}

// Stop terminates all active sessions and shuts down the game server.
func (server *server) Stop() {
	for _, session := range server.activeSessions {
		session.Close()
	}

	log.Println("Shutting down the game server...")
	os.Exit(0)
}

func (server *server) sessionsFromSameAddr(conn net.Conn) int {
	count := 0

	for _, session := range server.activeSessions {
		if conn.LocalAddr().String() == session.connection.LocalAddr().String() {
			count++
		}
	}

	return count
}
