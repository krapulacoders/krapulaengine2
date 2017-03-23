package gui

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/krapulacoders/krapulaengine2/graphics"
	"github.com/krapulacoders/krapulaengine2/windows"
)

var TopContainer *TopContainerType

func InitTopContainer() {
	TopContainer = newTopContainer()
}

// An invisible container that covers the entire window
// It manages the shader programs.

type TopContainerType struct {
	windows.SimpleSceneImpl
	children      []Component
	width, height int
	program       uint32
	model         int
	modelUniform  mgl32.Mat4
	vao, vbo      int
}

func newTopContainer() *TopContainerType {
	c := new(TopContainerType)
	c.SetState(windows.StateUninited)

	c.width, c.height = windows.GetSize()
	c.children = make([]Component, 0, 10)
	return c
}

func (self *TopContainerType) Tick(timedelta float64, key_states []bool) {
	for _, child := range self.children {
		child.Tick(timedelta, key_states)
	}
}

func (self *TopContainerType) Render() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	for _, child := range self.children {
		child.Render()
	}
}

func (self *TopContainerType) HandleInput(key_events []windows.KeyboardInputEvent,
	mouse_events []windows.MouseInputEvent) windows.WindowAction {
	// nothing for now
	for _, event := range key_events {
		_ = event
	}
	return windows.WindowActionNone
}

func (self *TopContainerType) Program() uint32 {
	return self.program
}

var vertexShader = `
#version 330
in vec4 vertAndTexCoord;
out vec2 fragTexCoord;
void main() {
    fragTexCoord = vec2(vertAndTexCoord[2], vertAndTexCoord[3]);
    gl_Position = vec4(vertAndTexCoord[0], vertAndTexCoord[1], 0, 1);
}
` + "\x00"

var fragmentShader = `
#version 330
uniform sampler2D tex;
in vec2 fragTexCoord;
out vec4 outputColor;
void main() {
    outputColor = texture(tex, fragTexCoord);
}
` + "\x00"

func (self *TopContainerType) Init() {
	// Configure the vertex and fragment shaders
	program, err := graphics.NewProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	self.program = program

	gl.UseProgram(program)
	// use texture 0
	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)
	self.SetState(windows.StateInited)
}
