package server

import (
	"github.com/jtieri/HabbGo/habbgo/game/player"
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
	case 228: // GET_SOUND_SETTING
		handlers.HandleGetSoundSetting(player, packet)
	case 315: // TEST_LATENCY ---> Init Latency Test
		handlers.HandleTestLatency(player, packet)

	// Navigator Packets ----------------------------------------------------------------------------------------------
	case 150: // NAVIGATE
		//handlers.HandleNavigate(player, packet)
	case 151: // GETUSERFLATCATS

	case 21: // GETFLATINFO

	case 23: // DELETEFLAT

	case 24: // UPDATEFLAT

	case 25: // SETFLATINFO

	case 13: // SBUSYF

	case 152: // GETFLATCAT

	case 153: // SETFLATCAT

	case 155: // REMOVEALLRIGHTS

	case 156: // GETPARENTCHAIN

	case 16: // SUSERF

	case 264: // GET_RECOMMENDED_ROOMS

	case 17: //SRCHF

	case 154: // GETSPACENODEUSERS

	case 18: // GETFVRF

	case 19: // ADD_FAVORITE_ROOM

	case 20: // DEL_FAVORITE_ROOM

	// ----------------------------------------------------------------------------------------------------------------
	default:
		log.Printf("No registered handler for packet [%v - %v], it's payload contained %v ",
			packet.Header, packet.HeaderId, packet.Payload.String())
	}
}
