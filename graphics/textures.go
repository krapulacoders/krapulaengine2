package graphics

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"image"
	"image/draw"
	_ "image/png"
	"os"
)

func NewTextureFromFile(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return NewTextureFromImage(rgba)

}

func NewTextureFromImage(img *image.RGBA) (uint32, error) {

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(img.Rect.Dx()),
		int32(img.Rect.Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(img.Pix))

	return texture, nil
}
