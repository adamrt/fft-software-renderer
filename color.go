package main

import "math/rand"

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

func randColor() Color {
	return Color{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}
