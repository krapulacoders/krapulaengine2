package rendergroups

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"

	gl "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/krapulacoders/krapulaengine2/graphics"
	"github.com/krapulacoders/krapulaengine2/graphics/errors"
	"github.com/krapulacoders/krapulaengine2/graphics/shaders"
)

// BasicRenderGroup2DAttribute represents different kinds of functionality
type BasicRenderGroup2DAttribute int

const (
	// RotationEnabled attribute
	RotationEnabled = iota
	// ColorEnabled attribute
	ColorEnabled = iota
	// TexturesEnabled attribute
	TexturesEnabled = iota
)
const (
	maxArrayElements = 1 << 30
)

// GenericObject2D struct containing rendering details
type GenericObject2D struct {
	Coords        []mgl32.Vec3
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
	vao             uint32
	vbo             uint32
	attributes      map[BasicRenderGroup2DAttribute]bool
	rendering       bool
	texture         uint32
	hasChanged      bool
}

// NewBasicRenderGroup2D creates a new basic 2d render group
func NewBasicRenderGroup2D(id string, glType uint32, expectedSize int32,
	changeRatio float32, texture uint32) (*graphics.RenderGroup, *BasicRenderGroup2D) {
	manager := new(BasicRenderGroup2D)
	manager.objects = make([]*GenericObject2D, 0, expectedSize)
	manager.freeIndexes = make([]int, 0, int(float32(expectedSize)*changeRatio))
	manager.shaderVars = shaders.NewShaderVariableHandler()
	manager.texture = texture
	manager.renderType = glType
	manager.attributes = make(map[BasicRenderGroup2DAttribute]bool)

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
	g.SetShaderFile("graphics/shaders/2d/basic.vert")
	g.SetShaderFile("graphics/shaders/2d/basic.frag")
	manager.rg = g

	return g, manager
}

// SetAttribute sets an attribute, see *Enabled constants
func (g *BasicRenderGroup2D) SetAttribute(attribute BasicRenderGroup2DAttribute, value bool) {
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
	errors.AssertGLError(errors.Normal, "glBindFragDataLocation")

	g.shaderVars.ReadUniformLocations(g.rg.GetShaderProgram(),
		[]string{"normalMatrix", "tex"})
	errors.AssertGLError(errors.Normal, "after read uniforms")

	g.shaderVars.ReadAttributeLocations(g.rg.GetShaderProgram(),
		[]string{"vert", "vertTexCoord", "rotation", "centerPoint", "inColor"})
	errors.AssertGLError(errors.Normal, "after read attributes")

	gl.GenVertexArrays(1, &g.vao)
	errors.AssertGLError(errors.Critical, "glGenVertexArrays")
	gl.GenBuffers(1, &g.vbo)
	errors.AssertGLError(errors.Critical, "glGenBuffers")

	//gl.Uniform1i((int32)(g.shaderVars.Get("tex")), 0)
	//errors.AssertGLError(errors.Normal, fmt.Sprintf("Uniform1i(%v, 0)", int32(g.shaderVars.Get("tex"))))
}

// Render implements the rendering
func (g *BasicRenderGroup2D) Render() {
	errors.AssertGLError(errors.Debug, "BasicRenderGroup2D.Render")
	g.rendering = true
	if g.hasChanged {
		g.setupRendering()
		g.hasChanged = false
	}
	gl.BindVertexArray(g.vao)
	if g.attributes[TexturesEnabled] {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, g.texture)
		errors.AssertGLError(errors.Normal, "glBindTexture")
	}
	gl.LineWidth(5)

	// uniforms
	normalMatrixUniform := g.shaderVars.GetUniform("normalMatrix")
	textureUniform := g.shaderVars.GetUniform("tex")
	normalMatrix := graphics.GetSimpleNormalMatrixMat2()
	fmt.Printf("normalMatrix: %v", normalMatrix)
	gl.UniformMatrix2fv(normalMatrixUniform, 1, false, &normalMatrix[0])
	errors.AssertGLError(errors.Debug, fmt.Sprintf("normalMatrix: %v", textureUniform))

	//gl.Uniform1i(textureUniform, 0)
	//errors.AssertGLError(errors.Debug, fmt.Sprintf("tex (sampler2D). id: %v", textureUniform))

	errors.AssertGLError(errors.Debug, "before glDrawArrays")
	// just render everything
	gl.DrawArrays(g.renderType, 0, (int32)(len(g.objects)*g.coordsPerObject))
	errors.AssertGLError(errors.Normal, "glDrawArrays")

	//fmt.Printf("drew %v vertices\n", len(g.objects)*g.coordsPerObject)
	g.rendering = false
}

func (g *BasicRenderGroup2D) setupRendering() {
	errors.AssertGLError(errors.Normal, "BasicRenderGroup2D.setupRendering")
	vertexSize := 0
	colorSize := 0
	textureCoordSize := 0
	rotationSize := 0
	centerPointSize := 0
	if g.coordsPerObject == 0 {
		panic("not implemented")
	}
	totalCoords := len(g.objects) * g.coordsPerObject

	colorEnabled := g.attributes[ColorEnabled]
	rotationEnabled := g.attributes[RotationEnabled]
	texturesEnabled := g.attributes[TexturesEnabled]

	vertA := uint32(g.shaderVars.GetAttribute("vert"))
	InColorA := uint32(g.shaderVars.GetAttribute("inColor"))
	rotationA := uint32(g.shaderVars.GetAttribute("rotation"))
	centerPointA := uint32(g.shaderVars.GetAttribute("centerPoint"))
	vertTexCoordA := uint32(g.shaderVars.GetAttribute("vertTexCoord"))

	gl.BindVertexArray(g.vao)

	// one coord: 3 float32: 2*4.
	vertexSize = totalCoords * 3 * 4

	if colorEnabled {
		// RGBA: 4* float32 = 4*4
		colorSize = totalCoords * 4 * 4
	}

	if rotationEnabled {
		// a single float32
		rotationSize = totalCoords * 4
		// two float32 per coord
		centerPointSize = totalCoords * 2 * 4
	}

	if texturesEnabled {
		// 1 tex coord per vertex coord, each having a vec2 = 2*float32 = 2*4
		textureCoordSize = totalCoords * 2 * 4
	}

	vertexIndex := 0
	colorIndex := vertexIndex + vertexSize
	rotationIndex := colorIndex + colorSize
	centerPointIndex := rotationIndex + rotationSize
	textureCoordIndex := centerPointIndex + centerPointSize

	totalSize := vertexSize + colorSize + rotationSize + centerPointSize + textureCoordSize

	//fmt.Printf("==sizes== for %v objects with %v coords\n vertex: %v ___ color: %v ___ rotation:%v ___ center:%v ___ texcoord:%v\n",
	//	len(g.objects), g.coordsPerObject, vertexSize, colorSize, rotationSize, centerPointSize, textureCoordSize)
	gl.BindBuffer(gl.ARRAY_BUFFER, g.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, totalSize, nil, gl.DYNAMIC_DRAW)
	bufferPointer := gl.MapBufferRange(gl.ARRAY_BUFFER, 0, totalSize, // gl.MAP_READ_BIT|
		gl.MAP_WRITE_BIT|gl.MAP_INVALIDATE_BUFFER_BIT)
	fmt.Printf("bufferPointer: %v. shaderPgm: %v\n", bufferPointer, g.rg.GetShaderProgram())
	if bufferPointer == nil {
		panic("mapping buffer failed")
	}
	// convert the buffer pointer to go arrays
	vertexArray := (*((*[maxArrayElements]mgl32.Vec3)(bufferPointer)))[:totalCoords:totalCoords]
	colorPointer := unsafe.Pointer(((uintptr)(bufferPointer)) + (uintptr)(colorIndex))
	colorArray := (*((*[maxArrayElements]mgl32.Vec4)(colorPointer)))[:totalCoords:totalCoords]
	rotationPointer := unsafe.Pointer(((uintptr)(bufferPointer)) + (uintptr)(rotationIndex))
	rotationArray := (*((*[maxArrayElements]float32)(rotationPointer)))[:totalCoords:totalCoords]
	centerPointPointer := unsafe.Pointer(((uintptr)(bufferPointer)) + (uintptr)(centerPointIndex))
	centerPointArray := (*((*[maxArrayElements]mgl32.Vec2)(centerPointPointer)))[:totalCoords:totalCoords]
	textureCoordPointer := unsafe.Pointer(((uintptr)(bufferPointer)) + (uintptr)(textureCoordIndex))
	textureCoordArray := (*((*[maxArrayElements]mgl32.Vec2)(textureCoordPointer)))[:totalCoords:totalCoords]
	/*
		fmt.Printf("==index==\nvertex: %v \ncolor: %v\n rotation: %v\ncenter: %v\n tex: %v\n",
			(int)((uintptr)(bufferPointer))-(int)((uintptr)(bufferPointer)),
			(int)((uintptr)(colorPointer))-(int)((uintptr)(bufferPointer)),
			(int)((uintptr)(rotationPointer))-(int)((uintptr)(bufferPointer)),
			(int)((uintptr)(centerPointPointer))-(int)((uintptr)(bufferPointer)),
			(int)((uintptr)(textureCoordPointer))-(int)((uintptr)(bufferPointer)),
		)
	*/
	for oi, obj := range g.objects {
		for ci := 0; ci < g.coordsPerObject; ci++ {
			arrayIndex := oi*g.coordsPerObject + ci
			vertexArray[arrayIndex] = obj.Coords[ci]

			if texturesEnabled {
				//fmt.Printf("texture coords set\n")
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
		//fmt.Printf("object: %v\n", obj)
	}
	//fmt.Printf("vertex: %v ___ color: %v ___ rotation:%v ___ center:%v ___ texcoord:%v\n", vertexArray, colorArray, rotationArray, centerPointArray, textureCoordArray)
	//fmt.Printf("raw %v: %v\n", totalSize, (*((*[maxArrayElements]byte)(bufferPointer)))[:totalSize:totalSize])
	//fmt.Printf("raw %v: %v\n", totalSize, (*((*[maxArrayElements]float32)(bufferPointer)))[:totalSize/4:totalSize/4])
	gl.UnmapBuffer(gl.ARRAY_BUFFER)
	errors.AssertGLError(errors.Critical, "glUnmapBuffer")

	errors.AssertGLError(errors.Normal, fmt.Sprintf("glBindBuffer(gl.ARRAY_BUFFER, %v)", g.vbo))
	gl.VertexAttribPointer(vertA, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(vertA)
	errors.AssertGLError(errors.Normal, "vertex attribute vertex")

	if colorEnabled {
		gl.VertexAttribPointer(InColorA, 4, gl.FLOAT, false, 0, gl.PtrOffset(colorIndex))
		gl.EnableVertexAttribArray(InColorA)
	} else {
		gl.DisableVertexAttribArray(InColorA)
		// default color is white
		gl.VertexAttrib4f(InColorA, 1, 1, 1, 1)
	}
	errors.AssertGLError(errors.Normal, "vertex attribute colors")

	if vertTexCoordA > 0 {
		if texturesEnabled {
			gl.VertexAttribPointer(vertTexCoordA, 2, gl.FLOAT, false, 0, gl.PtrOffset(textureCoordIndex))
			gl.EnableVertexAttribArray(vertTexCoordA)
		} else {
			gl.DisableVertexAttribArray(vertTexCoordA)
			// default rotation is 0
			gl.VertexAttrib2f(vertTexCoordA, 0, 0)
		}
	}

	errors.AssertGLError(errors.Normal, "vertex attribute textures")

	if rotationEnabled {
		gl.VertexAttribPointer(rotationA, 1, gl.FLOAT, false, 0, gl.PtrOffset(rotationIndex))
		gl.VertexAttribPointer(centerPointA, 2, gl.FLOAT, false, 0, gl.PtrOffset(centerPointIndex))
		gl.EnableVertexAttribArray(rotationA)
		gl.EnableVertexAttribArray(centerPointA)
		errors.AssertGLError(errors.Normal, "vertex attribute rotation")
	} else {
		gl.DisableVertexAttribArray(rotationA)
		gl.DisableVertexAttribArray(centerPointA)
		// default rotation is 0
		gl.VertexAttrib1f(rotationA, 0)
		gl.VertexAttrib2f(centerPointA, 0, 0)
	}

	errors.AssertGLError(errors.Normal, "vertex attribute rotation")
}

// Deinit does nothing
func (g *BasicRenderGroup2D) Deinit() {

}
