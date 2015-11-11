package windows

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type KeyboardInputEvent struct {
	Key      glfw.Key
	Scancode int
	Action   glfw.Action
	Mod      glfw.ModifierKey
}

type MouseInputEvent struct {
	Button glfw.MouseButton
	Action glfw.Action
	Mod    glfw.ModifierKey
	X, Y   float32
}
