package database

import (
	"github.com/jtieri/HabbGo/habbgo/game/model"
	"log"
)

func Login(player *model.Player, username string, password string) bool {
	var pw, name string
	err := player.Session.Database().QueryRow("SELECT P.passwrd, P.username FROM Players P WHERE P.username = ?", username).Scan(&pw, &name)

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

func LoadBadges(player *model.Player) {
	rows, err := player.Session.Database().Query("SELECT P.badge FROM Players_Badges P WHERE P.pid = ?", player.Details.Id)
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

func fillDetails(p *model.Player) {
	query := "SELECT P.id, P.username, P.sex, P.figure, P.pool_figure, P.film, P.credits, P.tickets, P.motto, " +
		"P.console_motto, P.current_badge, P.display_badge, P.last_online, P.sound_enabled " +
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
