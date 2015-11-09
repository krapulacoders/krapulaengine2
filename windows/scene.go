package windows

type Scene interface {
	// Window creates an input channel that is passed to the scene if necessary
	// All relevant input will be sent to it.

	// Returns true if the Scene cares about input
	AcceptsInput() bool
	SetInputChannel(ic chan Input)

	// Processes input events
	HandleInput()

	Tick(timedelta float64)

	Render()

	// Some scenes want to ignore input or pause if they aren't focused
	SetFocused(is_focused bool)
	IsFocused() bool
}

/* Partial Scene Implementation that implements shared basic fields and functionality
 *
 */
type SceneImpl struct {
	is_focused    bool
	input_channel chan Input
}

func (self *SceneImpl) SetInputChannel(ic chan Input) {
	self.input_channel = ic
}

func (self *SceneImpl) SetFocused(is_focused bool) {
	self.is_focused = is_focused
}

func (self *SceneImpl) IsFocused() bool {
	return self.is_focused
}
