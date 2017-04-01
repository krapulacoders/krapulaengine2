package shaders

import (
	"fmt"

	gl "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/krapulacoders/krapulaengine2/graphics/errors"
)

// ShaderVariableHandler keeps track of uniform locations
type ShaderVariableHandler struct {
	uniformLocations   map[string]int32
	attributeLocations map[string]int32
}

// ReadUniformLocations fetches the uniform locations from opengl
func (h *ShaderVariableHandler) ReadUniformLocations(shaderPgm uint32, uniforms []string) {
	for _, id := range uniforms {
		if loc := getUniformLocation(shaderPgm, id); loc != -1 {
			h.uniformLocations[id] = loc
		} else {
			errors.LogError(errors.Normal, fmt.Sprintf("couldn't find uniform %v in shaderPgm %v", id, shaderPgm))
		}
	}
}

func getUniformLocation(shaderPgm uint32, id string) int32 {
	result := gl.GetUniformLocation(shaderPgm, gl.Str(id+"\x00"))
	errors.AssertGLError(errors.Normal, fmt.Sprintf("glGetAttribLocation(%v, %v) -> %v", shaderPgm, id, result))
	return result
}

// ReadAttributeLocations fetches the uniform locations from opengl
func (h *ShaderVariableHandler) ReadAttributeLocations(shaderPgm uint32, attributes []string) {
	for _, id := range attributes {
		if loc := getAttributeLocation(shaderPgm, id); loc != -1 {
			h.attributeLocations[id] = loc
		} else {
			errors.LogError(errors.Normal, fmt.Sprintf("couldn't find attribute %v in shaderPgm %v", id, shaderPgm))
		}
	}
}

// RefreshLocations refetches the locations
func (h *ShaderVariableHandler) RefreshLocations(shaderPgm uint32) {
	for id := range h.uniformLocations {
		h.uniformLocations[id] = getUniformLocation(shaderPgm, id)
	}
	for id := range h.attributeLocations {
		h.attributeLocations[id] = getAttributeLocation(shaderPgm, id)
	}
}

func getAttributeLocation(shaderPgm uint32, id string) int32 {
	result := gl.GetAttribLocation(shaderPgm, gl.Str(id+"\x00"))
	errors.AssertGLError(errors.Normal, fmt.Sprintf("glGetAttribLocation(%v, %v) -> %v", shaderPgm, id, result))
	return result
}

// GetAttribute returns the location of a shader variable
func (h *ShaderVariableHandler) GetAttribute(id string) int32 {
	if loc, ok := h.attributeLocations[id]; ok {
		return loc
	}
	errors.LogError(errors.Normal, fmt.Sprintf("uniform not found: %v", id))
	return 0
}

// GetUniform returns the location of a shader variable
func (h *ShaderVariableHandler) GetUniform(id string) int32 {
	if loc, ok := h.uniformLocations[id]; ok {
		return loc
	}
	errors.LogError(errors.Normal, fmt.Sprintf("uniform not found: %v", id))
	return 0
}

// NewShaderVariableHandler inits and returns
func NewShaderVariableHandler() *ShaderVariableHandler {
	u := new(ShaderVariableHandler)
	u.uniformLocations = make(map[string]int32)
	u.attributeLocations = make(map[string]int32)
	return u
}
