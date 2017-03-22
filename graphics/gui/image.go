package gui

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Image struct {
	StaticComponentImpl
	texture uint32
}

func (img *Image) Render() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, img.texture)

	//gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
}
