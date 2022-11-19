package actions

type Action struct {
	Name      string // key sent to client "status code"
	Params    int    // value sent to client (e.g. drink id, food id etc)
	Lifecycle int    // how long does this action run in seconds

	SwapAction    *Action // An optional action to switch to
	TimeTillSwap  int     // how long in seconds till we switch to optional swap action
	SwapLifecycle int     // how long till we switch back to the original action after swapping
}

func NewAction(name string, params, lifecycle, swapLifecycle int, swapAction *Action) Action {
	return Action{
		Name:          name,
		Params:        params,
		Lifecycle:     lifecycle,
		SwapAction:    swapAction,
		SwapLifecycle: swapLifecycle,
	}
}

// Set a status with a limited lifetime, and optional swap to action every x seconds which lasts for
// * x seconds. Use -1 and 'null' for action and lifetimes to make it last indefinitely.
