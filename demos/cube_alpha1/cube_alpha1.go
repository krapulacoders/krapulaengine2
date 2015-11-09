package main

// mostly copied from https://github.com/go-gl/examples/blob/master/glfw31-gl41core-cube/cube.go

import (
	"fmt"
	gl "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	. "github.com/krapulacoders/krapulaengine2/windows"
	"runtime"
)

var window_width = 1200
var window_height = 700

func main() {

	runtime.LockOSThread()

	InitWindowing()

	defer glfw.Terminate()

	window := NewWindow(window_width, window_height)
	window.Init()

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	scene := newHelloWorldScene()

	window.AddScene("cube", scene)
	window.SetCurrentScene("cube")

	window.MainLoop()

}

type HelloWorldScene struct {
	SceneImpl
	program      uint32
	previousTime float64
	angle        float32
	modelUniform int32
	vao, vbo     uint32
}

func newHelloWorldScene() *HelloWorldScene {
	sc := new(HelloWorldScene)
	// Configure the vertex and fragment shaders
	setupShaders(sc)

	return sc
}

func setupShaders(sc *HelloWorldScene) {
	program, err := NewProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(window_width)/float32(window_height), 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{5, 5, 5}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	sc.modelUniform = gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(sc.modelUniform, 1, false, &model[0])

	uniformColor := int32(gl.GetUniformLocation(program, gl.Str("triangleColor\x00")))
	gl.Uniform4f(uniformColor, 1.0, 0.0, 1.0, 1.0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Configure the vertex data
	gl.GenVertexArrays(1, &sc.vao)
	gl.BindVertexArray(sc.vao)

	gl.GenBuffers(1, &sc.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, sc.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
}

func (self *HelloWorldScene) Render() {
	// Update
	time := glfw.GetTime()
	elapsed := time - self.previousTime
	self.previousTime = time

	self.angle += float32(elapsed) * 10
	//fmt.Println(self.angle)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	//color := []float32{1, 0, 0, 1}
	//gl.ClearBufferfv(gl.COLOR, 0, &color[0])

	model := mgl32.HomogRotate3D(self.angle, mgl32.Vec3{0, 1, 0})

	// Render
	gl.UseProgram(self.program)
	gl.UniformMatrix4fv(self.modelUniform, 1, false, &model[0])
	gl.BindVertexArray(self.vao)
	gl.PointSize(30)
	gl.LineWidth(30)
	gl.DrawArrays(gl.LINES, 0, 6*2*3)
}

// Processes input events
func (self *HelloWorldScene) HandleInput() {

}

func (self *HelloWorldScene) Tick(timedelta float64) {

}

func (self *HelloWorldScene) AcceptsInput() bool {
	return true
}

var vertexShader = `
#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform vec4 triangleColor;

in vec3 vert;
out vec4 color;
void main() {
	color = triangleColor;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

var fragmentShader = `
#version 330
in vec4 color;
out vec4 outputColor;
void main() {
    outputColor = color;
    //outputColor = vec4(1.0, 1.0, 1.0, 1.0);
}
` + "\x00"

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-0.5, -0.5, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}
