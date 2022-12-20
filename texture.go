package main

import (
	"image"
)

type Texture struct {
	width, height int
	data          []Color
}

func NewTexture(width, height int, data []Color) Texture {
	return Texture{width, height, data}
}

func NewTextureFromImage(image image.Image) Texture {
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()
	data := make([]Color, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := image.At(x, y).RGBA()
			color := Color{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}
			data[(y*width)+x] = color
		}
	}
	return Texture{width, height, data}
}

// Palette represents the 16-color palette to use during rendering a polygon.  This is due
// to FFT texture storage. The raw texture pixel value is an index for a palettes. Each
// map has 16 palettes of 16 colors each. Each polygon references one of the 16 palettes
// to use.  Eventually Renderer.DrawTexel() function uses uses the pallet.
type Palette []Color
