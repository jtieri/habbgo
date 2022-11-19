package room

import "errors"

var (
	ErrTileNotExist     = errors.New("tile at specified (x, y) coordinate does not exist in the room map")
	ErrIndexOutOfBounds = errors.New("(x, y) coordinate pair is outside of valid range for room map")
	ErrTileInaccessible = errors.New("tile at specified (x, y) coordinate is not accessible")
)

// Map represents the in-game map for a Room.
type Map struct {
	tiles [][]Tile
	sizeX int
	sizeY int
}

func (m *Map) CollisionMap() {

}

func buildTileMap(sizeX, sizeY int) [][]Tile {
	matrix := make([][]Tile, sizeX)
	for x := 0; x < sizeX; x++ {
		matrix[x] = make([]Tile, sizeY)
		for y := 0; y < sizeY; y++ {
			matrix[x][y] = NewTile()
		}
	}
	return matrix
}
