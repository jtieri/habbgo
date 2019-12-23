package service

import (
	"github.com/jtieri/HabbGo/habbgo/database"
	"github.com/jtieri/HabbGo/habbgo/game/model"
	"github.com/jtieri/HabbGo/habbgo/protocol/composers"
)

func Login(player *model.Player) {
	// Set player logged in & ping ready for latency test
	// Possibly add player to a list of online players? Health endpoint with server stats?
	// Save current time to DB for players last online time

	// Check if player is banned & if so send USER_BANNED
	// Log IP address to DB

	database.LoadBadges(player)
	go player.Session.Send(composers.ComposeLoginOk())

	// If Config has alerts enabled, send player ALERT

	// Check if player gets club gift & update club status
}
