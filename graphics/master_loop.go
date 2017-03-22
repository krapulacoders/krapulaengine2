package graphics

import (
	"time"

	gl "github.com/go-gl/gl/v3.3-core/gl"
)

const (
	quit    = iota
	inited  = iota
	running = iota
	paused  = iota
)

type masterLoop struct {
	managers          map[int]RenderGroup
	nextFreeIndex     int
	state             int
	swapBuffers       func()
	makeContextActive func()
}

var mLoop masterLoop

// InitMasterLoop must be called before starting the graphics system using Start().
func InitMasterLoop(swapBuffers, makeContextActive func()) {
	mLoop = masterLoop{make(map[int]RenderGroup), 0, inited, swapBuffers, makeContextActive}

	// context must be set before gl.Init
	makeContextActive()
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}
}

// AddManager assigns an id to the manager and returns it.
// This is NOT threadsafe and should not be ran while the graphics loop is running for now
func AddManager(m RenderGroup) int {
	mLoop.managers[mLoop.nextFreeIndex] = m
	mLoop.nextFreeIndex++
	return mLoop.nextFreeIndex
}

// RemoveManager removes the specified id.
// This is NOT threadsafe and should not be ran while the graphics loop is running for now
func RemoveManager(id int) {
	delete(mLoop.managers, id)
}

// GetManager returns the specified manager
func GetManager(id int) RenderGroup {
	return mLoop.managers[id]
}

// Start starts a goroutine for the graphics loop if one doesn't already exist.
func Start() {
	if mLoop.state == inited {
		go mainLoop()
	} else if mLoop.state == paused {
		mLoop.state = running
	} else {
		panic("started master loop in state " + string(mLoop.state))
	}
}

// Pause the render loop. It will still upload new textures on demand while paused.
func Pause() {
	if mLoop.state == quit {
		panic("started master loop in state " + string(mLoop.state))
	}
	mLoop.state = paused
}

// DeinitMasterLoop stops the graphics thread. It cannot be restarted after this operation.
func DeinitMasterLoop() {
	mLoop.state = quit
}

func mainLoop() {
	for mLoop.state != quit {
		if mLoop.state == paused {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		mLoop.makeContextActive()
		mLoop.swapBuffers()
		// clear screen
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, manager := range mLoop.managers {
			manager.Render()
		}
	}
}
