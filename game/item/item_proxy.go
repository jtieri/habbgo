package item

import (
	"fmt"

	"github.com/jtieri/habbgo/game/service/query"
)

// ItemServiceProxy provides a high level API for communicating with a running instance of item.ItemService.
type ItemServiceProxy struct {
	service *ServiceChannels
}

func NewProxy(channels *ServiceChannels) *ItemServiceProxy {
	return &ItemServiceProxy{
		service: channels,
	}
}

func (p *ItemServiceProxy) Init() {

}

func (p *ItemServiceProxy) ItemDefinition(id int) Definition {
	req := query.NewRequest(id, Definition{})
	p.service.DefinitionChan <- req

	resp := <-req.Response
	return resp.Value
}

func (p *ItemServiceProxy) PublicItems(roomID int, modelName string) []Item {
	pr := publicRoom{
		roomID:    roomID,
		modelName: modelName,
	}

	fmt.Println("Requesting public items")
	req := query.NewRequest(pr, []Item{})
	p.service.PublicItemChan <- req

	resp := <-req.Response
	fmt.Println("Got public items")
	return resp.Value
}
