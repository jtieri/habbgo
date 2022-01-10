package player

import (
	"github.com/jtieri/habbgo/crypto"
	"log"
)

func Register(player *Player, username, figure, gender, email, birthday, createdAt, password string, salt []byte) error {
	stmt, err := player.Database.Prepare(
		"INSERT INTO Players(username, figure, sex, email, birthday, created_on, password_hash, password_salt) VALUES($1, $2, $3, $4, $5, $6, $7, $8)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(username, figure, gender, email, birthday, createdAt, password, salt)
	if err != nil {
		return err
	}

	return nil
}

func LoginDB(player *Player, username string, password string) bool {
	var (
		psswrdHash, uname string
		psswrdSalt        []byte
	)

	err := player.Database.QueryRow(
		"SELECT P.password_hash, P.password_salt, P.username FROM Players P WHERE P.username = $1", username).
		Scan(&psswrdHash, &psswrdSalt, &uname)

	if err != nil {
		player.LogErr(err)
	}

	if crypto.HashPassword(password, psswrdSalt) == psswrdHash {
		player.Details.Username = uname
		fillDetails(player)
		return true
	}

	return false
}

func LoadBadges(player *Player) {
	rows, err := player.Database.Query("SELECT P.badge_id FROM player_badges P WHERE P.player_id = $1", player.Details.Id)
	if err != nil {
		log.Printf("%v ", err) // TODO properly log error
	}
	defer rows.Close()

	var badges []string
	for rows.Next() {
		var badge string
		err := rows.Scan(&badge)
		if err != nil {
			log.Printf("%v ", err) // TODO properly log error
		}

		badges = append(badges, badge)
	}

	player.Details.Badges = badges
}

func PlayerExists(p *Player, username string) bool {
	rows, err := p.Database.Query("SELECT P.id FROM Players P WHERE P.username = $1", username)
	if err != nil {
		log.Printf("%s", err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	if rows.Err() != nil {
		log.Printf("%s", err)
	}

	return false
}

func UpdateLastOnline(datetime string) {

}

func fillDetails(p *Player) {
	query := "SELECT P.id, P.username, P.sex, P.figure, P.pool_figure, P.film, P.credits, P.tickets, P.motto, " +
		"P.console_motto, P.last_online, P.sound_enabled, P.Rank " +
		"FROM Players P " +
		"WHERE P.username = $1"

	var tmpRank string
	err := p.Database.QueryRow(query, p.Details.Username).Scan(&p.Details.Id, &p.Details.Username,
		&p.Details.Sex, &p.Details.Figure, &p.Details.PoolFigure, &p.Details.Film, &p.Details.Credits,
		&p.Details.Tickets, &p.Details.Motto, &p.Details.ConsoleMotto, &p.Details.LastOnline, &p.Details.SoundEnabled,
		&tmpRank)

	if err != nil {
		log.Printf("%v ", err) // TODO log database errors properly
	}

	p.Details.PlayerRank = PlayerRank(tmpRank)
}
