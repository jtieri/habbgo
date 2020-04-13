package model

type Room struct {
	Details *Data
	Model   *Model
	//Map     *Map
}

type Data struct {
	Id int
	CatId int
	Name string
	Desc string
	CCTs string
	Wallpaper int
	Floor int
	Landscape float32
	Owner_Id int
	Owner_Name string
	ShowOwner bool
	SudoUsers bool
	HideRoom bool
	AccessType int
	Password string
	CurrentVisitors int
	MaxVisitors int
	Rating int
	ChildRooms []*Room
}

type Model struct {
	Id            int
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
