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
	mesh              = Mesh{}
)

type Engine struct {
	isRunning bool
	window    *Window
	renderer  *Renderer
}

func NewEngine(window *Window, renderer *Renderer) *Engine {
	return &Engine{window: window, renderer: renderer}
}

func (e *Engine) loadObj(file string) {
	mesh = NewMeshFromObj(file)
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
	mesh.rotation.x += 0.003
	mesh.rotation.y += 0.005
	mesh.rotation.z += 0.002

	for _, tri := range mesh.triangles {
		var t Triangle
		for i, vertex := range tri.points {
			// Transform
			rotated := rotate(vertex, mesh.rotation)

			// Project
			projected := project(rotated)
			// Invert the Y asis to compensate for the Y axis of the model and
			// the color buffer being different (+Y up vs +Y down, respectively).
			projected.y *= -1

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
		a, b, c := t.points[0], t.points[1], t.points[2]
		e.renderer.DrawLine(int(a.x), int(a.y), int(b.x), int(b.y), White)
		e.renderer.DrawLine(int(b.x), int(b.y), int(c.x), int(c.y), White)
		e.renderer.DrawLine(int(c.x), int(c.y), int(a.x), int(a.y), White)

		e.renderer.DrawRect(int(a.x)-2, int(a.y)-2, 4, 4, Red)
		e.renderer.DrawRect(int(b.x)-2, int(b.y)-2, 4, 4, Red)
		e.renderer.DrawRect(int(c.x)-2, int(c.y)-2, 4, 4, Red)

	}

	// Present
	e.window.Present()
}

func project(v Vec3) Vec3 {
	fov := 128.0 // arbitrary fov to scale the small points
	return Vec3{
		x: (v.x * fov),
		y: (v.y * fov),
		z: v.z,
	}
}

func rotate(v, rotation Vec3) Vec3 {
	rotated := rotate_x(v, rotation.x)
	rotated = rotate_y(rotated, rotation.y)
	rotated = rotate_z(rotated, rotation.z)
	return rotated
}
