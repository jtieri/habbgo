package server

import (
	"github.com/jtieri/HabbGo/habbgo/game/player"
	"github.com/jtieri/HabbGo/habbgo/protocol/handlers"
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
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
	r.RegisterPlayerHandlers()
	r.RegisterNavigatorHandlers()

	return
}

func (r *Router) RegisterHandshakeHandlers() {
	r.RegisteredPackets[206] = handlers.InitCrypto
	r.RegisteredPackets[202] = handlers.GenerateKey  // older clients
	r.RegisteredPackets[2002] = handlers.GenerateKey // newer clients
	// 207 - SECRETKEY
	// 5 - VERSIONCHECK in older clients
	// 1170 - VERSIONCHECK in later clients? v26+?
	// TODO figure out exact client revisions when these packet headers change
	// 6 - UNIQUEID
	r.RegisteredPackets[181] = handlers.GetSessionParams
	r.RegisteredPackets[204] = handlers.SSO
	r.RegisteredPackets[4] = handlers.TryLogin
}

func (r *Router) RegisterPlayerHandlers() {
	r.RegisteredPackets[7] = handlers.GetInfo
	r.RegisteredPackets[8] = handlers.GetCredits
	r.RegisteredPackets[157] = handlers.GetAvailableBadges
	r.RegisteredPackets[228] = handlers.GetSoundSetting
	r.RegisteredPackets[315] = handlers.TestLatency
}

func (r *Router) RegisterNavigatorHandlers() {
	r.RegisteredPackets[150] = handlers.Navigate
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
