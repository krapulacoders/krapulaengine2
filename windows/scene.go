package windows

const (
	STATE_UNINITED   = iota
	STATE_INITED     = iota
	STATE_RUNNING    = iota
	STATE_PAUSED     = iota
	STATE_TERMINATED = iota
)

// Scenes are responsible for handling input and rendering their content.
// They can be in the following states:
// 1. Uninited - only bare minimum setup at this point
// 2. Inited - resources should be allocated
// 3. Running - Running normally
// 4. Paused - Paused, can be resumed to return to the Running state
// 5. Terminated - all resources should be released
type Scene interface {
	// Window creates an input channel that is passed to the scene if necessary
	// All relevant input will be sent to it.

	// Returns true if the Scene cares about input
	AcceptsInput() bool

	// Processes input events
	HandleInput(key_events []KeyboardInputEvent, mouse_events []MouseInputEvent) WindowAction

	// Update the game frame and process time-dependant input
	Tick(timedelta float64, key_states []bool)

	// Render the scene
	Render()

	// Some scenes might want to ignore input or pause if they aren't focused
	SetFocused(is_focused bool)
	IsFocused() bool

	// Init the scene
	Init()
	// Returns true unless the scene is in state STATE_UNINITED or STATE_TERMINATED
	IsInited() bool

	// Set the scene in running mode
	Run()
	// Returns true if the scene is in state STATE_RUNNING
	IsRunning() bool

	// Pause the scene
	Pause()
	// Returns true if the scene is in state STATE_PAUSED
	IsPaused() bool

	Exit()
	// returns true if this scene is in STATE_TERMINATED.
	IsTerminated() bool
}

// Partial Scene Implementation that implements shared basic fields and functionality
// You still need to implement the following functions yourself:
// Init()
// HandleInput(key_events []KeyboardInputEvent, mouse_events []MouseInputEvent) WindowAction
// Tick(timedelta float64, key_states []bool)
// Render()

type BasicSceneImpl struct {
	focused bool
	state   int
}

func (self *BasicSceneImpl) SetFocused(focused bool) {
	self.focused = focused
}

func (self *BasicSceneImpl) IsFocused() bool {
	return self.focused
}

func (self *BasicSceneImpl) IsInited() bool {
	return self.state != STATE_TERMINATED && self.state != STATE_UNINITED
}

func (self *BasicSceneImpl) IsRunning() bool {
	return self.state == STATE_RUNNING
}

func (self *BasicSceneImpl) IsPaused() bool {
	return self.state == STATE_PAUSED
}

func (self *BasicSceneImpl) IsTerminated() bool {
	return self.state == STATE_TERMINATED
}

func (self *BasicSceneImpl) SetState(state int) {
	self.state = state
}

// A simple scene implementation to use when you want default state transitions and input
type SimpleSceneImpl struct {
	BasicSceneImpl
}

func (self *SimpleSceneImpl) AcceptsInput() bool {
	return true
}

func (self *SimpleSceneImpl) Run() {
	if !self.IsInited() {
		panic("Cannot run scene that hasn't been inited. state: " + string(self.state))
	} else {
		self.SetState(STATE_RUNNING)
	}
}

func (self *SimpleSceneImpl) Pause() {
	if !self.IsRunning() {
		panic("Cannot pause scene that isn't running")
	} else {
		self.SetState(STATE_PAUSED)
	}
}

func (self *SimpleSceneImpl) Exit() {
	self.SetState(STATE_TERMINATED)
}
