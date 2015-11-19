package gui

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	. "github.com/krapulacoders/krapulaengine2/windows"
)

type Component interface {

	// Returns true if the Scene cares about input
	AcceptsInput() bool

	// Processes input events
	HandleInput(key_events []KeyboardInputEvent, mouse_events []MouseInputEvent) WindowAction

	// Update the game frame and process time-dependant input
	Tick(timedelta float64, key_states []bool)

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

func (self *TexturedComponentImpl) Render() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE0, self.texture)
	gl.DrawArrays(gl.TRIANGLES, 0, 6) // 6 vertices starting from position 0

}

// assign a texture and texture coordinates to this component
func (self *TexturedComponentImpl) SetTexture(texture uint32, topleft, topright, bottomleft, bottomright mgl32.Vec2) {
	self.texture = texture
	self.tex_coords = make([]mgl32.Vec2, 4)
	self.tex_coords[TOP_LEFT] = topleft
	self.tex_coords[TOP_RIGHT] = topright
	self.tex_coords[BOTTOM_LEFT] = bottomleft
	self.tex_coords[BOTTOM_RIGHT] = bottomright
	self.updateTexture()
}

func (self *TexturedComponentImpl) updateTexture() {
	topleft := self.tex_coords[TOP_LEFT]
	topright := self.tex_coords[TOP_RIGHT]
	bottomleft := self.tex_coords[BOTTOM_LEFT]
	bottomright := self.tex_coords[BOTTOM_RIGHT]
	triangles := []float32{
		float32(self.X), float32(self.Y), topleft.X(), topleft.Y(),
		float32(self.X + self.W), float32(self.Y), topright.X(), topright.Y(),
		float32(self.X), float32(self.Y + self.H), bottomleft.X(), bottomright.Y(),

		float32(self.X + self.W), float32(self.Y), topright.X(), topright.Y(),
		float32(self.X + self.W), float32(self.Y + self.H), bottomright.X(), bottomright.Y(),
		float32(self.X), float32(self.Y + self.H), bottomleft.X(), bottomright.Y(),
	}

	if self.vao == 0 {
		gl.GenVertexArrays(1, &self.vao)
	}
	if self.vbo == 0 {
		gl.GenBuffers(gl.ARRAY_BUFFER, &self.vbo)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triangles)*4, gl.Ptr(triangles), gl.DYNAMIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(TopContainer.Program(), gl.Str("vertAndTexCoord\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
}

func (self *TexturedComponentImpl) Move(x, y, w, h int) {
	self.X = x
	self.Y = y
	self.W = w
	self.H = h
	self.updateTexture()
}

// A partial component implementation that ignores input and ticks
type StaticComponentImpl struct {
}

func (self *StaticComponentImpl) AcceptsInput() bool {
	return false
}

func (self *StaticComponentImpl) HandleInput(key_events []KeyboardInputEvent, mouse_events []MouseInputEvent) WindowAction {
	return WINDOW_ACTION_NONE
}

func (self *StaticComponentImpl) Tick(timedelta float64, key_states []bool) {

}

// A basic component that takes no input and renders a texture
type StaticTexturedComponent struct {
	TexturedComponentImpl
	StaticComponentImpl
}
