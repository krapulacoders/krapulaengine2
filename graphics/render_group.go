package graphics

import (
	"fmt"
	"strings"

	gl "github.com/go-gl/gl/v3.3-core/gl"
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
		gl.UseProgram(0)
		vertexShaderPgm := getCachedShader(g.vertexShader, gl.VERTEX_SHADER)
		fragmentShaderPgm := getCachedShader(g.fragmentShader, gl.FRAGMENT_SHADER)

		g.shaderPgm = gl.CreateProgram()
		gl.AttachShader(g.shaderPgm, vertexShaderPgm)
		gl.AttachShader(g.shaderPgm, fragmentShaderPgm)
		gl.LinkProgram(g.shaderPgm)

		var status int32
		gl.GetProgramiv(g.shaderPgm, gl.LINK_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetProgramiv(g.shaderPgm, gl.INFO_LOG_LENGTH, &logLength)

			log := strings.Repeat("\x00", int(logLength+1))
			gl.GetProgramInfoLog(g.shaderPgm, logLength, nil, gl.Str(log))
			panic(fmt.Errorf("failed to link program: %v", log))
		}
		// then do first-time binding
		g.impl.InitShader()
	}

	gl.UseProgram(g.shaderPgm)

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
	g.activateShaderProgram()

	if g.depthTestEnabled {
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(g.depthTestFunc)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}

	if g.blendEnabled {
		gl.Enable(gl.BLEND)
		gl.BlendFunc(g.blendFunc1, g.blendFunc2)
	} else {
		gl.Disable(gl.BLEND)
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
