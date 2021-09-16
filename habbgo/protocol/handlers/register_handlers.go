package handlers

import (
	"github.com/jtieri/HabbGo/habbgo/date"
	"github.com/jtieri/HabbGo/habbgo/game/player"
	"github.com/jtieri/HabbGo/habbgo/protocol/composers"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
	"github.com/jtieri/HabbGo/habbgo/text"
	"strings"
)

const (
	OK              = 0
	TOOLONG         = 1
	TOOSHORT        = 2
	UNACCEPTABLE    = 3
	ALREADYRESERVED = 4
)

// checkName takes in a proposed username and returns an integer representing the approval status of the given name
func checkName(username string) int {
	allowedChars := "1234567890qwertyuiopasdfghjklzxcvbnm_-+=?!@:.,$" // TODO make this a config option
	switch {
	//TODO add case for username already being taken
	case len(username) > 16:
		return TOOLONG
	case len(username) < 1:
		return TOOSHORT
	case !text.ContainsAllowedChars(strings.ToLower(username), allowedChars) || strings.Contains(username, " "):
		return UNACCEPTABLE
	case strings.Contains(strings.ToUpper(username), "MOD-"):
		return UNACCEPTABLE
	default:
		return OK
	}
}

func GDATE(p *player.Player, packet *packets.IncomingPacket) {
	p.Session.Send(composers.DATE(date.GetCurrentDate()))
}

func APPROVENAME(p *player.Player, packet *packets.IncomingPacket) {
	name := text.Filter(packet.ReadString())
	p.Session.Send(composers.APPROVENAMEREPLY(checkName(name)))
}

//func APPROVE_PASSWORD(p *player.Player, packet *packets.IncomingPacket) {
//	username := packet.ReadString()
//	password := packet.ReadString()
//
//	errorCode := 0
//
//	p.Session.Send(composers.PASSWORD_APPROVED(errorCode))
//}
