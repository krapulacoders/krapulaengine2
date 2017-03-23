package gui

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/krapulacoders/krapulaengine2/windows"
)

type Component interface {

	// Returns true if the Scene cares about input
	AcceptsInput() bool

	// Processes input events
	HandleInput(keyEvents []windows.KeyboardInputEvent, mouseEvents []windows.MouseInputEvent) windows.WindowAction

	// Update the game frame and process time-dependant input
	Tick(timedelta float64, keyStates []bool)

	// Render the component
	Render()
}

type BasicComponentImpl struct {
	X, Y, W, H int
}

type TexturedComponentImpl struct {
	BasicComponentImpl
	texture    uint32
	vbo        uint32
	vao        uint32
	tex_coords []mgl32.Vec2
}

const (
	TOP_LEFT     = 0
	TOP_RIGHT    = 1
	BOTTOM_LEFT  = 2
	BOTTOM_RIGHT = 3
)

func (c *TexturedComponentImpl) Render() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE0, c.texture)
	gl.DrawArrays(gl.TRIANGLES, 0, 6) // 6 vertices starting from position 0

}

// assign a texture and texture coordinates to this component
func (c *TexturedComponentImpl) SetTexture(texture uint32, topleft, topright, bottomleft, bottomright mgl32.Vec2) {
	c.texture = texture
	c.tex_coords = make([]mgl32.Vec2, 4)
	c.tex_coords[TOP_LEFT] = topleft
	c.tex_coords[TOP_RIGHT] = topright
	c.tex_coords[BOTTOM_LEFT] = bottomleft
	c.tex_coords[BOTTOM_RIGHT] = bottomright
	c.updateTexture()
}

func (c *TexturedComponentImpl) updateTexture() {
	topleft := c.tex_coords[TOP_LEFT]
	topright := c.tex_coords[TOP_RIGHT]
	bottomleft := c.tex_coords[BOTTOM_LEFT]
	bottomright := c.tex_coords[BOTTOM_RIGHT]
	triangles := []float32{
		float32(c.X), float32(c.Y), topleft.X(), topleft.Y(),
		float32(c.X + c.W), float32(c.Y), topright.X(), topright.Y(),
		float32(c.X), float32(c.Y + c.H), bottomleft.X(), bottomright.Y(),

		float32(c.X + c.W), float32(c.Y), topright.X(), topright.Y(),
		float32(c.X + c.W), float32(c.Y + c.H), bottomright.X(), bottomright.Y(),
		float32(c.X), float32(c.Y + c.H), bottomleft.X(), bottomright.Y(),
	}

	if c.vao == 0 {
		gl.GenVertexArrays(1, &c.vao)
	}
	if c.vbo == 0 {
		gl.GenBuffers(gl.ARRAY_BUFFER, &c.vbo)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triangles)*4, gl.Ptr(triangles), gl.DYNAMIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(TopContainer.Program(), gl.Str("vertAndTexCoord\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
}

func (c *TexturedComponentImpl) Move(x, y, w, h int) {
	c.X = x
	c.Y = y
	c.W = w
	c.H = h
	c.updateTexture()
}

// A partial component implementation that ignores input and ticks
type StaticComponentImpl struct {
}

func (self *StaticComponentImpl) AcceptsInput() bool {
	return false
}

func (self *StaticComponentImpl) HandleInput(keyEvents []windows.KeyboardInputEvent, mouseEvents []windows.MouseInputEvent) windows.WindowAction {
	return windows.WindowActionNone
}

func (self *StaticComponentImpl) Tick(timedelta float64, keyStates []bool) {

}

// StaticTexturedComponent is a basic component that takes no input and renders a texture
type StaticTexturedComponent struct {
	TexturedComponentImpl
	StaticComponentImpl
}
