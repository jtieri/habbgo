package scheduler

import (
	"context"

	"go.uber.org/zap"
)

// Scheduler describes the functions needed to process a new Job that must be executed on some interval of time.
type Scheduler interface {
	ScheduleJob(job Job)
	Start()
	Stop()
}

// Job describes the functions needed to run some process on an interval of time.
type Job interface {
	Execute()
	Stop()
	Running() bool
	Name() string
}

// GameScheduler processes new Job's that need to be executed on a regular time interval.
// GameScheduler runs on its own goroutine and has new jobs passed to it over a buffered channel.
type GameScheduler struct {
	ctx    context.Context
	cancel context.CancelFunc

	jobQueue chan Job

	log *zap.Logger
}

// NewGameScheduler instantiates a new Scheduler for processing Job's.
func NewGameScheduler(ctx context.Context, cancel context.CancelFunc, log *zap.Logger) *GameScheduler {
	return &GameScheduler{
		ctx:      ctx,
		cancel:   cancel,
		jobQueue: make(chan Job, 100), // TODO make the buffer size configurable
		log:      log,
	}
}

// ScheduleJob passes a new Job into the Job queue to be executed by the GameScheduler.
func (gs *GameScheduler) ScheduleJob(job Job) {
	gs.jobQueue <- job
}

// Start will begin listening for new Job's that need to be executed,
// while also listening for the context to be cancelled so that it can gracefully shut down.
func (gs *GameScheduler) Start() {
	for {
		select {
		case j := <-gs.jobQueue:
			// when a new job is scheduled we read from the channel and execute the job
			go j.Execute()
		case <-gs.ctx.Done():
			// cancelled
			// TODO wait for running tasks to finish? propagate cancel call to running jobs?
			close(gs.jobQueue)
			for job := range gs.jobQueue {
				job.Stop()
			}
			return
		}
	}
}

// Stop calls the GameScheduler's cancel function.
func (gs *GameScheduler) Stop() {
	gs.log.Info("Stopping the game scheduler")
	gs.cancel()
}
