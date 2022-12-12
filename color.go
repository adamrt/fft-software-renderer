package main

var (
	Transparent = Color{0, 0, 0, 0}

	Black     = Color{0, 0, 0, 255}
	White     = Color{255, 255, 255, 255}
	LightGray = Color{38, 38, 38, 255}
	DarkGray  = Color{36, 36, 36, 255}

	Red    = Color{255, 0, 0, 255}
	Green  = Color{0, 255, 0, 255}
	Blue   = Color{0, 0, 255, 255}
	Yellow = Color{255, 255, 0, 255}
)

type Color struct {
	R, G, B, A uint8
}
