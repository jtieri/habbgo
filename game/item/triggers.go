package item

import "strings"

var _ Interactioner = &DefaultInteraction{}

// Interactioner is an interface that is implemented by various item interaction types (e.g. beds, chairs, teleporter, etc).
type Interactioner interface {
	OnRoomEntry()
	OnRoomLeave()
	OnEntityLeave()
	OnEntityStep()
	OnEntityStop()
	OnItemPlaced()
	OnItemPickup()
}

/*  public void onRoomEntry(Entity entity, room room, boolean firstEntry, Object... customArgs) { }
    public void onItemPickup(Player player, room room, Item item) { }
    public void onRoomLeave(Entity entity, room room, Object... customArgs) { }
    public void onEntityStep(Entity entity, RoomState roomEntity, Item item, Position oldPosition) { }
    public void onEntityStop(Entity entity, RoomState roomEntity, Item item, boolean isRotation) { }
    public void onItemPlaced(Player player, room room, Item item) { }
    public void onItemMoved(Player player, room room, Item item, boolean isRotation, Position oldPosition, Item itemBelow, Item itemAbove) { }
    public void onEntityLeave(Entity entity, RoomState roomEntity, Item item) { }
*/

// interactionByName will return the appropriate Interactioner implementation for the specified name.
func interactionByName(name string) Interactioner {
	switch strings.ToLower(name) {
	case "default":
		return &DefaultInteraction{}
	default:
		return &DefaultInteraction{}
	}
}

// DefaultInteraction is the default Interactioner implementation
type DefaultInteraction struct {
}

func (d DefaultInteraction) OnRoomEntry() {
	//TODO implement me
	panic("implement me")
}

func (d DefaultInteraction) OnRoomLeave() {
	//TODO implement me
	panic("implement me")
}

func (d DefaultInteraction) OnEntityLeave() {
	//TODO implement me
	panic("implement me")
}

func (d DefaultInteraction) OnEntityStep() {
	//TODO implement me
	panic("implement me")
}

func (d DefaultInteraction) OnEntityStop() {
	//TODO implement me
	panic("implement me")
}

func (d DefaultInteraction) OnItemPlaced() {
	//TODO implement me
	panic("implement me")
}

func (d DefaultInteraction) OnItemPickup() {
	//TODO implement me
	panic("implement me")
}
