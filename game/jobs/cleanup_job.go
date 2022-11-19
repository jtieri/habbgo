package jobs

import (
	"context"
	"time"

	"github.com/jtieri/habbgo/game/room"
)

const (
	cleanupJobInterval = 60 * time.Second // TODO make this configurable
	cleanupJobName     = "Cleanup Job"
)

type CleanupJob struct {
	ctx    context.Context
	cancel context.CancelFunc

	interval time.Duration
	ticker   *time.Ticker
	running  bool
	name     string

	room *room.Room
}

func NewCleanupJob(ctx context.Context, cancel context.CancelFunc, room *room.Room) CleanupJob {
	return CleanupJob{
		ctx:      ctx,
		cancel:   cancel,
		interval: cleanupJobInterval,
		ticker:   time.NewTicker(cleanupJobInterval),
		running:  false,
		name:     cleanupJobName,
		room:     room,
	}
}

func (cj *CleanupJob) Execute() {
	cj.running = true
	ticker := time.NewTicker(cj.interval)
	for {
		select {
		case <-ticker.C:
			// duration has passed, execute logic for this job
			cj.runJob()
		case <-cj.ctx.Done():
			// the job was cancelled, clean up appropriately and stop this job
			cj.running = false
			return
		}
	}
}

func (cj *CleanupJob) Stop() {
	cj.cancel()
}

func (cj *CleanupJob) Running() bool {
	return cj.running
}

func (cj *CleanupJob) Name() string {
	return cj.name
}

func (cj *CleanupJob) runJob() {
	// We only want to run this job if there are no players in the room.
	if cj.room.PlayerCount() > 0 {
		return
	}

	// TODO reset item states?

	cj.room.Initialized = false
	cj.room.StopRoomJobs()

	// TODO clear items, players, votes, etc.
	cj.room.ClearPlayers()
}
