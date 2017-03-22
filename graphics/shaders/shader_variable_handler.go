package shaders

import (
	gl "github.com/go-gl/gl/v3.3-core/gl"
)

// ShaderVariableHandler keeps track of uniform locations
type ShaderVariableHandler struct {
	uniformLocations   map[string]uint32
	attributeLocations map[string]uint32
}

// ReadUniformLocations fetches the uniform locations from opengl
func (h *ShaderVariableHandler) ReadUniformLocations(shaderPgm uint32, uniforms []string) {
	for _, id := range uniforms {
		h.uniformLocations[id] = getUniformLocation(shaderPgm, id)
	}
}

func getUniformLocation(shaderPgm uint32, id string) uint32 {
	return uint32(gl.GetUniformLocation(shaderPgm, gl.Str(id+"\x00")))
}

// ReadAttributeLocations fetches the uniform locations from opengl
func (h *ShaderVariableHandler) ReadAttributeLocations(shaderPgm uint32, attributes []string) {
	for _, id := range attributes {
		h.attributeLocations[id] = getAttributeLocation(shaderPgm, id)
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

func getAttributeLocation(shaderPgm uint32, id string) uint32 {
	return uint32(gl.GetAttribLocation(shaderPgm, gl.Str(id+"\x00")))
}

// Get returns the location of a shader variable
func (h *ShaderVariableHandler) Get(id string) uint32 {
	if loc, ok := h.uniformLocations[id]; ok {
		return loc
	}
	// else an attribute
	return h.attributeLocations[id]
}

// NewShaderVariableHandler inits and returns
func NewShaderVariableHandler() *ShaderVariableHandler {
	u := new(ShaderVariableHandler)
	u.uniformLocations = make(map[string]uint32)
	u.attributeLocations = make(map[string]uint32)
	return u
}
