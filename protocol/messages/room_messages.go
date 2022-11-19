package messages

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jtieri/habbgo/game/item"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/num"
	"github.com/jtieri/habbgo/protocol/packets"
	"github.com/jtieri/habbgo/text"
)

// INTERSTITIALDATA is sent from the server as a response to commands.GETINTERST.
// It's payload can contain two strings where one is a URL to an image to be used
// during the loading room UI and the second one is a URL to a destination users will go if the image is clicked.
func INTERSTITIALDATA() packets.OutgoingPacket {
	p := packets.NewOutgoing(258)
	/*
		This is the client lingo code for handling this message.
		We should be checking if there is an interstitial set for loading and if so send
		the URL to the image and the URL to the destination site.

		if tMsg.content.length > 1 then
		    tDelim = the itemDelimiter
		    the itemDelimiter = "\t"
		    tSourceURL = tMsg.content.getProp(#item, 1)
		    tTargetURL = tMsg.content.getProp(#item, 2)
		    the itemDelimiter = tDelim
		    me.getComponent().getInterstitial().Init(tSourceURL, tTargetURL)
	*/
	return p
}

// OPEN_CONNECTION_OK is sent from the server as a response to commands.ROOM_DIRECTORY.
// It will confirm to the client that the room the player is currently trying to enter is
// actually a private room.
func OPEN_CONNECTION_OK() packets.OutgoingPacket {
	return packets.NewOutgoing(19) // @S
}

// ROOM_READY is sent from the server as one of many packets in response to commands.ROOM_DIRECTORY.
// It sends the appropriate room ID and room model name to the client for initializing the room.
func ROOM_READY(roomID int, modelName string) packets.OutgoingPacket {
	p := packets.NewOutgoing(69) // AE
	p.WriteString(modelName)
	p.WriteString(" ")
	p.WriteInt(roomID)
	return p
}

// FLATPROPERTY is sent from the server as one of the many packets in response to commands.ROOM_DIRECTORY.
// It contains a room property and value which is used in updating the wallpaper and floor elements
// in the room.
func FLATPROPERTY(property string, value int) packets.OutgoingPacket {
	p := packets.NewOutgoing(46) // @n
	p.WriteDelim([]byte(property), []byte("/"))
	p.Write(strconv.Itoa(value))
	return p
}

// ROOM_RATING is sent from the server to update a rooms ratings in the client.
func ROOM_RATING(rating int) packets.OutgoingPacket {
	p := packets.NewOutgoing(345) //
	p.WriteInt(rating)
	/*
		TODO
		Lingo code from client where it reads the data from this packet.
		Looks like this actually sends the room ratings and another number called roomRatingPercent
		I would imagine we store roomTotalVotes & roomPositiveVotes then calculate the percentage of positive votes

		  tRoomRating = tConn.GetIntFrom()
		  tRoomRatingPercent = tConn.GetIntFrom()
	*/
	return p
}

// ROOMAD is sent from the server as a response to commands.GETROOMAD.
// It will send the rooms ad image URL and the target URL for where users end up if they click the ad.
func ROOMAD() packets.OutgoingPacket {
	p := packets.NewOutgoing(208) // CP

	/*
		TODO this needs to be implemented properly still. Here is the lingo from the client where this packet is processed
			if tMsg.content.length > 1 then
			    tDelim = the itemDelimiter
			    the itemDelimiter = "\t"
			    tSourceURL = tMsg.content.getProp(#item, 1)
			    tTargetURL = tMsg.content.getProp(#item, 2)
			    the itemDelimiter = tDelim
			    tLayoutID = me.getInterface().getRoomVisualizer().pLayout
			    me.getComponent().getAd().Init(tSourceURL, tTargetURL, tLayoutID)
			  else
			    me.getComponent().getAd().Init(0)
			  end if
	*/

	return p
}

// HEIGHTMAP is sent from the server as a response to commands.G_HMAP.
// It sends the room's heightmap as a string to the client to for rendering.
func HEIGHTMAP(heightmap string) packets.OutgoingPacket {
	p := packets.NewOutgoing(31) // @_
	p.Write(heightmap)
	return p
}

// USER_OBJECTS is sent from the server as a response to commands.G_USRS.
// It serializes the state for each entity in the current room and sends the data to the client for rendering.
func USER_OBJECTS(players []player.Player) packets.OutgoingPacket {
	p := packets.NewOutgoing(28) // @\
	p.Write("\r")

	// TODO this will need revised to account for bots and pets in the future
	for _, plyr := range players {
		state := plyr.State()
		p.WriteKeyValue([]byte("i"), []byte(strconv.Itoa(state.InstanceID)))
		p.WriteKeyValue([]byte("a"), []byte(strconv.Itoa(plyr.Details.Id)))
		p.WriteKeyValue([]byte("n"), []byte(plyr.Details.Username))
		p.WriteKeyValue([]byte("f"), []byte(plyr.Details.Figure))
		p.WriteKeyValue([]byte("s"), []byte(plyr.Details.Sex))
		p.WriteKeyValue([]byte("l"), []byte(fmt.Sprintf("%d %d %.2f", state.Position.X, state.Position.Y, num.Round(state.Position.Z))))

		if len(plyr.Details.Motto) > 0 {
			p.WriteKeyValue([]byte("c"), []byte(plyr.Details.Motto))
		}

		if plyr.Details.DisplayBadge {
			p.WriteKeyValue([]byte("b"), []byte(plyr.Details.CurrentBadge))
		}

		// If the room is one of the public poolrooms then be sure to write the players pool figure
		if strings.HasPrefix(state.ModelName, "pool_") || state.ModelName == "md_a" {
			if plyr.Details.PoolFigure != "" {
				p.WriteKeyValue([]byte("p"), []byte(plyr.Details.PoolFigure))
			}
		}

		// the client checks for [bot] on L239 in the room handler class, from the client src code,
		// to determine if the entity being rendered is a bot.
		if state.StateType == "bot" {
			p.WriteDelim([]byte("[bot]"), []byte(string(rune(13))))
		}
	}

	return p
}

// OBJECTS is sent from the server as a response to commands.G_OBJS.
// It serializes the public room items in the current room and sends them to the client for rendering.
func OBJECTS(items []item.Item) packets.OutgoingPacket {
	p := packets.NewOutgoing(30) // @^

	if items != nil {
		p.WriteInt(len(items))
		SerializeItems(items, p)
	}

	return p
}

// ACTIVE_OBJECTS is sent from the server as a response to commands.G_OBJS.
// It serializes the floor items in the current room and sends them to the client for rendering.
func ACTIVE_OBJECTS(items []item.Item) packets.OutgoingPacket {
	p := packets.NewOutgoing(32) // @`
	p.WriteInt(len(items))
	SerializeItems(items, p)
	return p
}

func ITEMS() packets.OutgoingPacket {
	p := packets.NewOutgoing(45) // @m

	return p
}

func STATUS(players []player.Player) packets.OutgoingPacket {
	p := packets.NewOutgoing(34) // @b

	for _, plyr := range players {
		pState := plyr.State()
		p.WriteDelim(num.IntToBytes(pState.InstanceID), []byte(" "))
		p.WriteDelim(num.IntToBytes(pState.Position.X), []byte(","))
		p.WriteDelim(num.IntToBytes(pState.Position.Y), []byte(","))
		p.WriteDelim([]byte(num.Float64ToString(pState.Position.Z)), []byte(","))
		p.WriteDelim(num.IntToBytes(pState.Position.HeadRotation), []byte(","))
		p.WriteDelim(num.IntToBytes(pState.Position.BodyRotation), []byte("/"))

		for _, action := range pState.Actions {
			p.Write(action.Name)

			if action.Params > 0 {
				p.Write(" ")
				p.Write(strconv.Itoa(action.Params))
			}

			p.Write("/")
		}

		p.Write(string(rune(13)))
	}

	return p
}

// SerializeItems writes the Item data, in the appropriate format for the type of Definition,
// to the packets.OutgoingPacket.
// Used in the messages package when Item data must be sent over the wire.
func SerializeItems(items []item.Item, packet packets.OutgoingPacket) packets.OutgoingPacket {
	for _, i := range items {
		itemState := i.State

		switch {
		// handle public room items
		case i.Definition.ContainsBehavior(item.PublicSpaceObject):
			packet.WriteDelim([]byte(i.CustomData), []byte(" "))
			packet.WriteString(i.Definition.Sprite)
			packet.WriteDelim([]byte(strconv.Itoa(itemState.Position.X)), []byte(" "))
			packet.WriteDelim([]byte(strconv.Itoa(itemState.Position.Y)), []byte(" "))
			packet.WriteDelim([]byte(strconv.Itoa(int(itemState.Position.Z))), []byte(" "))
			packet.Write(strconv.Itoa(itemState.Position.BodyRotation))

			if i.Definition.ContainsBehavior(item.ExtraParameter) {
				packet.Write(" 2")
			}

			packet.Write(string(rune(13)))

		// handle wall item
		case i.Definition.ContainsBehavior(item.WallItem):
			packet.WriteDelim([]byte(strconv.Itoa(i.ID)), []byte(string(rune(9))))
			packet.WriteDelim([]byte(i.Definition.Sprite), []byte(string(rune(9))))
			packet.WriteDelim([]byte(" "), []byte(string(rune(9))))
			packet.WriteDelim([]byte(i.WallPosition), []byte(string(rune(9))))

			if i.CustomData != "" {
				if i.Definition.ContainsBehavior(item.PostIt) {
					packet.Write(text.Substr(i.CustomData, 0, 6)) // Get color of Post-It note
				} else {
					packet.Write(i.CustomData)
				}
			}

			packet.Write(string(rune(13)))

		// handle default case
		default:
			packet.WriteString(strconv.Itoa(i.ID))
			packet.WriteString(i.Definition.Sprite)
			packet.WriteInt(itemState.Position.X)
			packet.WriteInt(itemState.Position.Y)
			packet.WriteInt(i.Definition.Length)
			packet.WriteInt(i.Definition.Width)
			packet.WriteInt(itemState.Position.BodyRotation)
			packet.WriteString(fmt.Sprintf("%.2f", num.Round(itemState.Position.Z)))
			packet.WriteString(i.Definition.Color)
			packet.WriteString("")

			// In order for animations to work with rollers we have to write a 2 here.
			if i.Definition.ContainsBehavior(item.Roller) {
				packet.WriteInt(2)
			} else {
				packet.WriteInt(0)
			}

			// Write the item's custom data appropriately if it is a present.
			if i.Definition.ContainsBehavior(item.Present) {
				presentData := strings.Split(i.CustomData, item.PresentDelimiter)
				if len(presentData) >= 3 {
					packet.WriteString("!" + presentData[2])
				} else {
					packet.WriteString("")
				}
			} else {
				packet.WriteString(i.CustomData)
			}
		}
	}
	return packet
}
