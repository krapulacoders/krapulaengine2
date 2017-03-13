package graphics

import (
	"fmt"
	gl "github.com/go-gl/gl/v3.3-core/gl"
	"strings"
)

// Manager holds a collection of graphical objects that can be efficiently rendered together
type Manager interface {

	// Init allocates required resources
	Init() error

	// Render renders the objects.
	Render()

	// Deinit deallocates any resources
	Deinit() error

	test()
}

// ManagerImpl implements basic shading handling
type ManagerImpl struct {
	shaderPgm            uint32
	vertexShader         string
	fragmentShader       string
	shaderPgmNeedsRelink bool

	shaderBinder ShaderBinder
}

// SetShaderFile sets the shader file according to file ending.
// possible endings: .frag or .vert
// The programs are not compiled until the first time they are run.
// The reason it does not compile the programs is to stay away from the graphics thread's business
func (self *ManagerImpl) SetShaderFile(shaderFile string) {
	if strings.HasSuffix(shaderFile, ".frag") {
		self.fragmentShader = shaderFile
	} else if strings.HasSuffix(shaderFile, ".vert") {
		self.vertexShader = shaderFile
	} else {
		panic("unknown shader file extension " + shaderFile)
	}
	self.shaderPgmNeedsRelink = true
}

func (self *ManagerImpl) LinkProgram() error {
	if self.shaderPgmNeedsRelink {
		vertexShaderPgm := getCachedShader(self.fragmentShader, gl.VERTEX_SHADER)
		fragmentShaderPgm := getCachedShader(self.vertexShader, gl.FRAGMENT_SHADER)

		self.shaderPgm = gl.CreateProgram()
		gl.AttachShader(self.shaderPgm, vertexShaderPgm)
		gl.AttachShader(self.shaderPgm, fragmentShaderPgm)
		gl.LinkProgram(self.shaderPgm)

		var status int32
		gl.GetProgramiv(self.shaderPgm, gl.LINK_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetProgramiv(self.shaderPgm, gl.INFO_LOG_LENGTH, &logLength)

			log := strings.Repeat("\x00", int(logLength+1))
			gl.GetProgramInfoLog(self.shaderPgm, logLength, nil, gl.Str(log))
			panic(fmt.Errorf("failed to link program: %v", log))
		}
		// then do first-time binding
		self.shaderBinder.BindFirstTime()
	}

	gl.UseProgram(self.shaderPgm)

	return nil

}
