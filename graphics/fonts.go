package graphics

import (
	"image"
	"image/draw"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// GenerateImageFromFont generates an image from string, font and size.
func GenerateImageFromFont(text string, font *truetype.Font, fontSize float64) (*image.RGBA, error) {
	fg, bg := image.Black, image.Transparent
	rgba := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	dpi := 150.0
	//line_spacing := 1.5
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	//c.SetHinting(font.HintingNone)
	//c.SetHinting(font.HintingFull)

	// Draw the text.
	pt := freetype.Pt(10, 128+int(c.PointToFixed(fontSize)>>6))
	_, err := c.DrawString(text, pt)
	if err != nil {
		return nil, err
	}
	//pt.Y += c.PointToFixed(font_size * line_spacing)
	return rgba, nil
}

// ReadFont reads a file and returns a truetype Font object
func ReadFont(fontFile string) (*truetype.Font, error) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	return font, nil
}
