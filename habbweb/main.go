package main

import (
	"github.com/jtieri/HabbGo/habbweb/server"
	"log"
)

func main() {
	log.Println("Starting the web server.... ")
	webServer := server.New()
	err := webServer.Start("127.0.0.1", 8080)
	if err != nil {
		log.Fatal("Failed to start the web server due to:  " + err.Error())
	}
}
