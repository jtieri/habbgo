package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jtieri/HabbGo/habbgo/config"
	"github.com/jtieri/HabbGo/habbgo/server"
	"log"
	"strconv"
)

func main() {
	log.Println("Booting up HabbGo... ")

	log.Println("Loading config file... ")
	c := config.LoadConfig()

	log.Println("Attempting to make connection with the database... ")
	db, err := sql.Open("mysql", c.Database.User+":"+c.Database.Password+"@tcp"+
		"("+c.Database.Host+":"+strconv.Itoa(int(c.Database.Port))+")"+"/"+c.Database.Name)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database %v at %v:%v %v", c.Database.Name, c.Database.Host, c.Database.Port, err)
	}
	defer db.Close()
	log.Printf("Successfully connected to database %v at %v:%v ", c.Database.Name, c.Database.Host, c.Database.Port)

	log.Println("Starting the game server... ")
	gameServer := server.New(&c, db)
	gameServer.Start()
	defer gameServer.Stop()
}
