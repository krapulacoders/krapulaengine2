package windows

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/krapulacoders/krapulaengine2/graphics"
)

// MainWindow is _the_ window
var MainWindow = newWindow(800, 600)

// Inits glfw
func initGLFW() {

	if err := glfw.Init(); err != nil {
		panic("Can't init glfw: " + err.Error())
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

}

type WindowAction int

const (
	WINDOW_ACTION_NONE = iota
	WINDOW_ACTION_EXIT = iota
)

// Window is an actual window in the OS.
type Window struct {
	window *glfw.Window

	// general window info
	width, height int

	// Loaded scenes. Only one can be active at any time
	scenes        map[string]Scene
	current_scene string

	// A set of overlays that can be rendered on top of the actual game scene
	overlays map[string]Scene
	// many can be active at the same time, keep track of which ones are active
	active_overlays map[string]bool

	// set to true when the main loop should quit
	quit bool

	// internal input channels
	keyInput       chan KeyboardInputEvent
	mouse_input    chan MouseInputEvent
	mouseX, mouseY float32
	keyEvents      []KeyboardInputEvent
	mouseEvents    []MouseInputEvent
	keyStates      []bool
}

// Allocates a new Window
func newWindow(width, height int) *Window {
	w := new(Window)
	w.scenes = make(map[string]Scene)
	w.overlays = make(map[string]Scene)
	w.active_overlays = make(map[string]bool)
	w.current_scene = ""
	w.width = width
	w.height = height
	w.keyInput = make(chan KeyboardInputEvent, 100)
	w.mouse_input = make(chan MouseInputEvent, 100)
	w.keyEvents = make([]KeyboardInputEvent, 0, 100)
	w.mouseEvents = make([]MouseInputEvent, 0, 100)
	w.keyStates = make([]bool, glfw.KeyLast+1)
	return w
}

// Init Creates the window and initializes the graphics thread
func Init() {
	initGLFW()

	window, err := glfw.CreateWindow(MainWindow.width, MainWindow.height, "Testing", nil, nil)
	if err != nil {
		panic("Can't create window")
	}
	MainWindow.window = window
	graphics.InitMasterLoop(
		func() {
			MainWindow.window.SwapBuffers()
		}, func() {
			MainWindow.window.MakeContextCurrent()
		})
	window.SetKeyCallback(keyEventHandler())
	window.SetCursorPosCallback(cursorPosHandler())
	window.SetMouseButtonCallback(mouseButtonEventHandler())

}
func keyEventHandler() glfw.KeyCallback {
	return func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		MainWindow.keyInput <- KeyboardInputEvent{key, scancode, action, mods}
	}
}

func mouseButtonEventHandler() glfw.MouseButtonCallback {
	return func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		MainWindow.mouse_input <- MouseInputEvent{button, action, mod, MainWindow.mouseX, MainWindow.mouseY}
	}
}

func cursorPosHandler() glfw.CursorPosCallback {
	return func(w *glfw.Window, xpos float64, ypos float64) {
		MainWindow.mouseX = float32(xpos)
		MainWindow.mouseY = float32(ypos)
	}
}

// The mainloop initializes its scenes and then enters the loop that runs the game.
// It directs input events to the current scene and any overlays and renders them.
// When the loop ends, the scenes are terminated
func MainLoop() {

	if MainWindow.window == nil {
		Init()
	}

	defer func() {
		// Terminate all Scenes at exit
		for _, scene := range MainWindow.scenes {
			scene.Exit()
		}
		for _, scene := range MainWindow.overlays {
			scene.Exit()
		}
		MainWindow.window = nil
		glfw.Terminate()
	}()

	initScenes()

	MainWindow.quit = false
	old_time := glfw.GetTime()
	for !MainWindow.quit {
		new_time := glfw.GetTime()
		timedelta := new_time - old_time
		old_time = new_time

		processInput()
		tick(timedelta)
		render()

	}
}

// Init() scenes, then Run() them
func initScenes() {
	// Init all scenes
	for _, scene := range MainWindow.scenes {
		if !scene.IsInited() {
			scene.Init()
		}
	}
	for _, scene := range MainWindow.overlays {
		if !scene.IsInited() {
			scene.Init()
		}
	}

	// Mark all scenes as running.
	// This is done after Init because scenes may depend on each other being inited.
	for _, scene := range MainWindow.scenes {
		scene.Run()
	}
	for _, scene := range MainWindow.overlays {
		scene.Run()
	}
}

func Exit() {
	MainWindow.quit = true
	graphics.DeinitMasterLoop()
}

func AddScene(id string, scene Scene) {
	if _, exists := MainWindow.scenes[id]; exists {
		panic("Tried adding scene twice to window")
	}
	MainWindow.scenes[id] = scene
	if MainWindow.current_scene == "" {
		SetCurrentScene(id)
	}
}

func SetCurrentScene(id string) {
	MainWindow.current_scene = id
}

func AddOverlay(id string, scene Scene) {
	if _, exists := MainWindow.overlays[id]; exists {
		panic("Tried adding overlay twice to window")
	}
	MainWindow.overlays[id] = scene
}

func processInput() {

	// process input
	glfw.PollEvents()
	// clear the slice
	MainWindow.keyEvents = MainWindow.keyEvents[:0]
	// key events
	has_input := true
	for has_input {
		select {
		case keyEvent := <-MainWindow.keyInput:
			MainWindow.keyEvents = append(MainWindow.keyEvents, keyEvent)
			if keyEvent.Action == glfw.Release {
				MainWindow.keyStates[keyEvent.Key] = false
			} else {
				MainWindow.keyStates[keyEvent.Key] = true
			}
		default:
			has_input = false
		}
	}

	// mouse events
	// clear the slice
	MainWindow.mouseEvents = MainWindow.mouseEvents[:0]
	has_input = true
	for has_input {
		select {
		case mouseEvent := <-MainWindow.mouse_input:
			MainWindow.mouseEvents = append(MainWindow.mouseEvents, mouseEvent)

		default:
			has_input = false
		}
	}

	// Send inputs to scenes
	// Overlays get input first.

	for _, overlay := range MainWindow.overlays {
		if overlay.AcceptsInput() {
			action := overlay.HandleInput(MainWindow.keyEvents, MainWindow.mouseEvents)
			switch action {
			case WINDOW_ACTION_EXIT:
				Exit()
			}
		}
	}

	for _, scene := range MainWindow.scenes {
		if scene.AcceptsInput() {
			action := scene.HandleInput(MainWindow.keyEvents, MainWindow.mouseEvents)
			switch action {
			case WINDOW_ACTION_EXIT:
				Exit()
			}
		}
	}
}

func tick(timedelta float64) {
	for _, scene := range MainWindow.scenes {
		scene.Tick(timedelta, MainWindow.keyStates)
	}
	for _, overlay := range MainWindow.overlays {
		overlay.Tick(timedelta, MainWindow.keyStates)
	}
}

func render() {

	if MainWindow.window.ShouldClose() {
		MainWindow.quit = true
	}
}

func GetSize() (int, int) {
	return MainWindow.window.GetSize()
}
