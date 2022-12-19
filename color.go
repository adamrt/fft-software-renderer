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

func (orig Color) Mul(factor float64) Color {
	factor = clamp(factor, 0.0, 1.0)
	return Color{
		R: uint8(float64(orig.R) * factor),
		G: uint8(float64(orig.G) * factor),
		B: uint8(float64(orig.B) * factor),
		A: orig.A,
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
