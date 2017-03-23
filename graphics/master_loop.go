package graphics

import (
	gl "github.com/go-gl/gl/v3.3-core/gl"
)

type masterLoop struct {
	rendergroups      map[int]*RenderGroup
	nextFreeIndex     int
	clearColorChanged bool
	clearColor        [4]float32
}

var mLoop masterLoop

// InitMasterLoop must be called before starting the graphics system using Start().
func InitMasterLoop() {
	mLoop = masterLoop{
		rendergroups: make(map[int]*RenderGroup),
	}

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.ClearColor(1, 0, 1, 0)
}

// SetClearColor sets the clear color.
func SetClearColor(r, g, b, a float32) {
	mLoop.clearColorChanged = true
	mLoop.clearColor[0] = r
	mLoop.clearColor[1] = g
	mLoop.clearColor[2] = b
	mLoop.clearColor[3] = a
}

// AddRenderGroup assigns an id to the manager and returns it.
// This is NOT threadsafe and should not be ran while the graphics loop is running for now
func AddRenderGroup(g *RenderGroup) int {
	mLoop.rendergroups[mLoop.nextFreeIndex] = g
	mLoop.nextFreeIndex++
	return mLoop.nextFreeIndex
}

// RemoveRenderGroup removes the specified id.
// This is NOT threadsafe and should not be ran while the graphics loop is running for now
func RemoveRenderGroup(id int) {
	delete(mLoop.rendergroups, id)
}

// GetRenderGroup returns the specified manager
func GetRenderGroup(id int) *RenderGroup {
	return mLoop.rendergroups[id]
}

// DeinitMasterLoop deinits all rendergroups
func DeinitMasterLoop() {
	for _, g := range mLoop.rendergroups {
		g.Deinit()
	}
}

// Render the registered render groups
func Render() {
	if mLoop.clearColorChanged {
		gl.ClearColor(mLoop.clearColor[0], mLoop.clearColor[1], mLoop.clearColor[2], mLoop.clearColor[3])
		mLoop.clearColorChanged = false
	}
	// clear screen
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, g := range mLoop.rendergroups {
		g.Render()
	}
}
