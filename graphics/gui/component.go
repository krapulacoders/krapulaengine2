package gui

import (
	. "github.com/krapulacoders/krapulaengine2/windows"
	"image"
)

type Component interface {

	// Returns true if the Scene cares about input
	AcceptsInput() bool

	// Processes input events
	HandleInput(key_events []KeyboardInputEvent, mouse_events []MouseInputEvent) WindowAction

	// Update the game frame and process time-dependant input
	Tick(timedelta float64, key_states []bool)

	// Render the component
	Render()
}

type Rectangle struct {
	X, Y, W, H int
}

// converts the rectangle to an image Rectangle.
func (self *Rectangle) Bounds() image.Rectangle {
	return image.Rect(self.X, self.Y, self.W, self.H)
}

// A partial component implementation that ignores input and ticks
type StaticComponentImpl struct {
	Rectangle
}

func (self *StaticComponentImpl) AcceptsInput() bool {
	return false
}

func (self *StaticComponentImpl) HandleInput(key_events []KeyboardInputEvent, mouse_events []MouseInputEvent) WindowAction {
	return WINDOW_ACTION_NONE
}

func (self *StaticComponentImpl) Tick(timedelta float64, key_states []bool) {

}
