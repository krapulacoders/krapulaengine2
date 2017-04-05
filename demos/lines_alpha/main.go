package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
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

func getVec3FromAngle(angle, size, z float32) mgl32.Vec3 {
	a64 := (float64)(angle)
	return mgl32.Vec3{size * (float32)(math.Cos(a64)), size * (float32)(math.Sin(a64)), z}
}

type backgroundScene struct {
	windows.SimpleSceneImpl
	renderGroups    [3]*graphics.RenderGroup
	lineManager     *rendergroups.BasicRenderGroup2D
	line1, line2    rendergroups.GenericObject2D
	triangleManager *rendergroups.BasicRenderGroup2D
	triangle1       rendergroups.GenericObject2D
	pointManager    *rendergroups.BasicRenderGroup2D
}

func (s *backgroundScene) Init() {

	s.line1 = rendergroups.GenericObject2D{
		Coords: []mgl32.Vec3{getVec3FromAngle(0, 300, -0.5), getVec3FromAngle(math.Pi, 300, -0.5)},
		Color:  mgl32.Vec4{1, 0, 0, 1},
	}

	s.line2 = rendergroups.GenericObject2D{
		Coords: []mgl32.Vec3{getVec3FromAngle(math.Pi/2, 700, 0.2), getVec3FromAngle(3*math.Pi/2, 800, -0.5)},
		Color:  mgl32.Vec4{0, 1, 0, 1},
	}

	s.renderGroups[0], s.lineManager = rendergroups.NewBasicRenderGroup2D("test", gl.LINES, 2, 0, 0)
	s.lineManager.AddObject(&s.line1)
	s.lineManager.AddObject(&s.line2)

	s.triangle1 = rendergroups.GenericObject2D{
		Coords: []mgl32.Vec3{getVec3FromAngle(0, 500, 0), getVec3FromAngle(2*math.Pi/3, 500, 0),
			getVec3FromAngle(4*math.Pi/3, 600, 0)},
		Color: mgl32.Vec4{0, 0, 1, 1},
	}
	s.renderGroups[1], s.triangleManager = rendergroups.NewBasicRenderGroup2D("test", gl.TRIANGLES, 1, 0, 0)
	s.triangleManager.AddObject(&s.triangle1)

	s.renderGroups[2], s.pointManager = rendergroups.NewBasicRenderGroup2D("test", gl.POINTS, 1, 0, 0)
	for xi := float32(-800); xi <= 800; xi += 100 {
		for yi := float32(-800); yi <= 800; yi += 100 {
			s.pointManager.AddObject(&rendergroups.GenericObject2D{
				Coords: []mgl32.Vec3{mgl32.Vec3{xi, yi, rand.Float32()*2 - 1}},
				Color:  mgl32.Vec4{1, 1, 1, 1},
			})
		}
	}

	for _, g := range s.renderGroups {
		graphics.AddRenderGroup(g)
		g.SetDepthTestMode(true, gl.LESS)
	}

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
