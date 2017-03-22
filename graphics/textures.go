package graphics

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var texCache = make(map[string]uint32)

// GetTextureByID returns an opengl texture id identified by a string id.
func GetTextureByID(id string) (uint32, error) {
	if texid, ok := texCache[id]; ok {
		return texid, nil
	}
	return 0, errors.New("texture not found: " + id)
}

// RegisterTextureFromFile first checks the cache for the specified texture id,
// then it reads it from the file.
func RegisterTextureFromFile(id string, file string) (uint32, error) {
	if texid, err := GetTextureByID(id); err == nil {
		return texid, nil
	}
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
	return RegisterTextureFromImage(id, rgba)

}

// RegisterTextureFromImage first checks the cache for the specified texture id,
// then it creates one from the image.
func RegisterTextureFromImage(id string, img *image.RGBA) (uint32, error) {
	if texid, err := GetTextureByID(id); err == nil {
		return texid, nil
	}

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
