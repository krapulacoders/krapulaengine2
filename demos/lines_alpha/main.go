package main

import (
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/krapulacoders/krapulaengine2/graphics"
	"github.com/krapulacoders/krapulaengine2/graphics/rendergroups"
	"github.com/krapulacoders/krapulaengine2/windows"
)

func main() {
	windows.Init()

	lineScene := new(backgroundScene)
	windows.AddScene("lines", lineScene)
	windows.SetCurrentScene("lines")
	windows.MainLoop()

}

func getVec2FromAngle(angle float32) mgl32.Vec2 {
	a64 := (float64)(angle)
	return mgl32.Vec2{(float32)(math.Cos(a64)), (float32)(math.Sin(a64))}
}

type backgroundScene struct {
	windows.SimpleSceneImpl
	renderGroup   *graphics.RenderGroup
	objectManager *rendergroups.BasicRenderGroup2D
	line1, line2  rendergroups.GenericObject2D
}

func (s *backgroundScene) Init() {

	s.line1 = rendergroups.GenericObject2D{
		Coords: []mgl32.Vec2{getVec2FromAngle(0), getVec2FromAngle(math.Pi)},
		Color:  mgl32.Vec4{1, 0, 1, 1},
	}

	s.line2 = rendergroups.GenericObject2D{
		Coords: []mgl32.Vec2{getVec2FromAngle(math.Pi / 2), getVec2FromAngle(3 * math.Pi / 2)},
		Color:  mgl32.Vec4{0, 1, 1, 1},
	}

	s.renderGroup, s.objectManager = rendergroups.NewBasicRenderGroup2D("test", gl.LINES, 2, 0, 0)
	s.objectManager.AddObject(&s.line1)
	s.objectManager.AddObject(&s.line2)
	graphics.AddManager(s.renderGroup)
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
	// rotate in different directions 1/4 rotations per second
	s.line1.Rotation += (float32)(timedelta * math.Pi / 4)
	s.line2.Rotation -= (float32)(timedelta * math.Pi / 4)
}
