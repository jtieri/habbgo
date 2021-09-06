package service

import (
	"github.com/jtieri/HabbGo/habbgo/game/model"
)

func Login(player *model.Player) {
	// Set player logged in & ping ready for latency test
	// Possibly add player to a list of online players? Health endpoint with server stats?
	// Save current time to Conn for players last online time

	// Check if player is banned & if so send USER_BANNED
	// Log IP address to Conn

	// database.LoadBadges(player)

	// If Config has alerts enabled, send player ALERT

	// Check if player gets club gift & update club status
}
