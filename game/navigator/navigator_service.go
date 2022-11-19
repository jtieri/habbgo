package navigator

import (
	"context"
	"database/sql"
	"sync"

	"github.com/jtieri/habbgo/game/scheduler"
	"github.com/jtieri/habbgo/game/service/query"
	"go.uber.org/zap"
)

const channelBufferSize = 100

type NavigatorService struct {
	ctx    context.Context
	cancel context.CancelFunc

	repo      NavRepo
	scheduler *scheduler.GameScheduler
	navigator Navigator
	channels  *ServiceChannels
	running   bool

	log *zap.Logger
}

type ServiceChannels struct {
	CategoryByIDChan       chan *query.Request[int, Category]
	CategoryByParentIDChan chan *query.Request[int, []Category]
	CategoriesChan         chan *query.Request[int, []Category]
}

func (ns *NavigatorService) Channels() *ServiceChannels {
	return ns.channels
}

func NewNavigatorService(ctx context.Context, log *zap.Logger, db *sql.DB, scheduler *scheduler.GameScheduler, cancel context.CancelFunc) *NavigatorService {
	return &NavigatorService{
		ctx:    ctx,
		cancel: cancel,

		repo:      NewNavRepo(db),
		scheduler: scheduler,
		navigator: newNavigator(),
		channels:  newServiceChannels(),
		running:   false,

		log: log,
	}
}

func newServiceChannels() *ServiceChannels {
	return &ServiceChannels{
		CategoryByIDChan:       make(chan *query.Request[int, Category], channelBufferSize),
		CategoryByParentIDChan: make(chan *query.Request[int, []Category], channelBufferSize),
		CategoriesChan:         make(chan *query.Request[int, []Category], channelBufferSize),
	}
}

// Start retrieves the room categories from the database and builds the in-game Navigator with them.
func (ns *NavigatorService) Start() {
	cats, err := ns.repo.Categories()
	if err != nil {
		panic(err)
	}

	for _, cat := range cats {
		ns.navigator.categoryCache.Set(cat.ID, cat)
	}

	ns.log.Debug(
		"Loaded navigator categories from database",
		zap.Int("public_rooms_loaded", ns.navigator.categoryCache.Count()),
	)

	ns.running = true
	wg := &sync.WaitGroup{}

	for {
		// When the service starts we need to spin up a new goroutine that handles reading/writing
		// for one specific channel. This will allow us to concurrently handle requests from each channel at once.
		for _, handle := range ns.handlers() {
			go handle(wg)
			wg.Add(1)
		}

		// Block here until the context is cancelled and all the worker goroutines die.
		wg.Wait()

		// TODO finish cleaning up the navigator service and gracefully shutdown on context cancellation

		ns.running = false
		return
	}
}

func (ns *NavigatorService) handlers() []func(wg *sync.WaitGroup) {
	return []func(wg *sync.WaitGroup){
		ns.handleCategoryByID,
		ns.handleCategoriesByParentID,
		ns.handleCategories,
	}
}

func (ns *NavigatorService) handleCategoryByID(wg *sync.WaitGroup) {
	defer close(ns.channels.CategoryByIDChan)
	defer wg.Done()

	for {
		select {
		case <-ns.ctx.Done():
			return
		case req := <-ns.channels.CategoryByIDChan:
			ns.CategoryById(req)
		}
	}
}

func (ns *NavigatorService) handleCategoriesByParentID(wg *sync.WaitGroup) {
	defer close(ns.channels.CategoryByParentIDChan)
	defer wg.Done()

	for {
		select {
		case <-ns.ctx.Done():
			return
		case req := <-ns.channels.CategoryByParentIDChan:
			ns.CategoriesByParentId(req)
		}
	}
}

func (ns *NavigatorService) handleCategories(wg *sync.WaitGroup) {
	defer close(ns.channels.CategoriesChan)
	defer wg.Done()

	for {
		select {
		case <-ns.ctx.Done():
			return
		case req := <-ns.channels.CategoriesChan:
			ns.Categories(req)
		}
	}
}

// CategoryById retrieves a navigator category given the int parameter id and returns it if there is a match.
func (ns *NavigatorService) CategoryById(req *query.Request[int, Category]) {
	cat, ok := ns.navigator.categoryByID(req.Query.Key)
	if ok {
		req.Query.Value = cat
		req.Response <- req.Query
		return
	}

	req.Response <- nil
}

// CategoriesByParentId retrieves a slice of sub-categories given the int parameter pid and returns it if there is a match.
func (ns *NavigatorService) CategoriesByParentId(req *query.Request[int, []Category]) {
	cats := ns.navigator.categoryByParentID(req.Query.Key)
	req.Query.Value = cats
	req.Response <- req.Query
}

func (ns *NavigatorService) Categories(req *query.Request[int, []Category]) {
	req.Query.Value = ns.navigator.categories()
	req.Response <- req.Query
}
