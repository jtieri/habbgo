package server

import (
	"github.com/jtieri/HabbGo/habbgo/game/model/player"
	"github.com/jtieri/HabbGo/habbgo/protocol/handlers"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
	"log"
)

func Handle(player *player.Player, packet *packets.IncomingPacket) {
	switch packet.HeaderId {
	// Handshake Packets ----------------------------------------------------------------------------------------------
	case 206: // INIT_CRYPTO
		handlers.HandleInitCrypto(player, packet)
	case 202: // GENERATEKEY
		handlers.HandleGenerateKey(player, packet)
	case 207: // SECRETKEY

	case 5: // VERSIONCHECK

	case 6: // UNIQUEID

	case 181: // GET_SESSION_PARAMETERS
		handlers.HandleGetSessionParams(player, packet)
	case 204: // SSO
		handlers.HandleSSO(player, packet)
	case 4: // TRY LOGIN - used when SSO is disabled
		handlers.HandleTryLogin(player, packet)

	// Player Packets -------------------------------------------------------------------------------------------------
	case 7: // GET_INFO
		handlers.HandleGetInfo(player, packet)
	case 8: // GET_CREDITS
		handlers.HandleGetCredits(player, packet)
	case 157: // GETAVAILABLEBADGES
		handlers.HandleGetAvailableBadges(player, packet)
	case 228: // Get Sound setting

	case 315: // TEST_LATENCY ---> Init Latency Test

	// ----------------------------------------------------------------------------------------------------------------
	default:
		log.Printf("No registered handler for packet [%v - %v], it's payload contained %v ",
			packet.Header, packet.HeaderId, packet.Payload.String())
	}
}
