// Based on https://github.com/go-gl/examples/blob/master/glfw31-gl41core-cube/cube.go

package main

import (
	"image"
	_ "image/png"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/krapulacoders/krapulaengine2/graphics"
	"github.com/krapulacoders/krapulaengine2/windows"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	fontRegular, err := graphics.ReadFont("RobotoMono-Regular.ttf")
	img, err := graphics.GenerateImageFromFont("Hello World", fontRegular, 16)
	if err != nil {
		panic(err.Error())
	}

	scene := newCubeScene(img)
	windows.Init()
	windows.AddScene("cube", scene)
	windows.SetCurrentScene("cube")
	windows.MainLoop()

}

var vertexShader = `
#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
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

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
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

type cubeScene struct {
	windows.SimpleSceneImpl

	angle                     float32
	model                     mgl32.Mat4
	program                   uint32
	modelUniform              int32
	vao, vbo                  uint32
	textureText, textureImage uint32
	rotateDirection           float32
	image                     *image.RGBA
	showText                  bool
}

func newCubeScene(img *image.RGBA) *cubeScene {
	scene := new(cubeScene)
	scene.SetState(windows.StateUninited)
	scene.rotateDirection = 1
	scene.image = img
	return scene
}

func (s *cubeScene) Tick(timedelta float64, keyStates []bool) {
	if keyStates[glfw.KeyLeft] {
		s.angle -= float32(timedelta)
	}
	if keyStates[glfw.KeyRight] {
		s.angle += float32(timedelta)
	}
}

func (s *cubeScene) Render() {

	s.model = mgl32.HomogRotate3D(float32(s.angle), mgl32.Vec3{0, 1, 0})

	// Render
	gl.UseProgram(s.program)
	gl.UniformMatrix4fv(s.modelUniform, 1, false, &s.model[0])

	gl.BindVertexArray(s.vao)

	gl.ActiveTexture(gl.TEXTURE0)
	if s.showText {
		gl.BindTexture(gl.TEXTURE_2D, s.textureText)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, s.textureImage)
	}

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
}

func (s *cubeScene) HandleInput(keyEvents []windows.KeyboardInputEvent, mouseEvents []windows.MouseInputEvent) windows.WindowAction {
	// nothing for now
	for _, event := range keyEvents {
		switch event.Key {
		case glfw.KeyEscape:
			return windows.WindowActionExit
		case glfw.KeySpace:
			if event.Action == glfw.Press {
				s.showText = !s.showText
			}
		}
	}
	return windows.WindowActionNone
}

func (s *cubeScene) Init() {
	// Configure the vertex and fragment shaders
	program, err := graphics.NewProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	s.program = program

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	s.model = mgl32.Ident4()
	s.modelUniform = gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(s.modelUniform, 1, false, &s.model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the textures
	s.textureText, err = graphics.RegisterTextureFromImage("hello_world", s.image)
	if err != nil {
		panic(err)
	}
	s.textureImage, err = graphics.RegisterTextureFromFile("square", "square.png")
	if err != nil {
		panic(err)
	}

	// Configure the vertex data
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)

	gl.GenBuffers(1, &s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

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

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	s.angle = 0.0
	s.SetState(windows.StateInited)
}
