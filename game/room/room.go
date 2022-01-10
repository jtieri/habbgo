package room

import "time"

type Room struct {
	Details *Details
	Model   *Model
}

type Details struct {
	Id              int
	CategoryID      int
	Name            string
	Description     string
	CCTs            string
	Wallpaper       int
	Floor           int
	Landscape       float32
	OwnerId         int
	OwnerName       string
	ShowOwner       bool
	SudoUsers       bool
	Hidden          bool
	AccessType      Access
	Password        string
	CurrentVisitors int
	MaxVisitors     int
	Rating          int
	ChildRooms      []*Room
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewRoom() *Room {
	return &Room{
		Details: &Details{},
		Model:   &Model{},
	}
}

type Access int

const (
	Open Access = iota
	Closed
	Password
)

func (a Access) String() string {
	switch a {
	case Open:
		return "open"
	case Closed:
		return "closed"
	case Password:
		return "password"
	default:
		return "open"
	}
}

func AccessType(accessString string) Access {
	switch accessString {
	case "open":
		return Open
	case "closed":
		return Closed
	case "password":
		return Password
	default:
		return Open
	}
}

type Model struct {
	ID            int
	Name          string
	DoorX         int
	DoorY         int
	DoorZ         float32
	DoorDirection int
	Heightmap     string
}

type Map struct {
	mapping [][]*Tile
}

type Tile struct {
	// May incorporate this into pathfinding later down the line
}
