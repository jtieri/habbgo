package types

type Service interface {
	Start()
}

type Proxy interface {
	Init()
}

type ServiceProxies interface {
	RoomService() Proxy
	NavigatorService() Proxy
	ItemService() Proxy
	PlayerService() Proxy
}

type RoomState interface {
	GetInstanceID() int
	Type() string
	MoveTo(x, y int)
}
