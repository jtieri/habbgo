package log

import (
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/protocol/packets"
)

func LogOutgoingPacket(username string, p *packets.OutgoingPacket) {
	log.Printf("[OUTGOING] [%v] [%v - %v]: %v ", username, p.Header, p.HeaderId, p.Payload.String())
}

func LogIncomingPacket(username string, handler func(*player.Player, *packets.IncomingPacket), p *packets.IncomingPacket) {
	hName := getHandlerName(runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
	log.Printf("[INCOMING] [%v] [%v - %v|%v]: %v ", username, hName, p.Header, p.HeaderId, p.Payload.String())
}

func LogUnknownPacket(username string, p *packets.IncomingPacket) {
	log.Printf("[UNKNOWN] [%v] [%v - %v]: %v ", username, p.Header, p.HeaderId, p.Payload.String())
}

func getHandlerName(handler string) string {
	sp := strings.Split(handler, "/") // e.g. github.com/jtieri/habbgo/protocol/handlers.GenerateKey
	s2 := sp[len(sp)-1]               // e.g. handlers.GenerateKey
	return strings.Split(s2, ".")[1]  // e.g. GenerateKey
}