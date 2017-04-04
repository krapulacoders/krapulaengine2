package graphics

import (
	"fmt"
	"strings"

	gl "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/krapulacoders/krapulaengine2/graphics/errors"
)

// RenderGroup implements basic shading handling
type RenderGroup struct {
	id                   string
	shaderPgm            uint32
	vertexShader         string
	fragmentShader       string
	shaderPgmNeedsRelink bool
	depthTestEnabled     bool
	depthTestFunc        uint32
	blendEnabled         bool
	blendFunc1           uint32
	blendFunc2           uint32

	impl RenderGroupImplementation
}

// RenderGroupImplementation is an actual implementation of a render group
type RenderGroupImplementation interface {

	// Deinit deallocates any resources
	Deinit()

	InitShader()

	Render()
}

// SetShaderFile sets the shader file according to file ending.
// possible endings: .frag or .vert
// The programs are not compiled until the first time they are run.
// The reason it does not compile the programs is to stay away from the graphics thread's business
func (g *RenderGroup) SetShaderFile(shaderFile string) {
	if strings.HasSuffix(shaderFile, ".frag") {
		g.fragmentShader = shaderFile
	} else if strings.HasSuffix(shaderFile, ".vert") {
		g.vertexShader = shaderFile
	} else {
		panic("unknown shader file extension " + shaderFile)
	}
	g.shaderPgmNeedsRelink = true
}

// GetShaderProgram returns the shader pgm
func (g *RenderGroup) GetShaderProgram() uint32 {
	return g.shaderPgm
}

func (g *RenderGroup) activateShaderProgram() {
	if g.shaderPgmNeedsRelink {
		fmt.Println("initing shader program")
		gl.UseProgram(0)
		errors.AssertGLError(errors.Critical, "gl.UseProgram")
		vertexShaderPgm := getCachedShader(g.vertexShader, gl.VERTEX_SHADER)
		fragmentShaderPgm := getCachedShader(g.fragmentShader, gl.FRAGMENT_SHADER)

		g.shaderPgm = gl.CreateProgram()
		gl.AttachShader(g.shaderPgm, vertexShaderPgm)
		errors.AssertGLError(errors.Critical, "gl.AttachShader1")
		gl.AttachShader(g.shaderPgm, fragmentShaderPgm)
		errors.AssertGLError(errors.Critical, "gl.AttachShader2")
		gl.LinkProgram(g.shaderPgm)
		errors.AssertGLError(errors.Critical, "gl.LinkProgram")

		var status int32
		gl.GetProgramiv(g.shaderPgm, gl.LINK_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetProgramiv(g.shaderPgm, gl.INFO_LOG_LENGTH, &logLength)

			log := strings.Repeat("\x00", int(logLength+1))
			gl.GetProgramInfoLog(g.shaderPgm, logLength, nil, gl.Str(log))
			panic(fmt.Errorf("failed to link program: %v", log))
		}

		if !gl.IsProgram(g.shaderPgm) {
			panic("newly compiled shader pgm is not a shader pgm")
		}

		// then do first-time binding
		g.impl.InitShader()
		g.shaderPgmNeedsRelink = false
	}

	if !gl.IsProgram(g.shaderPgm) {
		panic("previously compiled shader pgm is not a shader pgm")
	}

	errors.AssertGLError(errors.Debug, fmt.Sprintf("before gl.UseProgram(%v)", g.shaderPgm))
	gl.UseProgram(g.shaderPgm)
	errors.AssertGLError(errors.Critical, fmt.Sprintf("gl.UseProgram(%v)", g.shaderPgm))

}

// Deinit calls the initer to deinitiate
func (g *RenderGroup) Deinit() {
	if g.shaderPgm != 0 {
		gl.UseProgram(0)
		gl.DeleteProgram(g.shaderPgm)

	}
	g.shaderPgmNeedsRelink = true
	g.impl.Deinit()
}

// Render activates the shader program and calls the renderer
func (g *RenderGroup) Render() {
	errors.AssertGLError(errors.Debug, "RenderGroup.Render")
	g.activateShaderProgram()

	if g.depthTestEnabled {
		gl.Enable(gl.DEPTH_TEST)
		errors.AssertGLError(errors.Debug, "gl enable depth test")
		gl.DepthFunc(g.depthTestFunc)
		errors.AssertGLError(errors.Debug, "gl set depth func")
	} else {
		gl.Disable(gl.DEPTH_TEST)
		errors.AssertGLError(errors.Debug, "gl disable depth test")
	}

	if g.blendEnabled {
		gl.Enable(gl.BLEND)
		errors.AssertGLError(errors.Debug, "gl enable blending")
		gl.BlendFunc(g.blendFunc1, g.blendFunc2)
		errors.AssertGLError(errors.Debug, "gl set blend func")
	} else {
		gl.Disable(gl.BLEND)
		errors.AssertGLError(errors.Debug, "gl disable blending")
	}

	g.impl.Render()
}

// SetDepthTestMode sets the depth test details used when rendering this group
func (g *RenderGroup) SetDepthTestMode(enabled bool, depthFunc uint32) {
	g.depthTestEnabled = enabled
	g.depthTestFunc = depthFunc
}

// SetBlendingMode sets the blending mode used when rendering this group
func (g *RenderGroup) SetBlendingMode(enabled bool, blendFunc1, blendFunc2 uint32) {
	g.blendEnabled = enabled
	g.blendFunc1 = blendFunc1
	g.blendFunc2 = blendFunc2
}

// NewRenderGroup creates a RenderGroup.
func NewRenderGroup(id string, impl RenderGroupImplementation) *RenderGroup {
	g := new(RenderGroup)
	g.impl = impl
	g.shaderPgm = 0
	g.shaderPgmNeedsRelink = true
	return g
}
