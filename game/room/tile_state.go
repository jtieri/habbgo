package room

//go:generate go run golang.org/x/tools/cmd/stringer -type=TileState

// TileState represents the state of a Tile (e.g. is the tile occupied or is it unoccupied).
type TileState int

const (
	Accessible TileState = iota
	Inaccessible
)

// TileStates returns a slice of all the available TileState values.
func TileStates() []TileState {
	return []TileState{
		Accessible,
		Inaccessible,
	}
}
