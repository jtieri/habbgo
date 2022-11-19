package navigator

import "github.com/jtieri/habbgo/game/service/query"

// NavigatorServiceProxy provides a high level API for communicating with a running instance of navigator.NavigatorService.
type NavigatorServiceProxy struct {
	service *ServiceChannels
}

func NewProxy(channels *ServiceChannels) *NavigatorServiceProxy {
	return &NavigatorServiceProxy{service: channels}
}

func (np *NavigatorServiceProxy) Init() {

}

func (np *NavigatorServiceProxy) CategoryByID(catID int) Category {
	req := query.NewRequest(catID, Category{})
	np.service.CategoryByIDChan <- req

	resp := <-req.Response
	if resp == nil {
		return Category{}
	}

	return resp.Value
}

func (np *NavigatorServiceProxy) CategoriesByParentID(parentID int) []Category {
	req := query.NewRequest(parentID, []Category{})
	np.service.CategoryByParentIDChan <- req

	resp := <-req.Response
	return resp.Value
}

func (np *NavigatorServiceProxy) Categories() []Category {
	req := query.NewRequest(0, []Category{})
	np.service.CategoriesChan <- req

	resp := <-req.Response
	return resp.Value
}
