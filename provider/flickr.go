package provider

import (
	"image/color"
)

type FlickrParameter struct {
	Width  int
	Height int
	Tag    string
}

func Provide(p *FlickrParameter) ([]byte, err) {

}
