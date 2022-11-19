package actor

import (
	"fmt"
	"testing"

	"github.com/jtieri/habbgo/collections"
)

const (
	taskQueueSize  = 600
	actorQueueSize = 300
)

func TestActor(t *testing.T) {
	fmt.Println("Hello World")

	c := make(chan string, 5)
	d := make(chan string, 5)
	shutdown := make(chan bool, 1)

	for i := 0; i < 5; i++ {
		c <- fmt.Sprintf("msg-c%d", i)
		d <- fmt.Sprintf("msg-d%d", i)
	}

	//go func() {
	//	time.Sleep(50 * time.Millisecond)
	//	shutdown <- true
	//}()
	shutdown <- true

	var stop bool
	for {
		select {
		case <-shutdown:
			stop = true
			close(c)
			close(d)
		case m := <-c:
			fmt.Println(m)
		case m := <-d:
			fmt.Println(m)
		}

		fmt.Println("After select")

		if stop {
			for msg := range c {
				fmt.Println(msg)
			}

			for msg := range d {
				fmt.Println(msg)
			}

			break
		}
	}

	fmt.Println("After for")

}

// DemoActorService handles the lifecycle of all Actor objects in the system.
// It will process incoming Tasks and pass them to the proper Actor.
type DemoActorService struct {
	// actorPool is a map of cached Actor objects that are running.
	actorPool collections.Cache[string, Actor]

	// addActorQueue listens for an incoming Actor to add to the pool.
	addActorQueue chan Actor

	// removeActorQueue listens for an incoming Actor to remove from the pool.
	removeActorQueue chan Actor

	// taskQueue listens for an incoming Task that should be delegated to a specific Actor in the pool.
	taskQueue chan Task

	// shutdown listens for incoming signals to shut down so the service can close gracefully.
	shutdown chan bool
}

func NewDemoActorService() *DemoActorService {
	return &DemoActorService{
		actorPool:        collections.NewCache(make(map[string]Actor)),
		addActorQueue:    make(chan Actor, actorQueueSize),
		removeActorQueue: make(chan Actor, actorQueueSize),
		taskQueue:        make(chan Task, taskQueueSize),
		shutdown:         make(chan bool, 1),
	}
}

func (ds DemoActorService) SubmitTask(t Task) {
	ds.taskQueue <- t
}

func (ds DemoActorService) handleActorTask(task Task) {
	actor, found := ds.actorPool.Get(task.ActorID())
	if found {
		actor.AddTask(task)
	}

	// If the Actor is not in the pool but the ActorService has inbound tasks for it
	// that means the tasks were likely queued before the Actor was stopped by some
	// task that arrived before it or other factors e.g. network conditions/liveness
}

func (ds DemoActorService) handleAddActor(actor Actor) {
	modified := ds.actorPool.SetIfAbsent(actor.ID(), actor)
	if modified {
		go actor.Start()
	}

	// If IDs are being managed properly there should never be a reason why there
	// are multiple attempts to add the same Actor to the pool.
}

// Need to listen for
//   - Incoming Tasks
//   - Actors being added and removed
//   - Queries about the cache size
func (ds DemoActorService) Run() {
	var stop bool
	defer close(ds.shutdown)

	for {
		// Process incoming data on the actor and task queue.
		// If we get a signal to shutdown then the channels are closed and
		// the queues will be cleared before shutdown is complete.
		select {
		case <-ds.shutdown:
			stop = true
			close(ds.addActorQueue)
			close(ds.removeActorQueue)
			close(ds.taskQueue)
		case actor := <-ds.addActorQueue:
			// listens for calls to RegisterActor.
			ds.handleAddActor(actor)
		case actor := <-ds.removeActorQueue:
			// remove actor from pool and stop
			ds.actorPool.Remove(actor.ID())
			actor.Stop()
		case task := <-ds.taskQueue:
			// As Tasks come in find the Actor in the pool that the Task is meant for and add it to their task queue.
			ds.handleActorTask(task)
		}

		// If shutdown signal hasn't come in keep reading from channels above
		if !stop {
			continue
		}

		// If ActorService gets signal to stop we empty the actor and task queues.
		// Then we stop each actor and wait until they are all done emptying their
		// task queues before we return.

		// Process all actors in the actor queue
		for actor := range ds.addActorQueue {
			ds.handleAddActor(actor)
		}

		// Process all tasks in the task queue
		for task := range ds.taskQueue {
			ds.handleActorTask(task)
		}

		// If we get this far that means the ActorService has been shutdown and the actor queue and the task queue are empty.
		// We can now stop each running Actor which will trigger them clearing out their own task queues.
		actors := ds.actorPool.Items()
		for _, actor := range actors {
			actor.Stop()
		}

		// We need to ideally block here until all Actors are done?
		// we could have another channel for remove actor queue
		// when all the actors have finished cleaning up after the Stop calls
		// they can write to the channel, once the cache is empty we know they all finished
		// so we can break from for range chan loop and be done
		for actor := range ds.removeActorQueue {
			ds.actorPool.Remove(actor.ID())

			// actor pool is empty so we can break and return
			if ds.actorPool.Count() == 0 {
				break
			}
		}

		break
	}

	return
}

func (ds DemoActorService) Shutdown() {
	ds.shutdown <- true
}

func (ds DemoActorService) RegisterActor(actor Actor) {
	ds.addActorQueue <- actor
}

func (ds DemoActorService) RemoveActor(actor Actor) {
	ds.removeActorQueue <- actor
}

type DemoActor struct {
	id        string
	parent    ActorService
	taskQueue chan Task
	shutdown  chan bool
	state     int
}

func NewDemoActor() DemoActor {
	return DemoActor{
		id:        "test",
		taskQueue: make(chan Task, taskQueueSize),
		shutdown:  make(chan bool),
		state:     0,
	}
}

func (da DemoActor) handleTask(task Task) {
	switch t := task.(type) {
	case AddTask:
		// handle
		t.Execute()
	case MinusTask:
		// handle
		t.Execute()
	}
}

func (da DemoActor) Start() {
	defer close(da.shutdown)
	var stop bool

	for {
		select {
		case <-da.shutdown:
			stop = true
			close(da.taskQueue)
		case task := <-da.taskQueue:
			da.handleTask(task)
		}

		if !stop {
			continue
		}

		// Empty task queue
		for task := range da.taskQueue {
			// determine what type of task and handle it appropriately.
			da.handleTask(task)
		}

		da.parent.RemoveActor(da)

		return
	}

	// only way we get here is if the task queue is closed which should only happen if stop is called
	// gracefully clean up actor state
	// need to signal to the ActorService that this Actor is dead somehow so it can be removed from actor pool
}

func (da DemoActor) ID() string {
	return "test-actor"
}

func (da DemoActor) AddTask(t Task) {
	da.taskQueue <- t
}

func (da DemoActor) Stop() {
	da.shutdown <- true
}

type AddTask struct {
	actorID string
	x       int

	// if we need a response to the callsite of where this task was sent, we can utilize a channel
	resp chan int
}

func (t AddTask) ActorID() string {
	return t.actorID
}

func (t AddTask) Execute() error {
	// take in currecnt actor state needed
	// change state and return the new state
	//
	//
	newState := 0
	t.resp <- newState
	return nil
}

type MinusTask struct {
	actorID string
	x       int
}

func (t MinusTask) ActorID() string {
	return t.actorID
}

func (t MinusTask) Execute() error {
	//TODO implement me
	panic("implement me")
}
