package main

import (
	"go/build"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

var (
	fontFile   = "luxi-fonts/luxisr.ttf"
	dpi        = 72.0
	txtPortion = 25 // percentage
)

func drawString(c *canvas, s string) error {
	font, err := loadFont()
	if err != nil {
		log.Println(err)
	}

	var sider int
	if c.h > c.w {
		sider = c.w
	} else {
		sider = c.h
	}
	size := float64((sider * txtPortion) / 100)
	scale := size / float64(font.FUnitsPerEm())
	bounds := font.Bounds(int32(scale))
	height := float64(bounds.YMax-bounds.YMin) * scale

	_width := 0
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := font.Index(rune)
		if hasPrev {
			_width += int(font.Kerning(font.FUnitsPerEm(), prev, index))
		}
		_width += int(font.HMetric(font.FUnitsPerEm(), index).AdvanceWidth)
	}
	width := float64(_width) * scale

	log.Printf("w: %v, h: %v", width, height)
	log.Printf("size: %v", size)
	log.Printf("scale: %v", scale)

	ctx := freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(font)
	ctx.SetFontSize(size)
	ctx.SetClip(c.img.Bounds())
	ctx.SetDst(c.img)
	ctx.SetSrc(image.NewUniform(c.fg))

	x, y := (c.w-int(width))/2, (c.h-int(height))/2

	log.Printf("x: %v, y: %v", x, y)
	// Draw the text
	pt := freetype.Pt(x, y)
	_, err = ctx.DrawString(s, pt)
	if err != nil {
		return err
	}
	return nil
}

func loadFont() (font *truetype.Font, err error) {
	if pkg, err := build.Import("github.com/gedex/go-imgplaceholder", "", build.FindOnly); err == nil {
		p := filepath.Join(pkg.Dir, fontFile)
		if _, err := os.Stat(p); err == nil {
			fontFile = p
		}
	}

	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	return
}
