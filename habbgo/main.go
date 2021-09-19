package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jtieri/HabbGo/habbgo/app"
	"github.com/jtieri/HabbGo/habbgo/config"
	"github.com/jtieri/HabbGo/habbgo/server"
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

	// Check that the connection to the DB is alive
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database %v at %v:%v %v", c.DB.Name, c.DB.Host, c.DB.Port, err)
	}
	defer db.Close()
	log.Printf("Successfully connected to database %v at %v:%v ", c.DB.Name, c.DB.Host, c.DB.Port)

	// Create the global App context for accessing Config and DB across the server
	app.New(c, db)

	log.Println("Starting the game server... ")
	gameServer := server.New()
	gameServer.Start()

	defer gameServer.Stop()
}
