package room

import (
	"github.com/jtieri/habbgo/collections"
	"github.com/jtieri/habbgo/game/item"
	"github.com/jtieri/habbgo/game/pathfinder/position"
	"github.com/jtieri/habbgo/game/player"
)

type Tile struct {
	Position position.Position
	Height   float64
	State    TileState

	players collections.Cache[int, player.Player]
	items   collections.Cache[int, item.Item]

	TopItem item.Item
}

func NewTile() Tile {
	return Tile{
		Position: position.Position{},
		Height:   0,
		State:    Inaccessible,
		players:  collections.NewCache(make(map[int]player.Player)),
		items:    collections.NewCache(make(map[int]item.Item)),
		TopItem:  item.Item{},
	}
}

func (t *Tile) AddPlayer(player player.Player) {
	t.players.SetIfAbsent(player.Details.Id, player)
}

func (t *Tile) RemovePlayer(playerID int) {
	t.players.Remove(playerID)
}

func (t *Tile) ContainsPlayer(playerID int) bool {
	return t.players.Has(playerID)
}

func (t *Tile) AddItem(item item.Item) {
	t.items.SetIfAbsent(item.ID, item)
}

func (t *Tile) RemoveItem(itemID int) {
	t.items.Remove(itemID)
}

func (t *Tile) ValidTile() bool {
	return true
}

func (t *Tile) ValidDiagonalTile() bool {
	return true
}

// ItemOnTopIsWalkable returns true if the item.Item at the top of the Tile is able to be walked on to.
func (t *Tile) ItemOnTopIsWalkable() bool {
	return t.TopItem.Walkable()
}

func (t *Tile) ItemOnTop() {
	var itemAtTop item.Item

	for _, i := range t.items.Items() {
		if i.Height() > t.Height {
			itemAtTop = i
		}
	}

	t.TopItem = itemAtTop
}

// WalkingHeight returns the height of the tile without the added height
// of beds and chairs.
func (t *Tile) WalkingHeight() float64 {
	height := t.Height

	if t.TopItem.Definition.ContainsBehavior(item.CanSitOnTop) || t.TopItem.Definition.ContainsBehavior(item.CanLayOnTop) {
		height = height - t.TopItem.Definition.TopHeight
	}

	return height
}

// DropInHeight returns true if the specified Tile's walking height is lower than the calling Tile's walking height.
func (t *Tile) DropInHeight(tile Tile) bool {
	return t.WalkingHeight() > tile.WalkingHeight()
}

// IncreaseInHeight returns true if the specified Tile's walking height is higher than the calling Tile's walking height.
func (t *Tile) IncreaseInHeight(tile Tile) bool {
	return t.WalkingHeight() < tile.WalkingHeight()
}

func ValidTile(position position.Position, room Room) bool {

	return true
}
