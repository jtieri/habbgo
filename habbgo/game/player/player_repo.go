package player

import (
	"github.com/jtieri/HabbGo/habbgo/app"
	"github.com/jtieri/HabbGo/habbgo/crypto"
	"log"
)

func Register(username, figure, gender, email, birthday, createdAt, password string, salt []byte) error {
	stmt, err := app.HabbGo().Database.Prepare(
		"INSERT INTO Players(Username, Figure, Sex, Email, Birthday, CreatedOn, PasswordHash, PasswordSalt) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")

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

	err := app.HabbGo().Database.QueryRow(
		"SELECT P.PasswordHash, P.PasswordSalt, P.Username FROM Players P WHERE P.Username = ?", username).
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
	rows, err := app.HabbGo().Database.Query("SELECT P.Badge FROM PlayerBadges P WHERE P.PlayerID = ?", player.Details.Id)
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
	rows, err := app.HabbGo().Database.Query("SELECT P.ID FROM Players P WHERE P.Username = ?", username)
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
	query := "SELECT P.ID, P.Username, P.Sex, P.Figure, P.PoolFigure, P.Film, P.Credits, P.Tickets, P.Motto, " +
		"P.ConsoleMotto, P.DisplayBadge, P.LastOnline, P.SoundEnabled " +
		"FROM Players P " +
		"WHERE P.username = ?"

	err := app.HabbGo().Database.QueryRow(query, p.Details.Username).Scan(&p.Details.Id, &p.Details.Username,
		&p.Details.Sex, &p.Details.Figure, &p.Details.PoolFigure, &p.Details.Film, &p.Details.Credits, &p.Details.Tickets,
		&p.Details.Motto, &p.Details.ConsoleMotto, &p.Details.DisplayBadge, &p.Details.LastOnline, &p.Details.SoundEnabled)

	if err != nil {
		log.Printf("%v ", err) // TODO log database errors properly
	}
}
