package commands

import (
	"context"
	"fmt"

	"github.com/jtieri/habbgo/game/jobs"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/room"
	"github.com/jtieri/habbgo/game/scheduler"
	"github.com/jtieri/habbgo/protocol/messages"
	"github.com/jtieri/habbgo/protocol/packets"
)

// GETINTERST is sent from the client when attempting to enter a room from the Navigator.
// It attempts to get interstitial data used when displaying the room loading UI.
func GETINTERST(player player.Player, packet packets.IncomingPacket) {
	player.Session.Send(messages.INTERSTITIALDATA, messages.INTERSTITIALDATA())
}

// ROOM_DIRECTORY is sent from the client when trying to enter a room.
// It will determine if the room is public or private and perform the necessary
// checks before initiating the server side logic for entering a room and,
// sending the response packets.
func ROOM_DIRECTORY(player player.Player, packet packets.IncomingPacket) {
	publicRoom := packet.ReadBytes(1)
	roomID := packet.ReadInt()
	doorID := packet.ReadInt()
	_ = doorID // TODO find out what doorID is for and then remove this

	fmt.Println("At top")

	// Check if the room is a private room and if so send OPEN_CONNECTION_OK
	if string(publicRoom) != "A" {
		player.Session.Send(messages.OPEN_CONNECTION_OK, messages.OPEN_CONNECTION_OK())
		return
	}

	roomService := player.Services.RoomService().(*room.RoomServiceProxy)
	r := roomService.RoomByID(roomID)
	if !r.Ready {
		fmt.Println()
		fmt.Println("Room was not ready.")
		fmt.Println()
		return
	}

	fmt.Println("After get room by ID")
	// TODO check if r is Habbo Club only and if user has Habbo Club membership

	// Check if the r is currently full
	if r.Details.CurrentVisitors >= r.Details.MaxVisitors {
		// TODO send CANTCONNECT indicating r is full
		fmt.Println()
		fmt.Println("Room was full")
		fmt.Println()
		return
	}

	fmt.Println("Before enter room")
	// enter r
	// this should initialize the r and perform all server side logic like
	// adding the player to the r and adjusting player count etc.
	if err := enterRoom(roomService, r, player); err != nil {
		// TODO possibly log here somehow? need better insight into errors
		fmt.Println()
		fmt.Println("Error on enter room")
		fmt.Println()
		return
	}
	fmt.Println("After enter room")
}

func GETROOMAD(player player.Player, packet packets.IncomingPacket) {
	player.Session.Send(messages.ROOMAD, messages.ROOMAD())
}

func G_HMAP(player player.Player, packet packets.IncomingPacket) {
	proxy := player.Services.RoomService().(*room.RoomServiceProxy)
	room := proxy.RoomByID(player.State().RoomID)

	player.Session.Send(messages.HEIGHTMAP, messages.HEIGHTMAP(room.Model.Heightmap))
}

func G_USRS(player player.Player, packet packets.IncomingPacket) {
	if !player.InRoom() {
		return
	}

	proxy := player.Services.RoomService().(*room.RoomServiceProxy)
	room := proxy.RoomByID(player.State().RoomID)

	player.Session.Send(messages.USER_OBJECTS, messages.USER_OBJECTS(room.Players()))
}

func G_OBJS(player player.Player, packet packets.IncomingPacket) {
	if !player.InRoom() {
		return
	}

	proxy := player.Services.RoomService().(*room.RoomServiceProxy)
	publicItems := proxy.RoomPublicItems(player.State().RoomID)
	floorItems := proxy.RoomFloorItems(player.State().RoomID)

	fmt.Printf("G_OBJS - loaded %d public items \n", len(publicItems))

	player.Session.Send(messages.OBJECTS, messages.OBJECTS(publicItems))
	player.Session.Send(messages.ACTIVE_OBJECTS, messages.ACTIVE_OBJECTS(floorItems))
}

func G_STAT(player player.Player, packet packets.IncomingPacket) {
	// This packet is sent from the client to refresh the room state when a player is attempting to enter the room.
	// Private rooms will need to account for setting rights in the case that the player has room rights.

	// Send USER_OBJECTS to every player in the room to update the client with the new player who is currently entering the room.
	// Send USER_OBJECTS to the player entering the room. This may already be done and not necessary so double check.

	// Send UESR_STATUSES with room.Entities as argument
	// Make sure EntityAction is ran for this player.

	// roomPlayer := player.room.(*room.RoomPlayer)

	// For every item in the room send SHOWPROGRAM which I believe updates the items state and applies animations.
}

func GOAWAY(p player.Player, packet packets.IncomingPacket) {
	roomService := p.Services.RoomService().(*room.RoomServiceProxy)

	roomService.LeaveRoom(p.State().RoomID, p)

	// TODO possibly more to do here for resetting a players state when leaving a room

	/*
				Position doorLocation = this.getRoom().getModel().getDoorLocation();

		        if (doorLocation == null) {
		            this.getRoom().getEntityManager().leaveRoom(this.entity, true);
		            return;
		        }

		        // If we're standing in the door, immediately leave room
		        if (this.getPosition().equals(doorLocation)) {
		            this.getRoom().getEntityManager().leaveRoom(this.entity, true);
		            return;
		        }

		        // Attempt to walk to the door
		        this.walkTo(doorLocation.getX(), doorLocation.getY());
		        this.isWalkingAllowed = allowWalking;
		        this.beingKicked = true;

		        // If user isn't walking, leave immediately
		        if (!this.isWalking) {
		            this.getRoom().getEntityManager().leaveRoom(this.entity, true);
		        }
	*/
}

func MOVE(player player.Player, packet packets.IncomingPacket) {
	if !player.CanMove() {
		return
	}

	xCoord := packet.ReadB64()
	yCoord := packet.ReadB64()

	roomService := player.Services.RoomService().(*room.RoomServiceProxy)

	// TODO remove indiretion after refactor away from pointers
	roomService.MovePlayer(player.State().RoomID, room.NewPlayerMove(player, xCoord, yCoord))
}

func enterRoom(rs *room.RoomServiceProxy, r room.Room, player player.Player) error {
	ctx, cancel := context.WithCancel(player.Ctx)
	roomStartupJobs := []scheduler.Job{jobs.NewPlayerJob(ctx, cancel, r.Details.Id, *player.Services.RoomService().(*room.RoomServiceProxy))}

	fmt.Println("Before room service proxy call EnterRoom")
	// TODO remove indiretion after refactor away from pointers
	rs.EnterRoom(r.Details.Id, room.NewPlayerEnter(player, roomStartupJobs))

	fmt.Println("After room service proxy call EnterRoom")
	// TODO handle teleporters in the future here since doorPos may be different than the room models door position.

	player.Session.Send(messages.ROOM_READY, messages.ROOM_READY(r.Details.Id, r.Model.Name))

	if r.Details.Wallpaper > 0 {
		player.Session.Send(messages.FLATPROPERTY, messages.FLATPROPERTY(room.WallpaperProperty, r.Details.Wallpaper))
	}

	if r.Details.Floor > 0 {
		player.Session.Send(messages.FLATPROPERTY, messages.FLATPROPERTY(room.FloorProperty, r.Details.Floor))
	}

	// TODO check if this user is the room owner or if they have voted on the room already
	// if so don't let them vote in next send packet call
	// send ROOM_RATING with appropriate value
	voted := false

	if voted {
		player.Session.Send(messages.ROOM_RATING, messages.ROOM_RATING(r.Details.Rating))
	} else {
		player.Session.Send(messages.ROOM_RATING, messages.ROOM_RATING(-1))
	}

	fmt.Println("After sending all packets in enterRoom")
	// TODO send room event info when room events are implemented

	// TODO save new room visitor count to the database
	return nil
}

func STOP(player player.Player, packet packets.IncomingPacket) {
	/*
	 mapping stopWhat = reader.contents();

	        if (stopWhat.equals("Dance")) {
	            player.getRoomUser().removeStatus(StatusType.DANCE);
	            player.getRoomUser().setNeedsUpdate(true);
	        }

	        if (stopWhat.equals("CarryItem")) {
	            player.getRoomUser().removeStatus(StatusType.CARRY_ITEM);
	            player.getRoomUser().setNeedsUpdate(true);
	        }

	        player.getRoomUser().getTimerManager().resetRoomTimer();
	*/
}
