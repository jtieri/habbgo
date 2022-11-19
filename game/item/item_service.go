package item

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/jtieri/habbgo/collections"
	"github.com/jtieri/habbgo/game/scheduler"
	"github.com/jtieri/habbgo/game/service/query"
	"go.uber.org/zap"
)

const channelBufferSize = 100

type ItemService struct {
	ctx    context.Context
	cancel context.CancelFunc

	repo            ItemRepo
	scheduler       *scheduler.GameScheduler
	definitionCache collections.Cache[int, Definition]
	publicItemCache collections.Cache[int, publicItem]
	channels        *ServiceChannels
	running         bool

	log *zap.Logger
}

func NewItemService(ctx context.Context, log *zap.Logger, db *sql.DB, scheduler *scheduler.GameScheduler, cancel context.CancelFunc) *ItemService {
	return &ItemService{
		ctx:    ctx,
		cancel: cancel,

		repo:            NewItemRepo(db),
		scheduler:       scheduler,
		definitionCache: collections.NewCache(make(map[int]Definition)),
		publicItemCache: collections.NewCache(make(map[int]publicItem)),
		channels:        newServiceChannel(),
		running:         false,

		log: log,
	}
}

// ServiceChannels is a wrapper type for all the channels needed to send and receive requests from/to the service.
type ServiceChannels struct {
	DefinitionChan chan *query.Request[int, Definition]
	PublicItemChan chan *query.Request[publicRoom, []Item]
}

// newServiceChannel creates a new ServiceChannels object with all of its channels initialized.
func newServiceChannel() *ServiceChannels {
	return &ServiceChannels{
		DefinitionChan: make(chan *query.Request[int, Definition], channelBufferSize),
		PublicItemChan: make(chan *query.Request[publicRoom, []Item], channelBufferSize),
	}
}

// Start will load the appropriate data and perform the necessary actions to setup the ItemService on startup.
func (is *ItemService) Start() {
	definitions, err := is.repo.LoadItemDefinitions()
	if err != nil {
		panic(err)
	}

	for _, d := range definitions {
		is.definitionCache.Set(d.ID, d)
	}

	is.log.Debug(
		"Loaded item definitions from database",
		zap.Int("definitions_loaded", is.definitionCache.Count()),
	)

	publicItems, err := is.repo.LoadPublicRoomItemData()
	if err != nil {
		panic(err)
	}

	for _, i := range publicItems {
		is.publicItemCache.Set(i.id, i)
	}

	is.log.Debug(
		"Loaded public room item data from database",
		zap.Int("public_room_items_loaded", is.publicItemCache.Count()),
	)

	is.running = true
	wg := &sync.WaitGroup{}
	for {
		// When the service starts we need to spin up a new goroutine that handles reading/writing
		// for one specific channel. This will allow us to concurrently handle requests from each channel at once.
		for _, handle := range is.handlers() {
			go handle(wg)
			wg.Add(1)
		}

		// Block here until the context is cancelled and all the worker goroutines die.
		wg.Wait()

		// TODO finish gracefully closing out the item service.

		is.running = false
		return
	}
}

func (is *ItemService) Channels() *ServiceChannels {
	return is.channels
}

func (is *ItemService) handlers() []func(*sync.WaitGroup) {
	return []func(wg *sync.WaitGroup){
		is.handleDefinition,
		is.handlePublicItems,
	}
}

func (is *ItemService) handleDefinition(wg *sync.WaitGroup) {
	defer close(is.channels.DefinitionChan)
	defer wg.Done()

	for {
		select {
		case <-is.ctx.Done():
			return
		case req := <-is.channels.DefinitionChan:
			is.Definition(req)
		}
	}
}

func (is *ItemService) handlePublicItems(wg *sync.WaitGroup) {
	defer close(is.channels.PublicItemChan)
	defer wg.Done()

	for {
		select {
		case <-is.ctx.Done():
			return
		case req := <-is.channels.PublicItemChan:
			is.PublicItems(req)
		}
	}
}

func (is *ItemService) Definition(req *query.Request[int, Definition]) {
	def, ok := is.definitionCache.Get(req.Query.Key)
	if ok {
		req.Query.Value = def
	}
	req.Response <- req.Query
}

type publicRoom struct {
	roomID    int
	modelName string
}

// PublicItems builds a slice of item.Item objects
func (is *ItemService) PublicItems(req *query.Request[publicRoom, []Item]) {
	fmt.Println("Inside PublicItems")
	// Find all the cached public room item definitions for the specified room model.
	var itemData []publicItem
	for _, i := range is.publicItemCache.Items() {
		if i.roomModel == req.Query.Key.modelName {
			itemData = append(itemData, i)
		}
	}

	var (
		items   []Item
		usedIDs []string
	)

	// Build an Item object for every public item definition being requested.
	for _, data := range itemData {

		// Generate a new unique ID for this public item.
		customID := randomPublicID(usedIDs)
		usedIDs = append(usedIDs, customID)

		// Build Item from the public item data.
		newItem := data.build(req.Query.Key.roomID, customID)

		// TODO finish adding behaviors & interactions
		// Add appropriate behaviors to the item.
		if !strings.Contains(newItem.Definition.Sprite, "queue_tile2") {
			newItem.Definition.AddBehavior(PublicSpaceObject)
		}

		items = append(items, newItem)
	}

	fmt.Println("Leaving PublicItems")
	fmt.Printf("There were %d public items for this room \n", len(items))
	req.Query.Value = items
	req.Response <- req.Query
}
