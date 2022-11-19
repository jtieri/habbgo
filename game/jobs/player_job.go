package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/room"
	"github.com/jtieri/habbgo/protocol/messages"
)

const (
	playerJobInterval = 500 * time.Millisecond
	playerJobName     = "Player Job"
)

type PlayerJob struct {
	ctx    context.Context
	cancel context.CancelFunc

	interval time.Duration
	ticker   *time.Ticker
	running  bool
	name     string

	roomID  int
	service room.RoomServiceProxy
}

func NewPlayerJob(ctx context.Context, cancel context.CancelFunc, roomID int, service room.RoomServiceProxy) PlayerJob {
	return PlayerJob{
		ctx:      ctx,
		cancel:   cancel,
		interval: playerJobInterval,
		ticker:   time.NewTicker(playerJobInterval),
		running:  false,
		name:     playerJobName,
		roomID:   roomID,
		service:  service,
	}
}

func (pj PlayerJob) Execute() {
	pj.running = true
	ticker := time.NewTicker(pj.interval)
	for {
		select {
		case <-ticker.C:
			// duration has passed, execute logic for this job
			pj.runJob()
		case <-pj.ctx.Done():
			// the job was cancelled, clean up appropriately and stop this job
			fmt.Println("Cancelling Player Room Job")
			pj.running = false
			return
		}
	}
}

func (pj PlayerJob) Stop() {
	pj.cancel()
}

func (pj PlayerJob) Running() bool {
	return pj.running
}

func (pj PlayerJob) Name() string {
	return pj.name
}

func (pj PlayerJob) runJob() {
	// If there are no players in the room we don't need to do anything.
	if pj.service.RoomPlayerCount(pj.roomID) == 0 {
		return
	}

	var players []player.Player

	// Respect context cancellation mid job.
	select {
	case <-pj.ctx.Done():
		return
	default:

		for _, p := range pj.service.Players(pj.roomID) {
			// Sanity checks to make sure player is still in this room.
			if p.State().RoomID != pj.roomID {
				continue
			}

			// TODO apply pathfinding logic and handle player actions e.g. drinks, carry items etc.

			// If player needs update add to the slice of players we will send the STATUS message to.
			if p.NeedsUpdate() {
				p.SetUpdate(false)
				ps := p.Services.PlayerService().(*player.PlayerServiceProxy)
				ps.UpdatePlayer(p)

				// TODO need to update
				players = append(players, p)
			}
		}

		if len(players) > 0 {
			pj.service.RoomSendPacket(pj.roomID, messages.STATUS(players))
		}
	}
}
