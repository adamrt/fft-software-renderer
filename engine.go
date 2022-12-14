package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS        = 60
	MSPerFrame = (1000 / FPS)
)

// Data
var (
	trianglesToRender = []Triangle{}
	mesh              = Mesh{}

	// Timing
	previous uint32 = 0

	// arbitrary fov to scale the small points
	fov            = 640.0
	cameraPosition = Vec3{0, 0, 0}
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

	// Rotate more each frame
	// mesh.rotation.x += 0.03
	mesh.rotation.y += 0.03
	// mesh.rotation.z += 0.03

	for _, triangle := range mesh.triangles {
		var vertices [3]Vec3

		// Transform
		for i, vertex := range triangle.vertices {
			// Rotate then move away from the camera (0,0,0)
			vertex = rotate(vertex, mesh.rotation)
			vertex.z += 5.0
			vertices[i] = vertex
		}

		// Backface culling
		a, b, c := vertices[0], vertices[1], vertices[2]
		ab, ac := b.Sub(a), c.Sub(a)
		normal := ab.Cross(ac).Normalize()
		ray := cameraPosition.Sub(a)
		visibility := normal.Dot(ray)

		if visibility < 0.0 {
			continue
		}

		// Projection
		for i, vertex := range vertices {
			// Project
			point := project(vertex)

			// Invert the Y asis to compensate for the Y axis of the model and
			// the color buffer being different (+Y up vs +Y down, respectively).
			point.y *= -1

			// Scale and translate to middle of screen
			point.x += float64(e.window.width / 2)
			point.y += float64(e.window.height / 2)

			triangle.points[i] = point
		}
		trianglesToRender = append(trianglesToRender, triangle)
	}
}

func (e *Engine) render() {
	// Draw
	for _, t := range trianglesToRender {
		// Draw triangles
		a, b, c := t.points[0], t.points[1], t.points[2]
		e.renderer.DrawLine(int(a.x), int(a.y), int(b.x), int(b.y), White)
		e.renderer.DrawLine(int(b.x), int(b.y), int(c.x), int(c.y), White)
		e.renderer.DrawLine(int(c.x), int(c.y), int(a.x), int(a.y), White)

		// Draw vertices
		// e.renderer.DrawRect(int(a.x)-2, int(a.y)-2, 4, 4, Red)
		// e.renderer.DrawRect(int(b.x)-2, int(b.y)-2, 4, 4, Red)
		// e.renderer.DrawRect(int(c.x)-2, int(c.y)-2, 4, 4, Red)
	}

	// Present
	e.window.Present()

	// Clear triangles from last frame
	trianglesToRender = trianglesToRender[:0]
}

func project(v Vec3) Vec2 {
	return Vec2{
		x: (v.x * fov) / v.z,
		y: (v.y * fov) / v.z,
	}
}

func rotate(v, rotation Vec3) Vec3 {
	rotated := rotate_x(v, rotation.x)
	rotated = rotate_y(rotated, rotation.y)
	rotated = rotate_z(rotated, rotation.z)
	return rotated
}
