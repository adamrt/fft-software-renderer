package main

type Texture struct {
	width, height int
	data          []Color
}

func NewTexture(width, height int, data []Color) Texture {
	return Texture{width, height, data}
}

// Palette represents the 16-color palette to use during rendering a polygon.  This is due
// to FFT texture storage. The raw texture pixel value is an index for a palettes. Each
// map has 16 palettes of 16 colors each. Each polygon references one of the 16 palettes
// to use.  Eventually Renderer.DrawTexel() function uses uses the pallet.
type Palette []Color
