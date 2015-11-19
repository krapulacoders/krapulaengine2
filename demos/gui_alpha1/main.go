package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/krapulacoders/krapulaengine2/graphics/gui"
	"github.com/krapulacoders/krapulaengine2/windows"
)

func main() {
	windows.InitWindowing()
	windows.MainWindow.Init()
	gui.InitTopContainer()
	dummy_scene := new(backgroundScene)
	windows.MainWindow.AddScene("dummy", dummy_scene)
	windows.MainWindow.SetCurrentScene("dummy")
	windows.MainWindow.AddOverlay("gui", gui.TopContainer)
	windows.MainWindow.MainLoop()
}

type backgroundScene struct {
	windows.SimpleSceneImpl
}

func (self *backgroundScene) Init() {
	self.SetState(windows.STATE_INITED)
}

func (self *backgroundScene) HandleInput(key_events []windows.KeyboardInputEvent, mouse_events []windows.MouseInputEvent) windows.WindowAction {
	for _, event := range key_events {
		switch event.Key {
		case glfw.KeyEscape:
			return windows.WINDOW_ACTION_EXIT
		}
	}
	return windows.WINDOW_ACTION_NONE
}

func (self *backgroundScene) Tick(timedelta float64, key_states []bool) {}

func (self *backgroundScene) Render() {}
