package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS        = 60
	MSPerFrame = (1000 / FPS)
)

// Timing
var previous uint32

// Data
var (
	projectedPoints []Vec2
	cube            = shapeCube()
	cubeRotation    = Vec3{0, 0, 0}
)

type Engine struct {
	isRunning bool
	window    *Window
	renderer  *Renderer
}

func NewEngine(window *Window, renderer *Renderer) *Engine {
	return &Engine{window: window, renderer: renderer}
}

func (e *Engine) setup() {
	e.isRunning = true
	previous = sdl.GetTicks()
}

func (e *Engine) processInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			e.isRunning = false
		case *sdl.KeyboardEvent:
			if t.Type != sdl.KEYDOWN {
				continue
			}
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE:
				e.isRunning = false
			}
		}
	}
}

func (e *Engine) update() {
	if wait := MSPerFrame - (sdl.GetTicks() - previous); wait > 0 && wait <= MSPerFrame {
		sdl.Delay(wait)
	}
	previous = sdl.GetTicks()

	cubeRotation.x += 0.001
	cubeRotation.y += 0.005
	cubeRotation.z += 0.002
	for _, point := range cube {
		transformedPoint := rotate_x(point, cubeRotation.x)
		transformedPoint = rotate_y(transformedPoint, cubeRotation.y)
		transformedPoint = rotate_z(transformedPoint, cubeRotation.z)

		projectedPoint := project(transformedPoint)
		projectedPoints = append(projectedPoints, projectedPoint)
	}
}

func (e *Engine) render() {
	// Draw
	hw := e.window.width / 2
	hh := e.window.height / 2
	for _, point := range projectedPoints {
		e.renderer.DrawRect(
			int(point.x)+hw,
			int(point.y)+hh,
			4,
			4,
			Yellow,
		)
	}

	// Present
	projectedPoints = projectedPoints[:0]
	e.window.Present()
}

func project(v Vec3) Vec2 {
	fov := 128.0 // arbitrary fov to scale the small points
	return Vec2{
		x: (v.x * fov),
		y: (v.y * fov),
	}
}
