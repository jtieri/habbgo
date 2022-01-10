package player

import "fmt"

func (p *Player) Login() {
	// Set player logged in & ping ready for latency test
	// Possibly add player to a list of online players? Health endpoint with server stats?
	// Save current time to Conn for players last online time

	// Check if player is banned & if so send USER_BANNED
	// Log IP address to Conn

	LoadBadges(p)

	// If Config has alerts enabled, send player ALERT

	// Check if player gets club gift & update club status
}

func (p *Player) Register(username, figure, gender, email, birthday, createdAt, password string, salt []byte) {
	err := Register(p, username, figure, gender, email, birthday, createdAt, password, salt)
	if err != nil {
		p.LogErr(err)
	}
}

func (p *Player) LogErr(err error) {
	fmt.Printf("[%s-%s] Player encountered error: %s \n", p.Details.Id, p.Details.Username, err)
}
