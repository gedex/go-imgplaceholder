package provider

import (
	"bytes"
	"code.google.com/p/freetype-go/freetype"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

const (
	defaultBackgroundColor = color.RGBAModel{}
	defaultTextColor       = color.RGBAModel{255, 255, 255, 255}
)

type BlankParameter struct {
	Width           int
	Height          int
	BackgroundColor color.RGBAModel
	TextColor       color.RGBAModel
	Text            string
	FontFile        string
}

func BlankProvide(p *BlankParameter) ([]byte, error) {
	m := image.NewRGBA(image.Rect(0, 0, p.Width, p.Height))

	draw.Draw(m, m.Bounds(), &image.Uniform{p.BackgroundColor}, image.ZP, draw.Src)

	var buf bytes.Buffer
	err := png.Encode(&buf, m)
	return buf.Bytes(), err
}

// func setFont() error {
// 	fontBytes, err :=
// }
