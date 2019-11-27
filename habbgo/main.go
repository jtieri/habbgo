package main

import (
	"github.com/jtieri/HabbGo/habbgo/server"
	"github.com/jtieri/HabbGo/habbgo/utils"
)

func main() {
	//log.Println(encoding.DecodeVl64([]byte{'@','A'}))
	//buffer := bytes.Buffer{}
	//packet := packets.IncomingPacket{"@@", 1, buffer}
	//packet.Payload.Write([]byte{'h','e','l','l','o'})
	//log.Println(string(packet.Bytes()))
	//log.Println(packet.Payload.String())


	config := utils.LoadConfig()
	gameServer := server.New(config.Server.Port, config.Server.Host)
	gameServer.Start()
	defer gameServer.Stop()
}
