package windows

import (
	"testing"
)

type TestScene struct {
	SimpleSceneImpl
}

func newTestScene() *TestScene {
	scene := new(TestScene)
	scene.state = 999
	scene.SetState(StateUninited)
	return scene
}

func (s *TestScene) Init() {
	s.SetState(StateInited)
}

func TestSimpleSceneImpl(t *testing.T) {
	scene := newTestScene()
	if scene.state != StateUninited {
		t.Errorf("scene was %v. expected %v", scene.state, StateUninited)
	}
	scene.Init()
	if scene.state != StateInited {
		t.Errorf("scene was %v. expected %v", scene.state, StateInited)
	}

	scene.Run()
	if scene.state != StateRunning {
		t.Errorf("scene was %v. expected %v", scene.state, StateRunning)
	}
}
