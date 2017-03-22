package graphics

import (
	"strconv"
	"time"
	"unsafe"

	gl "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/krapulacoders/krapulaengine2/graphics"
	"github.com/krapulacoders/krapulaengine2/graphics/shaders"
)

const (
	// RotationEnabled attribute
	RotationEnabled = iota
	// ColorEnabled attribute
	ColorEnabled = iota
	// TexturesEnabled attribute
	TexturesEnabled = iota

	maxArrayElements = 1024 * 1024 * 1024
)

// GenericObject2D struct containing rendering details
type GenericObject2D struct {
	Coords        []mgl32.Vec2
	Color         mgl32.Vec4
	Rotation      float32
	centerPoint   mgl32.Vec2
	TextureCoords []mgl32.Vec2
}

// BasicRenderGroup2D is a 2D render group that supports colors, textures and rotation.
type BasicRenderGroup2D struct {
	rg              *graphics.RenderGroup
	shaderVars      *shaders.ShaderVariableHandler
	renderType      uint32
	coordsPerObject int
	objects         []*GenericObject2D
	freeIndexes     []int
	vaos            []uint32
	vbo             int
	attributes      map[int]bool
	rendering       bool
	texture         uint32
	hasChanged      bool
}

// NewBasicRenderGroup2D creates a new basic 2d render group
func NewBasicRenderGroup2D(id string, glType uint32, expectedSize int32, changeRatio float32, texture uint32) *graphics.RenderGroup {
	manager := new(BasicRenderGroup2D)
	manager.objects = make([]*GenericObject2D, 0, expectedSize)
	manager.freeIndexes = make([]int, 0, int(float32(expectedSize)*changeRatio))
	manager.shaderVars = shaders.NewShaderVariableHandler()
	manager.texture = texture
	manager.renderType = glType

	switch glType {
	case gl.TRIANGLES:
		manager.coordsPerObject = 3
	case gl.LINES:
		manager.coordsPerObject = 2
	case gl.POINTS:
		manager.coordsPerObject = 1
	case gl.TRIANGLE_FAN, gl.TRIANGLE_STRIP:
		// special case: variable amount of coords for one object.
		manager.coordsPerObject = 0
		fallthrough
	default:
		panic("unsupported glType: " + strconv.Itoa((int)(glType)))

	}
	manager.SetAttribute(ColorEnabled, true)
	manager.SetAttribute(RotationEnabled, true)
	manager.SetAttribute(TexturesEnabled, texture != 0)

	g := graphics.NewRenderGroup(id, manager)
	g.SetShaderFile("basic.vert")
	g.SetShaderFile("basic.frag")

	return g
}

// SetAttribute sets an attribute, see *Enabled constants
func (g *BasicRenderGroup2D) SetAttribute(attribute int, value bool) {
	g.attributes[attribute] = value
}

// AddObject adds an object and returns an id for it.
func (g *BasicRenderGroup2D) AddObject(obj *GenericObject2D) int {
	for g.rendering {
		time.Sleep(100 * time.Microsecond)
	}
	if len(g.freeIndexes) > 0 {
		// remove last element
		freeIndex := g.freeIndexes[len(g.freeIndexes)]
		g.freeIndexes = g.freeIndexes[:len(g.freeIndexes)-1]
		g.objects[freeIndex] = obj
		return freeIndex
	}
	// else
	g.objects = append(g.objects, obj)
	g.NotifyObjectChanged()
	return len(g.objects) - 1
}

// RemoveObject removes an object by id
func (g *BasicRenderGroup2D) RemoveObject(id int) {
	for g.rendering {
		time.Sleep(100 * time.Microsecond)
	}
	if len(g.freeIndexes) > id && g.objects[id] != nil {
		// Add to freeIndexes
		g.freeIndexes = append(g.freeIndexes, id)
		g.objects[id] = nil
	} else {
		panic("Tried removing non-existing object")
	}
	g.NotifyObjectChanged()
}

// NotifyObjectChanged tells the rendergroup that an object has changed
func (g *BasicRenderGroup2D) NotifyObjectChanged() {
	g.hasChanged = true
}

// InitShader is run once per program
func (g *BasicRenderGroup2D) InitShader() {
	gl.BindFragDataLocation(g.rg.GetShaderProgram(), 0, gl.Str("outputColor\x00"))
	g.shaderVars.ReadUniformLocations(g.rg.GetShaderProgram(), []string{"normalMatrix"})
	g.shaderVars.ReadAttributeLocations(g.rg.GetShaderProgram(), []string{"vert", "vertTexCoord", "rotation", "centerPoint"})
}

// Render implements the rendering
func (g *BasicRenderGroup2D) Render() {
	g.rendering = true
	if g.hasChanged {
		g.setupRendering()
		g.hasChanged = false
	}
	if g.attributes[TexturesEnabled] {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, g.texture)
	}

	// just render everything
	gl.DrawArrays(g.renderType, 0, (int32)(len(g.objects)*g.coordsPerObject))

	g.rendering = false
}

func (g *BasicRenderGroup2D) setupRendering() {
	vertexSize := 0
	colorSize := 0
	textureCoordSize := 0
	rotationSize := 0
	centerPointSize := 0

	totalCoords := len(g.objects) * g.coordsPerObject

	colorEnabled := g.attributes[ColorEnabled]
	rotationEnabled := g.attributes[RotationEnabled]
	texturesEnabled := g.attributes[TexturesEnabled]

	vertA := g.shaderVars.Get("vert")
	InColorA := g.shaderVars.Get("inColor")
	rotationA := g.shaderVars.Get("rotation")
	centerPointA := g.shaderVars.Get("centerPoint")
	vertTexCoordA := g.shaderVars.Get("vertTexCoord")

	if g.coordsPerObject == 0 {
		panic("not implemented")
	} else {
		// one coord: 2 float32: 2*4.
		vertexSize = totalCoords * 2 * 4
		gl.EnableVertexAttribArray(vertA)
	}

	if colorEnabled {
		// RGBA: 4* float32 = 4*4
		colorSize = totalCoords * 4 * 4
		gl.EnableVertexAttribArray(InColorA)
	} else {
		gl.DisableVertexAttribArray(InColorA)
		// default color is white
		gl.VertexAttrib4f(InColorA, 1, 1, 1, 1)
	}

	if rotationEnabled {
		// a single float32
		rotationSize = totalCoords * 4
		gl.EnableVertexAttribArray(rotationA)
		centerPointSize = totalCoords * 2
		gl.EnableVertexAttribArray(centerPointA)
	} else {
		gl.DisableVertexAttribArray(rotationA)
		gl.DisableVertexAttribArray(centerPointA)
		// default rotation is 0
		gl.VertexAttrib1f(rotationA, 0)
		gl.VertexAttrib2f(centerPointA, 0, 0)
	}

	if texturesEnabled {
		// 1 tex coord per vertex coord, each having a vec2 = 2*float32 = 2*4
		textureCoordSize = totalCoords * 2 * 4
		gl.EnableVertexAttribArray(vertTexCoordA)
	} else {
		gl.DisableVertexAttribArray(vertTexCoordA)
		// default rotation is 0
		gl.VertexAttrib2f(vertTexCoordA, 0, 0)
	}

	vertexIndex := 0
	colorIndex := vertexIndex + vertexSize
	rotationIndex := colorIndex + colorSize
	centerPointIndex := rotationIndex + rotationSize
	textureCoordIndex := centerPointIndex + centerPointSize

	totalSize := vertexSize + colorSize + textureCoordSize + rotationSize + centerPointSize

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, totalSize, nil, gl.DYNAMIC_DRAW)

	bufferPointer := gl.MapBuffer(vbo, gl.WRITE_ONLY)
	// convert the buffer pointer to go arrays
	vertexArray := *((*[maxArrayElements]mgl32.Vec2)(bufferPointer))
	colorPointer := unsafe.Pointer(((uintptr)(bufferPointer)) + (uintptr)(colorIndex))
	colorArray := *((*[maxArrayElements]mgl32.Vec4)(colorPointer))
	rotationPointer := unsafe.Pointer(((uintptr)(bufferPointer)) + (uintptr)(rotationIndex))
	rotationArray := *((*[maxArrayElements]float32)(rotationPointer))
	centerPointPointer := unsafe.Pointer(((uintptr)(bufferPointer)) + (uintptr)(textureCoordIndex))
	centerPointArray := *((*[maxArrayElements]mgl32.Vec2)(centerPointPointer))
	textureCoordPointer := unsafe.Pointer(((uintptr)(bufferPointer)) + (uintptr)(textureCoordIndex))
	textureCoordArray := *((*[maxArrayElements]mgl32.Vec2)(textureCoordPointer))

	for oi, obj := range g.objects {
		for ci := 0; ci < g.coordsPerObject; ci++ {
			arrayIndex := oi*g.coordsPerObject + ci
			vertexArray[arrayIndex] = obj.Coords[ci]

			if texturesEnabled {
				textureCoordArray[arrayIndex] = obj.TextureCoords[ci]
			}
			if rotationEnabled {
				rotationArray[arrayIndex] = obj.Rotation
				centerPointArray[arrayIndex] = obj.centerPoint
			}
			if colorEnabled {
				colorArray[arrayIndex] = obj.Color
			}
		}
	}
	gl.UnmapBuffer(gl.ARRAY_BUFFER)

	gl.VertexAttribPointer(vertA, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	if colorEnabled {
		gl.VertexAttribPointer(InColorA, 4, gl.FLOAT, false, 0, gl.PtrOffset(colorIndex))
	}
	if texturesEnabled {
		gl.VertexAttribPointer(vertTexCoordA, 2, gl.FLOAT, false, 0, gl.PtrOffset(textureCoordIndex))
	}
	if rotationEnabled {
		gl.VertexAttribPointer(rotationA, 1, gl.FLOAT, false, 0, gl.PtrOffset(rotationIndex))
		gl.VertexAttribPointer(centerPointA, 2, gl.FLOAT, false, 0, gl.PtrOffset(centerPointIndex))
	}

}

// Deinit does nothing
func (g *BasicRenderGroup2D) Deinit() {

}
