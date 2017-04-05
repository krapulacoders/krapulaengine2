package graphics

import (
	gl "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/krapulacoders/krapulaengine2/graphics/errors"
)

type masterLoop struct {
	rendergroups                map[int]*RenderGroup
	nextFreeIndex               int
	clearColorChanged           bool
	clearColor                  [4]float32
	width, height               float32
	viewportChanged             bool
	precalculatedNormalMatrices [2]mgl32.Mat4
}

var mLoop masterLoop

// Normal matrix IDs
const (
	NormalMatrixOrthoOrigo       = iota
	NormalMatrixOrthoScreenCords = iota
)

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
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.VERTEX_PROGRAM_POINT_SIZE)
	gl.Enable(gl.LINE_SMOOTH)
	gl.DepthMask(false)
	gl.Hint(gl.LINE_SMOOTH_HINT, gl.NICEST)
}

func reCalculateNormalMatrices() {
	mLoop.precalculatedNormalMatrices[NormalMatrixOrthoOrigo] = getNormalMatrixOrthoOrigo()
	mLoop.precalculatedNormalMatrices[NormalMatrixOrthoScreenCords] = getNormalMatrixOrthoScreenCords()
}

// SetClearColor sets the clear color.
func SetClearColor(r, g, b, a float32) {
	mLoop.clearColor[0] = r
	mLoop.clearColor[1] = g
	mLoop.clearColor[2] = b
	mLoop.clearColor[3] = a
	mLoop.clearColorChanged = true
}

// SetViewPortSize sets the size of the rendering area.
func SetViewPortSize(x, y int) {
	mLoop.width, mLoop.height = (float32)(x), (float32)(y)
	mLoop.viewportChanged = true
	reCalculateNormalMatrices()
}

// GetViewPortSize returns the size of the rendering area
func GetViewPortSize() mgl32.Vec2 {
	return mgl32.Vec2{(float32)(mLoop.width), (float32)(mLoop.height)}
}

func GetNormalMatrix(id int) mgl32.Mat4 {
	return mLoop.precalculatedNormalMatrices[id]
}

// getNormalMatrixOrthoOrigo returns a normal matrix centered around origo, with the size of the framebuffer
func getNormalMatrixOrthoOrigo() mgl32.Mat4 {
	return mgl32.Ortho(-mLoop.width/2, mLoop.width/2, -mLoop.height/2, mLoop.height/2, -1, 1)
}

// GetNormalMatrixOrthoScreenCords returns a normal matrix covering (0, 0) -> (width, height)
func getNormalMatrixOrthoScreenCords() mgl32.Mat4 {
	return mgl32.Ortho(0, mLoop.width, -mLoop.height, 0, -1, 1)
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
	if mLoop.viewportChanged {
		gl.Viewport(0, 0, int32(mLoop.width), int32(mLoop.height))
		mLoop.viewportChanged = false
	}
	// clear screen
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	errors.AssertGLError(errors.Critical, "glClear")

	//var depthBits int32
	//gl.GetFramebufferAttachmentParameteriv(gl.DRAW_FRAMEBUFFER, gl.DEPTH, gl.FRAMEBUFFER_ATTACHMENT_DEPTH_SIZE, &depthBits)
	//fmt.Printf("depth bits: %v\n", depthBits)
	//errors.AssertGLError(errors.Debug, "glGetFramebufferAttachmentParameteriv")

	for _, g := range mLoop.rendergroups {
		g.Render()
	}
}
