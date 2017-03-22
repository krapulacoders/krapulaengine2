package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/krapulacoders/krapulaengine2/windows"
)

func main() {
	windows.Init()

	lineScene := new(backgroundScene)
	windows.AddScene("lines", lineScene)
	windows.SetCurrentScene("lines")
	windows.MainLoop()

	//line1 := rendergroups.GenericObject2D{[]mgl32.Vec2{-1, -1, 1, 1}, }
}

type backgroundScene struct {
	windows.SimpleSceneImpl
}

func (s *backgroundScene) Init() {
	s.SetState(windows.STATE_INITED)
}

func (s *backgroundScene) HandleInput(keyEvents []windows.KeyboardInputEvent, mouseEvents []windows.MouseInputEvent) windows.WindowAction {
	for _, event := range keyEvents {
		switch event.Key {
		case glfw.KeyEscape:
			return windows.WINDOW_ACTION_EXIT
		}
	}
	return windows.WINDOW_ACTION_NONE
}

func (s *backgroundScene) Tick(timedelta float64, keyStates []bool) {

}
