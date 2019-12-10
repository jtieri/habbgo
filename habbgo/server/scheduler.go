package server

import (
	"github.com/jtieri/HabbGo/habbgo/game/model/player"
	"github.com/jtieri/HabbGo/habbgo/server/protocol/handlers"
	"github.com/jtieri/HabbGo/habbgo/server/protocol/packets"
	"log"
)

func Handle(player *player.Player, packet *packets.IncomingPacket) {
	switch packet.HeaderId {
	// Handshake Packets
	case 206:
		handlers.HandleInitCrypto(player, packet)
	case 202:
		handlers.HandleGenerateKey(player, packet)
	case 207: // SECRETKEY

	case 5: // VERSIONCHECK

	case 6: // UNIQUEID

	case 181:
		handlers.HandleGetSessionParams(player, packet)
	case 204: // SSO
		// // TODO Login user if credentials are correct if not send disconnect from server
		// Send LOGIN_OK Header:3
	case 4: // TRY LOGIN - used when SSO is disabled
		// TODO Login user if credentials are correct if not send LOCALISED_ERROR("Incorrect Login Details.")
		// Send LOGIN_OK Header:3
	case 7: // GET_INFO
		// Send USER_OBJ Header:5
	case 8: // GET_CREDITS
		// Send CREDIT_BALANCE
	case 157: // Get Badges

	case 228: // Get Sound setting

	case 315: // TEST_LATENCY ---> Init Latency Test




	default:
		log.Printf("No registered handler for packet [%v - %v], it's payload contained %v ",
			packet.Header, packet.HeaderId, packet.Payload.String())
	}
}
