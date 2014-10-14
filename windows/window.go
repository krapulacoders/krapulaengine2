package windows

import(
	glfw "github.com/go-gl/glfw3"
	"time"
)


type Window struct {
	window *glfw.Window

	// Loaded scenes. Only one can be active at any time
	scenes map[string] Scene
	current_scene string

	// A set of overlays that can be rendered on top of the actual game scene
	overlays map[string] Scene
	// many can be active at the same time, keep track of which ones are active
	active_overlays map[string] bool

	//
	quit bool
}

func NewWindow() {
	w := new(Window)
	w.scenes = make(map[string] Scene)
	w.overlays = make(map[string] Scene)
	w.active_overlays = make(map[string] bool)
}

func (self *Window) init() {
	glfw.SetErrorCallback(func(code glfw.ErrorCode, desc string) {
        panic(desc)
    })

    if !glfw.Init() {
        panic("Can't init glfw!")
    }

    window, err := glfw.CreateWindow(800, 600, "Testing", nil, nil)
    if err != nil {
        panic("Can't create window")
    }
    self.window = window
}

func (self *Window) MainLoop() {
	// create window and stuff
	self.init()

	defer glfw.Terminate()

	self.quit = false
	old_time := time.Now().UnixNano()
	for !self.quit {
		new_time := time.Now().UnixNano()
		timedelta := float64(new_time - old_time) / 1e9
		old_time = new_time

		// process input
		glfw.PollEvents()

		// tick() all scenes
		for _, sc := range self.scenes {
			sc.Tick(timedelta)
		}
		for _, sc := range self.overlays {
			sc.Tick(timedelta)
		}

		// render current scene

		// and any overlays


		self.window.SwapBuffers()
    	
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
}

func (self *Window) AddOverlay(id string, scene Scene) {
	if _, exists := self.overlays[id]; exists {
		panic("Tried adding overlay twice to window")
	}
	self.overlays[id] = scene
}