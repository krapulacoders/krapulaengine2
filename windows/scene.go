package windows

const STATE_UNINITED = 0
const STATE_INITED = 0
const STATE_RUNNING = 0
const STATE_PAUSED = 0
const STATE_TERMINATED = 0

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
	HandleInput(key_events []KeyboardInputEvent, mouse_events []MouseInputEvent)

	// Update the game frame
	Tick(timedelta float64)

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
	isRunning() bool

	// Pause the scene
	Pause()
	// Returns true if the scene is in state STATE_PAUSED
	isPaused() bool

	Exit()
	// returns true if this scene is in STATE_TERMINATED.
	IsTerminated() bool
}

// Partial Scene Implementation that implements shared basic fields and functionality

type SceneImpl struct {
	focused bool
	state   int
}

func (self *SceneImpl) SetFocused(focused bool) {
	self.focused = focused
}

func (self *SceneImpl) IsFocused() bool {
	return self.focused
}

func (self *SceneImpl) IsInited() bool {
	return self.state != STATE_TERMINATED && self.state != STATE_UNINITED
}

func (self *SceneImpl) IsRunning() bool {
	return self.state == STATE_RUNNING
}

func (self *SceneImpl) IsPaused() bool {
	return self.state == STATE_PAUSED
}

func (self *SceneImpl) IsTerminated() bool {
	return self.state == STATE_TERMINATED
}
