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
	. "github.com/krapulacoders/krapulaengine2/windows"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	InitWindowing()
	defer glfw.Terminate()

	font_regular, err := graphics.ReadFont("RobotoMono-Regular.ttf")
	img, err := graphics.GenerateImageFromFont("Hello World", font_regular, 16)
	if err != nil {
		panic(err.Error())
	}

	scene := NewCubeScene(img)
	MainWindow.AddScene("cube", scene)
	MainWindow.SetCurrentScene("cube")
	MainWindow.MainLoop()

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

type CubeScene struct {
	SimpleSceneImpl

	angle                       float32
	model                       mgl32.Mat4
	program                     uint32
	modelUniform                int32
	vao, vbo                    uint32
	texture_text, texture_image uint32
	rotate_direction            float32
	image                       *image.RGBA
	show_text                   bool
}

func NewCubeScene(img *image.RGBA) *CubeScene {
	scene := new(CubeScene)
	scene.SetState(STATE_UNINITED)
	scene.rotate_direction = 1
	scene.image = img
	return scene
}

func (self *CubeScene) Tick(timedelta float64, key_states []bool) {
	if key_states[glfw.KeyLeft] {
		self.angle -= float32(timedelta)
	}
	if key_states[glfw.KeyRight] {
		self.angle += float32(timedelta)
	}
}

func (self *CubeScene) Render() {

	self.model = mgl32.HomogRotate3D(float32(self.angle), mgl32.Vec3{0, 1, 0})

	// Render
	gl.UseProgram(self.program)
	gl.UniformMatrix4fv(self.modelUniform, 1, false, &self.model[0])

	gl.BindVertexArray(self.vao)

	gl.ActiveTexture(gl.TEXTURE0)
	if self.show_text {
		gl.BindTexture(gl.TEXTURE_2D, self.texture_text)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, self.texture_image)
	}

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
}

func (self *CubeScene) HandleInput(key_events []KeyboardInputEvent, mouse_events []MouseInputEvent) WindowAction {
	// nothing for now
	for _, event := range key_events {
		switch event.Key {
		case glfw.KeyEscape:
			return WINDOW_ACTION_EXIT
		case glfw.KeySpace:
			if event.Action == glfw.Press {
				self.show_text = !self.show_text
			}
		}
	}
	return WINDOW_ACTION_NONE
}

func (self *CubeScene) Init() {
	// Configure the vertex and fragment shaders
	program, err := graphics.NewProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	self.program = program

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
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

	// Load the textures
	self.texture_text, err = graphics.NewTextureFromImage(self.image)
	if err != nil {
		panic(err)
	}
	self.texture_image, err = graphics.NewTextureFromFile("square.png")
	if err != nil {
		panic(err)
	}

	// Configure the vertex data
	gl.GenVertexArrays(1, &self.vao)
	gl.BindVertexArray(self.vao)

	gl.GenBuffers(1, &self.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo)
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

	self.angle = 0.0
	self.SetState(STATE_INITED)
}
