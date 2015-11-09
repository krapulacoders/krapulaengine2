package windows

import (
	gl "github.com/go-gl/gl/v3.3-core/gl"
	glfw "github.com/go-gl/glfw3"
	"time"
)

func InitWindowing() {
	glfw.SetErrorCallback(func(code glfw.ErrorCode, desc string) {
		panic(desc)
	})

	if !glfw.Init() {
		panic("Can't init glfw!")
	}

}

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

	//
	quit bool
}

func NewWindow(width, height int) *Window {
	w := new(Window)
	w.scenes = make(map[string]Scene)
	w.overlays = make(map[string]Scene)
	w.active_overlays = make(map[string]bool)
	w.current_scene = ""
	w.width = width
	w.height = height
	return w
}

func (self *Window) Init() {

	window, err := glfw.CreateWindow(self.width, self.height, "Testing", nil, nil)
	if err != nil {
		panic("Can't create window")
	}
	self.window = window
	self.window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}
}

func (self *Window) MainLoop() {

	if self.window == nil {
		self.Init()
	}

	defer glfw.Terminate()

	self.quit = false
	old_time := time.Now().UnixNano()
	for !self.quit {
		new_time := time.Now().UnixNano()
		timedelta := float64(new_time-old_time) / 1e9
		old_time = new_time

		// process input
		glfw.PollEvents()

		// tick() all scenes
		for _, scene := range self.scenes {
			scene.Tick(timedelta)
		}
		for _, overlay := range self.overlays {
			overlay.Tick(timedelta)
		}

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
