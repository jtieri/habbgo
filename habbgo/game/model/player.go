package model

import "github.com/jtieri/HabbGo/habbgo/server"

type Player struct {
	session       *server.Session
	playerDetails *PlayerDetails
}

type PlayerDetails struct {
	username string
}
