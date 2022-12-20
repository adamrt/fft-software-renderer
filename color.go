package main

import (
	"math/rand"
)

var (
	Transparent = Color{0, 0, 0, 0}

	Black     = Color{0, 0, 0, 255}
	White     = Color{255, 255, 255, 255}
	LightGray = Color{38, 38, 38, 255}
	DarkGray  = Color{36, 36, 36, 255}

	Red   = Color{255, 0, 0, 255}
	Green = Color{0, 255, 0, 255}
	Blue  = Color{0, 0, 255, 255}

	Yellow  = Color{255, 255, 0, 255}
	Magenta = Color{255, 0, 255, 255}
	Cyan    = Color{0, 255, 255, 255}
)

type Color struct {
	R, G, B, A uint8
}

func (c Color) isTrans() bool {
	return c.A == 0 && c.R == 0 && c.G == 0 && c.B == 0
}

func (c Color) Mul(factor float64) Color {
	factor = clamp(factor, 0.0, 1.0)
	return Color{
		R: uint8(float64(c.R) * factor),
		G: uint8(float64(c.G) * factor),
		B: uint8(float64(c.B) * factor),
		A: c.A,
	}
}

func randColor() Color {
	return Color{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

// Gradient background color for FFT Map
type Background struct {
	Top    Color
	Bottom Color
}

// These are vertical gradients so we don't care about x.  The colors need to be
// float64s before subtraction so there isn't uint8 overflow.
func (bg Background) At(y int, height int) Color {
	d := float64(y) / float64(height)
	r := float64(bg.Bottom.R) + d*(float64(bg.Top.R)-float64(bg.Bottom.R))
	g := float64(bg.Bottom.G) + d*(float64(bg.Top.G)-float64(bg.Bottom.G))
	b := float64(bg.Bottom.B) + d*(float64(bg.Top.B)-float64(bg.Bottom.B))
	return Color{uint8(r), uint8(g), uint8(b), 255}
}

// This returns a checkerboard for background use. It will be copied into a sdl.Texture.
func GenerateCheckerboard(width, height int, a, b Color) []Color {
	// Draw checkerboard to buffer
	bgBuffer := make([]Color, width*height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if (y%(64) < 32) == (x%(64) < 32) {
				bgBuffer[y*width+x] = a
			} else {
				bgBuffer[y*width+x] = b
			}
		}
	}
	return bgBuffer
}
