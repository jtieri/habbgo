package main

import (
	"database/sql"
	"fmt"
	"github.com/jtieri/habbgo/config"
	"github.com/jtieri/habbgo/game/navigator"
	"github.com/jtieri/habbgo/game/room"
	"github.com/jtieri/habbgo/server"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	log.Println("Booting up habbgo... ")

	log.Println("Loading config file... ")
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Attempting to make connection with the database... ")
	host := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)

	db, err := sql.Open(c.DBDriver, host)
	if err != nil {
		log.Fatal(err)
	}

	// Check that the connection to the DB is alive
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database %v at %v:%v %v", c.DBName, c.DBHost, c.DBPort, err)
	}
	defer db.Close()
	log.Printf("Successfully connected to database %v at %v:%v ", c.DBName, c.DBHost, c.DBPort)

	log.Println("Preparing the in game services...")
	navigatorService := navigator.NavigatorService()
	navigatorService.SetDBCon(db)
	navigatorService.BuildNavigator()

	roomService := room.RoomService()
	roomService.SetDBConn(db)

	log.Println("Starting the game server... ")
	gameServer := server.New(c, db)
	gameServer.Start()

	defer gameServer.Stop()
}
