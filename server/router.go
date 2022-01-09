package server

import (
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/protocol/commands"
	"github.com/jtieri/habbgo/protocol/packets"
)

type Router struct {
	RegisteredPackets map[int]func(*player.Player, *packets.IncomingPacket)
}

func (r *Router) GetHandler(headerId int) (func(*player.Player, *packets.IncomingPacket), bool) {
	h, found := r.RegisteredPackets[headerId]
	return h, found
}

func RegisterHandlers() (r *Router) {
	r = &Router{RegisteredPackets: make(map[int]func(p *player.Player, packet *packets.IncomingPacket))}

	r.RegisterHandshakeHandlers()
	r.RegisterRegistrationHandlers()
	r.RegisterPlayerHandlers()
	r.RegisterNavigatorHandlers()

	return
}

func (r *Router) RegisterHandshakeHandlers() {
	r.RegisteredPackets[206] = commands.INIT_CRYPTO
	r.RegisteredPackets[202] = commands.GENERATEKEY  // older clients
	r.RegisteredPackets[2002] = commands.GENERATEKEY // newer clients
	r.RegisteredPackets[5] = commands.VERSIONCHECK   // 1170 - VERSIONCHECK in later clients? v26+? // TODO figure out exact client revisions when these packet headers change
	r.RegisteredPackets[6] = commands.UNIQUEID
	r.RegisteredPackets[181] = commands.GET_SESSION_PARAMETERS
	r.RegisteredPackets[204] = commands.SSO
	r.RegisteredPackets[4] = commands.TRY_LOGIN
	r.RegisteredPackets[207] = commands.SECRETKEY
}

func (r *Router) RegisterRegistrationHandlers() {
	r.RegisteredPackets[9] = commands.GETAVAILABLESETS
	r.RegisteredPackets[49] = commands.GDATE
	r.RegisteredPackets[42] = commands.APPROVENAME
	r.RegisteredPackets[203] = commands.APPROVE_PASSWORD
	r.RegisteredPackets[197] = commands.APPROVEEMAIL
	r.RegisteredPackets[43] = commands.REGISTER
}

func (r *Router) RegisterPlayerHandlers() {
	r.RegisteredPackets[7] = commands.GET_INFO
	r.RegisteredPackets[8] = commands.GET_CREDITS
	r.RegisteredPackets[157] = commands.GETAVAILABLEBADGES
	r.RegisteredPackets[228] = commands.GET_SOUND_SETTING
	r.RegisteredPackets[315] = commands.TestLatency
}

func (r *Router) RegisterNavigatorHandlers() {
	r.RegisteredPackets[150] = commands.Navigate
	// 151: GETUSERFLATCATS
	// 21: GETFLATINFO
	// 23: DELETEFLAT
	// 24: UPDATEFLAT
	// 25: SETFLATINFO
	// 13: SBUSYF
	// 152: GETFLATCAT
	// 153: SETFLATCAT
	// 155: REMOVEALLRIGHTS
	// 156: GETPARENTCHAIN
	// 16: SUSERF
	// 264: GET_RECOMMENDED_ROOMS
	// 17: SRCHF
	// 154: GETSPACENODEUSERS
	// 18: GETFVRF
	// 19: ADD_FAVORITE_ROOM
	// 20: DEL_FAVORITE_ROOM
}
