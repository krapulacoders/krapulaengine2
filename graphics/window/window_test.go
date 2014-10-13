package window

import(
	"testing"
	glfw "github.com/go-gl/glfw3"
)

/* A temporary test to figure out how windows work and what is needed for glfw.
 *
 */
func TestWindow(t *testing.T) {
    glfw.SetErrorCallback(errorCallback)

    if !glfw.Init() {
        panic("Can't init glfw!")
    }
    defer glfw.Terminate()

    window, err := glfw.CreateWindow(800, 600, "Testing", nil, nil)
    if err != nil {
        panic(err)
    }

    window.MakeContextCurrent()

    for !window.ShouldClose() {
        //Do OpenGL stuff
        window.SwapBuffers()
        glfw.PollEvents()
    }
}