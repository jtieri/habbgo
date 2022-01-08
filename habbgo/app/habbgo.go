package app

import (
	"database/sql"
	"fmt"
	"github.com/jtieri/HabbGo/habbgo/config"
)

var habbGo *App

type App struct {
	Config   *config.Config
	Database *sql.DB
}

func New(cfg *config.Config, db *sql.DB) {
	if habbGo == nil {
		habbGo = &App{
			Config:   cfg,
			Database: db,
		}
	}
}

func HabbGo() *App {
	return habbGo
}

func LogErr(err error) {
	fmt.Printf("HabbGo encountered error: %s \n", err)
}
