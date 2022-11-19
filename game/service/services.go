package service

import (
	"github.com/jtieri/habbgo/game/item"
	"github.com/jtieri/habbgo/game/navigator"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/room"
	"github.com/jtieri/habbgo/game/types"
)

type Proxies struct {
	Rooms     *room.RoomServiceProxy
	Items     *item.ItemServiceProxy
	Navigator *navigator.NavigatorServiceProxy
	Players   *player.PlayerServiceProxy
}

func (p *Proxies) RoomService() types.Proxy {
	return p.Rooms
}
func (p *Proxies) ItemService() types.Proxy {
	return p.Items
}

func (p *Proxies) NavigatorService() types.Proxy {
	return p.Navigator
}

func (p *Proxies) PlayerService() types.Proxy {
	return p.Players
}
