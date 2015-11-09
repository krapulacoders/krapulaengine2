package windows

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type Input interface{}

type KeyboardInputEvent struct {
	key      glfw.Key
	scancode int
	action   glfw.Action
	mod      glfw.ModifierKey
}

type MouseInputEvent struct {
	button glfw.MouseButton
	action glfw.Action
	mod    glfw.ModifierKey
}
