package main

import (
	"math"
)

type DirectionalLight struct {
	position Vec3
	target   Vec3
	color    Color
}

func NewDirectionLight(position, target Vec3) DirectionalLight {
	return DirectionalLight{position: position, target: position}
}

func (l *DirectionalLight) direction() Vec3 {
	return l.target.Sub(l.position)
}

type AmbientLight struct {
	color Color
}

func applyLightIntensity(orig Color, factor float64) Color {
	// Clamp from 0.0 to 1.0
	factor = math.Max(0, math.Min(factor, 1.0))

	return Color{
		R: uint8(float64(orig.R) * factor),
		G: uint8(float64(orig.G) * factor),
		B: uint8(float64(orig.B) * factor),
		A: orig.A,
	}
}
