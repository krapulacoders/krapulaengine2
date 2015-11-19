package graphics

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/draw"
	_ "image/png"
	"io/ioutil"
)

func GenerateImageFromFont(text string, font *truetype.Font, font_size float64) (*image.RGBA, error) {
	fg, bg := image.Black, image.Transparent
	rgba := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	dpi := 150.0
	//line_spacing := 1.5
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(font_size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	//c.SetHinting(font.HintingNone)
	//c.SetHinting(font.HintingFull)

	// Draw the text.
	pt := freetype.Pt(10, 128+int(c.PointToFixed(font_size)>>6))
	_, err := c.DrawString(text, pt)
	if err != nil {
		return nil, err
	}
	//pt.Y += c.PointToFixed(font_size * line_spacing)
	return rgba, nil
}

func ReadFont(font_file string) (*truetype.Font, error) {
	// Read the font data.
	font_bytes, err := ioutil.ReadFile(font_file)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(font_bytes)
	if err != nil {
		return nil, err
	}

	return font, nil
}
