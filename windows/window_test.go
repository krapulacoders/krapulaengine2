package windows

import(
	"testing"
	glfw "github.com/go-gl/glfw3"
)

/* Tests that a glfw3 window can be created.
 */
func TestGLFW(t *testing.T) {
    glfw.SetErrorCallback(func(code glfw.ErrorCode, desc string) {
        t.Error(desc)
    })

    if !glfw.Init() {
        t.Error("Can't init glfw!")
        t.FailNow()
    }
    defer glfw.Terminate()

    window, err := glfw.CreateWindow(800, 600, "Testing", nil, nil)
    if err != nil {
        t.Errorf("Can't create window: %v", err.Error())
        t.FailNow()
    }

    window.MakeContextCurrent()

    //Do OpenGL stuff
    window.SwapBuffers()
    glfw.PollEvents()
}