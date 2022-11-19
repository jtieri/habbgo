package player

import (
	"github.com/jtieri/habbgo/game/service/query"
)

type PlayerServiceProxy struct {
	service *ServiceChannels
}

func NewProxy(channels *ServiceChannels) *PlayerServiceProxy {
	return &PlayerServiceProxy{service: channels}
}

func (ps *PlayerServiceProxy) Init() {

}

func (ps *PlayerServiceProxy) AddPlayer(p Player) {
	req := query.NewRequest(p.Details.Username, p)
	ps.service.AddPlayer <- *req
	<-req.Response
}

func (ps *PlayerServiceProxy) RemovePlayer(p Player) {
	req := query.NewRequest(p.Details.Username, p)
	ps.service.RemovePlayer <- *req
	<-req.Response
}

func (ps *PlayerServiceProxy) UpdatePlayer(p Player) {
	req := query.NewRequest(p.Details.Username, p)
	ps.service.UpdatePlayer <- *req
	<-req.Response
}

func (ps *PlayerServiceProxy) GetPlayer(p Player) Player {
	req := query.NewRequest(p.Details.Username, p)
	ps.service.GetPlayer <- *req

	resp := <-req.Response
	if resp == nil {
		return Player{}
	}

	return resp.Value
}
