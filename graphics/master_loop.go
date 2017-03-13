package graphics

const (
	quit    = iota
	inited  = iota
	running = iota
	paused  = iota
)

type masterLoop_t struct {
	managers []Manager
	state    int
}

var masterLoop masterLoop_t

func InitMasterLoop() {
	masterLoop = masterLoop_t{make([]Manager, 1000), inited}

}

func Start() {
	if masterLoop.state == inited {

	} else if masterLoop.state == paused {
		masterLoop.state = running
	} else {
		panic("started master loop in state " + string(masterLoop.state))
	}
}

func Pause() {
	if masterLoop.state == quit {
		panic("started master loop in state " + string(masterLoop.state))
	}
	masterLoop.state = paused
}

func DeinitMasterLoop() {
	masterLoop.state = quit
}

func mainLoop() {
	for masterLoop.state != quit {

		for _, manager := range masterLoop.managers {
			manager.Render()
		}
	}
}
