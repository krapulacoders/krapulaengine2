package gui

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	//"github.com/krapulacoders/krapulaengine2/graphics"
	"github.com/krapulacoders/krapulaengine2/windows"
)

type TopContainer struct {
	windows.BasicSceneImpl
	children      []Component
	width, height int
	program       uint32
	model         int
	modelUniform  mgl32.Mat4
	vao, vbo      int
}

func NewTopContainer(window *windows.Window) *TopContainer {
	c := new(TopContainer)
	c.SetState(windows.STATE_UNINITED)

	c.width, c.height = window.GetSize()
	c.children = make([]Component, 0, 10)
	return c
}

func (self *TopContainer) Tick(timedelta float64, key_states []bool) {
}

func (self *TopContainer) Render() {
}

func (self *TopContainer) HandleInput(key_events []windows.KeyboardInputEvent,
	mouse_events []windows.MouseInputEvent) windows.WindowAction {
	// nothing for now
	for _, event := range key_events {
		_ = event
	}
	return windows.WINDOW_ACTION_NONE
}

const (
	vertexShader   = ""
	fragmentShader = ""
)

func (self *TopContainer) Init() {
	// Configure the vertex and fragment shaders
	program, err := windows.NewProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	self.program = program

	gl.UseProgram(program)
	/*
		projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(self.width)/float32(self.height), 0.1, 10.0)
		projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

		camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
		cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

		self.model = mgl32.Ident4()
		self.modelUniform = gl.GetUniformLocation(program, gl.Str("model\x00"))
		gl.UniformMatrix4fv(self.modelUniform, 1, false, &self.model[0])

		textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
		gl.Uniform1i(textureUniform, 0)

		gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

		// Configure the vertex data
		gl.GenVertexArrays(1, &self.vao)
		gl.BindVertexArray(self.vao)

		gl.GenBuffers(1, &self.vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(QuadVertices), gl.STATIC_DRAW)

		vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
		gl.EnableVertexAttribArray(vertAttrib)
		gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

		texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
		gl.EnableVertexAttribArray(texCoordAttrib)
		gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

		// Configure global settings
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LESS)
		gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	*/
	self.SetState(windows.STATE_INITED)
}

var quadVertices = []float32{}
