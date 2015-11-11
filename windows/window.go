package windows

import (
	gl "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

// Inits glfw
func InitWindowing() {

	if err := glfw.Init(); err != nil {
		panic("Can't init glfw: " + err.Error())
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

}

//
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
	key_input        chan KeyboardInputEvent
	mouse_input      chan MouseInputEvent
	mouse_x, mouse_y float32
	key_events       []KeyboardInputEvent
	mouse_events     []MouseInputEvent
}

// Allocates a new Window
func NewWindow(width, height int) *Window {
	w := new(Window)
	w.scenes = make(map[string]Scene)
	w.overlays = make(map[string]Scene)
	w.active_overlays = make(map[string]bool)
	w.current_scene = ""
	w.width = width
	w.height = height
	w.key_input = make(chan KeyboardInputEvent, 100)
	w.mouse_input = make(chan MouseInputEvent, 100)
	w.key_events = make([]KeyboardInputEvent, 0, 100)
	w.mouse_events = make([]MouseInputEvent, 0, 100)
	return w
}

// Creates the window and initializes OpenGL
func (self *Window) Init() {

	window, err := glfw.CreateWindow(self.width, self.height, "Testing", nil, nil)
	if err != nil {
		panic("Can't create window")
	}
	self.window = window
	self.window.MakeContextCurrent()
	window.SetKeyCallback(keyEventHandler(self))
	window.SetCursorPosCallback(cursorPosHandler(self))
	window.SetMouseButtonCallback(mouseButtonEventHandler(self))

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}
}
func keyEventHandler(window *Window) glfw.KeyCallback {
	return func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		window.key_input <- KeyboardInputEvent{key, scancode, action, mods}
	}
}

func mouseButtonEventHandler(window *Window) glfw.MouseButtonCallback {
	return func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		window.mouse_input <- MouseInputEvent{button, action, mod, window.mouse_x, window.mouse_y}
	}
}

func cursorPosHandler(window *Window) glfw.CursorPosCallback {
	return func(w *glfw.Window, xpos float64, ypos float64) {
		window.mouse_x = float32(xpos)
		window.mouse_y = float32(ypos)
	}
}

// The mainloop initializes its scenes and then enters the loop that runs the game.
// It directs input events to the current scene and any overlays and renders them.
// When the loop ends, the scenes are terminated
func (self *Window) MainLoop() {

	if self.window == nil {
		self.Init()
	}

	defer func() {
		// Terminate all Scenes at exit
		for _, scene := range self.scenes {
			scene.Exit()
		}
		for _, scene := range self.overlays {
			scene.Exit()
		}
		self.window = nil
		glfw.Terminate()
	}()

	self.initScenes()

	self.quit = false
	old_time := glfw.GetTime()
	for !self.quit {
		new_time := glfw.GetTime()
		timedelta := new_time - old_time
		old_time = new_time

		self.processInput()
		self.tick(timedelta)
		self.render()

	}
}

// Init() scenes, then Run() them
func (self *Window) initScenes() {
	// Init all scenes
	for _, scene := range self.scenes {
		if !scene.IsInited() {
			scene.Init()
		}
	}
	for _, scene := range self.overlays {
		if !scene.IsInited() {
			scene.Init()
		}
	}

	// Mark all scenes as running.
	// This is done after Init because scenes may depend on each other being inited.
	for _, scene := range self.scenes {
		scene.Run()
	}
	for _, scene := range self.overlays {
		scene.Run()
	}
}

func (self *Window) Exit() {
	self.quit = true
}

func (self *Window) AddScene(id string, scene Scene) {
	if _, exists := self.scenes[id]; exists {
		panic("Tried adding scene twice to window")
	}
	self.scenes[id] = scene
	if self.current_scene == "" {
		self.SetCurrentScene(id)
	}
}

func (self *Window) SetCurrentScene(id string) {
	self.current_scene = id
}

func (self *Window) AddOverlay(id string, scene Scene) {
	if _, exists := self.overlays[id]; exists {
		panic("Tried adding overlay twice to window")
	}
	self.overlays[id] = scene
}

func (self *Window) processInput() {

	// process input
	glfw.PollEvents()
	// clear the slice
	self.key_events = self.key_events[:0]
	// key events
	has_input := true
	for has_input {
		select {
		case keyEvent := <-self.key_input:
			self.key_events = append(self.key_events, keyEvent)

		default:
			has_input = false
		}
	}

	// mouse events
	// clear the slice
	self.mouse_events = self.mouse_events[:0]
	has_input = true
	for has_input {
		select {
		case mouseEvent := <-self.mouse_input:
			self.mouse_events = append(self.mouse_events, mouseEvent)

		default:
			has_input = false
		}
	}

	// Send inputs to scenes
	// Overlays get input first.

	for _, overlay := range self.overlays {
		if overlay.AcceptsInput() {
			overlay.HandleInput(self.key_events, self.mouse_events)
		}
	}

	for _, scene := range self.scenes {
		if scene.AcceptsInput() {
			scene.HandleInput(self.key_events, self.mouse_events)
		}
	}
}

func (self *Window) tick(timedelta float64) {
	for _, scene := range self.scenes {
		scene.Tick(timedelta)
	}
	for _, overlay := range self.overlays {
		overlay.Tick(timedelta)
	}
}

func (self *Window) render() {

	self.window.MakeContextCurrent()

	// clear screen
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// render current scene
	if self.current_scene != "" {
		self.scenes[self.current_scene].Render()
	}

	// and any overlays
	for _, overlay := range self.overlays {
		overlay.Render()
	}

	self.window.SwapBuffers()
	if self.window.ShouldClose() {
		self.quit = true
	}
}
