package player

import (
	"log"
)

func LoginDB(player *Player, username string, password string) bool {
	var pw, name string
	err := player.Session.Database().QueryRow("SELECT P.Password, P.Username FROM Player P WHERE P.username = ?", username).Scan(&pw, &name)

	if err != nil {
		log.Printf("%v ", err) // TODO log database errors properly
	}

	if password == pw {
		player.Details.Username = name
		fillDetails(player)
		return true
	}

	return false
}

func LoadBadges(player *Player) {
	rows, err := player.Session.Database().Query("SELECT P.Badge FROM PlayerBadges P WHERE P.PlayerID = ?", player.Details.Id)
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
	rows, err := p.Session.Database().Query("SELECT P.ID FROM Player P WHERE P.Username = ?", username)
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

func fillDetails(p *Player) {
	query := "SELECT P.ID, P.Username, P.Sex, P.Figure, P.PoolFigure, P.Film, P.Credits, P.Tickets, P.Motto, " +
		"P.ConsoleMotto, P.DisplayBadge, P.LastOnline, P.SoundEnabled " +
		"FROM Players P " +
		"WHERE P.username = ?"

	err := p.Session.Database().QueryRow(query, p.Details.Username).Scan(&p.Details.Id, &p.Details.Username,
		&p.Details.Sex, &p.Details.Figure, &p.Details.PoolFigure, &p.Details.Film, &p.Details.Credits,
		&p.Details.Tickets, &p.Details.Motto, &p.Details.ConsoleMotto, &p.Details.CurrentBadge, &p.Details.DisplayBadge,
		&p.Details.LastOnline, &p.Details.SoundEnabled)

	if err != nil {
		log.Printf("%v ", err) // TODO log database errors properly
	}
}
