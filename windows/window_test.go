package windows

import (
	"testing"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
)

/* Tests that glfw3 is properly set up and working on the system.
 */
func TestGLFW(t *testing.T) {

	if err := glfw.Init(); err != nil {
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

// Create a window and have the mainloop run for 50 milliseconds,
// then close it.
func TestWindow(t *testing.T) {
	Init()
	go func() {
		time.Sleep(50 * time.Millisecond)
		Exit()
	}()
	MainLoop()
}
