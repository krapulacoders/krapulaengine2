package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/krapulacoders/krapulaengine2/graphics/gui"
	"github.com/krapulacoders/krapulaengine2/windows"
)

func main() {
	windows.Init()
	gui.InitTopContainer()
	dummyScene := new(backgroundScene)
	windows.AddScene("dummy", dummyScene)
	windows.SetCurrentScene("dummy")
	windows.AddOverlay("gui", gui.TopContainer)
	windows.MainLoop()
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

func (s *backgroundScene) Tick(timedelta float64, keyStates []bool) {}

func (s *backgroundScene) Render() {}
