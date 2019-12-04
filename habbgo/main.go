package main

import (
	"github.com/jtieri/HabbGo/habbgo/server"
	"github.com/jtieri/HabbGo/habbgo/utils"
	"log"
)

func main() {
	//log.Println(string(encoding.EncodeB64(6, 3)))
	log.Println("Booting up HabbGo... ")
	config := utils.LoadConfig()
	gameServer := server.New(config.Server.Port, config.Server.Host, config.Server.MaxConns)
	gameServer.Start()
	defer gameServer.Stop()
}
