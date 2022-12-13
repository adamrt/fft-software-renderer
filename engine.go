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
	trianglesToRender = []Triangle{}
	mesh              = cubeMesh()
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
	// Variable timestep
	if wait := MSPerFrame - (sdl.GetTicks() - previous); wait > 0 && wait <= MSPerFrame {
		sdl.Delay(wait)
	}
	previous = sdl.GetTicks()

	// Clear triangles from last frame
	trianglesToRender = trianglesToRender[:0]

	// Rotate more each frame
	mesh.rotation.x += 0.025
	mesh.rotation.y += 0.005
	mesh.rotation.z += 0.015

	for _, face := range mesh.faces {
		var t Triangle
		for i, idx := range face.indexes {
			vertex := mesh.vertices[idx-1]

			// Transform
			rotated := rotate(vertex, mesh.rotation)

			// Project
			projected := project(rotated)

			// Scale and translate to middle of screen
			projected.x += float64(e.window.width / 2)
			projected.y += float64(e.window.height / 2)

			t.points[i] = projected
		}
		trianglesToRender = append(trianglesToRender, t)
	}
}

func (e *Engine) render() {
	// Draw
	for _, t := range trianglesToRender {
		e.renderer.DrawRect(int(t.points[0].x), int(t.points[0].y), 4, 4, Yellow)
		e.renderer.DrawRect(int(t.points[1].x), int(t.points[1].y), 4, 4, Yellow)
		e.renderer.DrawRect(int(t.points[2].x), int(t.points[2].y), 4, 4, Yellow)

	}

	// Present
	e.window.Present()
}

func project(v Vec3) Vec2 {
	fov := 128.0 // arbitrary fov to scale the small points
	return Vec2{
		x: (v.x * fov),
		y: (v.y * fov),
	}
}

func rotate(v, rotation Vec3) Vec3 {
	rotated := rotate_x(v, rotation.x)
	rotated = rotate_y(rotated, rotation.y)
	rotated = rotate_z(rotated, rotation.z)
	return rotated
}
