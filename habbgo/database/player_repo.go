package database

import (
	"github.com/jtieri/HabbGo/habbgo/game/model/player"
	"log"
)

func Login(player *player.Player, username string, password string) bool {
	var pw string
	err := player.Session.Database().QueryRow("SELECT P.passwrd FROM Players P WHERE P.username = ?", username).Scan(&pw)

	if err != nil {
		log.Printf("%v ", err) // TODO log database errors properly
	}

	if password == pw {
		fillDetails(player.Details)
		return true
	}

	return false
}

func fillDetails(pd *player.Details) {

}
