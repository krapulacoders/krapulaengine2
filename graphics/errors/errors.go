package errors

import (
	"fmt"

	gl "github.com/go-gl/gl/v3.3-core/gl"
)

// ErrorImportance defines how important/dangerous an error is.
// possible values: Critical, Normal, Debug
type ErrorImportance int

// values for ErrorImportance
const (
	Critical ErrorImportance = iota
	Normal
	Debug
)

var errorImportanceMap = map[ErrorImportance]string{
	Critical: "Critical",
	Normal:   "Normal",
	Debug:    "Debug",
}

var debug = true

var errorMappings = map[uint32]string{
	gl.INVALID_ENUM:                  "GL_INVALID_ENUM",
	gl.INVALID_VALUE:                 "GL_INVALID_VALUE",
	gl.INVALID_OPERATION:             "GL_INVALID_OPERATION",
	gl.STACK_OVERFLOW:                "GL_STACK_OVERFLOW",
	gl.STACK_UNDERFLOW:               "GL_STACK_UNDERFLOW",
	gl.OUT_OF_MEMORY:                 "GL_OUT_OF_MEMORY",
	gl.INVALID_FRAMEBUFFER_OPERATION: "GL_INVALID_FRAMEBUFFER_OPERATION",
	gl.CONTEXT_LOST:                  "GL_CONTEXT_LOST",
}

// SetDebugMode enabled/disables debug mode.
func SetDebugMode(enabled bool) {
	debug = enabled
}

// GetDebugMode returns current debug status
func GetDebugMode() bool {
	return debug
}

// AssertGLError panics or logs depending on debug mode and given importance of error
func AssertGLError(importance ErrorImportance, comment string) {
	if !debug && importance == Debug {
		return
	}

	err := gl.GetError()
	if err == gl.NO_ERROR {
		return
	}

	errString := "glError: " + errorMappings[err] + ". " + comment

	switch importance {
	case Critical:
		panic(errString)
	case Normal:
		if debug {
			panic(errString)
		} else {
			fmt.Println(errString)
		}
	case Debug:
		fmt.Println(errString)
	}

}

// LogError logs a error in the graphics log
func LogError(importance ErrorImportance, text string) {
	if debug && importance != Debug {
		return
	}
	fmt.Println(errorImportanceMap[importance] + ": " + text)
}
