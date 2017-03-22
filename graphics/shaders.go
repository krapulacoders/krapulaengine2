package graphics

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var shaderCache = make(map[string]uint32)

func getCachedShader(file string, shaderType uint32) uint32 {
	if shaderPgm, ok := shaderCache[file]; ok {
		return shaderPgm
	}

	shaderSource, err := ioutil.ReadFile(file)
	if err != nil {
		panic("failed to open shader " + file + " " + err.Error())
	}
	shaderSourceWithNull := string(shaderSource) + "\x00"

	vertexShaderPgm, err := compileShader(shaderSourceWithNull, shaderType)
	if err != nil {
		panic(err.Error())
	}
	return vertexShaderPgm
}

func clearShaderCache() {
	for _, shader := range shaderCache {
		gl.DeleteShader(shader)
	}
}

// compileShader compiles a shader program from a source string
func compileShader(source string, shaderType uint32) (uint32, error) {

	shader := gl.CreateShader(shaderType)

	// go 1.6 cgo workaround
	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		sts := "unknown"
		switch shaderType {
		case gl.VERTEX_SHADER:
			sts = "vertex"
		case gl.FRAGMENT_SHADER:
			sts = "fragment"
		}
		return 0, fmt.Errorf("failed to compile %v shader %v: \"%v\"", sts, source, log)
	}

	return shader, nil
}

// NewProgram compiles a shader program som vertex and fragment sources
// Deprecated: will be removed at some points
func NewProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}
