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
	scene.SetState(STATE_UNINITED)
	return scene
}

func (s *TestScene) Init() {
	s.SetState(STATE_INITED)
}

func TestSimpleSceneImpl(t *testing.T) {
	scene := newTestScene()
	if scene.state != STATE_UNINITED {
		t.Errorf("scene was %v. expected %v", scene.state, STATE_UNINITED)
	}
	scene.Init()
	if scene.state != STATE_INITED {
		t.Errorf("scene was %v. expected %v", scene.state, STATE_INITED)
	}

	scene.Run()
	if scene.state != STATE_RUNNING {
		t.Errorf("scene was %v. expected %v", scene.state, STATE_RUNNING)
	}
}
