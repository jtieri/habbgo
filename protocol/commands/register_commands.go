package commands

import (
	"net/mail"
	"strings"
	"time"

	"github.com/jtieri/habbgo/crypto"
	"github.com/jtieri/habbgo/date"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/protocol/messages"
	"github.com/jtieri/habbgo/protocol/packets"
	"github.com/jtieri/habbgo/text"
)

const (
	ALLOWEDCHARS          = "1234567890qwertyuiopasdfghjklzxcvbnm_-+=?!@:.,$" // TODO make this a config option
	OK                    = 0
	NAMETOOLONG           = 1
	NAMETOOSHORT          = 2
	NAMEUNACCEPTABLE      = 3
	NAMEALREADYRESERVED   = 4
	PASSWORDTOOSHORT      = 1
	PASSWORDTOOLONG       = 2
	PASSWORDUNACCEPTABLE  = 3
	PASSWORDHASNONUM      = 4
	PASSWORDSIMILARTONAME = 5
)

func GETAVAILABLESETS(p *player.Player, packet *packets.IncomingPacket) {
	p.Session.Send(messages.AVAILABLESETS, messages.AVAILABLESETS())
}

func GDATE(p *player.Player, packet *packets.IncomingPacket) {
	p.Session.Send(messages.DATE, messages.DATE(date.GetCurrentDate()))
}

func APPROVENAME(p *player.Player, packet *packets.IncomingPacket) {
	name := text.Filter(packet.ReadString())
	p.Session.Send(messages.APPROVENAMEREPLY, messages.APPROVENAMEREPLY(checkName(p, name)))
}

func APPROVE_PASSWORD(p *player.Player, packet *packets.IncomingPacket) {
	username := packet.ReadString()
	password := packet.ReadString()
	p.Session.Send(messages.PASSWORD_APPROVED, messages.PASSWORD_APPROVED(checkPassword(p, username, password)))
}

func APPROVEEMAIL(p *player.Player, packet *packets.IncomingPacket) {
	email := packet.ReadString()

	if _, err := mail.ParseAddress(email); err != nil {
		p.Session.Send(messages.EMAIL_REJECTED, messages.EMAIL_REJECTED())
	} else {
		p.Session.Send(messages.EMAIL_APPROVED, messages.EMAIL_APPROVED())
	}
}

func REGISTER(p *player.Player, packet *packets.IncomingPacket) {
	packet.ReadB64()
	username := packet.ReadString()

	packet.ReadB64()
	figure := packet.ReadString()

	packet.ReadB64()
	gender := packet.ReadString()

	packet.ReadB64()
	packet.ReadB64()

	packet.ReadB64()
	email := packet.ReadString()

	packet.ReadB64()
	birthday := packet.ReadString()

	packet.ReadBytes(11)
	password := packet.ReadString()

	// hash password before storing in db
	randSalt := crypto.GenerateRandomSalt(crypto.SALTSIZE)
	hPsswrd := crypto.HashPassword(password, randSalt)

	// generate date time stamp in UTC with format YYYY-MM-DD T HH-MM-SS
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	createdAt := now.Format("2006-01-02 15:04:05")

	// put birthday in YYYY-MM-DD format before storing in db
	var bday string
	for i, s := range strings.Split(birthday, ".") {
		if i == 0 {
			bday = s + bday
		} else {
			bday = s + "-" + bday
		}
	}

	p.Register(username, figure, gender, email, bday, createdAt, hPsswrd, randSalt)
	/*
		2021/09/16 22:28:48 [127.0.0.1] [UNK] [@k - 43]: @B@Itreebeard@D@Y1000118001270012900121001@E@AM@F@@@G@Mboob@none.com@H@J27.01.1995@JA@A@@I@@C@Jtreebeard1
		2021/09/16 22:29:04 [127.0.0.1] [INCOMING] [TRY_LOGIN - @D|4]: @Itreebeard@Jtreebeard1
	*/
}

// checkName takes in a proposed username and returns an integer representing the approval status of the given name
func checkName(p *player.Player, username string) int {
	switch {
	case player.PlayerExists(p, username):
		return NAMEALREADYRESERVED
	case len(username) > 16:
		return NAMETOOLONG
	case len(username) < 1:
		return NAMETOOSHORT
	case !text.ContainsAllowedChars(strings.ToLower(username), ALLOWEDCHARS) || strings.Contains(username, " "):
		return NAMEUNACCEPTABLE
	case strings.Contains(strings.ToUpper(username), "MOD-"):
		return NAMEUNACCEPTABLE
	default:
		return OK
	}
}

// checkPassword takes in a proposed password and returns an integer representing the approval status of the given password
func checkPassword(p *player.Player, username, password string) int {
	switch {
	case len(password) < 6:
		return PASSWORDTOOSHORT // too short
	case len(password) > 16:
		return PASSWORDTOOLONG // too long
	case !text.ContainsAllowedChars(strings.ToLower(password), ALLOWEDCHARS) || strings.Contains(username, " "):
		return PASSWORDUNACCEPTABLE // using non-permitted characters
	case !text.ContainsANumber(password):
		return PASSWORDHASNONUM // password does not contain a number
	case username == password:
		return PASSWORDSIMILARTONAME // name and pass too similar
	default:
		return OK
	}
}
