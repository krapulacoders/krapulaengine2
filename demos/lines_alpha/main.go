package main

import (
	"fmt"
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
	graphics.SetClearColor(0.5, 0.5, 0.5, 1)
	windows.MainLoop()

}

func getVec2FromAngle(angle float32, size float32) mgl32.Vec2 {
	a64 := (float64)(angle)
	return mgl32.Vec2{size * (float32)(math.Cos(a64)), size * (float32)(math.Sin(a64))}
}

type backgroundScene struct {
	windows.SimpleSceneImpl
	renderGroup1, renderGroup2 *graphics.RenderGroup
	lineManager                *rendergroups.BasicRenderGroup2D
	line1, line2               rendergroups.GenericObject2D
	triangleManager            *rendergroups.BasicRenderGroup2D
	triangle1                  rendergroups.GenericObject2D
}

func (s *backgroundScene) Init() {

	s.line1 = rendergroups.GenericObject2D{
		Coords: []mgl32.Vec2{getVec2FromAngle(0, 0.3), getVec2FromAngle(math.Pi, 0.3)},
		Color:  mgl32.Vec4{1, 0, 0, 1},
	}

	s.line2 = rendergroups.GenericObject2D{
		Coords: []mgl32.Vec2{getVec2FromAngle(math.Pi/2, 0.7), getVec2FromAngle(3*math.Pi/2, 0.7)},
		Color:  mgl32.Vec4{0, 1, 0, 1},
	}

	s.renderGroup1, s.lineManager = rendergroups.NewBasicRenderGroup2D("test", gl.LINES, 2, 0, 0)
	s.lineManager.AddObject(&s.line1)
	s.lineManager.AddObject(&s.line2)

	s.triangle1 = rendergroups.GenericObject2D{
		Coords: []mgl32.Vec2{getVec2FromAngle(0, 0.5), getVec2FromAngle(2*math.Pi/3, 0.5), getVec2FromAngle(4*math.Pi/3, 0.6)},
		Color:  mgl32.Vec4{0, 0, 1, 1},
	}
	s.renderGroup2, s.triangleManager = rendergroups.NewBasicRenderGroup2D("test", gl.TRIANGLES, 1, 0, 0)
	s.triangleManager.AddObject(&s.triangle1)
	s.triangleManager.NotifyObjectChanged()

	graphics.AddRenderGroup(s.renderGroup1)
	graphics.AddRenderGroup(s.renderGroup2)
	s.SetState(windows.StateInited)
}

func (s *backgroundScene) HandleInput(keyEvents []windows.KeyboardInputEvent, mouseEvents []windows.MouseInputEvent) windows.WindowAction {
	for _, event := range keyEvents {
		switch event.Key {
		case glfw.KeyEscape:
			return windows.WindowActionExit
		}
	}
	return windows.WindowActionNone
}

var color float32 = 0.3
var increment float32 = 0.001

func (s *backgroundScene) Tick(timedelta float64, keyStates []bool) {
	// rotate in different directions 1/4 rotations per second
	s.line1.Rotation += (float32)(timedelta * math.Pi / 4)
	s.line2.Rotation -= (float32)(timedelta * math.Pi / 4)
	s.lineManager.NotifyObjectChanged()
	s.triangleManager.NotifyObjectChanged()
	fmt.Println("tick!")
	color += increment
	if color > 0.3 {
		increment = -increment
	} else if color < 0 {
		increment = -increment
	}
	graphics.SetClearColor(color, color, color, 1)
}
