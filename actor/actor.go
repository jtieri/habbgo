package actor

/*

When the server starts up it will initialize each ActorService and maintain a reference to each one. As Sessions are
handled by the server they get the reference to each ActorService passed into their constructor call. When a Session
initializes a Player they first generate a unique user ID so that they know how to reference the Player when needed.
After Players are initialized the Session will call RegisterActor on the PlayerActorService using the unique ID as the key.
Now the Player is a globally accessible object. As Sessions handle incoming packets and call the appropriate handler,
they will pass the reference to each ActorService into the command function so that appropriate state can be fetched
on a per packet basis, which ensures that we are always working with the current state.

ActorServices need to be able to handle context cancellation so that term signals can be observed from the server.
ActorServices need to be able to be informed when Actors stop and need removed from the Actor cache.

When new Actors are initialized (players, rooms, catalog, etc.) they need to be registered in the appropriate ActorPool
which is maintained in each ActorService. Actors maintain a single buffered channel of Tasks which are added by calling AddTask.
Actors process Tasks that are pending one at a time. This ensures that each Actor is maintaining its own state and that
Tasks are processed in the order that they are received. Actors need an ID so that they can be retrieved from the ActorPool
when a Task designated for it is added to the ActorService. When Actors are stopped they finish processing their Task queue
and then they shutdown.


ActorServices have a single buffered channel serving as a queue of incoming Tasks. Tasks are meant to be delivered to a single
Actor which is inside an ActorPool.
The ActorServices are run when the server starts up. The ActorService listens for incoming Tasks and passes
them to the appropriate Actor which they are meant for. On Shutdown it closes the tasks channel blocking any new incoming tasks,
waits for all received tasks to be assigned to Actors. Then it invokes Stop on each Actor and waits on them to finish.


If an Actor expects a response a chan must be included in the Task for the receiving Actor to send the response
*/

type ActorService interface {
	RegisterActor(actor Actor)
	RemoveActor(actor Actor)
	Run()
	Shutdown()
	SubmitTask(t Task)
}

type Actor interface {
	ID() string
	AddTask(t Task)
	Start()
	Stop()
}

type Task interface {
	ActorID() string
	Execute() error
}
