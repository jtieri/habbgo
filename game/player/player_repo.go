package player

import (
	"database/sql"
	"log"

	"github.com/jtieri/habbgo/crypto"
	"github.com/jtieri/habbgo/game/ranks"
	"go.uber.org/zap"
)

type PlayerRepo struct {
	database *sql.DB
}

// NewPlayerRepo returns a new instance of PlayerRepo which Player's utilize for accessing the database.
func NewPlayerRepo(db *sql.DB) PlayerRepo {
	return PlayerRepo{database: db}
}

func (pr *PlayerRepo) Register(username, figure, gender, email, birthday, createdAt, password string, salt []byte) error {
	stmt, err := pr.database.Prepare(
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

func (pr *PlayerRepo) LoginDB(player Player, username string, password string) bool {
	var (
		psswrdHash, uname string
		psswrdSalt        []byte
	)

	err := pr.database.QueryRow(
		"SELECT P.password_hash, P.password_salt, P.username FROM Players P WHERE P.username = $1", username).
		Scan(&psswrdHash, &psswrdSalt, &uname)

	if err != nil {
		player.log.Warn("Failed to query database during login",
			zap.String("username", username),
			zap.Error(err),
		)
	}

	if crypto.HashPassword(password, psswrdSalt) == psswrdHash {
		player.Details.Username = uname
		fillDetails(&player)
		return true
	}

	return false
}

func (pr *PlayerRepo) LoadBadges(player *Player) {
	rows, err := pr.database.Query("SELECT P.badge_id FROM player_badges P WHERE P.player_id = $1", player.Details.Id)
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

func (pr *PlayerRepo) PlayerExists(username string) bool {
	rows, err := pr.database.Query("SELECT P.id FROM Players P WHERE P.username = $1", username)
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
	query := "SELECT P.id, P.username, P.sex, P.figure, P.pool_figure, P.film, P.credits, P.tickets, P.motto, " +
		"P.console_motto, P.last_online, P.sound_enabled, P.Rank " +
		"FROM Players P " +
		"WHERE P.username = $1"

	var tmpRank int
	err := p.Repo.database.QueryRow(query, p.Details.Username).Scan(&p.Details.Id, &p.Details.Username,
		&p.Details.Sex, &p.Details.Figure, &p.Details.PoolFigure, &p.Details.Film, &p.Details.Credits,
		&p.Details.Tickets, &p.Details.Motto, &p.Details.ConsoleMotto, &p.Details.LastOnline, &p.Details.SoundEnabled,
		&tmpRank)

	if err != nil {
		log.Printf("%v ", err) // TODO log database errors properly
	}

	p.Details.PlayerRank = ranks.Rank(tmpRank)
}
