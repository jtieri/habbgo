package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jtieri/HabbGo/habbgo/config"
	"github.com/jtieri/HabbGo/habbgo/game/navigator"
	"github.com/jtieri/HabbGo/habbgo/game/room"
	"github.com/jtieri/HabbGo/habbgo/server"
	"log"
)

func main() {
	log.Println("Booting up HabbGo... ")

	log.Println("Loading config file... ")
	c := config.LoadConfig("config.yml")

	log.Println("Attempting to make connection with the database... ")
	host := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name)

	db, err := sql.Open("mysql", host)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database %v at %v:%v %v", c.DB.Name, c.DB.Host, c.DB.Port, err)
	}
	defer db.Close()
	log.Printf("Successfully connected to database %v at %v:%v ", c.DB.Name, c.DB.Host, c.DB.Port)

	log.Printf("Setting up in-game services and models...")
	navigator.NavigatorService().SetDBCon(db)
	navigator.NavigatorService().BuildNavigator()

	room.RoomService().SetDBConn(db)

	log.Println("Starting the game server... ")
	gameServer := server.New(c, db)
	gameServer.Start()

	defer gameServer.Stop()
}
