package room

import (
	"github.com/jtieri/habbgo/game/pathfinder/position"
)

// heightmapDelimiter is used to mark the end of a line in the text based representation of a Model's Heightmap.
const heightmapDelimiter = "|"

// Model represents a Room's model data.
type Model struct {
	id        int
	Name      string
	Door      Door
	Heightmap string
}

// Door represents the entrypoint into a Room.
type Door struct {
	X         int
	Y         int
	Z         float64
	Direction int
}

// Position returns the pathfinder.Position representing the coordinates of a Room's entrypoint.
func (d Door) Position() position.Position {
	return position.Position{
		X:            d.X,
		Y:            d.Y,
		Z:            d.Z,
		BodyRotation: d.Direction,
		HeadRotation: d.Direction,
	}
}

// TODO may need to add something for handling public room actions
// e.g. Handling actions for special rooms like the pools, game rooms, infobus, etc
