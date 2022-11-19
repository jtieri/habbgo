package room

import (
	"database/sql"
	"strconv"
	"strings"
	"unicode"
)

type RoomRepo struct {
	database *sql.DB
}

func NewRoomRepo(db *sql.DB) RoomRepo {
	return RoomRepo{database: db}
}

func (rr *RoomRepo) LoadPublicRooms() ([]Room, error) {
	stmt, err := rr.database.Prepare(
		"SELECT r.*, rm.* FROM rooms r LEFT JOIN room_categories rc ON r.category_id = rc.id LEFT JOIN room_models rm on r.model_id = rm.id WHERE rc.is_public=true",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		room, err := fillRoomData(rows)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (rr *RoomRepo) RoomsByPlayerId(id int) ([]Room, error) {
	stmt, err := rr.database.Prepare("SELECT * FROM rooms r LEFT JOIN room_models rm on r.model_id = rm.id WHERE r.owner_id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		room, err := fillRoomData(rows)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, err
}

func (rr *RoomRepo) RoomByID(roomID int) (Room, error) {
	stmt, err := rr.database.Prepare("SELECT * FROM rooms r LEFT JOIN room_models rm on r.model_id = rm.id WHERE r.id = $1")
	if err != nil {
		return Room{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(roomID)
	if err != nil {
		return Room{}, err
	}
	defer rows.Close()

	var room Room
	if rows.Next() {
		room, err = fillRoomData(rows)
		if err != nil {
			return Room{}, err
		}
	}

	return room, nil
}

func fillRoomData(rows *sql.Rows) (Room, error) {
	var tmpAccessType string

	room := NewRoom()

	err := rows.Scan(
		&room.Details.Id,
		&room.Details.CategoryID,
		&room.Details.Name,
		&room.Details.Description,
		&room.Details.OwnerId,
		&room.Model.id,
		&room.Details.CCTs,
		&room.Details.Wallpaper,
		&room.Details.Floor,
		&room.Details.ShowOwner,
		&room.Details.Password,
		&tmpAccessType,
		&room.Details.SudoUsers,
		&room.Details.CurrentVisitors,
		&room.Details.MaxVisitors,
		&room.Details.Rating,
		&room.Details.Hidden,
		&room.Details.CreatedAt,
		&room.Details.UpdatedAt,
		&room.Model.id,
		&room.Model.Name,
		&room.Model.Door.X,
		&room.Model.Door.Y,
		&room.Model.Door.Z,
		&room.Model.Door.Direction,
		&room.Model.Heightmap,
	)
	if err != nil {
		return Room{}, err
	}

	// find the right Access for the string representation of an Access
	// that we loaded from the database.
	for _, access := range AccessTypes() {
		if strings.ToLower(access.String()) == tmpAccessType {
			room.Details.AccessType = access
		}
	}

	// When rooms are loaded from the database we want to be sure that we are building their map
	// from the room model's heightmap.
	room, err = parseHeightMap(room)
	if err != nil {
		return Room{}, err
	}
	return room, nil
}

// parseHeightMap will read the string representing the Room's Heightmap and build the Room's Map.
// NOTE: This function must be called when Room's are loaded from the database because it alters the heightmap
//       in a way that is necessary for the client to be able to render it properly.
func parseHeightMap(r Room) (Room, error) {
	tmpHeightmapRows := strings.Split(r.Model.Heightmap, heightmapDelimiter)

	// We check every row to see if it's empty because if the heightmap string
	// ends with the delimiter for some reason, then we end up with an extra empty row,
	// and it causes index out of range errors in the nested for loop below.
	var heightmapRows []string
	for _, hm := range tmpHeightmapRows {
		if hm != "" {
			heightmapRows = append(heightmapRows, hm)
		}
	}

	mapSizeX := len(heightmapRows[0])
	mapSizeY := len(heightmapRows)
	r.mapping.sizeX = mapSizeX
	r.mapping.sizeY = mapSizeY
	r.mapping.tiles = buildTileMap(mapSizeX, mapSizeY)

	var sb strings.Builder

	// Iterate over every (x, y) coordinate by indexing into the string representation of the heightmap,
	// and determine if this tile is occupied or unoccupied and what the height of the tile is.
	for y := 0; y < mapSizeY; y++ {
		heightmapRow := heightmapRows[y]
		for x := 0; x < mapSizeX; x++ {
			tile := heightmapRow[x]

			// If the value at this tile is a numerical value we want to set the tile height using this value.
			// Otherwise, we just set the tile height to 0.
			if unicode.IsDigit(rune(tile)) {
				tileHeightVal, err := strconv.ParseFloat(string(tile), 64)
				if err != nil {
					return Room{}, err
				}
				r.mapping.tiles[x][y].Height = tileHeightVal
				r.mapping.tiles[x][y].State = Accessible
			} else {
				r.mapping.tiles[x][y].Height = 0
				r.mapping.tiles[x][y].State = Inaccessible
			}

			// Check if the current loop iterator pair (x, y) is equal to this
			// Room's door location and set the height and state appropriately.
			if x == r.Model.Door.X && y == r.Model.Door.Y {
				r.mapping.tiles[x][y].Height = r.Model.Door.Z
				r.mapping.tiles[x][y].State = Accessible
			}

			sb.WriteString(string(tile))
		}

		// add the delimiter for newlines, this is used by the client to render the heightmap.
		// NOTE: this is crucial for when we send the heightmap to the client in the G_HMAP message.
		sb.WriteString("\r")
	}
	r.Model.Heightmap = sb.String()
	return r, nil
}
