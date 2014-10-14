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
	SetFocus(is_focused bool)
	IsFocused() bool

}