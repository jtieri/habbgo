package app

import (
	"database/sql"
	"fmt"
	"github.com/jtieri/habbgo/config"
)

var habbgo *App

type App struct {
	Config   *config.Config
	Database *sql.DB
}

func New(cfg *config.Config, db *sql.DB) {
	if habbgo == nil {
		habbgo = &App{
			Config:   cfg,
			Database: db,
		}
	}
}

func Habbgo() *App {
	return habbgo
}

func LogErr(err error) {
	fmt.Printf("Habbgo encountered error: %s \n", err)
}
