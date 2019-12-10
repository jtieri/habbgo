package main

import (
	"github.com/jtieri/HabbGo/habbgo/server"
	"github.com/jtieri/HabbGo/habbgo/utils"
	"log"
)

func main() {
	//log.Println(string(encoding.EncodeB64(206, 2)))
	log.Println("Booting up HabbGo... ")
	config := utils.LoadConfig()
	gameServer := server.New(&config)
	gameServer.Start()
	defer gameServer.Stop()
}
