package server

import (
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/protocol/commands"
	"github.com/jtieri/habbgo/protocol/packets"
)

// Router maps incoming packet header ID's to their appropriate Command handlers.
type Router struct {
	RegisteredCommands map[int]func(*player.Player, *packets.IncomingPacket)
}

// GetCommand returns the Command handler function associated with the specified headerId,
// if there is no registered Command with that headerId false is returned.
func (r *Router) GetCommand(headerId int) (func(*player.Player, *packets.IncomingPacket), bool) {
	h, found := r.RegisteredCommands[headerId]
	return h, found
}

// RegisterCommands initializes the Router and registers the Command handler functions.
func RegisterCommands() (r *Router) {
	r = &Router{RegisteredCommands: make(map[int]func(p *player.Player, packet *packets.IncomingPacket))}

	r.RegisterHandshakeCommands()
	r.RegisterRegistrationCommands()
	r.RegisterPlayerCommands()
	r.RegisterNavigatorCommands()

	return
}

// RegisterHandshakeCommands registers the handshake related Command handlers.
func (r *Router) RegisterHandshakeCommands() {
	r.RegisteredCommands[206] = commands.INIT_CRYPTO
	r.RegisteredCommands[202] = commands.GENERATEKEY  // older clients
	r.RegisteredCommands[2002] = commands.GENERATEKEY // newer clients
	r.RegisteredCommands[5] = commands.VERSIONCHECK   // 1170 - VERSIONCHECK in later clients? v26+? // TODO figure out exact client revisions when these packet headers change
	r.RegisteredCommands[6] = commands.UNIQUEID
	r.RegisteredCommands[181] = commands.GET_SESSION_PARAMETERS
	r.RegisteredCommands[204] = commands.SSO
	r.RegisteredCommands[4] = commands.TRY_LOGIN
	r.RegisteredCommands[207] = commands.SECRETKEY
}

// RegisterRegistrationCommands registers the registration related Command handlers.
func (r *Router) RegisterRegistrationCommands() {
	r.RegisteredCommands[9] = commands.GETAVAILABLESETS
	r.RegisteredCommands[49] = commands.GDATE
	r.RegisteredCommands[42] = commands.APPROVENAME
	r.RegisteredCommands[203] = commands.APPROVE_PASSWORD
	r.RegisteredCommands[197] = commands.APPROVEEMAIL
	r.RegisteredCommands[43] = commands.REGISTER
}

// RegisterPlayerCommands registers the player related Command handlers.
func (r *Router) RegisterPlayerCommands() {
	r.RegisteredCommands[7] = commands.GET_INFO
	r.RegisteredCommands[8] = commands.GET_CREDITS
	r.RegisteredCommands[157] = commands.GETAVAILABLEBADGES
	r.RegisteredCommands[228] = commands.GET_SOUND_SETTING
	r.RegisteredCommands[315] = commands.TestLatency
}

// RegisterNavigatorCommands registers the Navigator related Command handlers.
func (r *Router) RegisterNavigatorCommands() {
	r.RegisteredCommands[150] = commands.Navigate
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
