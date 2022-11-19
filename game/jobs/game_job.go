package jobs

import (
	"context"
	"time"

	"github.com/jtieri/habbgo/game/item"
	"github.com/jtieri/habbgo/game/player"
)

const (
	gameJobName = "Game Job"
)

type GameJob struct {
	ctx    context.Context
	cancel context.CancelFunc

	interval time.Duration
	ticker   *time.Ticker
	running  bool
	name     string

	rewardsQueue    chan player.Player
	itemSaveQueue   chan item.Item
	itemDeleteQueue chan item.Item
}

func NewGameJob(ctx context.Context, cancel context.CancelFunc, interval time.Duration) *GameJob {
	return &GameJob{
		ctx:             ctx,
		cancel:          cancel,
		interval:        interval,
		ticker:          time.NewTicker(interval),
		running:         false,
		name:            gameJobName,
		rewardsQueue:    make(chan player.Player),
		itemSaveQueue:   make(chan item.Item),
		itemDeleteQueue: make(chan item.Item),
	}
}

func (gj *GameJob) Execute() {
	gj.running = true
	defer gj.ticker.Stop()
	for {
		select {
		case <-gj.ticker.C:
			// executing
		case <-gj.ctx.Done():
			// cancelled
			// TODO do anything before the job is stopped e.g. save state to db

			gj.running = false
			return
		}
	}
}

func (gj *GameJob) Stop() {
	gj.cancel()
}

func (gj *GameJob) Running() bool {
	return gj.running
}

func (gj *GameJob) Name() string {
	return gj.name
}
