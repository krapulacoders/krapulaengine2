package windows

// states
const (
	StateUninited   = iota
	StateInited     = iota
	StateRunning    = iota
	StatePaused     = iota
	StateTerminated = iota
)

// Scene are responsible for handling input and rendering their content.
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
	HandleInput(keyEvents []KeyboardInputEvent, mouseEvents []MouseInputEvent) WindowAction

	// Update the game frame and process time-dependant input
	Tick(timedelta float64, keyStates []bool)

	// Some scenes might want to ignore input or pause if they aren't focused
	SetFocused(isFocused bool)
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

// BasicSceneImpl is a partial Scene implementation that implements shared basic fields and functionality
// You still need to implement the following functions yourself:
// Init()
// HandleInput(key_events []KeyboardInputEvent, mouse_events []MouseInputEvent) WindowAction
// Tick(timedelta float64, key_states []bool)
type BasicSceneImpl struct {
	focused bool
	state   int
}

// SetFocused ...
func (s *BasicSceneImpl) SetFocused(focused bool) {
	s.focused = focused
}

// IsFocused ...
func (s *BasicSceneImpl) IsFocused() bool {
	return s.focused
}

// IsInited ...
func (s *BasicSceneImpl) IsInited() bool {
	return s.state != StateTerminated && s.state != StateUninited
}

// IsRunning ...
func (s *BasicSceneImpl) IsRunning() bool {
	return s.state == StateRunning
}

// IsPaused ...
func (s *BasicSceneImpl) IsPaused() bool {
	return s.state == StatePaused
}

// IsTerminated ...
func (s *BasicSceneImpl) IsTerminated() bool {
	return s.state == StateTerminated
}

// SetState sets state
func (s *BasicSceneImpl) SetState(state int) {
	s.state = state
}

// SimpleSceneImpl is a simple scene implementation to use when you want default state transitions and input
type SimpleSceneImpl struct {
	BasicSceneImpl
}

// AcceptsInput true
func (s *SimpleSceneImpl) AcceptsInput() bool {
	return true
}

// Run checks that the scene is inited and sets the state to running
func (s *SimpleSceneImpl) Run() {
	if !s.IsInited() {
		panic("Cannot run scene that hasn't been inited. state: " + string(s.state))
	} else {
		s.SetState(StateRunning)
	}
}

// Pause pauses the scene
func (s *SimpleSceneImpl) Pause() {
	if !s.IsRunning() {
		panic("Cannot pause scene that isn't running")
	} else {
		s.SetState(StatePaused)
	}
}

// Exit terminates the scene
func (s *SimpleSceneImpl) Exit() {
	s.SetState(StateTerminated)
}
