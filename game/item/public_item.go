package item

import (
	"github.com/jtieri/habbgo/game/pathfinder/position"
	"github.com/jtieri/habbgo/text"
)

// publicItem represents a public room item.
type publicItem struct {
	id             int
	inRoomID       string
	roomModel      string
	sprite         string
	x              int
	y              int
	z              float64
	rotation       int
	topHeight      float64
	length         int
	width          int
	behavior       Behavior
	currentProgram string
	teleportTo     string
	swimTo         string
}

// build returns an Item built with the public item's data.
func (p *publicItem) build(roomID int, customID string) Item {
	return Item{
		ID:     p.id,
		RoomID: roomID,
		Definition: Definition{
			Sprite:      p.sprite,
			TopHeight:   p.topHeight,
			Length:      p.length,
			Width:       p.width,
			Behaviors:   Behaviors{p.behavior},
			Interaction: nil,
		},
		State: ItemState{
			Position: position.Position{
				X:            p.x,
				Y:            p.y,
				Z:            p.z,
				BodyRotation: p.rotation,
				HeadRotation: p.rotation,
			},
		},
		CustomData:     customID,
		CurrentProgram: p.currentProgram,
	}
}

// randomPublicID generates a new unique ID to be used for public items.
func randomPublicID(usedIDs []string) string {
	id := ""

	// Build a unique random ID for each public furni item.
	for id == "" {
		tmpID := text.RandomString(20)

		// Ensure that the newly generated ID is not already in use.
		idAlreadyUsed := false
		for _, i := range usedIDs {
			if tmpID == i {
				idAlreadyUsed = true
				break
			}
		}

		// If the ID is not already in use then we can use it.
		if idAlreadyUsed {
			continue
		}

		id = tmpID
	}

	return id
}
