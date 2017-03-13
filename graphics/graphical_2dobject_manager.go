package graphics

// Graphical2DObjectID ID for 2d objects
type Object2DID uint32

// Graphical2DObject struct containing rendering details
type Graphical2DObject struct {
	X, Y, W, H float32
	texture    uint32
	rotation   float32
}

type Manager2D interface {
	Manager

	AddObject(id Object2DID, o Graphical2DObject)

	removeObject(id Object2DID)
}

type Manager2DImpl struct {
	ManagerImpl
}

func (self *Manager2DImpl) Init() error {
	return nil
}

func (self *Manager2DImpl) DeInit() error {
	return nil
}

func NewManager2DImpl(shaderFile string) *Manager2DImpl {
	manager := new(Manager2DImpl)

	return manager
}
