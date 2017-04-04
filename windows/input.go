package windows

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// KeyboardInputEvent represents a glfw keyboard event
type KeyboardInputEvent struct {
	Key      glfw.Key
	Scancode int
	Action   glfw.Action
	Mod      glfw.ModifierKey
}

// MouseInputEvent represents a glfw mouse event
type MouseInputEvent struct {
	Button glfw.MouseButton
	Action glfw.Action
	Mod    glfw.ModifierKey
	X, Y   float32
}
