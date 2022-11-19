package item

import (
	"github.com/jtieri/habbgo/game/pathfinder/position"
)

// DefaultTopHeight is the default top height used when building an Item's Definition.
const DefaultTopHeight = 0.001

const PresentDelimiter = "|"

// Item represents an in game item.
type Item struct {
	ID           int
	OrderID      int
	OwnerID      int
	RoomID       int
	TeleporterID int

	Definition Definition

	ItemAbove *Item
	ItemBelow *Item

	State ItemState

	WallPosition        string
	CustomData          string
	CurrentProgram      string
	CurrentProgramValue string

	RequiresUpdate     bool
	CurrentRollBlocked bool
	Hidden             bool
}

// Height returns the height of the item given its current
func (i *Item) Height() float64 {
	return i.State.Position.Z + i.Definition.TopHeight
}

func (i *Item) Walkable() bool {
	return true
}

func (i *Item) GateOpen() bool {
	if i.Definition.ContainsBehavior(Gate) {
		return i.CustomData == "O"
	}

	return false
}

type ItemState struct {
	InstanceID int
	StateType  string
	Position   position.Position
}

func (is *ItemState) GetInstanceID() int {
	return is.InstanceID
}

func (is *ItemState) Type() string {
	return is.StateType
}

func (is *ItemState) MoveTo(x, y int) {

}
