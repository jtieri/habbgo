package server

import (
	"sync"

	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/protocol/commands"
	"github.com/jtieri/habbgo/protocol/packets"
)

// Router maps incoming packet header ID's to their appropriate Command handlers.
type  Router struct {
	RegisteredCommands map[int]func(player.Player, packets.IncomingPacket)
	mutex              sync.RWMutex
}

// GetCommand returns the Command handler function associated with the specified headerId,
// if there is no registered Command with that headerId false is returned.
func (r *Router) GetCommand(headerId int) (func(player.Player, packets.IncomingPacket), bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	h, found := r.RegisteredCommands[headerId]
	return h, found
}

// RegisterCommands initializes the Router and registers the Command handler functions.
func RegisterCommands() (r *Router) {
	r = &Router{
		RegisteredCommands: make(map[int]func(p player.Player, packet packets.IncomingPacket)),
	}

	r.RegisterHandshakeCommands()
	r.RegisterRegistrationCommands()
	r.RegisterPlayerCommands()
	r.RegisterNavigatorCommands()
	r.RegisterRoomCommands()
	r.RegisterRoomUserCommands()

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
	r.RegisteredCommands[150] = commands.NAVIGATE
	r.RegisteredCommands[151] = commands.GETUSERFLATCATS
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

// incoming packet header 230 is GET_GROUP_BADGES
func (r *Router) RegisterRoomCommands() {
	r.RegisteredCommands[182] = commands.GETINTERST
	r.RegisteredCommands[2] = commands.ROOM_DIRECTORY
	r.RegisteredCommands[126] = commands.GETROOMAD
	r.RegisteredCommands[60] = commands.G_HMAP
	r.RegisteredCommands[61] = commands.G_USRS
	r.RegisteredCommands[62] = commands.G_OBJS
	r.RegisteredCommands[64] = commands.G_STAT
	r.RegisteredCommands[115] = commands.GOAWAY
	r.RegisteredCommands[75] = commands.MOVE
	/*
			tCmds.setaProp(#room_directory, 2)
			  tCmds.setaProp("GETDOORFLAT", 28)
			  tCmds.setaProp("CHAT", 52)
			  tCmds.setaProp("SHOUT", 55)
			  tCmds.setaProp("WHISPER", 56)
			  tCmds.setaProp("QUIT", 53)
			  tCmds.setaProp("GOVIADOOR", 54)
			  tCmds.setaProp("TRYFLAT", 57)
			  tCmds.setaProp("GOTOFLAT", 59)
			  tCmds.setaProp("G_HMAP", 60)
			  tCmds.setaProp("G_USRS", 61)
			  tCmds.setaProp("G_OBJS", 62)
			  tCmds.setaProp("G_ITEMS", 63)
			  tCmds.setaProp("G_STAT", 64)
			  tCmds.setaProp("GETSTRIP", 65)
			  tCmds.setaProp("FLATPROPBYITEM", 66)
			  tCmds.setaProp("ADDSTRIPITEM", 67)
			  tCmds.setaProp("TRADE_UNACCEPT", 68)
			  tCmds.setaProp("TRADE_ACCEPT", 69)
			  tCmds.setaProp("TRADE_CLOSE", 70)
			  tCmds.setaProp("TRADE_OPEN", 71)
			  tCmds.setaProp("TRADE_ADDITEM", 72)
			  tCmds.setaProp("MOVESTUFF", 73)
			  tCmds.setaProp("SETSTUFFDATA", 74)
			  tCmds.setaProp("MOVE", 75)
			  tCmds.setaProp("THROW_DICE", 76)
			  tCmds.setaProp("DICE_OFF", 77)
			  tCmds.setaProp("PRESENTOPEN", 78)
			  tCmds.setaProp("LOOKTO", 79)
			  tCmds.setaProp("CARRYDRINK", 80)
			  tCmds.setaProp("INTODOOR", 81)
			  tCmds.setaProp("DOORGOIN", 82)
			  tCmds.setaProp("G_IDATA", 83)
			  tCmds.setaProp("SETITEMDATA", 84)
			  tCmds.setaProp("REMOVEITEM", 85)
			  tCmds.setaProp("CARRYITEM", 87)
			  tCmds.setaProp("STOP", 88)
			  tCmds.setaProp("USEITEM", 89)
			  tCmds.setaProp("PLACESTUFF", 90)
			  tCmds.setaProp("DANCE", 93)
			  tCmds.setaProp("WAVE", 94)
			  tCmds.setaProp("KICKUSER", 95)
			  tCmds.setaProp("ASSIGNRIGHTS", 96)
			  tCmds.setaProp("REMOVERIGHTS", 97)
			  tCmds.setaProp("LETUSERIN", 98)
			  tCmds.setaProp("REMOVESTUFF", 99)
			  tCmds.setaProp("GOAWAY", 115)
			  tCmds.setaProp("GETROOMAD", 126)
			  tCmds.setaProp("GETPETSTAT", 128)
			  tCmds.setaProp("SETBADGE", 158)
			  tCmds.setaProp("GETINTERST", 182)
			  tCmds.setaProp("CONVERT_FURNI_TO_CREDITS", 183)
			  tCmds.setaProp("ROOM_QUEUE_CHANGE", 211)
			  tCmds.setaProp("SETITEMSTATE", 214)
			  tCmds.setaProp("GET_SPECTATOR_AMOUNT", 216)
			  tCmds.setaProp("GET_GROUP_BADGES", 230)
			  tCmds.setaProp("GET_GROUP_DETAILS", 231)
			  tCmds.setaProp("SPIN_WHEEL_OF_FORTUNE", 247)
			  tCmds.setaProp("RATEFLAT", 261)

		tMsgs.setaProp(-1, #handle_disconnect)
		  tMsgs.setaProp(18, #handle_clc)
		  tMsgs.setaProp(19, #handle_opc_ok)
		  tMsgs.setaProp(24, #handle_chat)
		  tMsgs.setaProp(25, #handle_chat)
		  tMsgs.setaProp(26, #handle_chat)
		  tMsgs.setaProp(28, #handle_users)
		  tMsgs.setaProp(29, #handle_logout)
		  tMsgs.setaProp(30, #handle_OBJECTS)
		  tMsgs.setaProp(31, #handle_heightmap)
		  tMsgs.setaProp(32, #handle_activeobjects)
		  tMsgs.setaProp(33, #handle_error)
		  tMsgs.setaProp(34, #handle_status)
		  tMsgs.setaProp(41, #handle_flat_letin)
		  tMsgs.setaProp(45, #handle_items)
		  tMsgs.setaProp(42, #handle_room_rights)
		  tMsgs.setaProp(43, #handle_room_rights)
		  tMsgs.setaProp(46, #handle_flatproperty)
		  tMsgs.setaProp(47, #handle_room_rights)
		  tMsgs.setaProp(48, #handle_idata)
		  tMsgs.setaProp(62, #handle_doorflat)
		  tMsgs.setaProp(63, #handle_doordeleted)
		  tMsgs.setaProp(64, #handle_doordeleted)
		  tMsgs.setaProp(69, #handle_room_ready)
		  tMsgs.setaProp(70, #handle_youaremod)
		  tMsgs.setaProp(71, #handle_showprogram)
		  tMsgs.setaProp(76, #handle_no_user_for_gift)
		  tMsgs.setaProp(83, #handle_items)
		  tMsgs.setaProp(84, #handle_removeitem)
		  tMsgs.setaProp(85, #handle_updateitem)
		  tMsgs.setaProp(88, #handle_stuffdataupdate)
		  tMsgs.setaProp(89, #handle_door_out)
		  tMsgs.setaProp(90, #handle_dice_value)
		  tMsgs.setaProp(91, #handle_doorbell_ringing)
		  tMsgs.setaProp(92, #handle_door_in)
		  tMsgs.setaProp(93, #handle_activeobject_add)
		  tMsgs.setaProp(94, #handle_activeobject_remove)
		  tMsgs.setaProp(95, #handle_activeobject_update)
		  tMsgs.setaProp(98, #handle_stripinfo)
		  tMsgs.setaProp(99, #handle_removestripitem)
		  tMsgs.setaProp(101, #handle_stripupdated)
		  tMsgs.setaProp(102, #handle_youarenotallowed)
		  tMsgs.setaProp(103, #handle_othernotallowed)
		  tMsgs.setaProp(105, #handle_trade_completed)
		  tMsgs.setaProp(108, #handle_trade_items)
		  tMsgs.setaProp(109, #handle_trade_accept)
		  tMsgs.setaProp(110, #handle_trade_close)
		  tMsgs.setaProp(112, #handle_trade_completed)
		  tMsgs.setaProp(129, #handle_presentopen)
		  tMsgs.setaProp(131, #handle_flatnotallowedtoenter)
		  tMsgs.setaProp(140, #handle_stripinfo)
		  tMsgs.setaProp(208, #handle_roomad)
		  tMsgs.setaProp(210, #handle_petstat)
		  tMsgs.setaProp(219, #handle_heightmapupdate)
		  tMsgs.setaProp(228, #handle_userbadge)
		  tMsgs.setaProp(230, #handle_slideobjectbundle)
		  tMsgs.setaProp(258, #handle_interstitialdata)
		  tMsgs.setaProp(259, #handle_roomqueuedata)
		  tMsgs.setaProp(254, #handle_youarespectator)
		  tMsgs.setaProp(283, #handle_removespecs)
		  tMsgs.setaProp(266, #handle_figure_change)
		  tMsgs.setaProp(298, #handle_spectator_amount)
		  tMsgs.setaProp(309, #handle_group_badges)
		  tMsgs.setaProp(310, #handle_group_membership_update)
		  tMsgs.setaProp(311, #handle_group_details)
		  tMsgs.setaProp(345, #handle_room_rating)
	*/
}

func (r *Router) RegisterRoomUserCommands() {
	r.RegisteredCommands[88] = commands.STOP
}
